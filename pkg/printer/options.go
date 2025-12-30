package printer

// PrintOptions contains all settings for a print job
type PrintOptions struct {
	PaperSize string // e.g., "4x6.Borderless", "A4.Borderless"
	Tray      string // e.g., "Photo", "Main", "Rear", "Auto"
	MediaType string // e.g., "photographic-glossy", "photographic-matte", "stationery"
	Quality   int    // 3=draft, 4=high, 5=best
	PageRange string // e.g., "1-5", "all", ":5", "5:"
	Copies    int    // Number of copies (default: 1)
}

// DefaultPrintOptions returns options optimized for photo printing
func DefaultPrintOptions() PrintOptions {
	return PrintOptions{
		PaperSize: "4x6.Borderless",
		Tray:      "Photo",
		MediaType: "photographic-glossy",
		Quality:   5,
		PageRange: "all",
		Copies:    1,
	}
}

// TestPrintOptions returns options for test printing on A4 plain paper
func TestPrintOptions() PrintOptions {
	return PrintOptions{
		PaperSize: "A4",
		Tray:      "Main",
		MediaType: "stationery",
		Quality:   3, // Draft quality for testing
		PageRange: "all",
		Copies:    1,
	}
}
