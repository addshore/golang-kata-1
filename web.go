package main

import (
	"html/template"
	"net/http"
	"sort"
	"strings"
	"time"
)

type WebServer struct {
	library *Library
}

type SearchResult struct {
	Type        string
	Title       string
	ISBN        string
	Authors     string
	Description string
	Published   string
}

func startWebServer(library *Library) {
	server := &WebServer{library: library}

	http.HandleFunc("/", server.handleHome)
	http.HandleFunc("/books", server.handleBooks)
	http.HandleFunc("/magazines", server.handleMagazines)
	http.HandleFunc("/all", server.handleAll)
	http.HandleFunc("/sorted", server.handleSorted)
	http.HandleFunc("/search", server.handleSearch)
	http.HandleFunc("/add", server.handleAdd)
	http.HandleFunc("/add-book", server.handleAddBook)
	http.HandleFunc("/add-magazine", server.handleAddMagazine)

	println("Web server starting on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func (ws *WebServer) handleHome(w http.ResponseWriter, r *http.Request) {
	tmpl := `
<!DOCTYPE html>
<html>
<head>
    <title>Library System</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 20px; }
        .nav { margin-bottom: 20px; }
        .nav a { margin-right: 15px; text-decoration: none; color: #0066cc; }
        .nav a:hover { text-decoration: underline; }
        .stats { background: #f5f5f5; padding: 15px; border-radius: 5px; margin: 20px 0; }
    </style>
</head>
<body>
    <h1>Library Management System</h1>
    
    <div class="nav">
        <a href="/books">View Books</a>
        <a href="/magazines">View Magazines</a>
        <a href="/all">View All Items</a>
        <a href="/sorted">View All (Sorted by Title)</a>
        <a href="/search">Search</a>
        <a href="/add">Add New Item</a>
    </div>
    
    <div class="stats">
        <h3>Library Statistics</h3>
        <p>Total Books: {{.BookCount}}</p>
        <p>Total Magazines: {{.MagazineCount}}</p>
        <p>Total Authors: {{.AuthorCount}}</p>
        <p>Total Items: {{.TotalCount}}</p>
    </div>
    
    <h3>Welcome to the Library System</h3>
    <p>Use the navigation links above to browse books, magazines, search for items, or add new content to the library.</p>
</body>
</html>`

	t := template.Must(template.New("home").Parse(tmpl))

	data := map[string]interface{}{
		"BookCount":     len(ws.library.Books),
		"MagazineCount": len(ws.library.Magazines),
		"AuthorCount":   len(ws.library.Authors),
		"TotalCount":    len(ws.library.Books) + len(ws.library.Magazines),
	}

	t.Execute(w, data)
}

func (ws *WebServer) handleBooks(w http.ResponseWriter, r *http.Request) {
	tmpl := `
<!DOCTYPE html>
<html>
<head>
    <title>Books - Library System</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 20px; }
        .nav { margin-bottom: 20px; }
        .nav a { margin-right: 15px; text-decoration: none; color: #0066cc; }
        .item { border: 1px solid #ddd; margin: 10px 0; padding: 15px; border-radius: 5px; }
        .title { font-weight: bold; font-size: 1.1em; color: #333; }
        .isbn { color: #666; }
        .authors { color: #0066cc; }
        .description { margin-top: 10px; }
    </style>
</head>
<body>
    <div class="nav">
        <a href="/">Home</a>
        <a href="/magazines">Magazines</a>
        <a href="/all">All Items</a>
        <a href="/sorted">Sorted</a>
        <a href="/search">Search</a>
        <a href="/add">Add New</a>
    </div>
    
    <h1>Books ({{len .}})</h1>
    
    {{range $index, $book := .}}
    <div class="item">
        <div class="title">{{$book.Title}}</div>
        <div class="isbn">ISBN: {{$book.ISBN}}</div>
        <div class="authors">Authors: {{$book.AuthorsString}}</div>
        <div class="description">{{$book.Description}}</div>
    </div>
    {{end}}
</body>
</html>`

	t := template.Must(template.New("books").Parse(tmpl))

	type BookDisplay struct {
		Title         string
		ISBN          string
		AuthorsString string
		Description   string
	}

	var books []BookDisplay
	for _, book := range ws.library.Books {
		books = append(books, BookDisplay{
			Title:         book.Title,
			ISBN:          book.ISBN,
			AuthorsString: formatAuthors(book.Authors),
			Description:   book.Description,
		})
	}

	t.Execute(w, books)
}

func (ws *WebServer) handleMagazines(w http.ResponseWriter, r *http.Request) {
	tmpl := `
<!DOCTYPE html>
<html>
<head>
    <title>Magazines - Library System</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 20px; }
        .nav { margin-bottom: 20px; }
        .nav a { margin-right: 15px; text-decoration: none; color: #0066cc; }
        .item { border: 1px solid #ddd; margin: 10px 0; padding: 15px; border-radius: 5px; }
        .title { font-weight: bold; font-size: 1.1em; color: #333; }
        .isbn { color: #666; }
        .authors { color: #0066cc; }
        .published { color: #888; }
    </style>
</head>
<body>
    <div class="nav">
        <a href="/">Home</a>
        <a href="/books">Books</a>
        <a href="/all">All Items</a>
        <a href="/sorted">Sorted</a>
        <a href="/search">Search</a>
        <a href="/add">Add New</a>
    </div>
    
    <h1>Magazines ({{len .}})</h1>
    
    {{range $index, $magazine := .}}
    <div class="item">
        <div class="title">{{$magazine.Title}}</div>
        <div class="isbn">ISBN: {{$magazine.ISBN}}</div>
        <div class="authors">Authors: {{$magazine.AuthorsString}}</div>
        {{if $magazine.Published}}
        <div class="published">Published: {{$magazine.Published}}</div>
        {{end}}
    </div>
    {{end}}
</body>
</html>`

	t := template.Must(template.New("magazines").Parse(tmpl))

	type MagazineDisplay struct {
		Title         string
		ISBN          string
		AuthorsString string
		Published     string
	}

	var magazines []MagazineDisplay
	for _, magazine := range ws.library.Magazines {
		published := ""
		if !magazine.PublishedAt.IsZero() {
			published = magazine.PublishedAt.Format("02.01.2006")
		}

		magazines = append(magazines, MagazineDisplay{
			Title:         magazine.Title,
			ISBN:          magazine.ISBN,
			AuthorsString: formatAuthors(magazine.Authors),
			Published:     published,
		})
	}

	t.Execute(w, magazines)
}

func (ws *WebServer) handleAll(w http.ResponseWriter, r *http.Request) {
	ws.renderAllItems(w, false)
}

func (ws *WebServer) handleSorted(w http.ResponseWriter, r *http.Request) {
	ws.renderAllItems(w, true)
}

func (ws *WebServer) renderAllItems(w http.ResponseWriter, sorted bool) {
	tmpl := `
<!DOCTYPE html>
<html>
<head>
    <title>{{if .Sorted}}All Items (Sorted by Title){{else}}All Items{{end}} - Library System</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 20px; }
        .nav { margin-bottom: 20px; }
        .nav a { margin-right: 15px; text-decoration: none; color: #0066cc; }
        .item { border: 1px solid #ddd; margin: 10px 0; padding: 15px; border-radius: 5px; }
        .title { font-weight: bold; font-size: 1.1em; color: #333; }
        .type { background: #e6f3ff; color: #0066cc; padding: 2px 8px; border-radius: 3px; font-size: 0.8em; margin-right: 10px; }
        .isbn { color: #666; }
        .authors { color: #0066cc; }
        .description { margin-top: 10px; }
        .published { color: #888; }
    </style>
</head>
<body>
    <div class="nav">
        <a href="/">Home</a>
        <a href="/books">Books</a>
        <a href="/magazines">Magazines</a>
        {{if .Sorted}}
        <a href="/all">All Items</a>
        {{else}}
        <a href="/sorted">Sorted</a>
        {{end}}
        <a href="/search">Search</a>
        <a href="/add">Add New</a>
    </div>
    
    <h1>{{if .Sorted}}All Items (Sorted by Title){{else}}All Items{{end}} ({{len .Items}})</h1>
    
    {{range $index, $item := .Items}}
    <div class="item">
        <span class="type">{{$item.Type}}</span>
        <span class="title">{{$item.Title}}</span>
        <div class="isbn">ISBN: {{$item.ISBN}}</div>
        <div class="authors">Authors: {{$item.Authors}}</div>
        {{if $item.Description}}
        <div class="description">{{$item.Description}}</div>
        {{end}}
        {{if $item.Published}}
        <div class="published">Published: {{$item.Published}}</div>
        {{end}}
    </div>
    {{end}}
</body>
</html>`

	t := template.Must(template.New("all").Parse(tmpl))

	var items []SearchResult

	// Add books
	for _, book := range ws.library.Books {
		items = append(items, SearchResult{
			Type:        "BOOK",
			Title:       book.Title,
			ISBN:        book.ISBN,
			Authors:     formatAuthors(book.Authors),
			Description: book.Description,
			Published:   "",
		})
	}

	// Add magazines
	for _, magazine := range ws.library.Magazines {
		published := ""
		if !magazine.PublishedAt.IsZero() {
			published = magazine.PublishedAt.Format("02.01.2006")
		}

		items = append(items, SearchResult{
			Type:        "MAGAZINE",
			Title:       magazine.Title,
			ISBN:        magazine.ISBN,
			Authors:     formatAuthors(magazine.Authors),
			Description: "",
			Published:   published,
		})
	}

	if sorted {
		sort.Slice(items, func(i, j int) bool {
			return strings.ToLower(items[i].Title) < strings.ToLower(items[j].Title)
		})
	}

	data := map[string]interface{}{
		"Items":  items,
		"Sorted": sorted,
	}

	t.Execute(w, data)
}

func (ws *WebServer) handleSearch(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		ws.handleSearchResults(w, r)
		return
	}

	tmpl := `
<!DOCTYPE html>
<html>
<head>
    <title>Search - Library System</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 20px; }
        .nav { margin-bottom: 20px; }
        .nav a { margin-right: 15px; text-decoration: none; color: #0066cc; }
        .form { background: #f9f9f9; padding: 20px; border-radius: 5px; margin: 20px 0; }
        .form input, .form select { padding: 8px; margin: 5px 0; width: 300px; }
        .form button { padding: 10px 20px; background: #0066cc; color: white; border: none; border-radius: 3px; cursor: pointer; }
    </style>
</head>
<body>
    <div class="nav">
        <a href="/">Home</a>
        <a href="/books">Books</a>
        <a href="/magazines">Magazines</a>
        <a href="/all">All Items</a>
        <a href="/sorted">Sorted</a>
        <a href="/add">Add New</a>
    </div>
    
    <h1>Search Library</h1>
    
    <div class="form">
        <form method="POST">
            <h3>Search by:</h3>
            <select name="searchType">
                <option value="isbn">ISBN</option>
                <option value="author">Author Email</option>
            </select><br>
            <input type="text" name="query" placeholder="Enter search term..." required><br>
            <button type="submit">Search</button>
        </form>
    </div>
</body>
</html>`

	t := template.Must(template.New("search").Parse(tmpl))
	t.Execute(w, nil)
}

