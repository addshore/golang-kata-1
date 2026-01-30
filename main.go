package main

import (
	"flag"
	"fmt"
	"log"
	"path/filepath"
)

func main() {
	cliMode := flag.Bool("cli", false, "Run in CLI mode")
	flag.Parse()

	fmt.Println("Starting Library Application...")

	// Define paths to resources
	resourcesDir := "resources"
	authorsPath := filepath.Join(resourcesDir, "authors.csv")
	booksPath := filepath.Join(resourcesDir, "books.csv")
	magsPath := filepath.Join(resourcesDir, "magazines.csv")

	// Load Data
	fmt.Println("Loading data...")
	authors, err := LoadAuthors(authorsPath)
	if err != nil {
		log.Fatalf("Failed to load authors: %v", err)
	}

	books, err := LoadBooks(booksPath, authors)
	if err != nil {
		log.Fatalf("Failed to load books: %v", err)
	}

	magazines, err := LoadMagazines(magsPath, authors)
	if err != nil {
		log.Fatalf("Failed to load magazines: %v", err)
	}

	// Calculate counts for startup log
	fmt.Printf("Loaded %d authors, %d books, %d magazines\n", len(authors), len(books), len(magazines))

	data := &LibraryData{
		Authors:   authors,
		Books:     books,
		Magazines: magazines,
	}

	if *cliMode {
		StartCLI(data)
	} else {
		// Start Server
		server := &Server{Data: data}
		addr := ":8099"
		fmt.Printf("Server starting on http://localhost%s\n", addr)
		if err := server.Start(addr); err != nil {
			log.Fatalf("Server failed: %v", err)
		}
	}
}
