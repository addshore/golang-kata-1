package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strings"
)

var templates *template.Template
var webLibrary *Library

// StartWebServer starts the web server on the given port
func StartWebServer(library *Library, port int) error {
	webLibrary = library

	// Parse templates
	var err error
	templates, err = template.New("").Funcs(template.FuncMap{
		"formatAuthors": func(emails []string) string {
			return formatAuthors(webLibrary, emails)
		},
		"truncate": truncateString,
	}).Parse(htmlTemplates)
	if err != nil {
		return fmt.Errorf("failed to parse templates: %w", err)
	}

	// Routes
	http.HandleFunc("/", handleHome)
	http.HandleFunc("/books", handleBooks)
	http.HandleFunc("/magazines", handleMagazines)
	http.HandleFunc("/all", handleAll)
	http.HandleFunc("/sorted", handleSorted)
	http.HandleFunc("/search", handleSearch)
	http.HandleFunc("/add-book", handleAddBook)
	http.HandleFunc("/add-magazine", handleAddMagazine)
	http.HandleFunc("/api/authors", handleAPIAuthors)

	fmt.Printf("Starting web server at http://localhost:%d\n", port)
	fmt.Println("Press Ctrl+C to stop the server")
	return http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	templates.ExecuteTemplate(w, "home", nil)
}

func handleBooks(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Books []Book
		Added bool
	}{
		Books: webLibrary.Books,
		Added: r.URL.Query().Get("added") == "1",
	}
	templates.ExecuteTemplate(w, "books", data)
}

func handleMagazines(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Magazines []Magazine
		Added     bool
	}{
		Magazines: webLibrary.Magazines,
		Added:     r.URL.Query().Get("added") == "1",
	}
	templates.ExecuteTemplate(w, "magazines", data)
}

func handleAll(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Books     []Book
		Magazines []Magazine
	}{
		Books:     webLibrary.Books,
		Magazines: webLibrary.Magazines,
	}
	templates.ExecuteTemplate(w, "all", data)
}

func handleSorted(w http.ResponseWriter, r *http.Request) {
	direction := r.URL.Query().Get("dir")
	ascending := direction != "desc"
	items := getAllItemsSortedByTitle(webLibrary, ascending)
	data := struct {
		Items     []LibraryItem
		Ascending bool
	}{
		Items:     items,
		Ascending: ascending,
	}
	templates.ExecuteTemplate(w, "sorted", data)
}

func handleSearch(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	searchType := r.URL.Query().Get("type")

	var results []SearchResult
	if query != "" {
		if searchType == "isbn" {
			results = findByISBN(webLibrary, query)
		} else if searchType == "email" {
			results = findByAuthorEmail(webLibrary, query)
		}
	}

	// Get author info if searching by email
	var author *Author
	if searchType == "email" && query != "" {
		if a, exists := webLibrary.Authors[query]; exists {
			author = &a
		}
	}

	data := struct {
		Query      string
		SearchType string
		Results    []SearchResult
		Author     *Author
	}{
		Query:      query,
		SearchType: searchType,
		Results:    results,
		Author:     author,
	}
	templates.ExecuteTemplate(w, "search", data)
}

