package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
)

// Author represents a library author
type Author struct {
	Email     string
	FirstName string
	LastName  string
}

// Book represents a library book
type Book struct {
	Title        string
	ISBN         string
	AuthorEmails []string
	Description  string
}

// Magazine represents a library magazine
type Magazine struct {
	Title        string
	ISBN         string
	AuthorEmails []string
	PublishedAt  string
}

// Library holds all the data
type Library struct {
	Authors   map[string]Author
	Books     []Book
	Magazines []Magazine
}

// SearchResult represents a search result item
type SearchResult struct {
	Type string // "Book" or "Magazine"
	Book *Book
	Magazine *Magazine
}

func main() {
	webMode := flag.Bool("web", false, "Start web server instead of CLI")
	port := flag.Int("port", 8080, "Port for web server")
	flag.Parse()

	library, err := loadLibrary()
	if err != nil {
		fmt.Printf("Error loading library: %v\n", err)
		os.Exit(1)
	}

	if *webMode {
		if err := StartWebServer(library, *port); err != nil {
			fmt.Printf("Error starting web server: %v\n", err)
			os.Exit(1)
		}
	} else {
		runUI(library)
	}
}

func welcomeMessage() string {
	return "Welcome to the Library Application!"
}

// loadLibrary loads all data from CSV files
func loadLibrary() (*Library, error) {
	library := &Library{
		Authors:   make(map[string]Author),
		Books:     []Book{},
		Magazines: []Magazine{},
	}

	// Load authors
	authorsFile, err := os.Open("resources/authors.csv")
	if err != nil {
		return nil, fmt.Errorf("failed to open authors file: %w", err)
	}
	defer authorsFile.Close()

	authors, err := parseAuthors(authorsFile)
	if err != nil {
		return nil, fmt.Errorf("failed to parse authors: %w", err)
	}
	library.Authors = authors

	// Load books
	booksFile, err := os.Open("resources/books.csv")
	if err != nil {
		return nil, fmt.Errorf("failed to open books file: %w", err)
	}
	defer booksFile.Close()

	books, err := parseBooks(booksFile)
	if err != nil {
		return nil, fmt.Errorf("failed to parse books: %w", err)
	}
	library.Books = books

	// Load magazines
	magazinesFile, err := os.Open("resources/magazines.csv")
	if err != nil {
		return nil, fmt.Errorf("failed to open magazines file: %w", err)
	}
	defer magazinesFile.Close()

	magazines, err := parseMagazines(magazinesFile)
	if err != nil {
		return nil, fmt.Errorf("failed to parse magazines: %w", err)
	}
	library.Magazines = magazines

	return library, nil
}

// parseAuthors parses authors from a reader
func parseAuthors(r io.Reader) (map[string]Author, error) {
	reader := csv.NewReader(r)
	reader.Comma = ';'

	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	authors := make(map[string]Author)
	for i, record := range records {
		if i == 0 { // Skip header
			continue
		}
		if len(record) >= 3 {
			author := Author{
				Email:     record[0],
				FirstName: record[1],
				LastName:  record[2],
			}
			authors[author.Email] = author
		}
	}

	return authors, nil
}

// parseBooks parses books from a reader
func parseBooks(r io.Reader) ([]Book, error) {
	reader := csv.NewReader(r)
	reader.Comma = ';'

	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var books []Book
	for i, record := range records {
		if i == 0 { // Skip header
			continue
		}
		if len(record) >= 4 {
			authorEmails := strings.Split(record[2], ",")
			book := Book{
				Title:        record[0],
				ISBN:         record[1],
				AuthorEmails: authorEmails,
				Description:  record[3],
			}
			books = append(books, book)
		}
	}

	return books, nil
}

// parseMagazines parses magazines from a reader
func parseMagazines(r io.Reader) ([]Magazine, error) {
	reader := csv.NewReader(r)
	reader.Comma = ';'

	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var magazines []Magazine
	for i, record := range records {
		if i == 0 { // Skip header
			continue
		}
		if len(record) >= 4 {
			authorEmails := strings.Split(record[2], ",")
			magazine := Magazine{
				Title:        record[0],
				ISBN:         record[1],
				AuthorEmails: authorEmails,
				PublishedAt:  record[3],
			}
			magazines = append(magazines, magazine)
		}
	}

	return magazines, nil
}

