package main

import "fmt"

func main() {
	config := newConfiguration("config.json")
	fmt.Println(config)

	err := generateCpp(config)
	if err != nil {
		fmt.Printf("There was an error: %+v\n", err)
	}
	// config.save("config.json")
}
