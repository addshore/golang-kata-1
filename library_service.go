package main

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

// LibraryService handles business logic for library operations
type LibraryService struct {
	library *Library
}

// NewLibraryService creates a new library service
func NewLibraryService(library *Library) *LibraryService {
	return &LibraryService{library: library}
}

// SearchByISBN searches for books and magazines by ISBN
func (s *LibraryService) SearchByISBN(searchTerm string) ([]Book, []Magazine) {
	var books []Book
	var magazines []Magazine

	searchTerm = strings.ToLower(strings.TrimSpace(searchTerm))
	if searchTerm == "" {
		return books, magazines
	}

	// Search books
	for _, book := range s.library.Books {
		if strings.Contains(strings.ToLower(book.ISBN), searchTerm) {
			books = append(books, book)
		}
	}

	// Search magazines
	for _, magazine := range s.library.Magazines {
		if strings.Contains(strings.ToLower(magazine.ISBN), searchTerm) {
			magazines = append(magazines, magazine)
		}
	}

	return books, magazines
}

// SearchByAuthorEmail searches for books and magazines by author email
func (s *LibraryService) SearchByAuthorEmail(searchTerm string) ([]Book, []Magazine) {
	var books []Book
	var magazines []Magazine

	searchTerm = strings.ToLower(strings.TrimSpace(searchTerm))
	if searchTerm == "" {
		return books, magazines
	}

	// Search books by author email
	for _, book := range s.library.Books {
		for _, author := range book.Authors {
			if strings.Contains(strings.ToLower(author.Email), searchTerm) {
				books = append(books, book)
				break // Avoid duplicates if multiple authors match
			}
		}
	}

	// Search magazines by author email
	for _, magazine := range s.library.Magazines {
		for _, author := range magazine.Authors {
			if strings.Contains(strings.ToLower(author.Email), searchTerm) {
				magazines = append(magazines, magazine)
				break // Avoid duplicates if multiple authors match
			}
		}
	}

	return books, magazines
}

// GetAllItemsSorted returns all books and magazines sorted by title
func (s *LibraryService) GetAllItemsSorted(ascending bool) []LibraryItem {
	var allItems []LibraryItem

	// Add books
	for _, book := range s.library.Books {
		allItems = append(allItems, LibraryItem{
			Title:       book.Title,
			ISBN:        book.ISBN,
			AuthorNames: s.formatAuthors(book.Authors),
			Type:        "BOOK",
			Description: book.Description,
		})
	}

	// Add magazines
	for _, magazine := range s.library.Magazines {
		allItems = append(allItems, LibraryItem{
			Title:       magazine.Title,
			ISBN:        magazine.ISBN,
			AuthorNames: s.formatAuthors(magazine.Authors),
			Type:        "MAGAZINE",
			PublishedAt: magazine.PublishedAt.Format("02.01.2006"),
		})
	}

	// Sort by title
	sort.Slice(allItems, func(i, j int) bool {
		if ascending {
			return strings.ToLower(allItems[i].Title) < strings.ToLower(allItems[j].Title)
		}
		return strings.ToLower(allItems[i].Title) > strings.ToLower(allItems[j].Title)
	})

	return allItems
}

// GetBooks returns all books
func (s *LibraryService) GetBooks() []Book {
	return s.library.Books
}

// GetMagazines returns all magazines
func (s *LibraryService) GetMagazines() []Magazine {
	return s.library.Magazines
}

// GetAuthors returns all authors
func (s *LibraryService) GetAuthors() map[string]Author {
	return s.library.Authors
}

// GetLibraryStats returns statistics about the library
func (s *LibraryService) GetLibraryStats() (int, int, int) {
	return len(s.library.Books), len(s.library.Magazines), len(s.library.Authors)
}

// formatAuthors formats a slice of authors into a readable string
func (s *LibraryService) formatAuthors(authors []Author) string {
	if len(authors) == 0 {
		return "Unknown"
	}

	var names []string
	for _, author := range authors {
		names = append(names, author.FirstName+" "+author.LastName)
	}
	return strings.Join(names, ", ")
}

// FormatAuthors is a public method for formatting authors (used by UI)
func (s *LibraryService) FormatAuthors(authors []Author) string {
	return s.formatAuthors(authors)
}