// findByISBN searches for books and magazines by ISBN
func findByISBN(library *Library, isbn string) []SearchResult {
	var results []SearchResult

	for i := range library.Books {
		if library.Books[i].ISBN == isbn {
			results = append(results, SearchResult{
				Type: "Book",
				Book: &library.Books[i],
			})
		}
	}

	for i := range library.Magazines {
		if library.Magazines[i].ISBN == isbn {
			results = append(results, SearchResult{
				Type:     "Magazine",
				Magazine: &library.Magazines[i],
			})
		}
	}

	return results
}

// findByAuthorEmail searches for books and magazines by author email
func findByAuthorEmail(library *Library, email string) []SearchResult {
	var results []SearchResult

	for i := range library.Books {
		for _, authorEmail := range library.Books[i].AuthorEmails {
			if strings.TrimSpace(authorEmail) == email {
				results = append(results, SearchResult{
					Type: "Book",
					Book: &library.Books[i],
				})
				break
			}
		}
	}

	for i := range library.Magazines {
		for _, authorEmail := range library.Magazines[i].AuthorEmails {
			if strings.TrimSpace(authorEmail) == email {
				results = append(results, SearchResult{
					Type:     "Magazine",
					Magazine: &library.Magazines[i],
				})
				break
			}
		}
	}

	return results
}

// getAllItemsSortedByTitle returns all books and magazines sorted by title
// ascending=true for A-Z, ascending=false for Z-A
func getAllItemsSortedByTitle(library *Library, ascending bool) []LibraryItem {
	var items []LibraryItem

	for _, book := range library.Books {
		items = append(items, LibraryItem{
			Title:       book.Title,
			ISBN:        book.ISBN,
			Authors:     formatAuthors(library, book.AuthorEmails),
			ItemType:    "Book",
			Description: book.Description,
		})
	}

	for _, magazine := range library.Magazines {
		items = append(items, LibraryItem{
			Title:       magazine.Title,
			ISBN:        magazine.ISBN,
			Authors:     formatAuthors(library, magazine.AuthorEmails),
			ItemType:    "Magazine",
			PublishedAt: magazine.PublishedAt,
		})
	}

	sort.Slice(items, func(i, j int) bool {
		if ascending {
			return strings.ToLower(items[i].Title) < strings.ToLower(items[j].Title)
		}
		return strings.ToLower(items[i].Title) > strings.ToLower(items[j].Title)
	})

	return items
}

// runUI runs the main user interface loop
func runUI(library *Library) {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		printMenu()
		fmt.Print("Enter your choice: ")

		if !scanner.Scan() {
			break
		}

		choice := strings.TrimSpace(scanner.Text())

		switch choice {
		case "1":
			displayAllBooks(library)
		case "2":
			displayAllMagazines(library)
		case "3":
			displayAll(library)
		case "4":
			searchByISBN(library, scanner)
		case "5":
			searchByAuthorEmail(library, scanner)
		case "6":
			displayAllSortedByTitle(library, scanner)
		case "7":
			addNewBook(library, scanner)
		case "8":
			addNewMagazine(library, scanner)
		case "9":
			fmt.Println("\nGoodbye!")
			return
		default:
			fmt.Println("\nInvalid choice. Please try again.")
		}
	}
}

// printMenu displays the main menu
func printMenu() {
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println(welcomeMessage())
	fmt.Println(strings.Repeat("=", 60))
	fmt.Println("\n1. Display all Books")
	fmt.Println("2. Display all Magazines")
	fmt.Println("3. Display all Books and Magazines")
	fmt.Println("4. Search by ISBN")
	fmt.Println("5. Search by Author Email")
	fmt.Println("6. Display all sorted by Title")
	fmt.Println("7. Add a new Book")
	fmt.Println("8. Add a new Magazine")
	fmt.Println("9. Exit")
	fmt.Println()
}