func handleAddBook(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		templates.ExecuteTemplate(w, "add-book", nil)
		return
	}

	// POST - add the book
	r.ParseForm()
	title := strings.TrimSpace(r.FormValue("title"))
	isbn := strings.TrimSpace(r.FormValue("isbn"))
	description := strings.TrimSpace(r.FormValue("description"))
	authorEmails := strings.Split(r.FormValue("authors"), ",")

	// Clean author emails
	var cleanEmails []string
	for _, email := range authorEmails {
		email = strings.TrimSpace(email)
		if email != "" {
			cleanEmails = append(cleanEmails, email)
		}
	}

	// Validate
	var errorMsg string
	if title == "" {
		errorMsg = "Title is required"
	} else if isbn == "" {
		errorMsg = "ISBN is required"
	} else if len(cleanEmails) == 0 {
		errorMsg = "At least one author email is required"
	} else if len(findByISBN(webLibrary, isbn)) > 0 {
		errorMsg = "ISBN already exists"
	}

	if errorMsg != "" {
		data := struct {
			Error       string
			Title       string
			ISBN        string
			Description string
			Authors     string
		}{
			Error:       errorMsg,
			Title:       title,
			ISBN:        isbn,
			Description: description,
			Authors:     r.FormValue("authors"),
		}
		templates.ExecuteTemplate(w, "add-book", data)
		return
	}

	// Check for new authors
	newAuthorEmails := r.FormValue("new_author_emails")
	newAuthorFirstNames := r.FormValue("new_author_firstnames")
	newAuthorLastNames := r.FormValue("new_author_lastnames")

	if newAuthorEmails != "" {
		emails := strings.Split(newAuthorEmails, ",")
		firstNames := strings.Split(newAuthorFirstNames, ",")
		lastNames := strings.Split(newAuthorLastNames, ",")

		for i, email := range emails {
			email = strings.TrimSpace(email)
			if email != "" && i < len(firstNames) && i < len(lastNames) {
				author := Author{
					Email:     email,
					FirstName: strings.TrimSpace(firstNames[i]),
					LastName:  strings.TrimSpace(lastNames[i]),
				}
				webLibrary.Authors[email] = author
			}
		}
		saveAuthors(webLibrary)
	}

	// Create book
	book := Book{
		Title:        title,
		ISBN:         isbn,
		AuthorEmails: cleanEmails,
		Description:  description,
	}
	webLibrary.Books = append(webLibrary.Books, book)
	saveBooks(webLibrary)

	http.Redirect(w, r, "/books?added=1", http.StatusSeeOther)
}

func handleAddMagazine(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		templates.ExecuteTemplate(w, "add-magazine", nil)
		return
	}

	// POST - add the magazine
	r.ParseForm()
	title := strings.TrimSpace(r.FormValue("title"))
	isbn := strings.TrimSpace(r.FormValue("isbn"))
	publishedAt := strings.TrimSpace(r.FormValue("published_at"))
	authorEmails := strings.Split(r.FormValue("authors"), ",")

	// Clean author emails
	var cleanEmails []string
	for _, email := range authorEmails {
		email = strings.TrimSpace(email)
		if email != "" {
			cleanEmails = append(cleanEmails, email)
		}
	}

	// Validate
	var errorMsg string
	if title == "" {
		errorMsg = "Title is required"
	} else if isbn == "" {
		errorMsg = "ISBN is required"
	} else if publishedAt == "" {
		errorMsg = "Publication date is required"
	} else if len(cleanEmails) == 0 {
		errorMsg = "At least one author email is required"
	} else if len(findByISBN(webLibrary, isbn)) > 0 {
		errorMsg = "ISBN already exists"
	}

	if errorMsg != "" {
		data := struct {
			Error       string
			Title       string
			ISBN        string
			PublishedAt string
			Authors     string
		}{
			Error:       errorMsg,
			Title:       title,
			ISBN:        isbn,
			PublishedAt: publishedAt,
			Authors:     r.FormValue("authors"),
		}
		templates.ExecuteTemplate(w, "add-magazine", data)
		return
	}

	// Check for new authors
	newAuthorEmails := r.FormValue("new_author_emails")
	newAuthorFirstNames := r.FormValue("new_author_firstnames")
	newAuthorLastNames := r.FormValue("new_author_lastnames")

	if newAuthorEmails != "" {
		emails := strings.Split(newAuthorEmails, ",")
		firstNames := strings.Split(newAuthorFirstNames, ",")
		lastNames := strings.Split(newAuthorLastNames, ",")

		for i, email := range emails {
			email = strings.TrimSpace(email)
			if email != "" && i < len(firstNames) && i < len(lastNames) {
				author := Author{
					Email:     email,
					FirstName: strings.TrimSpace(firstNames[i]),
					LastName:  strings.TrimSpace(lastNames[i]),
				}
				webLibrary.Authors[email] = author
			}
		}
		saveAuthors(webLibrary)
	}

	// Create magazine
	magazine := Magazine{
		Title:        title,
		ISBN:         isbn,
		AuthorEmails: cleanEmails,
		PublishedAt:  publishedAt,
	}
	webLibrary.Magazines = append(webLibrary.Magazines, magazine)
	saveMagazines(webLibrary)

	http.Redirect(w, r, "/magazines?added=1", http.StatusSeeOther)
}

