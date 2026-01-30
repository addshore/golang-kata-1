package library

import (
	"testing"
	. "github.com/onsi/gomega"
)

func TestLibrary_SearchAndSort(t *testing.T) {
	g := NewGomegaWithT(t)

	lib := &Library{
		Books: []Book{
			{Title: "B", ISBN: "1", Authors: []string{"a1@test.com"}},
			{Title: "A", ISBN: "2", Authors: []string{"a2@test.com"}},
		},
		Magazines: []Magazine{
			{Title: "C", ISBN: "3", Authors: []string{"a1@test.com"}},
		},
	}

	// Test FindItemByISBN
	item, found := lib.FindItemByISBN("1")
	g.Expect(found).To(BeTrue())
	g.Expect(item.(Book).Title).To(Equal("B"))

	item, found = lib.FindItemByISBN("3")
	g.Expect(found).To(BeTrue())
	g.Expect(item.(Magazine).Title).To(Equal("C"))

	_, found = lib.FindItemByISBN("non-existent")
	g.Expect(found).To(BeFalse())

	// Test FindItemsByAuthorEmail
	items := lib.FindItemsByAuthorEmail("a1@test.com")
	g.Expect(items).To(HaveLen(2))

	items = lib.FindItemsByAuthorEmail("a2@test.com")
	g.Expect(items).To(HaveLen(1))

	// Test GetSortedItems
	sorted := lib.GetSortedItems()
	g.Expect(sorted).To(HaveLen(3))
	g.Expect(getTitle(sorted[0])).To(Equal("A"))
	g.Expect(getTitle(sorted[1])).To(Equal("B"))
	g.Expect(getTitle(sorted[2])).To(Equal("C"))
}

func getTitle(item interface{}) string {
	switch v := item.(type) {
	case Book:
		return v.Title
	case Magazine:
		return v.Title
	}
	return ""
}
