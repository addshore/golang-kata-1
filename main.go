package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	fmt.Println("Loading library data...")

	// Load library data from CSV files
	library, err := LoadLibrary()
	if err != nil {
		log.Fatalf("Failed to load library data: %v", err)
	}

	fmt.Printf("Loaded %d books, %d magazines, and %d authors\n",
		len(library.Books), len(library.Magazines), len(library.Authors))

	// Check command line arguments to determine interface mode
	if len(os.Args) > 1 && os.Args[1] == "cli" {
		// Start CLI interface
		fmt.Println("\nStarting Library Management System - CLI Mode...")
		cliHandler := NewCLIHandler(library)
		cliHandler.Start()
	} else {
		// Start web interface (default)
		fmt.Println("\nStarting Library Management System - Web Mode...")
		fmt.Println("Tip: Use './golang-kata-1 cli' to start in CLI mode")
		uiHandler := NewUIHandler(library)
		if err := uiHandler.StartServer("8080"); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}
}
