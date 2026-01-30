package main

import (
	"testing"
)

// TestNewLibrary tests that a new library is properly initialized
func TestNewLibrary(t *testing.T) {
	lib := NewLibrary()

	if lib.Books == nil || len(lib.Books) != 0 {
		t.Error("Books should be initialized as an empty slice")
	}

	if lib.Magazines == nil || len(lib.Magazines) != 0 {
		t.Error("Magazines should be initialized as an empty slice")
	}

	if lib.Authors == nil {
		t.Error("Authors should be initialized as an empty map")
	}
}

// TestAddBooksAndMagazines tests adding books and magazines to the library
func TestAddBooksAndMagazines(t *testing.T) {
	lib := NewLibrary()

	// Add authors
	author1 := Author{
		Email:     "john@example.com",
		FirstName: "John",
		LastName:  "Doe",
	}
	author2 := Author{
		Email:     "jane@example.com",
		FirstName: "Jane",
		LastName:  "Smith",
	}
	lib.Authors[author1.Email] = author1
	lib.Authors[author2.Email] = author2

	// Add books
	book1 := Book{
		Title:       "Test Book 1",
		ISBN:        "978-1234567890",
		Authors:     []Author{author1},
		Description: "A test book",
	}
	book2 := Book{
		Title:       "Test Book 2",
		ISBN:        "978-0987654321",
		Authors:     []Author{author2},
		Description: "Another test book",
	}
	lib.Books = append(lib.Books, book1, book2)

	// Add magazines
	mag1 := Magazine{
		Title:       "Test Magazine 1",
		ISBN:        "123-456",
		Authors:     []Author{author1},
		PublishedAt: "2025-01-01",
	}
	lib.Magazines = append(lib.Magazines, mag1)

	if len(lib.Books) != 2 {
		t.Errorf("Expected 2 books, got %d", len(lib.Books))
	}

	if len(lib.Magazines) != 1 {
		t.Errorf("Expected 1 magazine, got %d", len(lib.Magazines))
	}
}

// TestSearchBooksByISBN tests searching for books by ISBN
func TestSearchBooksByISBN(t *testing.T) {
	lib := NewLibrary()

	author := Author{
		Email:     "author@example.com",
		FirstName: "Test",
		LastName:  "Author",
	}
	lib.Authors[author.Email] = author

	book1 := Book{
		Title:       "Book A",
		ISBN:        "978-1111111111",
		Authors:     []Author{author},
		Description: "First book",
	}
	book2 := Book{
		Title:       "Book B",
		ISBN:        "978-2222222222",
		Authors:     []Author{author},
		Description: "Second book",
	}
	book3 := Book{
		Title:       "Book C",
		ISBN:        "111-2222222222",
		Authors:     []Author{author},
		Description: "Third book",
	}
	lib.Books = append(lib.Books, book1, book2, book3)

	// Test exact ISBN search
	results := lib.SearchBooksByISBN("978-1111111111")
	if len(results) != 1 || results[0].Title != "Book A" {
		t.Errorf("Expected to find Book A, got %v", results)
	}

	// Test partial ISBN search
	results = lib.SearchBooksByISBN("978-")
	if len(results) != 2 {
		t.Errorf("Expected to find 2 books with partial ISBN '978-', got %d", len(results))
	}

	// Test case-insensitive search
	results = lib.SearchBooksByISBN("111")
	if len(results) != 2 {
		t.Errorf("Expected to find 2 books with ISBN containing '111', got %d", len(results))
	}

	// Test search with no results
	results = lib.SearchBooksByISBN("999-9999999999")
	if len(results) != 0 {
		t.Errorf("Expected no results for non-existent ISBN, got %d", len(results))
	}
}

// TestSearchMagazinesByISBN tests searching for magazines by ISBN
func TestSearchMagazinesByISBN(t *testing.T) {
	lib := NewLibrary()

	author := Author{
		Email:     "author@example.com",
		FirstName: "Test",
		LastName:  "Author",
	}
	lib.Authors[author.Email] = author

	mag1 := Magazine{
		Title:       "Magazine A",
		ISBN:        "555-1111",
		Authors:     []Author{author},
		PublishedAt: "2025-01-01",
	}
	mag2 := Magazine{
		Title:       "Magazine B",
		ISBN:        "666-2222",
		Authors:     []Author{author},
		PublishedAt: "2025-02-01",
	}
	lib.Magazines = append(lib.Magazines, mag1, mag2)

	results := lib.SearchMagazinesByISBN("555")
	if len(results) != 1 || results[0].Title != "Magazine A" {
		t.Errorf("Expected to find Magazine A, got %v", results)
	}

	results = lib.SearchMagazinesByISBN("non-existent")
	if len(results) != 0 {
		t.Errorf("Expected no results, got %d", len(results))
	}
}

