package main

import (
	"os"
	"strings"
	"testing"

	. "github.com/onsi/gomega"
)

func TestLoadLibrary(t *testing.T) {
	g := NewGomegaWithT(t)

	// Load library from test resources
	lib, err := LoadLibrary("./resources")
	g.Expect(err).To(BeNil())
	g.Expect(lib).ToNot(BeNil())

	// Verify authors are loaded
	g.Expect(len(lib.Authors)).To(Equal(6))
	g.Expect(lib.Authors["null-walter@echocat.org"]).ToNot(BeNil())
	g.Expect(lib.Authors["null-walter@echocat.org"].FirstName).To(Equal("Paul"))
	g.Expect(lib.Authors["null-walter@echocat.org"].LastName).To(Equal("Walter"))

	// Verify books are loaded
	g.Expect(len(lib.Books)).To(Equal(8))

	// Verify magazines are loaded
	g.Expect(len(lib.Magazines)).To(Equal(6))
}

func TestLoadAuthors(t *testing.T) {
	g := NewGomegaWithT(t)

	lib := &Library{
		Authors: make(map[string]*Author),
	}

	err := lib.loadAuthors("./resources/authors.csv")
	g.Expect(err).To(BeNil())
	g.Expect(len(lib.Authors)).To(Equal(6))

	// Check specific author
	author := lib.Authors["null-mueller@echocat.org"]
	g.Expect(author).ToNot(BeNil())
	g.Expect(author.FirstName).To(Equal("Max"))
	g.Expect(author.LastName).To(Equal("Müller"))
}

func TestLoadBooks(t *testing.T) {
	g := NewGomegaWithT(t)

	lib := &Library{
		Authors: make(map[string]*Author),
	}

	err := lib.loadBooks("./resources/books.csv")
	g.Expect(err).To(BeNil())
	g.Expect(len(lib.Books)).To(Equal(8))

	// Check a book with single author
	book1 := lib.Books[0]
	g.Expect(book1.Title).To(Equal("Ich helfe dir kochen. Das erfolgreiche Universalkochbuch mit großem Backteil"))
	g.Expect(book1.ISBN).To(Equal("5554-5545-4518"))
	g.Expect(len(book1.Authors)).To(Equal(1))
	g.Expect(book1.Authors[0]).To(Equal("null-walter@echocat.org"))

	// Check a book with multiple authors
	book2 := lib.Books[1]
	g.Expect(len(book2.Authors)).To(Equal(2))
}

func TestLoadMagazines(t *testing.T) {
	g := NewGomegaWithT(t)

	lib := &Library{
		Authors: make(map[string]*Author),
	}

	err := lib.loadMagazines("./resources/magazines.csv")
	g.Expect(err).To(BeNil())
	g.Expect(len(lib.Magazines)).To(Equal(6))

	// Check a magazine
	mag1 := lib.Magazines[0]
	g.Expect(mag1.Title).To(Equal("Beautiful cooking"))
	g.Expect(mag1.ISBN).To(Equal("5454-5587-3210"))
	g.Expect(mag1.PublishedAt).To(Equal("21.05.2011"))
}

func TestGetAuthorNames(t *testing.T) {
	g := NewGomegaWithT(t)

	lib := &Library{
		Authors: map[string]*Author{
			"test1@example.com": {Email: "test1@example.com", FirstName: "John", LastName: "Doe"},
			"test2@example.com": {Email: "test2@example.com", FirstName: "Jane", LastName: "Smith"},
		},
	}

	names := lib.GetAuthorNames([]string{"test1@example.com", "test2@example.com"})
	g.Expect(len(names)).To(Equal(2))
	g.Expect(names[0]).To(Equal("John Doe"))
	g.Expect(names[1]).To(Equal("Jane Smith"))
}

