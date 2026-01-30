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

// SearchByISBN searches and displays results by ISBN
func (cli *CLI) SearchByISBN(isbn string) {
	books, magazines := cli.library.SearchByISBN(isbn)

	if len(books) == 0 && len(magazines) == 0 {
		fmt.Printf("\nNo results found for ISBN: %s\n", isbn)
		return
	}

	fmt.Printf("\n=== SEARCH RESULTS FOR ISBN: %s ===\n", isbn)

	if len(books) > 0 {
		fmt.Println("\n--- BOOKS ---")
		for _, book := range books {
			cli.displayBook(book)
		}
	}

	if len(magazines) > 0 {
		fmt.Println("\n--- MAGAZINES ---")
		for _, magazine := range magazines {
			cli.displayMagazine(magazine)
		}
	}
}

// SearchByAuthorEmail searches and displays results by author email
func (cli *CLI) SearchByAuthorEmail(email string) {
	books, magazines := cli.library.SearchByAuthorEmail(email)

	if len(books) == 0 && len(magazines) == 0 {
		fmt.Printf("\nNo results found for author email: %s\n", email)
		return
	}

	authorName := "Unknown"
	if author, ok := cli.library.Authors[email]; ok {
		authorName = fmt.Sprintf("%s %s", author.FirstName, author.LastName)
	}

	fmt.Printf("\n=== SEARCH RESULTS FOR AUTHOR: %s (%s) ===\n", authorName, email)

	if len(books) > 0 {
		fmt.Println("\n--- BOOKS ---")
		for _, book := range books {
			cli.displayBook(book)
		}
	}

	if len(magazines) > 0 {
		fmt.Println("\n--- MAGAZINES ---")
		for _, magazine := range magazines {
			cli.displayMagazine(magazine)
		}
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
