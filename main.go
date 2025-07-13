package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"time"
)

// Author represents an author
type Author struct {
	Email     string
	FirstName string
	LastName  string
}

// Book represents a book
type Book struct {
	Title       string
	ISBN        string
	Authors     []Author
	Description string
}

// Magazine represents a magazine
type Magazine struct {
	Title       string
	ISBN        string
	Authors     []Author
	PublishedAt time.Time
}

// LibraryItem represents a generic library item for sorting
type LibraryItem struct {
	Type        string // "BOOK" or "MAGAZINE"
	Title       string
	ISBN        string
	Authors     []Author
	Description string
	PublishedAt *time.Time // Only for magazines
}

// Library holds all the library data
type Library struct {
	Authors   map[string]Author
	Books     []Book
	Magazines []Magazine
}

func main() {
	library, err := loadLibrary()
	if err != nil {
		log.Fatalf("Error loading library: %v", err)
	}

	fmt.Println("=== Welcome to the Library System ===")
	fmt.Println()

	for {
		showMenu()
		choice := getUserInput("Enter your choice: ")

		switch choice {
		case "1":
			displayBooks(library.Books)
		case "2":
			displayMagazines(library.Magazines)
		case "3":
			displayAllItems(library.Books, library.Magazines)
		case "4":
			searchByISBN(library)
		case "5":
			searchByAuthorEmail(library)
		case "6":
			displayAllItemsSorted(library.Books, library.Magazines)
		case "7":
			addNewBook(library)
		case "8":
			addNewMagazine(library)
		case "9":
			startWebServer(library)
		case "10":
			fmt.Println("Thank you for using the Library System!")
			return
		default:
			fmt.Println("Invalid choice. Please try again.")
		}
		fmt.Println()
	}
}

func showMenu() {
	fmt.Println("=== Library Menu ===")
	fmt.Println("1. Display all books")
	fmt.Println("2. Display all magazines")
	fmt.Println("3. Display all items")
	fmt.Println("4. Search by ISBN")
	fmt.Println("5. Search by author email")
	fmt.Println("6. Display all items sorted by title")
	fmt.Println("7. Add new book")
	fmt.Println("8. Add new magazine")
	fmt.Println("9. Start web interface")
	fmt.Println("10. Exit")
	fmt.Println()
}

func getUserInput(prompt string) string {
	fmt.Print(prompt)
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	return strings.TrimSpace(scanner.Text())
}

func searchByISBN(library *Library) {
	isbn := getUserInput("Enter ISBN to search for: ")

	if isbn == "" {
		fmt.Println("ISBN cannot be empty.")
		return
	}

	fmt.Printf("\n=== SEARCH RESULTS FOR ISBN: %s ===\n", isbn)

	found := false

	// Search in books
	for _, book := range library.Books {
		if strings.EqualFold(book.ISBN, isbn) {
			fmt.Println("\n--- BOOK FOUND ---")
			displaySingleBook(book)
			found = true
		}
	}

	// Search in magazines
	for _, magazine := range library.Magazines {
		if strings.EqualFold(magazine.ISBN, isbn) {
			fmt.Println("\n--- MAGAZINE FOUND ---")
			displaySingleMagazine(magazine)
			found = true
		}
	}

	if !found {
		fmt.Printf("No book or magazine found with ISBN: %s\n", isbn)
	}
}

func displaySingleBook(book Book) {
	fmt.Printf("Title: %s\n", book.Title)
	fmt.Printf("ISBN: %s\n", book.ISBN)
	fmt.Printf("Authors: %s\n", formatAuthors(book.Authors))
	fmt.Printf("Description: %s\n", book.Description)
	fmt.Println(strings.Repeat("-", 60))
}

func displaySingleMagazine(magazine Magazine) {
	fmt.Printf("Title: %s\n", magazine.Title)
	fmt.Printf("ISBN: %s\n", magazine.ISBN)
	fmt.Printf("Authors: %s\n", formatAuthors(magazine.Authors))
	if !magazine.PublishedAt.IsZero() {
		fmt.Printf("Published: %s\n", magazine.PublishedAt.Format("02.01.2006"))
	}
	fmt.Println(strings.Repeat("-", 60))
}

