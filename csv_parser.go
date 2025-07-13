package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"
	"time"
)

// LoadLibrary loads all data from CSV files and returns a Library struct
func LoadLibrary() (*Library, error) {
	library := &Library{
		Authors:   make(map[string]Author),
		Books:     []Book{},
		Magazines: []Magazine{},
	}

	// Load authors first
	if err := loadAuthors(library); err != nil {
		return nil, fmt.Errorf("failed to load authors: %w", err)
	}

	// Load books
	if err := loadBooks(library); err != nil {
		return nil, fmt.Errorf("failed to load books: %w", err)
	}

	// Load magazines
	if err := loadMagazines(library); err != nil {
		return nil, fmt.Errorf("failed to load magazines: %w", err)
	}

	return library, nil
}

// loadAuthors loads authors from authors.csv
func loadAuthors(library *Library) error {
	file, err := os.Open("resources/authors.csv")
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = ';'

	records, err := reader.ReadAll()
	if err != nil {
		return err
	}

	// Skip header row
	for i := 1; i < len(records); i++ {
		record := records[i]
		if len(record) < 3 {
			continue
		}

		author := Author{
			Email:     strings.TrimSpace(record[0]),
			FirstName: strings.TrimSpace(record[1]),
			LastName:  strings.TrimSpace(record[2]),
		}

		library.Authors[author.Email] = author
	}

	return nil
}

// loadBooks loads books from books.csv
func loadBooks(library *Library) error {
	file, err := os.Open("resources/books.csv")
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = ';'

	records, err := reader.ReadAll()
	if err != nil {
		return err
	}

	// Skip header row
	for i := 1; i < len(records); i++ {
		record := records[i]
		if len(record) < 4 {
			continue
		}

		book := Book{
			Title:       strings.TrimSpace(record[0]),
			ISBN:        strings.TrimSpace(record[1]),
			Description: strings.TrimSpace(record[3]),
		}

		// Parse authors
		authorsStr := strings.TrimSpace(record[2])
		if authorsStr != "" {
			authorEmails := strings.Split(authorsStr, ",")
			for _, email := range authorEmails {
				email = strings.TrimSpace(email)
				if author, exists := library.Authors[email]; exists {
					book.Authors = append(book.Authors, author)
				}
			}
		}

		library.Books = append(library.Books, book)
	}

	return nil
}

// loadMagazines loads magazines from magazines.csv
func loadMagazines(library *Library) error {
	file, err := os.Open("resources/magazines.csv")
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = ';'

	records, err := reader.ReadAll()
	if err != nil {
		return err
	}

	// Skip header row
	for i := 1; i < len(records); i++ {
		record := records[i]
		if len(record) < 4 {
			continue
		}

		magazine := Magazine{
			Title: strings.TrimSpace(record[0]),
			ISBN:  strings.TrimSpace(record[1]),
		}

		// Parse published date
		publishedStr := strings.TrimSpace(record[3])
		if publishedStr != "" {
			publishedAt, err := time.Parse("02.01.2006", publishedStr)
			if err != nil {
				fmt.Printf("Warning: Could not parse date '%s' for magazine '%s': %v\n", publishedStr, magazine.Title, err)
			} else {
				magazine.PublishedAt = publishedAt
			}
		}

		// Parse authors
		authorsStr := strings.TrimSpace(record[2])
		if authorsStr != "" {
			authorEmails := strings.Split(authorsStr, ",")
			for _, email := range authorEmails {
				email = strings.TrimSpace(email)
				if author, exists := library.Authors[email]; exists {
					magazine.Authors = append(magazine.Authors, author)
				}
			}
		}

		library.Magazines = append(library.Magazines, magazine)
	}

	return nil
}

// SaveLibrary saves the library data back to CSV files
func SaveLibrary(library *Library) error {
	if err := saveAuthors(library.Authors); err != nil {
		return fmt.Errorf("failed to save authors: %v", err)
	}
	if err := saveBooks(library.Books); err != nil {
		return fmt.Errorf("failed to save books: %v", err)
	}
	if err := saveMagazines(library.Magazines); err != nil {
		return fmt.Errorf("failed to save magazines: %v", err)
	}
	return nil
}

// saveAuthors saves authors to CSV file
func saveAuthors(authors map[string]Author) error {
	file, err := os.Create("resources/authors.csv")
	if err != nil {
		return err
	}
	defer file.Close()

	// Write BOM for UTF-8
	file.WriteString("\ufeff")

	writer := csv.NewWriter(file)
	writer.Comma = ';'
	defer writer.Flush()

	// Write header
	if err := writer.Write([]string{"email", "firstname", "lastname"}); err != nil {
		return err
	}

	// Write data
	for _, author := range authors {
		record := []string{author.Email, author.FirstName, author.LastName}
		if err := writer.Write(record); err != nil {
			return err
		}
	}

	return nil
}

// saveBooks saves books to CSV file
func saveBooks(books []Book) error {
	file, err := os.Create("resources/books.csv")
	if err != nil {
		return err
	}
	defer file.Close()

	// Write BOM for UTF-8
	file.WriteString("\ufeff")

	writer := csv.NewWriter(file)
	writer.Comma = ';'
	defer writer.Flush()

	// Write header
	if err := writer.Write([]string{"title", "isbn", "authors", "description"}); err != nil {
		return err
	}

	// Write data
	for _, book := range books {
		var authorEmails []string
		for _, author := range book.Authors {
			authorEmails = append(authorEmails, author.Email)
		}
		authorsStr := strings.Join(authorEmails, ",")
		record := []string{book.Title, book.ISBN, authorsStr, book.Description}
		if err := writer.Write(record); err != nil {
			return err
		}
	}

	return nil
}

// saveMagazines saves magazines to CSV file
func saveMagazines(magazines []Magazine) error {
	file, err := os.Create("resources/magazines.csv")
	if err != nil {
		return err
	}
	defer file.Close()

	// Write BOM for UTF-8
	file.WriteString("\ufeff")

	writer := csv.NewWriter(file)
	writer.Comma = ';'
	defer writer.Flush()

	// Write header
	if err := writer.Write([]string{"title", "isbn", "authors", "publishedAt"}); err != nil {
		return err
	}

	// Write data
	for _, magazine := range magazines {
		var authorEmails []string
		for _, author := range magazine.Authors {
			authorEmails = append(authorEmails, author.Email)
		}
		authorsStr := strings.Join(authorEmails, ",")
		publishedAtStr := magazine.PublishedAt.Format("02.01.2006")
		record := []string{magazine.Title, magazine.ISBN, authorsStr, publishedAtStr}
		if err := writer.Write(record); err != nil {
			return err
		}
	}

	return nil
}
