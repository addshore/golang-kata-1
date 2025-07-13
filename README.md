# Library Application

This is a simple library application with both a web and CLI interface. It uses CSV files in the `resources/` directory to store books, magazines, and authors.

## Features

- Display all books and magazines with details
- Search by ISBN or author email
- Print all items sorted by title
- Add books, magazines, and authors (with correct CSV updates)
- Web UI and interactive CLI

## Running the Web Server

**Important:**  
To avoid Go's "multiple main functions" error, ensure you only run the web server files:

```sh
go run main.go library.go
```

Open your browser and visit:

- [http://localhost:8080/](http://localhost:8080/) — Home page
- `/books` — List/search books
- `/magazines` — List/search magazines
- `/all-sorted` — View all sorted by title
- `/add-item` — Add a book or magazine

## Using the CLI

**Important:**  
To avoid Go's "multiple main functions" error, ensure you only run the CLI files:

```sh
go run cli.go library.go
```

You will see a menu:

```
Library CLI
1. List all books
2. List all magazines
3. Search by ISBN
4. Search by author email
5. List all sorted by title
6. Add book or magazine
0. Exit
Choose an option:
```

### Example CLI Usage

- **List all books:** Enter `1`
- **List all magazines:** Enter `2`
- **Search by ISBN:** Enter `3`, then type the ISBN
- **Search by author email:** Enter `4`, then type the email
- **List all sorted by title:** Enter `5`
- **Add book or magazine:** Enter `6` and follow the prompts

All data is persisted to the CSV files in `resources/`.

## Requirements

- Go 1.18+
- The `resources/` directory with `books.csv`, `magazines.csv`, and `authors.csv` present
