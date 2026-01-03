package cmd

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/Eric-Eklund/epson-printing/pkg/printer"
	"github.com/spf13/cobra"
)

var (
	// Persistent flags (available to all commands)
	printerURI string

	// Print command flags
	profileFlag string
	pagesFlag   string
	qualityFlag int
	paperFlag   string
	trayFlag    string
	mediaFlag   string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "print <file> [profile]",
	Short: "Print files with profile-based settings",
	Long: `Print files (PDF, images, documents) to your Epson ET-8550 printer
using predefined profiles or custom settings.

Profile can be specified as:
  - Numeric ID (0-18): print file.pdf 14
  - Profile name: print file.pdf photo-a3plus-borderless-matte
  - Default (0) if not specified

Use 'print list' to see all available profiles.`,
	Example: `  # Print with profile ID
  print document.pdf 14
  print photo.jpg 1

  # Print with profile name
  print calendar.pdf photo-a3plus-borderless-matte

  # Override settings
  print document.pdf 14 --pages "2:"
  print document.pdf 7 --quality 3`,
	Args: cobra.MinimumNArgs(1),
	Run:  runPrint,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	// Persistent flags (available to all subcommands)
	rootCmd.PersistentFlags().StringVar(&printerURI, "printer", os.Getenv("PRINTER_URI"),
		"Printer URI (default from PRINTER_URI env var)")

	// Local flags (only for print command)
	rootCmd.Flags().StringVarP(&profileFlag, "profile", "p", "",
		"Profile name or ID (alternative to positional)")
	rootCmd.Flags().StringVar(&pagesFlag, "pages", "",
		"Page range: '1', '1-5', '2:' (from 2), ':5' (to 5)")
	rootCmd.Flags().IntVarP(&qualityFlag, "quality", "q", 0,
		"Quality: 3 (draft), 4 (normal), 5 (best)")
	rootCmd.Flags().StringVar(&paperFlag, "paper", "",
		"Paper size override")
	rootCmd.Flags().StringVar(&trayFlag, "tray", "",
		"Tray override: Photo, Main, Rear, Auto")
	rootCmd.Flags().StringVar(&mediaFlag, "media", "",
		"Media type override")
}

func runPrint(_ *cobra.Command, args []string) {
	// Check printer URI
	if printerURI == "" {
		log.Fatal("Error: PRINTER_URI environment variable not set\n\n" +
			"Please set the printer URI:\n" +
			"  export PRINTER_URI=\"http://localhost:631/printers/EPSON_ET-8550_Series\"\n" +
			"Or use --printer flag")
	}

	// Get file and profile from arguments
	pdfFile := args[0]
	profile := profileFlag

	// If no --profile flag, check for positional profile argument
	if profile == "" && len(args) >= 2 {
		profile = args[1]
	}

	// Default profile if none specified
	if profile == "" {
		profile = "default"
	}

	// Get print options
	opts, err := getOptionsFromProfile(profile)
	if err != nil {
		log.Fatalf("Error: %v\n", err)
	}

	// Apply overrides
	if pagesFlag != "" {
		opts.PageRange = pagesFlag
	}
	if qualityFlag > 0 {
		if qualityFlag < 3 || qualityFlag > 5 {
			log.Fatal("Error: Quality must be 3 (draft), 4 (normal), or 5 (best)")
		}
		opts.Quality = qualityFlag
	}
	if paperFlag != "" {
		opts.PaperSize = paperFlag
	}
	if trayFlag != "" {
		opts.Tray = trayFlag
	}
	if mediaFlag != "" {
		opts.MediaType = mediaFlag
	}

	// Print info
	fmt.Println("=========================================")
	fmt.Println("PDF PRINT")
	fmt.Println("=========================================")
	fmt.Printf("File:        %s\n", pdfFile)
	fmt.Printf("Profile:     %s\n", profile)
	fmt.Printf("Paper size:  %s\n", opts.PaperSize)
	fmt.Printf("Tray:        %s\n", opts.Tray)
	fmt.Printf("Media type:  %s\n", opts.MediaType)
	fmt.Printf("Quality:     %d (3=draft, 4=normal, 5=best)\n", opts.Quality)
	fmt.Printf("Pages:       %s\n", opts.PageRange)
	fmt.Println("=========================================")
	fmt.Println()

	// Print the PDF
	jobID, err := printer.PrintPDF(printerURI, pdfFile, opts)
	if err != nil {
		log.Fatalf("Print failed: %v\n", err)
	}

	fmt.Printf("✓ Print job sent successfully! (Job ID: %d)\n", jobID)
}

func getOptionsFromProfile(profile string) (printer.PrintOptions, error) {
	// Try to parse as numeric ID first
	if id, err := strconv.Atoi(profile); err == nil {
		profileName, err := printer.GetProfileByID(id)
		if err != nil {
			return printer.PrintOptions{}, err
		}
		return printer.GetPrintOptions(profileName)
	}

	// Otherwise treat as profile name
	return printer.GetPrintOptions(printer.PrintProfile(profile))
}
