package configuration

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/hheld/ResPP/pkg/errors"
)

// File consists of a file name and its base64 encoded content.
type File struct {
	Name           string `json:"name"`
	EncodedContent string `json:"content"`
}

// Resource holds a list of all files sorted under a prefix string.
type Resource struct {
	Prefix string `json:"prefix"`
	Files  []File `json:"files"`
}

// Configuration contains all information that is stored in the config file.
type Configuration struct {
	Contents []Resource `json:"contents"`
}

// OpenConfiguration opens an existing config file specified by its path.
func OpenConfiguration(configFile string) (*Configuration, error) {
	var config Configuration

	err := config.LoadFromFile(configFile)

	if err != nil {
		return nil, err
	}

	err = config.checkFileExistences()

	if err != nil {
		fmt.Printf("There was an error: %v\n", err)
		return nil, err
	}

	err = config.readFiles()

	if err != nil {
		fmt.Printf("There was an error: %v\n", err)
		return nil, err
	}

	return &config, nil
}

func (c *Configuration) save(fileName string) error {
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

// LoadFromFile loads and fills a Configuration object from the given file.
func (c *Configuration) LoadFromFile(fileName string) error {
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

func (c *Configuration) checkFileExistences() error {
	var missingFiles errors.MissingFilesError

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

func (c *Configuration) readFiles() error {
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
