package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func StartCLI(data *LibraryData) {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println("\n--- EchoCat Library CLI ---")
		fmt.Println("1. View All Items (Sorted by Title)")
		fmt.Println("2. Search by ISBN")
		fmt.Println("3. Search by Author Email")
		fmt.Println("4. Add Item")
		fmt.Println("Q. Quit")
		fmt.Print("Select an option: ")

		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		switch strings.ToUpper(input) {
		case "1":
			viewAll(data)
		case "2":
			searchISBN(data, reader)
		case "3":
			searchEmail(data, reader)
		case "4":
			addItem(data, reader)
		case "Q":
			fmt.Println("Goodbye!")
			return
		default:
			fmt.Println("Invalid option. Please try again.")
		}
	}
}

func printItems(items []*DisplayItem) {
	if len(items) == 0 {
		fmt.Println("No items found.")
		return
	}
	fmt.Printf("%-10s | %-40s | %-30s | %s\n", "Type", "Title", "Authors", "ISBN")
	fmt.Println(strings.Repeat("-", 100))
	for _, item := range items {
		var authors []string
		for _, a := range item.Authors {
			authors = append(authors, fmt.Sprintf("%s %s", a.Firstname, a.Lastname))
		}
		authorStr := strings.Join(authors, ", ")
		if len(authorStr) > 30 {
			authorStr = authorStr[:27] + "..."
		}
		title := item.Title
		if len(title) > 40 {
			title = title[:37] + "..."
		}
		fmt.Printf("%-10s | %-40s | %-30s | %s\n", item.Type, title, authorStr, item.ISBN)
	}
}

func viewAll(data *LibraryData) {
	// Reusing service logic: empty filters, sort by title asc
	result := ProcessRequest(data, "", "", "title", "asc")
	printItems(result.SortedItems)
}

func searchISBN(data *LibraryData, reader *bufio.Reader) {
	fmt.Print("Enter ISBN: ")
	isbn, _ := reader.ReadString('\n')
	isbn = strings.TrimSpace(isbn)

	result := ProcessRequest(data, isbn, "", "title", "asc")
	printItems(result.SortedItems)
}

func searchEmail(data *LibraryData, reader *bufio.Reader) {
	fmt.Print("Enter Author Email: ")
	email, _ := reader.ReadString('\n')
	email = strings.TrimSpace(email)

	result := ProcessRequest(data, "", email, "title", "asc")
	printItems(result.SortedItems)
}

func addItem(data *LibraryData, reader *bufio.Reader) {
	fmt.Println("\nAdding New Item")

	fmt.Print("Item Type (Book/Magazine): ")
	itemType, _ := reader.ReadString('\n')
	itemType = strings.TrimSpace(itemType)
	if !strings.EqualFold(itemType, "Book") && !strings.EqualFold(itemType, "Magazine") {
		fmt.Println("Invalid item type.")
		return
	}

	fmt.Print("Title: ")
	title, _ := reader.ReadString('\n')
	title = strings.TrimSpace(title)

	fmt.Print("ISBN: ")
	isbn, _ := reader.ReadString('\n')
	isbn = strings.TrimSpace(isbn)

	fmt.Print("Extra Info (Description/Date): ")
	extra, _ := reader.ReadString('\n')
	extra = strings.TrimSpace(extra)

	fmt.Println("\nAuthor Details")
	fmt.Print("Email: ")
	email, _ := reader.ReadString('\n')
	email = strings.TrimSpace(email)

	fmt.Print("First Name: ")
	firstname, _ := reader.ReadString('\n')
	firstname = strings.TrimSpace(firstname)

	fmt.Print("Last Name: ")
	lastname, _ := reader.ReadString('\n')
	lastname = strings.TrimSpace(lastname)

	// Logic similar to server.go handleAdd
	author, exists := data.Authors[email]
	if !exists {
		author = &Author{
			Email:     email,
			Firstname: firstname,
			Lastname:  lastname,
		}
		if err := AppendAuthor(filepath.Join("resources", "authors.csv"), author); err != nil {
			fmt.Printf("Error saving author: %v\n", err)
			return
		}
		data.Authors[email] = author
	}

	authors := []*Author{author}

	if strings.EqualFold(itemType, "Book") {
		book := &Book{
			Title:       title,
			ISBN:        isbn,
			Authors:     authors,
			Description: extra,
		}
		if err := AppendBook(filepath.Join("resources", "books.csv"), book); err != nil {
			fmt.Printf("Error saving book: %v\n", err)
			return
		}
		data.Books = append(data.Books, book)
	} else {
		mag := &Magazine{
			Title:       title,
			ISBN:        isbn,
			Authors:     authors,
			PublishedAt: extra,
		}
		if err := AppendMagazine(filepath.Join("resources", "magazines.csv"), mag); err != nil {
			fmt.Printf("Error saving magazine: %v\n", err)
			return
		}
		data.Magazines = append(data.Magazines, mag)
	}
	fmt.Println("Item added successfully!")
}
