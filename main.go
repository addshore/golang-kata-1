package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println(welcomeMessage())

	// Load library data
	library, err := LoadLibrary("./resources")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading library: %v\n", err)
		os.Exit(1)
	}

	// Create and use CLI
	cli := NewCLI(library)
	cli.DisplayAll()
}

func welcomeMessage() string {
	return "Hello world!"
}
