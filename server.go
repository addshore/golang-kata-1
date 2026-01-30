package main

import (
	"html/template"
	"net/http"
	"path/filepath"
)

type Server struct {
	Data *LibraryData
}

func (s *Server) Start(addr string) error {
	http.HandleFunc("/", s.handleIndex)
	http.HandleFunc("/add", s.handleAdd)
	return http.ListenAndServe(addr, nil)
}

func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	// Parse the template on every request for development convenience (hot-reloading)
	tmplPath := filepath.Join("templates", "index.html")
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		http.Error(w, "Could not load template: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Filter and Sort data using the service layer
	isbn := r.URL.Query().Get("isbn")
	email := r.URL.Query().Get("email")
	sortParam := r.URL.Query().Get("sort")
	dirParam := r.URL.Query().Get("dir")

	processedData := ProcessRequest(s.Data, isbn, email, sortParam, dirParam)

	if err := tmpl.Execute(w, processedData); err != nil {
		http.Error(w, "Failed to render template: "+err.Error(), http.StatusInternalServerError)
	}
}

func (s *Server) handleAdd(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form: "+err.Error(), http.StatusBadRequest)
		return
	}

	itemType := r.FormValue("type")
	title := r.FormValue("title")
	isbn := r.FormValue("isbn")
	extra := r.FormValue("extra") // Description or PublishedAt
	email := r.FormValue("email")
	firstname := r.FormValue("firstname")
	lastname := r.FormValue("lastname")

	if title == "" || isbn == "" || email == "" || firstname == "" || lastname == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	// 1. Handle Author
	author, exists := s.Data.Authors[email]
	if !exists {
		author = &Author{
			Email:     email,
			Firstname: firstname,
			Lastname:  lastname,
		}
		// Persist new author
		if err := AppendAuthor(filepath.Join("resources", "authors.csv"), author); err != nil {
			http.Error(w, "Failed to save author: "+err.Error(), http.StatusInternalServerError)
			return
		}
		// Update memory
		s.Data.Authors[email] = author
	}

	// 2. Handle Item
	authors := []*Author{author}

	if itemType == "Book" {
		book := &Book{
			Title:       title,
			ISBN:        isbn,
			Authors:     authors,
			Description: extra,
		}
		// Persist book
		if err := AppendBook(filepath.Join("resources", "books.csv"), book); err != nil {
			http.Error(w, "Failed to save book: "+err.Error(), http.StatusInternalServerError)
			return
		}
		// Update memory
		s.Data.Books = append(s.Data.Books, book)
	} else if itemType == "Magazine" {
		mag := &Magazine{
			Title:       title,
			ISBN:        isbn,
			Authors:     authors,
			PublishedAt: extra,
		}
		// Persist magazine
		if err := AppendMagazine(filepath.Join("resources", "magazines.csv"), mag); err != nil {
			http.Error(w, "Failed to save magazine: "+err.Error(), http.StatusInternalServerError)
			return
		}
		// Update memory
		s.Data.Magazines = append(s.Data.Magazines, mag)
	} else {
		http.Error(w, "Invalid item type", http.StatusBadRequest)
		return
	}

	// Redirect back to home
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