// displayAllBooks shows all books with details
func displayAllBooks(library *Library) {
	fmt.Println("\n" + strings.Repeat("-", 60))
	fmt.Println("                      BOOKS")
	fmt.Println(strings.Repeat("-", 60))

	if len(library.Books) == 0 {
		fmt.Println("No books found.")
		return
	}

	for i, book := range library.Books {
		fmt.Printf("\n[%d] %s\n", i+1, book.Title)
		fmt.Printf("    ISBN: %s\n", book.ISBN)
		fmt.Printf("    Authors: %s\n", formatAuthors(library, book.AuthorEmails))
		fmt.Printf("    Description: %s\n", truncateString(book.Description, 100))
	}
}

// displayAllMagazines shows all magazines with details
func displayAllMagazines(library *Library) {
	fmt.Println("\n" + strings.Repeat("-", 60))
	fmt.Println("                    MAGAZINES")
	fmt.Println(strings.Repeat("-", 60))

	if len(library.Magazines) == 0 {
		fmt.Println("No magazines found.")
		return
	}

	for i, magazine := range library.Magazines {
		fmt.Printf("\n[%d] %s\n", i+1, magazine.Title)
		fmt.Printf("    ISBN: %s\n", magazine.ISBN)
		fmt.Printf("    Authors: %s\n", formatAuthors(library, magazine.AuthorEmails))
		fmt.Printf("    Published: %s\n", magazine.PublishedAt)
	}
}

// displayAll shows all books and magazines
func displayAll(library *Library) {
	displayAllBooks(library)
	displayAllMagazines(library)
}

// formatAuthors converts author emails to formatted names
func formatAuthors(library *Library, emails []string) string {
	var names []string
	for _, email := range emails {
		email = strings.TrimSpace(email)
		if author, exists := library.Authors[email]; exists {
			names = append(names, fmt.Sprintf("%s %s", author.FirstName, author.LastName))
		} else {
			names = append(names, email)
		}
	}
	return strings.Join(names, ", ")
}

// truncateString truncates a string to a maximum length
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}

// searchByISBN prompts for ISBN and searches for matching book or magazine
func searchByISBN(library *Library, scanner *bufio.Scanner) {
	fmt.Print("\nEnter ISBN to search: ")
	if !scanner.Scan() {
		return
	}

	isbn := strings.TrimSpace(scanner.Text())
	if isbn == "" {
		fmt.Println("No ISBN entered.")
		return
	}

	results := findByISBN(library, isbn)

	if len(results) == 0 {
		fmt.Printf("\nNo book or magazine found with ISBN: %s\n", isbn)
		return
	}

	for _, result := range results {
		if result.Type == "Book" {
			fmt.Println("\n" + strings.Repeat("-", 60))
			fmt.Println("                  BOOK FOUND")
			fmt.Println(strings.Repeat("-", 60))
			fmt.Printf("\nTitle: %s\n", result.Book.Title)
			fmt.Printf("ISBN: %s\n", result.Book.ISBN)
			fmt.Printf("Authors: %s\n", formatAuthors(library, result.Book.AuthorEmails))
			fmt.Printf("Description: %s\n", result.Book.Description)
		} else {
			fmt.Println("\n" + strings.Repeat("-", 60))
			fmt.Println("                MAGAZINE FOUND")
			fmt.Println(strings.Repeat("-", 60))
			fmt.Printf("\nTitle: %s\n", result.Magazine.Title)
			fmt.Printf("ISBN: %s\n", result.Magazine.ISBN)
			fmt.Printf("Authors: %s\n", formatAuthors(library, result.Magazine.AuthorEmails))
			fmt.Printf("Published: %s\n", result.Magazine.PublishedAt)
		}
	}
}

