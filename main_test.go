package main

import (
	"testing"

	. "github.com/onsi/gomega"
)

func TestWelcomeMessage(t *testing.T) {
	g := NewGomegaWithT(t)
	expected := "Hello world!"
	actual := welcomeMessage()
	g.Expect(actual).To(Equal(expected))
}

func TestCombineBooksAndMagazines(t *testing.T) {
	g := NewGomegaWithT(t)
	books := []Book{{Title: "A"}, {Title: "C"}}
	mags := []Magazine{{Title: "B"}}
	items := combineBooksAndMagazines(books, mags)
	g.Expect(len(items)).To(Equal(3))
	g.Expect(items[0].Title).To(Equal("A"))
	g.Expect(items[1].Title).To(Equal("C"))
	g.Expect(items[2].Title).To(Equal("B"))
}

func TestSortCombinedItemsByTitle(t *testing.T) {
	g := NewGomegaWithT(t)
	items := []CombinedItem{
		{Title: "C"}, {Title: "A"}, {Title: "B"},
	}
	sortCombinedItemsByTitle(items, true)
	g.Expect(items[0].Title).To(Equal("A"))
	g.Expect(items[1].Title).To(Equal("B"))
	g.Expect(items[2].Title).To(Equal("C"))
	sortCombinedItemsByTitle(items, false)
	g.Expect(items[0].Title).To(Equal("C"))
	g.Expect(items[1].Title).To(Equal("B"))
	g.Expect(items[2].Title).To(Equal("A"))
}

func TestContainsAuthor(t *testing.T) {
	g := NewGomegaWithT(t)
	authors := "a@b.com,c@d.com"
	g.Expect(containsAuthor(authors, "a@b.com")).To(BeTrue())
	g.Expect(containsAuthor(authors, "c@d.com")).To(BeTrue())
	g.Expect(containsAuthor(authors, "x@y.com")).To(BeFalse())
}
