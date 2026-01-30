package main

import (
	"strings"
	"testing"

	. "github.com/onsi/gomega"
)

func TestWelcomeMessage(t *testing.T) {
	g := NewGomegaWithT(t)

	// given
	expected := "Welcome to the Library Application!"

	// when
	actual := welcomeMessage()

	// then
	g.Expect(actual).To(Equal(expected))
}

func TestTruncateString(t *testing.T) {
	g := NewGomegaWithT(t)

	t.Run("string shorter than max length", func(t *testing.T) {
		result := truncateString("Hello", 10)
		g.Expect(result).To(Equal("Hello"))
	})

	t.Run("string equal to max length", func(t *testing.T) {
		result := truncateString("Hello", 5)
		g.Expect(result).To(Equal("Hello"))
	})

	t.Run("string longer than max length", func(t *testing.T) {
		result := truncateString("Hello World", 5)
		g.Expect(result).To(Equal("Hello..."))
	})

	t.Run("empty string", func(t *testing.T) {
		result := truncateString("", 10)
		g.Expect(result).To(Equal(""))
	})
}

func TestFormatAuthors(t *testing.T) {
	g := NewGomegaWithT(t)

	library := &Library{
		Authors: map[string]Author{
			"john@example.com": {Email: "john@example.com", FirstName: "John", LastName: "Doe"},
			"jane@example.com": {Email: "jane@example.com", FirstName: "Jane", LastName: "Smith"},
		},
	}

	t.Run("single known author", func(t *testing.T) {
		result := formatAuthors(library, []string{"john@example.com"})
		g.Expect(result).To(Equal("John Doe"))
	})

	t.Run("multiple known authors", func(t *testing.T) {
		result := formatAuthors(library, []string{"john@example.com", "jane@example.com"})
		g.Expect(result).To(Equal("John Doe, Jane Smith"))
	})

	t.Run("unknown author returns email", func(t *testing.T) {
		result := formatAuthors(library, []string{"unknown@example.com"})
		g.Expect(result).To(Equal("unknown@example.com"))
	})

	t.Run("mixed known and unknown authors", func(t *testing.T) {
		result := formatAuthors(library, []string{"john@example.com", "unknown@example.com"})
		g.Expect(result).To(Equal("John Doe, unknown@example.com"))
	})

	t.Run("trims whitespace from emails", func(t *testing.T) {
		result := formatAuthors(library, []string{" john@example.com "})
		g.Expect(result).To(Equal("John Doe"))
	})

	t.Run("empty email list", func(t *testing.T) {
		result := formatAuthors(library, []string{})
		g.Expect(result).To(Equal(""))
	})
}

func TestParseAuthors(t *testing.T) {
	g := NewGomegaWithT(t)

	t.Run("parses valid CSV", func(t *testing.T) {
		csv := "email;firstname;lastname\njohn@example.com;John;Doe\njane@example.com;Jane;Smith"
		reader := strings.NewReader(csv)

		authors, err := parseAuthors(reader)

		g.Expect(err).To(BeNil())
		g.Expect(authors).To(HaveLen(2))
		g.Expect(authors["john@example.com"]).To(Equal(Author{
			Email:     "john@example.com",
			FirstName: "John",
			LastName:  "Doe",
		}))
		g.Expect(authors["jane@example.com"]).To(Equal(Author{
			Email:     "jane@example.com",
			FirstName: "Jane",
			LastName:  "Smith",
		}))
	})

	t.Run("skips header row", func(t *testing.T) {
		csv := "email;firstname;lastname"
		reader := strings.NewReader(csv)

		authors, err := parseAuthors(reader)

		g.Expect(err).To(BeNil())
		g.Expect(authors).To(HaveLen(0))
	})

	t.Run("handles empty CSV", func(t *testing.T) {
		csv := ""
		reader := strings.NewReader(csv)

		authors, err := parseAuthors(reader)

		g.Expect(err).To(BeNil())
		g.Expect(authors).To(HaveLen(0))
	})

	t.Run("returns error for malformed CSV", func(t *testing.T) {
		csv := "email;firstname;lastname\njohn@example.com;John"
		reader := strings.NewReader(csv)

		_, err := parseAuthors(reader)

		g.Expect(err).ToNot(BeNil())
	})
}

