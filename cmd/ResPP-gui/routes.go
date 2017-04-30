package main

import (
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
