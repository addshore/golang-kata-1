package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"
	"sort"
	"strings"
)

type Book struct {
	Title  string
	Author string
	ISBN   string
	Genre  string
}

type Magazine struct {
	Title       string
	Publisher   string
	IssueNumber string
	PublishedAt string
}

func parseBooks(filePath string) []Book {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		return nil
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = ';' // Set semicolon as the delimiter
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Printf("Error reading CSV: %v\n", err)
		return nil
	}

	var books []Book
	for i, record := range records {
		if i == 0 { // Skip header row
			continue
		}
		if len(record) < 4 {
			fmt.Printf("Skipping invalid record on line %d: %v\n", i+1, record)
			continue
		}

		isbn := strings.TrimSpace(record[1])

		books = append(books, Book{
			Title:  strings.TrimSpace(record[0]),
			Author: strings.TrimSpace(record[2]),
			ISBN:   isbn,
			Genre:  strings.TrimSpace(record[3]),
		})
	}

	return books
}

func parseMagazines(filePath string) []Magazine {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		return nil
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = ';' // Set semicolon as the delimiter
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Printf("Error reading CSV: %v\n", err)
		return nil
	}

	var magazines []Magazine
	for i, record := range records {
		if i == 0 { // Skip header row
			continue
		}
		if len(record) < 4 {
			fmt.Printf("Skipping invalid record on line %d: %v\n", i+1, record)
			continue
		}

		magazines = append(magazines, Magazine{
			Title:       strings.TrimSpace(record[0]),
			Publisher:   strings.TrimSpace(record[2]),
			IssueNumber: strings.TrimSpace(record[1]),
			PublishedAt: strings.TrimSpace(record[3]),
		})
	}

	return magazines
}

func displayLibrary(books []Book, magazines []Magazine) {
	fmt.Println("========================")
	fmt.Println("       Library Books")
	fmt.Println("========================")
	for _, book := range books {
		fmt.Printf("Title: %s\nAuthor: %s\nISBN: %s\nDescription: %s\n\n", book.Title, book.Author, book.ISBN, book.Genre)
	}

	fmt.Println("========================")
	fmt.Println("     Library Magazines")
	fmt.Println("========================")
	for _, magazine := range magazines {
		fmt.Printf("Title: %s\nPublisher: %s\nISBN: %s\nPublished At: %s\n\n", magazine.Title, magazine.Publisher, magazine.IssueNumber, magazine.PublishedAt)
	}
}

func searchByISBN(books []Book, magazines []Magazine, isbn string) []string {
	var results []string

	for _, book := range books {
		if book.ISBN == isbn {
			results = append(results, fmt.Sprintf("Book - Title: %s\nAuthor: %s\nISBN: %s\nDescription: %s\n", book.Title, book.Author, book.ISBN, book.Genre))
		}
	}

	for _, magazine := range magazines {
		if magazine.IssueNumber == isbn {
			results = append(results, fmt.Sprintf("Magazine - Title: %s\nPublisher: %s\nISBN: %s\nPublished At: %s\n", magazine.Title, magazine.Publisher, magazine.IssueNumber, magazine.PublishedAt))
		}
	}

	return results
}

func searchByAuthorEmail(books []Book, magazines []Magazine, email string) []string {
	var results []string

	for _, book := range books {
		if strings.Contains(book.Author, email) {
			results = append(results, fmt.Sprintf("Book - Title: %s\nAuthor: %s\nISBN: %s\nDescription: %s\n", book.Title, book.Author, book.ISBN, book.Genre))
		}
	}

	for _, magazine := range magazines {
		if strings.Contains(magazine.Publisher, email) {
			results = append(results, fmt.Sprintf("Magazine - Title: %s\nPublisher: %s\nISBN: %s\nPublished At: %s\n", magazine.Title, magazine.Publisher, magazine.IssueNumber, magazine.PublishedAt))
		}
	}

	return results
}

func sortByTitle(books []Book, magazines []Magazine) []string {
	type LibraryItem struct {
		Title   string
		Details string
	}

	var libraryItems []LibraryItem

	for _, book := range books {
		libraryItems = append(libraryItems, LibraryItem{
			Title:   book.Title,
			Details: fmt.Sprintf("Book - Title: %s\nAuthor: %s\nISBN: %s\nDescription: %s\n", book.Title, book.Author, book.ISBN, book.Genre),
		})
	}

	for _, magazine := range magazines {
		libraryItems = append(libraryItems, LibraryItem{
			Title:   magazine.Title,
			Details: fmt.Sprintf("Magazine - Title: %s\nPublisher: %s\nISBN: %s\nPublished At: %s\n", magazine.Title, magazine.Publisher, magazine.IssueNumber, magazine.PublishedAt),
		})
	}

	sort.Slice(libraryItems, func(i, j int) bool {
		return libraryItems[i].Title < libraryItems[j].Title
	})

	var results []string
	for _, item := range libraryItems {
		results = append(results, item.Details)
	}

	return results
}

func printResults(results []string) {
	if len(results) == 0 {
		fmt.Println("No results found.")
	} else {
		for _, result := range results {
			fmt.Println(result)
		}
	}
}

