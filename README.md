# Simple Library System

A minimal Go application that displays books and magazines from CSV files.

## Features

- Loads library data from CSV files in the `resources/` directory
- Displays all books with title, ISBN, authors, and description
- Displays all magazines with title, ISBN, authors, and publication date
- Resolves author names from email addresses

## Usage

**List all books and magazines:**
```bash
go run main.go
```

**Search by ISBN:**
```bash
go run main.go search <ISBN>
```

**Search by author email:**
```bash
go run main.go author <email>
```

**Sort all books and magazines by title:**
```bash
go run main.go sort        # ascending (A-Z)
go run main.go sort desc   # descending (Z-A)
```

**Add new books or magazines:**
```bash
go run main.go add
```

**Start web interface:**
```bash
go run main.go web
```
Then open http://localhost:8080 in your browser.

## Data Files

- `resources/books.csv` - Book catalog
- `resources/magazines.csv` - Magazine catalog  
- `resources/authors.csv` - Author information