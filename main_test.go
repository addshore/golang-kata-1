package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestWelcomeMessage tests the welcome message function
func TestWelcomeMessage(t *testing.T) {
	msg := welcomeMessage()
	expected := "Hello world!"
	if msg != expected {
		t.Errorf("Expected '%s', got '%s'", expected, msg)
	}
}

// TestHandleBooks tests the /api/books endpoint
func TestHandleBooks(t *testing.T) {
	// Setup library with test data
	library = NewLibrary()
	author := Author{
		Email:     "test@example.com",
		FirstName: "Test",
		LastName:  "Author",
	}
	library.Authors[author.Email] = author
	library.Books = append(library.Books, Book{
		Title:       "Test Book",
		ISBN:        "123",
		Authors:     []Author{author},
		Description: "Test",
	})

	req := httptest.NewRequest("GET", "/api/books", nil)
	w := httptest.NewRecorder()
	handleBooks(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var books []Book
	err := json.NewDecoder(w.Body).Decode(&books)
	if err != nil {
		t.Errorf("Failed to decode response: %v", err)
	}

	if len(books) != 1 || books[0].Title != "Test Book" {
		t.Errorf("Expected 1 test book, got %v", books)
	}
}

// TestHandleMagazines tests the /api/magazines endpoint
func TestHandleMagazines(t *testing.T) {
	// Setup library with test data
	library = NewLibrary()
	author := Author{
		Email:     "test@example.com",
		FirstName: "Test",
		LastName:  "Author",
	}
	library.Authors[author.Email] = author
	library.Magazines = append(library.Magazines, Magazine{
		Title:       "Test Magazine",
		ISBN:        "456",
		Authors:     []Author{author},
		PublishedAt: "2025-01-01",
	})

	req := httptest.NewRequest("GET", "/api/magazines", nil)
	w := httptest.NewRecorder()
	handleMagazines(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var magazines []Magazine
	err := json.NewDecoder(w.Body).Decode(&magazines)
	if err != nil {
		t.Errorf("Failed to decode response: %v", err)
	}

	if len(magazines) != 1 || magazines[0].Title != "Test Magazine" {
		t.Errorf("Expected 1 test magazine, got %v", magazines)
	}
}

// TestHandleLibrary tests the /api/library endpoint
func TestHandleLibrary(t *testing.T) {
	// Setup library with test data
	library = NewLibrary()
	author := Author{
		Email:     "test@example.com",
		FirstName: "Test",
		LastName:  "Author",
	}
	library.Authors[author.Email] = author
	library.Books = append(library.Books, Book{
		Title:       "Test Book",
		ISBN:        "123",
		Authors:     []Author{author},
		Description: "Test",
	})
	library.Magazines = append(library.Magazines, Magazine{
		Title:       "Test Magazine",
		ISBN:        "456",
		Authors:     []Author{author},
		PublishedAt: "2025-01-01",
	})

	req := httptest.NewRequest("GET", "/api/library", nil)
	w := httptest.NewRecorder()
	handleLibrary(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var result map[string]interface{}
	err := json.NewDecoder(w.Body).Decode(&result)
	if err != nil {
		t.Errorf("Failed to decode response: %v", err)
	}

	books, ok := result["books"].([]interface{})
	if !ok || len(books) != 1 {
		t.Errorf("Expected 1 book in response")
	}

	magazines, ok := result["magazines"].([]interface{})
	if !ok || len(magazines) != 1 {
		t.Errorf("Expected 1 magazine in response")
	}
}

// TestHandleSearchBooks tests the /api/books/search endpoint
func TestHandleSearchBooks(t *testing.T) {
	// Setup library with test data
	library = NewLibrary()
	author := Author{
		Email:     "test@example.com",
		FirstName: "Test",
		LastName:  "Author",
	}
	library.Authors[author.Email] = author
	library.Books = append(library.Books,
		Book{
			Title:       "Book A",
			ISBN:        "978-1111111111",
			Authors:     []Author{author},
			Description: "First",
		},
		Book{
			Title:       "Book B",
			ISBN:        "978-2222222222",
			Authors:     []Author{author},
			Description: "Second",
		},
	)

	// Test with search parameter
	req := httptest.NewRequest("GET", "/api/books/search?isbn=978-1", nil)
	w := httptest.NewRecorder()
	handleSearchBooks(w, req)

	var books []Book
	err := json.NewDecoder(w.Body).Decode(&books)
	if err != nil {
		t.Errorf("Failed to decode response: %v", err)
	}

	if len(books) != 1 || books[0].Title != "Book A" {
		t.Errorf("Expected 1 book with ISBN '978-1', got %d books", len(books))
	}

	// Test without search parameter (should return all)
	req = httptest.NewRequest("GET", "/api/books/search", nil)
	w = httptest.NewRecorder()
	handleSearchBooks(w, req)

	err = json.NewDecoder(w.Body).Decode(&books)
	if err != nil {
		t.Errorf("Failed to decode response: %v", err)
	}

	if len(books) != 2 {
		t.Errorf("Expected 2 books without search parameter, got %d", len(books))
	}
}

// TestHandleSearchMagazines tests the /api/magazines/search endpoint
func TestHandleSearchMagazines(t *testing.T) {
	// Setup library with test data
	library = NewLibrary()
	author := Author{
		Email:     "test@example.com",
		FirstName: "Test",
		LastName:  "Author",
	}
	library.Authors[author.Email] = author
	library.Magazines = append(library.Magazines,
		Magazine{
			Title:       "Magazine A",
			ISBN:        "555-1111",
			Authors:     []Author{author},
			PublishedAt: "2025-01-01",
		},
		Magazine{
			Title:       "Magazine B",
			ISBN:        "666-2222",
			Authors:     []Author{author},
			PublishedAt: "2025-02-01",
		},
	)

	// Test with search parameter
	req := httptest.NewRequest("GET", "/api/magazines/search?isbn=555", nil)
	w := httptest.NewRecorder()
	handleSearchMagazines(w, req)

	var magazines []Magazine
	err := json.NewDecoder(w.Body).Decode(&magazines)
	if err != nil {
		t.Errorf("Failed to decode response: %v", err)
	}

	if len(magazines) != 1 || magazines[0].Title != "Magazine A" {
		t.Errorf("Expected 1 magazine with ISBN '555', got %d magazines", len(magazines))
	}
}

// TestHandleSearchBooksByAuthor tests the /api/books/search-by-author endpoint
func TestHandleSearchBooksByAuthor(t *testing.T) {
	// Setup library with test data
	library = NewLibrary()
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
	library.Authors[author1.Email] = author1
	library.Authors[author2.Email] = author2

	library.Books = append(library.Books,
		Book{
			Title:       "John's Book",
			ISBN:        "111",
			Authors:     []Author{author1},
			Description: "By John",
		},
		Book{
			Title:       "Jane's Book",
			ISBN:        "222",
			Authors:     []Author{author2},
			Description: "By Jane",
		},
	)

	// Test search by author
	req := httptest.NewRequest("GET", "/api/books/search-by-author?email=john@example.com", nil)
	w := httptest.NewRecorder()
	handleSearchBooksByAuthor(w, req)

	var books []Book
	err := json.NewDecoder(w.Body).Decode(&books)
	if err != nil {
		t.Errorf("Failed to decode response: %v", err)
	}

	if len(books) != 1 || books[0].Title != "John's Book" {
		t.Errorf("Expected 1 book by john@example.com, got %d", len(books))
	}
}

// TestHandleSearchMagazinesByAuthor tests the /api/magazines/search-by-author endpoint
func TestHandleSearchMagazinesByAuthor(t *testing.T) {
	// Setup library with test data
	library = NewLibrary()
	author1 := Author{
		Email:     "editor1@example.com",
		FirstName: "Editor",
		LastName:  "One",
	}
	library.Authors[author1.Email] = author1

	library.Magazines = append(library.Magazines,
		Magazine{
			Title:       "Magazine by Editor",
			ISBN:        "M-001",
			Authors:     []Author{author1},
			PublishedAt: "2025-01-01",
		},
	)

	// Test search by author
	req := httptest.NewRequest("GET", "/api/magazines/search-by-author?email=editor1@example.com", nil)
	w := httptest.NewRecorder()
	handleSearchMagazinesByAuthor(w, req)

	var magazines []Magazine
	err := json.NewDecoder(w.Body).Decode(&magazines)
	if err != nil {
		t.Errorf("Failed to decode response: %v", err)
	}

	if len(magazines) != 1 || magazines[0].Title != "Magazine by Editor" {
		t.Errorf("Expected 1 magazine by editor1@example.com, got %d", len(magazines))
	}
}

// TestHandleSortedItems tests the /api/items/sorted endpoint
func TestHandleSortedItems(t *testing.T) {
	// Setup library with test data
	library = NewLibrary()
	author := Author{
		Email:     "test@example.com",
		FirstName: "Test",
		LastName:  "Author",
	}
	library.Authors[author.Email] = author

	library.Books = append(library.Books,
		Book{
			Title:       "Zebra",
			ISBN:        "Z1",
			Authors:     []Author{author},
			Description: "",
		},
		Book{
			Title:       "Apple",
			ISBN:        "A1",
			Authors:     []Author{author},
			Description: "",
		},
	)

	library.Magazines = append(library.Magazines,
		Magazine{
			Title:       "Mango",
			ISBN:        "M1",
			Authors:     []Author{author},
			PublishedAt: "",
		},
	)

	// Test ascending order
	req := httptest.NewRequest("GET", "/api/items/sorted?direction=asc", nil)
	w := httptest.NewRecorder()
	handleSortedItems(w, req)

	var items []Item
	err := json.NewDecoder(w.Body).Decode(&items)
	if err != nil {
		t.Errorf("Failed to decode response: %v", err)
	}

	if len(items) != 3 {
		t.Errorf("Expected 3 items, got %d", len(items))
	}

	if items[0].Title != "Apple" || items[1].Title != "Mango" || items[2].Title != "Zebra" {
		t.Errorf("Expected ascending order: Apple, Mango, Zebra; got %s, %s, %s",
			items[0].Title, items[1].Title, items[2].Title)
	}

	// Test descending order
	req = httptest.NewRequest("GET", "/api/items/sorted?direction=desc", nil)
	w = httptest.NewRecorder()
	handleSortedItems(w, req)

	err = json.NewDecoder(w.Body).Decode(&items)
	if err != nil {
		t.Errorf("Failed to decode response: %v", err)
	}

	if items[0].Title != "Zebra" || items[1].Title != "Mango" || items[2].Title != "Apple" {
		t.Errorf("Expected descending order: Zebra, Mango, Apple; got %s, %s, %s",
			items[0].Title, items[1].Title, items[2].Title)
	}
}

// TestServeIndex tests that the index page is served as HTML
func TestServeIndex(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	serveIndex(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	contentType := w.Header().Get("Content-Type")
	if contentType != "text/html; charset=utf-8" {
		t.Errorf("Expected HTML content type, got %s", contentType)
	}

	body := w.Body.String()
	if len(body) == 0 {
		t.Error("Expected HTML content, got empty response")
	}

	if !contains(body, "<!DOCTYPE html>") {
		t.Error("Expected HTML document starting with <!DOCTYPE html>")
	}
}

// TestGetHTMLContent tests that HTML content is generated
func TestGetHTMLContent(t *testing.T) {
	html := getHTMLContent()

	if html == "" {
		t.Error("Expected HTML content, got empty string")
	}

	required := []string{"<!DOCTYPE html>", "<html", "<body>", "Library System"}
	for _, req := range required {
		if !contains(html, req) {
			t.Errorf("Expected HTML to contain '%s'", req)
		}
	}
}

// Helper function to check if a string contains a substring
func contains(s, substring string) bool {
	for i := 0; i <= len(s)-len(substring); i++ {
		if s[i:i+len(substring)] == substring {
			return true
		}
	}
	return false
}

// TestContentTypeHeaders tests that endpoints return correct content types
func TestContentTypeHeaders(t *testing.T) {
	testCases := []struct {
		name     string
		endpoint string
		handler  func(http.ResponseWriter, *http.Request)
		expected string
	}{
		{
			name:     "Books endpoint",
			endpoint: "/api/books",
			handler:  handleBooks,
			expected: "application/json",
		},
		{
			name:     "Magazines endpoint",
			endpoint: "/api/magazines",
			handler:  handleMagazines,
			expected: "application/json",
		},
		{
			name:     "Library endpoint",
			endpoint: "/api/library",
			handler:  handleLibrary,
			expected: "application/json",
		},
		{
			name:     "Search books endpoint",
			endpoint: "/api/books/search",
			handler:  handleSearchBooks,
			expected: "application/json",
		},
	}

	library = NewLibrary()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", tc.endpoint, nil)
			w := httptest.NewRecorder()
			tc.handler(w, req)

			contentType := w.Header().Get("Content-Type")
			if contentType != tc.expected {
				t.Errorf("Expected content type %s, got %s", tc.expected, contentType)
			}
		})
	}
}
