package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/hheld/ResPP/pkg/configuration"
)

func addLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s -- %s\n", r.Method, r.URL)
		next.ServeHTTP(w, r)
	})
}

type config configuration.Configuration

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

	enc := json.NewEncoder(w)
	enc.Encode(c)
}

func main() {
	log.Println("Server is about to listen at port 8000.")

	conf := &config{}

	mux := http.NewServeMux()

	mux.Handle("/", http.FileServer(http.Dir("./dist")))
	mux.HandleFunc("/load", conf.load)

	if err := http.ListenAndServe("localhost:8000", addLog(mux)); err != nil {
		log.Printf("Could not start server at port 8000: %v\n", err)
	}
}
