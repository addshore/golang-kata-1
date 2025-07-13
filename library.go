package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"sort"
	"strings"
)

type Author struct {
	Email     string
	FirstName string
	LastName  string
}

type Book struct {
	Title       string
	ISBN        string
	Authors     []Author
	Description string
}

type Magazine struct {
	Title       string
	ISBN        string
	Authors     []Author
	PublishedAt string
}

type LibraryItem struct {
	Title       string
	ISBN        string
	Authors     []Author
	Description string
	PublishedAt string
	Type        string
}

// Library represents the library system
type Library struct {
	Authors   map[string]Author
	Books     []Book
	Magazines []Magazine
}

// NewLibrary creates a new library instance
func NewLibrary() *Library {
	return &Library{
		Authors:   make(map[string]Author),
		Books:     []Book{},
		Magazines: []Magazine{},
	}
}

// LoadFromCSV loads data from CSV files
func (l *Library) LoadFromCSV(authorsFile, booksFile, magazinesFile string) error {
	if err := l.loadAuthors(authorsFile); err != nil {
		return fmt.Errorf("failed to load authors: %w", err)
	}
	
	if err := l.loadBooks(booksFile); err != nil {
		return fmt.Errorf("failed to load books: %w", err)
	}
	
	if err := l.loadMagazines(magazinesFile); err != nil {
		return fmt.Errorf("failed to load magazines: %w", err)
	}
	
	return nil
}

// loadAuthors loads authors from CSV file
func (l *Library) loadAuthors(filename string) error {
	file, err := os.Open(filename)
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

	for i, record := range records {
		if i == 0 {
			continue // Skip header
		}
		if len(record) >= 3 {
			l.Authors[record[0]] = Author{
				Email:     record[0],
				FirstName: record[1],
				LastName:  record[2],
			}
		}
	}
	return nil
}

// loadBooks loads books from CSV file
func (l *Library) loadBooks(filename string) error {
	file, err := os.Open(filename)
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

	for i, record := range records {
		if i == 0 {
			continue // Skip header
		}
		if len(record) >= 4 {
			authorEmails := strings.Split(record[2], ",")
			var bookAuthors []Author
			for _, email := range authorEmails {
				email = strings.TrimSpace(email)
				if author, exists := l.Authors[email]; exists {
					bookAuthors = append(bookAuthors, author)
				}
			}
			
			l.Books = append(l.Books, Book{
				Title:       record[0],
				ISBN:        record[1],
				Authors:     bookAuthors,
				Description: record[3],
			})
		}
	}
	return nil
}

// loadMagazines loads magazines from CSV file
func (l *Library) loadMagazines(filename string) error {
	file, err := os.Open(filename)
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

	for i, record := range records {
		if i == 0 {
			continue // Skip header
		}
		if len(record) >= 4 {
			authorEmails := strings.Split(record[2], ",")
			var magazineAuthors []Author
			for _, email := range authorEmails {
				email = strings.TrimSpace(email)
				if author, exists := l.Authors[email]; exists {
					magazineAuthors = append(magazineAuthors, author)
				}
			}
			
			l.Magazines = append(l.Magazines, Magazine{
				Title:       record[0],
				ISBN:        record[1],
				Authors:     magazineAuthors,
				PublishedAt: record[3],
			})
		}
	}
	return nil
}

// FindBookByISBN finds a book by ISBN
func (l *Library) FindBookByISBN(isbn string) *Book {
	for _, book := range l.Books {
		if book.ISBN == isbn {
			return &book
		}
	}
	return nil
}

// FindMagazineByISBN finds a magazine by ISBN
func (l *Library) FindMagazineByISBN(isbn string) *Magazine {
	for _, magazine := range l.Magazines {
		if magazine.ISBN == isbn {
			return &magazine
		}
	}
	return nil
}

// FindBooksByAuthorEmail finds books by author email
func (l *Library) FindBooksByAuthorEmail(email string) []Book {
	var foundBooks []Book
	for _, book := range l.Books {
		for _, author := range book.Authors {
			if author.Email == email {
				foundBooks = append(foundBooks, book)
				break
			}
		}
	}
	return foundBooks
}

// FindMagazinesByAuthorEmail finds magazines by author email
func (l *Library) FindMagazinesByAuthorEmail(email string) []Magazine {
	var foundMagazines []Magazine
	for _, magazine := range l.Magazines {
		for _, author := range magazine.Authors {
			if author.Email == email {
				foundMagazines = append(foundMagazines, magazine)
				break
			}
		}
	}
	return foundMagazines
}

