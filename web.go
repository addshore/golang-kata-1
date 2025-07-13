package main

import (
	"html/template"
	"net/http"
	"strings"
)

type WebServer struct {
	library *Library
}

type PageData struct {
	Title     string
	Books     []Book
	Magazines []Magazine
	Items     []LibraryItem
	Message   string
	Error     string
}

func NewWebServer(library *Library) *WebServer {
	return &WebServer{library: library}
}

func (ws *WebServer) Start(port string) {
	http.HandleFunc("/", ws.handleHome)
	http.HandleFunc("/books", ws.handleBooks)
	http.HandleFunc("/magazines", ws.handleMagazines)
	http.HandleFunc("/all", ws.handleAll)
	http.HandleFunc("/sorted", ws.handleSorted)
	http.HandleFunc("/search", ws.handleSearch)
	http.HandleFunc("/add", ws.handleAdd)
	http.HandleFunc("/add-book", ws.handleAddBook)
	http.HandleFunc("/add-magazine", ws.handleAddMagazine)
	
	println("Starting web server on http://localhost:" + port)
	http.ListenAndServe(":"+port, nil)
}

func (ws *WebServer) handleHome(w http.ResponseWriter, r *http.Request) {
	data := PageData{
		Title: "Library Management System",
	}
	ws.renderTemplate(w, "home", data)
}

func (ws *WebServer) handleBooks(w http.ResponseWriter, r *http.Request) {
	data := PageData{
		Title: "All Books",
		Books: ws.library.Books,
	}
	ws.renderTemplate(w, "books", data)
}

func (ws *WebServer) handleMagazines(w http.ResponseWriter, r *http.Request) {
	data := PageData{
		Title: "All Magazines",
		Magazines: ws.library.Magazines,
	}
	ws.renderTemplate(w, "magazines", data)
}

func (ws *WebServer) handleAll(w http.ResponseWriter, r *http.Request) {
	data := PageData{
		Title: "All Items",
		Books: ws.library.Books,
		Magazines: ws.library.Magazines,
	}
	ws.renderTemplate(w, "all", data)
}

func (ws *WebServer) handleSorted(w http.ResponseWriter, r *http.Request) {
	ascending := r.URL.Query().Get("order") != "desc"
	items := ws.library.GetAllItemsSorted(ascending)
	
	orderText := "A-Z"
	if !ascending {
		orderText = "Z-A"
	}
	
	data := PageData{
		Title: "All Items Sorted by Title (" + orderText + ")",
		Items: items,
	}
	ws.renderTemplate(w, "sorted", data)
}

func (ws *WebServer) handleSearch(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		data := PageData{
			Title: "Search Library",
		}
		ws.renderTemplate(w, "search", data)
		return
	}
	
	// Handle POST request
	searchType := r.FormValue("search_type")
	searchValue := strings.TrimSpace(r.FormValue("search_value"))
	
	data := PageData{
		Title: "Search Results",
	}
	
	if searchValue == "" {
		data.Error = "Search value cannot be empty"
		ws.renderTemplate(w, "search", data)
		return
	}
	
	switch searchType {
	case "isbn":
		book := ws.library.FindBookByISBN(searchValue)
		magazine := ws.library.FindMagazineByISBN(searchValue)
		
		if book != nil {
			data.Books = []Book{*book}
		}
		if magazine != nil {
			data.Magazines = []Magazine{*magazine}
		}
		
		if book == nil && magazine == nil {
			data.Error = "No items found with ISBN: " + searchValue
		}
		
	case "email":
		books := ws.library.FindBooksByAuthorEmail(searchValue)
		magazines := ws.library.FindMagazinesByAuthorEmail(searchValue)
		
		data.Books = books
		data.Magazines = magazines
		
		if len(books) == 0 && len(magazines) == 0 {
			data.Error = "No items found with author email: " + searchValue
		}
	}
	
	ws.renderTemplate(w, "search", data)
}

func (ws *WebServer) handleAdd(w http.ResponseWriter, r *http.Request) {
	data := PageData{
		Title: "Add New Item",
	}
	ws.renderTemplate(w, "add", data)
}

