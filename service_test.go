package main

import (
	"testing"
)

func TestProcessRequest(t *testing.T) {
	// Setup dummy data
	author := &Author{Email: "test@example.com", Firstname: "Test", Lastname: "User"}
	authors := map[string]*Author{"test@example.com": author}

	books := []*Book{
		{Title: "Book A", ISBN: "111", Authors: []*Author{author}, Description: "Desc A"},
		{Title: "Book B", ISBN: "222", Authors: []*Author{author}, Description: "Desc B"},
	}

	mags := []*Magazine{
		{Title: "Mag A", ISBN: "333", Authors: []*Author{author}, PublishedAt: "2021-01-01"},
		{Title: "Mag B", ISBN: "111", Authors: []*Author{author}, PublishedAt: "2021-02-01"},
	}

	data := &LibraryData{
		Authors:   authors,
		Books:     books,
		Magazines: mags,
	}

	// Test 1: ISBN Search
	t.Run("ISBN Search", func(t *testing.T) {
		result := ProcessRequest(data, "111", "", "", "")
		if len(result.Books) != 1 || result.Books[0].Title != "Book A" {
			t.Errorf("Expected 1 book (Book A), got %v", result.Books)
		}
		if len(result.Magazines) != 1 || result.Magazines[0].Title != "Mag B" {
			t.Errorf("Expected 1 magazine (Mag B), got %v", result.Magazines)
		}
	})

	// Test 2: Email Search
	t.Run("Email Search", func(t *testing.T) {
		result := ProcessRequest(data, "", "test@example.com", "", "")
		if len(result.Books) != 2 {
			t.Errorf("Expected 2 books, got %d", len(result.Books))
		}
		if len(result.Magazines) != 2 {
			t.Errorf("Expected 2 magazines, got %d", len(result.Magazines))
		}
	})

	// Test 3: Sort Title Ascending
	t.Run("Sort Title Ascending", func(t *testing.T) {
		result := ProcessRequest(data, "", "", "title", "asc")
		if len(result.SortedItems) != 4 {
			t.Fatalf("Expected 4 sorted items, got %d", len(result.SortedItems))
		}
		if result.SortedItems[0].Title != "Book A" {
			t.Errorf("Expected first item Book A, got %s", result.SortedItems[0].Title)
		}
		if result.SortedItems[3].Title != "Mag B" {
			t.Errorf("Expected last item Mag B, got %s", result.SortedItems[3].Title)
		}
	})

	// Test 4: Sort Title Descending
	t.Run("Sort Title Descending", func(t *testing.T) {
		result := ProcessRequest(data, "", "", "title", "desc")
		if len(result.SortedItems) != 4 {
			t.Fatalf("Expected 4 sorted items, got %d", len(result.SortedItems))
		}
		if result.SortedItems[0].Title != "Mag B" {
			t.Errorf("Expected first item Mag B, got %s", result.SortedItems[0].Title)
		}
		if result.SortedItems[3].Title != "Book A" {
			t.Errorf("Expected last item Book A, got %s", result.SortedItems[3].Title)
		}
	})

	// Test 5: Search AND Sort
	t.Run("Search ISBN and Sort", func(t *testing.T) {
		// Matches Book A (111) and Mag B (111)
		result := ProcessRequest(data, "111", "", "title", "asc")
		if len(result.SortedItems) != 2 {
			t.Fatalf("Expected 2 sorted items, got %d", len(result.SortedItems))
		}
		if result.SortedItems[0].Title != "Book A" {
			t.Errorf("Expected first item Book A, got %s", result.SortedItems[0].Title)
		}
		if result.SortedItems[1].Title != "Mag B" {
			t.Errorf("Expected second item Mag B, got %s", result.SortedItems[1].Title)
		}
	})
}