// GetAllItemsSorted returns all items sorted by title
func (l *Library) GetAllItemsSorted(ascending bool) []LibraryItem {
	var allItems []LibraryItem
	
	for _, book := range l.Books {
		allItems = append(allItems, LibraryItem{
			Title:       book.Title,
			ISBN:        book.ISBN,
			Authors:     book.Authors,
			Description: book.Description,
			PublishedAt: "",
			Type:        "Book",
		})
	}
	
	for _, magazine := range l.Magazines {
		allItems = append(allItems, LibraryItem{
			Title:       magazine.Title,
			ISBN:        magazine.ISBN,
			Authors:     magazine.Authors,
			Description: "",
			PublishedAt: magazine.PublishedAt,
			Type:        "Magazine",
		})
	}
	
	if ascending {
		sort.Slice(allItems, func(i, j int) bool {
			return strings.ToLower(allItems[i].Title) < strings.ToLower(allItems[j].Title)
		})
	} else {
		sort.Slice(allItems, func(i, j int) bool {
			return strings.ToLower(allItems[i].Title) > strings.ToLower(allItems[j].Title)
		})
	}
	
	return allItems
}

// FormatAuthors formats a slice of authors as a string
func FormatAuthors(authors []Author) string {
	if len(authors) == 0 {
		return "Unknown"
	}
	
	var authorNames []string
	for _, author := range authors {
		authorNames = append(authorNames, fmt.Sprintf("%s %s", author.FirstName, author.LastName))
	}
	return strings.Join(authorNames, ", ")
}

// AddAuthor adds a new author to the library
func (l *Library) AddAuthor(author Author) {
	l.Authors[author.Email] = author
}

// AddBook adds a new book to the library
func (l *Library) AddBook(book Book) {
	l.Books = append(l.Books, book)
}

// AddMagazine adds a new magazine to the library
func (l *Library) AddMagazine(magazine Magazine) {
	l.Magazines = append(l.Magazines, magazine)
}

// SaveToCSV saves all library data to CSV files
func (l *Library) SaveToCSV(authorsFile, booksFile, magazinesFile string) error {
	if err := l.saveAuthorsToCSV(authorsFile); err != nil {
		return fmt.Errorf("failed to save authors: %w", err)
	}
	
	if err := l.saveBooksToCSV(booksFile); err != nil {
		return fmt.Errorf("failed to save books: %w", err)
	}
	
	if err := l.saveMagazinesToCSV(magazinesFile); err != nil {
		return fmt.Errorf("failed to save magazines: %w", err)
	}
	
	return nil
}

// saveAuthorsToCSV saves authors to CSV file
func (l *Library) saveAuthorsToCSV(filename string) error {
	file, err := os.Create(filename)
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
	for _, author := range l.Authors {
		record := []string{author.Email, author.FirstName, author.LastName}
		if err := writer.Write(record); err != nil {
			return err
		}
	}

	return nil
}

// saveBooksToCSV saves books to CSV file
func (l *Library) saveBooksToCSV(filename string) error {
	file, err := os.Create(filename)
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
	for _, book := range l.Books {
		var authorEmails []string
		for _, author := range book.Authors {
			authorEmails = append(authorEmails, author.Email)
		}
		authorsStr := strings.Join(authorEmails, ",")
		
		record := []string{book.Title, book.ISBN, authorsStr, book.Description}
		if err := writer.Write(record); err != nil {
			return err
		}
	}

	return nil
}

// saveMagazinesToCSV saves magazines to CSV file
func (l *Library) saveMagazinesToCSV(filename string) error {
	file, err := os.Create(filename)
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
	for _, magazine := range l.Magazines {
		var authorEmails []string
		for _, author := range magazine.Authors {
			authorEmails = append(authorEmails, author.Email)
		}
		authorsStr := strings.Join(authorEmails, ",")
		
		record := []string{magazine.Title, magazine.ISBN, authorsStr, magazine.PublishedAt}
		if err := writer.Write(record); err != nil {
			return err
		}
	}

	return nil
}

// GetOrCreateAuthor gets an existing author or creates a new one
func (l *Library) GetOrCreateAuthor(email, firstName, lastName string) Author {
	if author, exists := l.Authors[email]; exists {
		return author
	}
	
	author := Author{
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
	}
	l.AddAuthor(author)
	return author
}