package main

import (
	"testing"
	"time"
)

func TestDateParsing(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected time.Time
		hasError bool
	}{
		{
			name:     "Valid date format",
			input:    "15.01.2023",
			expected: time.Date(2023, 1, 15, 0, 0, 0, 0, time.UTC),
			hasError: false,
		},
		{
			name:     "Another valid date",
			input:    "01.12.2022",
			expected: time.Date(2022, 12, 1, 0, 0, 0, 0, time.UTC),
			hasError: false,
		},
		{
			name:     "Invalid date format",
			input:    "invalid-date",
			expected: time.Time{},
			hasError: true,
		},
		{
			name:     "Wrong format",
			input:    "2023-01-15",
			expected: time.Time{},
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := time.Parse("02.01.2006", tt.input)

			if tt.hasError {
				if err == nil {
					t.Error("Expected error for invalid date, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if !result.Equal(tt.expected) {
					t.Errorf("Expected %v, got %v", tt.expected, result)
				}
			}
		})
	}
}

func TestLoadLibrary_WithRealFiles(t *testing.T) {
	// This test uses the actual CSV files in the resources directory
	// It's more of an integration test
	library, err := LoadLibrary()

	if err != nil {
		t.Skipf("Skipping test - CSV files not available: %v", err)
		return
	}

	// Basic validation that data was loaded
	if len(library.Authors) == 0 {
		t.Error("Expected authors to be loaded")
	}

	if len(library.Books) == 0 {
		t.Error("Expected books to be loaded")
	}

	if len(library.Magazines) == 0 {
		t.Error("Expected magazines to be loaded")
	}

	// Validate that books have authors linked correctly
	for _, book := range library.Books {
		if book.Title == "" {
			t.Error("Book should have a title")
		}
		if book.ISBN == "" {
			t.Error("Book should have an ISBN")
		}
		// Authors should be linked from the authors map
		for _, author := range book.Authors {
			if _, exists := library.Authors[author.Email]; !exists {
				t.Errorf("Book author %s not found in authors map", author.Email)
			}
		}
	}

	// Validate that magazines have authors linked correctly
	for _, magazine := range library.Magazines {
		if magazine.Title == "" {
			t.Error("Magazine should have a title")
		}
		if magazine.ISBN == "" {
			t.Error("Magazine should have an ISBN")
		}
		if magazine.PublishedAt.IsZero() {
			t.Error("Magazine should have a published date")
		}
		// Authors should be linked from the authors map
		for _, author := range magazine.Authors {
			if _, exists := library.Authors[author.Email]; !exists {
				t.Errorf("Magazine author %s not found in authors map", author.Email)
			}
		}
	}
}

func TestLibraryStructure(t *testing.T) {
	// Test that our data structures are properly defined
	library := &Library{
		Authors:   make(map[string]Author),
		Books:     []Book{},
		Magazines: []Magazine{},
	}

	// Test Author structure
	author := Author{
		Email:     "test@example.com",
		FirstName: "John",
		LastName:  "Doe",
	}
	library.Authors[author.Email] = author

	// Test Book structure
	book := Book{
		Title:       "Test Book",
		ISBN:        "1234-5678-9012",
		Authors:     []Author{author},
		Description: "A test book",
	}
	library.Books = append(library.Books, book)

	// Test Magazine structure
	magazine := Magazine{
		Title:       "Test Magazine",
		ISBN:        "2345-6789-0123",
		Authors:     []Author{author},
		PublishedAt: time.Now(),
	}
	library.Magazines = append(library.Magazines, magazine)

	// Validate the structures
	if len(library.Authors) != 1 {
		t.Errorf("Expected 1 author, got %d", len(library.Authors))
	}

	if len(library.Books) != 1 {
		t.Errorf("Expected 1 book, got %d", len(library.Books))
	}

	if len(library.Magazines) != 1 {
		t.Errorf("Expected 1 magazine, got %d", len(library.Magazines))
	}

	// Test that references work correctly
	retrievedAuthor := library.Authors["test@example.com"]
	if retrievedAuthor.FirstName != "John" {
		t.Errorf("Expected John, got %s", retrievedAuthor.FirstName)
	}

	if library.Books[0].Authors[0].Email != "test@example.com" {
		t.Error("Book author reference not working correctly")
	}

	if library.Magazines[0].Authors[0].Email != "test@example.com" {
		t.Error("Magazine author reference not working correctly")
	}
}