func handleAPIAuthors(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	w.Header().Set("Content-Type", "application/json")

	if author, exists := webLibrary.Authors[email]; exists {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"exists":    true,
			"email":     author.Email,
			"firstName": author.FirstName,
			"lastName":  author.LastName,
		})
	} else {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"exists": false,
		})
	}
}

const htmlTemplates = `
{{define "header"}}
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Library Application</title>
    <style>
        body { font-family: Arial, sans-serif; max-width: 1000px; margin: 0 auto; padding: 20px; }
        nav { background: #333; padding: 10px; margin-bottom: 20px; }
        nav a { color: white; margin-right: 15px; text-decoration: none; }
        nav a:hover { text-decoration: underline; }
        table { width: 100%; border-collapse: collapse; margin-top: 10px; }
        th, td { border: 1px solid #ddd; padding: 8px; text-align: left; }
        th { background: #333; color: white; }
        tr:nth-child(even) { background: #f9f9f9; }
        .item { border: 1px solid #ddd; padding: 15px; margin: 10px 0; }
        .item h3 { margin-top: 0; }
        .book { border-left: 4px solid #4CAF50; }
        .magazine { border-left: 4px solid #2196F3; }
        form { max-width: 500px; }
        label { display: block; margin-top: 10px; font-weight: bold; }
        input[type="text"], textarea { width: 100%; padding: 8px; margin-top: 5px; box-sizing: border-box; }
        button { margin-top: 15px; padding: 10px 20px; background: #333; color: white; border: none; cursor: pointer; }
        button:hover { background: #555; }
        .error { color: red; padding: 10px; background: #ffeeee; margin-bottom: 10px; }
        .success { color: green; padding: 10px; background: #eeffee; margin-bottom: 10px; }
        .search-form { margin-bottom: 20px; }
        .search-form input[type="text"] { width: 300px; display: inline-block; }
        .badge { display: inline-block; padding: 2px 8px; border-radius: 3px; font-size: 12px; }
        .badge-book { background: #4CAF50; color: white; }
        .badge-magazine { background: #2196F3; color: white; }
        .author-input { margin: 5px 0; padding: 5px; background: #f5f5f5; }
        #new-authors input { width: 45%; display: inline-block; margin: 2px; }
    </style>
</head>
<body>
    <h1>📚 Library Application</h1>
    <nav>
        <a href="/">Home</a>
        <a href="/books">Books</a>
        <a href="/magazines">Magazines</a>
        <a href="/all">All Items</a>
        <a href="/sorted">Sorted by Title</a>
        <a href="/search">Search</a>
        <a href="/add-book">Add Book</a>
        <a href="/add-magazine">Add Magazine</a>
    </nav>
{{end}}

{{define "footer"}}
</body>
</html>
{{end}}

{{define "home"}}
{{template "header"}}
<h2>Welcome to the Library</h2>
<p>Use the navigation above to browse books, magazines, search, or add new items.</p>
<h3>Quick Links</h3>
<ul>
    <li><a href="/books">View all Books</a></li>
    <li><a href="/magazines">View all Magazines</a></li>
    <li><a href="/sorted">View all items sorted by title</a></li>
    <li><a href="/search?type=isbn">Search by ISBN</a></li>
    <li><a href="/search?type=email">Search by Author Email</a></li>
</ul>
{{template "footer"}}
{{end}}

{{define "books"}}
{{template "header"}}
<h2>Books</h2>
{{if .Added}}<div class="success">Book added successfully!</div>{{end}}
{{if .Books}}
<table>
    <tr>
        <th>#</th>
        <th>Title</th>
        <th>ISBN</th>
        <th>Authors</th>
        <th>Description</th>
    </tr>
    {{range $i, $book := .Books}}
    <tr>
        <td>{{$i}}</td>
        <td>{{$book.Title}}</td>
        <td>{{$book.ISBN}}</td>
        <td>{{formatAuthors $book.AuthorEmails}}</td>
        <td>{{truncate $book.Description 80}}</td>
    </tr>
    {{end}}
</table>
{{else}}
<p>No books found.</p>
{{end}}
{{template "footer"}}
{{end}}

{{define "magazines"}}
{{template "header"}}
<h2>Magazines</h2>
{{if .Added}}<div class="success">Magazine added successfully!</div>{{end}}
{{if .Magazines}}
<table>
    <tr>
        <th>#</th>
        <th>Title</th>
        <th>ISBN</th>
        <th>Authors</th>
        <th>Published</th>
    </tr>
    {{range $i, $mag := .Magazines}}
    <tr>
        <td>{{$i}}</td>
        <td>{{$mag.Title}}</td>
        <td>{{$mag.ISBN}}</td>
        <td>{{formatAuthors $mag.AuthorEmails}}</td>
        <td>{{$mag.PublishedAt}}</td>
    </tr>
    {{end}}
</table>
{{else}}
<p>No magazines found.</p>
{{end}}
{{template "footer"}}
{{end}}

{{define "all"}}
{{template "header"}}
<h2>All Books and Magazines</h2>
<h3>Books ({{len .Books}})</h3>
{{range .Books}}
<div class="item book">
    <h3>{{.Title}}</h3>
    <p><strong>ISBN:</strong> {{.ISBN}}</p>
    <p><strong>Authors:</strong> {{formatAuthors .AuthorEmails}}</p>
    <p><strong>Description:</strong> {{.Description}}</p>
</div>
{{end}}
<h3>Magazines ({{len .Magazines}})</h3>
{{range .Magazines}}
<div class="item magazine">
    <h3>{{.Title}}</h3>
    <p><strong>ISBN:</strong> {{.ISBN}}</p>
    <p><strong>Authors:</strong> {{formatAuthors .AuthorEmails}}</p>
    <p><strong>Published:</strong> {{.PublishedAt}}</p>
</div>
{{end}}
{{template "footer"}}
{{end}}

{{define "sorted"}}
{{template "header"}}
<h2>All Items Sorted by Title</h2>
<div style="margin-bottom: 15px;">
    <strong>Sort direction:</strong>
    <a href="/sorted?dir=asc" style="margin-left: 10px; {{if .Ascending}}font-weight: bold;{{end}}">↑ A-Z (Ascending)</a>
    <a href="/sorted?dir=desc" style="margin-left: 10px; {{if not .Ascending}}font-weight: bold;{{end}}">↓ Z-A (Descending)</a>
</div>
{{if .Items}}
<table>
    <tr>
        <th>#</th>
        <th>Type</th>
        <th>Title</th>
        <th>ISBN</th>
        <th>Authors</th>
        <th>Details</th>
    </tr>
    {{range $i, $item := .Items}}
    <tr>
        <td>{{$i}}</td>
        <td><span class="badge badge-{{if eq $item.ItemType "Book"}}book{{else}}magazine{{end}}">{{$item.ItemType}}</span></td>
        <td>{{$item.Title}}</td>
        <td>{{$item.ISBN}}</td>
        <td>{{$item.Authors}}</td>
        <td>{{if eq $item.ItemType "Book"}}{{truncate $item.Description 50}}{{else}}{{$item.PublishedAt}}{{end}}</td>
    </tr>
    {{end}}
</table>
{{else}}
<p>No items found.</p>
{{end}}
{{template "footer"}}
{{end}}

{{define "search"}}
{{template "header"}}
<h2>Search</h2>
<div class="search-form">
    <form method="GET">
        <label style="display:inline;">
            <input type="radio" name="type" value="isbn" {{if eq .SearchType "isbn"}}checked{{else if eq .SearchType ""}}checked{{end}}> Search by ISBN
        </label>
        <label style="display:inline; margin-left:20px;">
            <input type="radio" name="type" value="email" {{if eq .SearchType "email"}}checked{{end}}> Search by Author Email
        </label>
        <br><br>
        <input type="text" name="q" value="{{.Query}}" placeholder="Enter ISBN or author email...">
        <button type="submit">Search</button>
    </form>
</div>

{{if .Query}}
<h3>Results for "{{.Query}}"</h3>
{{if .Author}}
<p><strong>Author:</strong> {{.Author.FirstName}} {{.Author.LastName}} ({{.Author.Email}})</p>
{{end}}
{{if .Results}}
{{range .Results}}
{{if eq .Type "Book"}}
<div class="item book">
    <span class="badge badge-book">Book</span>
    <h3>{{.Book.Title}}</h3>
    <p><strong>ISBN:</strong> {{.Book.ISBN}}</p>
    <p><strong>Authors:</strong> {{formatAuthors .Book.AuthorEmails}}</p>
    <p><strong>Description:</strong> {{.Book.Description}}</p>
</div>
{{else}}
<div class="item magazine">
    <span class="badge badge-magazine">Magazine</span>
    <h3>{{.Magazine.Title}}</h3>
    <p><strong>ISBN:</strong> {{.Magazine.ISBN}}</p>
    <p><strong>Authors:</strong> {{formatAuthors .Magazine.AuthorEmails}}</p>
    <p><strong>Published:</strong> {{.Magazine.PublishedAt}}</p>
</div>
{{end}}
{{end}}
{{else}}
<p>No results found.</p>
{{end}}
{{end}}
{{template "footer"}}
{{end}}

{{define "add-book"}}
{{template "header"}}
<h2>Add New Book</h2>
{{if .Error}}<div class="error">{{.Error}}</div>{{end}}
<form method="POST">
    <label>Title *</label>
    <input type="text" name="title" value="{{.Title}}" required>
    
    <label>ISBN *</label>
    <input type="text" name="isbn" value="{{.ISBN}}" required>
    
    <label>Description</label>
    <textarea name="description" rows="4">{{.Description}}</textarea>
    
    <label>Author Emails * (comma-separated)</label>
    <input type="text" name="authors" id="authors" value="{{.Authors}}" required placeholder="email1@example.com, email2@example.com">
    <small>Enter author emails. Click "Check Authors" to add new authors.</small>
    
    <div id="new-authors"></div>
    
    <input type="hidden" name="new_author_emails" id="new_author_emails">
    <input type="hidden" name="new_author_firstnames" id="new_author_firstnames">
    <input type="hidden" name="new_author_lastnames" id="new_author_lastnames">
    
    <button type="button" onclick="checkAuthors()">Check Authors</button>
    <button type="submit">Add Book</button>
</form>
<script>
var newAuthorEmails = [];
async function checkAuthors() {
    const emails = document.getElementById('authors').value.split(',').map(e => e.trim()).filter(e => e);
    const container = document.getElementById('new-authors');
    container.innerHTML = '';
    newAuthorEmails = [];
    
    for (const email of emails) {
        const resp = await fetch('/api/authors?email=' + encodeURIComponent(email));
        const data = await resp.json();
        if (data.exists) {
            container.innerHTML += '<div class="author-input">✓ ' + data.firstName + ' ' + data.lastName + ' (' + email + ')</div>';
        } else {
            const safeEmail = email.replace(/[^a-zA-Z0-9]/g, '_');
            container.innerHTML += '<div class="author-input"><strong>New author: ' + email + '</strong><br>' +
                '<input type="text" placeholder="First name" id="fn_' + safeEmail + '" data-email="' + email + '"> ' +
                '<input type="text" placeholder="Last name" id="ln_' + safeEmail + '"></div>';
            newAuthorEmails.push({email: email, safeEmail: safeEmail});
        }
    }
}

document.querySelector('form').onsubmit = function() {
    const emails = [], firsts = [], lasts = [];
    for (const item of newAuthorEmails) {
        const fn = document.getElementById('fn_' + item.safeEmail);
        const ln = document.getElementById('ln_' + item.safeEmail);
        if (fn && ln && fn.value && ln.value) {
            emails.push(item.email);
            firsts.push(fn.value);
            lasts.push(ln.value);
        }
    }
    document.getElementById('new_author_emails').value = emails.join(',');
    document.getElementById('new_author_firstnames').value = firsts.join(',');
    document.getElementById('new_author_lastnames').value = lasts.join(',');
};
</script>
{{template "footer"}}
{{end}}

{{define "add-magazine"}}
{{template "header"}}
<h2>Add New Magazine</h2>
{{if .Error}}<div class="error">{{.Error}}</div>{{end}}
<form method="POST">
    <label>Title *</label>
    <input type="text" name="title" value="{{.Title}}" required>
    
    <label>ISBN *</label>
    <input type="text" name="isbn" value="{{.ISBN}}" required>
    
    <label>Publication Date * (e.g., 01.01.2024)</label>
    <input type="text" name="published_at" value="{{.PublishedAt}}" required placeholder="DD.MM.YYYY">
    
    <label>Author Emails * (comma-separated)</label>
    <input type="text" name="authors" id="authors" value="{{.Authors}}" required placeholder="email1@example.com, email2@example.com">
    <small>Enter author emails. Click "Check Authors" to add new authors.</small>
    
    <div id="new-authors"></div>
    
    <input type="hidden" name="new_author_emails" id="new_author_emails">
    <input type="hidden" name="new_author_firstnames" id="new_author_firstnames">
    <input type="hidden" name="new_author_lastnames" id="new_author_lastnames">
    
    <button type="button" onclick="checkAuthors()">Check Authors</button>
    <button type="submit">Add Magazine</button>
</form>
<script>
var newAuthorEmails = [];
async function checkAuthors() {
    const emails = document.getElementById('authors').value.split(',').map(e => e.trim()).filter(e => e);
    const container = document.getElementById('new-authors');
    container.innerHTML = '';
    newAuthorEmails = [];
    
    for (const email of emails) {
        const resp = await fetch('/api/authors?email=' + encodeURIComponent(email));
        const data = await resp.json();
        if (data.exists) {
            container.innerHTML += '<div class="author-input">✓ ' + data.firstName + ' ' + data.lastName + ' (' + email + ')</div>';
        } else {
            const safeEmail = email.replace(/[^a-zA-Z0-9]/g, '_');
            container.innerHTML += '<div class="author-input"><strong>New author: ' + email + '</strong><br>' +
                '<input type="text" placeholder="First name" id="fn_' + safeEmail + '" data-email="' + email + '"> ' +
                '<input type="text" placeholder="Last name" id="ln_' + safeEmail + '"></div>';
            newAuthorEmails.push({email: email, safeEmail: safeEmail});
        }
    }
}

document.querySelector('form').onsubmit = function() {
    const emails = [], firsts = [], lasts = [];
    for (const item of newAuthorEmails) {
        const fn = document.getElementById('fn_' + item.safeEmail);
        const ln = document.getElementById('ln_' + item.safeEmail);
        if (fn && ln && fn.value && ln.value) {
            emails.push(item.email);
            firsts.push(fn.value);
            lasts.push(ln.value);
        }
    }
    document.getElementById('new_author_emails').value = emails.join(',');
    document.getElementById('new_author_firstnames').value = firsts.join(',');
    document.getElementById('new_author_lastnames').value = lasts.join(',');
};
</script>
{{template "footer"}}
{{end}}
`
