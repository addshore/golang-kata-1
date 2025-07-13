package main

import (
	"testing"

	. "github.com/onsi/gomega"
)

func TestFindByISBN(t *testing.T) {
	g := NewGomegaWithT(t)
	books := []Book{{Title: "A", ISBN: "1"}, {Title: "B", ISBN: "2"}}
	magazines := []Magazine{{Title: "M1", ISBN: "3"}}

	b, m := FindByISBN(books, magazines, "1")
	g.Expect(b).NotTo(BeNil())
	g.Expect(b.Title).To(Equal("A"))
	g.Expect(m).To(BeNil())

	b, m = FindByISBN(books, magazines, "3")
	g.Expect(b).To(BeNil())
	g.Expect(m).NotTo(BeNil())
	g.Expect(m.Title).To(Equal("M1"))

	b, m = FindByISBN(books, magazines, "notfound")
	g.Expect(b).To(BeNil())
	g.Expect(m).To(BeNil())
}

func TestFindByAuthorEmail(t *testing.T) {
	g := NewGomegaWithT(t)
	books := []Book{{Title: "A", Authors: []string{"a@b.com"}}, {Title: "B", Authors: []string{"c@d.com"}}}
	magazines := []Magazine{{Title: "M1", Authors: []string{"a@b.com"}}}

	b, m := FindByAuthorEmail(books, magazines, "a@b.com")
	g.Expect(len(b)).To(Equal(1))
	g.Expect(b[0].Title).To(Equal("A"))
	g.Expect(len(m)).To(Equal(1))
	g.Expect(m[0].Title).To(Equal("M1"))

	b, m = FindByAuthorEmail(books, magazines, "c@d.com")
	g.Expect(len(b)).To(Equal(1))
	g.Expect(b[0].Title).To(Equal("B"))
	g.Expect(len(m)).To(Equal(0))
}

func TestSortByTitle(t *testing.T) {
	g := NewGomegaWithT(t)
	books := []Book{{Title: "B"}, {Title: "A"}}
	magazines := []Magazine{{Title: "C"}}

	allAsc := SortByTitle(books, magazines, true)
	g.Expect(allAsc[0].(Book).Title).To(Equal("A"))
	g.Expect(allAsc[1].(Book).Title).To(Equal("B"))
	g.Expect(allAsc[2].(Magazine).Title).To(Equal("C"))

	allDesc := SortByTitle(books, magazines, false)
	g.Expect(allDesc[0].(Magazine).Title).To(Equal("C"))
	g.Expect(allDesc[1].(Book).Title).To(Equal("B"))
	g.Expect(allDesc[2].(Book).Title).To(Equal("A"))
}