func (ws *WebServer) handleSearchResults(w http.ResponseWriter, r *http.Request) {
	searchType := r.FormValue("searchType")
	query := r.FormValue("query")

	var results []SearchResult

	if searchType == "isbn" {
		// Search by ISBN
		for _, book := range ws.library.Books {
			if strings.EqualFold(book.ISBN, query) {
				results = append(results, SearchResult{
					Type:        "BOOK",
					Title:       book.Title,
					ISBN:        book.ISBN,
					Authors:     formatAuthors(book.Authors),
					Description: book.Description,
					Published:   "",
				})
			}
		}

		for _, magazine := range ws.library.Magazines {
			if strings.EqualFold(magazine.ISBN, query) {
				published := ""
				if !magazine.PublishedAt.IsZero() {
					published = magazine.PublishedAt.Format("02.01.2006")
				}

				results = append(results, SearchResult{
					Type:        "MAGAZINE",
					Title:       magazine.Title,
					ISBN:        magazine.ISBN,
					Authors:     formatAuthors(magazine.Authors),
					Description: "",
					Published:   published,
				})
			}
		}
	} else if searchType == "author" {
		// Search by author email
		for _, book := range ws.library.Books {
			for _, author := range book.Authors {
				if strings.EqualFold(author.Email, query) {
					results = append(results, SearchResult{
						Type:        "BOOK",
						Title:       book.Title,
						ISBN:        book.ISBN,
						Authors:     formatAuthors(book.Authors),
						Description: book.Description,
						Published:   "",
					})
					break
				}
			}
		}

		for _, magazine := range ws.library.Magazines {
			for _, author := range magazine.Authors {
				if strings.EqualFold(author.Email, query) {
					published := ""
					if !magazine.PublishedAt.IsZero() {
						published = magazine.PublishedAt.Format("02.01.2006")
					}

					results = append(results, SearchResult{
						Type:        "MAGAZINE",
						Title:       magazine.Title,
						ISBN:        magazine.ISBN,
						Authors:     formatAuthors(magazine.Authors),
						Description: "",
						Published:   published,
					})
					break
				}
			}
		}
	}

	tmpl := `
<!DOCTYPE html>
<html>
<head>
    <title>Search Results - Library System</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 20px; }
        .nav { margin-bottom: 20px; }
        .nav a { margin-right: 15px; text-decoration: none; color: #0066cc; }
        .item { border: 1px solid #ddd; margin: 10px 0; padding: 15px; border-radius: 5px; }
        .title { font-weight: bold; font-size: 1.1em; color: #333; }
        .type { background: #e6f3ff; color: #0066cc; padding: 2px 8px; border-radius: 3px; font-size: 0.8em; margin-right: 10px; }
        .isbn { color: #666; }
        .authors { color: #0066cc; }
        .description { margin-top: 10px; }
        .published { color: #888; }
        .no-results { color: #888; font-style: italic; }
    </style>
</head>
<body>
    <div class="nav">
        <a href="/">Home</a>
        <a href="/books">Books</a>
        <a href="/magazines">Magazines</a>
        <a href="/all">All Items</a>
        <a href="/sorted">Sorted</a>
        <a href="/search">Search</a>
        <a href="/add">Add New</a>
    </div>
    
    <h1>Search Results</h1>
    <p>Searching for "{{.Query}}" by {{.SearchType}}</p>
    
    {{if .Results}}
        <p>Found {{len .Results}} result(s):</p>
        {{range $index, $item := .Results}}
        <div class="item">
            <span class="type">{{$item.Type}}</span>
            <span class="title">{{$item.Title}}</span>
            <div class="isbn">ISBN: {{$item.ISBN}}</div>
            <div class="authors">Authors: {{$item.Authors}}</div>
            {{if $item.Description}}
            <div class="description">{{$item.Description}}</div>
            {{end}}
            {{if $item.Published}}
            <div class="published">Published: {{$item.Published}}</div>
            {{end}}
        </div>
        {{end}}
    {{else}}
        <p class="no-results">No results found for your search.</p>
    {{end}}
    
    <p><a href="/search">← Back to Search</a></p>
</body>
</html>`

	t := template.Must(template.New("results").Parse(tmpl))

	data := map[string]interface{}{
		"Results":    results,
		"Query":      query,
		"SearchType": searchType,
	}

	t.Execute(w, data)
}