func addToLibrary(booksFile, magazinesFile, authorsFile string, books *[]Book, magazines *[]Magazine) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("\nAdd to Library:")
	fmt.Print("Enter type (book/magazine): ")
	itemType, _ := reader.ReadString('\n')
	itemType = strings.TrimSpace(strings.ToLower(itemType))

	fmt.Print("Enter title: ")
	title, _ := reader.ReadString('\n')
	title = strings.TrimSpace(title)

	fmt.Print("Enter author email: ")
	authorEmail, _ := reader.ReadString('\n')
	authorEmail = strings.TrimSpace(authorEmail)

	fmt.Print("Enter ISBN: ")
	isbn, _ := reader.ReadString('\n')
	isbn = strings.TrimSpace(isbn)

	if itemType == "book" {
		fmt.Print("Enter genre: ")
		genre, _ := reader.ReadString('\n')
		genre = strings.TrimSpace(genre)

		file, err := os.OpenFile(booksFile, os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Printf("Error opening books file: %v\n", err)
			return
		}
		defer file.Close()

		_, err = file.WriteString(fmt.Sprintf("%s;%s;%s;%s\n", title, isbn, authorEmail, genre))
		if err != nil {
			fmt.Printf("Error writing to books file: %v\n", err)
			return
		}

		*books = parseBooks(booksFile)
	} else if itemType == "magazine" {
		fmt.Print("Enter published date (DD.MM.YYYY): ")
		publishedDate, _ := reader.ReadString('\n')
		publishedDate = strings.TrimSpace(publishedDate)

		file, err := os.OpenFile(magazinesFile, os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Printf("Error opening magazines file: %v\n", err)
			return
		}
		defer file.Close()

		_, err = file.WriteString(fmt.Sprintf("%s;%s;%s;%s\n", title, isbn, authorEmail, publishedDate))
		if err != nil {
			fmt.Printf("Error writing to magazines file: %v\n", err)
			return
		}

		*magazines = parseMagazines(magazinesFile)
	} else {
		fmt.Println("Invalid type. Please enter 'book' or 'magazine'.")
		return
	}

	authors := parseAuthors(authorsFile)
	authorExists := false
	for _, author := range authors {
		if author.Email == authorEmail {
			authorExists = true
			break
		}
	}

	if !authorExists {
		fmt.Print("Enter author's first name: ")
		firstName, _ := reader.ReadString('\n')
		firstName = strings.TrimSpace(firstName)

		fmt.Print("Enter author's last name: ")
		lastName, _ := reader.ReadString('\n')
		lastName = strings.TrimSpace(lastName)

		file, err := os.OpenFile(authorsFile, os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Printf("Error opening authors file: %v\n", err)
			return
		}
		defer file.Close()

		_, err = file.WriteString(fmt.Sprintf("%s;%s;%s\n", authorEmail, firstName, lastName))
		if err != nil {
			fmt.Printf("Error writing to authors file: %v\n", err)
			return
		}

		fmt.Println("Author added successfully!")
	} else {
		fmt.Println("Author already exists in the library.")
	}

	fmt.Println("Item added successfully!")
}

func parseAuthors(filePath string) []struct {
	Email     string
	FirstName string
	LastName  string
} {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		return nil
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = ';'
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Printf("Error reading CSV: %v\n", err)
		return nil
	}

	var authors []struct {
		Email     string
		FirstName string
		LastName  string
	}
	for i, record := range records {
		if i == 0 { // Skip header row
			continue
		}
		if len(record) < 3 {
			fmt.Printf("Skipping invalid record on line %d: %v\n", i+1, record)
			continue
		}

		authors = append(authors, struct {
			Email     string
			FirstName string
			LastName  string
		}{
			Email:     strings.TrimSpace(record[0]),
			FirstName: strings.TrimSpace(record[1]),
			LastName:  strings.TrimSpace(record[2]),
		})
	}

	return authors
}

func handleMenu(books []Book, magazines []Magazine) {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Println("\nLibrary Menu:")
		fmt.Println("1. List all books and magazines")
		fmt.Println("2. Search by ISBN")
		fmt.Println("3. Search by Author's Email")
		fmt.Println("4. Sort by Title")
		fmt.Println("5. Add to Library")
		fmt.Println("6. Exit")
		fmt.Print("Enter your choice: ")

		choice, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(choice)

		switch choice {
		case "1":
			displayLibrary(books, magazines)
		case "2":
			fmt.Print("Enter ISBN to search: ")
			isbn, _ := reader.ReadString('\n')
			isbn = strings.TrimSpace(isbn)
			printResults(searchByISBN(books, magazines, isbn))
		case "3":
			fmt.Print("Enter Author's Email to search: ")
			email, _ := reader.ReadString('\n')
			email = strings.TrimSpace(email)
			printResults(searchByAuthorEmail(books, magazines, email))
		case "4":
			printResults(sortByTitle(books, magazines))
		case "5":
			addToLibrary("resources/books.csv", "resources/magazines.csv", "resources/authors.csv", &books, &magazines)
		case "6":
			fmt.Println("Exiting...")
			return
		default:
			fmt.Println("Invalid choice. Please try again.")
		}
	}
}

func main() {
	books := parseBooks("resources/books.csv")
	magazines := parseMagazines("resources/magazines.csv")
	handleMenu(books, magazines)
}
