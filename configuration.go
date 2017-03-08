package main

import (
	"bufio"
	"encoding/json"
	"fmt"
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

	err = config.checkFileExistences()

	if err != nil {
		fmt.Printf("There was an error: %v\n", err)
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

func (c *configuration) checkFileExistences() error {
	var missingFiles missingFilesError

	for _, r := range c.Contents {
		for _, f := range r.Files {
			if _, err := os.Stat(f); os.IsNotExist(err) {
				missingFiles = append(missingFiles, f)
			}
		}
	}

	if len(missingFiles) > 0 {
		return &missingFiles
	}

	return nil
}
