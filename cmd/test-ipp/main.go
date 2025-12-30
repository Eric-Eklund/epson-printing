package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Eric-Eklund/epson-printing/pkg/printer"
)

func main() {
	fmt.Println("===========================================")
	fmt.Println("Epson ET-8550 IPP Connection Test")
	fmt.Println("===========================================\n")

	// Get printer URI from environment variable
	printerURI := os.Getenv("PRINTER_URI")
	if printerURI == "" {
		log.Fatal("Error: PRINTER_URI environment variable not set\n\n" +
			"Please set the printer URI:\n" +
			"  export PRINTER_URI=\"http://localhost:631/printers/EPSON_ET-8550_Series\"\n\n" +
			"Or for network printer:\n" +
			"  export PRINTER_URI=\"ipp://your-printer.local:631/ipp/print\"\n")
	}

	fmt.Printf("Connecting to: %s\n\n", printerURI)

	// Get all printer information
	info, err := printer.GetPrinterInfo(printerURI)
	if err != nil {
		log.Fatalf("Error getting printer info: %v\n", err)
	}

	// Display printer information
	info.Print()

	fmt.Println("\n===========================================")
	fmt.Println("✓ IPP Communication successful!")
	fmt.Println("===========================================")
}
