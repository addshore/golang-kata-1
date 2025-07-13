package main

import (
	. "github.com/onsi/gomega"
	"testing"
)

func TestWelcomeMessage(t *testing.T) {
	g := NewGomegaWithT(t)

	// given
	expected := "Hello world!"

	// than
	actual := welcomeMessage()

	// that
	g.Expect(actual).To(Equal(expected))
}

func TestLibraryCreation(t *testing.T) {
	g := NewGomegaWithT(t)
	
	// when
	library := NewLibrary()
	
	// then
	g.Expect(library).NotTo(BeNil())
	g.Expect(library.Authors).NotTo(BeNil())
	g.Expect(library.Books).NotTo(BeNil())
	g.Expect(library.Magazines).NotTo(BeNil())
	g.Expect(len(library.Authors)).To(Equal(0))
	g.Expect(len(library.Books)).To(Equal(0))
	g.Expect(len(library.Magazines)).To(Equal(0))
}

func TestFormatAuthors(t *testing.T) {
	g := NewGomegaWithT(t)
	
	tests := []struct {
		name     string
		authors  []Author
		expected string
	}{
		{
			name:     "Empty authors",
			authors:  []Author{},
			expected: "Unknown",
		},
		{
			name: "Single author",
			authors: []Author{
				{FirstName: "John", LastName: "Doe", Email: "john@example.com"},
			},
			expected: "John Doe",
		},
		{
			name: "Multiple authors",
			authors: []Author{
				{FirstName: "John", LastName: "Doe", Email: "john@example.com"},
				{FirstName: "Jane", LastName: "Smith", Email: "jane@example.com"},
			},
			expected: "John Doe, Jane Smith",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// when
			result := FormatAuthors(tt.authors)
			
			// then
			g.Expect(result).To(Equal(tt.expected))
		})
	}
}

func TestLibraryFindBookByISBN(t *testing.T) {
	g := NewGomegaWithT(t)
	
	// given
	library := NewLibrary()
	library.Books = []Book{
		{
			Title: "Test Book 1",
			ISBN:  "123-456-789",
			Authors: []Author{
				{FirstName: "John", LastName: "Doe", Email: "john@example.com"},
			},
			Description: "A test book",
		},
		{
			Title: "Test Book 2",
			ISBN:  "987-654-321",
			Authors: []Author{
				{FirstName: "Jane", LastName: "Smith", Email: "jane@example.com"},
			},
			Description: "Another test book",
		},
	}
	
	// when - found case
	foundBook := library.FindBookByISBN("123-456-789")
	
	// then
	g.Expect(foundBook).NotTo(BeNil())
	g.Expect(foundBook.Title).To(Equal("Test Book 1"))
	g.Expect(foundBook.ISBN).To(Equal("123-456-789"))
	
	// when - not found case
	notFoundBook := library.FindBookByISBN("000-000-000")
	
	// then
	g.Expect(notFoundBook).To(BeNil())
}

func TestLibraryFindMagazineByISBN(t *testing.T) {
	g := NewGomegaWithT(t)
	
	// given
	library := NewLibrary()
	library.Magazines = []Magazine{
		{
			Title: "Test Magazine 1",
			ISBN:  "111-222-333",
			Authors: []Author{
				{FirstName: "Alice", LastName: "Johnson", Email: "alice@example.com"},
			},
			PublishedAt: "2023-01-01",
		},
	}
	
	// when - found case
	foundMagazine := library.FindMagazineByISBN("111-222-333")
	
	// then
	g.Expect(foundMagazine).NotTo(BeNil())
	g.Expect(foundMagazine.Title).To(Equal("Test Magazine 1"))
	g.Expect(foundMagazine.ISBN).To(Equal("111-222-333"))
	
	// when - not found case
	notFoundMagazine := library.FindMagazineByISBN("000-000-000")
	
	// then
	g.Expect(notFoundMagazine).To(BeNil())
}