func loadLibrary() (*Library, error) {
	library := &Library{
		Authors: make(map[string]Author),
	}

	// Load authors first
	err := loadAuthors(library)
	if err != nil {
		return nil, fmt.Errorf("failed to load authors: %w", err)
	}

	// Load books
	err = loadBooks(library)
	if err != nil {
		return nil, fmt.Errorf("failed to load books: %w", err)
	}

	// Load magazines
	err = loadMagazines(library)
	if err != nil {
		return nil, fmt.Errorf("failed to load magazines: %w", err)
	}

	return library, nil
}

func loadAuthors(library *Library) error {
	file, err := os.Open("resources/authors.csv")
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = ';'

	records, err := reader.ReadAll()
	if err != nil {
		return err
	}

	// Skip header row
	for i := 1; i < len(records); i++ {
		record := records[i]
		if len(record) >= 3 {
			author := Author{
				Email:     record[0],
				FirstName: record[1],
				LastName:  record[2],
			}
			library.Authors[author.Email] = author
		}
	}

	return nil
}

func loadBooks(library *Library) error {
	file, err := os.Open("resources/books.csv")
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = ';'

	records, err := reader.ReadAll()
	if err != nil {
		return err
	}

	// Skip header row
	for i := 1; i < len(records); i++ {
		record := records[i]
		if len(record) >= 4 {
			book := Book{
				Title:       record[0],
				ISBN:        record[1],
				Authors:     parseAuthors(record[2], library.Authors),
				Description: record[3],
			}
			library.Books = append(library.Books, book)
		}
	}

	return nil
}

func loadMagazines(library *Library) error {
	file, err := os.Open("resources/magazines.csv")
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = ';'

	records, err := reader.ReadAll()
	if err != nil {
		return err
	}

	// Skip header row
	for i := 1; i < len(records); i++ {
		record := records[i]
		if len(record) >= 4 {
			publishedAt, err := time.Parse("02.01.2006", record[3])
			if err != nil {
				log.Printf("Warning: Could not parse date %s for magazine %s", record[3], record[0])
				publishedAt = time.Time{}
			}

			magazine := Magazine{
				Title:       record[0],
				ISBN:        record[1],
				Authors:     parseAuthors(record[2], library.Authors),
				PublishedAt: publishedAt,
			}
			library.Magazines = append(library.Magazines, magazine)
		}
	}

	return nil
}

func parseAuthors(authorEmails string, authorMap map[string]Author) []Author {
	var authors []Author
	emails := strings.Split(authorEmails, ",")

	for _, email := range emails {
		email = strings.TrimSpace(email)
		if author, exists := authorMap[email]; exists {
			authors = append(authors, author)
		}
	}

	return authors
}

func displayBooks(books []Book) {
	fmt.Println("=== BOOKS ===")
	fmt.Printf("Total books: %d\n\n", len(books))

	for i, book := range books {
		fmt.Printf("Book #%d\n", i+1)
		fmt.Printf("Title: %s\n", book.Title)
		fmt.Printf("ISBN: %s\n", book.ISBN)
		fmt.Printf("Authors: %s\n", formatAuthors(book.Authors))
		fmt.Printf("Description: %s\n", truncateDescription(book.Description, 150))
		fmt.Println(strings.Repeat("-", 80))
	}
}

func displayMagazines(magazines []Magazine) {
	fmt.Println("=== MAGAZINES ===")
	fmt.Printf("Total magazines: %d\n\n", len(magazines))

	for i, magazine := range magazines {
		fmt.Printf("Magazine #%d\n", i+1)
		fmt.Printf("Title: %s\n", magazine.Title)
		fmt.Printf("ISBN: %s\n", magazine.ISBN)
		fmt.Printf("Authors: %s\n", formatAuthors(magazine.Authors))
		if !magazine.PublishedAt.IsZero() {
			fmt.Printf("Published: %s\n", magazine.PublishedAt.Format("02.01.2006"))
		}
		fmt.Println(strings.Repeat("-", 80))
	}
}

