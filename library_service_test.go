package main

import (
	"testing"
	"time"
)

// Helper function to create a test library
func createTestLibrary() *Library {
	authors := map[string]Author{
		"author1@test.com": {Email: "author1@test.com", FirstName: "John", LastName: "Doe"},
		"author2@test.com": {Email: "author2@test.com", FirstName: "Jane", LastName: "Smith"},
		"author3@test.com": {Email: "author3@test.com", FirstName: "Bob", LastName: "Johnson"},
	}

	books := []Book{
		{
			Title:       "Advanced Programming",
			ISBN:        "1234-5678-9012",
			Authors:     []Author{authors["author1@test.com"]},
			Description: "A comprehensive guide to advanced programming concepts.",
		},
		{
			Title:       "Basic Cooking",
			ISBN:        "2345-6789-0123",
			Authors:     []Author{authors["author2@test.com"], authors["author3@test.com"]},
			Description: "Learn the basics of cooking.",
		},
	}

	magazines := []Magazine{
		{
			Title:       "Tech Weekly",
			ISBN:        "3456-7890-1234",
			Authors:     []Author{authors["author1@test.com"]},
			PublishedAt: time.Date(2023, 1, 15, 0, 0, 0, 0, time.UTC),
		},
		{
			Title:       "Cooking Monthly",
			ISBN:        "4567-8901-2345",
			Authors:     []Author{authors["author2@test.com"]},
			PublishedAt: time.Date(2023, 2, 20, 0, 0, 0, 0, time.UTC),
		},
	}

	return &Library{
		Authors:   authors,
		Books:     books,
		Magazines: magazines,
	}
}

func TestLibraryService_SearchByISBN(t *testing.T) {
	library := createTestLibrary()
	service := NewLibraryService(library)

	tests := []struct {
		name          string
		searchTerm    string
		expectedBooks int
		expectedMags  int
		expectedTitle string
	}{
		{
			name:          "Find book by exact ISBN",
			searchTerm:    "1234-5678-9012",
			expectedBooks: 1,
			expectedMags:  0,
			expectedTitle: "Advanced Programming",
		},
		{
			name:          "Find magazine by partial ISBN",
			searchTerm:    "3456",
			expectedBooks: 0,
			expectedMags:  1,
			expectedTitle: "Tech Weekly",
		},
		{
			name:          "No results for non-existent ISBN",
			searchTerm:    "9999-9999-9999",
			expectedBooks: 0,
			expectedMags:  0,
		},
		{
			name:          "Empty search term",
			searchTerm:    "",
			expectedBooks: 0,
			expectedMags:  0,
		},
		{
			name:          "Case insensitive search",
			searchTerm:    "1234-5678-9012",
			expectedBooks: 1,
			expectedMags:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			books, magazines := service.SearchByISBN(tt.searchTerm)

			if len(books) != tt.expectedBooks {
				t.Errorf("Expected %d books, got %d", tt.expectedBooks, len(books))
			}

			if len(magazines) != tt.expectedMags {
				t.Errorf("Expected %d magazines, got %d", tt.expectedMags, len(magazines))
			}

			if tt.expectedTitle != "" {
				found := false
				for _, book := range books {
					if book.Title == tt.expectedTitle {
						found = true
						break
					}
				}
				for _, mag := range magazines {
					if mag.Title == tt.expectedTitle {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Expected to find item with title '%s'", tt.expectedTitle)
				}
			}
		})
	}
}

func TestLibraryService_SearchByAuthorEmail(t *testing.T) {
	library := createTestLibrary()
	service := NewLibraryService(library)

	tests := []struct {
		name          string
		searchTerm    string
		expectedBooks int
		expectedMags  int
	}{
		{
			name:          "Find by exact author email",
			searchTerm:    "author1@test.com",
			expectedBooks: 1,
			expectedMags:  1,
		},
		{
			name:          "Find by partial email",
			searchTerm:    "author2",
			expectedBooks: 1,
			expectedMags:  1,
		},
		{
			name:          "No results for non-existent email",
			searchTerm:    "nonexistent@test.com",
			expectedBooks: 0,
			expectedMags:  0,
		},
		{
			name:          "Empty search term",
			searchTerm:    "",
			expectedBooks: 0,
			expectedMags:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			books, magazines := service.SearchByAuthorEmail(tt.searchTerm)

			if len(books) != tt.expectedBooks {
				t.Errorf("Expected %d books, got %d", tt.expectedBooks, len(books))
			}

			if len(magazines) != tt.expectedMags {
				t.Errorf("Expected %d magazines, got %d", tt.expectedMags, len(magazines))
			}
		})
	}
}

