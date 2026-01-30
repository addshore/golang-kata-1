package library

import (
	"testing"
	"path/filepath"
	. "github.com/onsi/gomega"
)

func TestLoadLibrary(t *testing.T) {
	g := NewGomegaWithT(t)

	// Use relative paths from the project root
	authorsPath := filepath.Join("..", "..", "resources", "authors.csv")
	booksPath := filepath.Join("..", "..", "resources", "books.csv")
	magazinesPath := filepath.Join("..", "..", "resources", "magazines.csv")

	lib, err := LoadLibrary(authorsPath, booksPath, magazinesPath)
	g.Expect(err).ToNot(HaveOccurred())
	g.Expect(lib.Authors).ToNot(BeEmpty())
	g.Expect(lib.Books).ToNot(BeEmpty())
	g.Expect(lib.Magazines).ToNot(BeEmpty())

	// Verify one author
	g.Expect(lib.Authors[0].Email).To(Equal("null-walter@echocat.org"))
	g.Expect(lib.Authors[0].FirstName).To(Equal("Paul"))

	// Verify one book
	g.Expect(lib.Books[0].ISBN).To(Equal("5554-5545-4518"))
	g.Expect(lib.Books[0].Authors).To(ContainElement("null-walter@echocat.org"))

	// Verify one magazine
	g.Expect(lib.Magazines[0].ISBN).To(Equal("5454-5587-3210"))
}
