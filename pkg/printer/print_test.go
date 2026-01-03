package printer

import (
	"os"
	"testing"
)

func TestConvertPageRange(t *testing.T) {
	tests := []struct {
		name      string
		pageRange string
		wantLower int
		wantUpper int
	}{
		{
			name:      "all pages",
			pageRange: "all",
			wantLower: 1,
			wantUpper: 999,
		},
		{
			name:      "single page",
			pageRange: "1",
			wantLower: 1,
			wantUpper: 1,
		},
		{
			name:      "single page 5",
			pageRange: "5",
			wantLower: 5,
			wantUpper: 5,
		},
		{
			name:      "range 1-5",
			pageRange: "1-5",
			wantLower: 1,
			wantUpper: 5,
		},
		{
			name:      "range 3-10",
			pageRange: "3-10",
			wantLower: 3,
			wantUpper: 10,
		},
		{
			name:      "first 5 pages",
			pageRange: ":5",
			wantLower: 1,
			wantUpper: 5,
		},
		{
			name:      "from page 5",
			pageRange: "5:",
			wantLower: 5,
			wantUpper: 999,
		},
		{
			name:      "specific pages (comma) - only first",
			pageRange: "1,3,5",
			wantLower: 1,
			wantUpper: 1,
		},
		{
			name:      "empty string",
			pageRange: "",
			wantLower: 1,
			wantUpper: 999,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := convertPageRange(tt.pageRange)

			if result.Lower != tt.wantLower {
				t.Errorf("convertPageRange(%q).Lower = %d, want %d", tt.pageRange, result.Lower, tt.wantLower)
			}
			if result.Upper != tt.wantUpper {
				t.Errorf("convertPageRange(%q).Upper = %d, want %d", tt.pageRange, result.Upper, tt.wantUpper)
			}
		})
	}
}

func TestPrintOptions_Validation(t *testing.T) {
	tests := []struct {
		name    string
		opts    PrintOptions
		wantErr bool
	}{
		{
			name: "valid options",
			opts: PrintOptions{
				PaperSize: "A4",
				Tray:      "Main",
				MediaType: "stationery",
				Quality:   4,
				PageRange: "all",
				Copies:    1,
			},
			wantErr: false,
		},
		{
			name: "valid borderless",
			opts: PrintOptions{
				PaperSize: "4x6.Borderless",
				Tray:      "Photo",
				MediaType: "photographic-glossy",
				Quality:   5,
				PageRange: "all",
				Copies:    1,
			},
			wantErr: false,
		},
		{
			name: "multiple copies",
			opts: PrintOptions{
				PaperSize: "A4",
				Tray:      "Main",
				MediaType: "stationery",
				Quality:   3,
				PageRange: "1-10",
				Copies:    5,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Just verify options can be created without panic
			// Actual validation would happen in PrintPDF
			_ = tt.opts
		})
	}
}

// TestPrintPDF_Integration tests actual printing to a real printer
// This test is skipped by default and only runs when explicitly enabled
func TestPrintPDF_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	if os.Getenv("INTEGRATION_TEST") == "" {
		t.Skip("Skipping integration test (set INTEGRATION_TEST=1 to run)")
	}

	// This would require a real printer and test PDF
	// For now, just verify the function signature is correct
	t.Skip("Integration test not implemented - requires real printer")
}

// TestPrintTestPDF_Integration tests the test print functionality
func TestPrintTestPDF_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	if os.Getenv("INTEGRATION_TEST") == "" {
		t.Skip("Skipping integration test (set INTEGRATION_TEST=1 to run)")
	}

	t.Skip("Integration test not implemented - requires real printer and test PDF")
}

// Benchmark convertPageRange
func BenchmarkConvertPageRange(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		convertPageRange("1-5")
	}
}
