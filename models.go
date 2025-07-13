package main

import (
	"time"
)

type Author struct {
	Email     string
	FirstName string
	LastName  string
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
	PublishedAt time.Time
}

type LibraryItem struct {
	Title       string
	ISBN        string
	Authors     []string
	Type        string // "book" or "magazine"
	Description string
	PublishedAt time.Time
}