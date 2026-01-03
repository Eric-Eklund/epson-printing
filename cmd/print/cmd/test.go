package cmd

import (
	"fmt"
	"log"

	"github.com/Eric-Eklund/epson-printing/pkg/printer"
	"github.com/spf13/cobra"
)

var jsonOutput bool

// testCmd represents the test command
var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Test IPP connection and display printer status",
	Long: `Test the IPP connection to the printer and display current status
including printer state, ink levels for all 6 tanks, and any messages.

Use --json flag to output in JSON format instead of formatted text.`,
	Example: `  # Test connection and show status
  print test

  # Output in JSON format
  print test --json`,
	Run: runTest,
}

func init() {
	rootCmd.AddCommand(testCmd)
	testCmd.Flags().BoolVar(&jsonOutput, "json", false, "Output in JSON format")
}

func runTest(_ *cobra.Command, _ []string) {
	if printerURI == "" {
		log.Fatal("Error: PRINTER_URI environment variable not set\n\n" +
			"Please set the printer URI:\n" +
			"  export PRINTER_URI=\"http://localhost:631/printers/EPSON_ET-8550_Series\"\n" +
			"Or use --printer flag")
	}

	if !jsonOutput {
		fmt.Println("=============================================")
		fmt.Println("Epson ET-8550 IPP Connection Test")
		fmt.Println("=============================================")
		fmt.Printf("Testing connection to: %s\n\n", printerURI)
	}

	// Get printer information
	info, err := printer.GetPrinterInfo(printerURI)
	if err != nil {
		log.Fatalf("Error: Failed to connect to printer\n%v\n", err)
	}

	// Output
	if jsonOutput {
		jsonStr, err := info.ToJSON()
		if err != nil {
			log.Fatalf("Error: Failed to generate JSON\n%v\n", err)
		}
		fmt.Println(jsonStr)
	} else {
		info.Print()
		fmt.Println("\n✓ Connection successful!")
	}
}
