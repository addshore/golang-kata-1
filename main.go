package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
)

var library *Library

// AddItemRequest represents a request to add a book or magazine
type AddItemRequest struct {
	Type        string `json:"type"` // "book" or "magazine"
	Title       string `json:"title"`
	ISBN        string `json:"isbn"`
	Description string `json:"description"` // for books
	PublishedAt string `json:"publishedAt"` // for magazines
	Authors     []struct {
		Email     string `json:"email"`
		FirstName string `json:"firstName"`
		LastName  string `json:"lastName"`
	} `json:"authors"`
}

func main() {
	var err error
	library, err = setupLibrary()
	if err != nil {
		log.Fatal(err)
	}

	// Parse command-line flags
	cliMode := flag.Bool("cli", false, "Run in CLI mode instead of web server")
	flag.Parse()

	if *cliMode {
		startCLI(library)
	} else {
		startWebServer()
	}
}

func startWebServer() {
	http.HandleFunc("/", serveIndex)
	http.HandleFunc("/api/books", handleBooks)
	http.HandleFunc("/api/magazines", handleMagazines)
	http.HandleFunc("/api/library", handleLibrary)
	http.HandleFunc("/api/books/search", handleSearchBooks)
	http.HandleFunc("/api/magazines/search", handleSearchMagazines)
	http.HandleFunc("/api/books/search-by-author", handleSearchBooksByAuthor)
	http.HandleFunc("/api/magazines/search-by-author", handleSearchMagazinesByAuthor)
	http.HandleFunc("/api/items/sorted", handleSortedItems)
	http.HandleFunc("/api/items/add", handleAddItem)

	fmt.Println("Library server starting at http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func welcomeMessage() string {
	return "Hello world!"
}

func setupLibrary() (*Library, error) {
	lib := NewLibrary()
	if err := lib.Load(); err != nil {
		return nil, err
	}
	return lib, nil
}

func serveIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(getHTMLContent()))
}

func handleBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(library.Books)
}

func handleMagazines(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(library.Magazines)
}

func handleLibrary(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	data := map[string]interface{}{
		"books":     library.Books,
		"magazines": library.Magazines,
	}
	json.NewEncoder(w).Encode(data)
}

func handleSearchBooks(w http.ResponseWriter, r *http.Request) {
	isbn := r.URL.Query().Get("isbn")
	if isbn == "" {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(library.Books)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	results := library.SearchBooksByISBN(isbn)
	json.NewEncoder(w).Encode(results)
}

func handleSearchMagazines(w http.ResponseWriter, r *http.Request) {
	isbn := r.URL.Query().Get("isbn")
	if isbn == "" {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(library.Magazines)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	results := library.SearchMagazinesByISBN(isbn)
	json.NewEncoder(w).Encode(results)
}

func handleSearchBooksByAuthor(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	if email == "" {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(library.Books)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	results := library.SearchBooksByAuthorEmail(email)
	json.NewEncoder(w).Encode(results)
}

func handleSearchMagazinesByAuthor(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	if email == "" {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(library.Magazines)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	results := library.SearchMagazinesByAuthorEmail(email)
	json.NewEncoder(w).Encode(results)
}

func handleSortedItems(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	direction := r.URL.Query().Get("direction")

	var items []Item
	if direction == "desc" {
		items = library.GetAllItemsSortedByTitleDescending()
	} else {
		items = library.GetAllItemsSortedByTitle()
	}

	json.NewEncoder(w).Encode(items)
}

func handleAddItem(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	var req AddItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request body"})
		return
	}

	// Validate required fields
	if req.Title == "" || req.Type == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Title and Type are required"})
		return
	}

	if req.Type != "book" && req.Type != "magazine" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Type must be 'book' or 'magazine'"})
		return
	}

	// Add or update authors
	for _, author := range req.Authors {
		if author.Email == "" || author.FirstName == "" || author.LastName == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "All author fields (Email, FirstName, LastName) are required"})
			return
		}

		authorObj := Author{
			Email:     author.Email,
			FirstName: author.FirstName,
			LastName:  author.LastName,
		}
		if err := library.AddAuthor(authorObj); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "Failed to add author"})
			return
		}
	}

	// Add book or magazine
	if req.Type == "book" {
		if req.ISBN == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "ISBN is required for books"})
			return
		}

		// Get author objects
		authorsToAdd := make([]Author, len(req.Authors))
		for i, author := range req.Authors {
			authorsToAdd[i] = Author{
				Email:     author.Email,
				FirstName: author.FirstName,
				LastName:  author.LastName,
			}
		}

		book := Book{
			Title:       req.Title,
			ISBN:        req.ISBN,
			Description: req.Description,
			Authors:     authorsToAdd,
		}

		if err := library.AddBook(book); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "Failed to add book"})
			return
		}
	} else {
		if req.ISBN == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "ISBN is required for magazines"})
			return
		}

		// Get author objects
		authorsToAdd := make([]Author, len(req.Authors))
		for i, author := range req.Authors {
			authorsToAdd[i] = Author{
				Email:     author.Email,
				FirstName: author.FirstName,
				LastName:  author.LastName,
			}
		}

		magazine := Magazine{
			Title:       req.Title,
			ISBN:        req.ISBN,
			PublishedAt: req.PublishedAt,
			Authors:     authorsToAdd,
		}

		if err := library.AddMagazine(magazine); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "Failed to add magazine"})
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Item added successfully"})
}

