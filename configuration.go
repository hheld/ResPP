package main

import (
	"bufio"
	"encoding/json"
	"os"
)

type resource struct {
	Prefix string   `json:"prefix"`
	Files  []string `json:"files"`
}

type configuration struct {
	Contents []resource `json:"contents"`
}

var config configuration

func init() {
	err := config.load("config.json")

	if err != nil {
		config = configuration{
			Contents: []resource{
				{
					Prefix: "/",
					Files:  []string{"file1", "file2"},
				},
				{
					Prefix: "/test",
					Files:  []string{"file3", "file4"},
				},
			},
		}
	}
}

func (c *configuration) save(fileName string) error {
	f, err := os.Create(fileName)

	if err != nil {
		return err
	}

	defer f.Close()

	w := bufio.NewWriter(f)

	enc := json.NewEncoder(w)

	enc.SetIndent("", "\t")
	enc.Encode(c)

	w.Flush()

	return nil
}

func (c *configuration) load(fileName string) error {
	f, err := os.Open(fileName)

	if err != nil {
		return err
	}

	defer f.Close()

	r := bufio.NewReader(f)

	dec := json.NewDecoder(r)

	dec.Decode(c)

	return nil
}