// searchByAuthorEmail prompts for author email and searches for matching books and magazines
func searchByAuthorEmail(library *Library, scanner *bufio.Scanner) {
	fmt.Print("\nEnter author email to search: ")
	if !scanner.Scan() {
		return
	}

	email := strings.TrimSpace(scanner.Text())
	if email == "" {
		fmt.Println("No email entered.")
		return
	}

	// Check if author exists
	author, authorExists := library.Authors[email]
	if authorExists {
		fmt.Printf("\nAuthor: %s %s (%s)\n", author.FirstName, author.LastName, author.Email)
	}

	results := findByAuthorEmail(library, email)

	if len(results) == 0 {
		fmt.Printf("\nNo books or magazines found for author: %s\n", email)
		return
	}

	// Separate books and magazines
	var foundBooks []Book
	var foundMagazines []Magazine
	for _, result := range results {
		if result.Type == "Book" {
			foundBooks = append(foundBooks, *result.Book)
		} else {
			foundMagazines = append(foundMagazines, *result.Magazine)
		}
	}

	if len(foundBooks) > 0 {
		fmt.Println("\n" + strings.Repeat("-", 60))
		fmt.Printf("              BOOKS BY AUTHOR (%d found)\n", len(foundBooks))
		fmt.Println(strings.Repeat("-", 60))
		for i, book := range foundBooks {
			fmt.Printf("\n[%d] %s\n", i+1, book.Title)
			fmt.Printf("    ISBN: %s\n", book.ISBN)
			fmt.Printf("    Authors: %s\n", formatAuthors(library, book.AuthorEmails))
			fmt.Printf("    Description: %s\n", truncateString(book.Description, 100))
		}
	}

	if len(foundMagazines) > 0 {
		fmt.Println("\n" + strings.Repeat("-", 60))
		fmt.Printf("            MAGAZINES BY AUTHOR (%d found)\n", len(foundMagazines))
		fmt.Println(strings.Repeat("-", 60))
		for i, magazine := range foundMagazines {
			fmt.Printf("\n[%d] %s\n", i+1, magazine.Title)
			fmt.Printf("    ISBN: %s\n", magazine.ISBN)
			fmt.Printf("    Authors: %s\n", formatAuthors(library, magazine.AuthorEmails))
			fmt.Printf("    Published: %s\n", magazine.PublishedAt)
		}
	}
}

// LibraryItem represents a generic item (book or magazine) for sorting
type LibraryItem struct {
	Title        string
	ISBN         string
	Authors      string
	ItemType     string // "Book" or "Magazine"
	Description  string // For books
	PublishedAt  string // For magazines
}

// displayAllSortedByTitle shows all books and magazines sorted by title
func displayAllSortedByTitle(library *Library, scanner *bufio.Scanner) {
	fmt.Println("\n" + strings.Repeat("-", 60))
	fmt.Println("          ALL ITEMS SORTED BY TITLE")
	fmt.Println(strings.Repeat("-", 60))

	fmt.Println("\nSort direction:")
	fmt.Println("1. Ascending (A-Z)")
	fmt.Println("2. Descending (Z-A)")
	fmt.Print("Enter choice [1]: ")

	ascending := true
	if scanner.Scan() {
		choice := strings.TrimSpace(scanner.Text())
		if choice == "2" {
			ascending = false
		}
	}

	items := getAllItemsSortedByTitle(library, ascending)

	if ascending {
		fmt.Println("\nShowing items sorted A-Z:")
	} else {
		fmt.Println("\nShowing items sorted Z-A:")
	}

	if len(items) == 0 {
		fmt.Println("No items found.")
		return
	}

	for i, item := range items {
		fmt.Printf("\n[%d] %s [%s]\n", i+1, item.Title, item.ItemType)
		fmt.Printf("    ISBN: %s\n", item.ISBN)
		fmt.Printf("    Authors: %s\n", item.Authors)
		if item.ItemType == "Book" {
			fmt.Printf("    Description: %s\n", truncateString(item.Description, 100))
		} else {
			fmt.Printf("    Published: %s\n", item.PublishedAt)
		}
	}
}