func getHTMLContent() string {
	return `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Library System</title>
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }

        body {
            font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            min-height: 100vh;
            padding: 20px;
        }

        .container {
            max-width: 1200px;
            margin: 0 auto;
        }

        header {
            background: white;
            padding: 30px;
            border-radius: 10px;
            box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
            margin-bottom: 30px;
            text-align: center;
        }

        header h1 {
            color: #667eea;
            margin-bottom: 10px;
        }

        header p {
            color: #666;
            font-size: 16px;
        }

        .tabs {
            display: flex;
            gap: 10px;
            margin-bottom: 20px;
            flex-wrap: wrap;
        }

        .tab-button {
            padding: 12px 24px;
            border: none;
            border-radius: 6px;
            background: white;
            color: #667eea;
            font-size: 16px;
            font-weight: 500;
            cursor: pointer;
            transition: all 0.3s ease;
        }

        .tab-button.active {
            background: #667eea;
            color: white;
            box-shadow: 0 4px 6px rgba(102, 126, 234, 0.4);
        }

        .tab-button:hover {
            transform: translateY(-2px);
            box-shadow: 0 4px 8px rgba(0, 0, 0, 0.15);
        }

        .tab-content {
            display: none;
        }

        .tab-content.active {
            display: block;
        }

        .items-grid {
            display: grid;
            grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
            gap: 20px;
        }

        .item-card {
            background: white;
            border-radius: 10px;
            padding: 20px;
            box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
            transition: all 0.3s ease;
        }

        .item-card:hover {
            transform: translateY(-5px);
            box-shadow: 0 8px 15px rgba(0, 0, 0, 0.2);
        }

        .item-card h3 {
            color: #667eea;
            margin-bottom: 12px;
            font-size: 18px;
            line-height: 1.4;
        }

        .item-card .isbn {
            color: #999;
            font-size: 12px;
            margin-bottom: 10px;
            font-family: monospace;
        }

        .item-card .description {
            color: #555;
            font-size: 14px;
            line-height: 1.6;
            margin-bottom: 15px;
        }

        .item-card .published-at {
            color: #999;
            font-size: 12px;
            margin-bottom: 15px;
        }

        .authors-section {
            margin-top: 15px;
            padding-top: 15px;
            border-top: 1px solid #eee;
        }

        .authors-label {
            color: #999;
            font-size: 12px;
            font-weight: 500;
            text-transform: uppercase;
            margin-bottom: 8px;
        }

        .author {
            display: inline-block;
            background: #f5f5f5;
            padding: 6px 10px;
            border-radius: 20px;
            font-size: 13px;
            color: #555;
            margin-right: 8px;
            margin-bottom: 6px;
        }

        .stats {
            background: white;
            padding: 20px;
            border-radius: 10px;
            margin-bottom: 20px;
            box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
        }

        .stats-grid {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(150px, 1fr));
            gap: 20px;
        }

        .stat-item {
            text-align: center;
        }

        .stat-number {
            font-size: 32px;
            font-weight: bold;
            color: #667eea;
        }

        .stat-label {
            color: #999;
            font-size: 14px;
            margin-top: 5px;
        }

        .loading {
            text-align: center;
            color: white;
            font-size: 18px;
        }

        .error {
            background: #ff6b6b;
            color: white;
            padding: 15px;
            border-radius: 6px;
            margin-bottom: 20px;
        }

        .search-box {
            background: white;
            padding: 20px;
            border-radius: 10px;
            margin-bottom: 20px;
            box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
            display: flex;
            gap: 10px;
            flex-wrap: wrap;
        }

        .search-box input {
            flex: 1;
            min-width: 200px;
            padding: 12px;
            border: 1px solid #ddd;
            border-radius: 6px;
            font-size: 14px;
        }

        .search-box button {
            padding: 12px 24px;
            background: #667eea;
            color: white;
            border: none;
            border-radius: 6px;
            font-size: 14px;
            font-weight: 500;
            cursor: pointer;
            transition: all 0.3s ease;
        }

        .search-box button:hover {
            background: #764ba2;
            transform: translateY(-2px);
            box-shadow: 0 4px 8px rgba(0, 0, 0, 0.15);
        }

        .search-box select {
            flex: 0 1 auto;
            padding: 12px;
            border: 1px solid #ddd;
            border-radius: 6px;
            font-size: 14px;
            background: white;
            cursor: pointer;
            min-width: 200px;
        }

        .search-box select:hover {
            border-color: #667eea;
        }

        .search-results-info {
            background: white;
            padding: 15px 20px;
            border-radius: 10px;
            margin-bottom: 20px;
            color: #666;
            font-size: 14px;
            box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
        }

        .items-list {
            background: white;
            border-radius: 10px;
            box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
            overflow: hidden;
        }

        .list-item {
            padding: 20px;
            border-bottom: 1px solid #eee;
            transition: all 0.3s ease;
        }

        .list-item:last-child {
            border-bottom: none;
        }

        .list-item:hover {
            background: #f9f9f9;
        }

        .list-item-header {
            display: flex;
            justify-content: space-between;
            align-items: start;
            margin-bottom: 10px;
        }

        .list-item-title {
            font-size: 18px;
            font-weight: 600;
            color: #667eea;
            flex: 1;
            margin-right: 10px;
        }

        .list-item-type {
            display: inline-block;
            background: #667eea;
            color: white;
            padding: 4px 12px;
            border-radius: 20px;
            font-size: 12px;
            font-weight: 600;
            text-transform: uppercase;
            white-space: nowrap;
        }

        .list-item-type.magazine {
            background: #764ba2;
        }

        .list-item-detail {
            color: #666;
            font-size: 14px;
            margin-bottom: 8px;
            display: flex;
            gap: 20px;
        }

        .list-item-detail-label {
            font-weight: 600;
            color: #999;
            min-width: 80px;
        }

        .list-item-description {
            color: #555;
            font-size: 14px;
            line-height: 1.6;
            margin-bottom: 12px;
        }

        .list-item-authors {
            display: flex;
            flex-wrap: wrap;
            gap: 8px;
        }

        .list-item-author {
            background: #f5f5f5;
            padding: 6px 10px;
            border-radius: 20px;
            font-size: 13px;
            color: #555;
        }

        .form-container {
            background: white;
            padding: 30px;
            border-radius: 10px;
            box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
            max-width: 600px;
            margin: 0 auto;
        }

        .form-group {
            margin-bottom: 20px;
        }

        .form-group label {
            display: block;
            margin-bottom: 8px;
            font-weight: 500;
            color: #333;
        }

        .form-group input,
        .form-group textarea,
        .form-group select {
            width: 100%;
            padding: 12px;
            border: 1px solid #ddd;
            border-radius: 6px;
            font-size: 14px;
            font-family: inherit;
        }

        .form-group textarea {
            resize: vertical;
            min-height: 100px;
        }

        .radio-group {
            display: flex;
            gap: 20px;
            margin-top: 8px;
        }

        .radio-group label {
            display: flex;
            align-items: center;
            margin-bottom: 0;
            cursor: pointer;
        }

        .radio-group input[type="radio"] {
            width: auto;
            margin-right: 8px;
            cursor: pointer;
        }

        .authors-list {
            background: #f9f9f9;
            padding: 15px;
            border-radius: 6px;
            margin-top: 10px;
        }

        .author-entry {
            background: white;
            padding: 15px;
            border-radius: 6px;
            margin-bottom: 10px;
            border: 1px solid #eee;
        }

        .author-entry .form-group {
            margin-bottom: 10px;
        }

        .author-entry .form-group:last-child {
            margin-bottom: 0;
        }

        .remove-author-btn {
            background: #ff6b6b;
            color: white;
            border: none;
            padding: 8px 12px;
            border-radius: 4px;
            cursor: pointer;
            font-size: 12px;
            margin-top: 10px;
        }

        .remove-author-btn:hover {
            background: #ff5252;
        }

        .add-author-btn {
            background: #667eea;
            color: white;
            border: none;
            padding: 10px 20px;
            border-radius: 6px;
            cursor: pointer;
            font-size: 14px;
            font-weight: 500;
            margin-top: 10px;
        }

        .add-author-btn:hover {
            background: #764ba2;
        }

        .form-buttons {
            display: flex;
            gap: 10px;
            margin-top: 30px;
        }

        .form-buttons button {
            flex: 1;
            padding: 14px;
            border: none;
            border-radius: 6px;
            font-size: 16px;
            font-weight: 500;
            cursor: pointer;
            transition: all 0.3s ease;
        }

        .submit-btn {
            background: #667eea;
            color: white;
        }

        .submit-btn:hover {
            background: #764ba2;
            transform: translateY(-2px);
            box-shadow: 0 4px 8px rgba(0, 0, 0, 0.15);
        }

        .clear-btn {
            background: #ddd;
            color: #333;
        }

        .clear-btn:hover {
            background: #ccc;
        }

        .form-message {
            padding: 15px;
            border-radius: 6px;
            margin-bottom: 20px;
            display: none;
        }

        .form-message.success {
            background: #d4edda;
            color: #155724;
            border: 1px solid #c3e6cb;
        }

        .form-message.error {
            background: #f8d7da;
            color: #721c24;
            border: 1px solid #f5c6cb;
        }
    </style>
</head>
<body>
    <div class="container">
        <header>
            <h1>📚 Library System</h1>
            <p>Browse all books and magazines in our collection</p>
        </header>

        <div class="tabs">
            <button class="tab-button active" onclick="switchTab('overview')">Overview</button>
            <button class="tab-button" onclick="switchTab('books')">Books</button>
            <button class="tab-button" onclick="switchTab('magazines')">Magazines</button>
            <button class="tab-button" onclick="switchTab('sorted')">Sorted by Title</button>
            <button class="tab-button" onclick="switchTab('addItems')">Add Items</button>
        </div>

        <div id="overview" class="tab-content active">
            <div class="stats">
                <div class="stats-grid">
                    <div class="stat-item">
                        <div class="stat-number" id="bookCount">0</div>
                        <div class="stat-label">Books</div>
                    </div>
                    <div class="stat-item">
                        <div class="stat-number" id="magazineCount">0</div>
                        <div class="stat-label">Magazines</div>
                    </div>
                    <div class="stat-item">
                        <div class="stat-number" id="totalCount">0</div>
                        <div class="stat-label">Total Items</div>
                    </div>
                </div>
            </div>
        </div>

        <div id="books" class="tab-content">
            <div class="search-box">
                <input type="text" id="booksSearchInput" placeholder="Search by ISBN..." />
                <button onclick="searchBooks()">Search ISBN</button>
                <input type="text" id="booksAuthorSearchInput" placeholder="Search by author email..." />
                <button onclick="searchBooksByAuthor()">Search Author</button>
                <button onclick="clearBooksSearch()" style="background: #999;">Clear</button>
            </div>
            <div id="booksSearchInfo" class="search-results-info" style="display: none;"></div>
            <div class="loading" id="booksLoading">Loading books...</div>
            <div id="booksContainer" class="items-grid"></div>
        </div>

        <div id="magazines" class="tab-content">
            <div class="search-box">
                <input type="text" id="magazinesSearchInput" placeholder="Search by ISBN..." />
                <button onclick="searchMagazines()">Search ISBN</button>
                <input type="text" id="magazinesAuthorSearchInput" placeholder="Search by author email..." />
                <button onclick="searchMagazinesByAuthor()">Search Author</button>
                <button onclick="clearMagazinesSearch()" style="background: #999;">Clear</button>
            </div>
            <div id="magazinesSearchInfo" class="search-results-info" style="display: none;"></div>
            <div class="loading" id="magazinesLoading">Loading magazines...</div>
            <div id="magazinesContainer" class="items-grid"></div>
        </div>

        <div id="sorted" class="tab-content">
            <div class="search-box">
                <select id="sortDirection" onchange="changeSortDirection()">
                    <option value="asc">Ascending (A-Z)</option>
                    <option value="desc">Descending (Z-A)</option>
                </select>
            </div>
            <div class="loading" id="sortedLoading">Loading sorted items...</div>
            <div id="sortedContainer"></div>
        </div>

        <div id="addItems" class="tab-content">
            <div class="form-container">
                <h2 style="margin-bottom: 20px; color: #667eea;">Add New Item</h2>
                <div id="formMessage" class="form-message"></div>

                <form id="addItemForm" onsubmit="handleAddItem(event)">
                    <div class="form-group">
                        <label>Item Type</label>
                        <div class="radio-group">
                            <label>
                                <input type="radio" name="type" value="book" checked onchange="updateFormForType(this.value)">
                                Book
                            </label>
                            <label>
                                <input type="radio" name="type" value="magazine" onchange="updateFormForType(this.value)">
                                Magazine
                            </label>
                        </div>
                    </div>

                    <div class="form-group">
                        <label for="title">Title *</label>
                        <input type="text" id="title" name="title" required placeholder="Enter book or magazine title">
                    </div>

                    <div class="form-group">
                        <label for="isbn">ISBN *</label>
                        <input type="text" id="isbn" name="isbn" required placeholder="Enter ISBN">
                    </div>

                    <div class="form-group" id="descriptionGroup">
                        <label for="description">Description</label>
                        <textarea id="description" name="description" placeholder="Enter book description"></textarea>
                    </div>

                    <div class="form-group" id="publishedAtGroup" style="display: none;">
                        <label for="publishedAt">Published Date</label>
                        <input type="text" id="publishedAt" name="publishedAt" placeholder="e.g., 2024-01-15">
                    </div>

                    <div class="form-group">
                        <label>Authors *</label>
                        <div id="authorsList" class="authors-list"></div>
                        <button type="button" class="add-author-btn" onclick="addAuthorField()">+ Add Author</button>
                    </div>

                    <div class="form-buttons">
                        <button type="submit" class="submit-btn">Add Item</button>
                        <button type="button" class="clear-btn" onclick="clearForm()">Clear</button>
                    </div>
                </form>
            </div>
        </div>
    </div>

    <script>
        let libraryData = { books: [], magazines: [] };
        let booksSearchActive = false;
        let magazinesSearchActive = false;

        async function loadLibrary() {
            try {
                const response = await fetch('/api/library');
                libraryData = await response.json();
                updateStats();
                renderBooks();
                renderMagazines();
                loadSortedItems();
            } catch (error) {
                console.error('Error loading library:', error);
                document.getElementById('booksContainer').innerHTML = '<div class="error">Error loading library data</div>';
            }
        }

        async function searchBooks() {
            const isbn = document.getElementById('booksSearchInput').value.trim();
            if (!isbn) {
                clearBooksSearch();
                return;
            }

            try {
                const response = await fetch('/api/books/search?isbn=' + encodeURIComponent(isbn));
                const results = await response.json();
                
                const infoDiv = document.getElementById('booksSearchInfo');
                infoDiv.textContent = 'Found ' + (results?.length || 0) + ' book(s) matching ISBN "' + isbn + '"';
                infoDiv.style.display = 'block';
                
                booksSearchActive = true;
                renderBooksResults(results || []);
            } catch (error) {
                console.error('Error searching books:', error);
            }
        }

        function clearBooksSearch() {
            document.getElementById('booksSearchInput').value = '';
            document.getElementById('booksAuthorSearchInput').value = '';
            document.getElementById('booksSearchInfo').style.display = 'none';
            booksSearchActive = false;
            renderBooks();
        }

        async function searchBooksByAuthor() {
            const email = document.getElementById('booksAuthorSearchInput').value.trim();
            if (!email) {
                clearBooksSearch();
                return;
            }

            try {
                const response = await fetch('/api/books/search-by-author?email=' + encodeURIComponent(email));
                const results = await response.json();
                
                const infoDiv = document.getElementById('booksSearchInfo');
                infoDiv.textContent = 'Found ' + (results?.length || 0) + ' book(s) by author "' + email + '"';
                infoDiv.style.display = 'block';
                
                booksSearchActive = true;
                renderBooksResults(results || []);
            } catch (error) {
                console.error('Error searching books by author:', error);
            }
        }

        async function searchMagazines() {
            const isbn = document.getElementById('magazinesSearchInput').value.trim();
            if (!isbn) {
                clearMagazinesSearch();
                return;
            }

            try {
                const response = await fetch('/api/magazines/search?isbn=' + encodeURIComponent(isbn));
                const results = await response.json();
                
                const infoDiv = document.getElementById('magazinesSearchInfo');
                infoDiv.textContent = 'Found ' + (results?.length || 0) + ' magazine(s) matching ISBN "' + isbn + '"';
                infoDiv.style.display = 'block';
                
                magazinesSearchActive = true;
                renderMagazinesResults(results || []);
            } catch (error) {
                console.error('Error searching magazines:', error);
            }
        }

        function clearMagazinesSearch() {
            document.getElementById('magazinesSearchInput').value = '';
            document.getElementById('magazinesAuthorSearchInput').value = '';
            document.getElementById('magazinesSearchInfo').style.display = 'none';
            magazinesSearchActive = false;
            renderMagazines();
        }

        async function searchMagazinesByAuthor() {
            const email = document.getElementById('magazinesAuthorSearchInput').value.trim();
            if (!email) {
                clearMagazinesSearch();
                return;
            }

            try {
                const response = await fetch('/api/magazines/search-by-author?email=' + encodeURIComponent(email));
                const results = await response.json();
                
                const infoDiv = document.getElementById('magazinesSearchInfo');
                infoDiv.textContent = 'Found ' + (results?.length || 0) + ' magazine(s) by author "' + email + '"';
                infoDiv.style.display = 'block';
                
                magazinesSearchActive = true;
                renderMagazinesResults(results || []);
            } catch (error) {
                console.error('Error searching magazines by author:', error);
            }
        }

        function updateStats() {
            document.getElementById('bookCount').textContent = libraryData.books?.length || 0;
            document.getElementById('magazineCount').textContent = libraryData.magazines?.length || 0;
            const total = (libraryData.books?.length || 0) + (libraryData.magazines?.length || 0);
            document.getElementById('totalCount').textContent = total;
        }

        function renderBooks() {
            const container = document.getElementById('booksContainer');
            const loading = document.getElementById('booksLoading');

            if (!libraryData.books || libraryData.books.length === 0) {
                container.innerHTML = '<p style="color: white; grid-column: 1/-1;">No books found</p>';
                loading.style.display = 'none';
                return;
            }

            loading.style.display = 'none';
            container.innerHTML = libraryData.books.map(book => createBookCard(book)).join('');
        }

        function renderBooksResults(results) {
            const container = document.getElementById('booksContainer');
            const loading = document.getElementById('booksLoading');

            if (!results || results.length === 0) {
                container.innerHTML = '<p style="color: white; grid-column: 1/-1;">No books found matching your search</p>';
                loading.style.display = 'none';
                return;
            }

            loading.style.display = 'none';
            container.innerHTML = results.map(book => createBookCard(book)).join('');
        }

        function renderMagazines() {
            const container = document.getElementById('magazinesContainer');
            const loading = document.getElementById('magazinesLoading');

            if (!libraryData.magazines || libraryData.magazines.length === 0) {
                container.innerHTML = '<p style="color: white; grid-column: 1/-1;">No magazines found</p>';
                loading.style.display = 'none';
                return;
            }

            loading.style.display = 'none';
            container.innerHTML = libraryData.magazines.map(magazine => createMagazineCard(magazine)).join('');
        }

        function renderMagazinesResults(results) {
            const container = document.getElementById('magazinesContainer');
            const loading = document.getElementById('magazinesLoading');

            if (!results || results.length === 0) {
                container.innerHTML = '<p style="color: white; grid-column: 1/-1;">No magazines found matching your search</p>';
                loading.style.display = 'none';
                return;
            }

            loading.style.display = 'none';
            container.innerHTML = results.map(magazine => createMagazineCard(magazine)).join('');
        }

        async function loadSortedItems() {
            const direction = document.getElementById('sortDirection').value;
            try {
                const response = await fetch('/api/items/sorted?direction=' + direction);
                const items = await response.json();
                renderSortedItems(items);
            } catch (error) {
                console.error('Error loading sorted items:', error);
                document.getElementById('sortedContainer').innerHTML = '<div class="error">Error loading sorted items</div>';
            }
        }

        function changeSortDirection() {
            loadSortedItems();
        }

        function renderSortedItems(items) {
            const container = document.getElementById('sortedContainer');
            const loading = document.getElementById('sortedLoading');

            if (!items || items.length === 0) {
                container.innerHTML = '<p style="color: white; margin-top: 20px;">No items found</p>';
                loading.style.display = 'none';
                return;
            }

            loading.style.display = 'none';
            const listHtml = '<div class="items-list">' + items.map(item => createItemListRow(item)).join('') + '</div>';
            container.innerHTML = listHtml;
        }

        function createItemListRow(item) {
            const authorsHtml = item.Authors && item.Authors.length > 0
                ? item.Authors.map(author => '<span class="list-item-author">' + author.FirstName + ' ' + author.LastName + ' (' + author.Email + ')</span>').join('')
                : '<span class="list-item-author">Unknown</span>';

            let detailsHtml = '';
            if (item.Type === 'book') {
                detailsHtml = '<div class="list-item-description">' + escapeHtml(item.Description) + '</div>';
            } else if (item.Type === 'magazine') {
                detailsHtml = '<div class="list-item-detail"><span class="list-item-detail-label">Published:</span><span>' + item.PublishedAt + '</span></div>';
            }

            const typeClass = item.Type === 'magazine' ? 'magazine' : '';

            return '<div class="list-item">' +
                '<div class="list-item-header">' +
                '<div class="list-item-title">' + escapeHtml(item.Title) + '</div>' +
                '<div class="list-item-type ' + typeClass + '">' + item.Type + '</div>' +
                '</div>' +
                '<div class="list-item-detail"><span class="list-item-detail-label">ISBN:</span><span>' + item.ISBN + '</span></div>' +
                detailsHtml +
                '<div class="list-item-authors">' + authorsHtml + '</div>' +
                '</div>';
        }

        function createBookCard(book) {
            const authorsHtml = book.Authors && book.Authors.length > 0
                ? book.Authors.map(author => '<span class="author">' + author.FirstName + ' ' + author.LastName + '</span>').join('')
                : '<span class="author">Unknown</span>';

            return '<div class="item-card">' +
                '<h3>' + escapeHtml(book.Title) + '</h3>' +
                '<div class="isbn">ISBN: ' + book.ISBN + '</div>' +
                '<div class="description">' + escapeHtml(book.Description) + '</div>' +
                '<div class="authors-section">' +
                '<div class="authors-label">Authors</div>' +
                '<div>' + authorsHtml + '</div>' +
                '</div>' +
                '</div>';
        }

        function createMagazineCard(magazine) {
            const authorsHtml = magazine.Authors && magazine.Authors.length > 0
                ? magazine.Authors.map(author => '<span class="author">' + author.FirstName + ' ' + author.LastName + '</span>').join('')
                : '<span class="author">Unknown</span>';

            return '<div class="item-card">' +
                '<h3>' + escapeHtml(magazine.Title) + '</h3>' +
                '<div class="isbn">ISBN: ' + magazine.ISBN + '</div>' +
                '<div class="published-at">Published: ' + magazine.PublishedAt + '</div>' +
                '<div class="authors-section">' +
                '<div class="authors-label">Contributors</div>' +
                '<div>' + authorsHtml + '</div>' +
                '</div>' +
                '</div>';
        }

        function switchTab(tabName) {
            document.querySelectorAll('.tab-content').forEach(tab => {
                tab.classList.remove('active');
            });
            document.querySelectorAll('.tab-button').forEach(btn => {
                btn.classList.remove('active');
            });

            document.getElementById(tabName).classList.add('active');
            event.target.classList.add('active');
            
            if (tabName === 'addItems') {
                initializeAddItemsForm();
            }
        }

        function initializeAddItemsForm() {
            if (document.getElementById('authorsList').children.length === 0) {
                addAuthorField();
            }
        }

        function updateFormForType(type) {
            if (type === 'book') {
                document.getElementById('descriptionGroup').style.display = 'block';
                document.getElementById('publishedAtGroup').style.display = 'none';
            } else {
                document.getElementById('descriptionGroup').style.display = 'none';
                document.getElementById('publishedAtGroup').style.display = 'block';
            }
        }

        function addAuthorField() {
            const authorsList = document.getElementById('authorsList');
            const authorIndex = authorsList.children.length;
            
            const authorHtml = '<div class="author-entry" id="author-' + authorIndex + '">' +
                '<div class="form-group">' +
                '<label for="author-email-' + authorIndex + '">Email *</label>' +
                '<input type="email" id="author-email-' + authorIndex + '" class="author-email" placeholder="author@example.com" required>' +
                '</div>' +
                '<div class="form-group">' +
                '<label for="author-firstname-' + authorIndex + '">First Name *</label>' +
                '<input type="text" id="author-firstname-' + authorIndex + '" class="author-firstname" placeholder="John" required>' +
                '</div>' +
                '<div class="form-group">' +
                '<label for="author-lastname-' + authorIndex + '">Last Name *</label>' +
                '<input type="text" id="author-lastname-' + authorIndex + '" class="author-lastname" placeholder="Doe" required>' +
                '</div>' +
                '<button type="button" class="remove-author-btn" onclick="removeAuthorField(' + authorIndex + ')">Remove Author</button>' +
                '</div>';
            
            authorsList.insertAdjacentHTML('beforeend', authorHtml);
        }

        function removeAuthorField(index) {
            const element = document.getElementById('author-' + index);
            if (element) {
                element.remove();
            }
        }

        function clearForm() {
            document.getElementById('addItemForm').reset();
            document.getElementById('authorsList').innerHTML = '';
            addAuthorField();
            document.getElementById('formMessage').style.display = 'none';
            updateFormForType('book');
        }

        async function handleAddItem(event) {
            event.preventDefault();

            const type = document.querySelector('input[name="type"]:checked').value;
            const title = document.getElementById('title').value.trim();
            const isbn = document.getElementById('isbn').value.trim();
            const description = document.getElementById('description').value.trim();
            const publishedAt = document.getElementById('publishedAt').value.trim();

            // Collect authors
            const authors = [];
            document.querySelectorAll('.author-entry').forEach(entry => {
                const email = entry.querySelector('.author-email').value.trim();
                const firstName = entry.querySelector('.author-firstname').value.trim();
                const lastName = entry.querySelector('.author-lastname').value.trim();

                if (email && firstName && lastName) {
                    authors.push({
                        email: email,
                        firstName: firstName,
                        lastName: lastName
                    });
                }
            });

            // Validate
            if (!title || !isbn) {
                showMessage('Title and ISBN are required', 'error');
                return;
            }

            if (authors.length === 0) {
                showMessage('At least one author is required', 'error');
                return;
            }

            // Prepare request
            const requestData = {
                type: type,
                title: title,
                isbn: isbn,
                description: description,
                publishedAt: publishedAt,
                authors: authors
            };

            try {
                const response = await fetch('/api/items/add', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json'
                    },
                    body: JSON.stringify(requestData)
                });

                const result = await response.json();

                if (!response.ok) {
                    showMessage(result.error || 'Failed to add item', 'error');
                    return;
                }

                showMessage('Item added successfully!', 'success');
                clearForm();
                
                // Reload the library data
                await loadLibrary();
            } catch (error) {
                console.error('Error adding item:', error);
                showMessage('Error adding item: ' + error.message, 'error');
            }
        }

        function showMessage(message, type) {
            const messageDiv = document.getElementById('formMessage');
            messageDiv.textContent = message;
            messageDiv.className = 'form-message ' + type;
            messageDiv.style.display = 'block';
        }

        function escapeHtml(text) {
            const div = document.createElement('div');
            div.textContent = text;
            return div.innerHTML;
        }

        document.addEventListener('DOMContentLoaded', function() {
            loadLibrary();
            
            // Add Enter key support for search inputs
            document.getElementById('booksSearchInput').addEventListener('keypress', function(e) {
                if (e.key === 'Enter') searchBooks();
            });
            
            document.getElementById('booksAuthorSearchInput').addEventListener('keypress', function(e) {
                if (e.key === 'Enter') searchBooksByAuthor();
            });
            
            document.getElementById('magazinesSearchInput').addEventListener('keypress', function(e) {
                if (e.key === 'Enter') searchMagazines();
            });
            
            document.getElementById('magazinesAuthorSearchInput').addEventListener('keypress', function(e) {
                if (e.key === 'Enter') searchMagazinesByAuthor();
            });
        });
    </script>
</body>
</html>`
}
