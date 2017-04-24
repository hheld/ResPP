package main

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type file struct {
	Name           string `json:"name"`
	EncodedContent string `json:"content"`
}

type resource struct {
	Prefix string `json:"prefix"`
	Files  []file `json:"files"`
}

type configuration struct {
	Contents []resource `json:"contents"`
}

func newConfiguration(configFile string) *configuration {
	var config configuration

	err := config.load(configFile)

	if err != nil {
		return nil
	}

	err = config.checkFileExistences()

	if err != nil {
		fmt.Printf("There was an error: %v\n", err)
		return nil
	}

	err = config.readFiles()

	if err != nil {
		fmt.Printf("There was an error: %v\n", err)
		return nil
	}

	return &config
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
			if _, err := os.Stat(f.Name); os.IsNotExist(err) {
				missingFiles = append(missingFiles, f.Name)
			}
		}
	}

	if len(missingFiles) > 0 {
		return &missingFiles
	}

	return nil
}

func (c *configuration) readFiles() error {
	for _, r := range c.Contents {
		for i, f := range r.Files {
			err := func() error {
				file, err := os.Open(f.Name)

				if err != nil {
					return err
				}

				defer file.Close()

				contents, err := ioutil.ReadAll(file)

				if err != nil {
					return err
				}

				r.Files[i].EncodedContent = base64.StdEncoding.EncodeToString(contents)

				return nil
			}()

			if err != nil {
				return err
			}
		}
	}

	return nil
}
