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

	// Create CLI
	cli := NewCLI(library)

	// Check for command-line arguments
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "search-isbn":
			if len(os.Args) < 3 {
				fmt.Println("Usage: go run . search-isbn <ISBN>")
				os.Exit(1)
			}
			cli.SearchByISBN(os.Args[2])
		case "search-author":
			if len(os.Args) < 3 {
				fmt.Println("Usage: go run . search-author <email>")
				os.Exit(1)
			}
			cli.SearchByAuthorEmail(os.Args[2])
		case "list":
			cli.DisplayAll()
		default:
			fmt.Printf("Unknown command: %s\n", os.Args[1])
			fmt.Println("Available commands: list, search-isbn <ISBN>, search-author <email>")
			os.Exit(1)
		}
	} else {
		// Default: display all
		cli.DisplayAll()
	}
}

func welcomeMessage() string {
	return "Hello world!"
}