func TestParseBooks(t *testing.T) {
	g := NewGomegaWithT(t)

	t.Run("parses valid CSV", func(t *testing.T) {
		csv := "title;isbn;authors;description\nTest Book;123-456;author1@example.com;A test book"
		reader := strings.NewReader(csv)

		books, err := parseBooks(reader)

		g.Expect(err).To(BeNil())
		g.Expect(books).To(HaveLen(1))
		g.Expect(books[0].Title).To(Equal("Test Book"))
		g.Expect(books[0].ISBN).To(Equal("123-456"))
		g.Expect(books[0].AuthorEmails).To(Equal([]string{"author1@example.com"}))
		g.Expect(books[0].Description).To(Equal("A test book"))
	})

	t.Run("parses multiple authors", func(t *testing.T) {
		csv := "title;isbn;authors;description\nTest Book;123-456;author1@example.com,author2@example.com;A test book"
		reader := strings.NewReader(csv)

		books, err := parseBooks(reader)

		g.Expect(err).To(BeNil())
		g.Expect(books).To(HaveLen(1))
		g.Expect(books[0].AuthorEmails).To(Equal([]string{"author1@example.com", "author2@example.com"}))
	})

	t.Run("handles empty CSV", func(t *testing.T) {
		csv := ""
		reader := strings.NewReader(csv)

		books, err := parseBooks(reader)

		g.Expect(err).To(BeNil())
		g.Expect(books).To(HaveLen(0))
	})

	t.Run("parses multiple books", func(t *testing.T) {
		csv := "title;isbn;authors;description\nBook 1;111;a@b.com;Desc 1\nBook 2;222;c@d.com;Desc 2"
		reader := strings.NewReader(csv)

		books, err := parseBooks(reader)

		g.Expect(err).To(BeNil())
		g.Expect(books).To(HaveLen(2))
		g.Expect(books[0].Title).To(Equal("Book 1"))
		g.Expect(books[1].Title).To(Equal("Book 2"))
	})
}

func TestParseMagazines(t *testing.T) {
	g := NewGomegaWithT(t)

	t.Run("parses valid CSV", func(t *testing.T) {
		csv := "title;isbn;authors;publishedAt\nTest Magazine;789-012;author@example.com;01.01.2024"
		reader := strings.NewReader(csv)

		magazines, err := parseMagazines(reader)

		g.Expect(err).To(BeNil())
		g.Expect(magazines).To(HaveLen(1))
		g.Expect(magazines[0].Title).To(Equal("Test Magazine"))
		g.Expect(magazines[0].ISBN).To(Equal("789-012"))
		g.Expect(magazines[0].AuthorEmails).To(Equal([]string{"author@example.com"}))
		g.Expect(magazines[0].PublishedAt).To(Equal("01.01.2024"))
	})

	t.Run("parses multiple authors", func(t *testing.T) {
		csv := "title;isbn;authors;publishedAt\nTest Magazine;789-012;a@b.com,c@d.com;01.01.2024"
		reader := strings.NewReader(csv)

		magazines, err := parseMagazines(reader)

		g.Expect(err).To(BeNil())
		g.Expect(magazines).To(HaveLen(1))
		g.Expect(magazines[0].AuthorEmails).To(Equal([]string{"a@b.com", "c@d.com"}))
	})

	t.Run("handles empty CSV", func(t *testing.T) {
		csv := ""
		reader := strings.NewReader(csv)

		magazines, err := parseMagazines(reader)

		g.Expect(err).To(BeNil())
		g.Expect(magazines).To(HaveLen(0))
	})
}

