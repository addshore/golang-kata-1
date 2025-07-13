package main

import (
	"testing"
)

func TestSearchByISBN(t *testing.T) {
	books := []Book{
		{Title: "Book1", Author: "Author1", ISBN: "123", Genre: "Genre1"},
		{Title: "Book2", Author: "Author2", ISBN: "456", Genre: "Genre2"},
	}
	magazines := []Magazine{
		{Title: "Magazine1", Publisher: "Publisher1", IssueNumber: "789", PublishedAt: "2025-01-01"},
	}

	results := searchByISBN(books, magazines, "123")
	if len(results) != 1 || results[0] != "Book - Title: Book1\nAuthor: Author1\nISBN: 123\nDescription: Genre1\n" {
		t.Errorf("Expected 1 result for ISBN '123', got %v", results)
	}

	results = searchByISBN(books, magazines, "789")
	if len(results) != 1 || results[0] != "Magazine - Title: Magazine1\nPublisher: Publisher1\nISBN: 789\nPublished At: 2025-01-01\n" {
		t.Errorf("Expected 1 result for ISBN '789', got %v", results)
	}

	results = searchByISBN(books, magazines, "000")
	if len(results) != 0 {
		t.Errorf("Expected 0 results for ISBN '000', got %v", results)
	}
}

func TestSearchByAuthorEmail(t *testing.T) {
	books := []Book{
		{Title: "Book1", Author: "author1@example.com", ISBN: "123", Genre: "Genre1"},
	}
	magazines := []Magazine{
		{Title: "Magazine1", Publisher: "publisher1@example.com", IssueNumber: "789", PublishedAt: "2025-01-01"},
	}

	results := searchByAuthorEmail(books, magazines, "author1@example.com")
	if len(results) != 1 || results[0] != "Book - Title: Book1\nAuthor: author1@example.com\nISBN: 123\nDescription: Genre1\n" {
		t.Errorf("Expected 1 result for email 'author1@example.com', got %v", results)
	}

	results = searchByAuthorEmail(books, magazines, "publisher1@example.com")
	if len(results) != 1 || results[0] != "Magazine - Title: Magazine1\nPublisher: publisher1@example.com\nISBN: 789\nPublished At: 2025-01-01\n" {
		t.Errorf("Expected 1 result for email 'publisher1@example.com', got %v", results)
	}

	results = searchByAuthorEmail(books, magazines, "unknown@example.com")
	if len(results) != 0 {
		t.Errorf("Expected 0 results for email 'unknown@example.com', got %v", results)
	}
}

func TestSortByTitle(t *testing.T) {
	books := []Book{
		{Title: "Book2", Author: "Author2", ISBN: "456", Genre: "Genre2"},
		{Title: "Book1", Author: "Author1", ISBN: "123", Genre: "Genre1"},
	}
	magazines := []Magazine{
		{Title: "Magazine1", Publisher: "Publisher1", IssueNumber: "789", PublishedAt: "2025-01-01"},
		{Title: "Magazine2", Publisher: "Publisher2", IssueNumber: "890", PublishedAt: "2025-02-01"},
	}

	results := sortByTitle(books, magazines)
	if len(results) != 4 || results[0] != "Book - Title: Book1\nAuthor: Author1\nISBN: 123\nDescription: Genre1\n" || results[1] != "Book - Title: Book2\nAuthor: Author2\nISBN: 456\nDescription: Genre2\n" || results[2] != "Magazine - Title: Magazine1\nPublisher: Publisher1\nISBN: 789\nPublished At: 2025-01-01\n" || results[3] != "Magazine - Title: Magazine2\nPublisher: Publisher2\nISBN: 890\nPublished At: 2025-02-01\n" {
		t.Errorf("Sorting by title failed, got %v", results)
	}
}
