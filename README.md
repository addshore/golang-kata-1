# Library Management System

A simple Go-based library management system that provides a user interface for browsing books and magazines stored in CSV files.

## Features

- **Interactive Menu System**: Navigate through different viewing options with a user-friendly menu
- **Book Display**: View all books with title, ISBN, authors, and description
- **Magazine Display**: Browse magazines with title, ISBN, authors, and publication date
- **Combined View**: See all library items (books and magazines) in one list
- **Author Management**: Automatically resolves author information from email references
- **Smart Formatting**: Truncates long descriptions and formats author names properly

## Data Structure

The library system reads from three CSV files in the `resources/` directory:

- **`books.csv`**: Contains book information (title, ISBN, authors, description)
- **`magazines.csv`**: Contains magazine information (title, ISBN, authors, publishedAt)
- **`authors.csv`**: Contains author details (email, firstname, lastname)

## Running the Application

### Interactive Mode
```bash
go run main.go
```

This will start the interactive menu where you can:
1. Display all books
2. Display all magazines  
3. Display all items
4. Exit

### Build and Run
```bash
go build -o library main.go
./library
```

### Demo Script
```bash
./demo.sh
```

Runs a demo script that builds the project, runs tests, and provides usage instructions.

## Testing

Run the test suite with:
```bash
go test
```

## Project Structure

```
├── main.go           # Main application with interactive UI
├── main_test.go      # Test suite
├── demo.sh          # Demo script
├── go.mod           # Go module definition
├── go.sum           # Go module checksums
└── resources/       # Data files
    ├── authors.csv  # Author information
    ├── books.csv    # Book catalog
    └── magazines.csv # Magazine catalog
```

## Data Format

### Books CSV (semicolon-separated)
```
title;isbn;authors;description
```

### Magazines CSV (semicolon-separated)  
```
title;isbn;authors;publishedAt
```

### Authors CSV (semicolon-separated)
```
email;firstname;lastname
```

## Features

- **CSV Parsing**: Robust CSV parsing with semicolon delimiters
- **Author Resolution**: Automatically links author emails to full names
- **Date Handling**: Parses publication dates in DD.MM.YYYY format
- **Error Handling**: Graceful handling of missing or malformed data
- **Multiple Authors**: Support for multiple authors per publication
- **Responsive Display**: Clean, formatted output with appropriate truncation