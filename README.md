# Simple library system.

This is a simple web application to manage a library of books and magazines.

## Usage

1. Run the application:
   ```sh
   go run .
   ```
2. Open your web browser and navigate to `http://localhost:8080`.

## Features

The web interface allows you to:
- **View all books and magazines** in a single, sortable list.
- **Sort all items** by title in ascending or descending order.
- **Search** for items by their ISBN or by an author's email address.
- **Add new items** (books or magazines) to the library via a simple form.
- **Add new authors** to the library at the same time as adding a new item.
- All additions are **persisted** to the CSV files in the `resources` directory.