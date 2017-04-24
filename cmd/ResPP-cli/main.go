package main

import (
	"flag"
	"fmt"

	"github.com/hheld/ResPP/pkg/codegen"
	"github.com/hheld/ResPP/pkg/configuration"
)

func main() {
	configFile := flag.String("configFile", "config.json", "the config file to load")
	outDir := flag.String("outDir", "cpp", "directory where the generated files will be written to")
	flag.Parse()

	var config *configuration.Configuration

	config, err := configuration.OpenConfiguration(*configFile)

	if err != nil {
		fmt.Printf("There was an error: %+v\n", err)
		return
	}

	err = codegen.GenerateCpp(config, *outDir)

	if err != nil {
		fmt.Printf("There was an error: %+v\n", err)
	}
}