func (ws *WebServer) handleAddBook(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		data := PageData{
			Title: "Add New Book",
		}
		ws.renderTemplate(w, "add-book", data)
		return
	}
	
	// Handle POST request
	title := strings.TrimSpace(r.FormValue("title"))
	isbn := strings.TrimSpace(r.FormValue("isbn"))
	description := strings.TrimSpace(r.FormValue("description"))
	authorEmails := strings.TrimSpace(r.FormValue("author_emails"))
	authorNames := strings.TrimSpace(r.FormValue("author_names"))
	
	data := PageData{
		Title: "Add New Book",
	}
	
	// Validate input
	if title == "" || isbn == "" || description == "" || authorEmails == "" || authorNames == "" {
		data.Error = "All fields are required"
		ws.renderTemplate(w, "add-book", data)
		return
	}
	
	// Check if ISBN already exists
	if ws.library.FindBookByISBN(isbn) != nil || ws.library.FindMagazineByISBN(isbn) != nil {
		data.Error = "ISBN already exists in the library"
		ws.renderTemplate(w, "add-book", data)
		return
	}
	
	// Process authors
	emails := strings.Split(authorEmails, ",")
	names := strings.Split(authorNames, ",")
	
	if len(emails) != len(names) {
		data.Error = "Number of emails and names must match"
		ws.renderTemplate(w, "add-book", data)
		return
	}
	
	var authors []Author
	for i, email := range emails {
		email = strings.TrimSpace(email)
		nameParts := strings.Fields(strings.TrimSpace(names[i]))
		
		if len(nameParts) < 2 {
			data.Error = "Each author name must have at least first and last name"
			ws.renderTemplate(w, "add-book", data)
			return
		}
		
		firstName := nameParts[0]
		lastName := strings.Join(nameParts[1:], " ")
		
		author := ws.library.GetOrCreateAuthor(email, firstName, lastName)
		authors = append(authors, author)
	}
	
	// Create and add book
	book := Book{
		Title:       title,
		ISBN:        isbn,
		Authors:     authors,
		Description: description,
	}
	
	ws.library.AddBook(book)
	
	// Save to CSV
	if err := ws.library.SaveToCSV("resources/authors.csv", "resources/books.csv", "resources/magazines.csv"); err != nil {
		data.Error = "Error saving to CSV files: " + err.Error()
		ws.renderTemplate(w, "add-book", data)
		return
	}
	
	data.Message = "Book '" + title + "' has been successfully added to the library!"
	ws.renderTemplate(w, "add-book", data)
}

func (ws *WebServer) handleAddMagazine(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		data := PageData{
			Title: "Add New Magazine",
		}
		ws.renderTemplate(w, "add-magazine", data)
		return
	}
	
	// Handle POST request
	title := strings.TrimSpace(r.FormValue("title"))
	isbn := strings.TrimSpace(r.FormValue("isbn"))
	publishedAt := strings.TrimSpace(r.FormValue("published_at"))
	authorEmails := strings.TrimSpace(r.FormValue("author_emails"))
	authorNames := strings.TrimSpace(r.FormValue("author_names"))
	
	data := PageData{
		Title: "Add New Magazine",
	}
	
	// Validate input
	if title == "" || isbn == "" || publishedAt == "" || authorEmails == "" || authorNames == "" {
		data.Error = "All fields are required"
		ws.renderTemplate(w, "add-magazine", data)
		return
	}
	
	// Check if ISBN already exists
	if ws.library.FindBookByISBN(isbn) != nil || ws.library.FindMagazineByISBN(isbn) != nil {
		data.Error = "ISBN already exists in the library"
		ws.renderTemplate(w, "add-magazine", data)
		return
	}
	
	// Process authors
	emails := strings.Split(authorEmails, ",")
	names := strings.Split(authorNames, ",")
	
	if len(emails) != len(names) {
		data.Error = "Number of emails and names must match"
		ws.renderTemplate(w, "add-magazine", data)
		return
	}
	
	var authors []Author
	for i, email := range emails {
		email = strings.TrimSpace(email)
		nameParts := strings.Fields(strings.TrimSpace(names[i]))
		
		if len(nameParts) < 2 {
			data.Error = "Each author name must have at least first and last name"
			ws.renderTemplate(w, "add-magazine", data)
			return
		}
		
		firstName := nameParts[0]
		lastName := strings.Join(nameParts[1:], " ")
		
		author := ws.library.GetOrCreateAuthor(email, firstName, lastName)
		authors = append(authors, author)
	}
	
	// Create and add magazine
	magazine := Magazine{
		Title:       title,
		ISBN:        isbn,
		Authors:     authors,
		PublishedAt: publishedAt,
	}
	
	ws.library.AddMagazine(magazine)
	
	// Save to CSV
	if err := ws.library.SaveToCSV("resources/authors.csv", "resources/books.csv", "resources/magazines.csv"); err != nil {
		data.Error = "Error saving to CSV files: " + err.Error()
		ws.renderTemplate(w, "add-magazine", data)
		return
	}
	
	data.Message = "Magazine '" + title + "' has been successfully added to the library!"
	ws.renderTemplate(w, "add-magazine", data)
}