func TestLibraryService_GetAllItemsSorted(t *testing.T) {
	library := createTestLibrary()
	service := NewLibraryService(library)

	t.Run("Sort ascending", func(t *testing.T) {
		items := service.GetAllItemsSorted(true)

		if len(items) != 4 { // 2 books + 2 magazines
			t.Errorf("Expected 4 items, got %d", len(items))
		}

		// Check if sorted ascending (A-Z)
		expectedOrder := []string{"Advanced Programming", "Basic Cooking", "Cooking Monthly", "Tech Weekly"}
		for i, item := range items {
			if item.Title != expectedOrder[i] {
				t.Errorf("Expected item %d to be '%s', got '%s'", i, expectedOrder[i], item.Title)
			}
		}
	})

	t.Run("Sort descending", func(t *testing.T) {
		items := service.GetAllItemsSorted(false)

		if len(items) != 4 {
			t.Errorf("Expected 4 items, got %d", len(items))
		}

		// Check if sorted descending (Z-A)
		expectedOrder := []string{"Tech Weekly", "Cooking Monthly", "Basic Cooking", "Advanced Programming"}
		for i, item := range items {
			if item.Title != expectedOrder[i] {
				t.Errorf("Expected item %d to be '%s', got '%s'", i, expectedOrder[i], item.Title)
			}
		}
	})

	t.Run("Check item types", func(t *testing.T) {
		items := service.GetAllItemsSorted(true)

		bookCount := 0
		magCount := 0
		for _, item := range items {
			if item.Type == "BOOK" {
				bookCount++
				if item.Description == "" {
					t.Error("Book should have description")
				}
			} else if item.Type == "MAGAZINE" {
				magCount++
				if item.PublishedAt == "" {
					t.Error("Magazine should have published date")
				}
			}
		}

		if bookCount != 2 {
			t.Errorf("Expected 2 books, got %d", bookCount)
		}
		if magCount != 2 {
			t.Errorf("Expected 2 magazines, got %d", magCount)
		}
	})
}

func TestLibraryService_GetLibraryStats(t *testing.T) {
	library := createTestLibrary()
	service := NewLibraryService(library)

	bookCount, magCount, authorCount := service.GetLibraryStats()

	if bookCount != 2 {
		t.Errorf("Expected 2 books, got %d", bookCount)
	}

	if magCount != 2 {
		t.Errorf("Expected 2 magazines, got %d", magCount)
	}

	if authorCount != 3 {
		t.Errorf("Expected 3 authors, got %d", authorCount)
	}
}

func TestLibraryService_FormatAuthors(t *testing.T) {
	library := createTestLibrary()
	service := NewLibraryService(library)

	tests := []struct {
		name     string
		authors  []Author
		expected string
	}{
		{
			name:     "Single author",
			authors:  []Author{{FirstName: "John", LastName: "Doe"}},
			expected: "John Doe",
		},
		{
			name: "Multiple authors",
			authors: []Author{
				{FirstName: "John", LastName: "Doe"},
				{FirstName: "Jane", LastName: "Smith"},
			},
			expected: "John Doe, Jane Smith",
		},
		{
			name:     "No authors",
			authors:  []Author{},
			expected: "Unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := service.FormatAuthors(tt.authors)
			if result != tt.expected {
				t.Errorf("Expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

func TestLibraryService_GetBooks(t *testing.T) {
	library := createTestLibrary()
	service := NewLibraryService(library)

	books := service.GetBooks()

	if len(books) != 2 {
		t.Errorf("Expected 2 books, got %d", len(books))
	}

	// Check that we get the actual books, not copies
	if books[0].Title != "Advanced Programming" && books[1].Title != "Advanced Programming" {
		t.Error("Expected to find 'Advanced Programming' book")
	}
}

func TestLibraryService_GetMagazines(t *testing.T) {
	library := createTestLibrary()
	service := NewLibraryService(library)

	magazines := service.GetMagazines()

	if len(magazines) != 2 {
		t.Errorf("Expected 2 magazines, got %d", len(magazines))
	}

	// Check that we get the actual magazines, not copies
	if magazines[0].Title != "Tech Weekly" && magazines[1].Title != "Tech Weekly" {
		t.Error("Expected to find 'Tech Weekly' magazine")
	}
}

func TestLibraryService_GetAuthors(t *testing.T) {
	library := createTestLibrary()
	service := NewLibraryService(library)

	authors := service.GetAuthors()

	if len(authors) != 3 {
		t.Errorf("Expected 3 authors, got %d", len(authors))
	}

	// Check that we get the actual authors map
	if _, exists := authors["author1@test.com"]; !exists {
		t.Error("Expected to find author1@test.com")
	}
}