func displayAllItems(books []Book, magazines []Magazine) {
	fmt.Printf("=== ALL LIBRARY ITEMS ===\n")
	fmt.Printf("Total items: %d (Books: %d, Magazines: %d)\n\n", len(books)+len(magazines), len(books), len(magazines))

	itemCount := 1

	fmt.Println("--- BOOKS ---")
	for _, book := range books {
		fmt.Printf("%d. [BOOK] %s (ISBN: %s)\n", itemCount, book.Title, book.ISBN)
		fmt.Printf("   Authors: %s\n", formatAuthors(book.Authors))
		itemCount++
	}

	fmt.Println("\n--- MAGAZINES ---")
	for _, magazine := range magazines {
		fmt.Printf("%d. [MAGAZINE] %s (ISBN: %s)\n", itemCount, magazine.Title, magazine.ISBN)
		fmt.Printf("   Authors: %s", formatAuthors(magazine.Authors))
		if !magazine.PublishedAt.IsZero() {
			fmt.Printf(" - Published: %s", magazine.PublishedAt.Format("02.01.2006"))
		}
		fmt.Println()
		itemCount++
	}
}

func displayAllItemsSorted(books []Book, magazines []Magazine) {
	fmt.Printf("=== ALL LIBRARY ITEMS (SORTED BY TITLE) ===\n")
	fmt.Printf("Total items: %d (Books: %d, Magazines: %d)\n\n", len(books)+len(magazines), len(books), len(magazines))

	// Convert all items to LibraryItem slice
	var allItems []LibraryItem

	// Add books
	for _, book := range books {
		item := LibraryItem{
			Type:        "BOOK",
			Title:       book.Title,
			ISBN:        book.ISBN,
			Authors:     book.Authors,
			Description: book.Description,
			PublishedAt: nil,
		}
		allItems = append(allItems, item)
	}

	// Add magazines
	for _, magazine := range magazines {
		var publishedAt *time.Time
		if !magazine.PublishedAt.IsZero() {
			publishedAt = &magazine.PublishedAt
		}

		item := LibraryItem{
			Type:        "MAGAZINE",
			Title:       magazine.Title,
			ISBN:        magazine.ISBN,
			Authors:     magazine.Authors,
			Description: "",
			PublishedAt: publishedAt,
		}
		allItems = append(allItems, item)
	}

	// Sort by title (case-insensitive)
	sort.Slice(allItems, func(i, j int) bool {
		return strings.ToLower(allItems[i].Title) < strings.ToLower(allItems[j].Title)
	})

	// Display sorted items with full details
	for i, item := range allItems {
		fmt.Printf("%d. [%s] %s\n", i+1, item.Type, item.Title)
		fmt.Printf("   ISBN: %s\n", item.ISBN)
		fmt.Printf("   Authors: %s\n", formatAuthors(item.Authors))

		if item.Type == "BOOK" && item.Description != "" {
			fmt.Printf("   Description: %s\n", truncateDescription(item.Description, 150))
		}

		if item.Type == "MAGAZINE" && item.PublishedAt != nil {
			fmt.Printf("   Published: %s\n", item.PublishedAt.Format("02.01.2006"))
		}

		fmt.Println(strings.Repeat("-", 80))
	}
}

func formatAuthors(authors []Author) string {
	if len(authors) == 0 {
		return "Unknown"
	}

	var authorNames []string
	for _, author := range authors {
		authorNames = append(authorNames, fmt.Sprintf("%s %s", author.FirstName, author.LastName))
	}

	return strings.Join(authorNames, ", ")
}

func truncateDescription(description string, maxLength int) string {
	if len(description) <= maxLength {
		return description
	}
	return description[:maxLength] + "..."
}

