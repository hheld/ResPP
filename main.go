package main

import "fmt"

func main() {
	fmt.Println(config)
	config.save("config.json")
}
