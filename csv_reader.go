package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

func parseAuthorsFromReader(r io.Reader) (map[string]Author, error) {
	reader := csv.NewReader(r)
	reader.Comma = ';'

	header, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("failed to read header: %w", err)
	}

	if len(header) > 0 && strings.HasPrefix(header[0], "\ufeff") {
		header[0] = strings.TrimPrefix(header[0], "\ufeff")
	}

	authors := make(map[string]Author)
	
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to read record: %w", err)
		}

		if len(record) < 3 {
			continue
		}

		author := Author{
			Email:     record[0],
			FirstName: record[1],
			LastName:  record[2],
		}
		authors[author.Email] = author
	}

	return authors, nil
}

func readAuthorsFromCSV(filepath string) (map[string]Author, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to open authors file: %w", err)
	}
	defer file.Close()

	return parseAuthorsFromReader(file)
}

func parseBooksFromReader(r io.Reader) ([]Book, error) {
	reader := csv.NewReader(r)
	reader.Comma = ';'

	header, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("failed to read header: %w", err)
	}

	if len(header) > 0 && strings.HasPrefix(header[0], "\ufeff") {
		header[0] = strings.TrimPrefix(header[0], "\ufeff")
	}

	var books []Book
	
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to read record: %w", err)
		}

		if len(record) < 4 {
			continue
		}

		authors := []string{}
		if record[2] != "" {
			authors = strings.Split(record[2], ",")
		}

		book := Book{
			Title:       record[0],
			ISBN:        record[1],
			Authors:     authors,
			Description: record[3],
		}
		books = append(books, book)
	}

	return books, nil
}

func readBooksFromCSV(filepath string) ([]Book, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to open books file: %w", err)
	}
	defer file.Close()

	return parseBooksFromReader(file)
}

func parseMagazinesFromReader(r io.Reader) ([]Magazine, error) {
	reader := csv.NewReader(r)
	reader.Comma = ';'

	header, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("failed to read header: %w", err)
	}

	if len(header) > 0 && strings.HasPrefix(header[0], "\ufeff") {
		header[0] = strings.TrimPrefix(header[0], "\ufeff")
	}

	var magazines []Magazine
	
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to read record: %w", err)
		}

		if len(record) < 4 {
			continue
		}

		authors := []string{}
		if record[2] != "" {
			authors = strings.Split(record[2], ",")
		}

		publishedAt, err := time.Parse("02.01.2006", record[3])
		if err != nil {
			publishedAt = time.Time{}
		}

		magazine := Magazine{
			Title:       record[0],
			ISBN:        record[1],
			Authors:     authors,
			PublishedAt: publishedAt,
		}
		magazines = append(magazines, magazine)
	}

	return magazines, nil
}

func readMagazinesFromCSV(filepath string) ([]Magazine, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to open magazines file: %w", err)
	}
	defer file.Close()

	return parseMagazinesFromReader(file)
}