func (ws *WebServer) handleAdd(w http.ResponseWriter, r *http.Request) {
	tmpl := `
<!DOCTYPE html>
<html>
<head>
    <title>Add New Item - Library System</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 20px; }
        .nav { margin-bottom: 20px; }
        .nav a { margin-right: 15px; text-decoration: none; color: #0066cc; }
        .option { background: #f9f9f9; padding: 20px; margin: 10px 0; border-radius: 5px; border: 1px solid #ddd; }
        .option h3 { margin-top: 0; }
        .option a { background: #0066cc; color: white; padding: 10px 20px; text-decoration: none; border-radius: 3px; }
    </style>
</head>
<body>
    <div class="nav">
        <a href="/">Home</a>
        <a href="/books">Books</a>
        <a href="/magazines">Magazines</a>
        <a href="/all">All Items</a>
        <a href="/sorted">Sorted</a>
        <a href="/search">Search</a>
    </div>
    
    <h1>Add New Item to Library</h1>
    
    <div class="option">
        <h3>Add New Book</h3>
        <p>Add a new book with title, ISBN, description, and authors.</p>
        <a href="/add-book">Add Book</a>
    </div>
    
    <div class="option">
        <h3>Add New Magazine</h3>
        <p>Add a new magazine with title, ISBN, publication date, and authors.</p>
        <a href="/add-magazine">Add Magazine</a>
    </div>
</body>
</html>`

	t := template.Must(template.New("add").Parse(tmpl))
	t.Execute(w, nil)
}

