# Library Management System

A comprehensive library management system built in Go that provides both web and command-line interfaces for managing books, magazines, and authors. The system uses CSV files for data persistence and offers full CRUD operations with search and sorting capabilities.

## Features

- **Dual Interface**: Both web UI and interactive CLI
- **Data Management**: Books, magazines, and authors with CSV persistence
- **Search Functionality**: Search by ISBN or author email
- **Sorting**: Display items sorted by title (ascending/descending)
- **Add Items**: Create new books, magazines, and authors
- **Inline Author Creation**: Add authors while creating books/magazines
- **Unit Tested**: Comprehensive test coverage

## Quick Start

### Prerequisites
- Go 1.19 or later

### Installation
```bash
git clone <repository-url>
cd golang-kata-1
go build
```

### Running the Application

#### Web Interface (Default)
```bash
./golang-kata-1
```
Then open http://localhost:8080 in your browser.

#### CLI Interface
```bash
./golang-kata-1 cli
```

## Usage

### Web Interface

The web interface provides a user-friendly HTML interface with the following pages:

#### Home Page (/)
- Library statistics (total books, magazines, authors)
- Sample books and magazines display
- Navigation to all features

#### View All Books (/books)
- Complete list of all books with details
- Shows title, ISBN, authors, and description

#### View All Magazines (/magazines)
- Complete list of all magazines with details
- Shows title, ISBN, authors, and publication date

#### Search Library (/search)
- Search by ISBN or author email
- Dropdown to select search type
- Results show both books and magazines

#### All Items Sorted (/all-sorted)
- Combined view of books and magazines
- Sort by title in ascending or descending order
- Toggle between A-Z and Z-A sorting

#### Add New Items (/add)
- Tabbed interface for adding books, magazines, or authors
- Inline author creation when adding books/magazines
- Form validation and success/error messages

### CLI Interface

The CLI provides an interactive menu system:

```
=== MAIN MENU ===
1. Display All Books
2. Display All Magazines
3. Display All Items (Sorted)
4. Search Library
5. Add New Item
6. Library Statistics
7. Exit
```

#### Menu Options:

1. **Display All Books**: Shows all books with complete details
2. **Display All Magazines**: Shows all magazines with complete details
3. **Display All Items (Sorted)**: Combined sorted view with direction choice
4. **Search Library**: Search by ISBN or author email
5. **Add New Item**: Sub-menu for adding books, magazines, or authors
6. **Library Statistics**: Shows counts of books, magazines, and authors
7. **Exit**: Quit the application

#### Adding Items via CLI:

When adding books or magazines, you'll see a numbered list of available authors:
```
Available authors:
1. Paul Walter (null-walter@echocat.org)
2. Max Müller (null-mueller@echocat.org)
3. Franz Ferdinand (null-ferdinand@echocat.org)

Enter author numbers (comma-separated, e.g., 1,3):
```

## Data Structure

The system uses three CSV files in the `resources/` directory:

### authors.csv
```csv
email;firstname;lastname
null-walter@echocat.org;Paul;Walter
null-mueller@echocat.org;Max;Müller
```

### books.csv
```csv
title;isbn;authors;description
Book Title;1234-5678-9012;author1@example.com,author2@example.com;Book description
```

### magazines.csv
```csv
title;isbn;authors;publishedAt
Magazine Title;1234-5678-9012;author1@example.com;21.05.2011
```

## Architecture

### Core Components

- **models.go**: Data structures (Author, Book, Magazine, Library)
- **csv_parser.go**: CSV file reading and writing functionality
- **library_service.go**: Business logic layer
- **ui.go**: Web interface handlers
- **cli.go**: Command-line interface
- **main.go**: Application entry point

### Service Layer

The `LibraryService` provides a clean API for:
- Searching by ISBN or author email
- Getting sorted items
- Adding new books, magazines, and authors
- Formatting author names for display

### Data Persistence

All changes are automatically saved to CSV files:
- UTF-8 encoding with BOM support
- Semicolon-separated values
- Proper escaping for special characters

## Testing

Run the comprehensive test suite:

```bash
go test
```

The tests cover:
- CSV parsing and writing
- Service layer functionality
- Data validation
- Error handling

## API Reference

### Web Endpoints

- `GET /` - Home page
- `GET /books` - All books
- `GET /magazines` - All magazines
- `GET /search?searchType=isbn&query=1234` - Search functionality
- `GET /all-sorted?direction=asc` - Sorted items
- `GET /add` - Add items form
- `POST /add` - Process add item form

### Service Methods

```go
// Search operations
SearchByISBN(searchTerm string) ([]Book, []Magazine)
SearchByAuthorEmail(searchTerm string) ([]Book, []Magazine)

// Data retrieval
GetBooks() []Book
GetMagazines() []Magazine
GetAllItemsSorted(ascending bool) []LibraryItem

// Add operations
AddBook(title, isbn, description string, authorEmails []string) error
AddMagazine(title, isbn string, publishedAt time.Time, authorEmails []string) error
AddAuthor(email, firstName, lastName string) error
```

## Error Handling

The system includes comprehensive error handling:
- CSV parsing errors
- Duplicate ISBN validation
- Required field validation
- Author existence validation
- Date format validation

## Contributing

1. Fork the repository
2. Create a feature branch
3. Add tests for new functionality
4. Ensure all tests pass
5. Submit a pull request

## License

This project is part of a coding kata exercise.