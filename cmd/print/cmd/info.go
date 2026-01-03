package cmd

import (
	"fmt"
	"log"

	"github.com/Eric-Eklund/epson-printing/pkg/printer"
	"github.com/spf13/cobra"
)

// infoCmd represents the info command
var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Generate and print a printer status report",
	Long: `Generate a PDF status report with printer information, ink levels,
and system details, then automatically print it.

The report includes:
  - Printer name, model, and state
  - All 6 ink tank levels with visual bars
  - Program information
  - Timestamp

The report is printed on A4 paper in normal quality.`,
	Example: `  # Generate and print status report
  print info`,
	Run: runInfo,
}

func init() {
	rootCmd.AddCommand(infoCmd)
}

func runInfo(_ *cobra.Command, _ []string) {
	if printerURI == "" {
		log.Fatal("Error: PRINTER_URI environment variable not set\n\n" +
			"Please set the printer URI:\n" +
			"  export PRINTER_URI=\"http://localhost:631/printers/EPSON_ET-8550_Series\"\n" +
			"Or use --printer flag")
	}

	fmt.Println("=============================================")
	fmt.Println("Epson ET-8550 Status Report")
	fmt.Println("=============================================\n")

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
