package main

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"
)

// LibraryItem defines a common interface for items in the library.
type LibraryItem interface {
	GetTitle() string
	GetISBN() string
	GetAuthors() []string
	GetType() string
	// Specific fields
	GetDescription() string
	GetPublishedAt() string
}

// Book represents a book with its details.
type Book struct {
	Title       string
	ISBN        string
	Authors     []string
	Description string
}

func (b Book) GetTitle() string {
	return b.Title
}

func (b Book) GetISBN() string {
	return b.ISBN
}

func (b Book) GetAuthors() []string {
	return b.Authors
}

func (b Book) GetType() string {
	return "Book"
}

func (b Book) GetDescription() string {
	return b.Description
}

func (b Book) GetPublishedAt() string { return "" }

// Magazine represents a magazine with its details.
type Magazine struct {
	Title       string
	ISBN        string
	Authors     []string
	PublishedAt string
}

// Author represents an author with their details.
type Author struct {
	Email     string
	FirstName string
	LastName  string
}

func (m Magazine) GetTitle() string {
	return m.Title
}

func (m Magazine) GetISBN() string {
	return m.ISBN
}

func (m Magazine) GetAuthors() []string {
	return m.Authors
}

func (m Magazine) GetType() string {
	return "Magazine"
}

func (m Magazine) GetDescription() string { return "" }

func (m Magazine) GetPublishedAt() string {
	return m.PublishedAt
}

var templates = template.Must(template.New("").Funcs(template.FuncMap{"join": strings.Join}).ParseGlob("templates/*.html"))