func searchByAuthorEmail(library *Library) {
	email := getUserInput("Enter author email to search for: ")

	if email == "" {
		fmt.Println("Author email cannot be empty.")
		return
	}

	fmt.Printf("\n=== SEARCH RESULTS FOR AUTHOR EMAIL: %s ===\n", email)

	found := false

	// Get author info first
	author, authorExists := library.Authors[email]
	if !authorExists {
		fmt.Printf("No author found with email: %s\n", email)
		return
	}

	fmt.Printf("Author: %s %s (%s)\n\n", author.FirstName, author.LastName, author.Email)

	// Search in books
	var booksFound []Book
	for _, book := range library.Books {
		for _, bookAuthor := range book.Authors {
			if strings.EqualFold(bookAuthor.Email, email) {
				booksFound = append(booksFound, book)
				found = true
				break
			}
		}
	}

	// Search in magazines
	var magazinesFound []Magazine
	for _, magazine := range library.Magazines {
		for _, magazineAuthor := range magazine.Authors {
			if strings.EqualFold(magazineAuthor.Email, email) {
				magazinesFound = append(magazinesFound, magazine)
				found = true
				break
			}
		}
	}

	if !found {
		fmt.Printf("No publications found for author: %s %s\n", author.FirstName, author.LastName)
		return
	}

	// Display results
	if len(booksFound) > 0 {
		fmt.Printf("--- BOOKS BY %s %s ---\n", author.FirstName, author.LastName)
		for i, book := range booksFound {
			fmt.Printf("%d. %s (ISBN: %s)\n", i+1, book.Title, book.ISBN)
			fmt.Printf("   Description: %s\n", truncateDescription(book.Description, 100))
			fmt.Println()
		}
	}

	if len(magazinesFound) > 0 {
		fmt.Printf("--- MAGAZINES BY %s %s ---\n", author.FirstName, author.LastName)
		for i, magazine := range magazinesFound {
			fmt.Printf("%d. %s (ISBN: %s)\n", i+1, magazine.Title, magazine.ISBN)
			if !magazine.PublishedAt.IsZero() {
				fmt.Printf("   Published: %s\n", magazine.PublishedAt.Format("02.01.2006"))
			}
			fmt.Println()
		}
	}

	totalPublications := len(booksFound) + len(magazinesFound)
	fmt.Printf("Total publications by %s %s: %d (Books: %d, Magazines: %d)\n",
		author.FirstName, author.LastName, totalPublications, len(booksFound), len(magazinesFound))
}

func addNewBook(library *Library) {
	fmt.Println("\n=== ADD NEW BOOK ===")

	// Get book details
	title := getUserInput("Enter book title: ")
	if title == "" {
		fmt.Println("Title cannot be empty.")
		return
	}

	isbn := getUserInput("Enter ISBN: ")
	if isbn == "" {
		fmt.Println("ISBN cannot be empty.")
		return
	}

	// Check if ISBN already exists
	if isbnExists(isbn, library) {
		fmt.Printf("A book or magazine with ISBN %s already exists.\n", isbn)
		return
	}

	description := getUserInput("Enter description: ")

	// Get authors
	authors := getAuthorsForItem(library)
	if len(authors) == 0 {
		fmt.Println("At least one author is required.")
		return
	}

	// Create new book
	newBook := Book{
		Title:       title,
		ISBN:        isbn,
		Authors:     authors,
		Description: description,
	}

	// Add to library
	library.Books = append(library.Books, newBook)

	// Save to CSV
	err := saveBooksToCSV(library.Books)
	if err != nil {
		fmt.Printf("Error saving book to CSV: %v\n", err)
		return
	}

	fmt.Printf("Book '%s' successfully added to the library!\n", title)
}

func addNewMagazine(library *Library) {
	fmt.Println("\n=== ADD NEW MAGAZINE ===")

	// Get magazine details
	title := getUserInput("Enter magazine title: ")
	if title == "" {
		fmt.Println("Title cannot be empty.")
		return
	}

	isbn := getUserInput("Enter ISBN: ")
	if isbn == "" {
		fmt.Println("ISBN cannot be empty.")
		return
	}

	// Check if ISBN already exists
	if isbnExists(isbn, library) {
		fmt.Printf("A book or magazine with ISBN %s already exists.\n", isbn)
		return
	}

	// Get publication date
	dateStr := getUserInput("Enter publication date (DD.MM.YYYY): ")
	publishedAt, err := time.Parse("02.01.2006", dateStr)
	if err != nil {
		fmt.Printf("Invalid date format. Please use DD.MM.YYYY format.\n")
		return
	}

	// Get authors
	authors := getAuthorsForItem(library)
	if len(authors) == 0 {
		fmt.Println("At least one author is required.")
		return
	}

	// Create new magazine
	newMagazine := Magazine{
		Title:       title,
		ISBN:        isbn,
		Authors:     authors,
		PublishedAt: publishedAt,
	}

	// Add to library
	library.Magazines = append(library.Magazines, newMagazine)

	// Save to CSV
	err = saveMagazinesToCSV(library.Magazines)
	if err != nil {
		fmt.Printf("Error saving magazine to CSV: %v\n", err)
		return
	}

	fmt.Printf("Magazine '%s' successfully added to the library!\n", title)
}

