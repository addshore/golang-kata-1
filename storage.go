package main

import (
	"encoding/csv"
	"os"
	"strings"
)

// appendToCSV appends a single record to the given file path using ';' as separator.
func appendToCSV(path string, record []string) error {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	w := csv.NewWriter(f)
	w.Comma = ';'

	if err := w.Write(record); err != nil {
		return err
	}
	w.Flush()
	return w.Error()
}

func AppendAuthor(path string, author *Author) error {
	record := []string{author.Email, author.Firstname, author.Lastname}
	return appendToCSV(path, record)
}

func AppendBook(path string, book *Book) error {
	authors := make([]string, len(book.Authors))
	for i, a := range book.Authors {
		authors[i] = a.Email
	}
	authorsStr := strings.Join(authors, ",")

	record := []string{book.Title, book.ISBN, authorsStr, book.Description}
	return appendToCSV(path, record)
}

func AppendMagazine(path string, mag *Magazine) error {
	authors := make([]string, len(mag.Authors))
	for i, a := range mag.Authors {
		authors[i] = a.Email
	}
	authorsStr := strings.Join(authors, ",")

	record := []string{mag.Title, mag.ISBN, authorsStr, mag.PublishedAt}
	return appendToCSV(path, record)
}
