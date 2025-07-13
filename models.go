package main

import (
	"time"
)

// Author represents an author from the authors.csv file
type Author struct {
	Email     string
	FirstName string
	LastName  string
}

// Book represents a book from the books.csv file
type Book struct {
	Title       string
	ISBN        string
	Authors     []Author
	Description string
}

// Magazine represents a magazine from the magazines.csv file
type Magazine struct {
	Title       string
	ISBN        string
	Authors     []Author
	PublishedAt time.Time
}

// Library holds all the library data
type Library struct {
	Authors   map[string]Author // keyed by email
	Books     []Book
	Magazines []Magazine
}