func getAuthorsForItem(library *Library) []Author {
	var authors []Author

	fmt.Println("\n--- AUTHORS ---")
	fmt.Println("You can add multiple authors. Press Enter with empty input when done.")

	for {
		email := getUserInput(fmt.Sprintf("Enter author email (author #%d, or press Enter if done): ", len(authors)+1))
		if email == "" {
			break
		}

		// Check if author already exists
		if existingAuthor, exists := library.Authors[email]; exists {
			fmt.Printf("Using existing author: %s %s (%s)\n", existingAuthor.FirstName, existingAuthor.LastName, existingAuthor.Email)
			authors = append(authors, existingAuthor)
		} else {
			// Create new author
			fmt.Printf("Author with email %s not found. Creating new author.\n", email)
			firstName := getUserInput("Enter first name: ")
			if firstName == "" {
				fmt.Println("First name cannot be empty.")
				continue
			}

			lastName := getUserInput("Enter last name: ")
			if lastName == "" {
				fmt.Println("Last name cannot be empty.")
				continue
			}

			newAuthor := Author{
				Email:     email,
				FirstName: firstName,
				LastName:  lastName,
			}

			// Add to library and save to CSV
			library.Authors[email] = newAuthor
			err := saveAuthorsToCSV(library.Authors)
			if err != nil {
				fmt.Printf("Error saving author to CSV: %v\n", err)
				continue
			}

			authors = append(authors, newAuthor)
			fmt.Printf("New author '%s %s' created and added.\n", firstName, lastName)
		}
	}

	return authors
}

func isbnExists(isbn string, library *Library) bool {
	// Check in books
	for _, book := range library.Books {
		if strings.EqualFold(book.ISBN, isbn) {
			return true
		}
	}

	// Check in magazines
	for _, magazine := range library.Magazines {
		if strings.EqualFold(magazine.ISBN, isbn) {
			return true
		}
	}

	return false
}

func saveAuthorsToCSV(authors map[string]Author) error {
	file, err := os.Create("resources/authors.csv")
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	writer.Comma = ';'
	defer writer.Flush()

	// Write header
	err = writer.Write([]string{"email", "firstname", "lastname"})
	if err != nil {
		return err
	}

	// Write authors
	for _, author := range authors {
		record := []string{author.Email, author.FirstName, author.LastName}
		err = writer.Write(record)
		if err != nil {
			return err
		}
	}

	return nil
}

func saveBooksToCSV(books []Book) error {
	file, err := os.Create("resources/books.csv")
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	writer.Comma = ';'
	defer writer.Flush()

	// Write header
	err = writer.Write([]string{"title", "isbn", "authors", "description"})
	if err != nil {
		return err
	}

	// Write books
	for _, book := range books {
		var authorEmails []string
		for _, author := range book.Authors {
			authorEmails = append(authorEmails, author.Email)
		}
		authorsStr := strings.Join(authorEmails, ",")

		record := []string{book.Title, book.ISBN, authorsStr, book.Description}
		err = writer.Write(record)
		if err != nil {
			return err
		}
	}

	return nil
}

func saveMagazinesToCSV(magazines []Magazine) error {
	file, err := os.Create("resources/magazines.csv")
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	writer.Comma = ';'
	defer writer.Flush()

	// Write header
	err = writer.Write([]string{"title", "isbn", "authors", "publishedAt"})
	if err != nil {
		return err
	}

	// Write magazines
	for _, magazine := range magazines {
		var authorEmails []string
		for _, author := range magazine.Authors {
			authorEmails = append(authorEmails, author.Email)
		}
		authorsStr := strings.Join(authorEmails, ",")

		var dateStr string
		if !magazine.PublishedAt.IsZero() {
			dateStr = magazine.PublishedAt.Format("02.01.2006")
		}

		record := []string{magazine.Title, magazine.ISBN, authorsStr, dateStr}
		err = writer.Write(record)
		if err != nil {
			return err
		}
	}

	return nil
}
