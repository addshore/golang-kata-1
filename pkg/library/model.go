package library

import (
	"fmt"
	"strings"
)

type Author struct {
	Email     string
	FirstName string
	LastName  string
}

type Book struct {
	Title       string
	ISBN        string
	Authors     []string // List of author emails
	Description string
}

type Magazine struct {
	Title       string
	ISBN        string
	Authors     []string // List of author emails
	PublishedAt string   // Keeping it as string for simplicity, can be parsed to time.Time if needed
}

type Library struct {
	Authors   []Author
	Books     []Book
	Magazines []Magazine
}

func (b Book) Print(authors map[string]Author) {
	fmt.Printf("Book: %s\n", b.Title)
	fmt.Printf("  ISBN: %s\n", b.ISBN)
	fmt.Printf("  Authors: %s\n", FormatAuthors(b.Authors, authors))
	fmt.Printf("  Description: %s\n", b.Description)
	fmt.Println()
}

func (m Magazine) Print(authors map[string]Author) {
	fmt.Printf("Magazine: %s\n", m.Title)
	fmt.Printf("  ISBN: %s\n", m.ISBN)
	fmt.Printf("  Authors: %s\n", FormatAuthors(m.Authors, authors))
	fmt.Printf("  Published At: %s\n", m.PublishedAt)
	fmt.Println()
}

func FormatAuthors(emails []string, authors map[string]Author) string {
	var names []string
	for _, email := range emails {
		if author, ok := authors[email]; ok {
			names = append(names, fmt.Sprintf("%s %s (%s)", author.FirstName, author.LastName, author.Email))
		} else {
			names = append(names, email)
		}
	}
	return strings.Join(names, ", ")
}
