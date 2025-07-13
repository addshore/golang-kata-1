package main

import (
	"encoding/csv"
	"os"
)

// Book, Magazine, CombinedItem structs
type Book struct {
	Title       string
	ISBN        string
	Authors     string
	Description string
}

type Magazine struct {
	Title       string
	ISBN        string
	Authors     string
	PublishedAt string
}

type CombinedItem struct {
	Type        string
	Title       string
	ISBN        string
	Authors     string
	Description string
	PublishedAt string
}

// CSV reading functions
func readBooksCSV(path string) ([]Book, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	reader := csv.NewReader(file)
	reader.Comma = ';'
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}
	var books []Book
	for i, rec := range records {
		if i == 0 {
			continue
		}
		books = append(books, Book{
			Title:       rec[0],
			ISBN:        rec[1],
			Authors:     rec[2],
			Description: rec[3],
		})
	}
	return books, nil
}

func readMagazinesCSV(path string) ([]Magazine, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	reader := csv.NewReader(file)
	reader.Comma = ';'
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}
	var magazines []Magazine
	for i, rec := range records {
		if i == 0 {
			continue
		}
		magazines = append(magazines, Magazine{
			Title:       rec[0],
			ISBN:        rec[1],
			Authors:     rec[2],
			PublishedAt: rec[3],
		})
	}
	return magazines, nil
}

// Combine and sort helpers
func combineBooksAndMagazines(books []Book, magazines []Magazine) []CombinedItem {
	var items []CombinedItem
	for _, b := range books {
		items = append(items, CombinedItem{
			Type:        "Book",
			Title:       b.Title,
			ISBN:        b.ISBN,
			Authors:     b.Authors,
			Description: b.Description,
			PublishedAt: "",
		})
	}
	for _, m := range magazines {
		items = append(items, CombinedItem{
			Type:        "Magazine",
			Title:       m.Title,
			ISBN:        m.ISBN,
			Authors:     m.Authors,
			Description: "",
			PublishedAt: m.PublishedAt,
		})
	}
	return items
}

func sortCombinedItemsByTitle(items []CombinedItem, asc bool) {
	for i := 0; i < len(items)-1; i++ {
		for j := i + 1; j < len(items); j++ {
			if asc {
				if items[i].Title > items[j].Title {
					items[i], items[j] = items[j], items[i]
				}
			} else {
				if items[i].Title < items[j].Title {
					items[i], items[j] = items[j], items[i]
				}
			}
		}
	}
}

// Author helpers
func containsAuthor(authors string, email string) bool {
	for _, a := range splitAuthors(authors) {
		if a == email {
			return true
		}
	}
	return false
}

func splitAuthors(authors string) []string {
	var res []string
	for _, a := range splitAndTrim(authors, ",") {
		res = append(res, a)
	}
	return res
}

func splitAndTrim(s, sep string) []string {
	var res []string
	for _, part := range split(s, sep) {
		res = append(res, trim(part))
	}
	return res
}

func split(s, sep string) []string {
	var res []string
	start := 0
	for i := 0; i+len(sep) <= len(s); i++ {
		if s[i:i+len(sep)] == sep {
			res = append(res, s[start:i])
			start = i + len(sep)
		}
	}
	res = append(res, s[start:])
	return res
}

func trim(s string) string {
	start, end := 0, len(s)
	for start < end && (s[start] == ' ' || s[start] == '\t' || s[start] == '\n') {
		start++
	}
	for end > start && (s[end-1] == ' ' || s[end-1] == '\t' || s[end-1] == '\n') {
		end--
	}
	return s[start:end]
}

// CSV append helper
func appendToCSV(path string, row []string) error {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	w := csv.NewWriter(f)
	w.Comma = ';'
	if err := w.Write(row); err != nil {
		return err
	}
	w.Flush()
	return w.Error()
}
