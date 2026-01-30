package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"sort"
	"strings"
	"sync"
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

type Library struct {
	Books     []Book
	Magazines []Magazine
	Authors   map[string]Author
	mu        sync.RWMutex
}

func NewLibrary() *Library {
	return &Library{
		Books:     []Book{},
		Magazines: []Magazine{},
		Authors:   make(map[string]Author),
	}
}

func (l *Library) LoadAuthors(filePath string) error {
	file, err := os.Open(filePath)
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
			continue
		}
		if len(record) < 3 {
			continue
		}

		author := Author{
			Email:     record[0],
			FirstName: record[1],
			LastName:  record[2],
		}
		l.Authors[author.Email] = author
	}

	return nil
}

func (l *Library) LoadBooks(filePath string) error {
	file, err := os.Open(filePath)
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
			continue
		}
		if len(record) < 4 {
			continue
		}

		var authors []Author
		authorEmails := strings.Split(record[2], ",")
		for _, email := range authorEmails {
			email = strings.TrimSpace(email)
			if author, exists := l.Authors[email]; exists {
				authors = append(authors, author)
			}
		}

		book := Book{
			Title:       record[0],
			ISBN:        record[1],
			Authors:     authors,
			Description: record[3],
		}
		l.Books = append(l.Books, book)
	}

	return nil
}

func (l *Library) LoadMagazines(filePath string) error {
	file, err := os.Open(filePath)
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
			continue
		}
		if len(record) < 4 {
			continue
		}

		var authors []Author
		authorEmails := strings.Split(record[2], ",")
		for _, email := range authorEmails {
			email = strings.TrimSpace(email)
			if author, exists := l.Authors[email]; exists {
				authors = append(authors, author)
			}
		}

		magazine := Magazine{
			Title:       record[0],
			ISBN:        record[1],
			Authors:     authors,
			PublishedAt: record[3],
		}
		l.Magazines = append(l.Magazines, magazine)
	}

	return nil
}

func (l *Library) Load() error {
	if err := l.LoadAuthors("resources/authors.csv"); err != nil {
		return fmt.Errorf("failed to load authors: %w", err)
	}

	if err := l.LoadBooks("resources/books.csv"); err != nil {
		return fmt.Errorf("failed to load books: %w", err)
	}

	if err := l.LoadMagazines("resources/magazines.csv"); err != nil {
		return fmt.Errorf("failed to load magazines: %w", err)
	}

	return nil
}

func (l *Library) SearchBooksByISBN(isbn string) []Book {
	var results []Book
	for _, book := range l.Books {
		if strings.Contains(strings.ToLower(book.ISBN), strings.ToLower(isbn)) {
			results = append(results, book)
		}
	}
	return results
}

func (l *Library) SearchMagazinesByISBN(isbn string) []Magazine {
	var results []Magazine
	for _, magazine := range l.Magazines {
		if strings.Contains(strings.ToLower(magazine.ISBN), strings.ToLower(isbn)) {
			results = append(results, magazine)
		}
	}
	return results
}

func (l *Library) SearchBooksByAuthorEmail(email string) []Book {
	var results []Book
	email = strings.ToLower(strings.TrimSpace(email))
	for _, book := range l.Books {
		for _, author := range book.Authors {
			if strings.EqualFold(author.Email, email) {
				results = append(results, book)
				break
			}
		}
	}
	return results
}

func (l *Library) SearchMagazinesByAuthorEmail(email string) []Magazine {
	var results []Magazine
	email = strings.ToLower(strings.TrimSpace(email))
	for _, magazine := range l.Magazines {
		for _, author := range magazine.Authors {
			if strings.EqualFold(author.Email, email) {
				results = append(results, magazine)
				break
			}
		}
	}
	return results
}

// Item represents a combined book or magazine for sorting
type Item struct {
	Type        string      // "book" or "magazine"
	Title       string
	ISBN        string
	Authors     []Author
	Description string // for books
	PublishedAt string  // for magazines
}

