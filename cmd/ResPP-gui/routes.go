package main

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/hheld/ResPP/pkg/configuration"
)

func (c *config) load(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	fileName := r.URL.Query().Get("file")

	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	err := (*configuration.Configuration)(c).LoadFromFile(fileName)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

type fileInfo struct {
	Prefix string   `json:"prefix"`
	Files  []string `json:"files"`
}

type fileInfoList struct {
	FileInfos []fileInfo `json:"file_infos"`
}

func (c *config) files(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var allFilesInConfig fileInfoList

	for _, r := range c.Contents {
		p := r.Prefix

		fi := fileInfo{
			Prefix: p,
		}

		for _, f := range r.Files {
			fi.Files = append(fi.Files, f.Name)
		}

		allFilesInConfig.FileInfos = append(allFilesInConfig.FileInfos, fi)
	}

	enc := json.NewEncoder(w)
	enc.Encode(allFilesInConfig)
}
