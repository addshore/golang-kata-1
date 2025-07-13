package main

// FilterBooksByISBN filters books by matching ISBN
func FilterBooksByISBN(books []Book, isbn string) []Book {
	if isbn == "" {
		return books
	}
	
	var filtered []Book
	for _, book := range books {
		if book.ISBN == isbn {
			filtered = append(filtered, book)
		}
	}
	return filtered
}

// FilterMagazinesByISBN filters magazines by matching ISBN
func FilterMagazinesByISBN(magazines []Magazine, isbn string) []Magazine {
	if isbn == "" {
		return magazines
	}
	
	var filtered []Magazine
	for _, magazine := range magazines {
		if magazine.ISBN == isbn {
			filtered = append(filtered, magazine)
		}
	}
	return filtered
}

// FilterBooksByAuthor filters books by author email
func FilterBooksByAuthor(books []Book, authorEmail string) []Book {
	if authorEmail == "" {
		return books
	}
	
	var filtered []Book
	for _, book := range books {
		for _, email := range book.Authors {
			if email == authorEmail {
				filtered = append(filtered, book)
				break
			}
		}
	}
	return filtered
}

// FilterMagazinesByAuthor filters magazines by author email
func FilterMagazinesByAuthor(magazines []Magazine, authorEmail string) []Magazine {
	if authorEmail == "" {
		return magazines
	}
	
	var filtered []Magazine
	for _, magazine := range magazines {
		for _, email := range magazine.Authors {
			if email == authorEmail {
				filtered = append(filtered, magazine)
				break
			}
		}
	}
	return filtered
}