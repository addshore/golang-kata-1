#!/bin/bash

echo "=== Library System Demo Script ==="
echo "This script demonstrates the library application functionality"
echo ""

echo "1. Building the application..."
go build -o library main.go
if [ $? -ne 0 ]; then
    echo "Build failed!"
    exit 1
fi

echo "2. Running tests..."
go test
if [ $? -ne 0 ]; then
    echo "Tests failed!"
    exit 1
fi

echo ""
echo "3. Application is ready! You can now run:"
echo "   ./library    # Interactive mode"
echo "   go run main.go    # Alternative way to run"

echo ""
echo "4. The interactive menu offers these options:"
echo "   - Display all books (8 books in the catalog)"
echo "   - Display all magazines (6 magazines in the catalog)" 
echo "   - Display all items combined"
echo "   - Exit"

echo ""
echo "5. Sample data includes:"
echo "   - Books: Cooking guides, technical books, recipes"
echo "   - Magazines: Cooking and wine publications"
echo "   - Authors: 6 authors with German names"

echo ""
echo "To start the application, run: ./library"