func TestFindByISBN(t *testing.T) {
	g := NewGomegaWithT(t)

	library := &Library{
		Authors: map[string]Author{},
		Books: []Book{
			{Title: "Book 1", ISBN: "111-111", AuthorEmails: []string{"a@b.com"}, Description: "Desc 1"},
			{Title: "Book 2", ISBN: "222-222", AuthorEmails: []string{"c@d.com"}, Description: "Desc 2"},
		},
		Magazines: []Magazine{
			{Title: "Magazine 1", ISBN: "333-333", AuthorEmails: []string{"e@f.com"}, PublishedAt: "01.01.2024"},
			{Title: "Magazine 2", ISBN: "444-444", AuthorEmails: []string{"g@h.com"}, PublishedAt: "02.02.2024"},
		},
	}

	t.Run("finds book by ISBN", func(t *testing.T) {
		results := findByISBN(library, "111-111")

		g.Expect(results).To(HaveLen(1))
		g.Expect(results[0].Type).To(Equal("Book"))
		g.Expect(results[0].Book.Title).To(Equal("Book 1"))
	})

	t.Run("finds magazine by ISBN", func(t *testing.T) {
		results := findByISBN(library, "333-333")

		g.Expect(results).To(HaveLen(1))
		g.Expect(results[0].Type).To(Equal("Magazine"))
		g.Expect(results[0].Magazine.Title).To(Equal("Magazine 1"))
	})

	t.Run("returns empty for non-existent ISBN", func(t *testing.T) {
		results := findByISBN(library, "999-999")

		g.Expect(results).To(HaveLen(0))
	})

	t.Run("returns empty for empty ISBN", func(t *testing.T) {
		results := findByISBN(library, "")

		g.Expect(results).To(HaveLen(0))
	})
}

func TestFindByAuthorEmail(t *testing.T) {
	g := NewGomegaWithT(t)

	library := &Library{
		Authors: map[string]Author{
			"author1@example.com": {Email: "author1@example.com", FirstName: "Author", LastName: "One"},
		},
		Books: []Book{
			{Title: "Book 1", ISBN: "111", AuthorEmails: []string{"author1@example.com"}, Description: "Desc 1"},
			{Title: "Book 2", ISBN: "222", AuthorEmails: []string{"author2@example.com"}, Description: "Desc 2"},
			{Title: "Book 3", ISBN: "333", AuthorEmails: []string{"author1@example.com", "author2@example.com"}, Description: "Desc 3"},
		},
		Magazines: []Magazine{
			{Title: "Magazine 1", ISBN: "444", AuthorEmails: []string{"author1@example.com"}, PublishedAt: "01.01.2024"},
			{Title: "Magazine 2", ISBN: "555", AuthorEmails: []string{"author3@example.com"}, PublishedAt: "02.02.2024"},
		},
	}

	t.Run("finds books and magazines by author email", func(t *testing.T) {
		results := findByAuthorEmail(library, "author1@example.com")

		g.Expect(results).To(HaveLen(3)) // Book 1, Book 3, Magazine 1

		bookCount := 0
		magazineCount := 0
		for _, r := range results {
			if r.Type == "Book" {
				bookCount++
			} else {
				magazineCount++
			}
		}
		g.Expect(bookCount).To(Equal(2))
		g.Expect(magazineCount).To(Equal(1))
	})

	t.Run("finds only books for book-only author", func(t *testing.T) {
		results := findByAuthorEmail(library, "author2@example.com")

		g.Expect(results).To(HaveLen(2)) // Book 2, Book 3
		for _, r := range results {
			g.Expect(r.Type).To(Equal("Book"))
		}
	})

	t.Run("finds only magazines for magazine-only author", func(t *testing.T) {
		results := findByAuthorEmail(library, "author3@example.com")

		g.Expect(results).To(HaveLen(1))
		g.Expect(results[0].Type).To(Equal("Magazine"))
		g.Expect(results[0].Magazine.Title).To(Equal("Magazine 2"))
	})

	t.Run("returns empty for unknown author", func(t *testing.T) {
		results := findByAuthorEmail(library, "unknown@example.com")

		g.Expect(results).To(HaveLen(0))
	})

	t.Run("handles whitespace in author emails", func(t *testing.T) {
		libraryWithSpaces := &Library{
			Authors: map[string]Author{},
			Books: []Book{
				{Title: "Book", ISBN: "111", AuthorEmails: []string{" author@example.com "}, Description: "Desc"},
			},
			Magazines: []Magazine{},
		}

		results := findByAuthorEmail(libraryWithSpaces, "author@example.com")

		g.Expect(results).To(HaveLen(1))
	})
}

