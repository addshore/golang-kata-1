package main

import (
	"fmt"
	"strings"
)

// CLI provides command-line interface functionality
type CLI struct {
	library *Library
}

// NewCLI creates a new CLI instance
func NewCLI(library *Library) *CLI {
	return &CLI{library: library}
}

// DisplayAll shows all books and magazines
func (cli *CLI) DisplayAll() {
	fmt.Println("\n=== BOOKS ===")
	for _, book := range cli.library.Books {
		cli.displayBook(book)
	}

	fmt.Println("\n=== MAGAZINES ===")
	for _, magazine := range cli.library.Magazines {
		cli.displayMagazine(magazine)
	}
}

func (cli *CLI) displayBook(book *Book) {
	authorNames := cli.library.GetAuthorNames(book.Authors)
	fmt.Printf("\nTitle: %s\n", book.Title)
	fmt.Printf("ISBN: %s\n", book.ISBN)
	fmt.Printf("Authors: %s\n", strings.Join(authorNames, ", "))
	fmt.Printf("Description: %s\n", book.Description)
	fmt.Println(strings.Repeat("-", 80))
}

func (cli *CLI) displayMagazine(magazine *Magazine) {
	authorNames := cli.library.GetAuthorNames(magazine.Authors)
	fmt.Printf("\nTitle: %s\n", magazine.Title)
	fmt.Printf("ISBN: %s\n", magazine.ISBN)
	fmt.Printf("Authors: %s\n", strings.Join(authorNames, ", "))
	fmt.Printf("Published At: %s\n", magazine.PublishedAt)
	fmt.Println(strings.Repeat("-", 80))
}
