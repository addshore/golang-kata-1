package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"sort"
	"strings"
)

// Author represents an author in the library
type Author struct {
	Email     string
	FirstName string
	LastName  string
}

// Book represents a book in the library
type Book struct {
	Title       string
	ISBN        string
	Authors     []string // List of author emails
	Description string
}

// Magazine represents a magazine in the library
type Magazine struct {
	Title       string
	ISBN        string
	Authors     []string // List of author emails
	PublishedAt string
}

// Library holds all library data
type Library struct {
	Authors   map[string]*Author // Keyed by email
	Books     []*Book
	Magazines []*Magazine
}

// LoadLibrary loads all data from CSV files
func LoadLibrary(resourcesDir string) (*Library, error) {
	lib := &Library{
		Authors: make(map[string]*Author),
	}

	// Load authors
	if err := lib.loadAuthors(resourcesDir + "/authors.csv"); err != nil {
		return nil, fmt.Errorf("failed to load authors: %w", err)
	}

	// Load books
	if err := lib.loadBooks(resourcesDir + "/books.csv"); err != nil {
		return nil, fmt.Errorf("failed to load books: %w", err)
	}

	// Load magazines
	if err := lib.loadMagazines(resourcesDir + "/magazines.csv"); err != nil {
		return nil, fmt.Errorf("failed to load magazines: %w", err)
	}

	return lib, nil
}

func (lib *Library) loadAuthors(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = ';'
	reader.LazyQuotes = true

	records, err := reader.ReadAll()
	if err != nil {
		return err
	}

	// Skip header (BOM + header row)
	for i, record := range records {
		if i == 0 {
			continue // Skip header
		}
		if len(record) >= 3 {
			author := &Author{
				Email:     strings.TrimSpace(record[0]),
				FirstName: strings.TrimSpace(record[1]),
				LastName:  strings.TrimSpace(record[2]),
			}
			lib.Authors[author.Email] = author
		}
	}

	return nil
}

func (lib *Library) loadBooks(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = ';'
	reader.LazyQuotes = true

	records, err := reader.ReadAll()
	if err != nil {
		return err
	}

	// Skip header
	for i, record := range records {
		if i == 0 {
			continue // Skip header
		}
		if len(record) >= 4 {
			authorsStr := strings.TrimSpace(record[2])
			authors := []string{}
			if authorsStr != "" {
				for _, email := range strings.Split(authorsStr, ",") {
					authors = append(authors, strings.TrimSpace(email))
				}
			}

			book := &Book{
				Title:       strings.TrimSpace(record[0]),
				ISBN:        strings.TrimSpace(record[1]),
				Authors:     authors,
				Description: strings.TrimSpace(record[3]),
			}
			lib.Books = append(lib.Books, book)
		}
	}

	return nil
}

func (lib *Library) loadMagazines(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = ';'
	reader.LazyQuotes = true

	records, err := reader.ReadAll()
	if err != nil {
		return err
	}

	// Skip header
	for i, record := range records {
		if i == 0 {
			continue // Skip header
		}
		if len(record) >= 4 {
			authorsStr := strings.TrimSpace(record[2])
			authors := []string{}
			if authorsStr != "" {
				for _, email := range strings.Split(authorsStr, ",") {
					authors = append(authors, strings.TrimSpace(email))
				}
			}

			magazine := &Magazine{
				Title:       strings.TrimSpace(record[0]),
				ISBN:        strings.TrimSpace(record[1]),
				Authors:     authors,
				PublishedAt: strings.TrimSpace(record[3]),
			}
			lib.Magazines = append(lib.Magazines, magazine)
		}
	}

	return nil
}

// GetAuthorNames returns formatted author names for a list of emails
func (lib *Library) GetAuthorNames(emails []string) []string {
	names := make([]string, 0, len(emails))
	for _, email := range emails {
		if author, ok := lib.Authors[email]; ok {
			names = append(names, fmt.Sprintf("%s %s", author.FirstName, author.LastName))
		}
	}
	return names
}

// SearchByISBN searches for books and magazines by ISBN
func (lib *Library) SearchByISBN(isbn string) ([]*Book, []*Magazine) {
	var books []*Book
	var magazines []*Magazine

	// Search books
	for _, book := range lib.Books {
		if book.ISBN == isbn {
			books = append(books, book)
		}
	}

	// Search magazines
	for _, magazine := range lib.Magazines {
		if magazine.ISBN == isbn {
			magazines = append(magazines, magazine)
		}
	}

	return books, magazines
}

