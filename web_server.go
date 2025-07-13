package main

import (
	"html/template"
	"log"
	"net/http"
	"sort"
)

type PageData struct {
	Books         []Book
	Magazines     []Magazine
	CombinedItems []LibraryItem
	AuthorsMap    map[string]Author
	BookCount     int
	MagazineCount int
	AuthorCount   int
	SearchISBN    string
	SearchAuthor  string
	ViewMode      string // "separate" or "combined"
	SortDir       string // "asc" or "desc"
}

func combineAndSortItems(books []Book, magazines []Magazine, sortDir string) []LibraryItem {
	var items []LibraryItem
	
	// Convert books to LibraryItems
	for _, book := range books {
		item := LibraryItem{
			Title:       book.Title,
			ISBN:        book.ISBN,
			Authors:     book.Authors,
			Type:        "book",
			Description: book.Description,
		}
		items = append(items, item)
	}
	
	// Convert magazines to LibraryItems
	for _, magazine := range magazines {
		item := LibraryItem{
			Title:       magazine.Title,
			ISBN:        magazine.ISBN,
			Authors:     magazine.Authors,
			Type:        "magazine",
			PublishedAt: magazine.PublishedAt,
		}
		items = append(items, item)
	}
	
	// Sort by title
	sort.Slice(items, func(i, j int) bool {
		if sortDir == "desc" {
			return items[i].Title > items[j].Title
		}
		return items[i].Title < items[j].Title
	})
	
	return items
}

func startWebServer(books []Book, magazines []Magazine, authors map[string]Author) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.New("library").Parse(htmlTemplate)
		if err != nil {
			http.Error(w, "Failed to parse template", http.StatusInternalServerError)
			return
		}

		// Get search queries and view mode
		searchISBN := r.URL.Query().Get("isbn")
		searchAuthor := r.URL.Query().Get("author")
		viewMode := r.URL.Query().Get("view")
		if viewMode == "" {
			viewMode = "separate"
		}
		sortDir := r.URL.Query().Get("sort")
		if sortDir == "" {
			sortDir = "asc"
		}
		
		// Filter books and magazines if search query exists
		displayBooks := books
		displayMagazines := magazines
		
		if searchISBN != "" {
			displayBooks = FilterBooksByISBN(books, searchISBN)
			displayMagazines = FilterMagazinesByISBN(magazines, searchISBN)
		} else if searchAuthor != "" {
			displayBooks = FilterBooksByAuthor(books, searchAuthor)
			displayMagazines = FilterMagazinesByAuthor(magazines, searchAuthor)
		}

		// Create combined and sorted items if in combined view mode
		var combinedItems []LibraryItem
		if viewMode == "combined" {
			combinedItems = combineAndSortItems(displayBooks, displayMagazines, sortDir)
		}

		data := PageData{
			Books:         displayBooks,
			Magazines:     displayMagazines,
			CombinedItems: combinedItems,
			AuthorsMap:    authors,
			BookCount:     len(books),
			MagazineCount: len(magazines),
			AuthorCount:   len(authors),
			SearchISBN:    searchISBN,
			SearchAuthor:  searchAuthor,
			ViewMode:      viewMode,
			SortDir:       sortDir,
		}

		err = tmpl.Execute(w, data)
		if err != nil {
			http.Error(w, "Failed to render template", http.StatusInternalServerError)
			return
		}
	})

	log.Println("Starting web server on http://localhost:8090")
	log.Fatal(http.ListenAndServe(":8090", nil))
}