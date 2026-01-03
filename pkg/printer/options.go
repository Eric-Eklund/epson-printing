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

// DefaultPrintOptions returns the default profile options (test/draft on A4)
// This is the default when no profile is specified
func DefaultPrintOptions() PrintOptions {
	return MustGetPrintOptions(ProfileDefault)
}