func (ws *WebServer) handleAddBook(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		ws.processAddBook(w, r)
		return
	}

	tmpl := `
<!DOCTYPE html>
<html>
<head>
    <title>Add New Book - Library System</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 20px; }
        .nav { margin-bottom: 20px; }
        .nav a { margin-right: 15px; text-decoration: none; color: #0066cc; }
        .form { background: #f9f9f9; padding: 20px; border-radius: 5px; margin: 20px 0; }
        .form input, .form textarea { padding: 8px; margin: 5px 0; width: 400px; }
        .form textarea { height: 100px; }
        .form button { padding: 10px 20px; background: #0066cc; color: white; border: none; border-radius: 3px; cursor: pointer; margin: 5px; }
        .author-section { margin: 20px 0; padding: 15px; background: #fff; border: 1px solid #ddd; border-radius: 5px; }
        .author-entry { margin: 10px 0; padding: 10px; background: #f5f5f5; border-radius: 3px; }
        .author-info { color: #0066cc; font-weight: bold; }
        .author-form { background: #fff; padding: 10px; border: 1px solid #ccc; border-radius: 3px; margin-top: 5px; }
        .remove-author { background: #cc0000; padding: 5px 10px; }
    </style>
    <script>
        let authorCount = 0;
        const authors = {{.AuthorsJSON}};
        
        function addAuthor() {
            authorCount++;
            const container = document.getElementById('authors-container');
            const div = document.createElement('div');
            div.className = 'author-entry';
            div.id = 'author-' + authorCount;
            
            div.innerHTML = 
                '<label>Author Email:</label><br>' +
                '<input type="email" id="email-' + authorCount + '" onblur="checkAuthor(' + authorCount + ')" required style="width: 300px;"><br>' +
                '<div id="author-info-' + authorCount + '"></div>' +
                '<button type="button" onclick="removeAuthor(' + authorCount + ')" class="remove-author">Remove Author</button>';
            
            container.appendChild(div);
        }
        
        function removeAuthor(id) {
            const element = document.getElementById('author-' + id);
            element.remove();
        }
        
        function checkAuthor(id) {
            const email = document.getElementById('email-' + id).value;
            const infoDiv = document.getElementById('author-info-' + id);
            
            if (!email) {
                infoDiv.innerHTML = '';
                return;
            }
            
            if (authors[email]) {
                infoDiv.innerHTML = 
                    '<div class="author-info">✓ Existing author: ' + authors[email].firstName + ' ' + authors[email].lastName + '</div>' +
                    '<input type="hidden" name="authors[]" value="' + email + '">';
            } else {
                infoDiv.innerHTML = 
                    '<div style="color: #cc6600;">⚠ New author - please provide details:</div>' +
                    '<div class="author-form">' +
                    '<label>First Name:</label><br>' +
                    '<input type="text" name="author_first_' + id + '" required style="width: 180px; margin-right: 10px;">' +
                    '<label>Last Name:</label><br>' +
                    '<input type="text" name="author_last_' + id + '" required style="width: 180px;">' +
                    '<input type="hidden" name="authors[]" value="' + email + '">' +
                    '</div>';
            }
        }
        
        window.onload = function() {
            addAuthor(); // Add first author field
        }
    </script>
</head>
<body>
    <div class="nav">
        <a href="/">Home</a>
        <a href="/books">Books</a>
        <a href="/magazines">Magazines</a>
        <a href="/all">All Items</a>
        <a href="/sorted">Sorted</a>
        <a href="/search">Search</a>
        <a href="/add">Add New</a>
    </div>
    
    <h1>Add New Book</h1>
    
    <div class="form">
        <form method="POST">
            <label>Title:</label><br>
            <input type="text" name="title" required><br>
            
            <label>ISBN:</label><br>
            <input type="text" name="isbn" required><br>
            
            <label>Description:</label><br>
            <textarea name="description"></textarea><br>
            
            <div class="author-section">
                <h3>Authors</h3>
                <div id="authors-container"></div>
                <button type="button" onclick="addAuthor()">Add Another Author</button>
            </div>
            
            <button type="submit">Add Book</button>
        </form>
    </div>
</body>
</html>`

	// Prepare authors data for JavaScript
	authorsJSON := "{"
	first := true
	for email, author := range ws.library.Authors {
		if !first {
			authorsJSON += ","
		}
		authorsJSON += "\"" + email + "\":{\"firstName\":\"" + author.FirstName + "\",\"lastName\":\"" + author.LastName + "\"}"
		first = false
	}
	authorsJSON += "}"

	t := template.Must(template.New("add-book").Parse(tmpl))
	data := map[string]interface{}{
		"AuthorsJSON": template.JS(authorsJSON),
	}
	t.Execute(w, data)
}

