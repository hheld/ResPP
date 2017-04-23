package main

import (
	"io/ioutil"
	"os"
	"strings"
	"text/template"
)

const headerTemplate = `#ifndef RESPPGEN_H
#define RESPP_H

namespace ResPP
{
const char* R(const char* prefix, const char* path);
}

#endif // RESPP_H
`

const sourceTemplate = `#include "ResPP.h"
#include <map>
#include <string>
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
}

namespace ResPP
{
const char* R(const char* prefix, const char* path)
{
	auto prefix_it = prefixToIdx.find(prefix);

	if(prefix_it!=prefixToIdx.cend())
	{
		auto m = prefix_it->second;

		auto path_it = m->find(path);

		if(path_it!=m->cend())
		{
			auto encodedContent = path_it->second;
			return encodedContent;
		}
	}

	return nullptr;
}
}
`

func generateCpp(config *configuration) error {
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
		os.Mkdir("cpp", 0755)
	}

	if err != nil {
		return err
	}

	ioutil.WriteFile("cpp/ResPP.h", []byte(headerTemplate), 0644)

	f, err := os.Create("cpp/ResPP.cpp")

	if err != nil {
		return err
	}

	return tmpl.Execute(f, config)
}
