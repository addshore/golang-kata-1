package main

import (
	"encoding/csv"
	"fmt"
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
	Authors     []string
	Description string
}

type Magazine struct {
	Title       string
	ISBN        string
	Authors     []string
	PublishedAt string
}

func main() {
	authors, err := loadAuthors("resources/authors.csv")
	if err != nil {
		fmt.Println("Error loading authors:", err)
		return
	}

	books, err := loadBooks("resources/books.csv")
	if err != nil {
		fmt.Println("Error loading books:", err)
		return
	}

	magazines, err := loadMagazines("resources/magazines.csv")
	if err != nil {
		fmt.Println("Error loading magazines:", err)
		return
	}

	for {
		fmt.Println("\nLibrary Menu:")
		fmt.Println("1. List all books")
		fmt.Println("2. List all magazines")
		fmt.Println("3. Search by ISBN")
		fmt.Println("4. Search by author email")
		fmt.Println("5. List all books and magazines sorted by title")
		fmt.Println("6. Exit")
		fmt.Println("7. Add a book, magazine, or author")
		fmt.Print("Select an option: ")

		var choice string
		fmt.Scanln(&choice)

		switch choice {
		case "1":
			fmt.Println("\n=== Books ===")
			for _, b := range books {
				printBook(b, authors)
			}
		case "2":
			fmt.Println("\n=== Magazines ===")
			for _, m := range magazines {
				printMagazine(m, authors)
			}
		case "3":
			fmt.Print("Enter ISBN to search: ")
			var isbn string
			fmt.Scanln(&isbn)
			b, m := FindByISBN(books, magazines, isbn)
			if b != nil {
				fmt.Println("\nBook found:")
				printBook(*b, authors)
			}
			if m != nil {
				fmt.Println("\nMagazine found:")
				printMagazine(*m, authors)
			}
			if b == nil && m == nil {
				fmt.Println("No book or magazine found with that ISBN.")
			}
		case "4":
			fmt.Print("Enter author email to search: ")
			var email string
			fmt.Scanln(&email)
			booksByAuthor, magazinesByAuthor := FindByAuthorEmail(books, magazines, email)
			if len(booksByAuthor) > 0 {
				fmt.Println("\nBooks by author:")
				for _, b := range booksByAuthor {
					printBook(b, authors)
				}
			}
			if len(magazinesByAuthor) > 0 {
				fmt.Println("\nMagazines by author:")
				for _, m := range magazinesByAuthor {
					printMagazine(m, authors)
				}
			}
			if len(booksByAuthor) == 0 && len(magazinesByAuthor) == 0 {
				fmt.Println("No book or magazine found for that author email.")
			}
		case "5":
			fmt.Print("Sort by title ascending (a) or descending (d)? [a/d]: ")
			var dir string
			fmt.Scanln(&dir)
			ascending := true
			if strings.ToLower(dir) == "d" {
				ascending = false
			}
			sorted := SortByTitle(books, magazines, ascending)
			fmt.Println("\n=== All Books and Magazines Sorted by Title ===")
			for _, item := range sorted {
				switch v := item.(type) {
				case Book:
					printBook(v, authors)
				case Magazine:
					printMagazine(v, authors)
				}
			}
		case "6":
			fmt.Println("Goodbye!")
			return
		case "7":
			fmt.Println("\nAdd to Library:")
			fmt.Println("a. Add Book")
			fmt.Println("b. Add Magazine")
			fmt.Println("c. Add Author")
			fmt.Print("Select type to add [a/b/c]: ")
			var addType string
			fmt.Scanln(&addType)
			switch strings.ToLower(addType) {
			case "a":
				// Add Book
				var title, isbn, authorEmails, description string
				fmt.Print("Enter book title: ")
				title = readLine()
				fmt.Print("Enter book ISBN: ")
				isbn = readLine()
				fmt.Print("Enter author email(s) (comma separated): ")
				authorEmails = readLine()
				fmt.Print("Enter book description: ")
				description = readLine()
				emails := parseEmails(authorEmails)
				// For each author, check if exists, if not, prompt for details and add
				for _, email := range emails {
					if _, ok := authors[email]; !ok {
						fmt.Printf("Author %s not found. Enter first name: ", email)
						firstName := readLine()
						fmt.Printf("Enter last name for %s: ", email)
						lastName := readLine()
						newAuthor := Author{Email: email, FirstName: firstName, LastName: lastName}
						authors[email] = newAuthor
						appendAuthorToCSV("resources/authors.csv", newAuthor)
						fmt.Println("Author added.")
					}
				}
				newBook := Book{Title: title, ISBN: isbn, Authors: emails, Description: description}
				books = append(books, newBook)
				appendBookToCSV("resources/books.csv", newBook)
				fmt.Println("Book added.")
			case "b":
				// Add Magazine
				var title, isbn, authorEmails, publishedAt string
				fmt.Print("Enter magazine title: ")
				title = readLine()
				fmt.Print("Enter magazine ISBN: ")
				isbn = readLine()
				fmt.Print("Enter author email(s) (comma separated): ")
				authorEmails = readLine()
				fmt.Print("Enter published at (date): ")
				publishedAt = readLine()
				emails := parseEmails(authorEmails)
				for _, email := range emails {
					if _, ok := authors[email]; !ok {
						fmt.Printf("Author %s not found. Enter first name: ", email)
						firstName := readLine()
						fmt.Printf("Enter last name for %s: ", email)
						lastName := readLine()
						newAuthor := Author{Email: email, FirstName: firstName, LastName: lastName}
						authors[email] = newAuthor
						appendAuthorToCSV("resources/authors.csv", newAuthor)
						fmt.Println("Author added.")
					}
				}
				newMagazine := Magazine{Title: title, ISBN: isbn, Authors: emails, PublishedAt: publishedAt}
				magazines = append(magazines, newMagazine)
				appendMagazineToCSV("resources/magazines.csv", newMagazine)
				fmt.Println("Magazine added.")
			case "c":
				// Add Author
				var email, firstName, lastName string
				fmt.Print("Enter author email: ")
				email = readLine()
				if _, ok := authors[email]; ok {
					fmt.Println("Author already exists.")
					break
				}
				fmt.Print("Enter first name: ")
				firstName = readLine()
				fmt.Print("Enter last name: ")
				lastName = readLine()
				newAuthor := Author{Email: email, FirstName: firstName, LastName: lastName}
				authors[email] = newAuthor
				appendAuthorToCSV("resources/authors.csv", newAuthor)
				fmt.Println("Author added.")
			default:
				fmt.Println("Invalid add type.")
			}
		default:
			fmt.Println("Invalid option. Please try again.")
		}
	}
}