func TestGetAllItemsSortedByTitle(t *testing.T) {
	g := NewGomegaWithT(t)

	library := &Library{
		Authors: map[string]Author{
			"a@b.com": {Email: "a@b.com", FirstName: "John", LastName: "Doe"},
		},
		Books: []Book{
			{Title: "Zebra Book", ISBN: "111", AuthorEmails: []string{"a@b.com"}, Description: "Desc"},
			{Title: "Apple Book", ISBN: "222", AuthorEmails: []string{"a@b.com"}, Description: "Desc"},
		},
		Magazines: []Magazine{
			{Title: "Mango Magazine", ISBN: "333", AuthorEmails: []string{"a@b.com"}, PublishedAt: "01.01.2024"},
			{Title: "Banana Magazine", ISBN: "444", AuthorEmails: []string{"a@b.com"}, PublishedAt: "02.02.2024"},
		},
	}

	t.Run("returns items sorted by title ascending", func(t *testing.T) {
		items := getAllItemsSortedByTitle(library, true)

		g.Expect(items).To(HaveLen(4))
		g.Expect(items[0].Title).To(Equal("Apple Book"))
		g.Expect(items[1].Title).To(Equal("Banana Magazine"))
		g.Expect(items[2].Title).To(Equal("Mango Magazine"))
		g.Expect(items[3].Title).To(Equal("Zebra Book"))
	})

	t.Run("returns items sorted by title descending", func(t *testing.T) {
		items := getAllItemsSortedByTitle(library, false)

		g.Expect(items).To(HaveLen(4))
		g.Expect(items[0].Title).To(Equal("Zebra Book"))
		g.Expect(items[1].Title).To(Equal("Mango Magazine"))
		g.Expect(items[2].Title).To(Equal("Banana Magazine"))
		g.Expect(items[3].Title).To(Equal("Apple Book"))
	})

	t.Run("sorts case-insensitively", func(t *testing.T) {
		libraryMixedCase := &Library{
			Authors: map[string]Author{},
			Books: []Book{
				{Title: "zebra", ISBN: "111", AuthorEmails: []string{}, Description: ""},
				{Title: "Apple", ISBN: "222", AuthorEmails: []string{}, Description: ""},
			},
			Magazines: []Magazine{},
		}

		items := getAllItemsSortedByTitle(libraryMixedCase, true)

		g.Expect(items).To(HaveLen(2))
		g.Expect(items[0].Title).To(Equal("Apple"))
		g.Expect(items[1].Title).To(Equal("zebra"))
	})

	t.Run("includes correct item types", func(t *testing.T) {
		items := getAllItemsSortedByTitle(library, true)

		bookCount := 0
		magazineCount := 0
		for _, item := range items {
			if item.ItemType == "Book" {
				bookCount++
			} else if item.ItemType == "Magazine" {
				magazineCount++
			}
		}
		g.Expect(bookCount).To(Equal(2))
		g.Expect(magazineCount).To(Equal(2))
	})

	t.Run("formats authors correctly", func(t *testing.T) {
		items := getAllItemsSortedByTitle(library, true)

		for _, item := range items {
			g.Expect(item.Authors).To(Equal("John Doe"))
		}
	})

	t.Run("returns empty for empty library", func(t *testing.T) {
		emptyLibrary := &Library{
			Authors:   map[string]Author{},
			Books:     []Book{},
			Magazines: []Magazine{},
		}

		items := getAllItemsSortedByTitle(emptyLibrary, true)

		g.Expect(items).To(HaveLen(0))
	})
}

func TestLoadLibraryIntegration(t *testing.T) {
	g := NewGomegaWithT(t)

	// This is an integration test that uses the actual CSV files
	t.Run("loads all data from CSV files", func(t *testing.T) {
		library, err := loadLibrary()

		g.Expect(err).To(BeNil())
		g.Expect(library.Authors).ToNot(BeEmpty())
		g.Expect(library.Books).ToNot(BeEmpty())
		g.Expect(library.Magazines).ToNot(BeEmpty())
	})

	t.Run("loads correct number of authors", func(t *testing.T) {
		library, err := loadLibrary()

		g.Expect(err).To(BeNil())
		g.Expect(library.Authors).To(HaveLen(6))
	})

	t.Run("loads correct number of books", func(t *testing.T) {
		library, err := loadLibrary()

		g.Expect(err).To(BeNil())
		g.Expect(library.Books).To(HaveLen(8))
	})

	t.Run("loads correct number of magazines", func(t *testing.T) {
		library, err := loadLibrary()

		g.Expect(err).To(BeNil())
		g.Expect(library.Magazines).To(HaveLen(6))
	})
}
