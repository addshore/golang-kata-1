package main

import (
	"encoding/csv"
	"errors"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

const (
	serverAddr     = ":8080"
	resourcesDir   = "resources"
	authorsCSVName = "authors.csv"
	booksCSVName   = "books.csv"
	magCSVName     = "magazines.csv"
)

type Author struct {
	Email     string
	FirstName string
	LastName  string
}

type Book struct {
	Title        string
	ISBN         string
	AuthorEmails []string
	Description  string
}

type Magazine struct {
	Title        string
	ISBN         string
	AuthorEmails []string
	PublishedAt  string
}

type BookView struct {
	Title       string
	ISBN        string
	Authors     string
	Description string
}

type MagazineView struct {
	Title       string
	ISBN        string
	Authors     string
	PublishedAt string
}

type ItemView struct {
	Title       string
	ItemType    string
	ISBN        string
	Authors     string
	Description string
	PublishedAt string
}

type AddRequest struct {
	ItemType       string
	Title          string
	ISBN           string
	AuthorEmails   []string
	Description    string
	PublishedAt    string
	NewAuthorEmail string
	NewAuthorFirst string
	NewAuthorLast  string
}

type AddPlan struct {
	AddAuthor   *Author
	AddBook     *Book
	AddMagazine *Magazine
}

type LibraryData struct {
	Authors   map[string]Author
	Books     []Book
	Magazines []Magazine
}

type PageData struct {
	Items         []ItemView
	Books         []BookView
	Magazines     []MagazineView
	SearchISBN    string
	SearchEmail   string
	HasSearch     bool
	NoResults     bool
	BookCount     int
	MagazineCount int
	AddMessage    string
	AddError      string
}

func main() {
	if len(os.Args) > 1 && os.Args[1] == "cli" {
		runCLI()
		return
	}
	tmpl := template.Must(template.New("index").Parse(indexHTML))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		switch r.Method {
		case http.MethodPost:
			handleAddSubmission(w, r, tmpl)
		case http.MethodGet:
			serveLibraryPage(w, r, tmpl)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	log.Printf("Library UI running at http://localhost%s", serverAddr)
	if err := http.ListenAndServe(serverAddr, nil); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

func serveLibraryPage(w http.ResponseWriter, r *http.Request, tmpl *template.Template) {
	libraryData, err := loadLibraryData()
	if err != nil {
		http.Error(w, "failed to load library data", http.StatusInternalServerError)
		return
	}
	isbn := strings.TrimSpace(r.URL.Query().Get("isbn"))
	authorEmail := strings.TrimSpace(r.URL.Query().Get("author"))
	pageData := buildPageData(libraryData, isbn, authorEmail)
	if r.URL.Query().Get("added") != "" {
		pageData.AddMessage = "Library updated successfully."
	}
	if err := tmpl.Execute(w, pageData); err != nil {
		log.Printf("template execute error: %v", err)
	}
}

func handleAddSubmission(w http.ResponseWriter, r *http.Request, tmpl *template.Template) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid form submission", http.StatusBadRequest)
		return
	}

	libraryData, err := loadLibraryData()
	if err != nil {
		http.Error(w, "failed to load library data", http.StatusInternalServerError)
		return
	}

	addRequest := parseAddRequest(r)
	plan, err := validateAddRequest(libraryData, addRequest)
	if err != nil {
		pageData := buildPageData(libraryData, "", "")
		pageData.AddError = err.Error()
		if err := tmpl.Execute(w, pageData); err != nil {
			log.Printf("template execute error: %v", err)
		}
		return
	}

	if err := applyAddPlan(plan); err != nil {
		pageData := buildPageData(libraryData, "", "")
		pageData.AddError = "Failed to save the new entry."
		if err := tmpl.Execute(w, pageData); err != nil {
			log.Printf("template execute error: %v", err)
		}
		return
	}

	http.Redirect(w, r, "/?added=1", http.StatusSeeOther)
}

func loadLibraryData() (LibraryData, error) {
	authors, err := loadAuthors()
	if err != nil {
		return LibraryData{}, err
	}

	books, err := loadBooks()
	if err != nil {
		return LibraryData{}, err
	}

	magazines, err := loadMagazines()
	if err != nil {
		return LibraryData{}, err
	}

	return LibraryData{
		Authors:   authors,
		Books:     books,
		Magazines: magazines,
	}, nil
}

func buildPageData(data LibraryData, isbn string, authorEmail string) PageData {
	filteredBooks := data.Books
	filteredMagazines := data.Magazines
	if isbn != "" {
		filteredBooks = filterBooksByISBN(filteredBooks, isbn)
		filteredMagazines = filterMagazinesByISBN(filteredMagazines, isbn)
	}
	if authorEmail != "" {
		filteredBooks = filterBooksByAuthorEmail(filteredBooks, authorEmail)
		filteredMagazines = filterMagazinesByAuthorEmail(filteredMagazines, authorEmail)
	}

	bookViews := make([]BookView, 0, len(filteredBooks))
	for _, book := range filteredBooks {
		bookViews = append(bookViews, BookView{
			Title:       book.Title,
			ISBN:        book.ISBN,
			Authors:     resolveAuthorNames(book.AuthorEmails, data.Authors),
			Description: book.Description,
		})
	}

	magViews := make([]MagazineView, 0, len(filteredMagazines))
	for _, mag := range filteredMagazines {
		magViews = append(magViews, MagazineView{
			Title:       mag.Title,
			ISBN:        mag.ISBN,
			Authors:     resolveAuthorNames(mag.AuthorEmails, data.Authors),
			PublishedAt: mag.PublishedAt,
		})
	}

	items := buildCombinedItems(filteredBooks, filteredMagazines, data.Authors)

	return PageData{
		Items:         items,
		Books:         bookViews,
		Magazines:     magViews,
		SearchISBN:    isbn,
		SearchEmail:   authorEmail,
		HasSearch:     isbn != "" || authorEmail != "",
		NoResults:     (isbn != "" || authorEmail != "") && len(bookViews) == 0 && len(magViews) == 0,
		BookCount:     len(bookViews),
		MagazineCount: len(magViews),
	}
}

func loadAuthors() (map[string]Author, error) {
	path := filepath.Join(resourcesDir, authorsCSVName)
	records, err := readCSV(path)
	if err != nil {
		return nil, err
	}
	if len(records) == 0 {
		return nil, errors.New("authors CSV is empty")
	}

	authors := make(map[string]Author, len(records))
	for _, row := range records[1:] {
		if len(row) < 3 {
			continue
		}
		authors[row[0]] = Author{
			Email:     row[0],
			FirstName: row[1],
			LastName:  row[2],
		}
	}
	return authors, nil
}

func loadBooks() ([]Book, error) {
	path := filepath.Join(resourcesDir, booksCSVName)
	records, err := readCSV(path)
	if err != nil {
		return nil, err
	}
	if len(records) == 0 {
		return nil, errors.New("books CSV is empty")
	}

	books := make([]Book, 0, len(records)-1)
	for _, row := range records[1:] {
		if len(row) < 4 {
			continue
		}
		books = append(books, Book{
			Title:        row[0],
			ISBN:         row[1],
			AuthorEmails: splitAuthors(row[2]),
			Description:  row[3],
		})
	}
	return books, nil
}

func loadMagazines() ([]Magazine, error) {
	path := filepath.Join(resourcesDir, magCSVName)
	records, err := readCSV(path)
	if err != nil {
		return nil, err
	}
	if len(records) == 0 {
		return nil, errors.New("magazines CSV is empty")
	}

	magazines := make([]Magazine, 0, len(records)-1)
	for _, row := range records[1:] {
		if len(row) < 4 {
			continue
		}
		magazines = append(magazines, Magazine{
			Title:        row[0],
			ISBN:         row[1],
			AuthorEmails: splitAuthors(row[2]),
			PublishedAt:  row[3],
		})
	}
	return magazines, nil
}

func readCSV(path string) ([][]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = ';'
	reader.FieldsPerRecord = -1
	return reader.ReadAll()
}

func splitAuthors(raw string) []string {
	if strings.TrimSpace(raw) == "" {
		return nil
	}
	parts := strings.Split(raw, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

func resolveAuthorNames(emails []string, authors map[string]Author) string {
	if len(emails) == 0 {
		return "Unknown"
	}
	resolved := make([]string, 0, len(emails))
	for _, email := range emails {
		if author, ok := authors[email]; ok {
			name := strings.TrimSpace(author.FirstName + " " + author.LastName)
			if name != "" {
				resolved = append(resolved, name)
				continue
			}
		}
		resolved = append(resolved, email)
	}
	return strings.Join(resolved, ", ")
}

func buildCombinedItems(books []Book, magazines []Magazine, authors map[string]Author) []ItemView {
	items := make([]ItemView, 0, len(books)+len(magazines))
	for _, book := range books {
		items = append(items, ItemView{
			Title:       book.Title,
			ItemType:    "Book",
			ISBN:        book.ISBN,
			Authors:     resolveAuthorNames(book.AuthorEmails, authors),
			Description: book.Description,
		})
	}
	for _, mag := range magazines {
		items = append(items, ItemView{
			Title:       mag.Title,
			ItemType:    "Magazine",
			ISBN:        mag.ISBN,
			Authors:     resolveAuthorNames(mag.AuthorEmails, authors),
			PublishedAt: mag.PublishedAt,
		})
	}
	if len(items) == 0 {
		return items
	}
	sort.Slice(items, func(i, j int) bool {
		left := strings.ToLower(strings.TrimSpace(items[i].Title))
		right := strings.ToLower(strings.TrimSpace(items[j].Title))
		if left == right {
			return strings.ToLower(items[i].ItemType) < strings.ToLower(items[j].ItemType)
		}
		return left < right
	})
	return items
}

func filterBooksByISBN(books []Book, isbn string) []Book {
	needle := normalizeISBN(isbn)
	if needle == "" {
		return books
	}
	filtered := make([]Book, 0, len(books))
	for _, book := range books {
		if normalizeISBN(book.ISBN) == needle {
			filtered = append(filtered, book)
		}
	}
	return filtered
}

func filterMagazinesByISBN(magazines []Magazine, isbn string) []Magazine {
	needle := normalizeISBN(isbn)
	if needle == "" {
		return magazines
	}
	filtered := make([]Magazine, 0, len(magazines))
	for _, mag := range magazines {
		if normalizeISBN(mag.ISBN) == needle {
			filtered = append(filtered, mag)
		}
	}
	return filtered
}

func normalizeISBN(isbn string) string {
	return strings.ToLower(strings.TrimSpace(isbn))
}

func filterBooksByAuthorEmail(books []Book, email string) []Book {
	needle := normalizeEmail(email)
	if needle == "" {
		return books
	}
	filtered := make([]Book, 0, len(books))
	for _, book := range books {
		if hasAuthorEmail(book.AuthorEmails, needle) {
			filtered = append(filtered, book)
		}
	}
	return filtered
}

func filterMagazinesByAuthorEmail(magazines []Magazine, email string) []Magazine {
	needle := normalizeEmail(email)
	if needle == "" {
		return magazines
	}
	filtered := make([]Magazine, 0, len(magazines))
	for _, mag := range magazines {
		if hasAuthorEmail(mag.AuthorEmails, needle) {
			filtered = append(filtered, mag)
		}
	}
	return filtered
}

func hasAuthorEmail(emails []string, needle string) bool {
	for _, email := range emails {
		if normalizeEmail(email) == needle {
			return true
		}
	}
	return false
}

func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

func parseAddRequest(r *http.Request) AddRequest {
	return AddRequest{
		ItemType:       strings.TrimSpace(r.FormValue("itemType")),
		Title:          strings.TrimSpace(r.FormValue("title")),
		ISBN:           strings.TrimSpace(r.FormValue("isbn")),
		AuthorEmails:   splitAuthors(r.FormValue("authors")),
		Description:    strings.TrimSpace(r.FormValue("description")),
		PublishedAt:    strings.TrimSpace(r.FormValue("publishedAt")),
		NewAuthorEmail: strings.TrimSpace(r.FormValue("newAuthorEmail")),
		NewAuthorFirst: strings.TrimSpace(r.FormValue("newAuthorFirst")),
		NewAuthorLast:  strings.TrimSpace(r.FormValue("newAuthorLast")),
	}
}

func validateAddRequest(data LibraryData, req AddRequest) (AddPlan, error) {
	itemType := strings.ToLower(strings.TrimSpace(req.ItemType))
	if itemType == "" {
		return AddPlan{}, errors.New("please choose what you want to add")
	}
	if itemType != "book" && itemType != "magazine" && itemType != "author" {
		return AddPlan{}, errors.New("invalid item type")
	}

	newAuthorEmail := normalizeEmail(req.NewAuthorEmail)
	newAuthorFirst := strings.TrimSpace(req.NewAuthorFirst)
	newAuthorLast := strings.TrimSpace(req.NewAuthorLast)
	addAuthor, err := buildAuthorToAdd(data, itemType, newAuthorEmail, newAuthorFirst, newAuthorLast)
	if err != nil {
		return AddPlan{}, err
	}

	if itemType == "author" {
		if addAuthor == nil {
			return AddPlan{}, errors.New("author email and name are required")
		}
		return AddPlan{AddAuthor: addAuthor}, nil
	}

	title := strings.TrimSpace(req.Title)
	if title == "" {
		return AddPlan{}, errors.New("title is required")
	}
	isbn := strings.TrimSpace(req.ISBN)
	if isbn == "" {
		return AddPlan{}, errors.New("ISBN is required")
	}
	if isbnExists(data, isbn) {
		return AddPlan{}, errors.New("ISBN already exists in the library")
	}

	authorEmails := normalizeEmailList(req.AuthorEmails)
	if len(authorEmails) == 0 {
		return AddPlan{}, errors.New("at least one author email is required")
	}
	if newAuthorEmail != "" && !containsEmail(authorEmails, newAuthorEmail) {
		authorEmails = append(authorEmails, newAuthorEmail)
	}
	if err := ensureAuthorsExist(data, authorEmails, newAuthorEmail, newAuthorFirst, newAuthorLast); err != nil {
		return AddPlan{}, err
	}

	plan := AddPlan{AddAuthor: addAuthor}
	if itemType == "book" {
		plan.AddBook = &Book{
			Title:        title,
			ISBN:         isbn,
			AuthorEmails: authorEmails,
			Description:  req.Description,
		}
		return plan, nil
	}

	if itemType == "magazine" {
		if strings.TrimSpace(req.PublishedAt) == "" {
			return AddPlan{}, errors.New("published date is required for magazines")
		}
		plan.AddMagazine = &Magazine{
			Title:        title,
			ISBN:         isbn,
			AuthorEmails: authorEmails,
			PublishedAt:  req.PublishedAt,
		}
		return plan, nil
	}

	return AddPlan{}, errors.New("unsupported item type")
}

func buildAuthorToAdd(data LibraryData, itemType string, email string, first string, last string) (*Author, error) {
	if email == "" && first == "" && last == "" {
		return nil, nil
	}
	if email == "" || first == "" || last == "" {
		return nil, errors.New("new author requires email, first name, and last name")
	}
	if _, exists := data.Authors[email]; exists {
		if itemType == "author" {
			return nil, errors.New("author already exists")
		}
		return nil, nil
	}
	return &Author{Email: email, FirstName: first, LastName: last}, nil
}

func isbnExists(data LibraryData, isbn string) bool {
	needle := normalizeISBN(isbn)
	for _, book := range data.Books {
		if normalizeISBN(book.ISBN) == needle {
			return true
		}
	}
	for _, mag := range data.Magazines {
		if normalizeISBN(mag.ISBN) == needle {
			return true
		}
	}
	return false
}

func normalizeEmailList(emails []string) []string {
	result := make([]string, 0, len(emails))
	for _, email := range emails {
		normalized := normalizeEmail(email)
		if normalized == "" {
			continue
		}
		if !containsEmail(result, normalized) {
			result = append(result, normalized)
		}
	}
	return result
}

func containsEmail(emails []string, email string) bool {
	needle := normalizeEmail(email)
	for _, entry := range emails {
		if normalizeEmail(entry) == needle {
			return true
		}
	}
	return false
}

func ensureAuthorsExist(data LibraryData, emails []string, newAuthorEmail string, newAuthorFirst string, newAuthorLast string) error {
	for _, email := range emails {
		if _, ok := data.Authors[email]; ok {
			continue
		}
		if email == newAuthorEmail && newAuthorFirst != "" && newAuthorLast != "" {
			continue
		}
		return errors.New("all author emails must already exist or be provided with a new author")
	}
	return nil
}

func applyAddPlan(plan AddPlan) error {
	if plan.AddAuthor != nil {
		if err := appendAuthor(*plan.AddAuthor); err != nil {
			return err
		}
	}
	if plan.AddBook != nil {
		if err := appendBook(*plan.AddBook); err != nil {
			return err
		}
	}
	if plan.AddMagazine != nil {
		if err := appendMagazine(*plan.AddMagazine); err != nil {
			return err
		}
	}
	return nil
}

func appendAuthor(author Author) error {
	path := filepath.Join(resourcesDir, authorsCSVName)
	return appendCSVRecord(path, []string{author.Email, author.FirstName, author.LastName})
}

func appendBook(book Book) error {
	path := filepath.Join(resourcesDir, booksCSVName)
	authors := strings.Join(book.AuthorEmails, ",")
	return appendCSVRecord(path, []string{book.Title, book.ISBN, authors, book.Description})
}

func appendMagazine(magazine Magazine) error {
	path := filepath.Join(resourcesDir, magCSVName)
	authors := strings.Join(magazine.AuthorEmails, ",")
	return appendCSVRecord(path, []string{magazine.Title, magazine.ISBN, authors, magazine.PublishedAt})
}

func appendCSVRecord(path string, record []string) error {
	file, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	writer.Comma = ';'
	if err := writer.Write(record); err != nil {
		return err
	}
	writer.Flush()
	return writer.Error()
}

const indexHTML = `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1" />
  <title>Library Catalog</title>
  <style>
    :root {
      color-scheme: light;
      --bg: #0f172a;
      --surface: #111827;
      --card: #1f2937;
      --text: #f8fafc;
      --muted: #cbd5f5;
      --accent: #38bdf8;
      --border: rgba(148, 163, 184, 0.25);
    }
    * { box-sizing: border-box; }
    body {
      margin: 0;
      font-family: "Segoe UI", system-ui, -apple-system, sans-serif;
      background: linear-gradient(135deg, #0b1120 0%, #111827 100%);
      color: var(--text);
    }
    header {
      padding: 3rem 2rem 2rem;
      text-align: center;
    }
    header h1 {
      margin: 0 0 0.5rem;
      font-size: 2.5rem;
    }
    header p {
      margin: 0;
      color: var(--muted);
      font-size: 1.1rem;
    }
		.search-form {
			margin-top: 1.5rem;
			display: flex;
			flex-wrap: wrap;
			justify-content: center;
			gap: 0.75rem;
		}
		.search-form input {
			min-width: 240px;
			padding: 0.6rem 0.9rem;
			border-radius: 10px;
			border: 1px solid var(--border);
			background: #0b1221;
			color: var(--text);
			font-size: 1rem;
		}
		.search-form button,
		.search-form a {
			padding: 0.6rem 1.2rem;
			border-radius: 999px;
			border: none;
			font-weight: 600;
			cursor: pointer;
			text-decoration: none;
			transition: transform 0.2s ease, box-shadow 0.2s ease;
		}
		.search-form button {
			background: var(--accent);
			color: #0b1221;
			box-shadow: 0 8px 20px rgba(56, 189, 248, 0.35);
		}
		.search-form a {
			background: rgba(148, 163, 184, 0.2);
			color: var(--text);
			border: 1px solid var(--border);
		}
		.search-form button:hover,
		.search-form a:hover {
			transform: translateY(-1px);
		}
		.notice {
			margin-top: 1rem;
			padding: 0.75rem 1rem;
			border-radius: 10px;
			background: rgba(239, 68, 68, 0.15);
			color: #fecaca;
			border: 1px solid rgba(239, 68, 68, 0.4);
		}
		.success {
			margin-top: 1rem;
			padding: 0.75rem 1rem;
			border-radius: 10px;
			background: rgba(34, 197, 94, 0.15);
			color: #bbf7d0;
			border: 1px solid rgba(34, 197, 94, 0.4);
		}
		.add-form {
			margin-top: 1rem;
			display: grid;
			gap: 1rem;
		}
		.add-form label {
			display: grid;
			gap: 0.5rem;
			font-size: 0.95rem;
			color: var(--muted);
		}
		.add-form input,
		.add-form select,
		.add-form textarea {
			width: 100%;
			padding: 0.6rem 0.9rem;
			border-radius: 10px;
			border: 1px solid var(--border);
			background: #0b1221;
			color: var(--text);
			font-size: 1rem;
		}
		.add-form textarea {
			resize: vertical;
		}
		.form-grid {
			display: grid;
			gap: 1rem;
			grid-template-columns: repeat(auto-fit, minmax(240px, 1fr));
		}
		.form-grid .full {
			grid-column: 1 / -1;
		}
		.add-form .primary {
			justify-self: start;
			padding: 0.7rem 1.4rem;
			border-radius: 999px;
			border: none;
			background: var(--accent);
			color: #0b1221;
			font-weight: 700;
			cursor: pointer;
			box-shadow: 0 8px 20px rgba(56, 189, 248, 0.35);
		}
    main {
      max-width: 1100px;
      margin: 0 auto 3rem;
      padding: 0 1.5rem;
      display: grid;
      gap: 2rem;
    }
    section {
      background: var(--surface);
      border-radius: 16px;
      padding: 1.5rem;
      border: 1px solid var(--border);
      box-shadow: 0 12px 30px rgba(15, 23, 42, 0.4);
    }
    section h2 {
      margin-top: 0;
      font-size: 1.6rem;
    }
    .grid {
      display: grid;
      grid-template-columns: repeat(auto-fit, minmax(260px, 1fr));
      gap: 1rem;
    }
    .card {
      background: var(--card);
      border-radius: 14px;
      padding: 1rem;
      border: 1px solid var(--border);
      display: flex;
      flex-direction: column;
      gap: 0.6rem;
    }
    .card h3 {
      margin: 0;
      font-size: 1.1rem;
    }
    .meta {
      font-size: 0.95rem;
      color: var(--muted);
    }
    .badge {
      display: inline-flex;
      align-items: center;
      gap: 0.4rem;
      background: rgba(56, 189, 248, 0.15);
      color: var(--accent);
      padding: 0.2rem 0.6rem;
      border-radius: 999px;
      font-size: 0.8rem;
      font-weight: 600;
    }
    .description {
      font-size: 0.95rem;
      line-height: 1.4;
      color: #e2e8f0;
    }
  </style>
</head>
<body>
  <header>
    <h1>Library Catalog</h1>
    <p>Browse every book and magazine stored in the library archive.</p>
		<form class="search-form" method="get" action="/">
			<input type="search" name="isbn" placeholder="Search by ISBN" value="{{.SearchISBN}}" />
			<input type="search" name="author" placeholder="Search by author email" value="{{.SearchEmail}}" />
			<button type="submit">Search</button>
			{{if .HasSearch}}
			<a href="/">Clear</a>
			{{end}}
		</form>
		{{if .HasSearch}}
		<p class="meta">
			Filters:
			{{if .SearchISBN}}ISBN {{.SearchISBN}}{{end}}
			{{if and .SearchISBN .SearchEmail}} · {{end}}
			{{if .SearchEmail}}Author {{.SearchEmail}}{{end}}
			({{.BookCount}} books, {{.MagazineCount}} magazines)
		</p>
		{{end}}
		{{if .NoResults}}
		<div class="notice">No items matched those filters. Try another search or clear the filter.</div>
		{{end}}
  </header>
  <main>
    <section>
			<h2>Library Items (Sorted by Title)</h2>
      <div class="grid">
				{{range .Items}}
        <article class="card">
					<div class="badge">{{.ItemType}}</div>
          <h3>{{.Title}}</h3>
          <div class="meta">ISBN: {{.ISBN}}</div>
          <div class="meta">Authors: {{.Authors}}</div>
					{{if .PublishedAt}}
					<div class="meta">Published: {{.PublishedAt}}</div>
					{{end}}
					{{if .Description}}
					<div class="description">{{.Description}}</div>
					{{end}}
        </article>
        {{end}}
      </div>
    </section>
		<section>
			<h2>Add to Library</h2>
			{{if .AddMessage}}
			<div class="success">{{.AddMessage}}</div>
			{{end}}
			{{if .AddError}}
			<div class="notice">{{.AddError}}</div>
			{{end}}
			<form class="add-form" method="post" action="/">
				<label>
					What are you adding?
					<select name="itemType" required>
						<option value="">Choose one...</option>
						<option value="book">Book</option>
						<option value="magazine">Magazine</option>
						<option value="author">Author only</option>
					</select>
				</label>
				<div class="form-grid">
					<label>
						Title
						<input type="text" name="title" placeholder="Title" />
					</label>
					<label>
						ISBN
						<input type="text" name="isbn" placeholder="ISBN" />
					</label>
					<label>
						Author emails (comma separated)
						<input type="text" name="authors" placeholder="author@example.com, coauthor@example.com" />
					</label>
					<label>
						Published date (for magazines)
						<input type="text" name="publishedAt" placeholder="DD.MM.YYYY" />
					</label>
					<label class="full">
						Description (for books)
						<textarea name="description" rows="3" placeholder="Short description"></textarea>
					</label>
				</div>
				<h3>New author details (optional)</h3>
				<p class="meta">If the author email does not exist yet, add their details here in the same action.</p>
				<div class="form-grid">
					<label>
						Author email
						<input type="email" name="newAuthorEmail" placeholder="new.author@example.com" />
					</label>
					<label>
						First name
						<input type="text" name="newAuthorFirst" placeholder="First name" />
					</label>
					<label>
						Last name
						<input type="text" name="newAuthorLast" placeholder="Last name" />
					</label>
				</div>
				<button class="primary" type="submit">Add to library</button>
			</form>
		</section>
  </main>
</body>
</html>`
