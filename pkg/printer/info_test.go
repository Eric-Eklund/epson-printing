package printer

import (
	"testing"

	"github.com/OpenPrinting/goipp"
)

// Test helper: create a mock IPP message with printer attributes
func createMockMessage() *goipp.Message {
	msg := &goipp.Message{
		Printer: []goipp.Attribute{
			goipp.MakeAttr("printer-info", goipp.TagText, goipp.String("Test Printer")),
			goipp.MakeAttr("printer-make-and-model", goipp.TagText, goipp.String("Test Model")),
			goipp.MakeAttr("printer-state", goipp.TagEnum, goipp.Integer(3)), // Idle
			goipp.MakeAttr("printer-state-reasons", goipp.TagKeyword, goipp.String("none")),
			goipp.MakeAttr("marker-names", goipp.TagName, goipp.String("Black"), goipp.String("Cyan")),
			goipp.MakeAttr("marker-levels", goipp.TagInteger, goipp.Integer(80), goipp.Integer(60)),
			goipp.MakeAttr("marker-colors", goipp.TagName, goipp.String("#000000"), goipp.String("#00FFFF")),
		},
	}
	return msg
}

func TestGetAttribute(t *testing.T) {
	msg := createMockMessage()

	tests := []struct {
		name          string
		attributeName string
		expectNil     bool
		expectedValue string
	}{
		{
			name:          "existing attribute",
			attributeName: "printer-info",
			expectNil:     false,
			expectedValue: "Test Printer",
		},
		{
			name:          "non-existing attribute",
			attributeName: "non-existent",
			expectNil:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			attr := getAttribute(msg, tt.attributeName)
			if tt.expectNil {
				if attr != nil {
					t.Errorf("expected nil, got %v", attr)
				}
			} else {
				if attr == nil {
					t.Fatal("expected attribute, got nil")
				}
				if attr.Name != tt.attributeName {
					t.Errorf("expected name %s, got %s", tt.attributeName, attr.Name)
				}
			}
		})
	}
}

func TestGetPrinterState(t *testing.T) {
	tests := []struct {
		name          string
		stateValue    int
		expectedState string
	}{
		{
			name:          "idle state",
			stateValue:    3,
			expectedState: "Idle",
		},
		{
			name:          "processing state",
			stateValue:    4,
			expectedState: "Processing",
		},
		{
			name:          "stopped state",
			stateValue:    5,
			expectedState: "Stopped",
		},
		{
			name:          "unknown state",
			stateValue:    99,
			expectedState: "Unknown (99)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := &goipp.Message{
				Printer: []goipp.Attribute{
					goipp.MakeAttr("printer-state", goipp.TagEnum, goipp.Integer(tt.stateValue)),
				},
			}

			state := getPrinterState(msg)
			if state != tt.expectedState {
				t.Errorf("expected state %s, got %s", tt.expectedState, state)
			}
		})
	}

	// Test missing printer-state attribute
	t.Run("missing state attribute", func(t *testing.T) {
		msg := &goipp.Message{Printer: []goipp.Attribute{}}
		state := getPrinterState(msg)
		if state != "Unknown" {
			t.Errorf("expected 'Unknown', got %s", state)
		}
	})
}

func TestGetInkLevels(t *testing.T) {
	msg := createMockMessage()
	levels := getInkLevels(msg)

	if len(levels) != 2 {
		t.Fatalf("expected 2 ink levels, got %d", len(levels))
	}

	// Test first ink level
	if levels[0].Name != "Black" {
		t.Errorf("expected name 'Black', got %s", levels[0].Name)
	}
	if levels[0].Level != 80 {
		t.Errorf("expected level 80, got %d", levels[0].Level)
	}
	if levels[0].Color != "#000000" {
		t.Errorf("expected color '#000000', got %s", levels[0].Color)
	}

	// Test second ink level
	if levels[1].Name != "Cyan" {
		t.Errorf("expected name 'Cyan', got %s", levels[1].Name)
	}
	if levels[1].Level != 60 {
		t.Errorf("expected level 60, got %d", levels[1].Level)
	}
}

func TestGetInkLevels_MissingAttributes(t *testing.T) {
	msg := &goipp.Message{Printer: []goipp.Attribute{}}
	levels := getInkLevels(msg)

	if len(levels) != 0 {
		t.Errorf("expected empty slice, got %d levels", len(levels))
	}
}

func BenchmarkGetPrinterState(b *testing.B) {
	msg := createMockMessage()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		getPrinterState(msg)
	}
}

func BenchmarkGetInkLevels(b *testing.B) {
	msg := createMockMessage()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		getInkLevels(msg)
	}
}