func appendAuthorToCSV(path string, author Author) {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening authors.csv for append:", err)
		return
	}
	defer f.Close()
	w := csv.NewWriter(f)
	w.Comma = ';'
	_ = w.Write([]string{author.Email, author.FirstName, author.LastName})
	w.Flush()
}

func appendBookToCSV(path string, book Book) {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening books.csv for append:", err)
		return
	}
	defer f.Close()
	w := csv.NewWriter(f)
	w.Comma = ';'
	_ = w.Write([]string{book.Title, book.ISBN, strings.Join(book.Authors, ","), book.Description})
	w.Flush()
}

func appendMagazineToCSV(path string, mag Magazine) {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening magazines.csv for append:", err)
		return
	}
	defer f.Close()
	w := csv.NewWriter(f)
	w.Comma = ';'
	_ = w.Write([]string{mag.Title, mag.ISBN, strings.Join(mag.Authors, ","), mag.PublishedAt})
	w.Flush()
}

// --- Refactored logic for testability ---
func FindByISBN(books []Book, magazines []Magazine, isbn string) (*Book, *Magazine) {
	for _, b := range books {
		if b.ISBN == isbn {
			return &b, nil
		}
	}
	for _, m := range magazines {
		if m.ISBN == isbn {
			return nil, &m
		}
	}
	return nil, nil
}