// addNewBook prompts the user to add a new book with authors
func addNewBook(library *Library, scanner *bufio.Scanner) {
	fmt.Println("\n" + strings.Repeat("-", 60))
	fmt.Println("                  ADD NEW BOOK")
	fmt.Println(strings.Repeat("-", 60))

	// Get book title
	fmt.Print("\nEnter book title: ")
	if !scanner.Scan() {
		return
	}
	title := strings.TrimSpace(scanner.Text())
	if title == "" {
		fmt.Println("Title cannot be empty. Operation cancelled.")
		return
	}

	// Get ISBN
	fmt.Print("Enter ISBN: ")
	if !scanner.Scan() {
		return
	}
	isbn := strings.TrimSpace(scanner.Text())
	if isbn == "" {
		fmt.Println("ISBN cannot be empty. Operation cancelled.")
		return
	}

	// Check if ISBN already exists
	if len(findByISBN(library, isbn)) > 0 {
		fmt.Printf("A book or magazine with ISBN %s already exists. Operation cancelled.\n", isbn)
		return
	}

	// Get description
	fmt.Print("Enter description: ")
	if !scanner.Scan() {
		return
	}
	description := strings.TrimSpace(scanner.Text())

	// Get authors
	authorEmails, newAuthors := collectAuthors(library, scanner)
	if len(authorEmails) == 0 {
		fmt.Println("At least one author is required. Operation cancelled.")
		return
	}

	// Create the book
	book := Book{
		Title:        title,
		ISBN:         isbn,
		AuthorEmails: authorEmails,
		Description:  description,
	}

	// Add new authors to library and save
	for _, author := range newAuthors {
		library.Authors[author.Email] = author
	}
	if len(newAuthors) > 0 {
		if err := saveAuthors(library); err != nil {
			fmt.Printf("Error saving authors: %v\n", err)
			return
		}
	}

	// Add book to library and save
	library.Books = append(library.Books, book)
	if err := saveBooks(library); err != nil {
		fmt.Printf("Error saving book: %v\n", err)
		return
	}

	fmt.Println("\n✓ Book added successfully!")
	fmt.Printf("  Title: %s\n", book.Title)
	fmt.Printf("  ISBN: %s\n", book.ISBN)
	fmt.Printf("  Authors: %s\n", formatAuthors(library, book.AuthorEmails))
}

// addNewMagazine prompts the user to add a new magazine with authors
func addNewMagazine(library *Library, scanner *bufio.Scanner) {
	fmt.Println("\n" + strings.Repeat("-", 60))
	fmt.Println("                ADD NEW MAGAZINE")
	fmt.Println(strings.Repeat("-", 60))

	// Get magazine title
	fmt.Print("\nEnter magazine title: ")
	if !scanner.Scan() {
		return
	}
	title := strings.TrimSpace(scanner.Text())
	if title == "" {
		fmt.Println("Title cannot be empty. Operation cancelled.")
		return
	}

	// Get ISBN
	fmt.Print("Enter ISBN: ")
	if !scanner.Scan() {
		return
	}
	isbn := strings.TrimSpace(scanner.Text())
	if isbn == "" {
		fmt.Println("ISBN cannot be empty. Operation cancelled.")
		return
	}

	// Check if ISBN already exists
	if len(findByISBN(library, isbn)) > 0 {
		fmt.Printf("A book or magazine with ISBN %s already exists. Operation cancelled.\n", isbn)
		return
	}

	// Get publication date
	fmt.Print("Enter publication date (e.g., 01.01.2024): ")
	if !scanner.Scan() {
		return
	}
	publishedAt := strings.TrimSpace(scanner.Text())
	if publishedAt == "" {
		fmt.Println("Publication date cannot be empty. Operation cancelled.")
		return
	}

	// Get authors
	authorEmails, newAuthors := collectAuthors(library, scanner)
	if len(authorEmails) == 0 {
		fmt.Println("At least one author is required. Operation cancelled.")
		return
	}

	// Create the magazine
	magazine := Magazine{
		Title:        title,
		ISBN:         isbn,
		AuthorEmails: authorEmails,
		PublishedAt:  publishedAt,
	}

	// Add new authors to library and save
	for _, author := range newAuthors {
		library.Authors[author.Email] = author
	}
	if len(newAuthors) > 0 {
		if err := saveAuthors(library); err != nil {
			fmt.Printf("Error saving authors: %v\n", err)
			return
		}
	}

	// Add magazine to library and save
	library.Magazines = append(library.Magazines, magazine)
	if err := saveMagazines(library); err != nil {
		fmt.Printf("Error saving magazine: %v\n", err)
		return
	}

	fmt.Println("\n✓ Magazine added successfully!")
	fmt.Printf("  Title: %s\n", magazine.Title)
	fmt.Printf("  ISBN: %s\n", magazine.ISBN)
	fmt.Printf("  Authors: %s\n", formatAuthors(library, magazine.AuthorEmails))
	fmt.Printf("  Published: %s\n", magazine.PublishedAt)
}

