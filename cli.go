package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func runCLI() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Library CLI")
	fmt.Println("Type a number and press Enter.")

	for {
		fmt.Println("")
		fmt.Println("1) List all items")
		fmt.Println("2) Search by ISBN")
		fmt.Println("3) Search by author email")
		fmt.Println("4) Add to library")
		fmt.Println("5) Quit")
		fmt.Print("> ")

		choice := readLine(reader)
		switch choice {
		case "1":
			data, err := loadLibraryData()
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				continue
			}
			items := buildCombinedItems(data.Books, data.Magazines, data.Authors)
			printItems(items)
		case "2":
			fmt.Print("ISBN: ")
			isbn := readLine(reader)
			data, err := loadLibraryData()
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				continue
			}
			books := filterBooksByISBN(data.Books, isbn)
			mags := filterMagazinesByISBN(data.Magazines, isbn)
			items := buildCombinedItems(books, mags, data.Authors)
			printItems(items)
		case "3":
			fmt.Print("Author email: ")
			email := readLine(reader)
			data, err := loadLibraryData()
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				continue
			}
			books := filterBooksByAuthorEmail(data.Books, email)
			mags := filterMagazinesByAuthorEmail(data.Magazines, email)
			items := buildCombinedItems(books, mags, data.Authors)
			printItems(items)
		case "4":
			if err := addFromCLI(reader); err != nil {
				fmt.Printf("Error: %v\n", err)
			} else {
				fmt.Println("Library updated successfully.")
			}
		case "5":
			fmt.Println("Goodbye!")
			return
		default:
			fmt.Println("Unknown option.")
		}
	}
}

func addFromCLI(reader *bufio.Reader) error {
	data, err := loadLibraryData()
	if err != nil {
		return err
	}

	fmt.Print("Add type (book/magazine/author): ")
	itemType := strings.ToLower(readLine(reader))

	request := AddRequest{ItemType: itemType}

	if itemType == "book" || itemType == "magazine" {
		fmt.Print("Title: ")
		request.Title = readLine(reader)
		fmt.Print("ISBN: ")
		request.ISBN = readLine(reader)
		fmt.Print("Author emails (comma separated): ")
		request.AuthorEmails = splitAuthors(readLine(reader))
		if itemType == "book" {
			fmt.Print("Description: ")
			request.Description = readLine(reader)
		}
		if itemType == "magazine" {
			fmt.Print("Published date (DD.MM.YYYY): ")
			request.PublishedAt = readLine(reader)
		}
	}

	fmt.Print("Add new author? (y/N): ")
	if strings.EqualFold(readLine(reader), "y") {
		fmt.Print("Author email: ")
		request.NewAuthorEmail = readLine(reader)
		fmt.Print("First name: ")
		request.NewAuthorFirst = readLine(reader)
		fmt.Print("Last name: ")
		request.NewAuthorLast = readLine(reader)
	}

	plan, err := validateAddRequest(data, request)
	if err != nil {
		return err
	}
	return applyAddPlan(plan)
}

func readLine(reader *bufio.Reader) string {
	line, _ := reader.ReadString('\n')
	return strings.TrimSpace(line)
}

func printItems(items []ItemView) {
	if len(items) == 0 {
		fmt.Println("No items found.")
		return
	}
	for _, item := range items {
		fmt.Printf("%s [%s]\n", item.Title, item.ItemType)
		fmt.Printf("  ISBN: %s\n", item.ISBN)
		fmt.Printf("  Authors: %s\n", item.Authors)
		if item.PublishedAt != "" {
			fmt.Printf("  Published: %s\n", item.PublishedAt)
		}
		if item.Description != "" {
			fmt.Printf("  Description: %s\n", item.Description)
		}
		fmt.Println("")
	}
}
