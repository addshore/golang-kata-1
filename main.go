package main

import (
	"fmt"
	"log"
)

func main() {
	authors, err := readAuthorsFromCSV("resources/authors.csv")
	if err != nil {
		log.Fatalf("Failed to read authors: %v", err)
	}
	fmt.Printf("Loaded %d authors\n", len(authors))

	books, err := readBooksFromCSV("resources/books.csv")
	if err != nil {
		log.Fatalf("Failed to read books: %v", err)
	}
	fmt.Printf("Loaded %d books\n", len(books))

	magazines, err := readMagazinesFromCSV("resources/magazines.csv")
	if err != nil {
		log.Fatalf("Failed to read magazines: %v", err)
	}
	fmt.Printf("Loaded %d magazines\n", len(magazines))

	startWebServer(books, magazines, authors)
}

func welcomeMessage() string {
	return "Hello world!"
}
