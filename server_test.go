package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleIndex(t *testing.T) {
	// Setup dummy data
	author := &Author{Email: "test@example.com", Firstname: "Test", Lastname: "User"}
	authors := map[string]*Author{"test@example.com": author}

	books := []*Book{
		{Title: "Book A", ISBN: "111", Authors: []*Author{author}, Description: "Desc A"},
	}

	mags := []*Magazine{
		{Title: "Mag A", ISBN: "333", Authors: []*Author{author}, PublishedAt: "2021-01-01"},
	}

	data := &LibraryData{
		Authors:   authors,
		Books:     books,
		Magazines: mags,
	}

	server := &Server{Data: data}

	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.handleIndex)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	// We can check if it contains the title of our book
	expected := "Book A"
	if !contains(rr.Body.String(), expected) {
		t.Errorf("handler returned unexpected body: does not contain %v",
			expected)
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && len(substr) > 0 && s[0:len(substr)] == substr ||
		len(s) > len(substr) && len(substr) > 0 && contains(s[1:], substr)
}
