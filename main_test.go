package main

import (
	. "github.com/onsi/gomega"
	"strings"
	"testing"
)

func TestLoadAuthors(t *testing.T) {
	g := NewGomegaWithT(t)

	// when
	authors, err := loadAuthors()

	// then
	g.Expect(err).To(BeNil())
	g.Expect(authors).To(HaveKey("null-walter@echocat.org"))
	g.Expect(authors["null-walter@echocat.org"].FirstName).To(Equal("Paul"))
	g.Expect(authors["null-walter@echocat.org"].LastName).To(Equal("Walter"))
}

func TestLoadBooks(t *testing.T) {
	g := NewGomegaWithT(t)

	// when
	books, err := loadBooks()

	// then
	g.Expect(err).To(BeNil())
	g.Expect(books).To(HaveLen(8))
	g.Expect(findBookByISBN(books, "5554-5545-4518")).ToNot(BeNil())
}

func TestLoadMagazines(t *testing.T) {
	g := NewGomegaWithT(t)

	// when
	magazines, err := loadMagazines()

	// then
	g.Expect(err).To(BeNil())
	g.Expect(magazines).To(HaveLen(6))
	g.Expect(findMagazineByISBN(magazines, "5454-5587-3210")).ToNot(BeNil())
}

func TestFilterBooksByISBN(t *testing.T) {
	g := NewGomegaWithT(t)

	// given
	books, err := loadBooks()
	g.Expect(err).To(BeNil())

	// when
	filtered := filterBooksByISBN(books, "5554-5545-4518")

	// then
	g.Expect(filtered).To(HaveLen(1))
	g.Expect(filtered[0].ISBN).To(Equal("5554-5545-4518"))
}

func TestFilterMagazinesByISBN(t *testing.T) {
	g := NewGomegaWithT(t)

	// given
	magazines, err := loadMagazines()
	g.Expect(err).To(BeNil())

	// when
	filtered := filterMagazinesByISBN(magazines, "5454-5587-3210")

	// then
	g.Expect(filtered).To(HaveLen(1))
	g.Expect(filtered[0].ISBN).To(Equal("5454-5587-3210"))
}

func TestFilterBooksByAuthorEmail(t *testing.T) {
	g := NewGomegaWithT(t)

	// given
	books, err := loadBooks()
	g.Expect(err).To(BeNil())

	// when
	filtered := filterBooksByAuthorEmail(books, "null-walter@echocat.org")

	// then
	g.Expect(filtered).ToNot(BeEmpty())
	for _, book := range filtered {
		g.Expect(book.AuthorEmails).To(ContainElement("null-walter@echocat.org"))
	}
}

func TestFilterMagazinesByAuthorEmail(t *testing.T) {
	g := NewGomegaWithT(t)

	// given
	magazines, err := loadMagazines()
	g.Expect(err).To(BeNil())

	// when
	filtered := filterMagazinesByAuthorEmail(magazines, "null-walter@echocat.org")

	// then
	g.Expect(filtered).ToNot(BeEmpty())
	for _, mag := range filtered {
		g.Expect(mag.AuthorEmails).To(ContainElement("null-walter@echocat.org"))
	}
}

func TestBuildCombinedItemsSortedByTitle(t *testing.T) {
	g := NewGomegaWithT(t)

	// given
	authors, err := loadAuthors()
	g.Expect(err).To(BeNil())
	books, err := loadBooks()
	g.Expect(err).To(BeNil())
	magazines, err := loadMagazines()
	g.Expect(err).To(BeNil())

	// when
	items := buildCombinedItems(books, magazines, authors)

	// then
	g.Expect(items).To(HaveLen(len(books) + len(magazines)))
	g.Expect(isSortedByTitle(items)).To(BeTrue())
}

func TestBuildPageDataFiltersAndCounts(t *testing.T) {
	g := NewGomegaWithT(t)

	// given
	data := sampleLibraryData()

	// when
	pageData := buildPageData(data, "111", "")

	// then
	g.Expect(pageData.BookCount).To(Equal(1))
	g.Expect(pageData.MagazineCount).To(Equal(0))
	g.Expect(pageData.HasSearch).To(BeTrue())
	g.Expect(pageData.NoResults).To(BeFalse())
	g.Expect(pageData.Items).To(HaveLen(1))
	g.Expect(pageData.Items[0].ItemType).To(Equal("Book"))

	// when
	pageData = buildPageData(data, "", "b@library.test")

	// then
	g.Expect(pageData.BookCount).To(Equal(1))
	g.Expect(pageData.MagazineCount).To(Equal(1))
	g.Expect(pageData.Items).To(HaveLen(2))
	g.Expect(isSortedByTitle(pageData.Items)).To(BeTrue())

	// when
	pageData = buildPageData(data, "999", "a@library.test")

	// then
	g.Expect(pageData.NoResults).To(BeTrue())
}

