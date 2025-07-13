package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)


func main() {
	// Check if web mode is requested
	if len(os.Args) > 1 && os.Args[1] == "web" {
		runWebMode()
		return
	}
	
	// Default to CLI mode
	runCLIMode()
}

func runWebMode() {
	fmt.Println("=== Library Management System - Web Mode ===")
	
	library := NewLibrary()
	err := library.LoadFromCSV("resources/authors.csv", "resources/books.csv", "resources/magazines.csv")
	if err != nil {
		fmt.Printf("Error loading library data: %v\n", err)
		return
	}
	
	webServer := NewWebServer(library)
	webServer.Start("8080")
}

func runCLIMode() {
	fmt.Println("=== Library Management System - CLI Mode ===")
	
	library := NewLibrary()
	err := library.LoadFromCSV("resources/authors.csv", "resources/books.csv", "resources/magazines.csv")
	if err != nil {
		fmt.Printf("Error loading library data: %v\n", err)
		return
	}
	
	for {
		fmt.Println("\nPlease select an option:")
		fmt.Println("1. Display all books")
		fmt.Println("2. Display all magazines")
		fmt.Println("3. Display all items")
		fmt.Println("4. Display all items sorted by title")
		fmt.Println("5. Search by ISBN")
		fmt.Println("6. Search by author email")
		fmt.Println("7. Add new item")
		fmt.Println("8. Exit")
		fmt.Print("Enter your choice (1-8): ")
		
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		choice := strings.TrimSpace(scanner.Text())
		
		switch choice {
		case "1":
			displayBooks(library.Books)
		case "2":
			displayMagazines(library.Magazines)
		case "3":
			displayBooks(library.Books)
			displayMagazines(library.Magazines)
		case "4":
			displayAllItemsSorted(library)
		case "5":
			searchByISBN(library)
		case "6":
			searchByAuthorEmail(library)
		case "7":
			addNewItem(library)
		case "8":
			fmt.Println("Goodbye!")
			return
		default:
			fmt.Println("Invalid choice. Please try again.")
		}
	}
}


func displayBooks(books []Book) {
	fmt.Println("\n=== BOOKS ===")
	for i, book := range books {
		fmt.Printf("%d. %s\n", i+1, book.Title)
		fmt.Printf("   ISBN: %s\n", book.ISBN)
		fmt.Printf("   Authors: %s\n", FormatAuthors(book.Authors))
		fmt.Printf("   Description: %s\n", book.Description)
		fmt.Println()
	}
}

func displayMagazines(magazines []Magazine) {
	fmt.Println("\n=== MAGAZINES ===")
	for i, magazine := range magazines {
		fmt.Printf("%d. %s\n", i+1, magazine.Title)
		fmt.Printf("   ISBN: %s\n", magazine.ISBN)
		fmt.Printf("   Authors: %s\n", FormatAuthors(magazine.Authors))
		fmt.Printf("   Published: %s\n", magazine.PublishedAt)
		fmt.Println()
	}
}

