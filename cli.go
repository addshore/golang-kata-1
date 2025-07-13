package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// CLIHandler handles the command-line interface
type CLIHandler struct {
	service *LibraryService
	scanner *bufio.Scanner
}

// NewCLIHandler creates a new CLI handler
func NewCLIHandler(library *Library) *CLIHandler {
	return &CLIHandler{
		service: NewLibraryService(library),
		scanner: bufio.NewScanner(os.Stdin),
	}
}

// Start starts the interactive CLI
func (c *CLIHandler) Start() {
	fmt.Println("=== Library Management System - CLI ===")
	fmt.Println()

	for {
		c.showMainMenu()
		choice := c.getInput("Enter your choice: ")

		switch choice {
		case "1":
			c.displayAllBooks()
		case "2":
			c.displayAllMagazines()
		case "3":
			c.displayAllItemsSorted()
		case "4":
			c.searchLibrary()
		case "5":
			c.addNewItem()
		case "6":
			c.displayStats()
		case "7":
			fmt.Println("Thank you for using the Library Management System!")
			return
		default:
			fmt.Println("Invalid choice. Please try again.")
		}

		fmt.Println()
		c.getInput("Press Enter to continue...")
		fmt.Println()
	}
}

// showMainMenu displays the main menu options
func (c *CLIHandler) showMainMenu() {
	fmt.Println("=== MAIN MENU ===")
	fmt.Println("1. Display All Books")
	fmt.Println("2. Display All Magazines")
	fmt.Println("3. Display All Items (Sorted)")
	fmt.Println("4. Search Library")
	fmt.Println("5. Add New Item")
	fmt.Println("6. Library Statistics")
	fmt.Println("7. Exit")
	fmt.Println()
}

// displayAllBooks shows all books
func (c *CLIHandler) displayAllBooks() {
	fmt.Println("=== ALL BOOKS ===")
	books := c.service.GetBooks()

	if len(books) == 0 {
		fmt.Println("No books found in the library.")
		return
	}

	for i, book := range books {
		fmt.Printf("%d. %s\n", i+1, book.Title)
		fmt.Printf("   ISBN: %s\n", book.ISBN)
		fmt.Printf("   Authors: %s\n", c.service.FormatAuthors(book.Authors))
		fmt.Printf("   Description: %s\n", book.Description)
		fmt.Println()
	}

	fmt.Printf("Total: %d books\n", len(books))
}

// displayAllMagazines shows all magazines
func (c *CLIHandler) displayAllMagazines() {
	fmt.Println("=== ALL MAGAZINES ===")
	magazines := c.service.GetMagazines()

	if len(magazines) == 0 {
		fmt.Println("No magazines found in the library.")
		return
	}

	for i, magazine := range magazines {
		fmt.Printf("%d. %s\n", i+1, magazine.Title)
		fmt.Printf("   ISBN: %s\n", magazine.ISBN)
		fmt.Printf("   Authors: %s\n", c.service.FormatAuthors(magazine.Authors))
		fmt.Printf("   Published: %s\n", magazine.PublishedAt.Format("02.01.2006"))
		fmt.Println()
	}

	fmt.Printf("Total: %d magazines\n", len(magazines))
}

// displayAllItemsSorted shows all items sorted by title
func (c *CLIHandler) displayAllItemsSorted() {
	fmt.Println("=== ALL ITEMS SORTED BY TITLE ===")

	direction := c.getInput("Sort direction (1=A-Z, 2=Z-A): ")
	ascending := direction != "2"

	items := c.service.GetAllItemsSorted(ascending)

	if len(items) == 0 {
		fmt.Println("No items found in the library.")
		return
	}

	sortDir := "A-Z"
	if !ascending {
		sortDir = "Z-A"
	}
	fmt.Printf("Sorted %s:\n\n", sortDir)

	for i, item := range items {
		fmt.Printf("%d. [%s] %s\n", i+1, item.Type, item.Title)
		fmt.Printf("   ISBN: %s\n", item.ISBN)
		fmt.Printf("   Authors: %s\n", item.AuthorNames)
		if item.Type == "BOOK" {
			fmt.Printf("   Description: %s\n", item.Description)
		} else {
			fmt.Printf("   Published: %s\n", item.PublishedAt)
		}
		fmt.Println()
	}

	bookCount, magazineCount, _ := c.service.GetLibraryStats()
	fmt.Printf("Total: %d items (%d books, %d magazines)\n", len(items), bookCount, magazineCount)
}

// searchLibrary handles search functionality
func (c *CLIHandler) searchLibrary() {
	fmt.Println("=== SEARCH LIBRARY ===")
	fmt.Println("1. Search by ISBN")
	fmt.Println("2. Search by Author Email")

	choice := c.getInput("Choose search type: ")
	searchTerm := c.getInput("Enter search term: ")

	var books []Book
	var magazines []Magazine

	switch choice {
	case "1":
		books, magazines = c.service.SearchByISBN(searchTerm)
		fmt.Printf("\n=== ISBN SEARCH RESULTS for '%s' ===\n", searchTerm)
	case "2":
		books, magazines = c.service.SearchByAuthorEmail(searchTerm)
		fmt.Printf("\n=== AUTHOR EMAIL SEARCH RESULTS for '%s' ===\n", searchTerm)
	default:
		fmt.Println("Invalid search type.")
		return
	}

	totalResults := len(books) + len(magazines)
	if totalResults == 0 {
		fmt.Println("No results found.")
		return
	}

	// Display books
	if len(books) > 0 {
		fmt.Println("\nBOOKS:")
		for i, book := range books {
			fmt.Printf("%d. %s\n", i+1, book.Title)
			fmt.Printf("   ISBN: %s\n", book.ISBN)
			fmt.Printf("   Authors: %s\n", c.service.FormatAuthors(book.Authors))
			fmt.Printf("   Description: %s\n", book.Description)
			fmt.Println()
		}
	}

	// Display magazines
	if len(magazines) > 0 {
		fmt.Println("MAGAZINES:")
		for i, magazine := range magazines {
			fmt.Printf("%d. %s\n", i+1, magazine.Title)
			fmt.Printf("   ISBN: %s\n", magazine.ISBN)
			fmt.Printf("   Authors: %s\n", c.service.FormatAuthors(magazine.Authors))
			fmt.Printf("   Published: %s\n", magazine.PublishedAt.Format("02.01.2006"))
			fmt.Println()
		}
	}

	fmt.Printf("Total results: %d (%d books, %d magazines)\n", totalResults, len(books), len(magazines))
}

