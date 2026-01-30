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

	// Check for command-line arguments
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "web":
			// Start web server
			port := ":8080"
			if len(os.Args) > 2 {
				port = os.Args[2]
			}
			
			webServer, err := NewWebServer(library, "./resources")
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error creating web server: %v\n", err)
				os.Exit(1)
			}
			
			webServer.SetupRoutes()
			fmt.Printf("Web UI available at http://localhost%s\n", port)
			if err := webServer.Start(port); err != nil {
				fmt.Fprintf(os.Stderr, "Error starting web server: %v\n", err)
				os.Exit(1)
			}
		case "search-isbn":
			if len(os.Args) < 3 {
				fmt.Println("Usage: go run . search-isbn <ISBN>")
				os.Exit(1)
			}
			cli := NewCLI(library)
			cli.SearchByISBN(os.Args[2])
		case "search-author":
			if len(os.Args) < 3 {
				fmt.Println("Usage: go run . search-author <email>")
				os.Exit(1)
			}
			cli := NewCLI(library)
			cli.SearchByAuthorEmail(os.Args[2])
		case "sorted":
			cli := NewCLI(library)
			cli.DisplayAllSortedByTitle()
		case "add":
			cli := NewCLI(library)
			if err := cli.AddItem("./resources"); err != nil {
				fmt.Fprintf(os.Stderr, "Error adding item: %v\n", err)
				os.Exit(1)
			}
		case "list":
			cli := NewCLI(library)
			cli.DisplayAll()
		default:
			fmt.Printf("Unknown command: %s\n", os.Args[1])
			fmt.Println("Available commands:")
			fmt.Println("  list                    - List all books and magazines")
			fmt.Println("  search-isbn <ISBN>      - Search by ISBN")
			fmt.Println("  search-author <email>   - Search by author email")
			fmt.Println("  sorted                  - Display all sorted by title")
			fmt.Println("  add                     - Add a new book or magazine")
			fmt.Println("  web [port]              - Start web UI (default :8080)")
			os.Exit(1)
		}
	} else {
		// Default: display all
		cli := NewCLI(library)
		cli.DisplayAll()
	}
}

func welcomeMessage() string {
	return "Hello world!"
}
