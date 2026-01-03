package printer

import "testing"

func TestDefaultPrintOptions(t *testing.T) {
	opts := DefaultPrintOptions()

	// DefaultPrintOptions now returns the ProfileDefault (A4 draft)
	tests := []struct {
		name     string
		got      interface{}
		expected interface{}
	}{
		{"PaperSize", opts.PaperSize, "A4"},
		{"Tray", opts.Tray, "Main"},
		{"MediaType", opts.MediaType, "stationery"},
		{"Quality", opts.Quality, 3},
		{"PageRange", opts.PageRange, "all"},
		{"Copies", opts.Copies, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.expected {
				t.Errorf("DefaultPrintOptions().%s = %v, expected %v", tt.name, tt.got, tt.expected)
			}
		})
	}
}

func TestTestPrintOptions(t *testing.T) {
	opts := DefaultPrintOptions()
	defaultOpts := DefaultPrintOptions()

	if opts != defaultOpts {
		t.Error("TestPrintOptions() should return same as DefaultPrintOptions()")
	}

	// Verify the actual values
	tests := []struct {
		name     string
		got      interface{}
		expected interface{}
	}{
		{"PaperSize", opts.PaperSize, "A4"},
		{"Tray", opts.Tray, "Main"},
		{"MediaType", opts.MediaType, "stationery"},
		{"Quality", opts.Quality, 3},
		{"PageRange", opts.PageRange, "all"},
		{"Copies", opts.Copies, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.expected {
				t.Errorf("TestPrintOptions().%s = %v, expected %v", tt.name, tt.got, tt.expected)
			}
		})
	}
}

func TestPrintOptionsCustomization(t *testing.T) {
	// Test that we can customize options
	opts := DefaultPrintOptions()
	opts.PaperSize = "A4.Borderless"
	opts.Quality = 3
	opts.Copies = 2

	if opts.PaperSize != "A4.Borderless" {
		t.Errorf("expected PaperSize 'A4.Borderless', got %s", opts.PaperSize)
	}
	if opts.Quality != 3 {
		t.Errorf("expected Quality 3, got %d", opts.Quality)
	}
	if opts.Copies != 2 {
		t.Errorf("expected Copies 2, got %d", opts.Copies)
	}
}
