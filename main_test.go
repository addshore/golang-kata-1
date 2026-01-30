package main

import (
	"strings"
	"testing"

	"github.com/onsi/gomega"
)

func TestLoadAuthors(t *testing.T) {
	g := gomega.NewGomegaWithT(t)

	csvData := `email;firstname;lastname
test@example.com;John;Doe
jane@example.com;Jane;Smith`

	authors, err := loadAuthorsFromReader(strings.NewReader(csvData))

	g.Expect(err).To(gomega.BeNil())
	g.Expect(len(authors)).To(gomega.Equal(2))
	g.Expect(authors["test@example.com"]).To(gomega.Equal(Author{Email: "test@example.com", FirstName: "John", LastName: "Doe"}))
	g.Expect(authors["jane@example.com"]).To(gomega.Equal(Author{Email: "jane@example.com", FirstName: "Jane", LastName: "Smith"}))
}

func TestLoadBooks(t *testing.T) {
	g := gomega.NewGomegaWithT(t)

	csvData := `title;isbn;authors;description
Test Book;123-456;test@example.com;Test description
Another Book;789-012;jane@example.com;Another description`

	books, err := loadBooksFromReader(strings.NewReader(csvData))

	g.Expect(err).To(gomega.BeNil())
	g.Expect(len(books)).To(gomega.Equal(2))
	g.Expect(books[0].Title).To(gomega.Equal("Test Book"))
	g.Expect(books[0].ISBN).To(gomega.Equal("123-456"))
	g.Expect(books[0].Authors).To(gomega.Equal([]string{"test@example.com"}))
	g.Expect(books[0].Description).To(gomega.Equal("Test description"))
}

func TestLoadMagazines(t *testing.T) {
	g := gomega.NewGomegaWithT(t)

	csvData := `title;isbn;authors;publishedAt
Test Mag;111-222;test@example.com;01.01.2020
Another Mag;333-444;jane@example.com;02.02.2020`

	magazines, err := loadMagazinesFromReader(strings.NewReader(csvData))

	g.Expect(err).To(gomega.BeNil())
	g.Expect(len(magazines)).To(gomega.Equal(2))
	g.Expect(magazines[0].Title).To(gomega.Equal("Test Mag"))
	g.Expect(magazines[0].ISBN).To(gomega.Equal("111-222"))
	g.Expect(magazines[0].Authors).To(gomega.Equal([]string{"test@example.com"}))
	g.Expect(magazines[0].PublishedAt).To(gomega.Equal("01.01.2020"))
}

func TestGetAuthorsNames(t *testing.T) {
	g := gomega.NewGomegaWithT(t)

	lib := &Library{
		Authors: map[string]Author{
			"test@example.com": {Email: "test@example.com", FirstName: "John", LastName: "Doe"},
			"jane@example.com": {Email: "jane@example.com", FirstName: "Jane", LastName: "Smith"},
		},
	}

	result := lib.getAuthorsNames([]string{"test@example.com", "jane@example.com"})
	g.Expect(result).To(gomega.Equal("John Doe, Jane Smith"))
}

func TestGetItems(t *testing.T) {
	g := gomega.NewGomegaWithT(t)

	lib := &Library{
		Authors: map[string]Author{
			"test@example.com": {Email: "test@example.com", FirstName: "John", LastName: "Doe"},
		},
		Books: []Book{
			{Title: "Book A", ISBN: "123", Authors: []string{"test@example.com"}, Description: "Desc A"},
			{Title: "Book B", ISBN: "456", Authors: []string{"test@example.com"}, Description: "Desc B"},
		},
		Magazines: []Magazine{
			{Title: "Mag A", ISBN: "789", Authors: []string{"test@example.com"}, PublishedAt: "01.01.2020"},
		},
	}

	// Test no search
	items := lib.getItems("", "")
	g.Expect(len(items)).To(gomega.Equal(3))
	// Check sorted by title
	g.Expect(items[0].Title).To(gomega.Equal("Book A"))
	g.Expect(items[1].Title).To(gomega.Equal("Book B"))
	g.Expect(items[2].Title).To(gomega.Equal("Mag A"))

	// Test search by ISBN
	items = lib.getItems("123", "")
	g.Expect(len(items)).To(gomega.Equal(1))
	g.Expect(items[0].Title).To(gomega.Equal("Book A"))

	// Test search by email
	items = lib.getItems("", "test@example.com")
	g.Expect(len(items)).To(gomega.Equal(3))
}