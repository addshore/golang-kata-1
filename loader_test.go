package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoader(t *testing.T) {
	// Create temporary directory for test files
	tmpDir, err := os.MkdirTemp("", "library-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create dummy authors.csv
	authorsContent := `email;firstname;lastname
test@example.com;Test;User`
	authorsPath := filepath.Join(tmpDir, "authors.csv")
	if err := os.WriteFile(authorsPath, []byte(authorsContent), 0644); err != nil {
		t.Fatal(err)
	}

	// Create dummy books.csv
	booksContent := `title;isbn;authors;description
Test Book;123-456;test@example.com;Description here`
	booksPath := filepath.Join(tmpDir, "books.csv")
	if err := os.WriteFile(booksPath, []byte(booksContent), 0644); err != nil {
		t.Fatal(err)
	}

	// Create dummy magazines.csv
	magsContent := `title;isbn;authors;publishedAt
Test Mag;987-654;test@example.com;21.05.2011`
	magsPath := filepath.Join(tmpDir, "magazines.csv")
	if err := os.WriteFile(magsPath, []byte(magsContent), 0644); err != nil {
		t.Fatal(err)
	}

	// Test LoadAuthors
	authors, err := LoadAuthors(authorsPath)
	if err != nil {
		t.Fatalf("LoadAuthors failed: %v", err)
	}
	if len(authors) != 1 {
		t.Errorf("Expected 1 author, got %d", len(authors))
	}
	if authors["test@example.com"].Firstname != "Test" {
		t.Errorf("Expected author firstname Test, got %s", authors["test@example.com"].Firstname)
	}

	// Test LoadBooks
	books, err := LoadBooks(booksPath, authors)
	if err != nil {
		t.Fatalf("LoadBooks failed: %v", err)
	}
	if len(books) != 1 {
		t.Errorf("Expected 1 book, got %d", len(books))
	}
	if books[0].Title != "Test Book" {
		t.Errorf("Expected book title Test Book, got %s", books[0].Title)
	}
	if len(books[0].Authors) != 1 || books[0].Authors[0].Email != "test@example.com" {
		t.Errorf("Author linking failed for Book")
	}

	// Test LoadMagazines
	mags, err := LoadMagazines(magsPath, authors)
	if err != nil {
		t.Fatalf("LoadMagazines failed: %v", err)
	}
	if len(mags) != 1 {
		t.Errorf("Expected 1 magazine, got %d", len(mags))
	}
	if mags[0].Title != "Test Mag" {
		t.Errorf("Expected mag title Test Mag, got %s", mags[0].Title)
	}
	if len(mags[0].Authors) != 1 || mags[0].Authors[0].Email != "test@example.com" {
		t.Errorf("Author linking failed for Magazine")
	}
}
