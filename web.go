package main

import (
	"html/template"
	"net/http"
	"strings"
)

type WebLibrary struct {
	*Library
}

func (wl *WebLibrary) reload() {
	authors := loadAuthors()
	books := loadBooks()
	magazines := loadMagazines()
	wl.Library = NewLibrary(authors, books, magazines)
}

func startWebServer() {
	authors := loadAuthors()
	books := loadBooks()
	magazines := loadMagazines()
	webLib := &WebLibrary{NewLibrary(authors, books, magazines)}

	http.HandleFunc("/", webLib.homeHandler)
	http.HandleFunc("/search", webLib.searchHandler)
	http.HandleFunc("/sort", webLib.sortHandler)
	http.HandleFunc("/add", webLib.addHandler)

	println("Web server starting at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func (wl *WebLibrary) homeHandler(w http.ResponseWriter, r *http.Request) {
	wl.reload()
	tmpl := `<!DOCTYPE html>
<html>
<head><title>Library System</title></head>
<body>
<h1>📚 Library System</h1>
<nav>
<a href="/">Home</a> | 
<a href="/search">Search</a> | 
<a href="/sort">Sort</a> | 
<a href="/add">Add Item</a>
</nav>

<h2>📚 Books ({{len .Books}})</h2>
{{range .Books}}
<div style="border:1px solid #ccc; margin:10px; padding:10px;">
<strong>{{.Title}}</strong><br>
ISBN: {{.ISBN}}<br>
Authors: {{formatAuthors .Authors $.Authors}}<br>
Description: {{.Description}}
</div>
{{end}}

<h2>📖 Magazines ({{len .Magazines}})</h2>
{{range .Magazines}}
<div style="border:1px solid #ccc; margin:10px; padding:10px;">
<strong>{{.Title}}</strong><br>
ISBN: {{.ISBN}}<br>
Authors: {{formatAuthors .Authors $.Authors}}<br>
Published: {{.PublishedAt}}
</div>
{{end}}
</body>
</html>`

	funcMap := template.FuncMap{
		"formatAuthors": func(emails []string, authors map[string]Author) string {
			return FormatAuthors(emails, authors)
		},
	}

	t := template.Must(template.New("home").Funcs(funcMap).Parse(tmpl))
	t.Execute(w, wl.Library)
}

func (wl *WebLibrary) searchHandler(w http.ResponseWriter, r *http.Request) {
	wl.reload()
	if r.Method == "POST" {
		r.ParseForm()
		searchType := r.FormValue("type")
		query := r.FormValue("query")

		tmpl := `<!DOCTYPE html>
<html>
<head><title>Search Results</title></head>
<body>
<h1>📚 Search Results</h1>
<nav><a href="/">Home</a> | <a href="/search">Search</a> | <a href="/sort">Sort</a> | <a href="/add">Add Item</a></nav>
<h2>Results for: {{.Query}}</h2>
{{if .Found}}
{{if .Book}}
<div style="border:1px solid #ccc; margin:10px; padding:10px;">
<strong>[BOOK] {{.Book.Title}}</strong><br>
ISBN: {{.Book.ISBN}}<br>
Authors: {{formatAuthors .Book.Authors .Authors}}<br>
Description: {{.Book.Description}}
</div>
{{end}}
{{if .Magazine}}
<div style="border:1px solid #ccc; margin:10px; padding:10px;">
<strong>[MAGAZINE] {{.Magazine.Title}}</strong><br>
ISBN: {{.Magazine.ISBN}}<br>
Authors: {{formatAuthors .Magazine.Authors .Authors}}<br>
Published: {{.Magazine.PublishedAt}}
</div>
{{end}}
{{if .Books}}
<h3>📚 Books</h3>
{{range .Books}}
<div style="border:1px solid #ccc; margin:10px; padding:10px;">
<strong>{{.Title}}</strong><br>
ISBN: {{.ISBN}}<br>
Authors: {{formatAuthors .Authors $.Authors}}<br>
Description: {{.Description}}
</div>
{{end}}
{{end}}
{{if .Magazines}}
<h3>📖 Magazines</h3>
{{range .Magazines}}
<div style="border:1px solid #ccc; margin:10px; padding:10px;">
<strong>{{.Title}}</strong><br>
ISBN: {{.ISBN}}<br>
Authors: {{formatAuthors .Authors $.Authors}}<br>
Published: {{.PublishedAt}}
</div>
{{end}}
{{end}}
{{else}}
<p>No results found.</p>
{{end}}
<a href="/search">Search Again</a>
</body>
</html>`

		data := struct {
			Query     string
			Found     bool
			Book      *Book
			Magazine  *Magazine
			Books     []Book
			Magazines []Magazine
			Authors   map[string]Author
		}{
			Query:   query,
			Authors: wl.Authors,
		}

		if searchType == "isbn" {
			if item, found := wl.FindByISBN(query); found {
				data.Found = true
				if book, ok := item.(Book); ok {
					data.Book = &book
				} else if magazine, ok := item.(Magazine); ok {
					data.Magazine = &magazine
				}
			}
		} else if searchType == "author" {
			books, magazines := wl.FindByAuthor(query)
			if len(books) > 0 || len(magazines) > 0 {
				data.Found = true
				data.Books = books
				data.Magazines = magazines
			}
		}

		funcMap := template.FuncMap{
			"formatAuthors": func(emails []string, authors map[string]Author) string {
				return FormatAuthors(emails, authors)
			},
		}

		t := template.Must(template.New("results").Funcs(funcMap).Parse(tmpl))
		t.Execute(w, data)
		return
	}

	tmpl := `<!DOCTYPE html>
<html>
<head><title>Search</title></head>
<body>
<h1>📚 Search Library</h1>
<nav><a href="/">Home</a> | <a href="/search">Search</a> | <a href="/sort">Sort</a> | <a href="/add">Add Item</a></nav>
<form method="post">
<p>
<input type="radio" name="type" value="isbn" checked> Search by ISBN<br>
<input type="radio" name="type" value="author"> Search by Author Email
</p>
<p>
<input type="text" name="query" placeholder="Enter ISBN or email" required style="width:300px;">
</p>
<p>
<input type="submit" value="Search">
</p>
</form>
</body>
</html>`

	t := template.Must(template.New("search").Parse(tmpl))
	t.Execute(w, nil)
}

func (wl *WebLibrary) sortHandler(w http.ResponseWriter, r *http.Request) {
	wl.reload()
	desc := r.URL.Query().Get("desc") == "true"
	items := wl.GetSortedItems(desc)

	tmpl := `<!DOCTYPE html>
<html>
<head><title>Sorted Items</title></head>
<body>
<h1>📚📖 Sorted Items</h1>
<nav><a href="/">Home</a> | <a href="/search">Search</a> | <a href="/sort">Sort</a> | <a href="/add">Add Item</a></nav>
<p>
<a href="/sort?desc=false">Sort A-Z</a> | 
<a href="/sort?desc=true">Sort Z-A</a>
</p>
<h2>{{if .Desc}}Descending{{else}}Ascending{{end}} Order ({{len .Items}} items)</h2>
{{range .Items}}
<div style="border:1px solid #ccc; margin:10px; padding:10px;">
{{if .IsBook}}
<strong>[BOOK] {{.Title}}</strong><br>
ISBN: {{.ISBN}}<br>
Authors: {{formatAuthors .Authors $.Authors}}<br>
Description: {{.Description}}
{{else}}
<strong>[MAGAZINE] {{.Title}}</strong><br>
ISBN: {{.ISBN}}<br>
Authors: {{formatAuthors .Authors $.Authors}}<br>
Published: {{.PublishedAt}}
{{end}}
</div>
{{end}}
</body>
</html>`

	data := struct {
		Items   []Item
		Desc    bool
		Authors map[string]Author
	}{
		Items:   items,
		Desc:    desc,
		Authors: wl.Authors,
	}

	funcMap := template.FuncMap{
		"formatAuthors": func(emails []string, authors map[string]Author) string {
			return FormatAuthors(emails, authors)
		},
	}

	t := template.Must(template.New("sort").Funcs(funcMap).Parse(tmpl))
	t.Execute(w, data)
}

func (wl *WebLibrary) addHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		r.ParseForm()
		itemType := r.FormValue("type")
		title := strings.TrimSpace(r.FormValue("title"))
		isbn := strings.TrimSpace(r.FormValue("isbn"))
		authorEmail := strings.TrimSpace(r.FormValue("author_email"))
		authorFirst := strings.TrimSpace(r.FormValue("author_first"))
		authorLast := strings.TrimSpace(r.FormValue("author_last"))

		if itemType == "book" {
			description := strings.TrimSpace(r.FormValue("description"))
			addBook(title, isbn, authorEmail, authorFirst, authorLast, description)
		} else {
			pubDate := strings.TrimSpace(r.FormValue("pub_date"))
			addMagazine(title, isbn, authorEmail, authorFirst, authorLast, pubDate)
		}

		tmpl := `<!DOCTYPE html>
<html>
<head><title>Item Added</title></head>
<body>
<h1>✅ Success!</h1>
<nav><a href="/">Home</a> | <a href="/search">Search</a> | <a href="/sort">Sort</a> | <a href="/add">Add Item</a></nav>
<p>{{.Type}} "{{.Title}}" has been added successfully!</p>
<p><a href="/">View Library</a> | <a href="/add">Add Another</a></p>
</body>
</html>`

		data := struct {
			Type  string
			Title string
		}{
			Type:  strings.Title(itemType),
			Title: title,
		}

		t := template.Must(template.New("success").Parse(tmpl))
		t.Execute(w, data)
		return
	}

	tmpl := `<!DOCTYPE html>
<html>
<head><title>Add Item</title></head>
<body>
<h1>➕ Add New Item</h1>
<nav><a href="/">Home</a> | <a href="/search">Search</a> | <a href="/sort">Sort</a> | <a href="/add">Add Item</a></nav>
<form method="post">
<p>
<input type="radio" name="type" value="book" checked onchange="toggleFields()"> Book<br>
<input type="radio" name="type" value="magazine" onchange="toggleFields()"> Magazine
</p>
<p>
<label>Title:</label><br>
<input type="text" name="title" required style="width:300px;">
</p>
<p>
<label>ISBN:</label><br>
<input type="text" name="isbn" required style="width:300px;">
</p>
<p>
<label>Author Email:</label><br>
<input type="email" name="author_email" required style="width:300px;">
</p>
<p>
<label>Author First Name:</label><br>
<input type="text" name="author_first" required style="width:300px;">
</p>
<p>
<label>Author Last Name:</label><br>
<input type="text" name="author_last" required style="width:300px;">
</p>
<div id="book-fields">
<p>
<label>Description:</label><br>
<textarea name="description" style="width:300px; height:100px;"></textarea>
</p>
</div>
<div id="magazine-fields" style="display:none;">
<p>
<label>Publication Date (DD.MM.YYYY):</label><br>
<input type="text" name="pub_date" placeholder="01.01.2024" style="width:300px;">
</p>
</div>
<p>
<input type="submit" value="Add Item">
</p>
</form>
<script>
function toggleFields() {
    var bookFields = document.getElementById('book-fields');
    var magazineFields = document.getElementById('magazine-fields');
    var isBook = document.querySelector('input[name="type"]:checked').value === 'book';
    
    bookFields.style.display = isBook ? 'block' : 'none';
    magazineFields.style.display = isBook ? 'none' : 'block';
}
</script>
</body>
</html>`

	t := template.Must(template.New("add").Parse(tmpl))
	t.Execute(w, nil)
}