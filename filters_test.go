package main

import (
	"testing"

	. "github.com/onsi/gomega"
)

func TestFilterBooksByISBN(t *testing.T) {
	g := NewGomegaWithT(t)

	books := []Book{
		{Title: "Book1", ISBN: "123", Authors: []string{"author1"}, Description: "Desc1"},
		{Title: "Book2", ISBN: "456", Authors: []string{"author2"}, Description: "Desc2"},
		{Title: "Book3", ISBN: "123", Authors: []string{"author3"}, Description: "Desc3"},
	}

	t.Run("filters books by ISBN", func(t *testing.T) {
		filtered := FilterBooksByISBN(books, "123")
		g.Expect(filtered).To(HaveLen(2))
		g.Expect(filtered[0].Title).To(Equal("Book1"))
		g.Expect(filtered[1].Title).To(Equal("Book3"))
	})

	t.Run("returns empty slice when no matches", func(t *testing.T) {
		filtered := FilterBooksByISBN(books, "999")
		g.Expect(filtered).To(BeEmpty())
	})

	t.Run("returns all books when ISBN is empty", func(t *testing.T) {
		filtered := FilterBooksByISBN(books, "")
		g.Expect(filtered).To(Equal(books))
	})
}

func TestFilterMagazinesByISBN(t *testing.T) {
	g := NewGomegaWithT(t)

	magazines := []Magazine{
		{Title: "Mag1", ISBN: "123", Authors: []string{"author1"}},
		{Title: "Mag2", ISBN: "456", Authors: []string{"author2"}},
		{Title: "Mag3", ISBN: "123", Authors: []string{"author3"}},
	}

	t.Run("filters magazines by ISBN", func(t *testing.T) {
		filtered := FilterMagazinesByISBN(magazines, "456")
		g.Expect(filtered).To(HaveLen(1))
		g.Expect(filtered[0].Title).To(Equal("Mag2"))
	})

	t.Run("returns all magazines when ISBN is empty", func(t *testing.T) {
		filtered := FilterMagazinesByISBN(magazines, "")
		g.Expect(filtered).To(Equal(magazines))
	})
}

func TestFilterBooksByAuthor(t *testing.T) {
	g := NewGomegaWithT(t)

	books := []Book{
		{Title: "Book1", ISBN: "123", Authors: []string{"author1@test.com", "author2@test.com"}},
		{Title: "Book2", ISBN: "456", Authors: []string{"author2@test.com"}},
		{Title: "Book3", ISBN: "789", Authors: []string{"author3@test.com"}},
	}

	t.Run("filters books by author email", func(t *testing.T) {
		filtered := FilterBooksByAuthor(books, "author2@test.com")
		g.Expect(filtered).To(HaveLen(2))
		g.Expect(filtered[0].Title).To(Equal("Book1"))
		g.Expect(filtered[1].Title).To(Equal("Book2"))
	})

	t.Run("returns empty slice when no matches", func(t *testing.T) {
		filtered := FilterBooksByAuthor(books, "unknown@test.com")
		g.Expect(filtered).To(BeEmpty())
	})

	t.Run("returns all books when author email is empty", func(t *testing.T) {
		filtered := FilterBooksByAuthor(books, "")
		g.Expect(filtered).To(Equal(books))
	})
}

func TestFilterMagazinesByAuthor(t *testing.T) {
	g := NewGomegaWithT(t)

	magazines := []Magazine{
		{Title: "Mag1", ISBN: "123", Authors: []string{"author1@test.com"}},
		{Title: "Mag2", ISBN: "456", Authors: []string{"author1@test.com", "author2@test.com"}},
		{Title: "Mag3", ISBN: "789", Authors: []string{"author3@test.com"}},
	}

	t.Run("filters magazines by author email", func(t *testing.T) {
		filtered := FilterMagazinesByAuthor(magazines, "author1@test.com")
		g.Expect(filtered).To(HaveLen(2))
		g.Expect(filtered[0].Title).To(Equal("Mag1"))
		g.Expect(filtered[1].Title).To(Equal("Mag2"))
	})

	t.Run("returns all magazines when author email is empty", func(t *testing.T) {
		filtered := FilterMagazinesByAuthor(magazines, "")
		g.Expect(filtered).To(Equal(magazines))
	})
}