package main

import (
	"testing"
	"time"

	. "github.com/onsi/gomega"
)

func TestCombineAndSortItems(t *testing.T) {
	g := NewGomegaWithT(t)

	books := []Book{
		{Title: "Zebra Book", ISBN: "123", Authors: []string{"author1"}, Description: "Desc1"},
		{Title: "Apple Book", ISBN: "456", Authors: []string{"author2"}, Description: "Desc2"},
		{Title: "Middle Book", ISBN: "789", Authors: []string{"author3"}, Description: "Desc3"},
	}

	magazines := []Magazine{
		{Title: "Beta Magazine", ISBN: "111", Authors: []string{"author4"}, PublishedAt: time.Now()},
		{Title: "Xenon Magazine", ISBN: "222", Authors: []string{"author5"}, PublishedAt: time.Now()},
		{Title: "Alpha Magazine", ISBN: "333", Authors: []string{"author6"}, PublishedAt: time.Now()},
	}

	t.Run("combines and sorts items in ascending order", func(t *testing.T) {
		items := combineAndSortItems(books, magazines, "asc")
		
		g.Expect(items).To(HaveLen(6))
		g.Expect(items[0].Title).To(Equal("Alpha Magazine"))
		g.Expect(items[0].Type).To(Equal("magazine"))
		g.Expect(items[1].Title).To(Equal("Apple Book"))
		g.Expect(items[1].Type).To(Equal("book"))
		g.Expect(items[2].Title).To(Equal("Beta Magazine"))
		g.Expect(items[3].Title).To(Equal("Middle Book"))
		g.Expect(items[4].Title).To(Equal("Xenon Magazine"))
		g.Expect(items[5].Title).To(Equal("Zebra Book"))
	})

	t.Run("combines and sorts items in descending order", func(t *testing.T) {
		items := combineAndSortItems(books, magazines, "desc")
		
		g.Expect(items).To(HaveLen(6))
		g.Expect(items[0].Title).To(Equal("Zebra Book"))
		g.Expect(items[1].Title).To(Equal("Xenon Magazine"))
		g.Expect(items[2].Title).To(Equal("Middle Book"))
		g.Expect(items[3].Title).To(Equal("Beta Magazine"))
		g.Expect(items[4].Title).To(Equal("Apple Book"))
		g.Expect(items[5].Title).To(Equal("Alpha Magazine"))
	})

	t.Run("preserves book properties", func(t *testing.T) {
		items := combineAndSortItems(books, []Magazine{}, "asc")
		
		g.Expect(items).To(HaveLen(3))
		g.Expect(items[0].Description).To(Equal("Desc2")) // Apple Book
		g.Expect(items[0].Type).To(Equal("book"))
		g.Expect(items[0].ISBN).To(Equal("456"))
		g.Expect(items[0].Authors).To(Equal([]string{"author2"}))
	})

	t.Run("preserves magazine properties", func(t *testing.T) {
		items := combineAndSortItems([]Book{}, magazines, "asc")
		
		g.Expect(items).To(HaveLen(3))
		g.Expect(items[0].Type).To(Equal("magazine"))
		g.Expect(items[0].PublishedAt).ToNot(BeZero())
	})

	t.Run("handles empty inputs", func(t *testing.T) {
		items := combineAndSortItems([]Book{}, []Magazine{}, "asc")
		g.Expect(items).To(BeEmpty())
	})
}