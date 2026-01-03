package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Eric-Eklund/epson-printing/pkg/printer"
)

const testPDF = "testprint_gopher.pdf"

func main() {
	// Get printer URI from environment variable
	printerURI := os.Getenv("PRINTER_URI")
	if printerURI == "" {
		log.Fatal("Error: PRINTER_URI environment variable not set\n\n" +
			"Please set the printer URI:\n" +
			"  export PRINTER_URI=\"http://localhost:631/printers/EPSON_ET-8550_Series\"\n\n" +
			"Or for network printer:\n" +
			"  export PRINTER_URI=\"ipp://your-printer.local:631/ipp/print\"\n")
	}

	// Check if test PDF exists
	if _, err := os.Stat(testPDF); os.IsNotExist(err) {
		log.Fatalf("Error: Test PDF not found: %s\n", testPDF)
	}

	fmt.Println("Epson ET-8550 Test Print")
	fmt.Println("This will print testprint_gopher.pdf on A4 plain paper")
	fmt.Printf("Using printer: %s\n\n", printerURI)

	// Example: Using the default profile (document-draft)
	// You can also use other profiles like:
	//   printer.ProfilePhoto4x6BorderlessGlossy
	//   printer.ProfilePhotoA3PlusBorderlessMatte
	//   printer.ProfileDocumentBest
	profile := printer.ProfileDefault
	fmt.Printf("Using profile: %s\n", profile)
	fmt.Printf("Settings: %s\n\n", printer.GetProfileDescription(profile))

	// Get options from profile
	opts, err := printer.GetPrintOptions(profile)
	if err != nil {
		log.Fatalf("Error getting profile: %v\n", err)
	}

	// Print the PDF
	jobID, err := printer.PrintPDF(printerURI, testPDF, opts)
	if err != nil {
		log.Fatalf("Error: %v\n", err)
	}

	fmt.Printf("✓ Print job sent successfully! (Job ID: %d)\n", jobID)
}
