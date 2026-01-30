package main

import (
	"sort"
)

// SearchAndSort filters the library data based on isbn or email,
// then returns a sorted list of DisplayItems if sortParam is "title".
// If no sort is requested, and no filter is applied, it returns nil (or empty slice depending on usage).
// However, the previous implementation only showed the "Unified List" if sort was requested.
// To maintain exact behavior:
// - ISBN/Email search returns a filtered LibraryData struct (subset of Books/Magazines).
// - Sort returns a slice of DisplayItems (unified list).
//
// Let's unify this. The UI logic is:
// - If .SortedItems is present, show that.
// - Else show .Books and .Magazines separately.
//
// So this function will return (LibraryData, []DisplayItem)
// But wait, the previous logic was:
// - If ISBN/Email: Filter .Books and .Magazines in place (new LibraryData instance).
// - If Sort: Create .SortedItems.
//
// We can encapsulate this into a single function that prepares the data for the view.
func ProcessRequest(data *LibraryData, isbn, email, sortParam, dirParam string) *LibraryData {
	// 1. Filtering (ISBN or Email)
	filteredData := data
	if isbn != "" {
		var filteredBooks []*Book
		for _, b := range data.Books {
			if b.ISBN == isbn {
				filteredBooks = append(filteredBooks, b)
			}
		}
		var filteredMags []*Magazine
		for _, m := range data.Magazines {
			if m.ISBN == isbn {
				filteredMags = append(filteredMags, m)
			}
		}
		filteredData = &LibraryData{
			Authors:   data.Authors,
			Books:     filteredBooks,
			Magazines: filteredMags,
		}
	} else if email != "" {
		var filteredBooks []*Book
		for _, b := range data.Books {
			for _, a := range b.Authors {
				if a.Email == email {
					filteredBooks = append(filteredBooks, b)
					break
				}
			}
		}
		var filteredMags []*Magazine
		for _, m := range data.Magazines {
			for _, a := range m.Authors {
				if a.Email == email {
					filteredMags = append(filteredMags, m)
					break
				}
			}
		}
		filteredData = &LibraryData{
			Authors:   data.Authors,
			Books:     filteredBooks,
			Magazines: filteredMags,
		}
	}

	// 2. Sorting (Unified List)
	if sortParam == "title" {
		var sortedItems []*DisplayItem
		for _, b := range filteredData.Books {
			sortedItems = append(sortedItems, &DisplayItem{
				Title:     b.Title,
				ISBN:      b.ISBN,
				Authors:   b.Authors,
				ExtraInfo: b.Description,
				Type:      "Book",
			})
		}
		for _, m := range filteredData.Magazines {
			sortedItems = append(sortedItems, &DisplayItem{
				Title:     m.Title,
				ISBN:      m.ISBN,
				Authors:   m.Authors,
				ExtraInfo: "Published: " + m.PublishedAt,
				Type:      "Magazine",
			})
		}

		sort.Slice(sortedItems, func(i, j int) bool {
			less := sortedItems[i].Title < sortedItems[j].Title
			if dirParam == "desc" {
				return !less && sortedItems[i].Title != sortedItems[j].Title
			}
			return less
		})

		// Return a new object with SortedItems populated
		return &LibraryData{
			Authors:     filteredData.Authors,
			Books:       filteredData.Books,
			Magazines:   filteredData.Magazines,
			SortedItems: sortedItems,
		}
	}

	return filteredData
}
