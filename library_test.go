package main

import (
	"strings"
	"testing"

	"github.com/onsi/gomega"
)

func TestFormatAuthors(t *testing.T) {
	g := gomega.NewWithT(t)
	
	authors := map[string]Author{
		"john@example.com": {Email: "john@example.com", Firstname: "John", Lastname: "Doe"},
		"jane@example.com": {Email: "jane@example.com", Firstname: "Jane", Lastname: "Smith"},
	}
	
	t.Run("single author", func(t *testing.T) {
		result := FormatAuthors([]string{"john@example.com"}, authors)
		g.Expect(result).To(gomega.Equal("John Doe"))
	})
	
	t.Run("multiple authors", func(t *testing.T) {
		result := FormatAuthors([]string{"john@example.com", "jane@example.com"}, authors)
		g.Expect(result).To(gomega.Equal("John Doe, Jane Smith"))
	})
	
	t.Run("unknown author", func(t *testing.T) {
		result := FormatAuthors([]string{"unknown@example.com"}, authors)
		g.Expect(result).To(gomega.Equal(""))
	})
}

func TestLibraryFindByISBN(t *testing.T) {
	g := gomega.NewWithT(t)
	
	books := []Book{
		{Title: "Test Book", ISBN: "123-456", Authors: []string{"author@test.com"}, Description: "Test"},
	}
	magazines := []Magazine{
		{Title: "Test Magazine", ISBN: "789-012", Authors: []string{"author@test.com"}, PublishedAt: "2023-01-01"},
	}
	library := NewLibrary(nil, books, magazines)
	
	t.Run("find existing book", func(t *testing.T) {
		item, found := library.FindByISBN("123-456")
		g.Expect(found).To(gomega.BeTrue())
		book, ok := item.(Book)
		g.Expect(ok).To(gomega.BeTrue())
		g.Expect(book.Title).To(gomega.Equal("Test Book"))
	})
	
	t.Run("find existing magazine", func(t *testing.T) {
		item, found := library.FindByISBN("789-012")
		g.Expect(found).To(gomega.BeTrue())
		magazine, ok := item.(Magazine)
		g.Expect(ok).To(gomega.BeTrue())
		g.Expect(magazine.Title).To(gomega.Equal("Test Magazine"))
	})
	
	t.Run("find non-existing item", func(t *testing.T) {
		_, found := library.FindByISBN("999-999")
		g.Expect(found).To(gomega.BeFalse())
	})
}

func TestLibraryFindByAuthor(t *testing.T) {
	g := gomega.NewWithT(t)
	
	books := []Book{
		{Title: "Book 1", ISBN: "123", Authors: []string{"author1@test.com"}, Description: "Test"},
		{Title: "Book 2", ISBN: "456", Authors: []string{"author2@test.com"}, Description: "Test"},
		{Title: "Book 3", ISBN: "789", Authors: []string{"author1@test.com", "author2@test.com"}, Description: "Test"},
	}
	magazines := []Magazine{
		{Title: "Magazine 1", ISBN: "111", Authors: []string{"author1@test.com"}, PublishedAt: "2023-01-01"},
		{Title: "Magazine 2", ISBN: "222", Authors: []string{"author3@test.com"}, PublishedAt: "2023-01-01"},
	}
	library := NewLibrary(nil, books, magazines)
	
	t.Run("find books and magazines by author", func(t *testing.T) {
		foundBooks, foundMagazines := library.FindByAuthor("author1@test.com")
		g.Expect(len(foundBooks)).To(gomega.Equal(2))
		g.Expect(len(foundMagazines)).To(gomega.Equal(1))
		g.Expect(foundBooks[0].Title).To(gomega.Equal("Book 1"))
		g.Expect(foundBooks[1].Title).To(gomega.Equal("Book 3"))
		g.Expect(foundMagazines[0].Title).To(gomega.Equal("Magazine 1"))
	})
	
	t.Run("find no items by unknown author", func(t *testing.T) {
		foundBooks, foundMagazines := library.FindByAuthor("unknown@test.com")
		g.Expect(len(foundBooks)).To(gomega.Equal(0))
		g.Expect(len(foundMagazines)).To(gomega.Equal(0))
	})
}

func TestLibraryGetSortedItems(t *testing.T) {
	g := gomega.NewWithT(t)
	
	books := []Book{
		{Title: "Z Book", ISBN: "123", Authors: []string{"author@test.com"}, Description: "Test"},
		{Title: "A Book", ISBN: "456", Authors: []string{"author@test.com"}, Description: "Test"},
	}
	magazines := []Magazine{
		{Title: "M Magazine", ISBN: "789", Authors: []string{"author@test.com"}, PublishedAt: "2023-01-01"},
		{Title: "B Magazine", ISBN: "012", Authors: []string{"author@test.com"}, PublishedAt: "2023-01-01"},
	}
	library := NewLibrary(nil, books, magazines)
	
	t.Run("sort ascending", func(t *testing.T) {
		items := library.GetSortedItems(false)
		g.Expect(len(items)).To(gomega.Equal(4))
		g.Expect(items[0].Title).To(gomega.Equal("A Book"))
		g.Expect(items[1].Title).To(gomega.Equal("B Magazine"))
		g.Expect(items[2].Title).To(gomega.Equal("M Magazine"))
		g.Expect(items[3].Title).To(gomega.Equal("Z Book"))
	})
	
	t.Run("sort descending", func(t *testing.T) {
		items := library.GetSortedItems(true)
		g.Expect(len(items)).To(gomega.Equal(4))
		g.Expect(items[0].Title).To(gomega.Equal("Z Book"))
		g.Expect(items[1].Title).To(gomega.Equal("M Magazine"))
		g.Expect(items[2].Title).To(gomega.Equal("B Magazine"))
		g.Expect(items[3].Title).To(gomega.Equal("A Book"))
	})
}

