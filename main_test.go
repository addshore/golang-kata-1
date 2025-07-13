package main

import (
	"reflect"
	"strings"
	"testing"
)

func TestParseBooks(t *testing.T) {
	csvData := `title;isbn;authors;description
Book 1;111-111;author1@a.com,author2@a.com;Desc 1
Book 2;222-222;author2@a.com;Desc 2`

	reader := strings.NewReader(csvData)
	books, err := parseBooks(reader, "test_books.csv")

	if err != nil {
		t.Fatalf("parseBooks() error = %v", err)
	}

	if len(books) != 2 {
		t.Fatalf("expected 2 books, got %d", len(books))
	}

	expectedBooks := []Book{
		{Title: "Book 1", ISBN: "111-111", Authors: []string{"author1@a.com", "author2@a.com"}, Description: "Desc 1"},
		{Title: "Book 2", ISBN: "222-222", Authors: []string{"author2@a.com"}, Description: "Desc 2"},
	}

	if !reflect.DeepEqual(books, expectedBooks) {
		t.Errorf("parsed books do not match expected books. got %+v, want %+v", books, expectedBooks)
	}
}

func TestParseMagazines(t *testing.T) {
	csvData := `title;isbn;authors;publishedAt
Magazine 1;333-333;author3@a.com;01.01.2022
Magazine 2;444-444;author1@a.com,author3@a.com;02.02.2022`

	reader := strings.NewReader(csvData)
	magazines, err := parseMagazines(reader, "test_magazines.csv")

	if err != nil {
		t.Fatalf("parseMagazines() error = %v", err)
	}

	if len(magazines) != 2 {
		t.Fatalf("expected 2 magazines, got %d", len(magazines))
	}

	expectedMagazines := []Magazine{
		{Title: "Magazine 1", ISBN: "333-333", Authors: []string{"author3@a.com"}, PublishedAt: "01.01.2022"},
		{Title: "Magazine 2", ISBN: "444-444", Authors: []string{"author1@a.com", "author3@a.com"}, PublishedAt: "02.02.2022"},
	}

	if !reflect.DeepEqual(magazines, expectedMagazines) {
		t.Errorf("parsed magazines do not match expected magazines. got %+v, want %+v", magazines, expectedMagazines)
	}
}

func TestParseAuthors(t *testing.T) {
	// The `\ufeff` is the UTF-8 BOM that might be at the start of the file.
	csvData := `﻿email;firstname;lastname
author1@a.com;First1;Last1
author2@a.com;First2;Last2`

	reader := strings.NewReader(csvData)
	authors, err := parseAuthors(reader, "test_authors.csv")

	if err != nil {
		t.Fatalf("parseAuthors() error = %v", err)
	}

	if len(authors) != 2 {
		t.Fatalf("expected 2 authors, got %d", len(authors))
	}

	expectedAuthors := map[string]Author{
		"author1@a.com": {Email: "author1@a.com", FirstName: "First1", LastName: "Last1"},
		"author2@a.com": {Email: "author2@a.com", FirstName: "First2", LastName: "Last2"},
	}

	if !reflect.DeepEqual(authors, expectedAuthors) {
		t.Errorf("parsed authors do not match expected authors. got %+v, want %+v", authors, expectedAuthors)
	}
}

var testBooks = []Book{
	{Title: "Book A", ISBN: "111", Authors: []string{"author1@a.com"}},
	{Title: "Book C", ISBN: "333", Authors: []string{"author2@a.com"}},
}
var testMagazines = []Magazine{
	{Title: "Magazine B", ISBN: "222", Authors: []string{"author1@a.com"}},
	{Title: "Magazine D", ISBN: "444", Authors: []string{"author3@a.com"}},
}

func TestFindItemByISBN(t *testing.T) {
	// Test finding a book
	item := findItemByISBN("111", testBooks, testMagazines)
	if item == nil || item.GetTitle() != "Book A" {
		t.Errorf("Expected to find 'Book A', got %v", item)
	}

	// Test finding a magazine
	item = findItemByISBN("222", testBooks, testMagazines)
	if item == nil || item.GetTitle() != "Magazine B" {
		t.Errorf("Expected to find 'Magazine B', got %v", item)
	}

	// Test not finding an item
	item = findItemByISBN("999", testBooks, testMagazines)
	if item != nil {
		t.Errorf("Expected to find nothing, got %v", item)
	}
}

func TestFindItemsByAuthor(t *testing.T) {
	// Test author with multiple items
	items := findItemsByAuthor("author1@a.com", testBooks, testMagazines)
	if len(items) != 2 {
		t.Fatalf("Expected 2 items for author1, got %d", len(items))
	}

	// Test author with one item
	items = findItemsByAuthor("author2@a.com", testBooks, testMagazines)
	if len(items) != 1 || items[0].GetTitle() != "Book C" {
		t.Errorf("Expected 1 item 'Book C' for author2, got %v", items)
	}

	// Test author with no items
	items = findItemsByAuthor("nonexistent@a.com", testBooks, testMagazines)
	if len(items) != 0 {
		t.Errorf("Expected 0 items for nonexistent author, got %d", len(items))
	}
}

func TestGetSortedItems(t *testing.T) {
	// Test ascending sort
	items, err := getSortedItems(testBooks, testMagazines, "title-asc")
	if err != nil {
		t.Fatalf("getSortedItems asc failed: %v", err)
	}
	expectedOrderAsc := []string{"Book A", "Book C", "Magazine B", "Magazine D"}
	for i, item := range items {
		if item.GetTitle() != expectedOrderAsc[i] {
			t.Errorf("Asc sort order incorrect. At index %d, got %s, want %s", i, item.GetTitle(), expectedOrderAsc[i])
		}
	}

	// Test descending sort
	items, err = getSortedItems(testBooks, testMagazines, "title-desc")
	if err != nil {
		t.Fatalf("getSortedItems desc failed: %v", err)
	}
	expectedOrderDesc := []string{"Magazine D", "Magazine B", "Book C", "Book A"}
	for i, item := range items {
		if item.GetTitle() != expectedOrderDesc[i] {
			t.Errorf("Desc sort order incorrect. At index %d, got %s, want %s", i, item.GetTitle(), expectedOrderDesc[i])
		}
	}

	// Test invalid sort order
	_, err = getSortedItems(testBooks, testMagazines, "invalid-sort")
	if err == nil {
		t.Error("Expected an error for invalid sort order, but got nil")
	}
}
