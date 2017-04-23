package main

import "fmt"

func main() {
	config := newConfiguration("config.json")
	fmt.Println(config)
	config.save("config.json")
}
