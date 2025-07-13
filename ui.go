package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"
	"time"
)

// UIHandler handles the web interface
type UIHandler struct {
	service *LibraryService
}

// NewUIHandler creates a new UI handler
func NewUIHandler(library *Library) *UIHandler {
	return &UIHandler{service: NewLibraryService(library)}
}

// StartServer starts the web server
func (h *UIHandler) StartServer(port string) error {
	http.HandleFunc("/", h.handleHome)
	http.HandleFunc("/books", h.handleBooks)
	http.HandleFunc("/magazines", h.handleMagazines)
	http.HandleFunc("/search", h.handleSearch)
	http.HandleFunc("/all-sorted", h.handleAllSorted)
	http.HandleFunc("/add", h.handleAdd)

	fmt.Printf("Starting server on http://localhost:%s\n", port)
	fmt.Println("Available endpoints:")
	fmt.Println("  / - Home page with overview")
	fmt.Println("  /books - List all books")
	fmt.Println("  /magazines - List all magazines")
	fmt.Println("  /search - Search by ISBN or Author Email")
	fmt.Println("  /all-sorted - All items sorted by title")
	fmt.Println("  /add - Add new books, magazines, and authors")

	return http.ListenAndServe(":"+port, nil)
}

// handleHome displays the home page with overview
func (h *UIHandler) handleHome(w http.ResponseWriter, r *http.Request) {
	tmpl := `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Library Management System</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 40px; background-color: #f5f5f5; }
        .container { max-width: 1200px; margin: 0 auto; background-color: white; padding: 30px; border-radius: 8px; box-shadow: 0 2px 10px rgba(0,0,0,0.1); }
        h1 { color: #333; text-align: center; margin-bottom: 30px; }
        .stats { display: flex; justify-content: space-around; margin: 30px 0; }
        .stat-card { background-color: #e8f4fd; padding: 20px; border-radius: 8px; text-align: center; min-width: 150px; }
        .stat-number { font-size: 2em; font-weight: bold; color: #2c5aa0; }
        .stat-label { color: #666; margin-top: 5px; }
        .navigation { text-align: center; margin: 30px 0; }
        .nav-button { display: inline-block; margin: 0 15px; padding: 12px 24px; background-color: #2c5aa0; color: white; text-decoration: none; border-radius: 5px; font-weight: bold; }
        .nav-button:hover { background-color: #1e3d6f; }
        .recent-items { margin-top: 40px; }
        .item-grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(300px, 1fr)); gap: 20px; margin-top: 20px; }
        .item-card { border: 1px solid #ddd; padding: 15px; border-radius: 5px; background-color: #fafafa; }
        .item-title { font-weight: bold; color: #333; margin-bottom: 8px; }
        .item-authors { color: #666; font-size: 0.9em; margin-bottom: 8px; }
        .item-isbn { color: #888; font-size: 0.8em; }
    </style>
</head>
<body>
    <div class="container">
        <h1>Library Management System</h1>
        
        <div class="stats">
            <div class="stat-card">
                <div class="stat-number">{{.BookCount}}</div>
                <div class="stat-label">Books</div>
            </div>
            <div class="stat-card">
                <div class="stat-number">{{.MagazineCount}}</div>
                <div class="stat-label">Magazines</div>
            </div>
            <div class="stat-card">
                <div class="stat-number">{{.AuthorCount}}</div>
                <div class="stat-label">Authors</div>
            </div>
        </div>

        <div class="navigation">
            <a href="/books" class="nav-button">View All Books</a>
            <a href="/magazines" class="nav-button">View All Magazines</a>
            <a href="/all-sorted" class="nav-button">All Items Sorted</a>
            <a href="/search" class="nav-button">Search Library</a>
            <a href="/add" class="nav-button" style="background-color: #28a745;">Add New Items</a>
        </div>

        <div class="recent-items">
            <h2>Sample Books</h2>
            <div class="item-grid">
                {{range .SampleBooks}}
                <div class="item-card">
                    <div class="item-title">{{.Title}}</div>
                    <div class="item-authors">Authors: {{.AuthorNames}}</div>
                    <div class="item-isbn">ISBN: {{.ISBN}}</div>
                </div>
                {{end}}
            </div>

            <h2>Sample Magazines</h2>
            <div class="item-grid">
                {{range .SampleMagazines}}
                <div class="item-card">
                    <div class="item-title">{{.Title}}</div>
                    <div class="item-authors">Authors: {{.AuthorNames}}</div>
                    <div class="item-isbn">ISBN: {{.ISBN}}</div>
                </div>
                {{end}}
            </div>
        </div>
    </div>
</body>
</html>`

	// Get library statistics
	bookCount, magazineCount, authorCount := h.service.GetLibraryStats()

	// Prepare data for template
	data := struct {
		BookCount       int
		MagazineCount   int
		AuthorCount     int
		SampleBooks     []BookDisplay
		SampleMagazines []MagazineDisplay
	}{
		BookCount:     bookCount,
		MagazineCount: magazineCount,
		AuthorCount:   authorCount,
	}

	// Get sample books (first 3)
	books := h.service.GetBooks()
	for i, book := range books {
		if i >= 3 {
			break
		}
		data.SampleBooks = append(data.SampleBooks, BookDisplay{
			Title:       book.Title,
			ISBN:        book.ISBN,
			AuthorNames: h.service.FormatAuthors(book.Authors),
			Description: book.Description,
		})
	}

	// Get sample magazines (first 3)
	magazines := h.service.GetMagazines()
	for i, magazine := range magazines {
		if i >= 3 {
			break
		}
		data.SampleMagazines = append(data.SampleMagazines, MagazineDisplay{
			Title:       magazine.Title,
			ISBN:        magazine.ISBN,
			AuthorNames: h.service.FormatAuthors(magazine.Authors),
			PublishedAt: magazine.PublishedAt.Format("02.01.2006"),
		})
	}

	t, err := template.New("home").Parse(tmpl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	t.Execute(w, data)
}

// BookDisplay represents a book for display purposes
type BookDisplay struct {
	Title       string
	ISBN        string
	AuthorNames string
	Description string
}

// MagazineDisplay represents a magazine for display purposes
type MagazineDisplay struct {
	Title       string
	ISBN        string
	AuthorNames string
	PublishedAt string
}

// handleBooks displays all books
func (h *UIHandler) handleBooks(w http.ResponseWriter, r *http.Request) {
	tmpl := `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>All Books - Library Management System</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 40px; background-color: #f5f5f5; }
        .container { max-width: 1200px; margin: 0 auto; background-color: white; padding: 30px; border-radius: 8px; box-shadow: 0 2px 10px rgba(0,0,0,0.1); }
        h1 { color: #333; text-align: center; margin-bottom: 30px; }
        .back-link { display: inline-block; margin-bottom: 20px; color: #2c5aa0; text-decoration: none; }
        .back-link:hover { text-decoration: underline; }
        .book-grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(400px, 1fr)); gap: 20px; }
        .book-card { border: 1px solid #ddd; padding: 20px; border-radius: 8px; background-color: #fafafa; }
        .book-title { font-size: 1.2em; font-weight: bold; color: #333; margin-bottom: 10px; }
        .book-authors { color: #666; margin-bottom: 8px; }
        .book-isbn { color: #888; font-size: 0.9em; margin-bottom: 12px; }
        .book-description { color: #555; line-height: 1.4; }
        .count { text-align: center; margin-bottom: 20px; color: #666; }
    </style>
</head>
<body>
    <div class="container">
        <a href="/" class="back-link">← Back to Home</a>
        <h1>All Books</h1>
        <div class="count">Total: {{.Count}} books</div>
        
        <div class="book-grid">
            {{range .Books}}
            <div class="book-card">
                <div class="book-title">{{.Title}}</div>
                <div class="book-authors"><strong>Authors:</strong> {{.AuthorNames}}</div>
                <div class="book-isbn"><strong>ISBN:</strong> {{.ISBN}}</div>
                <div class="book-description">{{.Description}}</div>
            </div>
            {{end}}
        </div>
    </div>
</body>
</html>`

	// Prepare data for template
	var books []BookDisplay
	for _, book := range h.service.GetBooks() {
		books = append(books, BookDisplay{
			Title:       book.Title,
			ISBN:        book.ISBN,
			AuthorNames: h.service.FormatAuthors(book.Authors),
			Description: book.Description,
		})
	}

	data := struct {
		Books []BookDisplay
		Count int
	}{
		Books: books,
		Count: len(books),
	}

	t, err := template.New("books").Parse(tmpl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	t.Execute(w, data)
}

// handleMagazines displays all magazines
func (h *UIHandler) handleMagazines(w http.ResponseWriter, r *http.Request) {
	tmpl := `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>All Magazines - Library Management System</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 40px; background-color: #f5f5f5; }
        .container { max-width: 1200px; margin: 0 auto; background-color: white; padding: 30px; border-radius: 8px; box-shadow: 0 2px 10px rgba(0,0,0,0.1); }
        h1 { color: #333; text-align: center; margin-bottom: 30px; }
        .back-link { display: inline-block; margin-bottom: 20px; color: #2c5aa0; text-decoration: none; }
        .back-link:hover { text-decoration: underline; }
        .magazine-grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(350px, 1fr)); gap: 20px; }
        .magazine-card { border: 1px solid #ddd; padding: 20px; border-radius: 8px; background-color: #fafafa; }
        .magazine-title { font-size: 1.2em; font-weight: bold; color: #333; margin-bottom: 10px; }
        .magazine-authors { color: #666; margin-bottom: 8px; }
        .magazine-isbn { color: #888; font-size: 0.9em; margin-bottom: 8px; }
        .magazine-published { color: #2c5aa0; font-weight: bold; }
        .count { text-align: center; margin-bottom: 20px; color: #666; }
    </style>
</head>
<body>
    <div class="container">
        <a href="/" class="back-link">← Back to Home</a>
        <h1>All Magazines</h1>
        <div class="count">Total: {{.Count}} magazines</div>
        
        <div class="magazine-grid">
            {{range .Magazines}}
            <div class="magazine-card">
                <div class="magazine-title">{{.Title}}</div>
                <div class="magazine-authors"><strong>Authors:</strong> {{.AuthorNames}}</div>
                <div class="magazine-isbn"><strong>ISBN:</strong> {{.ISBN}}</div>
                <div class="magazine-published"><strong>Published:</strong> {{.PublishedAt}}</div>
            </div>
            {{end}}
        </div>
    </div>
</body>
</html>`

	// Prepare data for template
	var magazines []MagazineDisplay
	for _, magazine := range h.service.GetMagazines() {
		magazines = append(magazines, MagazineDisplay{
			Title:       magazine.Title,
			ISBN:        magazine.ISBN,
			AuthorNames: h.service.FormatAuthors(magazine.Authors),
			PublishedAt: magazine.PublishedAt.Format("02.01.2006"),
		})
	}

	data := struct {
		Magazines []MagazineDisplay
		Count     int
	}{
		Magazines: magazines,
		Count:     len(magazines),
	}

	t, err := template.New("magazines").Parse(tmpl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	t.Execute(w, data)
}

// handleSearch handles ISBN search functionality
func (h *UIHandler) handleSearch(w http.ResponseWriter, r *http.Request) {
	tmpl := `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Search - Library Management System</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 40px; background-color: #f5f5f5; }
        .container { max-width: 1200px; margin: 0 auto; background-color: white; padding: 30px; border-radius: 8px; box-shadow: 0 2px 10px rgba(0,0,0,0.1); }
        h1 { color: #333; text-align: center; margin-bottom: 30px; }
        .back-link { display: inline-block; margin-bottom: 20px; color: #2c5aa0; text-decoration: none; }
        .back-link:hover { text-decoration: underline; }
        .search-form { background-color: #f8f9fa; padding: 20px; border-radius: 8px; margin-bottom: 30px; text-align: center; }
        .search-input { padding: 10px; font-size: 16px; border: 1px solid #ddd; border-radius: 4px; width: 300px; margin-right: 10px; }
        .search-button { padding: 10px 20px; font-size: 16px; background-color: #2c5aa0; color: white; border: none; border-radius: 4px; cursor: pointer; }
        .search-button:hover { background-color: #1e3d6f; }
        .clear-button { padding: 10px 20px; font-size: 16px; background-color: #6c757d; color: white; border: none; border-radius: 4px; cursor: pointer; margin-left: 10px; }
        .clear-button:hover { background-color: #545b62; }
        .results-section { margin-top: 30px; }
        .no-results { text-align: center; color: #666; font-style: italic; margin: 30px 0; }
        .result-item { border: 1px solid #ddd; padding: 20px; border-radius: 8px; background-color: #fafafa; margin-bottom: 15px; }
        .result-title { font-size: 1.2em; font-weight: bold; color: #333; margin-bottom: 10px; }
        .result-type { display: inline-block; padding: 4px 8px; border-radius: 4px; font-size: 0.8em; font-weight: bold; margin-bottom: 10px; }
        .result-type.book { background-color: #d4edda; color: #155724; }
        .result-type.magazine { background-color: #d1ecf1; color: #0c5460; }
        .result-authors { color: #666; margin-bottom: 8px; }
        .result-isbn { color: #888; font-size: 0.9em; margin-bottom: 8px; }
        .result-description { color: #555; line-height: 1.4; margin-bottom: 8px; }
        .result-published { color: #2c5aa0; font-weight: bold; }
    </style>
</head>
<body>
    <div class="container">
        <a href="/" class="back-link">← Back to Home</a>
        <h1>Search Library</h1>
        
        <div class="search-form">
            <form method="GET" action="/search">
                <select name="searchType" class="search-select" style="padding: 10px; font-size: 16px; border: 1px solid #ddd; border-radius: 4px; margin-right: 10px;">
                    <option value="isbn" {{if eq .SearchType "isbn"}}selected{{end}}>Search by ISBN</option>
                    <option value="email" {{if eq .SearchType "email"}}selected{{end}}>Search by Author Email</option>
                </select>
                <input type="text" name="query" placeholder="{{.Placeholder}}" value="{{.SearchTerm}}" class="search-input">
                <button type="submit" class="search-button">Search</button>
                <a href="/search" class="clear-button">Clear</a>
            </form>
        </div>

        {{if .SearchTerm}}
        <div class="results-section">
            <h2>{{if eq .SearchType "isbn"}}ISBN{{else}}Author Email{{end}} Search Results for "{{.SearchTerm}}"</h2>
            {{if .HasResults}}
                {{range .Books}}
                <div class="result-item">
                    <div class="result-type book">BOOK</div>
                    <div class="result-title">{{.Title}}</div>
                    <div class="result-authors"><strong>Authors:</strong> {{.AuthorNames}}</div>
                    <div class="result-isbn"><strong>ISBN:</strong> {{.ISBN}}</div>
                    <div class="result-description">{{.Description}}</div>
                </div>
                {{end}}
                {{range .Magazines}}
                <div class="result-item">
                    <div class="result-type magazine">MAGAZINE</div>
                    <div class="result-title">{{.Title}}</div>
                    <div class="result-authors"><strong>Authors:</strong> {{.AuthorNames}}</div>
                    <div class="result-isbn"><strong>ISBN:</strong> {{.ISBN}}</div>
                    <div class="result-published"><strong>Published:</strong> {{.PublishedAt}}</div>
                </div>
                {{end}}
            {{else}}
                <div class="no-results">No books or magazines found with ISBN "{{.SearchTerm}}"</div>
            {{end}}
        </div>
        {{else}}
        <div class="results-section">
            <p style="text-align: center; color: #666;">Select a search type and enter your search term above.</p>
            <p style="text-align: center; color: #888; font-size: 0.9em;">
                <strong>ISBN examples:</strong> 5554-5545-4518, 2145-8548-3325, 5454-5587-3210<br>
                <strong>Author email examples:</strong> null-walter@echocat.org, null-mueller@echocat.org, null-lieblich@echocat.org
            </p>
        </div>
        {{end}}
    </div>
</body>
</html>`

	// Get search parameters
	searchTerm := strings.TrimSpace(r.URL.Query().Get("query"))
	searchType := r.URL.Query().Get("searchType")

	// Default to ISBN search if not specified
	if searchType == "" {
		searchType = "isbn"
	}

	// Set placeholder text based on search type
	placeholder := "Enter ISBN (e.g., 5554-5545-4518)"
	if searchType == "email" {
		placeholder = "Enter author email (e.g., null-walter@echocat.org)"
	}

	// Prepare data for template
	data := struct {
		SearchTerm  string
		SearchType  string
		Placeholder string
		HasResults  bool
		Books       []BookDisplay
		Magazines   []MagazineDisplay
	}{
		SearchTerm:  searchTerm,
		SearchType:  searchType,
		Placeholder: placeholder,
	}

	// If there's a search term, perform the search
	if searchTerm != "" {
		var books []Book
		var magazines []Magazine

		if searchType == "isbn" {
			books, magazines = h.service.SearchByISBN(searchTerm)
		} else if searchType == "email" {
			books, magazines = h.service.SearchByAuthorEmail(searchTerm)
		}

		// Convert to display format
		for _, book := range books {
			data.Books = append(data.Books, BookDisplay{
				Title:       book.Title,
				ISBN:        book.ISBN,
				AuthorNames: h.service.FormatAuthors(book.Authors),
				Description: book.Description,
			})
		}

		for _, magazine := range magazines {
			data.Magazines = append(data.Magazines, MagazineDisplay{
				Title:       magazine.Title,
				ISBN:        magazine.ISBN,
				AuthorNames: h.service.FormatAuthors(magazine.Authors),
				PublishedAt: magazine.PublishedAt.Format("02.01.2006"),
			})
		}

		data.HasResults = len(data.Books) > 0 || len(data.Magazines) > 0
	}

	t, err := template.New("search").Parse(tmpl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	t.Execute(w, data)
}

// LibraryItem represents a combined item (book or magazine) for sorting
type LibraryItem struct {
	Title       string
	ISBN        string
	AuthorNames string
	Type        string // "BOOK" or "MAGAZINE"
	Description string // Only for books
	PublishedAt string // Only for magazines
}

// handleAllSorted displays all books and magazines sorted by title
func (h *UIHandler) handleAllSorted(w http.ResponseWriter, r *http.Request) {
	tmpl := `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>All Items Sorted - Library Management System</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 40px; background-color: #f5f5f5; }
        .container { max-width: 1200px; margin: 0 auto; background-color: white; padding: 30px; border-radius: 8px; box-shadow: 0 2px 10px rgba(0,0,0,0.1); }
        h1 { color: #333; text-align: center; margin-bottom: 30px; }
        .back-link { display: inline-block; margin-bottom: 20px; color: #2c5aa0; text-decoration: none; }
        .back-link:hover { text-decoration: underline; }
        .item-grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(400px, 1fr)); gap: 20px; }
        .item-card { border: 1px solid #ddd; padding: 20px; border-radius: 8px; background-color: #fafafa; }
        .item-title { font-size: 1.2em; font-weight: bold; color: #333; margin-bottom: 10px; }
        .item-type { display: inline-block; padding: 4px 8px; border-radius: 4px; font-size: 0.8em; font-weight: bold; margin-bottom: 10px; }
        .item-type.book { background-color: #d4edda; color: #155724; }
        .item-type.magazine { background-color: #d1ecf1; color: #0c5460; }
        .item-authors { color: #666; margin-bottom: 8px; }
        .item-isbn { color: #888; font-size: 0.9em; margin-bottom: 8px; }
        .item-description { color: #555; line-height: 1.4; margin-bottom: 8px; }
        .item-published { color: #2c5aa0; font-weight: bold; }
        .count { text-align: center; margin-bottom: 20px; color: #666; }
    </style>
</head>
<body>
    <div class="container">
        <a href="/" class="back-link">← Back to Home</a>
        <h1>All Library Items (Sorted by Title)</h1>
        
        <div style="background-color: #f8f9fa; padding: 15px; border-radius: 8px; margin-bottom: 20px; text-align: center;">
            <form method="GET" action="/all-sorted">
                <label for="direction" style="margin-right: 10px; font-weight: bold;">Sort Direction:</label>
                <select name="direction" id="direction" style="padding: 8px; font-size: 14px; border: 1px solid #ddd; border-radius: 4px; margin-right: 10px;">
                    <option value="asc" {{if eq .SortDirection "asc"}}selected{{end}}>A-Z (Ascending)</option>
                    <option value="desc" {{if eq .SortDirection "desc"}}selected{{end}}>Z-A (Descending)</option>
                </select>
                <button type="submit" style="padding: 8px 16px; font-size: 14px; background-color: #2c5aa0; color: white; border: none; border-radius: 4px; cursor: pointer;">Apply Sort</button>
            </form>
        </div>
        
        <div class="count">Total: {{.Count}} items ({{.BookCount}} books, {{.MagazineCount}} magazines) - Sorted {{if eq .SortDirection "desc"}}Z-A{{else}}A-Z{{end}}</div>
        
        <div class="item-grid">
            {{range .Items}}
            <div class="item-card">
                <div class="item-type {{if eq .Type "BOOK"}}book{{else}}magazine{{end}}">{{.Type}}</div>
                <div class="item-title">{{.Title}}</div>
                <div class="item-authors"><strong>Authors:</strong> {{.AuthorNames}}</div>
                <div class="item-isbn"><strong>ISBN:</strong> {{.ISBN}}</div>
                {{if eq .Type "BOOK"}}
                    <div class="item-description">{{.Description}}</div>
                {{else}}
                    <div class="item-published"><strong>Published:</strong> {{.PublishedAt}}</div>
                {{end}}
            </div>
            {{end}}
        </div>
    </div>
</body>
</html>`

	// Get sort direction from query parameter (default to ascending)
	sortDirection := r.URL.Query().Get("direction")
	if sortDirection == "" {
		sortDirection = "asc"
	}

	// Get sorted items from service
	ascending := sortDirection == "asc"
	allItems := h.service.GetAllItemsSorted(ascending)

	// Get library statistics
	bookCount, magazineCount, _ := h.service.GetLibraryStats()

	// Prepare data for template
	data := struct {
		Items         []LibraryItem
		Count         int
		BookCount     int
		MagazineCount int
		SortDirection string
	}{
		Items:         allItems,
		Count:         len(allItems),
		BookCount:     bookCount,
		MagazineCount: magazineCount,
		SortDirection: sortDirection,
	}

	t, err := template.New("all-sorted").Parse(tmpl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	t.Execute(w, data)
}

// handleAdd handles adding new books, magazines, and authors
func (h *UIHandler) handleAdd(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		h.handleAddPost(w, r)
		return
	}

	tmpl := `
<!DOCTYPE html>
<html>
<head>
	   <meta charset="UTF-8">
	   <title>Add New Items - Library Management System</title>
	   <style>
	       body { font-family: Arial, sans-serif; margin: 40px; background-color: #f5f5f5; }
	       .container { max-width: 800px; margin: 0 auto; background-color: white; padding: 30px; border-radius: 8px; box-shadow: 0 2px 10px rgba(0,0,0,0.1); }
	       h1 { color: #333; text-align: center; margin-bottom: 30px; }
	       .back-link { display: inline-block; margin-bottom: 20px; color: #2c5aa0; text-decoration: none; }
	       .back-link:hover { text-decoration: underline; }
	       .form-section { background-color: #f8f9fa; padding: 20px; border-radius: 8px; margin-bottom: 20px; }
	       .form-section h2 { color: #333; margin-top: 0; margin-bottom: 20px; }
	       .form-group { margin-bottom: 15px; }
	       .form-group label { display: block; margin-bottom: 5px; font-weight: bold; color: #555; }
	       .form-group input, .form-group textarea, .form-group select { width: 100%; padding: 10px; border: 1px solid #ddd; border-radius: 4px; font-size: 14px; box-sizing: border-box; }
	       .form-group textarea { height: 80px; resize: vertical; }
	       .form-group input[type="date"] { width: auto; }
	       .checkbox-group { display: flex; flex-wrap: wrap; gap: 10px; margin-top: 10px; }
	       .checkbox-item { display: flex; align-items: center; margin-right: 15px; margin-bottom: 5px; }
	       .checkbox-item input[type="checkbox"] { width: auto; margin-right: 5px; }
	       .submit-button { background-color: #28a745; color: white; padding: 12px 24px; border: none; border-radius: 4px; cursor: pointer; font-size: 16px; font-weight: bold; }
	       .submit-button:hover { background-color: #218838; }
	       .error { background-color: #f8d7da; color: #721c24; padding: 10px; border-radius: 4px; margin-bottom: 20px; }
	       .success { background-color: #d4edda; color: #155724; padding: 10px; border-radius: 4px; margin-bottom: 20px; }
	       .author-form { border-top: 2px solid #ddd; padding-top: 20px; margin-top: 20px; }
	       .item-type-selector { text-align: center; margin-bottom: 30px; }
	       .type-button { display: inline-block; margin: 0 10px; padding: 10px 20px; background-color: #6c757d; color: white; text-decoration: none; border-radius: 5px; cursor: pointer; }
	       .type-button.active { background-color: #2c5aa0; }
	       .type-button:hover { background-color: #1e3d6f; }
	       .form-content { display: none; }
	       .form-content.active { display: block; }
	   </style>
	   <script>
	       function showForm(type) {
	           // Hide all forms
	           document.querySelectorAll('.form-content').forEach(form => {
	               form.classList.remove('active');
	           });
	           
	           // Remove active class from all buttons
	           document.querySelectorAll('.type-button').forEach(btn => {
	               btn.classList.remove('active');
	           });
	           
	           // Show selected form
	           document.getElementById(type + '-form').classList.add('active');
	           document.querySelector('[onclick="showForm(\'' + type + '\')"]').classList.add('active');
	       }
	       
	       function updatePublishedDate() {
	           const today = new Date();
	           const dateStr = today.getFullYear() + '-' +
	                          String(today.getMonth() + 1).padStart(2, '0') + '-' +
	                          String(today.getDate()).padStart(2, '0');
	           document.getElementById('magazine-published').value = dateStr;
	       }
	       
	       function toggleNewAuthor(formType) {
	           const newAuthorSection = document.getElementById(formType + '-new-author');
	           const checkbox = document.getElementById(formType + '-add-new-author');
	           
	           if (checkbox.checked) {
	               newAuthorSection.style.display = 'block';
	           } else {
	               newAuthorSection.style.display = 'none';
	               // Clear the new author fields
	               document.getElementById(formType + '-new-email').value = '';
	               document.getElementById(formType + '-new-firstname').value = '';
	               document.getElementById(formType + '-new-lastname').value = '';
	           }
	       }
	       
	       window.onload = function() {
	           showForm('book'); // Show book form by default
	           updatePublishedDate(); // Set today's date for magazine
	       };
	   </script>
</head>
<body>
	   <div class="container">
	       <a href="/" class="back-link">← Back to Home</a>
	       <h1>Add New Items to Library</h1>
	       
	       {{if .Error}}
	       <div class="error">{{.Error}}</div>
	       {{end}}
	       
	       {{if .Success}}
	       <div class="success">{{.Success}}</div>
	       {{end}}
	       
	       <div class="item-type-selector">
	           <span class="type-button active" onclick="showForm('book')">Add Book</span>
	           <span class="type-button" onclick="showForm('magazine')">Add Magazine</span>
	           <span class="type-button" onclick="showForm('author')">Add Author</span>
	       </div>
	       
	       <!-- Book Form -->
	       <div id="book-form" class="form-content active">
	           <div class="form-section">
	               <h2>Add New Book</h2>
	               <form method="POST" action="/add">
	                   <input type="hidden" name="type" value="book">
	                   
	                   <div class="form-group">
	                       <label for="book-title">Title *</label>
	                       <input type="text" id="book-title" name="title" required>
	                   </div>
	                   
	                   <div class="form-group">
	                       <label for="book-isbn">ISBN *</label>
	                       <input type="text" id="book-isbn" name="isbn" placeholder="e.g., 5554-5545-4518" required>
	                   </div>
	                   
	                   <div class="form-group">
	                       <label for="book-description">Description</label>
	                       <textarea id="book-description" name="description" placeholder="Brief description of the book"></textarea>
	                   </div>
	                   
	                   <div class="form-group">
	                       <label>Authors *</label>
	                       <div class="checkbox-group">
	                           {{range .Authors}}
	                           <div class="checkbox-item">
	                               <input type="checkbox" id="book-author-{{.Email}}" name="authors" value="{{.Email}}">
	                               <label for="book-author-{{.Email}}">{{.FirstName}} {{.LastName}} ({{.Email}})</label>
	                           </div>
	                           {{end}}
	                       </div>
	                       
	                       <div style="margin-top: 15px; padding: 10px; background-color: #e9ecef; border-radius: 4px;">
	                           <div class="checkbox-item">
	                               <input type="checkbox" id="book-add-new-author" onchange="toggleNewAuthor('book')">
	                               <label for="book-add-new-author" style="font-weight: bold; color: #28a745;">Add New Author</label>
	                           </div>
	                           
	                           <div id="book-new-author" style="display: none; margin-top: 15px; padding: 15px; background-color: white; border-radius: 4px; border: 1px solid #ddd;">
	                               <h4 style="margin-top: 0; color: #333;">New Author Details</h4>
	                               <div class="form-group">
	                                   <label for="book-new-email">Email *</label>
	                                   <input type="email" id="book-new-email" name="new_author_email" placeholder="e.g., author@example.com">
	                               </div>
	                               <div class="form-group">
	                                   <label for="book-new-firstname">First Name *</label>
	                                   <input type="text" id="book-new-firstname" name="new_author_firstname">
	                               </div>
	                               <div class="form-group">
	                                   <label for="book-new-lastname">Last Name *</label>
	                                   <input type="text" id="book-new-lastname" name="new_author_lastname">
	                               </div>
	                           </div>
	                       </div>
	                   </div>
	                   
	                   <button type="submit" class="submit-button">Add Book</button>
	               </form>
	           </div>
	       </div>
	       
	       <!-- Magazine Form -->
	       <div id="magazine-form" class="form-content">
	           <div class="form-section">
	               <h2>Add New Magazine</h2>
	               <form method="POST" action="/add">
	                   <input type="hidden" name="type" value="magazine">
	                   
	                   <div class="form-group">
	                       <label for="magazine-title">Title *</label>
	                       <input type="text" id="magazine-title" name="title" required>
	                   </div>
	                   
	                   <div class="form-group">
	                       <label for="magazine-isbn">ISBN *</label>
	                       <input type="text" id="magazine-isbn" name="isbn" placeholder="e.g., 5454-5587-3210" required>
	                   </div>
	                   
	                   <div class="form-group">
	                       <label for="magazine-published">Published Date *</label>
	                       <input type="date" id="magazine-published" name="published" required>
	                   </div>
	                   
	                   <div class="form-group">
	                       <label>Authors *</label>
	                       <div class="checkbox-group">
	                           {{range .Authors}}
	                           <div class="checkbox-item">
	                               <input type="checkbox" id="magazine-author-{{.Email}}" name="authors" value="{{.Email}}">
	                               <label for="magazine-author-{{.Email}}">{{.FirstName}} {{.LastName}} ({{.Email}})</label>
	                           </div>
	                           {{end}}
	                       </div>
	                       
	                       <div style="margin-top: 15px; padding: 10px; background-color: #e9ecef; border-radius: 4px;">
	                           <div class="checkbox-item">
	                               <input type="checkbox" id="magazine-add-new-author" onchange="toggleNewAuthor('magazine')">
	                               <label for="magazine-add-new-author" style="font-weight: bold; color: #28a745;">Add New Author</label>
	                           </div>
	                           
	                           <div id="magazine-new-author" style="display: none; margin-top: 15px; padding: 15px; background-color: white; border-radius: 4px; border: 1px solid #ddd;">
	                               <h4 style="margin-top: 0; color: #333;">New Author Details</h4>
	                               <div class="form-group">
	                                   <label for="magazine-new-email">Email *</label>
	                                   <input type="email" id="magazine-new-email" name="new_author_email" placeholder="e.g., author@example.com">
	                               </div>
	                               <div class="form-group">
	                                   <label for="magazine-new-firstname">First Name *</label>
	                                   <input type="text" id="magazine-new-firstname" name="new_author_firstname">
	                               </div>
	                               <div class="form-group">
	                                   <label for="magazine-new-lastname">Last Name *</label>
	                                   <input type="text" id="magazine-new-lastname" name="new_author_lastname">
	                               </div>
	                           </div>
	                       </div>
	                   </div>
	                   
	                   <button type="submit" class="submit-button">Add Magazine</button>
	               </form>
	           </div>
	       </div>
	       
	       <!-- Author Form -->
	       <div id="author-form" class="form-content">
	           <div class="form-section">
	               <h2>Add New Author</h2>
	               <form method="POST" action="/add">
	                   <input type="hidden" name="type" value="author">
	                   
	                   <div class="form-group">
	                       <label for="author-email">Email *</label>
	                       <input type="email" id="author-email" name="email" placeholder="e.g., author@example.com" required>
	                   </div>
	                   
	                   <div class="form-group">
	                       <label for="author-firstname">First Name *</label>
	                       <input type="text" id="author-firstname" name="firstname" required>
	                   </div>
	                   
	                   <div class="form-group">
	                       <label for="author-lastname">Last Name *</label>
	                       <input type="text" id="author-lastname" name="lastname" required>
	                   </div>
	                   
	                   <button type="submit" class="submit-button">Add Author</button>
	               </form>
	           </div>
	       </div>
	       
	       <div style="margin-top: 30px; padding: 15px; background-color: #e9ecef; border-radius: 5px;">
	           <p><strong>Instructions:</strong></p>
	           <ul>
	               <li>Fields marked with * are required</li>
	               <li>When adding books or magazines, you must select at least one existing author OR create a new author</li>
	               <li>To create a new author while adding a book/magazine, check "Add New Author" and fill in the author details</li>
	               <li>You can also use the standalone "Add Author" form to add authors separately</li>
	               <li>ISBN must be unique across all books and magazines</li>
	               <li>All data will be saved to CSV files for persistence</li>
	           </ul>
	       </div>
	   </div>
</body>
</html>`

	// Get all authors for the checkboxes
	authors := h.service.GetAuthorsList()

	data := struct {
		Authors []Author
		Error   string
		Success string
	}{
		Authors: authors,
	}

	t, err := template.New("add").Parse(tmpl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	t.Execute(w, data)
}

// handleAddPost processes the form submission for adding new items
func (h *UIHandler) handleAddPost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		h.renderAddWithError(w, "Error parsing form data")
		return
	}

	itemType := r.FormValue("type")
	var resultErr error
	var successMsg string

	switch itemType {
	case "book":
		title := r.FormValue("title")
		isbn := r.FormValue("isbn")
		description := r.FormValue("description")
		authors := r.Form["authors"] // This gets all selected checkboxes

		// Check if we need to create a new author
		newAuthorEmail := strings.TrimSpace(r.FormValue("new_author_email"))
		if newAuthorEmail != "" {
			newAuthorFirstName := strings.TrimSpace(r.FormValue("new_author_firstname"))
			newAuthorLastName := strings.TrimSpace(r.FormValue("new_author_lastname"))

			// Create the new author first
			resultErr = h.service.AddAuthor(newAuthorEmail, newAuthorFirstName, newAuthorLastName)
			if resultErr != nil {
				break // Exit early if author creation failed
			}

			// Add the new author to the authors list
			authors = append(authors, newAuthorEmail)
		}

		resultErr = h.service.AddBook(title, isbn, description, authors)
		if resultErr == nil {
			successMsg = fmt.Sprintf("Book '%s' added successfully!", title)
		}

	case "magazine":
		title := r.FormValue("title")
		isbn := r.FormValue("isbn")
		publishedStr := r.FormValue("published")
		authors := r.Form["authors"]

		// Check if we need to create a new author
		newAuthorEmail := strings.TrimSpace(r.FormValue("new_author_email"))
		if newAuthorEmail != "" {
			newAuthorFirstName := strings.TrimSpace(r.FormValue("new_author_firstname"))
			newAuthorLastName := strings.TrimSpace(r.FormValue("new_author_lastname"))

			// Create the new author first
			resultErr = h.service.AddAuthor(newAuthorEmail, newAuthorFirstName, newAuthorLastName)
			if resultErr != nil {
				break // Exit early if author creation failed
			}

			// Add the new author to the authors list
			authors = append(authors, newAuthorEmail)
		}

		// Parse the published date
		publishedAt, err := time.Parse("2006-01-02", publishedStr)
		if err != nil {
			resultErr = fmt.Errorf("invalid published date format")
		} else {
			resultErr = h.service.AddMagazine(title, isbn, publishedAt, authors)
			if resultErr == nil {
				successMsg = fmt.Sprintf("Magazine '%s' added successfully!", title)
			}
		}

	case "author":
		email := r.FormValue("email")
		firstName := r.FormValue("firstname")
		lastName := r.FormValue("lastname")

		resultErr = h.service.AddAuthor(email, firstName, lastName)
		if resultErr == nil {
			successMsg = fmt.Sprintf("Author '%s %s' added successfully!", firstName, lastName)
		}

	default:
		resultErr = fmt.Errorf("invalid item type")
	}

	if resultErr != nil {
		h.renderAddWithError(w, resultErr.Error())
	} else {
		h.renderAddWithSuccess(w, successMsg)
	}
}

// renderAddWithError renders the add form with an error message
func (h *UIHandler) renderAddWithError(w http.ResponseWriter, errorMsg string) {
	authors := h.service.GetAuthorsList()

	data := struct {
		Authors []Author
		Error   string
		Success string
	}{
		Authors: authors,
		Error:   errorMsg,
	}

	h.renderAddForm(w, data)
}

// renderAddWithSuccess renders the add form with a success message
func (h *UIHandler) renderAddWithSuccess(w http.ResponseWriter, successMsg string) {
	authors := h.service.GetAuthorsList()

	data := struct {
		Authors []Author
		Error   string
		Success string
	}{
		Authors: authors,
		Success: successMsg,
	}

	h.renderAddForm(w, data)
}

// renderAddForm renders the add form with the given data
func (h *UIHandler) renderAddForm(w http.ResponseWriter, data interface{}) {
	tmpl := `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Add New Items - Library Management System</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 40px; background-color: #f5f5f5; }
        .container { max-width: 800px; margin: 0 auto; background-color: white; padding: 30px; border-radius: 8px; box-shadow: 0 2px 10px rgba(0,0,0,0.1); }
        h1 { color: #333; text-align: center; margin-bottom: 30px; }
        .back-link { display: inline-block; margin-bottom: 20px; color: #2c5aa0; text-decoration: none; }
        .back-link:hover { text-decoration: underline; }
        .form-section { background-color: #f8f9fa; padding: 20px; border-radius: 8px; margin-bottom: 20px; }
        .form-section h2 { color: #333; margin-top: 0; margin-bottom: 20px; }
        .form-group { margin-bottom: 15px; }
        .form-group label { display: block; margin-bottom: 5px; font-weight: bold; color: #555; }
        .form-group input, .form-group textarea, .form-group select { width: 100%; padding: 10px; border: 1px solid #ddd; border-radius: 4px; font-size: 14px; box-sizing: border-box; }
        .form-group textarea { height: 80px; resize: vertical; }
        .form-group input[type="date"] { width: auto; }
        .checkbox-group { display: flex; flex-wrap: wrap; gap: 10px; margin-top: 10px; }
        .checkbox-item { display: flex; align-items: center; margin-right: 15px; margin-bottom: 5px; }
        .checkbox-item input[type="checkbox"] { width: auto; margin-right: 5px; }
        .submit-button { background-color: #28a745; color: white; padding: 12px 24px; border: none; border-radius: 4px; cursor: pointer; font-size: 16px; font-weight: bold; }
        .submit-button:hover { background-color: #218838; }
        .error { background-color: #f8d7da; color: #721c24; padding: 10px; border-radius: 4px; margin-bottom: 20px; }
        .success { background-color: #d4edda; color: #155724; padding: 10px; border-radius: 4px; margin-bottom: 20px; }
        .author-form { border-top: 2px solid #ddd; padding-top: 20px; margin-top: 20px; }
        .item-type-selector { text-align: center; margin-bottom: 30px; }
        .type-button { display: inline-block; margin: 0 10px; padding: 10px 20px; background-color: #6c757d; color: white; text-decoration: none; border-radius: 5px; cursor: pointer; }
        .type-button.active { background-color: #2c5aa0; }
        .type-button:hover { background-color: #1e3d6f; }
        .form-content { display: none; }
        .form-content.active { display: block; }
    </style>
    <script>
        function showForm(type) {
            // Hide all forms
            document.querySelectorAll('.form-content').forEach(form => {
                form.classList.remove('active');
            });
            
            // Remove active class from all buttons
            document.querySelectorAll('.type-button').forEach(btn => {
                btn.classList.remove('active');
            });
            
            // Show selected form
            document.getElementById(type + '-form').classList.add('active');
            document.querySelector('[onclick="showForm(\'' + type + '\')"]').classList.add('active');
        }
        
        function updatePublishedDate() {
            const today = new Date();
            const dateStr = today.getFullYear() + '-' +
                           String(today.getMonth() + 1).padStart(2, '0') + '-' +
                           String(today.getDate()).padStart(2, '0');
            document.getElementById('magazine-published').value = dateStr;
        }
        
        function toggleNewAuthor(formType) {
            const newAuthorSection = document.getElementById(formType + '-new-author');
            const checkbox = document.getElementById(formType + '-add-new-author');
            
            if (checkbox.checked) {
                newAuthorSection.style.display = 'block';
            } else {
                newAuthorSection.style.display = 'none';
                // Clear the new author fields
                document.getElementById(formType + '-new-email').value = '';
                document.getElementById(formType + '-new-firstname').value = '';
                document.getElementById(formType + '-new-lastname').value = '';
            }
        }
        
        window.onload = function() {
            showForm('book'); // Show book form by default
            updatePublishedDate(); // Set today's date for magazine
        };
    </script>
</head>
<body>
    <div class="container">
        <a href="/" class="back-link">← Back to Home</a>
        <h1>Add New Items to Library</h1>
        
        {{if .Error}}
        <div class="error">{{.Error}}</div>
        {{end}}
        
        {{if .Success}}
        <div class="success">{{.Success}}</div>
        {{end}}
        
        <div class="item-type-selector">
            <span class="type-button active" onclick="showForm('book')">Add Book</span>
            <span class="type-button" onclick="showForm('magazine')">Add Magazine</span>
            <span class="type-button" onclick="showForm('author')">Add Author</span>
        </div>
        
        <!-- Book Form -->
        <div id="book-form" class="form-content active">
            <div class="form-section">
                <h2>Add New Book</h2>
                <form method="POST" action="/add">
                    <input type="hidden" name="type" value="book">
                    
                    <div class="form-group">
                        <label for="book-title">Title *</label>
                        <input type="text" id="book-title" name="title" required>
                    </div>
                    
                    <div class="form-group">
                        <label for="book-isbn">ISBN *</label>
                        <input type="text" id="book-isbn" name="isbn" placeholder="e.g., 5554-5545-4518" required>
                    </div>
                    
                    <div class="form-group">
                        <label for="book-description">Description</label>
                        <textarea id="book-description" name="description" placeholder="Brief description of the book"></textarea>
                    </div>
                    
                    <div class="form-group">
                        <label>Authors *</label>
                        <div class="checkbox-group">
                            {{range .Authors}}
                            <div class="checkbox-item">
                                <input type="checkbox" id="book-author-{{.Email}}" name="authors" value="{{.Email}}">
                                <label for="book-author-{{.Email}}">{{.FirstName}} {{.LastName}} ({{.Email}})</label>
                            </div>
                            {{end}}
                        </div>
                        
                        <div style="margin-top: 15px; padding: 10px; background-color: #e9ecef; border-radius: 4px;">
                            <div class="checkbox-item">
                                <input type="checkbox" id="book-add-new-author" onchange="toggleNewAuthor('book')">
                                <label for="book-add-new-author" style="font-weight: bold; color: #28a745;">Add New Author</label>
                            </div>
                            
                            <div id="book-new-author" style="display: none; margin-top: 15px; padding: 15px; background-color: white; border-radius: 4px; border: 1px solid #ddd;">
                                <h4 style="margin-top: 0; color: #333;">New Author Details</h4>
                                <div class="form-group">
                                    <label for="book-new-email">Email *</label>
                                    <input type="email" id="book-new-email" name="new_author_email" placeholder="e.g., author@example.com">
                                </div>
                                <div class="form-group">
                                    <label for="book-new-firstname">First Name *</label>
                                    <input type="text" id="book-new-firstname" name="new_author_firstname">
                                </div>
                                <div class="form-group">
                                    <label for="book-new-lastname">Last Name *</label>
                                    <input type="text" id="book-new-lastname" name="new_author_lastname">
                                </div>
                            </div>
                        </div>
                    </div>
                    
                    <button type="submit" class="submit-button">Add Book</button>
                </form>
            </div>
        </div>
        
        <!-- Magazine Form -->
        <div id="magazine-form" class="form-content">
            <div class="form-section">
                <h2>Add New Magazine</h2>
                <form method="POST" action="/add">
                    <input type="hidden" name="type" value="magazine">
                    
                    <div class="form-group">
                        <label for="magazine-title">Title *</label>
                        <input type="text" id="magazine-title" name="title" required>
                    </div>
                    
                    <div class="form-group">
                        <label for="magazine-isbn">ISBN *</label>
                        <input type="text" id="magazine-isbn" name="isbn" placeholder="e.g., 5454-5587-3210" required>
                    </div>
                    
                    <div class="form-group">
                        <label for="magazine-published">Published Date *</label>
                        <input type="date" id="magazine-published" name="published" required>
                    </div>
                    
                    <div class="form-group">
                        <label>Authors *</label>
                        <div class="checkbox-group">
                            {{range .Authors}}
                            <div class="checkbox-item">
                                <input type="checkbox" id="magazine-author-{{.Email}}" name="authors" value="{{.Email}}">
                                <label for="magazine-author-{{.Email}}">{{.FirstName}} {{.LastName}} ({{.Email}})</label>
                            </div>
                            {{end}}
                        </div>
                        
                        <div style="margin-top: 15px; padding: 10px; background-color: #e9ecef; border-radius: 4px;">
                            <div class="checkbox-item">
                                <input type="checkbox" id="magazine-add-new-author" onchange="toggleNewAuthor('magazine')">
                                <label for="magazine-add-new-author" style="font-weight: bold; color: #28a745;">Add New Author</label>
                            </div>
                            
                            <div id="magazine-new-author" style="display: none; margin-top: 15px; padding: 15px; background-color: white; border-radius: 4px; border: 1px solid #ddd;">
                                <h4 style="margin-top: 0; color: #333;">New Author Details</h4>
                                <div class="form-group">
                                    <label for="magazine-new-email">Email *</label>
                                    <input type="email" id="magazine-new-email" name="new_author_email" placeholder="e.g., author@example.com">
                                </div>
                                <div class="form-group">
                                    <label for="magazine-new-firstname">First Name *</label>
                                    <input type="text" id="magazine-new-firstname" name="new_author_firstname">
                                </div>
                                <div class="form-group">
                                    <label for="magazine-new-lastname">Last Name *</label>
                                    <input type="text" id="magazine-new-lastname" name="new_author_lastname">
                                </div>
                            </div>
                        </div>
                    </div>
                    
                    <button type="submit" class="submit-button">Add Magazine</button>
                </form>
            </div>
        </div>
        
        <!-- Author Form -->
        <div id="author-form" class="form-content">
            <div class="form-section">
                <h2>Add New Author</h2>
                <form method="POST" action="/add">
                    <input type="hidden" name="type" value="author">
                    
                    <div class="form-group">
                        <label for="author-email">Email *</label>
                        <input type="email" id="author-email" name="email" placeholder="e.g., author@example.com" required>
                    </div>
                    
                    <div class="form-group">
                        <label for="author-firstname">First Name *</label>
                        <input type="text" id="author-firstname" name="firstname" required>
                    </div>
                    
                    <div class="form-group">
                        <label for="author-lastname">Last Name *</label>
                        <input type="text" id="author-lastname" name="lastname" required>
                    </div>
                    
                    <button type="submit" class="submit-button">Add Author</button>
                </form>
            </div>
        </div>
        
        <div style="margin-top: 30px; padding: 15px; background-color: #e9ecef; border-radius: 5px;">
            <p><strong>Instructions:</strong></p>
            <ul>
                <li>Fields marked with * are required</li>
                <li>When adding books or magazines, you must select at least one existing author OR create a new author</li>
                <li>To create a new author while adding a book/magazine, check "Add New Author" and fill in the author details</li>
                <li>You can also use the standalone "Add Author" form to add authors separately</li>
                <li>ISBN must be unique across all books and magazines</li>
                <li>All data will be saved to CSV files for persistence</li>
            </ul>
        </div>
    </div>
</body>
</html>`

	t, err := template.New("add").Parse(tmpl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	t.Execute(w, data)
}