func TestBuildPageDataCombinesFilters(t *testing.T) {
	g := NewGomegaWithT(t)

	// given
	data := sampleLibraryData()

	// when
	pageData := buildPageData(data, "333", "a@library.test")

	// then
	g.Expect(pageData.BookCount).To(Equal(0))
	g.Expect(pageData.MagazineCount).To(Equal(1))
	g.Expect(pageData.Items).To(HaveLen(1))
	g.Expect(pageData.Items[0].ItemType).To(Equal("Magazine"))
}

func TestValidateAddRequestAuthorOnly(t *testing.T) {
	g := NewGomegaWithT(t)

	// given
	data := sampleLibraryData()
	request := AddRequest{
		ItemType:       "author",
		NewAuthorEmail: "new@library.test",
		NewAuthorFirst: "New",
		NewAuthorLast:  "Author",
	}

	// when
	plan, err := validateAddRequest(data, request)

	// then
	g.Expect(err).To(BeNil())
	g.Expect(plan.AddAuthor).ToNot(BeNil())
	g.Expect(plan.AddBook).To(BeNil())
	g.Expect(plan.AddMagazine).To(BeNil())
}

func TestValidateAddRequestBookWithNewAuthor(t *testing.T) {
	g := NewGomegaWithT(t)

	// given
	data := sampleLibraryData()
	request := AddRequest{
		ItemType:       "book",
		Title:          "Omega",
		ISBN:           "999",
		AuthorEmails:   []string{"new@library.test"},
		Description:    "A new book",
		NewAuthorEmail: "new@library.test",
		NewAuthorFirst: "New",
		NewAuthorLast:  "Author",
	}

	// when
	plan, err := validateAddRequest(data, request)

	// then
	g.Expect(err).To(BeNil())
	g.Expect(plan.AddBook).ToNot(BeNil())
	g.Expect(plan.AddAuthor).ToNot(BeNil())
}

func TestValidateAddRequestMissingAuthorDetails(t *testing.T) {
	g := NewGomegaWithT(t)

	// given
	data := sampleLibraryData()
	request := AddRequest{
		ItemType:     "book",
		Title:        "Omega",
		ISBN:         "999",
		AuthorEmails: []string{"missing@library.test"},
		Description:  "A new book",
	}

	// when
	_, err := validateAddRequest(data, request)

	// then
	g.Expect(err).ToNot(BeNil())
}

func TestValidateAddRequestDuplicateISBN(t *testing.T) {
	g := NewGomegaWithT(t)

	// given
	data := sampleLibraryData()
	request := AddRequest{
		ItemType:     "magazine",
		Title:        "Duplicate",
		ISBN:         "111",
		AuthorEmails: []string{"a@library.test"},
		PublishedAt:  "03.03.2020",
	}

	// when
	_, err := validateAddRequest(data, request)

	// then
	g.Expect(err).ToNot(BeNil())
}

func sampleLibraryData() LibraryData {
	authors := map[string]Author{
		"a@library.test": {Email: "a@library.test", FirstName: "Ada", LastName: "Lovelace"},
		"b@library.test": {Email: "b@library.test", FirstName: "Grace", LastName: "Hopper"},
	}
	books := []Book{
		{Title: "Alpha", ISBN: "111", AuthorEmails: []string{"a@library.test"}, Description: "Book A"},
		{Title: "Beta", ISBN: "222", AuthorEmails: []string{"b@library.test"}, Description: "Book B"},
	}
	magazines := []Magazine{
		{Title: "Gamma", ISBN: "333", AuthorEmails: []string{"a@library.test"}, PublishedAt: "01.01.2020"},
		{Title: "Delta", ISBN: "444", AuthorEmails: []string{"b@library.test"}, PublishedAt: "02.02.2020"},
	}

	return LibraryData{
		Authors:   authors,
		Books:     books,
		Magazines: magazines,
	}
}

func isSortedByTitle(items []ItemView) bool {
	if len(items) < 2 {
		return true
	}
	prev := normalizeTitle(items[0].Title)
	for i := 1; i < len(items); i++ {
		current := normalizeTitle(items[i].Title)
		if current < prev {
			return false
		}
		prev = current
	}
	return true
}

func normalizeTitle(title string) string {
	return strings.ToLower(strings.TrimSpace(title))
}

func findBookByISBN(books []Book, isbn string) *Book {
	for _, book := range books {
		if book.ISBN == isbn {
			return &book
		}
	}
	return nil
}

func findMagazineByISBN(magazines []Magazine, isbn string) *Magazine {
	for _, mag := range magazines {
		if mag.ISBN == isbn {
			return &mag
		}
	}
	return nil
}