func (l *Library) GetAllItemsSortedByTitle() []Item {
	var items []Item

	// Add all books
	for _, book := range l.Books {
		items = append(items, Item{
			Type:        "book",
			Title:       book.Title,
			ISBN:        book.ISBN,
			Authors:     book.Authors,
			Description: book.Description,
		})
	}

	// Add all magazines
	for _, magazine := range l.Magazines {
		items = append(items, Item{
			Type:        "magazine",
			Title:       magazine.Title,
			ISBN:        magazine.ISBN,
			Authors:     magazine.Authors,
			PublishedAt: magazine.PublishedAt,
		})
	}

	// Sort by title
	sort.Slice(items, func(i, j int) bool {
		return strings.ToLower(items[i].Title) < strings.ToLower(items[j].Title)
	})

	return items
}

func (l *Library) GetAllItemsSortedByTitleDescending() []Item {
	var items []Item

	// Add all books
	for _, book := range l.Books {
		items = append(items, Item{
			Type:        "book",
			Title:       book.Title,
			ISBN:        book.ISBN,
			Authors:     book.Authors,
			Description: book.Description,
		})
	}

	// Add all magazines
	for _, magazine := range l.Magazines {
		items = append(items, Item{
			Type:        "magazine",
			Title:       magazine.Title,
			ISBN:        magazine.ISBN,
			Authors:     magazine.Authors,
			PublishedAt: magazine.PublishedAt,
		})
	}

	// Sort by title descending
	sort.Slice(items, func(i, j int) bool {
		return strings.ToLower(items[i].Title) > strings.ToLower(items[j].Title)
	})

	return items
}

func (l *Library) AddAuthor(author Author) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.Authors[author.Email] = author
	return l.saveAuthorsToCSV()
}

func (l *Library) AddBook(book Book) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.Books = append(l.Books, book)
	return l.saveBooksToCSV()
}

func (l *Library) AddMagazine(magazine Magazine) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.Magazines = append(l.Magazines, magazine)
	return l.saveMagazinesToCSV()
}

func (l *Library) saveAuthorsToCSV() error {
	file, err := os.Create("resources/authors.csv")
	if err != nil {
		return fmt.Errorf("failed to create authors.csv: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	writer.Comma = ';'
	defer writer.Flush()

	// Write header
	if err := writer.Write([]string{"Email", "FirstName", "LastName"}); err != nil {
		return fmt.Errorf("failed to write header: %w", err)
	}

	// Write authors
	for _, author := range l.Authors {
		if err := writer.Write([]string{author.Email, author.FirstName, author.LastName}); err != nil {
			return fmt.Errorf("failed to write author: %w", err)
		}
	}

	return nil
}

func (l *Library) saveBooksToCSV() error {
	file, err := os.Create("resources/books.csv")
	if err != nil {
		return fmt.Errorf("failed to create books.csv: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	writer.Comma = ';'
	defer writer.Flush()

	// Write header
	if err := writer.Write([]string{"Title", "ISBN", "Authors", "Description"}); err != nil {
		return fmt.Errorf("failed to write header: %w", err)
	}

	// Write books
	for _, book := range l.Books {
		authorEmails := make([]string, len(book.Authors))
		for i, author := range book.Authors {
			authorEmails[i] = author.Email
		}
		if err := writer.Write([]string{book.Title, book.ISBN, strings.Join(authorEmails, ","), book.Description}); err != nil {
			return fmt.Errorf("failed to write book: %w", err)
		}
	}

	return nil
}

func (l *Library) saveMagazinesToCSV() error {
	file, err := os.Create("resources/magazines.csv")
	if err != nil {
		return fmt.Errorf("failed to create magazines.csv: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	writer.Comma = ';'
	defer writer.Flush()

	// Write header
	if err := writer.Write([]string{"Title", "ISBN", "Authors", "PublishedAt"}); err != nil {
		return fmt.Errorf("failed to write header: %w", err)
	}

	// Write magazines
	for _, magazine := range l.Magazines {
		authorEmails := make([]string, len(magazine.Authors))
		for i, author := range magazine.Authors {
			authorEmails[i] = author.Email
		}
		if err := writer.Write([]string{magazine.Title, magazine.ISBN, strings.Join(authorEmails, ","), magazine.PublishedAt}); err != nil {
			return fmt.Errorf("failed to write magazine: %w", err)
		}
	}

	return nil
}
