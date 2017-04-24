package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

const headerTemplate = `#ifndef RESPP_H
#define RESPP_H

#include <string>

namespace ResPP
{
std::string R(const char* prefix, const char* path);
}

#endif // RESPP_H
`

const sourceTemplate = `#include "ResPP.h"
#include <map>
#include <string.h>
{{$len := len .Contents}}
namespace
{
	typedef std::map<std::string, const char*> fileToContent;
	{{range .Contents}}
	fileToContent {{.Prefix | prefixToMapName}} = {
		{{- range .Files}}
		{"{{.Name}}", "{{.EncodedContent}}"},
		{{- end}}
	};
	{{- end}}

	std::map<std::string, const fileToContent* const> prefixToIdx = {
		{{- range .Contents}}
		{"{{.Prefix}}", &{{.Prefix | prefixToMapName}}},
		{{- end}}
	};

	const char alphabet[] = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
                                "abcdefghijklmnopqrstuvwxyz"
                                "0123456789+/=";

        long indexIn(const char *a, char c)
	{
	    const char *p = strchr(a, c);

	    if(p)
	    {
		return p - a;
	    }

	    return 0;
	}

	char *b64_decode(const char *s, const unsigned long numBytes, unsigned long *decodedNumBytes)
	{
	    if(!numBytes%4) return 0;

	    // This can only be 0, 1, or 2. Otherwise, something is wrong ...
	    size_t numPaddingChars = 0;

	    if(numBytes>=2)
	    {
		if(s[numBytes-1]=='=') ++numPaddingChars;
		if(s[numBytes-2]=='=') ++numPaddingChars;
	    }

	    const size_t out_len = 3*numBytes/4-numPaddingChars;

	    char *out = (char*)malloc(out_len+1);
	    out[out_len] = '\0';

	    int j = 0;
	    long b[4];
	    size_t i;

	    for(i=0; i<numBytes; i+=4)
	    {
		b[0] = indexIn(alphabet, s[i]);
		b[1] = indexIn(alphabet, s[i+1]);
		b[2] = indexIn(alphabet, s[i+2]);
		b[3] = indexIn(alphabet, s[i+3]);

		out[j++] = (char) ((b[0] << 2) | (b[1] >> 4));

		if(b[2] < 64)
		{
		    out[j++] = (char) ((b[1] << 4) | (b[2] >> 2));

		    if(b[3] < 64)
		    {
			out[j++] = (char) ((b[2] << 6) | b[3]);
		    }
		}
	    }

	    if(decodedNumBytes)
	    {
		*decodedNumBytes = out_len;
	    }

	    return out;
	}
}

namespace ResPP
{
std::string R(const char* prefix, const char* path)
{
	auto prefix_it = prefixToIdx.find(prefix);

	if(prefix_it!=prefixToIdx.cend())
	{
		auto m = prefix_it->second;

		auto path_it = m->find(path);

		if(path_it!=m->cend())
		{
			auto encodedContent = path_it->second;

			if(encodedContent)
			{
				std::string s(encodedContent);
				unsigned long decodedNumBytes;

				auto dec = b64_decode(s.c_str(), s.size(), &decodedNumBytes);
				auto decStr = std::string(dec);

				free(dec);

				return decStr;
			}

			return std::string(encodedContent);
		}
	}

	return std::string();
}
}
`

func generateCpp(config *configuration, outDir string) error {
	funcs := template.FuncMap{
		"addOne": func(i int) int {
			return i + 1
		},
		"prefixToMapName": func(prefix string) string {
			return strings.Replace(prefix, "/", "_", -1)
		},
	}

	tmpl, err := template.New("cpp").Funcs(funcs).Parse(sourceTemplate)

	if err != nil {
		return err
	}

	if _, err := os.Stat("cpp"); os.IsNotExist(err) {
		os.MkdirAll(outDir, 0755)
	}

	if err != nil {
		return err
	}

	ioutil.WriteFile(filepath.Join("cpp", "ResPP.h"), []byte(headerTemplate), 0644)

	f, err := os.Create(filepath.Join("cpp", "/ResPP.cpp"))

	if err != nil {
		return err
	}

	return tmpl.Execute(f, config)
}
