package main

import (
	"encoding/csv"
	"html/template"
	"log"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/books", booksHandler)
	http.HandleFunc("/magazines", magazinesHandler)
	http.HandleFunc("/all-sorted", allSortedHandler)
	http.HandleFunc("/add-item", addItemHandler)
	log.Println("Server running at http://localhost:8080/")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// Handler for add item form (GET and POST)
func addItemHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tmpl, err := template.ParseFiles("templates/add_item.html")
		if err != nil {
			http.Error(w, "Template error", http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, nil)
		return
	}

	r.ParseForm()

	// Step 2: If new_author_email fields are present, this is the second POST
	newAuthorEmails := r.Form["new_author_email"]
	if len(newAuthorEmails) > 0 {
		for _, email := range newAuthorEmails {
			first := r.FormValue("firstname_" + email)
			last := r.FormValue("lastname_" + email)
			appendToCSV("resources/authors.csv", []string{email, first, last})
		}
		// Restore previous form data for book/magazine
		itemType := r.FormValue("type")
		title := r.FormValue("title")
		isbn := r.FormValue("isbn")
		authors := r.FormValue("authors")
		description := r.FormValue("description")
		publishedAt := r.FormValue("publishedAt")
		if itemType == "book" {
			err := appendToCSV("resources/books.csv", []string{title, isbn, authors, description})
			if err != nil {
				http.Error(w, "Failed to add book", http.StatusInternalServerError)
				return
			}
		} else if itemType == "magazine" {
			err := appendToCSV("resources/magazines.csv", []string{title, isbn, authors, publishedAt})
			if err != nil {
				http.Error(w, "Failed to add magazine", http.StatusInternalServerError)
				return
			}
		} else {
			http.Error(w, "Invalid type", http.StatusBadRequest)
			return
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Step 1: First POST, check for new authors
	itemType := r.FormValue("type")
	title := r.FormValue("title")
	isbn := r.FormValue("isbn")
	authors := r.FormValue("authors")
	description := r.FormValue("description")
	publishedAt := r.FormValue("publishedAt")

	// Load existing author emails
	existing := map[string]bool{}
	file, err := openAuthorsCSV()
	if err == nil {
		records, _ := readAllAuthorsCSV(file)
		for _, rec := range records[1:] {
			if len(rec) > 0 {
				existing[rec[0]] = true // email is first column
			}
		}
		file.Close()
	}
	var newAuthors []string
	for _, email := range splitAuthors(authors) {
		email = trim(email)
		if email == "" || existing[email] {
			continue
		}
		newAuthors = append(newAuthors, email)
	}
	if len(newAuthors) > 0 {
		// Render form to collect names for new authors
		tmpl, err := template.ParseFiles("templates/add_author_names.html")
		if err != nil {
			http.Error(w, "Template error", http.StatusInternalServerError)
			return
		}
		// Pass previous form data as hidden fields
		prevData := map[string]string{
			"type":        itemType,
			"title":       title,
			"isbn":        isbn,
			"authors":     authors,
			"description": description,
			"publishedAt": publishedAt,
		}
		data := struct {
			NewAuthors []string
			PrevData   map[string]string
		}{
			NewAuthors: newAuthors,
			PrevData:   prevData,
		}
		tmpl.Execute(w, data)
		return
	}

	// All authors exist, proceed to add book or magazine
	if itemType == "book" {
		err := appendToCSV("resources/books.csv", []string{title, isbn, authors, description})
		if err != nil {
			http.Error(w, "Failed to add book", http.StatusInternalServerError)
			return
		}
	} else if itemType == "magazine" {
		err := appendToCSV("resources/magazines.csv", []string{title, isbn, authors, publishedAt})
		if err != nil {
			http.Error(w, "Failed to add magazine", http.StatusInternalServerError)
			return
		}
	} else {
		http.Error(w, "Invalid type", http.StatusBadRequest)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Helper for reading authors.csv (delegates to library.go)
func openAuthorsCSV() (*os.File, error) {
	return os.Open("resources/authors.csv")
}
func readAllAuthorsCSV(file *os.File) ([][]string, error) {
	reader := csv.NewReader(file)
	reader.Comma = ';'
	return reader.ReadAll()
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/home.html")
	if err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

func booksHandler(w http.ResponseWriter, r *http.Request) {
	books, err := readBooksCSV("resources/books.csv")
	if err != nil {
		http.Error(w, "Failed to load books", http.StatusInternalServerError)
		return
	}
	isbn := r.URL.Query().Get("isbn")
	author := r.URL.Query().Get("author")
	var filtered []Book
	if isbn != "" {
		for _, b := range books {
			if b.ISBN == isbn {
				filtered = append(filtered, b)
			}
		}
	} else if author != "" {
		for _, b := range books {
			if b.Authors != "" && containsAuthor(b.Authors, author) {
				filtered = append(filtered, b)
			}
		}
	} else {
		filtered = books
	}
	tmpl, err := template.ParseFiles("templates/books.html")
	if err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, filtered)
}

func magazinesHandler(w http.ResponseWriter, r *http.Request) {
	magazines, err := readMagazinesCSV("resources/magazines.csv")
	if err != nil {
		http.Error(w, "Failed to load magazines", http.StatusInternalServerError)
		return
	}
	isbn := r.URL.Query().Get("isbn")
	author := r.URL.Query().Get("author")
	var filtered []Magazine
	if isbn != "" {
		for _, m := range magazines {
			if m.ISBN == isbn {
				filtered = append(filtered, m)
			}
		}
	} else if author != "" {
		for _, m := range magazines {
			if m.Authors != "" && containsAuthor(m.Authors, author) {
				filtered = append(filtered, m)
			}
		}
	} else {
		filtered = magazines
	}
	tmpl, err := template.ParseFiles("templates/magazines.html")
	if err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, filtered)
}

func allSortedHandler(w http.ResponseWriter, r *http.Request) {
	books, err := readBooksCSV("resources/books.csv")
	if err != nil {
		http.Error(w, "Failed to load books", http.StatusInternalServerError)
		return
	}
	magazines, err := readMagazinesCSV("resources/magazines.csv")
	if err != nil {
		http.Error(w, "Failed to load magazines", http.StatusInternalServerError)
		return
	}
	items := combineBooksAndMagazines(books, magazines)
	dir := r.URL.Query().Get("dir")
	sortCombinedItemsByTitle(items, dir != "desc")
	tmpl, err := template.ParseFiles("templates/all_sorted.html")
	if err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, items)
}

func welcomeMessage() string {
	return "Hello world!"
}
