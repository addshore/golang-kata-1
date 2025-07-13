package main

import (
	"strings"
	"testing"
	"time"

	. "github.com/onsi/gomega"
)

func TestParseAuthorsFromReader(t *testing.T) {
	g := NewGomegaWithT(t)

	t.Run("parses authors correctly", func(t *testing.T) {
		csvContent := `email;firstname;lastname
null-walter@echocat.org;Paul;Walter
null-mueller@echocat.org;Max;Müller`

		reader := strings.NewReader(csvContent)
		authors, err := parseAuthorsFromReader(reader)

		g.Expect(err).ToNot(HaveOccurred())
		g.Expect(authors).To(HaveLen(2))
		g.Expect(authors["null-walter@echocat.org"]).To(Equal(Author{
			Email:     "null-walter@echocat.org",
			FirstName: "Paul",
			LastName:  "Walter",
		}))
		g.Expect(authors["null-mueller@echocat.org"]).To(Equal(Author{
			Email:     "null-mueller@echocat.org",
			FirstName: "Max",
			LastName:  "Müller",
		}))
	})

	t.Run("handles BOM prefix", func(t *testing.T) {
		csvContent := "\ufeffemail;firstname;lastname\nnull-walter@echocat.org;Paul;Walter"

		reader := strings.NewReader(csvContent)
		authors, err := parseAuthorsFromReader(reader)

		g.Expect(err).ToNot(HaveOccurred())
		g.Expect(authors).To(HaveLen(1))
	})

	t.Run("skips invalid records", func(t *testing.T) {
		csvContent := `email;firstname;lastname
null-walter@echocat.org;Paul
null-mueller@echocat.org;Max;Müller`

		reader := strings.NewReader(csvContent)
		authors, err := parseAuthorsFromReader(reader)

		g.Expect(err).ToNot(HaveOccurred())
		g.Expect(authors).To(HaveLen(1))
		g.Expect(authors).To(HaveKey("null-mueller@echocat.org"))
	})

	t.Run("handles empty file", func(t *testing.T) {
		csvContent := `email;firstname;lastname`

		reader := strings.NewReader(csvContent)
		authors, err := parseAuthorsFromReader(reader)

		g.Expect(err).ToNot(HaveOccurred())
		g.Expect(authors).To(BeEmpty())
	})
}

func TestParseBooksFromReader(t *testing.T) {
	g := NewGomegaWithT(t)

	t.Run("parses books correctly", func(t *testing.T) {
		csvContent := `title;isbn;authors;description
Das große GU-Kochbuch;2145-8548-3325;null-ferdinand@echocat.org,null-lieblich@echocat.org;A great cookbook
Schlank im Schlaf;4545-8558-3232;null-gustafsson@echocat.org;Sleep yourself slim`

		reader := strings.NewReader(csvContent)
		books, err := parseBooksFromReader(reader)

		g.Expect(err).ToNot(HaveOccurred())
		g.Expect(books).To(HaveLen(2))
		g.Expect(books[0]).To(Equal(Book{
			Title:       "Das große GU-Kochbuch",
			ISBN:        "2145-8548-3325",
			Authors:     []string{"null-ferdinand@echocat.org", "null-lieblich@echocat.org"},
			Description: "A great cookbook",
		}))
		g.Expect(books[1]).To(Equal(Book{
			Title:       "Schlank im Schlaf",
			ISBN:        "4545-8558-3232",
			Authors:     []string{"null-gustafsson@echocat.org"},
			Description: "Sleep yourself slim",
		}))
	})

	t.Run("handles empty authors", func(t *testing.T) {
		csvContent := `title;isbn;authors;description
My Book;1234-5678-9012;;No authors`

		reader := strings.NewReader(csvContent)
		books, err := parseBooksFromReader(reader)

		g.Expect(err).ToNot(HaveOccurred())
		g.Expect(books).To(HaveLen(1))
		g.Expect(books[0].Authors).To(BeEmpty())
	})

	t.Run("handles BOM prefix", func(t *testing.T) {
		csvContent := "\ufefftitle;isbn;authors;description\nBook;123;;Desc"

		reader := strings.NewReader(csvContent)
		books, err := parseBooksFromReader(reader)

		g.Expect(err).ToNot(HaveOccurred())
		g.Expect(books).To(HaveLen(1))
	})

	t.Run("skips invalid records", func(t *testing.T) {
		csvContent := `title;isbn;authors;description
Book1;123
Book2;456;author;Description`

		reader := strings.NewReader(csvContent)
		books, err := parseBooksFromReader(reader)

		g.Expect(err).ToNot(HaveOccurred())
		g.Expect(books).To(HaveLen(1))
		g.Expect(books[0].Title).To(Equal("Book2"))
	})
}

func TestParseMagazinesFromReader(t *testing.T) {
	g := NewGomegaWithT(t)

	t.Run("parses magazines correctly", func(t *testing.T) {
		csvContent := `title;isbn;authors;publishedAt
Beautiful cooking;5454-5587-3210;null-walter@echocat.org;21.05.2011
Vinum;1313-4545-8875;null-gustafsson@echocat.org;23.02.2012`

		reader := strings.NewReader(csvContent)
		magazines, err := parseMagazinesFromReader(reader)

		g.Expect(err).ToNot(HaveOccurred())
		g.Expect(magazines).To(HaveLen(2))
		
		expectedDate1, _ := time.Parse("02.01.2006", "21.05.2011")
		g.Expect(magazines[0]).To(Equal(Magazine{
			Title:       "Beautiful cooking",
			ISBN:        "5454-5587-3210",
			Authors:     []string{"null-walter@echocat.org"},
			PublishedAt: expectedDate1,
		}))
		
		expectedDate2, _ := time.Parse("02.01.2006", "23.02.2012")
		g.Expect(magazines[1]).To(Equal(Magazine{
			Title:       "Vinum",
			ISBN:        "1313-4545-8875",
			Authors:     []string{"null-gustafsson@echocat.org"},
			PublishedAt: expectedDate2,
		}))
	})

	t.Run("handles invalid date format", func(t *testing.T) {
		csvContent := `title;isbn;authors;publishedAt
Magazine;123;author;invalid-date`

		reader := strings.NewReader(csvContent)
		magazines, err := parseMagazinesFromReader(reader)

		g.Expect(err).ToNot(HaveOccurred())
		g.Expect(magazines).To(HaveLen(1))
		g.Expect(magazines[0].PublishedAt).To(Equal(time.Time{}))
	})

	t.Run("handles multiple authors", func(t *testing.T) {
		csvContent := `title;isbn;authors;publishedAt
Magazine;123;auth1@test.com,auth2@test.com;01.01.2020`

		reader := strings.NewReader(csvContent)
		magazines, err := parseMagazinesFromReader(reader)

		g.Expect(err).ToNot(HaveOccurred())
		g.Expect(magazines).To(HaveLen(1))
		g.Expect(magazines[0].Authors).To(Equal([]string{"auth1@test.com", "auth2@test.com"}))
	})
}