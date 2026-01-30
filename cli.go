package main

import (
	"bufio"
	"fmt"
	"os"
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

// DisplayAllSortedByTitle shows all books and magazines sorted by title
func (cli *CLI) DisplayAllSortedByTitle() {
	items := cli.library.GetAllSortedByTitle()

	fmt.Println("\n=== ALL ITEMS SORTED BY TITLE ===")

	for _, item := range items {
		if item.Type == "book" {
			cli.displayBook(item.Book)
		} else {
			cli.displayMagazine(item.Magazine)
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

// AddItem provides an interactive prompt to add a book or magazine
func (cli *CLI) AddItem(resourcesDir string) error {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("\n=== ADD ITEM TO LIBRARY ===")
	fmt.Println("What would you like to add?")
	fmt.Println("1. Book")
	fmt.Println("2. Magazine")
	fmt.Print("Enter choice (1 or 2): ")

	choice, _ := reader.ReadString('\n')
	choice = strings.TrimSpace(choice)

	if choice == "1" {
		return cli.addBook(reader, resourcesDir)
	} else if choice == "2" {
		return cli.addMagazine(reader, resourcesDir)
	} else {
		return fmt.Errorf("invalid choice")
	}
}

func (cli *CLI) addBook(reader *bufio.Reader, resourcesDir string) error {
	fmt.Println("\n--- Adding a Book ---")

	fmt.Print("Title: ")
	title, _ := reader.ReadString('\n')
	title = strings.TrimSpace(title)

	fmt.Print("ISBN: ")
	isbn, _ := reader.ReadString('\n')
	isbn = strings.TrimSpace(isbn)

	fmt.Print("Description: ")
	description, _ := reader.ReadString('\n')
	description = strings.TrimSpace(description)

	// Get authors
	authorEmails := cli.getAuthors(reader)

	// Add the book
	cli.library.AddBook(title, isbn, description, authorEmails)

	// Save to CSV
	if err := cli.library.SaveToCSV(resourcesDir); err != nil {
		return fmt.Errorf("failed to save: %w", err)
	}

	fmt.Printf("\n✓ Book '%s' added successfully!\n", title)
	return nil
}

func (cli *CLI) addMagazine(reader *bufio.Reader, resourcesDir string) error {
	fmt.Println("\n--- Adding a Magazine ---")

	fmt.Print("Title: ")
	title, _ := reader.ReadString('\n')
	title = strings.TrimSpace(title)

	fmt.Print("ISBN: ")
	isbn, _ := reader.ReadString('\n')
	isbn = strings.TrimSpace(isbn)

	fmt.Print("Published At (DD.MM.YYYY): ")
	publishedAt, _ := reader.ReadString('\n')
	publishedAt = strings.TrimSpace(publishedAt)

	// Get authors
	authorEmails := cli.getAuthors(reader)

	// Add the magazine
	cli.library.AddMagazine(title, isbn, publishedAt, authorEmails)

	// Save to CSV
	if err := cli.library.SaveToCSV(resourcesDir); err != nil {
		return fmt.Errorf("failed to save: %w", err)
	}

	fmt.Printf("\n✓ Magazine '%s' added successfully!\n", title)
	return nil
}

func (cli *CLI) getAuthors(reader *bufio.Reader) []string {
	var authorEmails []string

	fmt.Println("\nAuthors:")
	fmt.Println("Enter author details. Press Enter with empty email to finish.")

	for {
		fmt.Print("  Author email: ")
		email, _ := reader.ReadString('\n')
		email = strings.TrimSpace(email)

		if email == "" {
			break
		}

		// Check if author exists
		if _, exists := cli.library.Authors[email]; !exists {
			fmt.Print("  First name: ")
			firstName, _ := reader.ReadString('\n')
			firstName = strings.TrimSpace(firstName)

			fmt.Print("  Last name: ")
			lastName, _ := reader.ReadString('\n')
			lastName = strings.TrimSpace(lastName)

			cli.library.AddAuthor(email, firstName, lastName)
			fmt.Printf("  ✓ New author '%s %s' added\n", firstName, lastName)
		} else {
			author := cli.library.Authors[email]
			fmt.Printf("  ✓ Using existing author: %s %s\n", author.FirstName, author.LastName)
		}

		authorEmails = append(authorEmails, email)
	}

	return authorEmails
}