// addNewItem handles adding new items
func (c *CLIHandler) addNewItem() {
	fmt.Println("=== ADD NEW ITEM ===")
	fmt.Println("1. Add Book")
	fmt.Println("2. Add Magazine")
	fmt.Println("3. Add Author")

	choice := c.getInput("Choose item type: ")

	switch choice {
	case "1":
		c.addBook()
	case "2":
		c.addMagazine()
	case "3":
		c.addAuthor()
	default:
		fmt.Println("Invalid choice.")
	}
}

// addBook adds a new book
func (c *CLIHandler) addBook() {
	fmt.Println("\n=== ADD NEW BOOK ===")

	title := c.getInput("Enter book title: ")
	isbn := c.getInput("Enter ISBN: ")
	description := c.getInput("Enter description (optional): ")

	// Show available authors
	authors := c.service.GetAuthorsList()
	if len(authors) == 0 {
		fmt.Println("No authors available. Please add an author first.")
		return
	}

	fmt.Println("\nAvailable authors:")
	for i, author := range authors {
		fmt.Printf("%d. %s %s (%s)\n", i+1, author.FirstName, author.LastName, author.Email)
	}

	authorInput := c.getInput("Enter author numbers (comma-separated, e.g., 1,3): ")
	authorEmails := c.parseAuthorSelection(authorInput, authors)

	if len(authorEmails) == 0 {
		fmt.Println("At least one author must be selected.")
		return
	}

	err := c.service.AddBook(title, isbn, description, authorEmails)
	if err != nil {
		fmt.Printf("Error adding book: %s\n", err.Error())
	} else {
		fmt.Printf("Book '%s' added successfully!\n", title)
	}
}

// addMagazine adds a new magazine
func (c *CLIHandler) addMagazine() {
	fmt.Println("\n=== ADD NEW MAGAZINE ===")

	title := c.getInput("Enter magazine title: ")
	isbn := c.getInput("Enter ISBN: ")
	publishedStr := c.getInput("Enter published date (YYYY-MM-DD): ")

	publishedAt, err := time.Parse("2006-01-02", publishedStr)
	if err != nil {
		fmt.Println("Invalid date format. Please use YYYY-MM-DD.")
		return
	}

	// Show available authors
	authors := c.service.GetAuthorsList()
	if len(authors) == 0 {
		fmt.Println("No authors available. Please add an author first.")
		return
	}

	fmt.Println("\nAvailable authors:")
	for i, author := range authors {
		fmt.Printf("%d. %s %s (%s)\n", i+1, author.FirstName, author.LastName, author.Email)
	}

	authorInput := c.getInput("Enter author numbers (comma-separated, e.g., 1,3): ")
	authorEmails := c.parseAuthorSelection(authorInput, authors)

	if len(authorEmails) == 0 {
		fmt.Println("At least one author must be selected.")
		return
	}

	err = c.service.AddMagazine(title, isbn, publishedAt, authorEmails)
	if err != nil {
		fmt.Printf("Error adding magazine: %s\n", err.Error())
	} else {
		fmt.Printf("Magazine '%s' added successfully!\n", title)
	}
}

// addAuthor adds a new author
func (c *CLIHandler) addAuthor() {
	fmt.Println("\n=== ADD NEW AUTHOR ===")

	email := c.getInput("Enter author email: ")
	firstName := c.getInput("Enter first name: ")
	lastName := c.getInput("Enter last name: ")

	err := c.service.AddAuthor(email, firstName, lastName)
	if err != nil {
		fmt.Printf("Error adding author: %s\n", err.Error())
	} else {
		fmt.Printf("Author '%s %s' added successfully!\n", firstName, lastName)
	}
}

// displayStats shows library statistics
func (c *CLIHandler) displayStats() {
	fmt.Println("=== LIBRARY STATISTICS ===")

	bookCount, magazineCount, authorCount := c.service.GetLibraryStats()

	fmt.Printf("Books: %d\n", bookCount)
	fmt.Printf("Magazines: %d\n", magazineCount)
	fmt.Printf("Authors: %d\n", authorCount)
	fmt.Printf("Total Items: %d\n", bookCount+magazineCount)
}

// getInput prompts for user input and returns the trimmed result
func (c *CLIHandler) getInput(prompt string) string {
	fmt.Print(prompt)
	c.scanner.Scan()
	return strings.TrimSpace(c.scanner.Text())
}

// parseAuthorSelection parses comma-separated author numbers and returns email addresses
func (c *CLIHandler) parseAuthorSelection(input string, authors []Author) []string {
	var emails []string

	if input == "" {
		return emails
	}

	parts := strings.Split(input, ",")
	for _, part := range parts {
		numStr := strings.TrimSpace(part)
		num, err := strconv.Atoi(numStr)
		if err != nil || num < 1 || num > len(authors) {
			fmt.Printf("Invalid author number: %s\n", numStr)
			continue
		}
		emails = append(emails, authors[num-1].Email)
	}

	return emails
}