func TestParseAuthorsCSV(t *testing.T) {
	g := gomega.NewWithT(t)
	
	csvData := `email;firstname;lastname
john@example.com;John;Doe
jane@example.com;Jane;Smith`
	
	authors, err := ParseAuthorsCSV(strings.NewReader(csvData))
	g.Expect(err).To(gomega.BeNil())
	g.Expect(len(authors)).To(gomega.Equal(2))
	g.Expect(authors["john@example.com"].Firstname).To(gomega.Equal("John"))
	g.Expect(authors["john@example.com"].Lastname).To(gomega.Equal("Doe"))
	g.Expect(authors["jane@example.com"].Firstname).To(gomega.Equal("Jane"))
	g.Expect(authors["jane@example.com"].Lastname).To(gomega.Equal("Smith"))
}

func TestParseBooksCSV(t *testing.T) {
	g := gomega.NewWithT(t)
	
	csvData := `title;isbn;authors;description
Test Book;123-456;author1@test.com,author2@test.com;A test book
Another Book;789-012;author3@test.com;Another test book`
	
	books, err := ParseBooksCSV(strings.NewReader(csvData))
	g.Expect(err).To(gomega.BeNil())
	g.Expect(len(books)).To(gomega.Equal(2))
	g.Expect(books[0].Title).To(gomega.Equal("Test Book"))
	g.Expect(books[0].ISBN).To(gomega.Equal("123-456"))
	g.Expect(books[0].Authors).To(gomega.Equal([]string{"author1@test.com", "author2@test.com"}))
	g.Expect(books[0].Description).To(gomega.Equal("A test book"))
}

func TestParseMagazinesCSV(t *testing.T) {
	g := gomega.NewWithT(t)
	
	csvData := `title;isbn;authors;publishedAt
Test Magazine;123-456;author1@test.com;2023-01-01
Another Magazine;789-012;author2@test.com,author3@test.com;2023-02-01`
	
	magazines, err := ParseMagazinesCSV(strings.NewReader(csvData))
	g.Expect(err).To(gomega.BeNil())
	g.Expect(len(magazines)).To(gomega.Equal(2))
	g.Expect(magazines[0].Title).To(gomega.Equal("Test Magazine"))
	g.Expect(magazines[0].ISBN).To(gomega.Equal("123-456"))
	g.Expect(magazines[0].Authors).To(gomega.Equal([]string{"author1@test.com"}))
	g.Expect(magazines[0].PublishedAt).To(gomega.Equal("2023-01-01"))
	g.Expect(magazines[1].Authors).To(gomega.Equal([]string{"author2@test.com", "author3@test.com"}))
}

func TestLibraryAddBook(t *testing.T) {
	g := gomega.NewWithT(t)
	
	library := NewLibrary(make(map[string]Author), []Book{}, []Magazine{})
	book := Book{Title: "New Book", ISBN: "123", Authors: []string{"author@test.com"}, Description: "Test"}
	
	library.AddBook(book)
	
	g.Expect(len(library.Books)).To(gomega.Equal(1))
	g.Expect(library.Books[0].Title).To(gomega.Equal("New Book"))
}

func TestLibraryAddMagazine(t *testing.T) {
	g := gomega.NewWithT(t)
	
	library := NewLibrary(make(map[string]Author), []Book{}, []Magazine{})
	magazine := Magazine{Title: "New Magazine", ISBN: "456", Authors: []string{"author@test.com"}, PublishedAt: "2024-01-01"}
	
	library.AddMagazine(magazine)
	
	g.Expect(len(library.Magazines)).To(gomega.Equal(1))
	g.Expect(library.Magazines[0].Title).To(gomega.Equal("New Magazine"))
}

func TestLibraryAddAuthor(t *testing.T) {
	g := gomega.NewWithT(t)
	
	library := NewLibrary(make(map[string]Author), []Book{}, []Magazine{})
	author := Author{Email: "new@test.com", Firstname: "New", Lastname: "Author"}
	
	library.AddAuthor(author)
	
	g.Expect(len(library.Authors)).To(gomega.Equal(1))
	g.Expect(library.Authors["new@test.com"].Firstname).To(gomega.Equal("New"))
}

func TestWriteAuthorsCSV(t *testing.T) {
	g := gomega.NewWithT(t)
	
	authors := map[string]Author{
		"test@example.com": {Email: "test@example.com", Firstname: "Test", Lastname: "Author"},
	}
	
	var buf strings.Builder
	err := WriteAuthorsCSV(&buf, authors)
	
	g.Expect(err).To(gomega.BeNil())
	result := buf.String()
	g.Expect(result).To(gomega.ContainSubstring("email;firstname;lastname"))
	g.Expect(result).To(gomega.ContainSubstring("test@example.com;Test;Author"))
}