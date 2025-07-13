package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	for {
		fmt.Println("\nLibrary CLI")
		fmt.Println("1. List all books")
		fmt.Println("2. List all magazines")
		fmt.Println("3. Search by ISBN")
		fmt.Println("4. Search by author email")
		fmt.Println("5. List all sorted by title")
		fmt.Println("6. Add book or magazine")
		fmt.Println("0. Exit")
		fmt.Print("Choose an option: ")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		choice := scanner.Text()
		switch choice {
		case "1":
			books, _ := readBooksCSV("resources/books.csv")
			printBooks(books)
		case "2":
			mags, _ := readMagazinesCSV("resources/magazines.csv")
			printMagazines(mags)
		case "3":
			fmt.Print("Enter ISBN: ")
			scanner.Scan()
			isbn := scanner.Text()
			books, _ := readBooksCSV("resources/books.csv")
			mags, _ := readMagazinesCSV("resources/magazines.csv")
			printBooks(filterBooksByISBN(books, isbn))
			printMagazines(filterMagazinesByISBN(mags, isbn))
		case "4":
			fmt.Print("Enter author email: ")
			scanner.Scan()
			email := scanner.Text()
			books, _ := readBooksCSV("resources/books.csv")
			mags, _ := readMagazinesCSV("resources/magazines.csv")
			printBooks(filterBooksByAuthor(books, email))
			printMagazines(filterMagazinesByAuthor(mags, email))
		case "5":
			books, _ := readBooksCSV("resources/books.csv")
			mags, _ := readMagazinesCSV("resources/magazines.csv")
			items := combineBooksAndMagazines(books, mags)
			sortCombinedItemsByTitle(items, true)
			printCombinedItems(items)
		case "6":
			addCLIEntry()
		case "0":
			fmt.Println("Goodbye!")
			os.Exit(0)
		default:
			fmt.Println("Invalid option")
		}
	}
}

func printBooks(books []Book) {
	fmt.Println("\nBooks:")
	for _, b := range books {
		fmt.Printf("Title: %s | ISBN: %s | Authors: %s | Description: %s\n", b.Title, b.ISBN, b.Authors, b.Description)
	}
}

func printMagazines(mags []Magazine) {
	fmt.Println("\nMagazines:")
	for _, m := range mags {
		fmt.Printf("Title: %s | ISBN: %s | Authors: %s | Published At: %s\n", m.Title, m.ISBN, m.Authors, m.PublishedAt)
	}
}

func printCombinedItems(items []CombinedItem) {
	fmt.Println("\nAll Books and Magazines (Sorted):")
	for _, i := range items {
		fmt.Printf("%s | Title: %s | ISBN: %s | Authors: %s | Description: %s | Published At: %s\n",
			i.Type, i.Title, i.ISBN, i.Authors, i.Description, i.PublishedAt)
	}
}

func filterBooksByISBN(books []Book, isbn string) []Book {
	var res []Book
	for _, b := range books {
		if b.ISBN == isbn {
			res = append(res, b)
		}
	}
	return res
}

func filterMagazinesByISBN(mags []Magazine, isbn string) []Magazine {
	var res []Magazine
	for _, m := range mags {
		if m.ISBN == isbn {
			res = append(res, m)
		}
	}
	return res
}

func filterBooksByAuthor(books []Book, email string) []Book {
	var res []Book
	for _, b := range books {
		if containsAuthor(b.Authors, email) {
			res = append(res, b)
		}
	}
	return res
}

func filterMagazinesByAuthor(mags []Magazine, email string) []Magazine {
	var res []Magazine
	for _, m := range mags {
		if containsAuthor(m.Authors, email) {
			res = append(res, m)
		}
	}
	return res
}

func addCLIEntry() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Type (book/magazine): ")
	scanner.Scan()
	typ := scanner.Text()
	fmt.Print("Title: ")
	scanner.Scan()
	title := scanner.Text()
	fmt.Print("ISBN: ")
	scanner.Scan()
	isbn := scanner.Text()
	fmt.Print("Authors (comma-separated emails): ")
	scanner.Scan()
	authors := scanner.Text()
	var description, publishedAt string
	if typ == "book" {
		fmt.Print("Description: ")
		scanner.Scan()
		description = scanner.Text()
	} else if typ == "magazine" {
		fmt.Print("Published At (DD.MM.YYYY): ")
		scanner.Scan()
		publishedAt = scanner.Text()
	}
	if typ == "book" {
		appendToCSV("resources/books.csv", []string{title, isbn, authors, description})
	} else if typ == "magazine" {
		appendToCSV("resources/magazines.csv", []string{title, isbn, authors, publishedAt})
	}
	// Add authors if not present
	addAuthorsCLI(authors, scanner)
	fmt.Println("Entry added.")
}

func addAuthorsCLI(authors string, scanner *bufio.Scanner) {
	existing := map[string]bool{}
	file, err := os.Open("resources/authors.csv")
	if err == nil {
		reader := bufio.NewReader(file)
		lines := []string{}
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				break
			}
			lines = append(lines, line)
		}
		for _, line := range lines[1:] {
			parts := strings.Split(strings.TrimSpace(line), ";")
			if len(parts) > 0 {
				existing[parts[0]] = true
			}
		}
		file.Close()
	}
	for _, email := range splitAuthors(authors) {
		email = trim(email)
		if email == "" || existing[email] {
			continue
		}
		fmt.Printf("New author detected: %s\n", email)
		fmt.Print("First name: ")
		scanner.Scan()
		first := scanner.Text()
		fmt.Print("Last name: ")
		scanner.Scan()
		last := scanner.Text()
		appendToCSV("resources/authors.csv", []string{email, first, last})
	}
}