// SearchByAuthorEmail searches for books and magazines by author email
func (lib *Library) SearchByAuthorEmail(email string) ([]*Book, []*Magazine) {
	var books []*Book
	var magazines []*Magazine

	// Search books
	for _, book := range lib.Books {
		for _, authorEmail := range book.Authors {
			if authorEmail == email {
				books = append(books, book)
				break
			}
		}
	}

	// Search magazines
	for _, magazine := range lib.Magazines {
		for _, authorEmail := range magazine.Authors {
			if authorEmail == email {
				magazines = append(magazines, magazine)
				break
			}
		}
	}

	return books, magazines
}

// Item represents either a book or magazine for sorting
type Item struct {
	Title    string
	ISBN     string
	Type     string // "book" or "magazine"
	Book     *Book
	Magazine *Magazine
}

// GetAllSortedByTitle returns all books and magazines sorted by title
func (lib *Library) GetAllSortedByTitle() []Item {
	var items []Item

	// Add all books
	for _, book := range lib.Books {
		items = append(items, Item{
			Title: book.Title,
			ISBN:  book.ISBN,
			Type:  "book",
			Book:  book,
		})
	}

	// Add all magazines
	for _, magazine := range lib.Magazines {
		items = append(items, Item{
			Title:    magazine.Title,
			ISBN:     magazine.ISBN,
			Type:     "magazine",
			Magazine: magazine,
		})
	}

	// Sort by title
	sort.Slice(items, func(i, j int) bool {
		return strings.ToLower(items[i].Title) < strings.ToLower(items[j].Title)
	})

	return items
}

// SaveToCSV saves all library data back to CSV files
func (lib *Library) SaveToCSV(resourcesDir string) error {
	// Save authors
	if err := lib.saveAuthors(resourcesDir + "/authors.csv"); err != nil {
		return fmt.Errorf("failed to save authors: %w", err)
	}

	// Save books
	if err := lib.saveBooks(resourcesDir + "/books.csv"); err != nil {
		return fmt.Errorf("failed to save books: %w", err)
	}

	// Save magazines
	if err := lib.saveMagazines(resourcesDir + "/magazines.csv"); err != nil {
		return fmt.Errorf("failed to save magazines: %w", err)
	}

	return nil
}

func (lib *Library) saveAuthors(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	writer.Comma = ';'
	defer writer.Flush()

	// Write header with BOM
	if err := writer.Write([]string{"\ufeffemail", "firstname", "lastname"}); err != nil {
		return err
	}

	// Write authors
	for _, author := range lib.Authors {
		if err := writer.Write([]string{author.Email, author.FirstName, author.LastName}); err != nil {
			return err
		}
	}

	return nil
}

func (lib *Library) saveBooks(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	writer.Comma = ';'
	defer writer.Flush()

	// Write header with BOM
	if err := writer.Write([]string{"\ufefftitle", "isbn", "authors", "description"}); err != nil {
		return err
	}

	// Write books
	for _, book := range lib.Books {
		authors := strings.Join(book.Authors, ",")
		if err := writer.Write([]string{book.Title, book.ISBN, authors, book.Description}); err != nil {
			return err
		}
	}

	return nil
}

func (lib *Library) saveMagazines(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	writer.Comma = ';'
	defer writer.Flush()

	// Write header with BOM
	if err := writer.Write([]string{"\ufefftitle", "isbn", "authors", "publishedAt"}); err != nil {
		return err
	}

	// Write magazines
	for _, magazine := range lib.Magazines {
		authors := strings.Join(magazine.Authors, ",")
		if err := writer.Write([]string{magazine.Title, magazine.ISBN, authors, magazine.PublishedAt}); err != nil {
			return err
		}
	}

	return nil
}

// AddBook adds a new book to the library
func (lib *Library) AddBook(title, isbn, description string, authorEmails []string) {
	book := &Book{
		Title:       title,
		ISBN:        isbn,
		Authors:     authorEmails,
		Description: description,
	}
	lib.Books = append(lib.Books, book)
}

// AddMagazine adds a new magazine to the library
func (lib *Library) AddMagazine(title, isbn, publishedAt string, authorEmails []string) {
	magazine := &Magazine{
		Title:       title,
		ISBN:        isbn,
		Authors:     authorEmails,
		PublishedAt: publishedAt,
	}
	lib.Magazines = append(lib.Magazines, magazine)
}

// AddAuthor adds a new author to the library
func (lib *Library) AddAuthor(email, firstName, lastName string) {
	if _, exists := lib.Authors[email]; !exists {
		author := &Author{
			Email:     email,
			FirstName: firstName,
			LastName:  lastName,
		}
		lib.Authors[email] = author
	}
}
