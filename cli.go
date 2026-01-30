package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func startCLI(lib *Library) {
	reader := bufio.NewReader(os.Stdin)

	printCLIWelcome()

	for {
		fmt.Println("\n" + strings.Repeat("=", 50))
		printMainMenu()
		fmt.Print("Enter your choice (1-7): ")

		choice, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(choice)

		switch choice {
		case "1":
			displayAllBooks(lib)
		case "2":
			displayAllMagazines(lib)
		case "3":
			displaySortedByTitle(lib)
		case "4":
			searchByISBN(lib, reader)
		case "5":
			searchByAuthorEmail(lib, reader)
		case "6":
			addItem(lib, reader)
		case "7":
			fmt.Println("\nGoodbye!")
			return
		default:
			fmt.Println("Invalid choice. Please try again.")
		}
	}
}

func printCLIWelcome() {
	fmt.Println("\n" + strings.Repeat("=", 50))
	fmt.Println("     📚 Welcome to the Library System CLI 📚")
	fmt.Println(strings.Repeat("=", 50))
}

func printMainMenu() {
	fmt.Println("\nMain Menu:")
	fmt.Println("1. Display all books")
	fmt.Println("2. Display all magazines")
	fmt.Println("3. Display all items sorted by title")
	fmt.Println("4. Search by ISBN")
	fmt.Println("5. Search by author email")
	fmt.Println("6. Add a new item")
	fmt.Println("7. Exit")
}

func displayAllBooks(lib *Library) {
	lib.mu.RLock()
	defer lib.mu.RUnlock()

	if len(lib.Books) == 0 {
		fmt.Println("\nNo books found in the library.")
		return
	}

	fmt.Println("\n" + strings.Repeat("-", 50))
	fmt.Printf("Total Books: %d\n", len(lib.Books))
	fmt.Println(strings.Repeat("-", 50))

	for i, book := range lib.Books {
		printBookDetails(book, i+1)
	}
}

func displayAllMagazines(lib *Library) {
	lib.mu.RLock()
	defer lib.mu.RUnlock()

	if len(lib.Magazines) == 0 {
		fmt.Println("\nNo magazines found in the library.")
		return
	}

	fmt.Println("\n" + strings.Repeat("-", 50))
	fmt.Printf("Total Magazines: %d\n", len(lib.Magazines))
	fmt.Println(strings.Repeat("-", 50))

	for i, magazine := range lib.Magazines {
		printMagazineDetails(magazine, i+1)
	}
}

func displaySortedByTitle(lib *Library) {
	items := lib.GetAllItemsSortedByTitle()

	if len(items) == 0 {
		fmt.Println("\nNo items found in the library.")
		return
	}

	fmt.Println("\n" + strings.Repeat("-", 50))
	fmt.Printf("Total Items (sorted by title): %d\n", len(items))
	fmt.Println(strings.Repeat("-", 50))

	for i, item := range items {
		printItemDetails(item, i+1)
	}
}

func searchByISBN(lib *Library, reader *bufio.Reader) {
	fmt.Print("\nEnter ISBN to search: ")
	isbn, _ := reader.ReadString('\n')
	isbn = strings.TrimSpace(isbn)

	if isbn == "" {
		fmt.Println("ISBN cannot be empty.")
		return
	}

	bookResults := lib.SearchBooksByISBN(isbn)
	magazineResults := lib.SearchMagazinesByISBN(isbn)

	fmt.Println("\n" + strings.Repeat("-", 50))
	fmt.Printf("Search Results for ISBN: %s\n", isbn)
	fmt.Println(strings.Repeat("-", 50))

	if len(bookResults) == 0 && len(magazineResults) == 0 {
		fmt.Println("No items found matching this ISBN.")
		return
	}

	if len(bookResults) > 0 {
		fmt.Printf("\nBooks (%d):\n", len(bookResults))
		for i, book := range bookResults {
			printBookDetails(book, i+1)
		}
	}

	if len(magazineResults) > 0 {
		fmt.Printf("\nMagazines (%d):\n", len(magazineResults))
		for i, magazine := range magazineResults {
			printMagazineDetails(magazine, i+1)
		}
	}
}

func searchByAuthorEmail(lib *Library, reader *bufio.Reader) {
	fmt.Print("\nEnter author email to search: ")
	email, _ := reader.ReadString('\n')
	email = strings.TrimSpace(email)

	if email == "" {
		fmt.Println("Email cannot be empty.")
		return
	}

	bookResults := lib.SearchBooksByAuthorEmail(email)
	magazineResults := lib.SearchMagazinesByAuthorEmail(email)

	fmt.Println("\n" + strings.Repeat("-", 50))
	fmt.Printf("Search Results for Author: %s\n", email)
	fmt.Println(strings.Repeat("-", 50))

	if len(bookResults) == 0 && len(magazineResults) == 0 {
		fmt.Println("No items found for this author.")
		return
	}

	if len(bookResults) > 0 {
		fmt.Printf("\nBooks (%d):\n", len(bookResults))
		for i, book := range bookResults {
			printBookDetails(book, i+1)
		}
	}

	if len(magazineResults) > 0 {
		fmt.Printf("\nMagazines (%d):\n", len(magazineResults))
		for i, magazine := range magazineResults {
			printMagazineDetails(magazine, i+1)
		}
	}
}

