package main

import (
	"testing"

	. "github.com/onsi/gomega"
)

func TestMainFunctionality(t *testing.T) {
	g := NewGomegaWithT(t)

	// Test that we can create a library service
	// This is a basic integration test to ensure main components work together
	library, err := LoadLibrary()

	// If CSV files don't exist, skip the test
	if err != nil {
		t.Skipf("Skipping test - CSV files not available: %v", err)
		return
	}

	// Verify library was loaded
	g.Expect(library).ToNot(BeNil())
	g.Expect(len(library.Authors)).To(BeNumerically(">", 0))
	g.Expect(len(library.Books)).To(BeNumerically(">", 0))
	g.Expect(len(library.Magazines)).To(BeNumerically(">", 0))

	// Test that we can create a service
	service := NewLibraryService(library)
	g.Expect(service).ToNot(BeNil())

	// Test that we can create a UI handler
	uiHandler := NewUIHandler(library)
	g.Expect(uiHandler).ToNot(BeNil())
}

func TestLibraryServiceIntegration(t *testing.T) {
	g := NewGomegaWithT(t)

	// Create a minimal test library
	library := &Library{
		Authors: map[string]Author{
			"test@example.com": {Email: "test@example.com", FirstName: "Test", LastName: "Author"},
		},
		Books: []Book{
			{Title: "Test Book", ISBN: "123-456", Authors: []Author{{Email: "test@example.com", FirstName: "Test", LastName: "Author"}}, Description: "A test book"},
		},
		Magazines: []Magazine{
			{Title: "Test Magazine", ISBN: "789-012", Authors: []Author{{Email: "test@example.com", FirstName: "Test", LastName: "Author"}}},
		},
	}

	service := NewLibraryService(library)

	// Test basic functionality
	bookCount, magCount, authorCount := service.GetLibraryStats()
	g.Expect(bookCount).To(Equal(1))
	g.Expect(magCount).To(Equal(1))
	g.Expect(authorCount).To(Equal(1))

	// Test search functionality
	books, magazines := service.SearchByISBN("123")
	g.Expect(len(books)).To(Equal(1))
	g.Expect(len(magazines)).To(Equal(0))

	books, magazines = service.SearchByAuthorEmail("test@example.com")
	g.Expect(len(books)).To(Equal(1))
	g.Expect(len(magazines)).To(Equal(1))

	// Test sorting
	items := service.GetAllItemsSorted(true)
	g.Expect(len(items)).To(Equal(2))
	g.Expect(items[0].Title).To(Equal("Test Book")) // Should come first alphabetically
	g.Expect(items[1].Title).To(Equal("Test Magazine"))
}
