package main

import (
	"encoding/csv"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"sort"
	"strings"
)

type Book struct {
	Title  string
	Author string
	ISBN   string
	Genre  string
}

type Magazine struct {
	Title       string
	Publisher   string
	IssueNumber string
	PublishedAt string
}

type Author struct {
	Email     string
	FirstName string
	LastName  string
}

var books []Book
var magazines []Magazine
var authors []Author

func loadBooks(filePath string) []Book {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		return nil
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = ';'
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Printf("Error reading CSV: %v\n", err)
		return nil
	}

	var books []Book
	for i, record := range records {
		if i == 0 {
			continue
		}
		if len(record) < 4 {
			continue
		}
		books = append(books, Book{
			Title:  strings.TrimSpace(record[0]),
			Author: strings.TrimSpace(record[2]),
			ISBN:   strings.TrimSpace(record[1]),
			Genre:  strings.TrimSpace(record[3]),
		})
	}
	return books
}

func loadMagazines(filePath string) []Magazine {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		return nil
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = ';'
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Printf("Error reading CSV: %v\n", err)
		return nil
	}

	var magazines []Magazine
	for i, record := range records {
		if i == 0 {
			continue
		}
		if len(record) < 4 {
			continue
		}
		magazines = append(magazines, Magazine{
			Title:       strings.TrimSpace(record[0]),
			Publisher:   strings.TrimSpace(record[2]),
			IssueNumber: strings.TrimSpace(record[1]),
			PublishedAt: strings.TrimSpace(record[3]),
		})
	}
	return magazines
}

func loadAuthors(filePath string) []Author {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		return nil
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = ';'
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Printf("Error reading CSV: %v\n", err)
		return nil
	}

	var authors []Author
	for i, record := range records {
		if i == 0 {
			continue
		}
		if len(record) < 3 {
			continue
		}
		authors = append(authors, Author{
			Email:     strings.TrimSpace(record[0]),
			FirstName: strings.TrimSpace(record[1]),
			LastName:  strings.TrimSpace(record[2]),
		})
	}
	return authors
}

func displayLibraryHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("web/templates/library.html"))
	tmpl.Execute(w, struct {
		Books     []Book
		Magazines []Magazine
		Authors   []Author
	}{books, magazines, authors})
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("query")
	var results []string

	for _, book := range books {
		if book.ISBN == query || strings.Contains(book.Author, query) {
			results = append(results, fmt.Sprintf("Book - Title: %s, Author: %s, ISBN: %s, Genre: %s", book.Title, book.Author, book.ISBN, book.Genre))
		}
	}

	for _, magazine := range magazines {
		if magazine.IssueNumber == query || strings.Contains(magazine.Publisher, query) {
			results = append(results, fmt.Sprintf("Magazine - Title: %s, Publisher: %s, ISBN: %s, Published At: %s", magazine.Title, magazine.Publisher, magazine.IssueNumber, magazine.PublishedAt))
		}
	}

	tmpl := template.Must(template.ParseFiles("web/templates/search.html"))
	tmpl.Execute(w, results)
}

func sortHandler(w http.ResponseWriter, r *http.Request) {
	type LibraryItem struct {
		Title   string
		Details string
	}

	var libraryItems []LibraryItem

	for _, book := range books {
		libraryItems = append(libraryItems, LibraryItem{
			Title:   book.Title,
			Details: fmt.Sprintf("Book - Title: %s, Author: %s, ISBN: %s, Genre: %s", book.Title, book.Author, book.ISBN, book.Genre),
		})
	}

	for _, magazine := range magazines {
		libraryItems = append(libraryItems, LibraryItem{
			Title:   magazine.Title,
			Details: fmt.Sprintf("Magazine - Title: %s, Publisher: %s, ISBN: %s, Published At: %s", magazine.Title, magazine.Publisher, magazine.IssueNumber, magazine.PublishedAt),
		})
	}

	sort.Slice(libraryItems, func(i, j int) bool {
		return strings.ToLower(libraryItems[i].Title) < strings.ToLower(libraryItems[j].Title)
	})

	tmpl := template.Must(template.ParseFiles("web/templates/sort.html"))
	tmpl.Execute(w, libraryItems)
}

func addHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		itemType := r.FormValue("type")
		title := r.FormValue("title")
		authorEmail := r.FormValue("authorEmail")
		isbn := r.FormValue("isbn")

		authors := loadAuthors("resources/authors.csv")
		var existingAuthor *Author
		for _, author := range authors {
			if author.Email == authorEmail {
				existingAuthor = &author
				break
			}
		}

		if existingAuthor == nil {
			firstName := r.FormValue("firstName")
			lastName := r.FormValue("lastName")
			if firstName == "" || lastName == "" {
				http.Error(w, "Author first and last name are required.", http.StatusBadRequest)
				return
			}
			file, _ := os.OpenFile("resources/authors.csv", os.O_APPEND|os.O_WRONLY, 0644)
			defer file.Close()
			file.WriteString(fmt.Sprintf("%s;%s;%s\n", authorEmail, firstName, lastName))
		} else {
			authorEmail = existingAuthor.Email
		}

		if itemType == "book" {
			genre := r.FormValue("genre")
			file, _ := os.OpenFile("resources/books.csv", os.O_APPEND|os.O_WRONLY, 0644)
			defer file.Close()
			file.WriteString(fmt.Sprintf("%s;%s;%s;%s\n", title, isbn, authorEmail, genre))
			books = loadBooks("resources/books.csv")
		} else if itemType == "magazine" {
			publishedDate := r.FormValue("publishedDate")
			file, _ := os.OpenFile("resources/magazines.csv", os.O_APPEND|os.O_WRONLY, 0644)
			defer file.Close()
			file.WriteString(fmt.Sprintf("%s;%s;%s;%s\n", title, isbn, authorEmail, publishedDate))
			magazines = loadMagazines("resources/magazines.csv")
		} else {
			http.Error(w, "Invalid type. Must be 'book' or 'magazine'.", http.StatusBadRequest)
			return
		}
	}

	tmpl := template.Must(template.ParseFiles("web/templates/add.html"))
	tmpl.Execute(w, nil)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("web/templates/home.html"))
	tmpl.Execute(w, nil)
}

func main() {
	books = loadBooks("resources/books.csv")
	magazines = loadMagazines("resources/magazines.csv")
	authors = loadAuthors("resources/authors.csv")

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/library", displayLibraryHandler)
	http.HandleFunc("/search", searchHandler)
	http.HandleFunc("/sort", sortHandler)
	http.HandleFunc("/add", addHandler)

	fmt.Println("Starting server on :8080...")
	http.ListenAndServe(":8080", nil)
}
