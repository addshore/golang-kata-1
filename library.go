package main

import (
	"encoding/csv"
	"io"
	"sort"
	"strings"
)

type Library struct {
	Authors   map[string]Author
	Books     []Book
	Magazines []Magazine
}

func NewLibrary(authors map[string]Author, books []Book, magazines []Magazine) *Library {
	return &Library{
		Authors:   authors,
		Books:     books,
		Magazines: magazines,
	}
}

func (l *Library) FindByISBN(isbn string) (interface{}, bool) {
	for _, book := range l.Books {
		if book.ISBN == isbn {
			return book, true
		}
	}
	for _, magazine := range l.Magazines {
		if magazine.ISBN == isbn {
			return magazine, true
		}
	}
	return nil, false
}

func (l *Library) FindByAuthor(email string) ([]Book, []Magazine) {
	var books []Book
	var magazines []Magazine
	
	for _, book := range l.Books {
		for _, authorEmail := range book.Authors {
			if authorEmail == email {
				books = append(books, book)
				break
			}
		}
	}
	
	for _, magazine := range l.Magazines {
		for _, authorEmail := range magazine.Authors {
			if authorEmail == email {
				magazines = append(magazines, magazine)
				break
			}
		}
	}
	
	return books, magazines
}

func (l *Library) GetSortedItems(desc bool) []Item {
	var items []Item
	
	for _, book := range l.Books {
		items = append(items, Item{
			Title:       book.Title,
			ISBN:        book.ISBN,
			Authors:     book.Authors,
			Description: book.Description,
			IsBook:      true,
		})
	}
	
	for _, magazine := range l.Magazines {
		items = append(items, Item{
			Title:       magazine.Title,
			ISBN:        magazine.ISBN,
			Authors:     magazine.Authors,
			PublishedAt: magazine.PublishedAt,
			IsBook:      false,
		})
	}
	
	if desc {
		sort.Slice(items, func(i, j int) bool {
			return strings.ToLower(items[i].Title) > strings.ToLower(items[j].Title)
		})
	} else {
		sort.Slice(items, func(i, j int) bool {
			return strings.ToLower(items[i].Title) < strings.ToLower(items[j].Title)
		})
	}
	
	return items
}

func FormatAuthors(authorEmails []string, authors map[string]Author) string {
	var names []string
	for _, email := range authorEmails {
		if author, exists := authors[email]; exists {
			names = append(names, author.Firstname+" "+author.Lastname)
		}
	}
	return strings.Join(names, ", ")
}

func ParseAuthorsCSV(reader io.Reader) (map[string]Author, error) {
	csvReader := csv.NewReader(reader)
	csvReader.Comma = ';'
	records, err := csvReader.ReadAll()
	if err != nil {
		return nil, err
	}

	authors := make(map[string]Author)
	for i, record := range records {
		if i == 0 {
			continue // skip header
		}
		authors[record[0]] = Author{
			Email:     record[0],
			Firstname: record[1],
			Lastname:  record[2],
		}
	}
	return authors, nil
}

func ParseBooksCSV(reader io.Reader) ([]Book, error) {
	csvReader := csv.NewReader(reader)
	csvReader.Comma = ';'
	records, err := csvReader.ReadAll()
	if err != nil {
		return nil, err
	}

	var books []Book
	for i, record := range records {
		if i == 0 {
			continue // skip header
		}
		authors := strings.Split(record[2], ",")
		books = append(books, Book{
			Title:       record[0],
			ISBN:        record[1],
			Authors:     authors,
			Description: record[3],
		})
	}
	return books, nil
}

func ParseMagazinesCSV(reader io.Reader) ([]Magazine, error) {
	csvReader := csv.NewReader(reader)
	csvReader.Comma = ';'
	records, err := csvReader.ReadAll()
	if err != nil {
		return nil, err
	}

	var magazines []Magazine
	for i, record := range records {
		if i == 0 {
			continue // skip header
		}
		authors := strings.Split(record[2], ",")
		magazines = append(magazines, Magazine{
			Title:       record[0],
			ISBN:        record[1],
			Authors:     authors,
			PublishedAt: record[3],
		})
	}
	return magazines, nil
}

func (l *Library) AddBook(book Book) {
	l.Books = append(l.Books, book)
	for _, authorEmail := range book.Authors {
		if _, exists := l.Authors[authorEmail]; !exists {
			// Author will be added separately
		}
	}
}

func (l *Library) AddMagazine(magazine Magazine) {
	l.Magazines = append(l.Magazines, magazine)
	for _, authorEmail := range magazine.Authors {
		if _, exists := l.Authors[authorEmail]; !exists {
			// Author will be added separately
		}
	}
}

func (l *Library) AddAuthor(author Author) {
	l.Authors[author.Email] = author
}

func WriteAuthorsCSV(writers io.Writer, authors map[string]Author) error {
	w := csv.NewWriter(writers)
	w.Comma = ';'
	defer w.Flush()
	
	if err := w.Write([]string{"email", "firstname", "lastname"}); err != nil {
		return err
	}
	
	for _, author := range authors {
		if err := w.Write([]string{author.Email, author.Firstname, author.Lastname}); err != nil {
			return err
		}
	}
	return nil
}

func WriteBooksCSV(writers io.Writer, books []Book) error {
	w := csv.NewWriter(writers)
	w.Comma = ';'
	defer w.Flush()
	
	if err := w.Write([]string{"title", "isbn", "authors", "description"}); err != nil {
		return err
	}
	
	for _, book := range books {
		authors := strings.Join(book.Authors, ",")
		if err := w.Write([]string{book.Title, book.ISBN, authors, book.Description}); err != nil {
			return err
		}
	}
	return nil
}

func WriteMagazinesCSV(writers io.Writer, magazines []Magazine) error {
	w := csv.NewWriter(writers)
	w.Comma = ';'
	defer w.Flush()
	
	if err := w.Write([]string{"title", "isbn", "authors", "publishedAt"}); err != nil {
		return err
	}
	
	for _, magazine := range magazines {
		authors := strings.Join(magazine.Authors, ",")
		if err := w.Write([]string{magazine.Title, magazine.ISBN, authors, magazine.PublishedAt}); err != nil {
			return err
		}
	}
	return nil
}