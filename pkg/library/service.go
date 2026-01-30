package library

import "sort"

func (l *Library) FindItemByISBN(isbn string) (interface{}, bool) {
	for _, book := range l.Books {
		if book.ISBN == isbn {
			return book, true
		}
	}
	for _, magazine := range l.Magazines {
		if magazine.ISBN == isbn {
			return magazine, true
		}
	}
	return nil, false
}

func (l *Library) FindItemsByAuthorEmail(email string) []interface{} {
	var results []interface{}
	for _, book := range l.Books {
		for _, authorEmail := range book.Authors {
			if authorEmail == email {
				results = append(results, book)
				break
			}
		}
	}
	for _, magazine := range l.Magazines {
		for _, authorEmail := range magazine.Authors {
			if authorEmail == email {
				results = append(results, magazine)
				break
			}
		}
	}
	return results
}

func (l *Library) GetSortedItems() []interface{} {
	var items []interface{}
	for _, book := range l.Books {
		items = append(items, book)
	}
	for _, magazine := range l.Magazines {
		items = append(items, magazine)
	}

	sort.Slice(items, func(i, j int) bool {
		var titleI, titleJ string
		switch v := items[i].(type) {
		case Book:
			titleI = v.Title
		case Magazine:
			titleI = v.Title
		}
		switch v := items[j].(type) {
		case Book:
			titleJ = v.Title
		case Magazine:
			titleJ = v.Title
		}
		return titleI < titleJ
	})

	return items
}
