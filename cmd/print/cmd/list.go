package cmd

import (
	"fmt"

	"github.com/Eric-Eklund/epson-printing/pkg/printer"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all available print profiles",
	Long: `Display all available print profiles with their IDs, names, and settings.

Use profile IDs or names with the print command to quickly select
paper size, quality, and media type combinations.`,
	Example: `  # List all profiles
  print list

  # Then print using a profile ID
  print document.pdf 14`,
	Run: runList,
}

func init() {
	rootCmd.AddCommand(listCmd)
}

func runList(_ *cobra.Command, _ []string) {
	fmt.Println("Available Print Profiles:")
	fmt.Println("========================\n")

	infos := printer.ListProfilesWithInfo()

	// Group by category
	categories := []struct {
		name    string
		idStart int
		idEnd   int
	}{
		{"Default", 0, 0},
		{"4x6\" Borderless", 1, 3},
		{"5x7\" Borderless", 4, 6},
		{"A4 Borderless", 7, 9},
		{"A3 Borderless", 10, 12},
		{"A3+ Borderless (13x19\")", 13, 15},
		{"Documents", 16, 18},
	}

	for _, cat := range categories {
		if cat.name != "Default" {
			fmt.Println()
		}
		fmt.Printf("%s:\n", cat.name)

		for _, info := range infos {
			if info.ID >= cat.idStart && info.ID <= cat.idEnd {
				fmt.Printf("  %-2d  %-35s  %s, %s, %s, Quality %d\n",
					info.ID,
					info.Name,
					info.PaperSize,
					info.Tray,
					formatMediaType(info.MediaType),
					info.Quality)
			}
		}
	}

	fmt.Println("\nUsage:")
	fmt.Println("  print <file> <profile-id-or-name>")
	fmt.Println("\nExamples:")
	fmt.Println("  print document.pdf 14                               # Use profile ID 14")
	fmt.Println("  print photo.jpg photo-4x6-borderless-glossy         # Use full name")
	fmt.Println("  print document.pdf 14 --pages \"2:\"                  # Pages 2 to end")
	fmt.Println("  print calendar.pdf 7 --quality 3 --pages \"1-5\"      # Override settings")
}

func formatMediaType(media string) string {
	// Simplify media type names for display
	replacements := map[string]string{
		"photographic-glossy":     "Glossy",
		"photographic-matte":      "Matte",
		"photographic-semi-gloss": "Semi-gloss",
		"stationery":              "Plain",
		"stationery-coated":       "Coated",
	}

	if short, exists := replacements[media]; exists {
		return short
	}
	return media
}