func (ws *WebServer) handleAddMagazine(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		ws.processAddMagazine(w, r)
		return
	}

	tmpl := `
<!DOCTYPE html>
<html>
<head>
    <title>Add New Magazine - Library System</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 20px; }
        .nav { margin-bottom: 20px; }
        .nav a { margin-right: 15px; text-decoration: none; color: #0066cc; }
        .form { background: #f9f9f9; padding: 20px; border-radius: 5px; margin: 20px 0; }
        .form input { padding: 8px; margin: 5px 0; width: 400px; }
        .form button { padding: 10px 20px; background: #0066cc; color: white; border: none; border-radius: 3px; cursor: pointer; margin: 5px; }
        .author-section { margin: 20px 0; padding: 15px; background: #fff; border: 1px solid #ddd; border-radius: 5px; }
        .author-entry { margin: 10px 0; padding: 10px; background: #f5f5f5; border-radius: 3px; }
        .author-info { color: #0066cc; font-weight: bold; }
        .author-form { background: #fff; padding: 10px; border: 1px solid #ccc; border-radius: 3px; margin-top: 5px; }
        .remove-author { background: #cc0000; padding: 5px 10px; }
    </style>
    <script>
        let authorCount = 0;
        const authors = {{.AuthorsJSON}};
        
        function addAuthor() {
            authorCount++;
            const container = document.getElementById('authors-container');
            const div = document.createElement('div');
            div.className = 'author-entry';
            div.id = 'author-' + authorCount;
            
            div.innerHTML = 
                '<label>Author Email:</label><br>' +
                '<input type="email" id="email-' + authorCount + '" onblur="checkAuthor(' + authorCount + ')" required style="width: 300px;"><br>' +
                '<div id="author-info-' + authorCount + '"></div>' +
                '<button type="button" onclick="removeAuthor(' + authorCount + ')" class="remove-author">Remove Author</button>';
            
            container.appendChild(div);
        }
        
        function removeAuthor(id) {
            const element = document.getElementById('author-' + id);
            element.remove();
        }
        
        function checkAuthor(id) {
            const email = document.getElementById('email-' + id).value;
            const infoDiv = document.getElementById('author-info-' + id);
            
            if (!email) {
                infoDiv.innerHTML = '';
                return;
            }
            
            if (authors[email]) {
                infoDiv.innerHTML = 
                    '<div class="author-info">✓ Existing author: ' + authors[email].firstName + ' ' + authors[email].lastName + '</div>' +
                    '<input type="hidden" name="authors[]" value="' + email + '">';
            } else {
                infoDiv.innerHTML = 
                    '<div style="color: #cc6600;">⚠ New author - please provide details:</div>' +
                    '<div class="author-form">' +
                    '<label>First Name:</label><br>' +
                    '<input type="text" name="author_first_' + id + '" required style="width: 180px; margin-right: 10px;">' +
                    '<label>Last Name:</label><br>' +
                    '<input type="text" name="author_last_' + id + '" required style="width: 180px;">' +
                    '<input type="hidden" name="authors[]" value="' + email + '">' +
                    '</div>';
            }
        }
        
        window.onload = function() {
            addAuthor(); // Add first author field
        }
    </script>
</head>
<body>
    <div class="nav">
        <a href="/">Home</a>
        <a href="/books">Books</a>
        <a href="/magazines">Magazines</a>
        <a href="/all">All Items</a>
        <a href="/sorted">Sorted</a>
        <a href="/search">Search</a>
        <a href="/add">Add New</a>
    </div>
    
    <h1>Add New Magazine</h1>
    
    <div class="form">
        <form method="POST">
            <label>Title:</label><br>
            <input type="text" name="title" required><br>
            
            <label>ISBN:</label><br>
            <input type="text" name="isbn" required><br>
            
            <label>Publication Date (DD.MM.YYYY):</label><br>
            <input type="text" name="published" placeholder="01.01.2024" required><br>
            
            <div class="author-section">
                <h3>Authors</h3>
                <div id="authors-container"></div>
                <button type="button" onclick="addAuthor()">Add Another Author</button>
            </div>
            
            <button type="submit">Add Magazine</button>
        </form>
    </div>
</body>
</html>`

	// Prepare authors data for JavaScript
	authorsJSON := "{"
	first := true
	for email, author := range ws.library.Authors {
		if !first {
			authorsJSON += ","
		}
		authorsJSON += "\"" + email + "\":{\"firstName\":\"" + author.FirstName + "\",\"lastName\":\"" + author.LastName + "\"}"
		first = false
	}
	authorsJSON += "}"

	t := template.Must(template.New("add-magazine").Parse(tmpl))
	data := map[string]interface{}{
		"AuthorsJSON": template.JS(authorsJSON),
	}
	t.Execute(w, data)
}

