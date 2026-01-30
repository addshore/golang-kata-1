package main

type Author struct {
	Email     string
	Firstname string
	Lastname  string
}

type Book struct {
	Title       string
	ISBN        string
	Authors     []*Author
	Description string
}

type Magazine struct {
	Title       string
	ISBN        string
	Authors     []*Author
	PublishedAt string
}

type DisplayItem struct {
	Title     string
	ISBN      string
	Authors   []*Author
	ExtraInfo string // Description or PublishedAt
	Type      string // "Book" or "Magazine"
}

type LibraryData struct {
	Authors     map[string]*Author
	Books       []*Book
	Magazines   []*Magazine
	SortedItems []*DisplayItem
}
