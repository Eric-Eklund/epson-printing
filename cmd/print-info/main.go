package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Eric-Eklund/epson-printing/pkg/printer"
)

func main() {
	fmt.Println("=============================================")
	fmt.Println("Epson ET-8550 Status Report")
	fmt.Println("=============================================\n")

	// Get printer URI from environment variable
	printerURI := os.Getenv("PRINTER_URI")
	if printerURI == "" {
		log.Fatal("Error: PRINTER_URI environment variable not set\n\n" +
			"Please set the printer URI:\n" +
			"  export PRINTER_URI=\"http://localhost:631/printers/EPSON_ET-8550_Series\"\n\n" +
			"Or for network printer:\n" +
			"  export PRINTER_URI=\"ipp://your-printer.local:631/ipp/print\"\n")
	}

	fmt.Printf("Fetching printer information from: %s\n", printerURI)
	fmt.Println("Generating PDF report...")

	// Generate and print the status report
	pdfPath, jobID, err := printer.PrintStatusReport(printerURI)
	if err != nil {
		log.Fatalf("Error: %v\n", err)
	}

	fmt.Println("\n✓ PDF report created:", pdfPath)
	fmt.Printf("✓ Print job sent! (Job ID: %d)\n", jobID)
	fmt.Println("\nReport contains:")
	fmt.Println("  - Printer information")
	fmt.Println("  - Ink levels for all 6 tanks")
	fmt.Println("  - Program information")
	fmt.Println("  - Date and time")
	fmt.Println("\nPDF file saved for reference:", pdfPath)
}