func main() {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/search", searchHandler)
	http.HandleFunc("/add", addHandler)
	http.HandleFunc("/save", saveHandler)

	fmt.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func renderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	buf := new(bytes.Buffer)
	err := templates.ExecuteTemplate(buf, tmpl, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	buf.WriteTo(w)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	books, err := loadBooks("resources/books.csv")
	if err != nil {
		http.Error(w, "Error loading books", http.StatusInternalServerError)
		return
	}

	magazines, err := loadMagazines("resources/magazines.csv")
	if err != nil {
		http.Error(w, "Error loading magazines", http.StatusInternalServerError)
		return
	}

	sortOrder := r.URL.Query().Get("sort")
	if sortOrder == "" {
		sortOrder = "title-asc"
	}

	items, err := getSortedItems(books, magazines, sortOrder)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	data := map[string]interface{}{
		"Items":       items,
		"SortURLASC":  "/?sort=title-asc",
		"SortURLDESC": "/?sort=title-desc",
	}
	renderTemplate(w, "index.html", data)
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("query")
	searchType := r.URL.Query().Get("type")

	if query == "" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	books, err := loadBooks("resources/books.csv")
	if err != nil {
		http.Error(w, "Error loading books", http.StatusInternalServerError)
		return
	}
	magazines, err := loadMagazines("resources/magazines.csv")
	if err != nil {
		http.Error(w, "Error loading magazines", http.StatusInternalServerError)
		return
	}

	var results []LibraryItem
	var message string

	if searchType == "isbn" {
		item := findItemByISBN(query, books, magazines)
		if item != nil {
			results = append(results, item)
		}
	} else if searchType == "author" {
		results = findItemsByAuthor(query, books, magazines)
	}

	if len(results) == 0 {
		message = fmt.Sprintf("No items found for your search.")
	}

	data := map[string]interface{}{
		"Items":       results,
		"Message":     message,
		"Query":       query,
		"SearchType":  searchType,
		"SortURLASC":  fmt.Sprintf("/?sort=title-asc&query=%s&type=%s", query, searchType),
		"SortURLDESC": fmt.Sprintf("/?sort=title-desc&query=%s&type=%s", query, searchType),
	}
	renderTemplate(w, "index.html", data)
}

func addHandler(w http.ResponseWriter, r *http.Request) {
	authorsMap, err := loadAuthors("resources/authors.csv")
	if err != nil {
		http.Error(w, "Could not load authors", http.StatusInternalServerError)
		return
	}

	var authors []Author
	for _, author := range authorsMap {
		authors = append(authors, author)
	}
	sort.Slice(authors, func(i, j int) bool {
		return authors[i].LastName < authors[j].LastName
	})

	data := map[string]interface{}{
		"Authors": authors,
	}
	renderTemplate(w, "add.html", data)
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Could not parse form", http.StatusBadRequest)
		return
	}

	// Collect author emails
	authorEmails := r.Form["existing_authors"]

	// Handle new author
	newAuthorEmail := r.FormValue("new_author_email")
	if newAuthorEmail != "" {
		existingAuthors, err := loadAuthors("resources/authors.csv")
		if err != nil {
			http.Error(w, "Could not load authors for validation", http.StatusInternalServerError)
			return
		}

		if _, exists := existingAuthors[newAuthorEmail]; exists {
			http.Error(w, fmt.Sprintf("Author with email '%s' already exists.", newAuthorEmail), http.StatusBadRequest)
			return
		}

		newAuthor := Author{
			Email:     newAuthorEmail,
			FirstName: r.FormValue("new_author_firstname"),
			LastName:  r.FormValue("new_author_lastname"),
		}
		if newAuthor.FirstName == "" || newAuthor.LastName == "" {
			http.Error(w, "New author must have a first and last name.", http.StatusBadRequest)
			return
		}
		record := []string{newAuthor.Email, newAuthor.FirstName, newAuthor.LastName}
		if err := appendRecord("resources/authors.csv", record); err != nil {
			http.Error(w, "Could not save new author", http.StatusInternalServerError)
			return
		}
		authorEmails = append(authorEmails, newAuthorEmail)
	}

	if len(authorEmails) == 0 {
		http.Error(w, "An item must have at least one author.", http.StatusBadRequest)
		return
	}

	// Save book or magazine
	itemType := r.FormValue("itemType")
	title := r.FormValue("title")
	isbn := r.FormValue("isbn")

	if itemType == "book" {
		book := Book{
			Title:       title,
			ISBN:        isbn,
			Authors:     authorEmails,
			Description: r.FormValue("description"),
		}
		record := []string{book.Title, book.ISBN, strings.Join(book.Authors, ","), book.Description}
		if err := appendRecord("resources/books.csv", record); err != nil {
			http.Error(w, "Could not save book", http.StatusInternalServerError)
			return
		}
	} else if itemType == "magazine" {
		magazine := Magazine{
			Title:       title,
			ISBN:        isbn,
			Authors:     authorEmails,
			PublishedAt: r.FormValue("publishedAt"),
		}
		record := []string{magazine.Title, magazine.ISBN, strings.Join(magazine.Authors, ","), magazine.PublishedAt}
		if err := appendRecord("resources/magazines.csv", record); err != nil {
			http.Error(w, "Could not save magazine", http.StatusInternalServerError)
			return
		}
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func loadBooks(path string) ([]Book, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("could not open file %s: %w", path, err)
	}
	defer file.Close()
	return parseBooks(file, path)
}

func parseAuthors(reader io.Reader, sourceName string) (map[string]Author, error) {
	r := csv.NewReader(reader)
	r.Comma = ';'

	header, err := r.Read()
	if err != nil {
		if err == io.EOF {
			return map[string]Author{}, nil // Empty file is fine
		}
		return nil, fmt.Errorf("could not read header from %s: %w", sourceName, err)
	}
	// Handle UTF-8 BOM
	header[0] = strings.TrimPrefix(header[0], "\ufeff")

	authors := make(map[string]Author)
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("error reading record from %s: %w", sourceName, err)
		}
		if len(record) < 3 {
			log.Printf("Skipping invalid author record in %s: %v", sourceName, record)
			continue
		}
		author := Author{
			Email:     record[0],
			FirstName: record[1],
			LastName:  record[2],
		}
		authors[author.Email] = author
	}

	return authors, nil
}

