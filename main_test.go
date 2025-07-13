package main

import (
	"testing"
	"time"

	. "github.com/onsi/gomega"
)

func TestLoadLibrary(t *testing.T) {
	g := NewGomegaWithT(t)

	// when
	library, err := loadLibrary()

	// then
	g.Expect(err).To(BeNil())
	g.Expect(library).ToNot(BeNil())
	g.Expect(len(library.Authors)).To(BeNumerically(">", 0))
	g.Expect(len(library.Books)).To(BeNumerically(">", 0))
	g.Expect(len(library.Magazines)).To(BeNumerically(">", 0))
}

func TestFormatAuthors(t *testing.T) {
	g := NewGomegaWithT(t)

	// given
	authors := []Author{
		{Email: "test@example.com", FirstName: "John", LastName: "Doe"},
		{Email: "test2@example.com", FirstName: "Jane", LastName: "Smith"},
	}

	// when
	result := formatAuthors(authors)

	// then
	g.Expect(result).To(Equal("John Doe, Jane Smith"))
}

func TestFormatAuthorsEmpty(t *testing.T) {
	g := NewGomegaWithT(t)

	// given
	authors := []Author{}

	// when
	result := formatAuthors(authors)

	// then
	g.Expect(result).To(Equal("Unknown"))
}

func TestTruncateDescription(t *testing.T) {
	g := NewGomegaWithT(t)

	// given
	longDescription := "This is a very long description that should be truncated because it exceeds the maximum length limit set for display purposes."
	maxLength := 50

	// when
	result := truncateDescription(longDescription, maxLength)

	// then
	g.Expect(result).To(HaveLen(maxLength + 3)) // +3 for "..."
	g.Expect(result).To(HaveSuffix("..."))
}

func TestParseAuthors(t *testing.T) {
	g := NewGomegaWithT(t)

	// given
	authorMap := map[string]Author{
		"john@example.com": {Email: "john@example.com", FirstName: "John", LastName: "Doe"},
		"jane@example.com": {Email: "jane@example.com", FirstName: "Jane", LastName: "Smith"},
	}
	authorEmails := "john@example.com,jane@example.com"

	// when
	result := parseAuthors(authorEmails, authorMap)

	// then
	g.Expect(result).To(HaveLen(2))
	g.Expect(result[0].FirstName).To(Equal("John"))
	g.Expect(result[1].FirstName).To(Equal("Jane"))
}

func TestBookStructure(t *testing.T) {
	g := NewGomegaWithT(t)

	// given
	author := Author{Email: "test@example.com", FirstName: "Test", LastName: "Author"}
	book := Book{
		Title:       "Test Book",
		ISBN:        "123-456-789",
		Authors:     []Author{author},
		Description: "A test book description",
	}

	// then
	g.Expect(book.Title).To(Equal("Test Book"))
	g.Expect(book.ISBN).To(Equal("123-456-789"))
	g.Expect(book.Authors).To(HaveLen(1))
	g.Expect(book.Authors[0].FirstName).To(Equal("Test"))
}

func TestMagazineStructure(t *testing.T) {
	g := NewGomegaWithT(t)

	// given
	author := Author{Email: "test@example.com", FirstName: "Test", LastName: "Author"}
	publishedAt := time.Date(2023, 5, 15, 0, 0, 0, 0, time.UTC)
	magazine := Magazine{
		Title:       "Test Magazine",
		ISBN:        "987-654-321",
		Authors:     []Author{author},
		PublishedAt: publishedAt,
	}

	// then
	g.Expect(magazine.Title).To(Equal("Test Magazine"))
	g.Expect(magazine.ISBN).To(Equal("987-654-321"))
	g.Expect(magazine.Authors).To(HaveLen(1))
	g.Expect(magazine.PublishedAt.Year()).To(Equal(2023))
}

func TestSearchBookByISBN(t *testing.T) {
	g := NewGomegaWithT(t)

	// given
	library, err := loadLibrary()
	g.Expect(err).To(BeNil())

	// when searching for a known book ISBN
	foundBook := false
	for _, book := range library.Books {
		if book.ISBN == "5554-5545-4518" { // Known ISBN from test data
			foundBook = true
			g.Expect(book.Title).To(ContainSubstring("Ich helfe dir kochen"))
			break
		}
	}

	// then
	g.Expect(foundBook).To(BeTrue())
}

func TestSearchMagazineByISBN(t *testing.T) {
	g := NewGomegaWithT(t)

	// given
	library, err := loadLibrary()
	g.Expect(err).To(BeNil())

	// when searching for a known magazine ISBN
	foundMagazine := false
	for _, magazine := range library.Magazines {
		if magazine.ISBN == "5454-5587-3210" { // Known ISBN from test data
			foundMagazine = true
			g.Expect(magazine.Title).To(Equal("Beautiful cooking"))
			break
		}
	}

	// then
	g.Expect(foundMagazine).To(BeTrue())
}

func TestSearchNonExistentISBN(t *testing.T) {
	g := NewGomegaWithT(t)

	// given
	library, err := loadLibrary()
	g.Expect(err).To(BeNil())
	nonExistentISBN := "9999-9999-9999"

	// when searching for non-existent ISBN
	foundInBooks := false
	for _, book := range library.Books {
		if book.ISBN == nonExistentISBN {
			foundInBooks = true
			break
		}
	}

	foundInMagazines := false
	for _, magazine := range library.Magazines {
		if magazine.ISBN == nonExistentISBN {
			foundInMagazines = true
			break
		}
	}

	// then
	g.Expect(foundInBooks).To(BeFalse())
	g.Expect(foundInMagazines).To(BeFalse())
}