func TestLibraryFindBooksByAuthorEmail(t *testing.T) {
	g := NewGomegaWithT(t)
	
	// given
	library := NewLibrary()
	author1 := Author{FirstName: "John", LastName: "Doe", Email: "john@example.com"}
	author2 := Author{FirstName: "Jane", LastName: "Smith", Email: "jane@example.com"}
	
	library.Books = []Book{
		{
			Title:   "Book by John",
			ISBN:    "123-456-789",
			Authors: []Author{author1},
		},
		{
			Title:   "Book by Jane",
			ISBN:    "987-654-321",
			Authors: []Author{author2},
		},
		{
			Title:   "Book by John and Jane",
			ISBN:    "555-666-777",
			Authors: []Author{author1, author2},
		},
	}
	
	// when - find books by John
	johnBooks := library.FindBooksByAuthorEmail("john@example.com")
	
	// then
	g.Expect(len(johnBooks)).To(Equal(2))
	g.Expect(johnBooks[0].Title).To(Equal("Book by John"))
	g.Expect(johnBooks[1].Title).To(Equal("Book by John and Jane"))
	
	// when - find books by non-existent author
	noBooks := library.FindBooksByAuthorEmail("nobody@example.com")
	
	// then
	g.Expect(len(noBooks)).To(Equal(0))
}

func TestLibraryFindMagazinesByAuthorEmail(t *testing.T) {
	g := NewGomegaWithT(t)
	
	// given
	library := NewLibrary()
	author1 := Author{FirstName: "Alice", LastName: "Johnson", Email: "alice@example.com"}
	author2 := Author{FirstName: "Bob", LastName: "Wilson", Email: "bob@example.com"}
	
	library.Magazines = []Magazine{
		{
			Title:       "Magazine by Alice",
			ISBN:        "111-222-333",
			Authors:     []Author{author1},
			PublishedAt: "2023-01-01",
		},
		{
			Title:       "Magazine by Bob",
			ISBN:        "444-555-666",
			Authors:     []Author{author2},
			PublishedAt: "2023-02-01",
		},
		{
			Title:       "Magazine by Alice and Bob",
			ISBN:        "777-888-999",
			Authors:     []Author{author1, author2},
			PublishedAt: "2023-03-01",
		},
	}
	
	// when - find magazines by Alice
	aliceMagazines := library.FindMagazinesByAuthorEmail("alice@example.com")
	
	// then
	g.Expect(len(aliceMagazines)).To(Equal(2))
	g.Expect(aliceMagazines[0].Title).To(Equal("Magazine by Alice"))
	g.Expect(aliceMagazines[1].Title).To(Equal("Magazine by Alice and Bob"))
	
	// when - find magazines by non-existent author
	noMagazines := library.FindMagazinesByAuthorEmail("nobody@example.com")
	
	// then
	g.Expect(len(noMagazines)).To(Equal(0))
}

func TestLibraryGetAllItemsSorted(t *testing.T) {
	g := NewGomegaWithT(t)
	
	// given
	library := NewLibrary()
	author := Author{FirstName: "Test", LastName: "Author", Email: "test@example.com"}
	
	library.Books = []Book{
		{Title: "Zebra Book", ISBN: "123-456-789", Authors: []Author{author}, Description: "Last book"},
		{Title: "Alpha Book", ISBN: "987-654-321", Authors: []Author{author}, Description: "First book"},
	}
	
	library.Magazines = []Magazine{
		{Title: "Yankee Magazine", ISBN: "111-222-333", Authors: []Author{author}, PublishedAt: "2023-01-01"},
		{Title: "Beta Magazine", ISBN: "444-555-666", Authors: []Author{author}, PublishedAt: "2023-02-01"},
	}
	
	// when - ascending sort
	ascendingItems := library.GetAllItemsSorted(true)
	
	// then
	g.Expect(len(ascendingItems)).To(Equal(4))
	g.Expect(ascendingItems[0].Title).To(Equal("Alpha Book"))
	g.Expect(ascendingItems[1].Title).To(Equal("Beta Magazine"))
	g.Expect(ascendingItems[2].Title).To(Equal("Yankee Magazine"))
	g.Expect(ascendingItems[3].Title).To(Equal("Zebra Book"))
	
	// when - descending sort
	descendingItems := library.GetAllItemsSorted(false)
	
	// then
	g.Expect(len(descendingItems)).To(Equal(4))
	g.Expect(descendingItems[0].Title).To(Equal("Zebra Book"))
	g.Expect(descendingItems[1].Title).To(Equal("Yankee Magazine"))
	g.Expect(descendingItems[2].Title).To(Equal("Beta Magazine"))
	g.Expect(descendingItems[3].Title).To(Equal("Alpha Book"))
}