// AddBook adds a new book to the library and saves to CSV
func (s *LibraryService) AddBook(title, isbn, description string, authorEmails []string) error {
	// Validate required fields
	if strings.TrimSpace(title) == "" {
		return fmt.Errorf("title is required")
	}
	if strings.TrimSpace(isbn) == "" {
		return fmt.Errorf("ISBN is required")
	}

	// Check if ISBN already exists
	for _, book := range s.library.Books {
		if book.ISBN == isbn {
			return fmt.Errorf("book with ISBN %s already exists", isbn)
		}
	}
	for _, magazine := range s.library.Magazines {
		if magazine.ISBN == isbn {
			return fmt.Errorf("magazine with ISBN %s already exists", isbn)
		}
	}

	// Create book
	book := Book{
		Title:       strings.TrimSpace(title),
		ISBN:        strings.TrimSpace(isbn),
		Description: strings.TrimSpace(description),
		Authors:     []Author{},
	}

	// Add authors
	for _, email := range authorEmails {
		email = strings.TrimSpace(email)
		if email != "" {
			if author, exists := s.library.Authors[email]; exists {
				book.Authors = append(book.Authors, author)
			} else {
				return fmt.Errorf("author with email %s not found", email)
			}
		}
	}

	// Add to library
	s.library.Books = append(s.library.Books, book)

	// Save to CSV
	return SaveLibrary(s.library)
}

// AddMagazine adds a new magazine to the library and saves to CSV
func (s *LibraryService) AddMagazine(title, isbn string, publishedAt time.Time, authorEmails []string) error {
	// Validate required fields
	if strings.TrimSpace(title) == "" {
		return fmt.Errorf("title is required")
	}
	if strings.TrimSpace(isbn) == "" {
		return fmt.Errorf("ISBN is required")
	}

	// Check if ISBN already exists
	for _, book := range s.library.Books {
		if book.ISBN == isbn {
			return fmt.Errorf("book with ISBN %s already exists", isbn)
		}
	}
	for _, magazine := range s.library.Magazines {
		if magazine.ISBN == isbn {
			return fmt.Errorf("magazine with ISBN %s already exists", isbn)
		}
	}

	// Create magazine
	magazine := Magazine{
		Title:       strings.TrimSpace(title),
		ISBN:        strings.TrimSpace(isbn),
		PublishedAt: publishedAt,
		Authors:     []Author{},
	}

	// Add authors
	for _, email := range authorEmails {
		email = strings.TrimSpace(email)
		if email != "" {
			if author, exists := s.library.Authors[email]; exists {
				magazine.Authors = append(magazine.Authors, author)
			} else {
				return fmt.Errorf("author with email %s not found", email)
			}
		}
	}

	// Add to library
	s.library.Magazines = append(s.library.Magazines, magazine)

	// Save to CSV
	return SaveLibrary(s.library)
}

// AddAuthor adds a new author to the library and saves to CSV
func (s *LibraryService) AddAuthor(email, firstName, lastName string) error {
	// Validate required fields
	email = strings.TrimSpace(email)
	firstName = strings.TrimSpace(firstName)
	lastName = strings.TrimSpace(lastName)

	if email == "" {
		return fmt.Errorf("email is required")
	}
	if firstName == "" {
		return fmt.Errorf("first name is required")
	}
	if lastName == "" {
		return fmt.Errorf("last name is required")
	}

	// Check if author already exists
	if _, exists := s.library.Authors[email]; exists {
		return fmt.Errorf("author with email %s already exists", email)
	}

	// Create author
	author := Author{
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
	}

	// Add to library
	s.library.Authors[email] = author

	// Save to CSV
	return SaveLibrary(s.library)
}

// GetAuthorsList returns a slice of all authors for easier iteration
func (s *LibraryService) GetAuthorsList() []Author {
	var authors []Author
	for _, author := range s.library.Authors {
		authors = append(authors, author)
	}

	// Sort by last name, then first name
	sort.Slice(authors, func(i, j int) bool {
		if authors[i].LastName == authors[j].LastName {
			return authors[i].FirstName < authors[j].FirstName
		}
		return authors[i].LastName < authors[j].LastName
	})

	return authors
}
