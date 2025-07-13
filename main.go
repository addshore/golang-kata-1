package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Author struct {
	Email     string
	Firstname string
	Lastname  string
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

type Item struct {
	Title       string
	ISBN        string
	Authors     []string
	Description string
	PublishedAt string
	IsBook      bool
}

func main() {
	authors := loadAuthors()
	books := loadBooks()
	magazines := loadMagazines()
	library := NewLibrary(authors, books, magazines)

	fmt.Println("=== LIBRARY SYSTEM ===")
	if len(os.Args) > 1 {
		if os.Args[1] == "search" {
			if len(os.Args) < 3 {
				fmt.Println("Usage: go run main.go search <ISBN>")
				return
			}
			searchISBN(os.Args[2], library)
		} else if os.Args[1] == "author" {
			if len(os.Args) < 3 {
				fmt.Println("Usage: go run main.go author <email>")
				return
			}
			searchByAuthor(os.Args[2], library)
		} else if os.Args[1] == "sort" {
			desc := len(os.Args) > 2 && os.Args[2] == "desc"
			listSorted(library, desc)
		} else if os.Args[1] == "add" {
			addItem()
		} else if os.Args[1] == "web" {
			startWebServer()
		} else {
			listAll(library)
		}
	} else {
		listAll(library)
	}
}

func listAll(library *Library) {
	fmt.Printf("\n📚 BOOKS (%d):\n", len(library.Books))
	for _, book := range library.Books {
		displayBook(book, library.Authors)
	}

	fmt.Printf("\n📖 MAGAZINES (%d):\n", len(library.Magazines))
	for _, magazine := range library.Magazines {
		displayMagazine(magazine, library.Authors)
	}
}

func searchISBN(isbn string, library *Library) {
	item, found := library.FindByISBN(isbn)
	if found {
		switch v := item.(type) {
		case Book:
			fmt.Println("\n📚 BOOK FOUND:")
			displayBook(v, library.Authors)
		case Magazine:
			fmt.Println("\n📖 MAGAZINE FOUND:")
			displayMagazine(v, library.Authors)
		}
	} else {
		fmt.Printf("\nNo book or magazine found with ISBN: %s\n", isbn)
	}
}

func displayBook(book Book, authors map[string]Author) {
	fmt.Printf("\nTitle: %s\n", book.Title)
	fmt.Printf("ISBN: %s\n", book.ISBN)
	fmt.Printf("Authors: %s\n", FormatAuthors(book.Authors, authors))
	fmt.Printf("Description: %s\n", book.Description)
	fmt.Println(strings.Repeat("-", 50))
}

func displayMagazine(magazine Magazine, authors map[string]Author) {
	fmt.Printf("\nTitle: %s\n", magazine.Title)
	fmt.Printf("ISBN: %s\n", magazine.ISBN)
	fmt.Printf("Authors: %s\n", FormatAuthors(magazine.Authors, authors))
	fmt.Printf("Published: %s\n", magazine.PublishedAt)
	fmt.Println(strings.Repeat("-", 50))
}

func searchByAuthor(email string, library *Library) {
	books, magazines := library.FindByAuthor(email)
	found := len(books) > 0 || len(magazines) > 0
	
	if len(books) > 0 {
		fmt.Printf("\n📚 BOOKS BY %s:\n", email)
		for _, book := range books {
			displayBook(book, library.Authors)
		}
	}
	
	if len(magazines) > 0 {
		fmt.Printf("\n📖 MAGAZINES BY %s:\n", email)
		for _, magazine := range magazines {
			displayMagazine(magazine, library.Authors)
		}
	}
	
	if !found {
		fmt.Printf("\nNo books or magazines found by author: %s\n", email)
	}
}

func loadAuthors() map[string]Author {
	file, err := os.Open("resources/authors.csv")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	authors, err := ParseAuthorsCSV(file)
	if err != nil {
		panic(err)
	}
	return authors
}

func loadBooks() []Book {
	file, err := os.Open("resources/books.csv")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	books, err := ParseBooksCSV(file)
	if err != nil {
		panic(err)
	}
	return books
}

func loadMagazines() []Magazine {
	file, err := os.Open("resources/magazines.csv")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	magazines, err := ParseMagazinesCSV(file)
	if err != nil {
		panic(err)
	}
	return magazines
}



func listSorted(library *Library, desc bool) {
	items := library.GetSortedItems(desc)
	
	if desc {
		fmt.Printf("\n📚📖 ALL ITEMS SORTED BY TITLE (DESC) (%d):\n", len(items))
	} else {
		fmt.Printf("\n📚📖 ALL ITEMS SORTED BY TITLE (ASC) (%d):\n", len(items))
	}
	
	for _, item := range items {
		if item.IsBook {
			fmt.Printf("\n[BOOK] Title: %s\n", item.Title)
			fmt.Printf("ISBN: %s\n", item.ISBN)
			fmt.Printf("Authors: %s\n", FormatAuthors(item.Authors, library.Authors))
			fmt.Printf("Description: %s\n", item.Description)
		} else {
			fmt.Printf("\n[MAGAZINE] Title: %s\n", item.Title)
			fmt.Printf("ISBN: %s\n", item.ISBN)
			fmt.Printf("Authors: %s\n", FormatAuthors(item.Authors, library.Authors))
			fmt.Printf("Published: %s\n", item.PublishedAt)
		}
		fmt.Println(strings.Repeat("-", 50))
	}
}

func addItem() {
	scanner := bufio.NewScanner(os.Stdin)
	
	fmt.Print("Add (b)ook or (m)agazine? ")
	scanner.Scan()
	itemType := strings.ToLower(strings.TrimSpace(scanner.Text()))
	
	if itemType != "b" && itemType != "book" && itemType != "m" && itemType != "magazine" {
		fmt.Println("Invalid choice. Use 'b' for book or 'm' for magazine.")
		return
	}
	
	fmt.Print("Title: ")
	scanner.Scan()
	title := strings.TrimSpace(scanner.Text())
	
	fmt.Print("ISBN: ")
	scanner.Scan()
	isbn := strings.TrimSpace(scanner.Text())
	
	fmt.Print("Author email: ")
	scanner.Scan()
	authorEmail := strings.TrimSpace(scanner.Text())
	
	fmt.Print("Author first name: ")
	scanner.Scan()
	authorFirst := strings.TrimSpace(scanner.Text())
	
	fmt.Print("Author last name: ")
	scanner.Scan()
	authorLast := strings.TrimSpace(scanner.Text())
	
	if itemType == "b" || itemType == "book" {
		fmt.Print("Description: ")
		scanner.Scan()
		description := strings.TrimSpace(scanner.Text())
		addBook(title, isbn, authorEmail, authorFirst, authorLast, description)
	} else {
		fmt.Print("Publication date (DD.MM.YYYY): ")
		scanner.Scan()
		pubDate := strings.TrimSpace(scanner.Text())
		addMagazine(title, isbn, authorEmail, authorFirst, authorLast, pubDate)
	}
	
	fmt.Println("Item added successfully!")
}

func addBook(title, isbn, authorEmail, authorFirst, authorLast, description string) {
	addAuthor(authorEmail, authorFirst, authorLast)
	appendToCSV("resources/books.csv", []string{title, isbn, authorEmail, description})
}

func addMagazine(title, isbn, authorEmail, authorFirst, authorLast, pubDate string) {
	addAuthor(authorEmail, authorFirst, authorLast)
	appendToCSV("resources/magazines.csv", []string{title, isbn, authorEmail, pubDate})
}

func addAuthor(email, firstname, lastname string) {
	authors := loadAuthors()
	if _, exists := authors[email]; !exists {
		appendToCSV("resources/authors.csv", []string{email, firstname, lastname})
	}
}

func appendToCSV(filename string, record []string) {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	
	line := strings.Join(record, ";") + "\r\n"
	if _, err := file.WriteString(line); err != nil {
		panic(err)
	}
}
