# Simple Library System

A comprehensive library management system with both web and CLI interfaces for managing books and magazines.

## Building

```bash
go build -o library
```

## Running

### Web Interface (Default)
```bash
./library
```
This starts the web server at `http://localhost:8080` with a modern, interactive UI.

### CLI Interface
```bash
./library -cli
```
This starts an interactive command-line interface for managing the library.

## Features

Both interfaces support the following functionality:

### 1. Display all books
- View all books in the library with complete details
- See ISBN, description, and author information

### 2. Display all magazines
- View all magazines in the library
- See ISBN, publication date, and contributor information

### 3. Display items sorted by title
- View all items (books and magazines) sorted alphabetically by title
- Supports ascending and descending order

### 4. Search by ISBN
- Find books or magazines by ISBN number
- Supports partial ISBN matching

### 5. Search by author email
- Find books and magazines by author email address
- Returns all items authored by the specified person

### 6. Add items to library
- Add new books with title, ISBN, description, and authors
- Add new magazines with title, ISBN, publication date, and contributors
- New items are automatically persisted to the CSV files
- Authors are automatically added to the authors database

### 7. Data persistence
- All data is persisted to CSV files in the `resources/` directory:
  - `authors.csv` - Author information
  - `books.csv` - Book catalog
  - `magazines.csv` - Magazine catalog

## Web Interface Features

The web interface (`http://localhost:8080`) provides:
- Beautiful, responsive design with tabbed navigation
- Real-time search and filtering
- Visual overview with statistics
- Interactive forms for adding items
- Modern card-based layout for displaying items
- Search results showing matched items

## CLI Interface Features

The CLI interface provides:
- Simple text-based menu system
- Interactive prompts for all operations
- Formatted output with clear separation
- Support for adding multiple authors per item
- Immediate feedback on all operations

## API Endpoints

The following API endpoints are available when running the web interface:

- `GET /` - Web UI
- `GET /api/books` - Get all books (JSON)
- `GET /api/magazines` - Get all magazines (JSON)
- `GET /api/library` - Get all items (JSON)
- `GET /api/books/search?isbn=<ISBN>` - Search books by ISBN
- `GET /api/magazines/search?isbn=<ISBN>` - Search magazines by ISBN
- `GET /api/books/search-by-author?email=<EMAIL>` - Search books by author email
- `GET /api/magazines/search-by-author?email=<EMAIL>` - Search magazines by author email
- `GET /api/items/sorted?direction=asc|desc` - Get all items sorted by title
- `POST /api/items/add` - Add a new book or magazine

## Data Format

### Authors CSV
```
Email;FirstName;LastName
```

### Books CSV
```
Title;ISBN;Authors;Description
```

### Magazines CSV
```
Title;ISBN;Authors;PublishedAt
```

Authors are stored as comma-separated email addresses in the books and magazines files.