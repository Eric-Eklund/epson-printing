package printer

import (
	"os"
	"testing"
)

func TestGenerateStatusReport(t *testing.T) {
	tests := []struct {
		name       string
		info       *Info
		printerURI string
		wantError  bool
	}{
		{
			name: "complete printer info with all ink levels",
			info: &Info{
				Name:         "EPSON ET-8550 Series",
				Model:        "EPSON ET-8550",
				State:        "Idle",
				StateReasons: "none",
				InkLevels: []InkLevel{
					{Name: "Matte Black", Level: 95, Color: "#000000"},
					{Name: "Photo Black", Level: 88, Color: "#000000"},
					{Name: "Cyan", Level: 72, Color: "#00FFFF"},
					{Name: "Yellow", Level: 45, Color: "#FFFF00"},
					{Name: "Magenta", Level: 18, Color: "#FF00FF"},
					{Name: "Gray", Level: 91, Color: "#808080"},
				},
			},
			printerURI: "http://localhost:631/printers/EPSON_ET-8550_Series",
			wantError:  false,
		},
		{
			name: "minimal printer info with empty ink levels",
			info: &Info{
				Name:      "Test Printer",
				Model:     "Test Model",
				State:     "Processing",
				InkLevels: []InkLevel{},
			},
			printerURI: "http://localhost:631/printers/Test",
			wantError:  false,
		},
		{
			name: "printer with low ink levels",
			info: &Info{
				Name:  "EPSON ET-8550",
				Model: "ET-8550",
				State: "Idle",
				InkLevels: []InkLevel{
					{Name: "Cyan", Level: 5, Color: "#00FFFF"},
					{Name: "Magenta", Level: 12, Color: "#FF00FF"},
					{Name: "Yellow", Level: 8, Color: "#FFFF00"},
				},
			},
			printerURI: "http://localhost:631/printers/EPSON_ET-8550_Series",
			wantError:  false,
		},
		{
			name: "printer with state reasons",
			info: &Info{
				Name:         "EPSON ET-8550",
				Model:        "ET-8550",
				State:        "Processing",
				StateReasons: "media-low",
				InkLevels: []InkLevel{
					{Name: "Black", Level: 50, Color: "#000000"},
				},
			},
			printerURI: "http://localhost:631/printers/EPSON_ET-8550_Series",
			wantError:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temporary output file
			outputPath := "test-printer-status.pdf"
			defer os.Remove(outputPath) // Clean up after test

			// Generate report
			err := GenerateStatusReport(tt.info, tt.printerURI, outputPath)

			// Check error expectation
			if (err != nil) != tt.wantError {
				t.Errorf("GenerateStatusReport() error = %v, wantError %v", err, tt.wantError)
				return
			}

			// If we expect success, verify the file was created
			if !tt.wantError {
				// Check file exists
				info, err := os.Stat(outputPath)
				if err != nil {
					t.Errorf("PDF file was not created: %v", err)
					return
				}

				// Check file is not empty
				if info.Size() == 0 {
					t.Errorf("PDF file is empty")
					return
				}

				// Check file size is reasonable (between 1KB and 100KB)
				if info.Size() < 1000 || info.Size() > 100000 {
					t.Errorf("PDF file size %d bytes is outside expected range (1000-100000)", info.Size())
				}

				t.Logf("✓ PDF generated successfully: %d bytes", info.Size())
			}
		})
	}
}

func TestGenerateStatusReport_EdgeCases(t *testing.T) {
	t.Run("empty printer name", func(t *testing.T) {
		info := &Info{
			Name:      "",
			Model:     "",
			State:     "Unknown",
			InkLevels: []InkLevel{},
		}
		outputPath := "test-empty-status.pdf"
		defer os.Remove(outputPath)

		err := GenerateStatusReport(info, "http://localhost:631", outputPath)
		if err != nil {
			t.Errorf("Should handle empty printer info: %v", err)
		}
	})

	t.Run("very long printer URI", func(t *testing.T) {
		info := &Info{
			Name:  "Test",
			Model: "Model",
			State: "Idle",
		}
		longURI := "http://very-long-printer-hostname-that-might-wrap.local:631/printers/EPSON_ET-8550_Series_With_Very_Long_Name"
		outputPath := "test-long-uri.pdf"
		defer os.Remove(outputPath)

		err := GenerateStatusReport(info, longURI, outputPath)
		if err != nil {
			t.Errorf("Should handle long URIs: %v", err)
		}
	})

	t.Run("ink levels at boundaries", func(t *testing.T) {
		info := &Info{
			Name:  "Test",
			Model: "Model",
			State: "Idle",
			InkLevels: []InkLevel{
				{Name: "Empty", Level: 0, Color: "#000000"},
				{Name: "Full", Level: 100, Color: "#00FFFF"},
				{Name: "Low Threshold", Level: 20, Color: "#FFFF00"},
				{Name: "Medium Threshold", Level: 50, Color: "#FF00FF"},
			},
		}
		outputPath := "test-boundaries.pdf"
		defer os.Remove(outputPath)

		err := GenerateStatusReport(info, "http://localhost:631", outputPath)
		if err != nil {
			t.Errorf("Should handle boundary ink levels: %v", err)
		}
	})
}

func TestDrawInkLevelBar(t *testing.T) {
	// This is a smoke test to ensure the function doesn't panic
	t.Run("ink level bar doesn't panic", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("drawInkLevelBar panicked: %v", r)
			}
		}()

		info := &Info{
			Name:  "Test",
			Model: "Model",
			State: "Idle",
			InkLevels: []InkLevel{
				{Name: "Test Ink", Level: 75, Color: "#000000"},
			},
		}
		outputPath := "test-ink-bar.pdf"
		defer os.Remove(outputPath)

		err := GenerateStatusReport(info, "http://localhost:631", outputPath)
		if err != nil {
			t.Errorf("Ink level bar rendering failed: %v", err)
		}
	})
}