func (ws *WebServer) renderTemplate(w http.ResponseWriter, tmpl string, data PageData) {
	funcMap := template.FuncMap{
		"add": func(a, b int) int {
			return a + b
		},
		"formatAuthors": func(authors []Author) string {
			return FormatAuthors(authors)
		},
	}
	
	t, err := template.New(tmpl).Funcs(funcMap).Parse(ws.getTemplate(tmpl))
	if err != nil {
		http.Error(w, "Template error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	if err := t.Execute(w, data); err != nil {
		http.Error(w, "Template execution error: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func (ws *WebServer) getTemplate(name string) string {
	switch name {
	case "home":
		return ws.getHomeTemplate()
	case "books":
		return ws.getBooksTemplate()
	case "magazines":
		return ws.getMagazinesTemplate()
	case "all":
		return ws.getAllTemplate()
	case "sorted":
		return ws.getSortedTemplate()
	case "search":
		return ws.getSearchTemplate()
	case "add":
		return ws.getAddTemplate()
	case "add-book":
		return ws.getAddBookTemplate()
	case "add-magazine":
		return ws.getAddMagazineTemplate()
	default:
		return ws.getHomeTemplate()
	}
}

func (ws *WebServer) getHomeTemplate() string {
	return `
<!DOCTYPE html>
<html>
<head>
    <title>{{.Title}}</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 40px; }
        .nav { margin-bottom: 30px; }
        .nav a { margin-right: 15px; text-decoration: none; color: #0066cc; }
        .nav a:hover { text-decoration: underline; }
        h1 { color: #333; }
    </style>
</head>
<body>
    <h1>{{.Title}}</h1>
    <div class="nav">
        <a href="/">Home</a>
        <a href="/books">All Books</a>
        <a href="/magazines">All Magazines</a>
        <a href="/all">All Items</a>
        <a href="/sorted">Sorted Items</a>
        <a href="/search">Search</a>
        <a href="/add">Add New Item</a>
    </div>
    
    <h2>Welcome to the Library Management System</h2>
    <p>Use the navigation above to:</p>
    <ul>
        <li><strong>View Books:</strong> See all books in the library</li>
        <li><strong>View Magazines:</strong> See all magazines in the library</li>
        <li><strong>View All Items:</strong> See both books and magazines</li>
        <li><strong>Sorted Items:</strong> View all items sorted by title (A-Z or Z-A)</li>
        <li><strong>Search:</strong> Search by ISBN or author email</li>
        <li><strong>Add New Item:</strong> Add books or magazines to the library</li>
    </ul>
</body>
</html>
`
}

func (ws *WebServer) getBooksTemplate() string {
	return `
<!DOCTYPE html>
<html>
<head>
    <title>{{.Title}}</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 40px; }
        .nav { margin-bottom: 30px; }
        .nav a { margin-right: 15px; text-decoration: none; color: #0066cc; }
        .nav a:hover { text-decoration: underline; }
        .item { margin-bottom: 20px; padding: 15px; border: 1px solid #ddd; }
        .title { font-weight: bold; font-size: 18px; color: #333; }
        .isbn { color: #666; }
        .authors { color: #0066cc; }
        .description { margin-top: 10px; line-height: 1.5; }
    </style>
</head>
<body>
    <h1>{{.Title}}</h1>
    <div class="nav">
        <a href="/">Home</a>
        <a href="/books">All Books</a>
        <a href="/magazines">All Magazines</a>
        <a href="/all">All Items</a>
        <a href="/sorted">Sorted Items</a>
        <a href="/search">Search</a>
        <a href="/add">Add New Item</a>
    </div>
    
    {{range $index, $book := .Books}}
    <div class="item">
        <div class="title">{{add $index 1}}. {{$book.Title}}</div>
        <div class="isbn">ISBN: {{$book.ISBN}}</div>
        <div class="authors">Authors: {{formatAuthors $book.Authors}}</div>
        <div class="description">{{$book.Description}}</div>
    </div>
    {{end}}
</body>
</html>
`
}

func (ws *WebServer) getMagazinesTemplate() string {
	return `
<!DOCTYPE html>
<html>
<head>
    <title>{{.Title}}</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 40px; }
        .nav { margin-bottom: 30px; }
        .nav a { margin-right: 15px; text-decoration: none; color: #0066cc; }
        .nav a:hover { text-decoration: underline; }
        .item { margin-bottom: 20px; padding: 15px; border: 1px solid #ddd; }
        .title { font-weight: bold; font-size: 18px; color: #333; }
        .isbn { color: #666; }
        .authors { color: #0066cc; }
        .published { color: #888; }
    </style>
</head>
<body>
    <h1>{{.Title}}</h1>
    <div class="nav">
        <a href="/">Home</a>
        <a href="/books">All Books</a>
        <a href="/magazines">All Magazines</a>
        <a href="/all">All Items</a>
        <a href="/sorted">Sorted Items</a>
        <a href="/search">Search</a>
        <a href="/add">Add New Item</a>
    </div>
    
    {{range $index, $magazine := .Magazines}}
    <div class="item">
        <div class="title">{{add $index 1}}. {{$magazine.Title}}</div>
        <div class="isbn">ISBN: {{$magazine.ISBN}}</div>
        <div class="authors">Authors: {{formatAuthors $magazine.Authors}}</div>
        <div class="published">Published: {{$magazine.PublishedAt}}</div>
    </div>
    {{end}}
</body>
</html>
`
}

func (ws *WebServer) getAllTemplate() string {
	return `
<!DOCTYPE html>
<html>
<head>
    <title>{{.Title}}</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 40px; }
        .nav { margin-bottom: 30px; }
        .nav a { margin-right: 15px; text-decoration: none; color: #0066cc; }
        .nav a:hover { text-decoration: underline; }
        .item { margin-bottom: 20px; padding: 15px; border: 1px solid #ddd; }
        .title { font-weight: bold; font-size: 18px; color: #333; }
        .isbn { color: #666; }
        .authors { color: #0066cc; }
        .description { margin-top: 10px; line-height: 1.5; }
        .published { color: #888; }
        .section { margin-bottom: 40px; }
        .section h2 { color: #333; border-bottom: 2px solid #0066cc; padding-bottom: 10px; }
    </style>
</head>
<body>
    <h1>{{.Title}}</h1>
    <div class="nav">
        <a href="/">Home</a>
        <a href="/books">All Books</a>
        <a href="/magazines">All Magazines</a>
        <a href="/all">All Items</a>
        <a href="/sorted">Sorted Items</a>
        <a href="/search">Search</a>
        <a href="/add">Add New Item</a>
    </div>
    
    <div class="section">
        <h2>BOOKS</h2>
        {{range $index, $book := .Books}}
        <div class="item">
            <div class="title">{{add $index 1}}. {{$book.Title}}</div>
            <div class="isbn">ISBN: {{$book.ISBN}}</div>
            <div class="authors">Authors: {{formatAuthors $book.Authors}}</div>
            <div class="description">{{$book.Description}}</div>
        </div>
        {{end}}
    </div>
    
    <div class="section">
        <h2>MAGAZINES</h2>
        {{range $index, $magazine := .Magazines}}
        <div class="item">
            <div class="title">{{add $index 1}}. {{$magazine.Title}}</div>
            <div class="isbn">ISBN: {{$magazine.ISBN}}</div>
            <div class="authors">Authors: {{formatAuthors $magazine.Authors}}</div>
            <div class="published">Published: {{$magazine.PublishedAt}}</div>
        </div>
        {{end}}
    </div>
</body>
</html>
`
}

func (ws *WebServer) getSortedTemplate() string {
	return `
<!DOCTYPE html>
<html>
<head>
    <title>{{.Title}}</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 40px; }
        .nav { margin-bottom: 30px; }
        .nav a { margin-right: 15px; text-decoration: none; color: #0066cc; }
        .nav a:hover { text-decoration: underline; }
        .item { margin-bottom: 20px; padding: 15px; border: 1px solid #ddd; }
        .title { font-weight: bold; font-size: 18px; color: #333; }
        .isbn { color: #666; }
        .authors { color: #0066cc; }
        .description { margin-top: 10px; line-height: 1.5; }
        .published { color: #888; }
        .type { background: #f0f0f0; padding: 3px 8px; border-radius: 3px; font-size: 12px; }
        .controls { margin-bottom: 20px; }
        .controls a { margin-right: 10px; padding: 8px 16px; background: #0066cc; color: white; text-decoration: none; border-radius: 4px; }
        .controls a:hover { background: #0052a3; }
    </style>
</head>
<body>
    <h1>{{.Title}}</h1>
    <div class="nav">
        <a href="/">Home</a>
        <a href="/books">All Books</a>
        <a href="/magazines">All Magazines</a>
        <a href="/all">All Items</a>
        <a href="/sorted">Sorted Items</a>
        <a href="/search">Search</a>
        <a href="/add">Add New Item</a>
    </div>
    
    <div class="controls">
        <a href="/sorted?order=asc">Sort A-Z</a>
        <a href="/sorted?order=desc">Sort Z-A</a>
    </div>
    
    {{range $index, $item := .Items}}
    <div class="item">
        <div class="title">{{add $index 1}}. {{$item.Title}} <span class="type">{{$item.Type}}</span></div>
        <div class="isbn">ISBN: {{$item.ISBN}}</div>
        <div class="authors">Authors: {{formatAuthors $item.Authors}}</div>
        {{if eq $item.Type "Book"}}
            <div class="description">{{$item.Description}}</div>
        {{else}}
            <div class="published">Published: {{$item.PublishedAt}}</div>
        {{end}}
    </div>
    {{end}}
</body>
</html>
`
}

func (ws *WebServer) getSearchTemplate() string {
	return `
<!DOCTYPE html>
<html>
<head>
    <title>{{.Title}}</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 40px; }
        .nav { margin-bottom: 30px; }
        .nav a { margin-right: 15px; text-decoration: none; color: #0066cc; }
        .nav a:hover { text-decoration: underline; }
        .item { margin-bottom: 20px; padding: 15px; border: 1px solid #ddd; }
        .title { font-weight: bold; font-size: 18px; color: #333; }
        .isbn { color: #666; }
        .authors { color: #0066cc; }
        .description { margin-top: 10px; line-height: 1.5; }
        .published { color: #888; }
        .search-form { margin-bottom: 30px; padding: 20px; background: #f9f9f9; border-radius: 5px; }
        .search-form input, .search-form select { margin: 5px; padding: 8px; }
        .search-form button { padding: 8px 16px; background: #0066cc; color: white; border: none; border-radius: 4px; cursor: pointer; }
        .search-form button:hover { background: #0052a3; }
        .error { color: red; margin-bottom: 20px; }
        .message { color: green; margin-bottom: 20px; }
        .section { margin-bottom: 40px; }
        .section h2 { color: #333; border-bottom: 2px solid #0066cc; padding-bottom: 10px; }
    </style>
</head>
<body>
    <h1>{{.Title}}</h1>
    <div class="nav">
        <a href="/">Home</a>
        <a href="/books">All Books</a>
        <a href="/magazines">All Magazines</a>
        <a href="/all">All Items</a>
        <a href="/sorted">Sorted Items</a>
        <a href="/search">Search</a>
        <a href="/add">Add New Item</a>
    </div>
    
    <div class="search-form">
        <form method="POST">
            <select name="search_type">
                <option value="isbn">Search by ISBN</option>
                <option value="email">Search by Author Email</option>
            </select>
            <input type="text" name="search_value" placeholder="Enter search value..." required>
            <button type="submit">Search</button>
        </form>
    </div>
    
    {{if .Error}}
        <div class="error">{{.Error}}</div>
    {{end}}
    
    {{if .Books}}
    <div class="section">
        <h2>BOOKS FOUND</h2>
        {{range $index, $book := .Books}}
        <div class="item">
            <div class="title">{{add $index 1}}. {{$book.Title}}</div>
            <div class="isbn">ISBN: {{$book.ISBN}}</div>
            <div class="authors">Authors: {{formatAuthors $book.Authors}}</div>
            <div class="description">{{$book.Description}}</div>
        </div>
        {{end}}
    </div>
    {{end}}
    
    {{if .Magazines}}
    <div class="section">
        <h2>MAGAZINES FOUND</h2>
        {{range $index, $magazine := .Magazines}}
        <div class="item">
            <div class="title">{{add $index 1}}. {{$magazine.Title}}</div>
            <div class="isbn">ISBN: {{$magazine.ISBN}}</div>
            <div class="authors">Authors: {{formatAuthors $magazine.Authors}}</div>
            <div class="published">Published: {{$magazine.PublishedAt}}</div>
        </div>
        {{end}}
    </div>
    {{end}}
</body>
</html>
`
}

func (ws *WebServer) getAddTemplate() string {
	return `
<!DOCTYPE html>
<html>
<head>
    <title>{{.Title}}</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 40px; }
        .nav { margin-bottom: 30px; }
        .nav a { margin-right: 15px; text-decoration: none; color: #0066cc; }
        .nav a:hover { text-decoration: underline; }
        .add-options { margin-bottom: 30px; }
        .add-options a { display: inline-block; margin-right: 20px; padding: 15px 30px; background: #0066cc; color: white; text-decoration: none; border-radius: 5px; }
        .add-options a:hover { background: #0052a3; }
    </style>
</head>
<body>
    <h1>{{.Title}}</h1>
    <div class="nav">
        <a href="/">Home</a>
        <a href="/books">All Books</a>
        <a href="/magazines">All Magazines</a>
        <a href="/all">All Items</a>
        <a href="/sorted">Sorted Items</a>
        <a href="/search">Search</a>
        <a href="/add">Add New Item</a>
    </div>
    
    <h2>What would you like to add?</h2>
    <div class="add-options">
        <a href="/add-book">Add Book</a>
        <a href="/add-magazine">Add Magazine</a>
    </div>
</body>
</html>
`
}

func (ws *WebServer) getAddBookTemplate() string {
	return `
<!DOCTYPE html>
<html>
<head>
    <title>{{.Title}}</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 40px; }
        .nav { margin-bottom: 30px; }
        .nav a { margin-right: 15px; text-decoration: none; color: #0066cc; }
        .nav a:hover { text-decoration: underline; }
        .form { max-width: 600px; }
        .form label { display: block; margin-top: 15px; font-weight: bold; }
        .form input, .form textarea { width: 100%; padding: 8px; margin-top: 5px; border: 1px solid #ddd; border-radius: 4px; }
        .form textarea { height: 100px; resize: vertical; }
        .form button { margin-top: 20px; padding: 10px 20px; background: #0066cc; color: white; border: none; border-radius: 4px; cursor: pointer; }
        .form button:hover { background: #0052a3; }
        .error { color: red; margin-bottom: 20px; }
        .message { color: green; margin-bottom: 20px; }
        .help { color: #666; font-size: 14px; margin-top: 5px; }
    </style>
</head>
<body>
    <h1>{{.Title}}</h1>
    <div class="nav">
        <a href="/">Home</a>
        <a href="/books">All Books</a>
        <a href="/magazines">All Magazines</a>
        <a href="/all">All Items</a>
        <a href="/sorted">Sorted Items</a>
        <a href="/search">Search</a>
        <a href="/add">Add New Item</a>
    </div>
    
    {{if .Error}}
        <div class="error">{{.Error}}</div>
    {{end}}
    
    {{if .Message}}
        <div class="message">{{.Message}}</div>
    {{end}}
    
    <div class="form">
        <form method="POST">
            <label for="title">Book Title:</label>
            <input type="text" id="title" name="title" required>
            
            <label for="isbn">ISBN:</label>
            <input type="text" id="isbn" name="isbn" required>
            
            <label for="description">Description:</label>
            <textarea id="description" name="description" required></textarea>
            
            <label for="author_emails">Author Emails:</label>
            <input type="text" id="author_emails" name="author_emails" required>
            <div class="help">Enter email addresses separated by commas (e.g., john@example.com, jane@example.com)</div>
            
            <label for="author_names">Author Names:</label>
            <input type="text" id="author_names" name="author_names" required>
            <div class="help">Enter full names separated by commas (e.g., John Doe, Jane Smith)</div>
            
            <button type="submit">Add Book</button>
        </form>
    </div>
</body>
</html>
`
}

func (ws *WebServer) getAddMagazineTemplate() string {
	return `
<!DOCTYPE html>
<html>
<head>
    <title>{{.Title}}</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 40px; }
        .nav { margin-bottom: 30px; }
        .nav a { margin-right: 15px; text-decoration: none; color: #0066cc; }
        .nav a:hover { text-decoration: underline; }
        .form { max-width: 600px; }
        .form label { display: block; margin-top: 15px; font-weight: bold; }
        .form input, .form textarea { width: 100%; padding: 8px; margin-top: 5px; border: 1px solid #ddd; border-radius: 4px; }
        .form button { margin-top: 20px; padding: 10px 20px; background: #0066cc; color: white; border: none; border-radius: 4px; cursor: pointer; }
        .form button:hover { background: #0052a3; }
        .error { color: red; margin-bottom: 20px; }
        .message { color: green; margin-bottom: 20px; }
        .help { color: #666; font-size: 14px; margin-top: 5px; }
    </style>
</head>
<body>
    <h1>{{.Title}}</h1>
    <div class="nav">
        <a href="/">Home</a>
        <a href="/books">All Books</a>
        <a href="/magazines">All Magazines</a>
        <a href="/all">All Items</a>
        <a href="/sorted">Sorted Items</a>
        <a href="/search">Search</a>
        <a href="/add">Add New Item</a>
    </div>
    
    {{if .Error}}
        <div class="error">{{.Error}}</div>
    {{end}}
    
    {{if .Message}}
        <div class="message">{{.Message}}</div>
    {{end}}
    
    <div class="form">
        <form method="POST">
            <label for="title">Magazine Title:</label>
            <input type="text" id="title" name="title" required>
            
            <label for="isbn">ISBN:</label>
            <input type="text" id="isbn" name="isbn" required>
            
            <label for="published_at">Published Date:</label>
            <input type="text" id="published_at" name="published_at" placeholder="DD.MM.YYYY" required>
            
            <label for="author_emails">Author Emails:</label>
            <input type="text" id="author_emails" name="author_emails" required>
            <div class="help">Enter email addresses separated by commas (e.g., john@example.com, jane@example.com)</div>
            
            <label for="author_names">Author Names:</label>
            <input type="text" id="author_names" name="author_names" required>
            <div class="help">Enter full names separated by commas (e.g., John Doe, Jane Smith)</div>
            
            <button type="submit">Add Magazine</button>
        </form>
    </div>
</body>
</html>
`
}