func loadAuthors(path string) (map[string]Author, error) {
	file, err := os.Open(path)
	if err != nil {
		// It's ok if the file doesn't exist yet
		if os.IsNotExist(err) {
			return make(map[string]Author), nil
		}
		return nil, fmt.Errorf("could not open file %s: %w", path, err)
	}
	defer file.Close()
	return parseAuthors(file, path)
}

func parseBooks(reader io.Reader, sourceName string) ([]Book, error) {
	r := csv.NewReader(reader)
	r.Comma = ';'

	// Read header row
	if _, err := r.Read(); err != nil {
		if err == io.EOF {
			return []Book{}, nil // Empty file is fine
		}
		return nil, fmt.Errorf("could not read header from %s: %w", sourceName, err)
	}
	var books []Book
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("error reading record from %s: %w", sourceName, err)
		}
		if len(record) < 4 {
			log.Printf("Skipping invalid book record in %s: %v", sourceName, record)
			continue
		}

		books = append(books, Book{
			Title:       record[0],
			ISBN:        record[1],
			Authors:     strings.Split(record[2], ","),
			Description: record[3],
		})
	}

	return books, nil
}

func loadMagazines(path string) ([]Magazine, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("could not open file %s: %w", path, err)
	}
	defer file.Close()
	return parseMagazines(file, path)
}

func parseMagazines(reader io.Reader, sourceName string) ([]Magazine, error) {
	r := csv.NewReader(reader)
	r.Comma = ';'

	// Read header row
	if _, err := r.Read(); err != nil {
		if err == io.EOF {
			return []Magazine{}, nil // Empty file is fine
		}
		return nil, fmt.Errorf("could not read header from %s: %w", sourceName, err)
	}
	var magazines []Magazine
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("error reading record from %s: %w", sourceName, err)
		}
		if len(record) < 4 {
			log.Printf("Skipping invalid magazine record in %s: %v", sourceName, record)
			continue
		}

		magazines = append(magazines, Magazine{
			Title:       record[0],
			ISBN:        record[1],
			Authors:     strings.Split(record[2], ","),
			PublishedAt: record[3],
		})
	}

	return magazines, nil
}

func appendRecord(path string, record []string) error {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	writer.Comma = ';'
	defer writer.Flush()

	return writer.Write(record)
}

func findItemByISBN(isbn string, books []Book, magazines []Magazine) LibraryItem {
	for _, book := range books {
		if book.ISBN == isbn {
			return book
		}
	}
	for _, magazine := range magazines {
		if magazine.ISBN == isbn {
			return magazine
		}
	}
	return nil
}

func findItemsByAuthor(authorEmail string, books []Book, magazines []Magazine) []LibraryItem {
	var foundItems []LibraryItem
	for _, book := range books {
		for _, author := range book.Authors {
			if author == authorEmail {
				foundItems = append(foundItems, book)
				break
			}
		}
	}

	for _, magazine := range magazines {
		for _, author := range magazine.Authors {
			if author == authorEmail {
				foundItems = append(foundItems, magazine)
				break
			}
		}
	}
	return foundItems
}

func getSortedItems(books []Book, magazines []Magazine, sortOrder string) ([]LibraryItem, error) {
	if sortOrder != "title-asc" && sortOrder != "title-desc" {
		return nil, fmt.Errorf("invalid sort order: %s. Use 'title-asc' or 'title-desc'", sortOrder)
	}

	var allItems []LibraryItem
	for _, book := range books {
		allItems = append(allItems, book)
	}
	for _, magazine := range magazines {
		allItems = append(allItems, magazine)
	}

	sort.Slice(allItems, func(i, j int) bool {
		if sortOrder == "title-desc" {
			return allItems[i].GetTitle() > allItems[j].GetTitle()
		}
		return allItems[i].GetTitle() < allItems[j].GetTitle() // Default to ascending
	})

	return allItems, nil
}