func TestLibraryLoadFromCSV(t *testing.T) {
	g := NewGomegaWithT(t)
	
	// given
	library := NewLibrary()
	
	// when
	err := library.LoadFromCSV("resources/authors.csv", "resources/books.csv", "resources/magazines.csv")
	
	// then
	g.Expect(err).To(BeNil())
	g.Expect(len(library.Authors)).To(BeNumerically(">", 0))
	g.Expect(len(library.Books)).To(BeNumerically(">", 0))
	g.Expect(len(library.Magazines)).To(BeNumerically(">", 0))
	
	// verify specific data
	g.Expect(library.Authors["null-walter@echocat.org"].FirstName).To(Equal("Paul"))
	g.Expect(library.Authors["null-walter@echocat.org"].LastName).To(Equal("Walter"))
}

func TestLibraryLoadFromCSVWithInvalidFiles(t *testing.T) {
	g := NewGomegaWithT(t)
	
	// given
	library := NewLibrary()
	
	// when
	err := library.LoadFromCSV("nonexistent.csv", "resources/books.csv", "resources/magazines.csv")
	
	// then
	g.Expect(err).NotTo(BeNil())
	g.Expect(err.Error()).To(ContainSubstring("failed to load authors"))
}

func TestLibraryAddAuthor(t *testing.T) {
	g := NewGomegaWithT(t)
	
	// given
	library := NewLibrary()
	author := Author{
		Email:     "test@example.com",
		FirstName: "Test",
		LastName:  "User",
	}
	
	// when
	library.AddAuthor(author)
	
	// then
	g.Expect(len(library.Authors)).To(Equal(1))
	g.Expect(library.Authors["test@example.com"]).To(Equal(author))
}

func TestLibraryAddBook(t *testing.T) {
	g := NewGomegaWithT(t)
	
	// given
	library := NewLibrary()
	author := Author{FirstName: "Test", LastName: "Author", Email: "test@example.com"}
	book := Book{
		Title:       "Test Book",
		ISBN:        "123-456-789",
		Authors:     []Author{author},
		Description: "A test book",
	}
	
	// when
	library.AddBook(book)
	
	// then
	g.Expect(len(library.Books)).To(Equal(1))
	g.Expect(library.Books[0]).To(Equal(book))
}

func TestLibraryAddMagazine(t *testing.T) {
	g := NewGomegaWithT(t)
	
	// given
	library := NewLibrary()
	author := Author{FirstName: "Test", LastName: "Author", Email: "test@example.com"}
	magazine := Magazine{
		Title:       "Test Magazine",
		ISBN:        "111-222-333",
		Authors:     []Author{author},
		PublishedAt: "01.01.2023",
	}
	
	// when
	library.AddMagazine(magazine)
	
	// then
	g.Expect(len(library.Magazines)).To(Equal(1))
	g.Expect(library.Magazines[0]).To(Equal(magazine))
}

func TestLibraryGetOrCreateAuthor(t *testing.T) {
	g := NewGomegaWithT(t)
	
	// given
	library := NewLibrary()
	existingAuthor := Author{
		Email:     "existing@example.com",
		FirstName: "Existing",
		LastName:  "Author",
	}
	library.AddAuthor(existingAuthor)
	
	// when - get existing author
	retrievedAuthor := library.GetOrCreateAuthor("existing@example.com", "Should", "Ignore")
	
	// then
	g.Expect(retrievedAuthor).To(Equal(existingAuthor))
	g.Expect(len(library.Authors)).To(Equal(1))
	
	// when - create new author
	newAuthor := library.GetOrCreateAuthor("new@example.com", "New", "Author")
	
	// then
	g.Expect(newAuthor.Email).To(Equal("new@example.com"))
	g.Expect(newAuthor.FirstName).To(Equal("New"))
	g.Expect(newAuthor.LastName).To(Equal("Author"))
	g.Expect(len(library.Authors)).To(Equal(2))
}
