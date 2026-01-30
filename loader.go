package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"
)

func loadCSV(path string) ([][]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	r := csv.NewReader(f)
	r.Comma = ';'
	r.LazyQuotes = true // Be robust against quoting

	// Read all records
	records, err := r.ReadAll()
	if err != nil {
		return nil, err
	}

	if len(records) < 2 {
		return nil, fmt.Errorf("CSV file %s is empty or missing header", path)
	}

	// Return records skipping the header
	return records[1:], nil
}

func LoadAuthors(path string) (map[string]*Author, error) {
	records, err := loadCSV(path)
	if err != nil {
		return nil, err
	}

	authors := make(map[string]*Author)
	for _, record := range records {
		if len(record) < 3 {
			continue
		}
		email := record[0]
		firstname := record[1]
		lastname := record[2]

		authors[email] = &Author{
			Email:     email,
			Firstname: firstname,
			Lastname:  lastname,
		}
	}
	return authors, nil
}

func parseAuthors(emailsStr string, authorMap map[string]*Author) []*Author {
	if emailsStr == "" {
		return nil
	}
	emails := strings.Split(emailsStr, ",")
	var result []*Author
	for _, email := range emails {
		email = strings.TrimSpace(email)
		if author, ok := authorMap[email]; ok {
			result = append(result, author)
		}
	}
	return result
}

func LoadBooks(path string, authorMap map[string]*Author) ([]*Book, error) {
	records, err := loadCSV(path)
	if err != nil {
		return nil, err
	}

	var books []*Book
	for _, record := range records {
		if len(record) < 4 {
			continue
		}
		// title;isbn;authors;description
		books = append(books, &Book{
			Title:       record[0],
			ISBN:        record[1],
			Authors:     parseAuthors(record[2], authorMap),
			Description: record[3],
		})
	}
	return books, nil
}

func LoadMagazines(path string, authorMap map[string]*Author) ([]*Magazine, error) {
	records, err := loadCSV(path)
	if err != nil {
		return nil, err
	}

	var magazines []*Magazine
	for _, record := range records {
		if len(record) < 4 {
			continue
		}
		// title;isbn;authors;publishedAt
		magazines = append(magazines, &Magazine{
			Title:       record[0],
			ISBN:        record[1],
			Authors:     parseAuthors(record[2], authorMap),
			PublishedAt: record[3],
		})
	}
	return magazines, nil
}