// TestSearchBooksByAuthorEmail tests searching for books by author email
func TestSearchBooksByAuthorEmail(t *testing.T) {
	lib := NewLibrary()

	author1 := Author{
		Email:     "john@example.com",
		FirstName: "John",
		LastName:  "Doe",
	}
	author2 := Author{
		Email:     "jane@example.com",
		FirstName: "Jane",
		LastName:  "Smith",
	}
	lib.Authors[author1.Email] = author1
	lib.Authors[author2.Email] = author2

	book1 := Book{
		Title:       "John's Book",
		ISBN:        "111",
		Authors:     []Author{author1},
		Description: "Book by John",
	}
	book2 := Book{
		Title:       "Collaboration",
		ISBN:        "222",
		Authors:     []Author{author1, author2},
		Description: "Book by both",
	}
	book3 := Book{
		Title:       "Jane's Book",
		ISBN:        "333",
		Authors:     []Author{author2},
		Description: "Book by Jane",
	}
	lib.Books = append(lib.Books, book1, book2, book3)

	// Search for John's books
	results := lib.SearchBooksByAuthorEmail("john@example.com")
	if len(results) != 2 {
		t.Errorf("Expected 2 books by john@example.com, got %d", len(results))
	}

	// Search with case variations
	results = lib.SearchBooksByAuthorEmail("JOHN@EXAMPLE.COM")
	if len(results) != 2 {
		t.Errorf("Expected case-insensitive search to find 2 books, got %d", len(results))
	}

	// Search with whitespace
	results = lib.SearchBooksByAuthorEmail("  jane@example.com  ")
	if len(results) != 2 {
		t.Errorf("Expected whitespace to be trimmed, got %d results", len(results))
	}

	// Search for non-existent author
	results = lib.SearchBooksByAuthorEmail("nonexistent@example.com")
	if len(results) != 0 {
		t.Errorf("Expected no results for non-existent author, got %d", len(results))
	}
}

// TestSearchMagazinesByAuthorEmail tests searching for magazines by author email
func TestSearchMagazinesByAuthorEmail(t *testing.T) {
	lib := NewLibrary()

	author1 := Author{
		Email:     "editor1@example.com",
		FirstName: "Editor",
		LastName:  "One",
	}
	author2 := Author{
		Email:     "editor2@example.com",
		FirstName: "Editor",
		LastName:  "Two",
	}
	lib.Authors[author1.Email] = author1
	lib.Authors[author2.Email] = author2

	mag1 := Magazine{
		Title:       "Magazine 1",
		ISBN:        "M-001",
		Authors:     []Author{author1},
		PublishedAt: "2025-01-01",
	}
	mag2 := Magazine{
		Title:       "Magazine 2",
		ISBN:        "M-002",
		Authors:     []Author{author1, author2},
		PublishedAt: "2025-02-01",
	}
	lib.Magazines = append(lib.Magazines, mag1, mag2)

	results := lib.SearchMagazinesByAuthorEmail("editor1@example.com")
	if len(results) != 2 {
		t.Errorf("Expected 2 magazines, got %d", len(results))
	}

	results = lib.SearchMagazinesByAuthorEmail("editor2@example.com")
	if len(results) != 1 {
		t.Errorf("Expected 1 magazine, got %d", len(results))
	}
}

// TestGetAllItemsSortedByTitle tests sorting items by title in ascending order
func TestGetAllItemsSortedByTitle(t *testing.T) {
	lib := NewLibrary()

	author := Author{
		Email:     "author@example.com",
		FirstName: "Test",
		LastName:  "Author",
	}
	lib.Authors[author.Email] = author

	// Add books in non-alphabetical order
	lib.Books = append(lib.Books,
		Book{Title: "Zebra", ISBN: "Z1", Authors: []Author{author}, Description: ""},
		Book{Title: "Apple", ISBN: "A1", Authors: []Author{author}, Description: ""},
		Book{Title: "Mango", ISBN: "M1", Authors: []Author{author}, Description: ""},
	)

	// Add magazines in non-alphabetical order
	lib.Magazines = append(lib.Magazines,
		Magazine{Title: "Xray", ISBN: "X1", Authors: []Author{author}, PublishedAt: ""},
		Magazine{Title: "Banana", ISBN: "B1", Authors: []Author{author}, PublishedAt: ""},
	)

	items := lib.GetAllItemsSortedByTitle()

	expectedOrder := []string{"Apple", "Banana", "Mango", "Xray", "Zebra"}
	if len(items) != 5 {
		t.Errorf("Expected 5 items, got %d", len(items))
	}

	for i, expected := range expectedOrder {
		if items[i].Title != expected {
			t.Errorf("Item %d: expected '%s', got '%s'", i, expected, items[i].Title)
		}
	}
}

