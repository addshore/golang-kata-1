package library

import (
	"encoding/csv"
	"io"
	"os"
	"strings"
)

func ReadAuthors(filename string) ([]Author, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = ';'

	// Read header
	if _, err := reader.Read(); err != nil {
		return nil, err
	}

	var authors []Author
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		authors = append(authors, Author{
			Email:     record[0],
			FirstName: record[1],
			LastName:  record[2],
		})
	}

	return authors, nil
}

func ReadBooks(filename string) ([]Book, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = ';'

	// Read header
	if _, err := reader.Read(); err != nil {
		return nil, err
	}

	var books []Book
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		books = append(books, Book{
			Title:       record[0],
			ISBN:        record[1],
			Authors:     strings.Split(record[2], ","),
			Description: record[3],
		})
	}

	return books, nil
}

func ReadMagazines(filename string) ([]Magazine, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = ';'

	// Read header
	if _, err := reader.Read(); err != nil {
		return nil, err
	}

	var magazines []Magazine
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		magazines = append(magazines, Magazine{
			Title:       record[0],
			ISBN:        record[1],
			Authors:     strings.Split(record[2], ","),
			PublishedAt: record[3],
		})
	}

	return magazines, nil
}

func LoadLibrary(authorsPath, booksPath, magazinesPath string) (*Library, error) {
	authors, err := ReadAuthors(authorsPath)
	if err != nil {
		return nil, err
	}

	books, err := ReadBooks(booksPath)
	if err != nil {
		return nil, err
	}

	magazines, err := ReadMagazines(magazinesPath)
	if err != nil {
		return nil, err
	}

	return &Library{
		Authors:   authors,
		Books:     books,
		Magazines: magazines,
	}, nil
}

func (l *Library) SaveAuthors(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	writer.Comma = ';'
	defer writer.Flush()

	// Write header
	if err := writer.Write([]string{"email", "firstname", "lastname"}); err != nil {
		return err
	}

	for _, author := range l.Authors {
		if err := writer.Write([]string{author.Email, author.FirstName, author.LastName}); err != nil {
			return err
		}
	}

	return nil
}

func (l *Library) SaveBooks(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	writer.Comma = ';'
	defer writer.Flush()

	// Write header
	if err := writer.Write([]string{"title", "isbn", "authors", "description"}); err != nil {
		return err
	}

	for _, book := range l.Books {
		if err := writer.Write([]string{book.Title, book.ISBN, strings.Join(book.Authors, ","), book.Description}); err != nil {
			return err
		}
	}

	return nil
}

func (l *Library) SaveMagazines(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	writer.Comma = ';'
	defer writer.Flush()

	// Write header
	if err := writer.Write([]string{"title", "isbn", "authors", "publishedAt"}); err != nil {
		return err
	}

	for _, magazine := range l.Magazines {
		if err := writer.Write([]string{magazine.Title, magazine.ISBN, strings.Join(magazine.Authors, ","), magazine.PublishedAt}); err != nil {
			return err
		}
	}

	return nil
}

func (l *Library) Save(authorsPath, booksPath, magazinesPath string) error {
	if err := l.SaveAuthors(authorsPath); err != nil {
		return err
	}
	if err := l.SaveBooks(booksPath); err != nil {
		return err
	}
	if err := l.SaveMagazines(magazinesPath); err != nil {
		return err
	}
	return nil
}