func (ws *WebServer) processAddBook(w http.ResponseWriter, r *http.Request) {
	title := strings.TrimSpace(r.FormValue("title"))
	isbn := strings.TrimSpace(r.FormValue("isbn"))
	description := strings.TrimSpace(r.FormValue("description"))

	// Validate input
	if title == "" || isbn == "" {
		http.Error(w, "Title and ISBN are required", http.StatusBadRequest)
		return
	}

	// Check if ISBN exists
	if isbnExists(isbn, ws.library) {
		http.Error(w, "ISBN already exists", http.StatusBadRequest)
		return
	}

	// Process authors from new form structure
	authors := ws.processNewFormAuthors(r)
	if len(authors) == 0 {
		http.Error(w, "At least one valid author is required", http.StatusBadRequest)
		return
	}

	// Create book
	newBook := Book{
		Title:       title,
		ISBN:        isbn,
		Authors:     authors,
		Description: description,
	}

	// Add to library
	ws.library.Books = append(ws.library.Books, newBook)

	// Save to CSV
	err := saveBooksToCSV(ws.library.Books)
	if err != nil {
		http.Error(w, "Error saving book: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Redirect to books page
	http.Redirect(w, r, "/books", http.StatusSeeOther)
}

func (ws *WebServer) processAddMagazine(w http.ResponseWriter, r *http.Request) {
	title := strings.TrimSpace(r.FormValue("title"))
	isbn := strings.TrimSpace(r.FormValue("isbn"))
	publishedStr := strings.TrimSpace(r.FormValue("published"))

	// Validate input
	if title == "" || isbn == "" || publishedStr == "" {
		http.Error(w, "Title, ISBN, and Publication Date are required", http.StatusBadRequest)
		return
	}

	// Check if ISBN exists
	if isbnExists(isbn, ws.library) {
		http.Error(w, "ISBN already exists", http.StatusBadRequest)
		return
	}

	// Parse date
	publishedAt, err := time.Parse("02.01.2006", publishedStr)
	if err != nil {
		http.Error(w, "Invalid date format. Use DD.MM.YYYY", http.StatusBadRequest)
		return
	}

	// Process authors from new form structure
	authors := ws.processNewFormAuthors(r)
	if len(authors) == 0 {
		http.Error(w, "At least one valid author is required", http.StatusBadRequest)
		return
	}

	// Create magazine
	newMagazine := Magazine{
		Title:       title,
		ISBN:        isbn,
		Authors:     authors,
		PublishedAt: publishedAt,
	}

	// Add to library
	ws.library.Magazines = append(ws.library.Magazines, newMagazine)

	// Save to CSV
	err = saveMagazinesToCSV(ws.library.Magazines)
	if err != nil {
		http.Error(w, "Error saving magazine: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Redirect to magazines page
	http.Redirect(w, r, "/magazines", http.StatusSeeOther)
}

func (ws *WebServer) processNewFormAuthors(r *http.Request) []Author {
	var authors []Author

	// Parse form to get all values
	r.ParseForm()

	// Get all author emails from the form
	authorEmails := r.Form["authors[]"]

	for _, email := range authorEmails {
		email = strings.TrimSpace(email)
		if email == "" {
			continue
		}

		// Check if author already exists
		if existingAuthor, exists := ws.library.Authors[email]; exists {
			authors = append(authors, existingAuthor)
		} else {
			// Look for new author details by iterating through form fields
			var firstName, lastName string

			// Check all form values for matching patterns
			for key, values := range r.Form {
				if strings.HasPrefix(key, "author_first_") && len(values) > 0 {
					firstName = strings.TrimSpace(values[0])
					// Look for corresponding last name
					lastKey := strings.Replace(key, "author_first_", "author_last_", 1)
					if lastValues, exists := r.Form[lastKey]; exists && len(lastValues) > 0 {
						lastName = strings.TrimSpace(lastValues[0])
						break
					}
				}
			}

			if firstName != "" && lastName != "" {
				// Create new author
				newAuthor := Author{
					Email:     email,
					FirstName: firstName,
					LastName:  lastName,
				}

				ws.library.Authors[email] = newAuthor
				saveAuthorsToCSV(ws.library.Authors) // Save new author
				authors = append(authors, newAuthor)
			}
		}
	}

	return authors
}