// TestGetAllItemsSortedByTitleDescending tests sorting items by title in descending order
func TestGetAllItemsSortedByTitleDescending(t *testing.T) {
	lib := NewLibrary()

	author := Author{
		Email:     "author@example.com",
		FirstName: "Test",
		LastName:  "Author",
	}
	lib.Authors[author.Email] = author

	// Add books in non-alphabetical order
	lib.Books = append(lib.Books,
		Book{Title: "Zebra", ISBN: "Z1", Authors: []Author{author}, Description: ""},
		Book{Title: "Apple", ISBN: "A1", Authors: []Author{author}, Description: ""},
		Book{Title: "Mango", ISBN: "M1", Authors: []Author{author}, Description: ""},
	)

	// Add magazines
	lib.Magazines = append(lib.Magazines,
		Magazine{Title: "Xray", ISBN: "X1", Authors: []Author{author}, PublishedAt: ""},
		Magazine{Title: "Banana", ISBN: "B1", Authors: []Author{author}, PublishedAt: ""},
	)

	items := lib.GetAllItemsSortedByTitleDescending()

	expectedOrder := []string{"Zebra", "Xray", "Mango", "Banana", "Apple"}
	if len(items) != 5 {
		t.Errorf("Expected 5 items, got %d", len(items))
	}

	for i, expected := range expectedOrder {
		if items[i].Title != expected {
			t.Errorf("Item %d: expected '%s', got '%s'", i, expected, items[i].Title)
		}
	}
}

// TestItemTypeAssignment tests that items are correctly identified as books or magazines
func TestItemTypeAssignment(t *testing.T) {
	lib := NewLibrary()

	author := Author{
		Email:     "author@example.com",
		FirstName: "Test",
		LastName:  "Author",
	}
	lib.Authors[author.Email] = author

	lib.Books = append(lib.Books,
		Book{Title: "A Book", ISBN: "1", Authors: []Author{author}, Description: ""},
	)
	lib.Magazines = append(lib.Magazines,
		Magazine{Title: "A Magazine", ISBN: "2", Authors: []Author{author}, PublishedAt: ""},
	)

	items := lib.GetAllItemsSortedByTitle()

	if len(items) != 2 {
		t.Errorf("Expected 2 items, got %d", len(items))
	}

	if items[0].Type != "book" {
		t.Errorf("Expected first item to be book, got %s", items[0].Type)
	}

	if items[1].Type != "magazine" {
		t.Errorf("Expected second item to be magazine, got %s", items[1].Type)
	}
}

// TestCaseSensitiveSorting tests that sorting is case-insensitive
func TestCaseSensitiveSorting(t *testing.T) {
	lib := NewLibrary()

	author := Author{
		Email:     "author@example.com",
		FirstName: "Test",
		LastName:  "Author",
	}
	lib.Authors[author.Email] = author

	lib.Books = append(lib.Books,
		Book{Title: "zebra", ISBN: "1", Authors: []Author{author}, Description: ""},
		Book{Title: "APPLE", ISBN: "2", Authors: []Author{author}, Description: ""},
		Book{Title: "Mango", ISBN: "3", Authors: []Author{author}, Description: ""},
	)

	items := lib.GetAllItemsSortedByTitle()

	expectedOrder := []string{"APPLE", "Mango", "zebra"}
	for i, expected := range expectedOrder {
		if items[i].Title != expected {
			t.Errorf("Item %d: expected '%s', got '%s'", i, expected, items[i].Title)
		}
	}
}

// TestEmptyLibrarySorting tests sorting an empty library
func TestEmptyLibrarySorting(t *testing.T) {
	lib := NewLibrary()

	items := lib.GetAllItemsSortedByTitle()
	if len(items) != 0 {
		t.Errorf("Expected 0 items from empty library, got %d", len(items))
	}

	itemsDesc := lib.GetAllItemsSortedByTitleDescending()
	if len(itemsDesc) != 0 {
		t.Errorf("Expected 0 items from empty library (descending), got %d", len(itemsDesc))
	}
}

// TestSearchWithEmptyLibrary tests searching in an empty library
func TestSearchWithEmptyLibrary(t *testing.T) {
	lib := NewLibrary()

	results := lib.SearchBooksByISBN("any")
	if len(results) != 0 {
		t.Errorf("Expected no books in empty library, got %d", len(results))
	}

	results2 := lib.SearchMagazinesByISBN("any")
	if len(results2) != 0 {
		t.Errorf("Expected no magazines in empty library, got %d", len(results2))
	}

	resultsBooks := lib.SearchBooksByAuthorEmail("any@example.com")
	if len(resultsBooks) != 0 {
		t.Errorf("Expected no books in empty library, got %d", len(resultsBooks))
	}

	resultsMags := lib.SearchMagazinesByAuthorEmail("any@example.com")
	if len(resultsMags) != 0 {
		t.Errorf("Expected no magazines in empty library, got %d", len(resultsMags))
	}
}
