package main

import (
	"bufio"
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

type Author struct {
	Email     string
	FirstName string
	LastName  string
}

type Book struct {
	Title       string
	ISBN        string
	Authors     []string // emails
	Description string
}

type Magazine struct {
	Title       string
	ISBN        string
	Authors     []string // emails
	PublishedAt string
}

type Item struct {
	Title       string
	ISBN        string
	Authors     []string
	Type        string // "book" or "magazine"
	Description string // for books
	PublishedAt string // for magazines
}

type Library struct {
	Authors    map[string]Author
	Books      []Book
	Magazines  []Magazine
}

func loadAuthors(filename string) (map[string]Author, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return loadAuthorsFromReader(file)
}

func loadAuthorsFromReader(r io.Reader) (map[string]Author, error) {
	reader := csv.NewReader(r)
	reader.Comma = ';'
	reader.FieldsPerRecord = -1 // allow variable fields

	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	authors := make(map[string]Author)
	for i, record := range records {
		if i == 0 { // skip header
			continue
		}
		if len(record) < 3 {
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

func loadBooks(filename string) ([]Book, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return loadBooksFromReader(file)
}

func loadBooksFromReader(r io.Reader) ([]Book, error) {
	reader := csv.NewReader(r)
	reader.Comma = ';'
	reader.FieldsPerRecord = -1

	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var books []Book
	for i, record := range records {
		if i == 0 { // skip header
			continue
		}
		if len(record) < 4 {
			continue
		}
		authors := strings.Split(record[2], ",")
		for j, a := range authors {
			authors[j] = strings.TrimSpace(a)
		}
		book := Book{
			Title:       record[0],
			ISBN:        record[1],
			Authors:     authors,
			Description: record[3],
		}
		books = append(books, book)
	}
	return books, nil
}

func loadMagazines(filename string) ([]Magazine, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return loadMagazinesFromReader(file)
}

func loadMagazinesFromReader(r io.Reader) ([]Magazine, error) {
	reader := csv.NewReader(r)
	reader.Comma = ';'
	reader.FieldsPerRecord = -1

	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var magazines []Magazine
	for i, record := range records {
		if i == 0 { // skip header
			continue
		}
		if len(record) < 4 {
			continue
		}
		authors := strings.Split(record[2], ",")
		for j, a := range authors {
			authors[j] = strings.TrimSpace(a)
		}
		magazine := Magazine{
			Title:       record[0],
			ISBN:        record[1],
			Authors:     authors,
			PublishedAt: record[3],
		}
		magazines = append(magazines, magazine)
	}
	return magazines, nil
}

func (lib *Library) getAuthorsNames(emails []string) string {
	var names []string
	for _, email := range emails {
		names = append(names, lib.getAuthorName(email))
	}
	return strings.Join(names, ", ")
}

func (lib *Library) getAuthorName(email string) string {
	if author, ok := lib.Authors[email]; ok {
		return author.FirstName + " " + author.LastName
	}
	return email
}

func appendAuthor(filename string, author Author) error {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	writer.Comma = ';'
	defer writer.Flush()

	return writer.Write([]string{author.Email, author.FirstName, author.LastName})
}

func appendBook(filename string, book Book) error {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	writer.Comma = ';'
	defer writer.Flush()

	authors := strings.Join(book.Authors, ",")
	return writer.Write([]string{book.Title, book.ISBN, authors, book.Description})
}

func appendMagazine(filename string, mag Magazine) error {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	writer.Comma = ';'
	defer writer.Flush()

	authors := strings.Join(mag.Authors, ",")
	return writer.Write([]string{mag.Title, mag.ISBN, authors, mag.PublishedAt})
}

func cliMode(lib *Library) {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Println("\nLibrary CLI")
		fmt.Println("1. Display all items")
		fmt.Println("2. Search by ISBN")
		fmt.Println("3. Search by author email")
		fmt.Println("4. Add new item")
		fmt.Println("5. Exit")
		fmt.Print("Choose an option: ")

		if !scanner.Scan() {
			break
		}
		choice := strings.TrimSpace(scanner.Text())

		switch choice {
		case "1":
			items := lib.getItems("", "")
			printItems(items, lib)
		case "2":
			fmt.Print("Enter ISBN: ")
			if scanner.Scan() {
				isbn := strings.TrimSpace(scanner.Text())
				items := lib.getItems(isbn, "")
				printItems(items, lib)
			}
		case "3":
			fmt.Print("Enter author email: ")
			if scanner.Scan() {
				email := strings.TrimSpace(scanner.Text())
				items := lib.getItems("", email)
				printItems(items, lib)
			}
		case "4":
			addItemCLI(lib, scanner)
		case "5":
			return
		default:
			fmt.Println("Invalid option")
		}
	}
}

func printItems(items []Item, lib *Library) {
	if len(items) == 0 {
		fmt.Println("No items found.")
		return
	}
	for _, item := range items {
		fmt.Printf("Title: %s\n", item.Title)
		fmt.Printf("ISBN: %s\n", item.ISBN)
		fmt.Printf("Authors: %s\n", getAuthorsNamesCLI(item.Authors, lib))
		if item.Type == "book" {
			fmt.Printf("Description: %s\n", item.Description)
		} else {
			fmt.Printf("Published: %s\n", item.PublishedAt)
		}
		fmt.Println("---")
	}
}

func getAuthorsNamesCLI(emails []string, lib *Library) string {
	var names []string
	for _, email := range emails {
		names = append(names, lib.getAuthorName(email))
	}
	return strings.Join(names, ", ")
}

func addItemCLI(lib *Library, scanner *bufio.Scanner) {
	fmt.Print("Type (book/magazine): ")
	if !scanner.Scan() {
		return
	}
	itemType := strings.TrimSpace(scanner.Text())

	fmt.Print("Title: ")
	if !scanner.Scan() {
		return
	}
	title := strings.TrimSpace(scanner.Text())

	fmt.Print("ISBN: ")
	if !scanner.Scan() {
		return
	}
	isbn := strings.TrimSpace(scanner.Text())

	fmt.Print("Author Email: ")
	if !scanner.Scan() {
		return
	}
	authorEmail := strings.TrimSpace(scanner.Text())

	var authorFirst, authorLast string
	if _, exists := lib.Authors[authorEmail]; !exists {
		fmt.Print("Author First Name: ")
		if !scanner.Scan() {
			return
		}
		authorFirst = strings.TrimSpace(scanner.Text())

		fmt.Print("Author Last Name: ")
		if !scanner.Scan() {
			return
		}
		authorLast = strings.TrimSpace(scanner.Text())

		if authorFirst == "" || authorLast == "" {
			fmt.Println("Author first and last name required for new author")
			return
		}
	}

	var description, publishedAt string
	if itemType == "book" {
		fmt.Print("Description: ")
		if !scanner.Scan() {
			return
		}
		description = strings.TrimSpace(scanner.Text())
		if description == "" {
			fmt.Println("Description required for book")
			return
		}
	} else if itemType == "magazine" {
		fmt.Print("Published At: ")
		if !scanner.Scan() {
			return
		}
		publishedAt = strings.TrimSpace(scanner.Text())
		if publishedAt == "" {
			fmt.Println("Published date required for magazine")
			return
		}
	} else {
		fmt.Println("Invalid type")
		return
	}

	// Add author if new
	if _, exists := lib.Authors[authorEmail]; !exists {
		author := Author{Email: authorEmail, FirstName: authorFirst, LastName: authorLast}
		lib.Authors[authorEmail] = author
		if err := appendAuthor("resources/authors.csv", author); err != nil {
			fmt.Println("Failed to save author:", err)
			return
		}
	}

	// Add item
	if itemType == "book" {
		book := Book{Title: title, ISBN: isbn, Authors: []string{authorEmail}, Description: description}
		lib.Books = append(lib.Books, book)
		if err := appendBook("resources/books.csv", book); err != nil {
			fmt.Println("Failed to save book:", err)
			return
		}
	} else {
		mag := Magazine{Title: title, ISBN: isbn, Authors: []string{authorEmail}, PublishedAt: publishedAt}
		lib.Magazines = append(lib.Magazines, mag)
		if err := appendMagazine("resources/magazines.csv", mag); err != nil {
			fmt.Println("Failed to save magazine:", err)
			return
		}
	}

	fmt.Println("Item added successfully!")
}

func (lib *Library) handleAdd(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	itemType := r.FormValue("type")
	title := r.FormValue("title")
	isbn := r.FormValue("isbn")
	authorEmail := r.FormValue("author_email")
	authorFirst := r.FormValue("author_first")
	authorLast := r.FormValue("author_last")
	description := r.FormValue("description")
	publishedAt := r.FormValue("published_at")

	// Basic validation
	if itemType == "" || title == "" || isbn == "" || authorEmail == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	// Handle author
	if _, exists := lib.Authors[authorEmail]; !exists {
		if authorFirst == "" || authorLast == "" {
			http.Error(w, "New author requires first and last name", http.StatusBadRequest)
			return
		}
		author := Author{Email: authorEmail, FirstName: authorFirst, LastName: authorLast}
		lib.Authors[authorEmail] = author
		if err := appendAuthor("resources/authors.csv", author); err != nil {
			http.Error(w, "Failed to save author", http.StatusInternalServerError)
			return
		}
	}

	// Handle item
	if itemType == "book" {
		if description == "" {
			http.Error(w, "Book requires description", http.StatusBadRequest)
			return
		}
		book := Book{Title: title, ISBN: isbn, Authors: []string{authorEmail}, Description: description}
		lib.Books = append(lib.Books, book)
		if err := appendBook("resources/books.csv", book); err != nil {
			http.Error(w, "Failed to save book", http.StatusInternalServerError)
			return
		}
	} else if itemType == "magazine" {
		if publishedAt == "" {
			http.Error(w, "Magazine requires published date", http.StatusBadRequest)
			return
		}
		mag := Magazine{Title: title, ISBN: isbn, Authors: []string{authorEmail}, PublishedAt: publishedAt}
		lib.Magazines = append(lib.Magazines, mag)
		if err := appendMagazine("resources/magazines.csv", mag); err != nil {
			http.Error(w, "Failed to save magazine", http.StatusInternalServerError)
			return
		}
	} else {
		http.Error(w, "Invalid item type", http.StatusBadRequest)
		return
	}

	// Redirect back to the list
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (lib *Library) getItems(isbn, email string) []Item {
	var filteredBooks []Book
	var filteredMagazines []Magazine

	if isbn != "" {
		for _, book := range lib.Books {
			if book.ISBN == isbn {
				filteredBooks = append(filteredBooks, book)
			}
		}
		for _, mag := range lib.Magazines {
			if mag.ISBN == isbn {
				filteredMagazines = append(filteredMagazines, mag)
			}
		}
	} else if email != "" {
		for _, book := range lib.Books {
			for _, a := range book.Authors {
				if a == email {
					filteredBooks = append(filteredBooks, book)
					break
				}
			}
		}
		for _, mag := range lib.Magazines {
			for _, a := range mag.Authors {
				if a == email {
					filteredMagazines = append(filteredMagazines, mag)
					break
				}
			}
		}
	} else {
		filteredBooks = lib.Books
		filteredMagazines = lib.Magazines
	}

	// Create combined items
	var items []Item
	for _, book := range filteredBooks {
		items = append(items, Item{
			Title:       book.Title,
			ISBN:        book.ISBN,
			Authors:     book.Authors,
			Type:        "book",
			Description: book.Description,
		})
	}
	for _, mag := range filteredMagazines {
		items = append(items, Item{
			Title:       mag.Title,
			ISBN:        mag.ISBN,
			Authors:     mag.Authors,
			Type:        "magazine",
			PublishedAt: mag.PublishedAt,
		})
	}

	// Sort by title
	sort.Slice(items, func(i, j int) bool {
		return items[i].Title < items[j].Title
	})

	return items
}

func main() {
	lib := &Library{}

	var err error
	lib.Authors, err = loadAuthors("resources/authors.csv")
	if err != nil {
		log.Fatal(err)
	}

	lib.Books, err = loadBooks("resources/books.csv")
	if err != nil {
		log.Fatal(err)
	}

	lib.Magazines, err = loadMagazines("resources/magazines.csv")
	if err != nil {
		log.Fatal(err)
	}

	if len(os.Args) > 1 && os.Args[1] == "cli" {
		cliMode(lib)
	} else {
		http.HandleFunc("/", lib.indexHandler)
		log.Println("Starting server on :3005")
		log.Fatal(http.ListenAndServe(":3005", nil))
	}
}

func (lib *Library) indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		lib.handleAdd(w, r)
		return
	}

	isbn := r.URL.Query().Get("isbn")
	email := r.URL.Query().Get("email")

	var searchType string
	var searchValue string

	if isbn != "" {
		searchType = "ISBN"
		searchValue = isbn
	} else if email != "" {
		searchType = "author email"
		searchValue = email
	}

	items := lib.getItems(isbn, email)

	tmpl := `
<!DOCTYPE html>
<html>
<head>
    <title>Library</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 20px; }
        h1 { color: #333; }
        h2 { color: #666; }
        .item { margin-bottom: 20px; border: 1px solid #ddd; padding: 10px; }
        .title { font-weight: bold; }
        .isbn { font-style: italic; }
        .authors { color: #555; }
        .description { margin-top: 10px; }
        .published { color: #777; }
        form { margin-bottom: 20px; }
        input[type="text"], select, textarea { padding: 5px; width: 200px; margin: 5px 0; }
        textarea { width: 300px; height: 100px; }
        button { padding: 5px 10px; }
        .add-form { border: 1px solid #ccc; padding: 10px; margin-top: 20px; }
    </style>
</head>
<body>
    <h1>Library Collection</h1>

    <form action="/" method="get">
        <input type="text" name="isbn" placeholder="Enter ISBN to search" value="{{.SearchISBN}}">
        <input type="text" name="email" placeholder="Enter author email to search" value="{{.SearchEmail}}">
        <button type="submit">Search</button>
        <a href="/">Show All</a>
    </form>

    {{if .SearchType}}
    <p>Searching for {{.SearchType}}: {{.SearchValue}}</p>
    {{end}}

    <h2>Library Items (sorted by title)</h2>
    {{range .Items}}
    <div class="item">
        <div class="title">{{.Title}}</div>
        <div class="isbn">ISBN: {{.ISBN}}</div>
        <div class="authors">Authors: {{getAuthorsNames .Authors}}</div>
        {{if eq .Type "book"}}
        <div class="description">{{.Description}}</div>
        {{else}}
        <div class="published">Published: {{.PublishedAt}}</div>
        {{end}}
    </div>
    {{else}}
    <p>No items found.</p>
    {{end}}

    <div class="add-form">
        <h2>Add New Item</h2>
        <form action="/" method="post">
            <select name="type" required>
                <option value="">Select Type</option>
                <option value="book">Book</option>
                <option value="magazine">Magazine</option>
            </select><br>
            <input type="text" name="title" placeholder="Title" required><br>
            <input type="text" name="isbn" placeholder="ISBN" required><br>
            <input type="email" name="author_email" placeholder="Author Email" required><br>
            <input type="text" name="author_first" placeholder="Author First Name (if new)"><br>
            <input type="text" name="author_last" placeholder="Author Last Name (if new)"><br>
            <textarea name="description" placeholder="Description (for books)"></textarea><br>
            <input type="text" name="published_at" placeholder="Published At (for magazines, e.g. 01.01.2020)"><br>
            <button type="submit">Add Item</button>
        </form>
    </div>
</body>
</html>
`

	t, err := template.New("index").Funcs(template.FuncMap{"getAuthorsNames": lib.getAuthorsNames}).Parse(tmpl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := struct {
		Items       []Item
		SearchType  string
		SearchValue string
		SearchISBN  string
		SearchEmail string
	}{
		Items:       items,
		SearchType:  searchType,
		SearchValue: searchValue,
		SearchISBN:  isbn,
		SearchEmail: email,
	}

	err = t.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
