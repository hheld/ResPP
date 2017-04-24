package main

import "fmt"

func main() {
	config := newConfiguration("config.json")

	err := generateCpp(config, "cpp")
	if err != nil {
		fmt.Printf("There was an error: %+v\n", err)
	}
	// config.save("config.json")
}
