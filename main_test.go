package main

import (
	"testing"
	"github.com/onsi/gomega"
)

func TestNewLibrary(t *testing.T) {
	g := gomega.NewWithT(t)
	
	authors := map[string]Author{
		"test@example.com": {Email: "test@example.com", Firstname: "Test", Lastname: "Author"},
	}
	books := []Book{
		{Title: "Test Book", ISBN: "123", Authors: []string{"test@example.com"}, Description: "Test"},
	}
	magazines := []Magazine{
		{Title: "Test Magazine", ISBN: "456", Authors: []string{"test@example.com"}, PublishedAt: "2023-01-01"},
	}
	
	library := NewLibrary(authors, books, magazines)
	
	g.Expect(len(library.Authors)).To(gomega.Equal(1))
	g.Expect(len(library.Books)).To(gomega.Equal(1))
	g.Expect(len(library.Magazines)).To(gomega.Equal(1))
}
