package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"
)

// WebServer provides HTTP handlers for the web interface
type WebServer struct {
	library      *Library
	resourcesDir string
	templates    *template.Template
}

// NewWebServer creates a new web server instance
func NewWebServer(library *Library, resourcesDir string) (*WebServer, error) {
	// Parse templates
	tmpl, err := template.New("").Funcs(template.FuncMap{
		"join": strings.Join,
	}).ParseGlob("templates/*.html")
	if err != nil {
		return nil, err
	}

	return &WebServer{
		library:      library,
		resourcesDir: resourcesDir,
		templates:    tmpl,
	}, nil
}

// SetupRoutes configures HTTP routes
func (ws *WebServer) SetupRoutes() {
	http.HandleFunc("/", ws.handleHome)
	http.HandleFunc("/search", ws.handleSearch)
	http.HandleFunc("/sorted", ws.handleSorted)
	http.HandleFunc("/add", ws.handleAdd)
	http.HandleFunc("/add-submit", ws.handleAddSubmit)
}

func (ws *WebServer) handleHome(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Books     []*Book
		Magazines []*Magazine
		Library   *Library
	}{
		Books:     ws.library.Books,
		Magazines: ws.library.Magazines,
		Library:   ws.library,
	}

	if err := ws.templates.ExecuteTemplate(w, "home.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (ws *WebServer) handleSearch(w http.ResponseWriter, r *http.Request) {
	searchType := r.URL.Query().Get("type")
	query := r.URL.Query().Get("q")

	var books []*Book
	var magazines []*Magazine
	var searchLabel string

	if query != "" {
		if searchType == "isbn" {
			books, magazines = ws.library.SearchByISBN(query)
			searchLabel = fmt.Sprintf("ISBN: %s", query)
		} else if searchType == "author" {
			books, magazines = ws.library.SearchByAuthorEmail(query)
			if author, ok := ws.library.Authors[query]; ok {
				searchLabel = fmt.Sprintf("Author: %s %s (%s)", author.FirstName, author.LastName, query)
			} else {
				searchLabel = fmt.Sprintf("Author Email: %s", query)
			}
		}
	}

	data := struct {
		Books       []*Book
		Magazines   []*Magazine
		Library     *Library
		SearchLabel string
		Query       string
	}{
		Books:       books,
		Magazines:   magazines,
		Library:     ws.library,
		SearchLabel: searchLabel,
		Query:       query,
	}

	if err := ws.templates.ExecuteTemplate(w, "search.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (ws *WebServer) handleSorted(w http.ResponseWriter, r *http.Request) {
	items := ws.library.GetAllSortedByTitle()

	data := struct {
		Items   []Item
		Library *Library
	}{
		Items:   items,
		Library: ws.library,
	}

	if err := ws.templates.ExecuteTemplate(w, "sorted.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (ws *WebServer) handleAdd(w http.ResponseWriter, r *http.Request) {
	if err := ws.templates.ExecuteTemplate(w, "add.html", nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (ws *WebServer) handleAddSubmit(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	itemType := r.FormValue("type")
	title := strings.TrimSpace(r.FormValue("title"))
	isbn := strings.TrimSpace(r.FormValue("isbn"))

	// Validate required fields
	if itemType == "" {
		http.Error(w, "Item type is required", http.StatusBadRequest)
		return
	}
	if title == "" {
		http.Error(w, "Title is required", http.StatusBadRequest)
		return
	}
	if isbn == "" {
		http.Error(w, "ISBN is required", http.StatusBadRequest)
		return
	}

	// Parse authors
	authorsInput := r.FormValue("authors")
	authorEmails := []string{}
	if authorsInput != "" {
		for _, line := range strings.Split(authorsInput, "\n") {
			line = strings.TrimSpace(line)
			if line != "" {
				parts := strings.Split(line, ",")
				if len(parts) >= 3 {
					email := strings.TrimSpace(parts[0])
					firstName := strings.TrimSpace(parts[1])
					lastName := strings.TrimSpace(parts[2])
					ws.library.AddAuthor(email, firstName, lastName)
					authorEmails = append(authorEmails, email)
				}
			}
		}
	}

	if itemType == "book" {
		description := strings.TrimSpace(r.FormValue("description"))
		ws.library.AddBook(title, isbn, description, authorEmails)
	} else if itemType == "magazine" {
		publishedAt := strings.TrimSpace(r.FormValue("publishedAt"))
		ws.library.AddMagazine(title, isbn, publishedAt, authorEmails)
	} else {
		http.Error(w, "Invalid item type", http.StatusBadRequest)
		return
	}

	// Save to CSV
	if err := ws.library.SaveToCSV(ws.resourcesDir); err != nil {
		http.Error(w, "Failed to save changes to the library", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Start starts the web server
func (ws *WebServer) Start(addr string) error {
	fmt.Printf("Starting web server on %s\n", addr)
	return http.ListenAndServe(addr, nil)
}
