package main

import (
	"bufio"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/echocat/golang-kata-1/pkg/library"
)

func main() {
	lib, err := library.LoadLibrary("resources/authors.csv", "resources/books.csv", "resources/magazines.csv")
	if err != nil {
		log.Fatalf("Failed to load library: %v", err)
	}

	authorMap := make(map[string]library.Author)
	for _, author := range lib.Authors {
		authorMap[author.Email] = author
	}

	if len(os.Args) > 1 {
		cmd := os.Args[1]
		switch cmd {
		case "search-isbn":
			if len(os.Args) < 3 {
				log.Fatal("Usage: go run main.go search-isbn <isbn>")
			}
			isbn := os.Args[2]
			item, found := lib.FindItemByISBN(isbn)
			if !found {
				fmt.Printf("No item found with ISBN: %s\n", isbn)
			} else {
				switch v := item.(type) {
				case library.Book:
					v.Print(authorMap)
				case library.Magazine:
					v.Print(authorMap)
				}
			}
		case "search-author":
			if len(os.Args) < 3 {
				log.Fatal("Usage: go run main.go search-author <email>")
			}
			email := os.Args[2]
			items := lib.FindItemsByAuthorEmail(email)
			if len(items) == 0 {
				fmt.Printf("No items found for author email: %s\n", email)
			} else {
				for _, item := range items {
					switch v := item.(type) {
					case library.Book:
						v.Print(authorMap)
					case library.Magazine:
						v.Print(authorMap)
					}
				}
			}
		case "sort":
			items := lib.GetSortedItems()
			for _, item := range items {
				switch v := item.(type) {
				case library.Book:
					v.Print(authorMap)
				case library.Magazine:
					v.Print(authorMap)
				}
			}
		case "add":
			reader := bufio.NewReader(os.Stdin)
			addInteractive(lib, authorMap, reader)
		case "web":
			startWebServer(lib)
		default:
			fmt.Printf("Unknown command: %s\n", cmd)
		}
	} else {
		runInteractiveCLI(lib, authorMap)
	}
}

func runInteractiveCLI(lib *library.Library, authorMap map[string]library.Author) {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Println("\n--- Library Menu ---")
		fmt.Println("1. List all books and magazines")
		fmt.Println("2. Search by ISBN")
		fmt.Println("3. Search by Author Email")
		fmt.Println("4. Sort by Title")
		fmt.Println("5. Add Book/Magazine")
		fmt.Println("6. Start Web UI (this will block)")
		fmt.Println("0. Exit")
		fmt.Print("Choice: ")

		choiceStr, _ := reader.ReadString('\n')
		choice := strings.TrimSpace(choiceStr)

		switch choice {
		case "1":
			for _, book := range lib.Books {
				book.Print(authorMap)
			}
			for _, magazine := range lib.Magazines {
				magazine.Print(authorMap)
			}
		case "2":
			fmt.Print("ISBN: ")
			isbn, _ := reader.ReadString('\n')
			item, found := lib.FindItemByISBN(strings.TrimSpace(isbn))
			if !found {
				fmt.Println("Not found")
			} else {
				if b, ok := item.(library.Book); ok {
					b.Print(authorMap)
				} else if m, ok := item.(library.Magazine); ok {
					m.Print(authorMap)
				}
			}
		case "3":
			fmt.Print("Author Email: ")
			email, _ := reader.ReadString('\n')
			items := lib.FindItemsByAuthorEmail(strings.TrimSpace(email))
			for _, item := range items {
				if b, ok := item.(library.Book); ok {
					b.Print(authorMap)
				} else if m, ok := item.(library.Magazine); ok {
					m.Print(authorMap)
				}
			}
		case "4":
			items := lib.GetSortedItems()
			for _, item := range items {
				if b, ok := item.(library.Book); ok {
					b.Print(authorMap)
				} else if m, ok := item.(library.Magazine); ok {
					m.Print(authorMap)
				}
			}
		case "5":
			addInteractive(lib, authorMap, reader)
		case "6":
			startWebServer(lib)
		case "0":
			return
		default:
			fmt.Println("Invalid choice")
		}
	}
}