// collectAuthors collects author emails and creates new authors if needed
func collectAuthors(library *Library, scanner *bufio.Scanner) ([]string, []Author) {
	var authorEmails []string
	var newAuthors []Author

	fmt.Println("\nEnter author email(s). Type 'done' when finished.")

	for {
		fmt.Print("Author email (or 'done'): ")
		if !scanner.Scan() {
			break
		}
		input := strings.TrimSpace(scanner.Text())

		if strings.ToLower(input) == "done" {
			break
		}

		if input == "" {
			continue
		}

		// Check if author already exists
		if _, exists := library.Authors[input]; exists {
			authorEmails = append(authorEmails, input)
			fmt.Printf("  → Found existing author: %s\n", formatAuthors(library, []string{input}))
		} else {
			// Author doesn't exist, prompt to create
			fmt.Printf("  Author '%s' not found. Create new author?\n", input)
			fmt.Print("  First name: ")
			if !scanner.Scan() {
				break
			}
			firstName := strings.TrimSpace(scanner.Text())
			if firstName == "" {
				fmt.Println("  First name required. Skipping this author.")
				continue
			}

			fmt.Print("  Last name: ")
			if !scanner.Scan() {
				break
			}
			lastName := strings.TrimSpace(scanner.Text())
			if lastName == "" {
				fmt.Println("  Last name required. Skipping this author.")
				continue
			}

			newAuthor := Author{
				Email:     input,
				FirstName: firstName,
				LastName:  lastName,
			}
			newAuthors = append(newAuthors, newAuthor)
			authorEmails = append(authorEmails, input)
			fmt.Printf("  → New author will be created: %s %s (%s)\n", firstName, lastName, input)
		}
	}

	return authorEmails, newAuthors
}

// saveAuthors writes all authors to the CSV file
func saveAuthors(library *Library) error {
	file, err := os.Create("resources/authors.csv")
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	writer.Comma = ';'
	defer writer.Flush()

	// Write header
	if err := writer.Write([]string{"email", "firstname", "lastname"}); err != nil {
		return err
	}

	// Write authors
	for _, author := range library.Authors {
		record := []string{author.Email, author.FirstName, author.LastName}
		if err := writer.Write(record); err != nil {
			return err
		}
	}

	return nil
}

// saveBooks writes all books to the CSV file
func saveBooks(library *Library) error {
	file, err := os.Create("resources/books.csv")
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	writer.Comma = ';'
	defer writer.Flush()

	// Write header
	if err := writer.Write([]string{"title", "isbn", "authors", "description"}); err != nil {
		return err
	}

	// Write books
	for _, book := range library.Books {
		authors := strings.Join(book.AuthorEmails, ",")
		record := []string{book.Title, book.ISBN, authors, book.Description}
		if err := writer.Write(record); err != nil {
			return err
		}
	}

	return nil
}

// saveMagazines writes all magazines to the CSV file
func saveMagazines(library *Library) error {
	file, err := os.Create("resources/magazines.csv")
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	writer.Comma = ';'
	defer writer.Flush()

	// Write header
	if err := writer.Write([]string{"title", "isbn", "authors", "publishedAt"}); err != nil {
		return err
	}

	// Write magazines
	for _, magazine := range library.Magazines {
		authors := strings.Join(magazine.AuthorEmails, ",")
		record := []string{magazine.Title, magazine.ISBN, authors, magazine.PublishedAt}
		if err := writer.Write(record); err != nil {
			return err
		}
	}

	return nil
}