func searchByISBN(library *Library) {
	fmt.Print("Enter ISBN to search for: ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	isbn := strings.TrimSpace(scanner.Text())
	
	if isbn == "" {
		fmt.Println("ISBN cannot be empty.")
		return
	}
	
	found := false
	
	book := library.FindBookByISBN(isbn)
	if book != nil {
		fmt.Println("\n=== BOOK FOUND ===")
		fmt.Printf("Title: %s\n", book.Title)
		fmt.Printf("ISBN: %s\n", book.ISBN)
		fmt.Printf("Authors: %s\n", FormatAuthors(book.Authors))
		fmt.Printf("Description: %s\n", book.Description)
		fmt.Println()
		found = true
	}
	
	magazine := library.FindMagazineByISBN(isbn)
	if magazine != nil {
		fmt.Println("\n=== MAGAZINE FOUND ===")
		fmt.Printf("Title: %s\n", magazine.Title)
		fmt.Printf("ISBN: %s\n", magazine.ISBN)
		fmt.Printf("Authors: %s\n", FormatAuthors(magazine.Authors))
		fmt.Printf("Published: %s\n", magazine.PublishedAt)
		fmt.Println()
		found = true
	}
	
	if !found {
		fmt.Printf("No book or magazine found with ISBN: %s\n", isbn)
	}
}

func searchByAuthorEmail(library *Library) {
	fmt.Print("Enter author email to search for: ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	email := strings.TrimSpace(scanner.Text())
	
	if email == "" {
		fmt.Println("Email cannot be empty.")
		return
	}
	
	foundBooks := library.FindBooksByAuthorEmail(email)
	foundMagazines := library.FindMagazinesByAuthorEmail(email)
	
	if len(foundBooks) == 0 && len(foundMagazines) == 0 {
		fmt.Printf("No books or magazines found with author email: %s\n", email)
		return
	}
	
	if len(foundBooks) > 0 {
		fmt.Printf("\n=== BOOKS BY AUTHOR (%s) ===\n", email)
		for i, book := range foundBooks {
			fmt.Printf("%d. %s\n", i+1, book.Title)
			fmt.Printf("   ISBN: %s\n", book.ISBN)
			fmt.Printf("   Authors: %s\n", FormatAuthors(book.Authors))
			fmt.Printf("   Description: %s\n", book.Description)
			fmt.Println()
		}
	}
	
	if len(foundMagazines) > 0 {
		fmt.Printf("\n=== MAGAZINES BY AUTHOR (%s) ===\n", email)
		for i, magazine := range foundMagazines {
			fmt.Printf("%d. %s\n", i+1, magazine.Title)
			fmt.Printf("   ISBN: %s\n", magazine.ISBN)
			fmt.Printf("   Authors: %s\n", FormatAuthors(magazine.Authors))
			fmt.Printf("   Published: %s\n", magazine.PublishedAt)
			fmt.Println()
		}
	}
}

func displayAllItemsSorted(library *Library) {
	fmt.Println("\nSelect sort direction:")
	fmt.Println("1. Ascending (A-Z)")
	fmt.Println("2. Descending (Z-A)")
	fmt.Print("Enter your choice (1-2): ")
	
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	sortChoice := strings.TrimSpace(scanner.Text())
	
	var ascending bool
	switch sortChoice {
	case "1":
		ascending = true
	case "2":
		ascending = false
	default:
		fmt.Println("Invalid choice. Using ascending order (A-Z).")
		ascending = true
	}
	
	allItems := library.GetAllItemsSorted(ascending)
	
	if ascending {
		fmt.Println("\n=== ALL ITEMS SORTED BY TITLE (A-Z) ===")
	} else {
		fmt.Println("\n=== ALL ITEMS SORTED BY TITLE (Z-A) ===")
	}
	
	for i, item := range allItems {
		fmt.Printf("%d. %s (%s)\n", i+1, item.Title, item.Type)
		fmt.Printf("   ISBN: %s\n", item.ISBN)
		fmt.Printf("   Authors: %s\n", FormatAuthors(item.Authors))
		if item.Type == "Book" {
			fmt.Printf("   Description: %s\n", item.Description)
		} else {
			fmt.Printf("   Published: %s\n", item.PublishedAt)
		}
		fmt.Println()
	}
}

func addNewItem(library *Library) {
	fmt.Println("\nWhat would you like to add?")
	fmt.Println("1. Book")
	fmt.Println("2. Magazine")
	fmt.Print("Enter your choice (1-2): ")
	
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	choice := strings.TrimSpace(scanner.Text())
	
	switch choice {
	case "1":
		addNewBook(library)
	case "2":
		addNewMagazine(library)
	default:
		fmt.Println("Invalid choice. Returning to main menu.")
		return
	}
}

func addNewBook(library *Library) {
	fmt.Println("\n=== ADD NEW BOOK ===")
	scanner := bufio.NewScanner(os.Stdin)
	
	// Get book details
	fmt.Print("Enter book title: ")
	scanner.Scan()
	title := strings.TrimSpace(scanner.Text())
	if title == "" {
		fmt.Println("Title cannot be empty. Operation cancelled.")
		return
	}
	
	fmt.Print("Enter ISBN: ")
	scanner.Scan()
	isbn := strings.TrimSpace(scanner.Text())
	if isbn == "" {
		fmt.Println("ISBN cannot be empty. Operation cancelled.")
		return
	}
	
	// Check if ISBN already exists
	if library.FindBookByISBN(isbn) != nil || library.FindMagazineByISBN(isbn) != nil {
		fmt.Printf("ISBN %s already exists in the library. Operation cancelled.\n", isbn)
		return
	}
	
	fmt.Print("Enter description: ")
	scanner.Scan()
	description := strings.TrimSpace(scanner.Text())
	if description == "" {
		fmt.Println("Description cannot be empty. Operation cancelled.")
		return
	}
	
	// Get authors
	authors := getAuthorsFromUser(library, scanner)
	if len(authors) == 0 {
		fmt.Println("At least one author is required. Operation cancelled.")
		return
	}
	
	// Create and add the book
	book := Book{
		Title:       title,
		ISBN:        isbn,
		Authors:     authors,
		Description: description,
	}
	
	library.AddBook(book)
	
	// Save to CSV
	if err := library.SaveToCSV("resources/authors.csv", "resources/books.csv", "resources/magazines.csv"); err != nil {
		fmt.Printf("Error saving to CSV files: %v\n", err)
		return
	}
	
	fmt.Printf("Book '%s' has been successfully added to the library!\n", title)
}

func addNewMagazine(library *Library) {
	fmt.Println("\n=== ADD NEW MAGAZINE ===")
	scanner := bufio.NewScanner(os.Stdin)
	
	// Get magazine details
	fmt.Print("Enter magazine title: ")
	scanner.Scan()
	title := strings.TrimSpace(scanner.Text())
	if title == "" {
		fmt.Println("Title cannot be empty. Operation cancelled.")
		return
	}
	
	fmt.Print("Enter ISBN: ")
	scanner.Scan()
	isbn := strings.TrimSpace(scanner.Text())
	if isbn == "" {
		fmt.Println("ISBN cannot be empty. Operation cancelled.")
		return
	}
	
	// Check if ISBN already exists
	if library.FindBookByISBN(isbn) != nil || library.FindMagazineByISBN(isbn) != nil {
		fmt.Printf("ISBN %s already exists in the library. Operation cancelled.\n", isbn)
		return
	}
	
	fmt.Print("Enter published date (DD.MM.YYYY): ")
	scanner.Scan()
	publishedAt := strings.TrimSpace(scanner.Text())
	if publishedAt == "" {
		fmt.Println("Published date cannot be empty. Operation cancelled.")
		return
	}
	
	// Get authors
	authors := getAuthorsFromUser(library, scanner)
	if len(authors) == 0 {
		fmt.Println("At least one author is required. Operation cancelled.")
		return
	}
	
	// Create and add the magazine
	magazine := Magazine{
		Title:       title,
		ISBN:        isbn,
		Authors:     authors,
		PublishedAt: publishedAt,
	}
	
	library.AddMagazine(magazine)
	
	// Save to CSV
	if err := library.SaveToCSV("resources/authors.csv", "resources/books.csv", "resources/magazines.csv"); err != nil {
		fmt.Printf("Error saving to CSV files: %v\n", err)
		return
	}
	
	fmt.Printf("Magazine '%s' has been successfully added to the library!\n", title)
}

func getAuthorsFromUser(library *Library, scanner *bufio.Scanner) []Author {
	var authors []Author
	
	fmt.Println("\nEnter author information (you can add multiple authors):")
	
	for {
		fmt.Print("Enter author email: ")
		scanner.Scan()
		email := strings.TrimSpace(scanner.Text())
		if email == "" {
			fmt.Println("Email cannot be empty.")
			continue
		}
		
		// Check if author already exists
		if existingAuthor, exists := library.Authors[email]; exists {
			fmt.Printf("Author %s %s (%s) already exists in the library.\n", 
				existingAuthor.FirstName, existingAuthor.LastName, existingAuthor.Email)
			authors = append(authors, existingAuthor)
		} else {
			// Get new author details
			fmt.Print("Enter first name: ")
			scanner.Scan()
			firstName := strings.TrimSpace(scanner.Text())
			if firstName == "" {
				fmt.Println("First name cannot be empty.")
				continue
			}
			
			fmt.Print("Enter last name: ")
			scanner.Scan()
			lastName := strings.TrimSpace(scanner.Text())
			if lastName == "" {
				fmt.Println("Last name cannot be empty.")
				continue
			}
			
			// Create new author
			author := library.GetOrCreateAuthor(email, firstName, lastName)
			authors = append(authors, author)
			fmt.Printf("New author %s %s has been added.\n", firstName, lastName)
		}
		
		// Ask if user wants to add more authors
		fmt.Print("Add another author? (y/n): ")
		scanner.Scan()
		addMore := strings.ToLower(strings.TrimSpace(scanner.Text()))
		if addMore != "y" && addMore != "yes" {
			break
		}
	}
	
	return authors
}

func welcomeMessage() string {
	return "Hello world!"
}