func TestSearchByISBN(t *testing.T) {
	g := NewGomegaWithT(t)

	lib, err := LoadLibrary("./resources")
	g.Expect(err).To(BeNil())

	// Search for a book
	books, magazines := lib.SearchByISBN("1024-5245-8584")
	g.Expect(len(books)).To(Equal(1))
	g.Expect(len(magazines)).To(Equal(0))
	g.Expect(books[0].Title).To(Equal("Genial italienisch"))

	// Search for a magazine
	books, magazines = lib.SearchByISBN("5454-5587-3210")
	g.Expect(len(books)).To(Equal(0))
	g.Expect(len(magazines)).To(Equal(1))
	g.Expect(magazines[0].Title).To(Equal("Beautiful cooking"))

	// Search for non-existent ISBN
	books, magazines = lib.SearchByISBN("0000-0000-0000")
	g.Expect(len(books)).To(Equal(0))
	g.Expect(len(magazines)).To(Equal(0))
}

func TestSearchByAuthorEmail(t *testing.T) {
	g := NewGomegaWithT(t)

	lib, err := LoadLibrary("./resources")
	g.Expect(err).To(BeNil())

	// Search for Paul Walter
	books, magazines := lib.SearchByAuthorEmail("null-walter@echocat.org")
	g.Expect(len(books)).To(BeNumerically(">", 0))
	g.Expect(len(magazines)).To(BeNumerically(">", 0))

	// Verify a known book by this author
	foundBook := false
	for _, book := range books {
		if book.ISBN == "5554-5545-4518" {
			foundBook = true
			break
		}
	}
	g.Expect(foundBook).To(BeTrue())

	// Search for non-existent email
	books, magazines = lib.SearchByAuthorEmail("nonexistent@example.com")
	g.Expect(len(books)).To(Equal(0))
	g.Expect(len(magazines)).To(Equal(0))
}

func TestGetAllSortedByTitle(t *testing.T) {
	g := NewGomegaWithT(t)

	lib, err := LoadLibrary("./resources")
	g.Expect(err).To(BeNil())

	items := lib.GetAllSortedByTitle()
	g.Expect(len(items)).To(Equal(14)) // 8 books + 6 magazines

	// Verify items are sorted by title (case-insensitive)
	for i := 1; i < len(items); i++ {
		prev := strings.ToLower(items[i-1].Title)
		curr := strings.ToLower(items[i].Title)
		g.Expect(prev <= curr).To(BeTrue(), "Items should be sorted: '%s' should come before '%s'", items[i-1].Title, items[i].Title)
	}

	// Check that both books and magazines are mixed
	hasBooks := false
	hasMagazines := false
	for _, item := range items {
		if item.Type == "book" {
			hasBooks = true
		}
		if item.Type == "magazine" {
			hasMagazines = true
		}
	}
	g.Expect(hasBooks).To(BeTrue())
	g.Expect(hasMagazines).To(BeTrue())
}

func TestLoadLibraryInvalidPath(t *testing.T) {
	g := NewGomegaWithT(t)

	_, err := LoadLibrary("/nonexistent/path")
	g.Expect(err).ToNot(BeNil())
}

func TestCreateTempCSVAndLoad(t *testing.T) {
	g := NewGomegaWithT(t)

	// Create a temporary directory
	tmpDir, err := os.MkdirTemp("", "library-test-*")
	g.Expect(err).To(BeNil())
	defer os.RemoveAll(tmpDir)

	// Create test CSV files
	authorsCSV := `email;firstname;lastname
test1@example.com;Test;User
`
	err = os.WriteFile(tmpDir+"/authors.csv", []byte(authorsCSV), 0644)
	g.Expect(err).To(BeNil())

	booksCSV := `title;isbn;authors;description
Test Book;1234-5678-9012;test1@example.com;A test book
`
	err = os.WriteFile(tmpDir+"/books.csv", []byte(booksCSV), 0644)
	g.Expect(err).To(BeNil())

	magazinesCSV := `title;isbn;authors;publishedAt
Test Magazine;9876-5432-1098;test1@example.com;01.01.2020
`
	err = os.WriteFile(tmpDir+"/magazines.csv", []byte(magazinesCSV), 0644)
	g.Expect(err).To(BeNil())

	// Load the library
	lib, err := LoadLibrary(tmpDir)
	g.Expect(err).To(BeNil())
	g.Expect(len(lib.Authors)).To(Equal(1))
	g.Expect(len(lib.Books)).To(Equal(1))
	g.Expect(len(lib.Magazines)).To(Equal(1))

	g.Expect(lib.Books[0].Title).To(Equal("Test Book"))
	g.Expect(lib.Magazines[0].Title).To(Equal("Test Magazine"))
}