func FindByAuthorEmail(books []Book, magazines []Magazine, email string) ([]Book, []Magazine) {
	var booksByAuthor []Book
	var magazinesByAuthor []Magazine
	for _, b := range books {
		for _, a := range b.Authors {
			if strings.TrimSpace(a) == email {
				booksByAuthor = append(booksByAuthor, b)
				break
			}
		}
	}
	for _, m := range magazines {
		for _, a := range m.Authors {
			if strings.TrimSpace(a) == email {
				magazinesByAuthor = append(magazinesByAuthor, m)
				break
			}
		}
	}
	return booksByAuthor, magazinesByAuthor
}

func SortByTitle(books []Book, magazines []Magazine, ascending bool) []interface{} {
	var allItems []interface{}
	for _, b := range books {
		allItems = append(allItems, b)
	}
	for _, m := range magazines {
		allItems = append(allItems, m)
	}
	sort.Slice(allItems, func(i, j int) bool {
		var titleI, titleJ string
		switch v := allItems[i].(type) {
		case Book:
			titleI = v.Title
		case Magazine:
			titleI = v.Title
		}
		switch v := allItems[j].(type) {
		case Book:
			titleJ = v.Title
		case Magazine:
			titleJ = v.Title
		}
		if ascending {
			return strings.ToLower(titleI) < strings.ToLower(titleJ)
		}
		return strings.ToLower(titleI) > strings.ToLower(titleJ)
	})
	return allItems
}

func loadAuthors(path string) (map[string]Author, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	r := csv.NewReader(file)
	r.Comma = ';'
	records, err := r.ReadAll()
	if err != nil {
		return nil, err
	}
	authors := make(map[string]Author)
	for i, rec := range records {
		if i == 0 {
			continue // skip header
		}
		if len(rec) < 3 {
			continue
		}
		authors[rec[0]] = Author{
			Email:     rec[0],
			FirstName: rec[1],
			LastName:  rec[2],
		}
	}
	return authors, nil
}

func loadBooks(path string) ([]Book, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	r := csv.NewReader(file)
	r.Comma = ';'
	records, err := r.ReadAll()
	if err != nil {
		return nil, err
	}
	var books []Book
	for i, rec := range records {
		if i == 0 {
			continue // skip header
		}
		if len(rec) < 4 {
			continue
		}
		authors := strings.Split(rec[2], ",")
		books = append(books, Book{
			Title:       rec[0],
			ISBN:        rec[1],
			Authors:     authors,
			Description: rec[3],
		})
	}
	return books, nil
}

func loadMagazines(path string) ([]Magazine, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	r := csv.NewReader(file)
	r.Comma = ';'
	records, err := r.ReadAll()
	if err != nil {
		return nil, err
	}
	var magazines []Magazine
	for i, rec := range records {
		if i == 0 {
			continue // skip header
		}
		if len(rec) < 4 {
			continue
		}
		authors := strings.Split(rec[2], ",")
		magazines = append(magazines, Magazine{
			Title:       rec[0],
			ISBN:        rec[1],
			Authors:     authors,
			PublishedAt: rec[3],
		})
	}
	return magazines, nil
}

func resolveAuthors(emails []string, authors map[string]Author) string {
	var names []string
	for _, email := range emails {
		email = strings.TrimSpace(email)
		if a, ok := authors[email]; ok {
			names = append(names, a.FirstName+" "+a.LastName)
		} else {
			names = append(names, email)
		}
	}
	return strings.Join(names, ", ")
}

func printBook(b Book, authors map[string]Author) {
	fmt.Printf("Title: %s\nISBN: %s\nAuthors: %s\nDescription: %s\n\n",
		b.Title, b.ISBN, resolveAuthors(b.Authors, authors), b.Description)
}

func printMagazine(m Magazine, authors map[string]Author) {
	fmt.Printf("Title: %s\nISBN: %s\nAuthors: %s\nPublished At: %s\n\n",
		m.Title, m.ISBN, resolveAuthors(m.Authors, authors), m.PublishedAt)
}

// Helper to read a full line from stdin
func readLine() string {
	var input string
	fmt.Scanln(&input)
	return input
}

// Helper to parse comma-separated emails and trim whitespace
func parseEmails(s string) []string {
	parts := strings.Split(s, ",")
	var emails []string
	for _, p := range parts {
		email := strings.TrimSpace(p)
		if email != "" {
			emails = append(emails, email)
		}
	}
	return emails
}
