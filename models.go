package main

import (
	"encoding/csv"
	"fmt"
	"os"
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