func addItem(lib *Library, reader *bufio.Reader) {
	fmt.Println("\n" + strings.Repeat("-", 50))
	fmt.Println("Add a New Item")
	fmt.Println(strings.Repeat("-", 50))

	fmt.Println("\nSelect item type:")
	fmt.Println("1. Book")
	fmt.Println("2. Magazine")
	fmt.Print("Enter choice (1 or 2): ")

	itemType, _ := reader.ReadString('\n')
	itemType = strings.TrimSpace(itemType)

	var isBook bool
	if itemType == "1" {
		isBook = true
	} else if itemType == "2" {
		isBook = false
	} else {
		fmt.Println("Invalid choice.")
		return
	}

	// Get common fields
	fmt.Print("\nEnter title: ")
	title, _ := reader.ReadString('\n')
	title = strings.TrimSpace(title)

	if title == "" {
		fmt.Println("Title cannot be empty.")
		return
	}

	fmt.Print("Enter ISBN: ")
	isbn, _ := reader.ReadString('\n')
	isbn = strings.TrimSpace(isbn)

	if isbn == "" {
		fmt.Println("ISBN cannot be empty.")
		return
	}

	// Get type-specific fields
	var description, publishedAt string
	if isBook {
		fmt.Print("Enter description: ")
		description, _ = reader.ReadString('\n')
		description = strings.TrimSpace(description)
	} else {
		fmt.Print("Enter published date: ")
		publishedAt, _ = reader.ReadString('\n')
		publishedAt = strings.TrimSpace(publishedAt)
	}

	// Get authors
	var authors []Author
	for {
		fmt.Println("\nAdd an author:")
		fmt.Print("Enter author email (or press Enter to skip): ")
		email, _ := reader.ReadString('\n')
		email = strings.TrimSpace(email)

		if email == "" {
			if len(authors) == 0 {
				fmt.Println("At least one author is required.")
				continue
			}
			break
		}

		fmt.Print("Enter first name: ")
		firstName, _ := reader.ReadString('\n')
		firstName = strings.TrimSpace(firstName)

		fmt.Print("Enter last name: ")
		lastName, _ := reader.ReadString('\n')
		lastName = strings.TrimSpace(lastName)

		author := Author{
			Email:     email,
			FirstName: firstName,
			LastName:  lastName,
		}

		// Add author to library if not exists
		if err := lib.AddAuthor(author); err != nil {
			fmt.Printf("Error adding author: %v\n", err)
			continue
		}

		authors = append(authors, author)

		fmt.Print("\nAdd another author? (y/n): ")
		response, _ := reader.ReadString('\n')
		response = strings.TrimSpace(strings.ToLower(response))
		if response != "y" {
			break
		}
	}

	if len(authors) == 0 {
		fmt.Println("No authors added. Item not created.")
		return
	}

	// Add the item
	if isBook {
		book := Book{
			Title:       title,
			ISBN:        isbn,
			Authors:     authors,
			Description: description,
		}

		if err := lib.AddBook(book); err != nil {
			fmt.Printf("Error adding book: %v\n", err)
			return
		}

		fmt.Println("\n✓ Book added successfully!")
		printBookDetails(book, 1)
	} else {
		magazine := Magazine{
			Title:       title,
			ISBN:        isbn,
			Authors:     authors,
			PublishedAt: publishedAt,
		}

		if err := lib.AddMagazine(magazine); err != nil {
			fmt.Printf("Error adding magazine: %v\n", err)
			return
		}

		fmt.Println("\n✓ Magazine added successfully!")
		printMagazineDetails(magazine, 1)
	}
}

func printBookDetails(book Book, index int) {
	fmt.Printf("\n%d. [BOOK] %s\n", index, book.Title)
	fmt.Printf("   ISBN: %s\n", book.ISBN)
	fmt.Printf("   Description: %s\n", book.Description)
	fmt.Print("   Authors: ")
	if len(book.Authors) > 0 {
		authorNames := make([]string, len(book.Authors))
		for i, author := range book.Authors {
			authorNames[i] = fmt.Sprintf("%s %s (%s)", author.FirstName, author.LastName, author.Email)
		}
		fmt.Println(strings.Join(authorNames, ", "))
	} else {
		fmt.Println("Unknown")
	}
}

func printMagazineDetails(magazine Magazine, index int) {
	fmt.Printf("\n%d. [MAGAZINE] %s\n", index, magazine.Title)
	fmt.Printf("   ISBN: %s\n", magazine.ISBN)
	fmt.Printf("   Published: %s\n", magazine.PublishedAt)
	fmt.Print("   Contributors: ")
	if len(magazine.Authors) > 0 {
		authorNames := make([]string, len(magazine.Authors))
		for i, author := range magazine.Authors {
			authorNames[i] = fmt.Sprintf("%s %s (%s)", author.FirstName, author.LastName, author.Email)
		}
		fmt.Println(strings.Join(authorNames, ", "))
	} else {
		fmt.Println("Unknown")
	}
}

func printItemDetails(item Item, index int) {
	itemType := strings.ToUpper(item.Type)
	fmt.Printf("\n%d. [%s] %s\n", index, itemType, item.Title)
	fmt.Printf("   ISBN: %s\n", item.ISBN)

	if item.Type == "book" && item.Description != "" {
		fmt.Printf("   Description: %s\n", item.Description)
	} else if item.Type == "magazine" && item.PublishedAt != "" {
		fmt.Printf("   Published: %s\n", item.PublishedAt)
	}

	fmt.Print("   Authors: ")
	if len(item.Authors) > 0 {
		authorNames := make([]string, len(item.Authors))
		for i, author := range item.Authors {
			authorNames[i] = fmt.Sprintf("%s %s (%s)", author.FirstName, author.LastName, author.Email)
		}
		fmt.Println(strings.Join(authorNames, ", "))
	} else {
		fmt.Println("Unknown")
	}
}