func addInteractive(lib *library.Library, authorMap map[string]library.Author, reader *bufio.Reader) {
	fmt.Print("Add (B)ook or (M)agazine? ")
	choice, _ := reader.ReadString('\n')
	choice = strings.TrimSpace(strings.ToUpper(choice))

	if choice != "B" && choice != "M" {
		fmt.Println("Invalid choice")
		return
	}

	fmt.Print("Title: ")
	title, _ := reader.ReadString('\n')
	title = strings.TrimSpace(title)

	fmt.Print("ISBN: ")
	isbn, _ := reader.ReadString('\n')
	isbn = strings.TrimSpace(isbn)

	fmt.Print("Authors (comma separated emails): ")
	authorsStr, _ := reader.ReadString('\n')
	authors := strings.Split(strings.TrimSpace(authorsStr), ",")
	for i, a := range authors {
		authors[i] = strings.TrimSpace(a)
	}

	for _, email := range authors {
		if _, ok := authorMap[email]; !ok {
			fmt.Printf("Author %s not found. Please provide details:\n", email)
			fmt.Print("First Name: ")
			fname, _ := reader.ReadString('\n')
			fmt.Print("Last Name: ")
			lname, _ := reader.ReadString('\n')
			author := library.Author{
				Email:     email,
				FirstName: strings.TrimSpace(fname),
				LastName:  strings.TrimSpace(lname),
			}
			lib.Authors = append(lib.Authors, author)
			authorMap[email] = author
		}
	}

	if choice == "B" {
		fmt.Print("Description: ")
		desc, _ := reader.ReadString('\n')
		lib.Books = append(lib.Books, library.Book{
			Title:       title,
			ISBN:        isbn,
			Authors:     authors,
			Description: strings.TrimSpace(desc),
		})
	} else {
		fmt.Print("Published At (DD.MM.YYYY): ")
		pub, _ := reader.ReadString('\n')
		lib.Magazines = append(lib.Magazines, library.Magazine{
			Title:       title,
			ISBN:        isbn,
			Authors:     authors,
			PublishedAt: strings.TrimSpace(pub),
		})
	}

	err := lib.Save("resources/authors.csv", "resources/books.csv", "resources/magazines.csv")
	if err != nil {
		fmt.Printf("Failed to save library: %v\n", err)
	} else {
		fmt.Println("Added successfully and saved to CSV.")
	}
}

type ViewItem struct {
	Type             string
	Title            string
	ISBN             string
	AuthorsFormatted string
	Details          string
}

func startWebServer(lib *library.Library) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		authorMap := make(map[string]library.Author)
		for _, author := range lib.Authors {
			authorMap[author.Email] = author
		}

		var rawItems []interface{}
		isbn := r.URL.Query().Get("isbn")
		email := r.URL.Query().Get("email")
		sortOrder := r.URL.Query().Get("sort")

		if isbn != "" {
			if item, found := lib.FindItemByISBN(isbn); found {
				rawItems = append(rawItems, item)
			}
		} else if email != "" {
			rawItems = lib.FindItemsByAuthorEmail(email)
		} else if sortOrder == "title" {
			rawItems = lib.GetSortedItems()
		} else {
			for _, b := range lib.Books {
				rawItems = append(rawItems, b)
			}
			for _, m := range lib.Magazines {
				rawItems = append(rawItems, m)
			}
		}

		var viewItems []ViewItem
		for _, item := range rawItems {
			switch v := item.(type) {
			case library.Book:
				viewItems = append(viewItems, ViewItem{
					Type:             "Book",
					Title:            v.Title,
					ISBN:             v.ISBN,
					AuthorsFormatted: library.FormatAuthors(v.Authors, authorMap),
					Details:          v.Description,
				})
			case library.Magazine:
				viewItems = append(viewItems, ViewItem{
					Type:             "Magazine",
					Title:            v.Title,
					ISBN:             v.ISBN,
					AuthorsFormatted: library.FormatAuthors(v.Authors, authorMap),
					Details:          "Published at: " + v.PublishedAt,
				})
			}
		}

		tmpl, err := template.ParseFiles("templates/index.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, struct{ Items []ViewItem }{Items: viewItems})
	})

	http.HandleFunc("/add", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		r.ParseForm()
		itemType := r.FormValue("type")
		title := r.FormValue("title")
		isbn := r.FormValue("isbn")
		authorsStr := r.FormValue("authors")
		authors := strings.Split(authorsStr, ",")
		for i, a := range authors {
			authors[i] = strings.TrimSpace(a)
		}

		// Handle new author
		newEmail := r.FormValue("new_author_email")
		if newEmail != "" {
			lib.Authors = append(lib.Authors, library.Author{
				Email:     newEmail,
				FirstName: r.FormValue("new_author_firstname"),
				LastName:  r.FormValue("new_author_lastname"),
			})
		}

		if itemType == "book" {
			lib.Books = append(lib.Books, library.Book{
				Title:       title,
				ISBN:        isbn,
				Authors:     authors,
				Description: r.FormValue("description"),
			})
		} else {
			lib.Magazines = append(lib.Magazines, library.Magazine{
				Title:       title,
				ISBN:        isbn,
				Authors:     authors,
				PublishedAt: r.FormValue("publishedAt"),
			})
		}

		lib.Save("resources/authors.csv", "resources/books.csv", "resources/magazines.csv")
		http.Redirect(w, r, "/", http.StatusSeeOther)
	})

	fmt.Println("Server starting at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
