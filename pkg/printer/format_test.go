package printer

import (
	"encoding/json"
	"strings"
	"testing"
	"unicode/utf8"
)

func TestInfo_Print(t *testing.T) {
	info := &Info{
		Name:         "Test Printer",
		Model:        "Test Model",
		State:        "Idle",
		StateReasons: "none",
		StateMessage: "Ready",
		InkLevels: []InkLevel{
			{Name: "Black", Level: 80, Color: "#000000"},
			{Name: "Cyan", Level: 60, Color: "#00FFFF"},
		},
	}

	// Print() outputs to stdout, we can't easily capture it
	// But we can verify it doesn't panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Print() panicked: %v", r)
		}
	}()

	info.Print()
}

func TestInfo_Print_EmptyMessage(t *testing.T) {
	info := &Info{
		Name:         "Test Printer",
		Model:        "Test Model",
		State:        "Idle",
		StateReasons: "none",
		StateMessage: "", // Empty message should not be printed
		InkLevels:    []InkLevel{},
	}

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Print() panicked with empty message: %v", r)
		}
	}()

	info.Print()
}

func TestInfo_ToJSON(t *testing.T) {
	info := &Info{
		Name:         "Test Printer",
		Model:        "Test Model",
		State:        "Idle",
		StateReasons: "none",
		StateMessage: "",
		InkLevels: []InkLevel{
			{Name: "Black", Level: 80, Color: "#000000"},
			{Name: "Cyan", Level: 60, Color: "#00FFFF"},
		},
	}

	jsonStr, err := info.ToJSON()
	if err != nil {
		t.Fatalf("ToJSON() error: %v", err)
	}

	// Verify it's valid JSON
	var decoded Info
	err = json.Unmarshal([]byte(jsonStr), &decoded)
	if err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}

	// Verify values
	if decoded.Name != "Test Printer" {
		t.Errorf("expected name 'Test Printer', got %s", decoded.Name)
	}
	if decoded.State != "Idle" {
		t.Errorf("expected state 'Idle', got %s", decoded.State)
	}
	if len(decoded.InkLevels) != 2 {
		t.Errorf("expected 2 ink levels, got %d", len(decoded.InkLevels))
	}

	// Verify state_message is omitted when empty (omitempty tag)
	if strings.Contains(jsonStr, "state_message") {
		t.Error("expected state_message to be omitted when empty")
	}
}

func TestInfo_ToJSON_WithMessage(t *testing.T) {
	info := &Info{
		Name:         "Test Printer",
		Model:        "Test Model",
		State:        "Idle",
		StateReasons: "none",
		StateMessage: "Paper loaded",
		InkLevels:    []InkLevel{},
	}

	jsonStr, err := info.ToJSON()
	if err != nil {
		t.Fatalf("ToJSON() error: %v", err)
	}

	// Verify state_message is included when non-empty
	if !strings.Contains(jsonStr, "state_message") {
		t.Error("expected state_message to be included when non-empty")
	}
	if !strings.Contains(jsonStr, "Paper loaded") {
		t.Error("expected state_message value 'Paper loaded' in JSON")
	}
}

func TestCreateBar(t *testing.T) {
	tests := []struct {
		name     string
		level    int
		expected string
	}{
		{
			name:     "0% level",
			level:    0,
			expected: "░░░░░░░░░░░░░░░░░░░░",
		},
		{
			name:     "50% level",
			level:    50,
			expected: "██████████░░░░░░░░░░",
		},
		{
			name:     "100% level",
			level:    100,
			expected: "████████████████████",
		},
		{
			name:     "75% level",
			level:    75,
			expected: "███████████████░░░░░",
		},
		{
			name:     "25% level",
			level:    25,
			expected: "█████░░░░░░░░░░░░░░░",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bar := createBar(tt.level)
			if bar != tt.expected {
				t.Errorf("expected bar:\n%s\ngot:\n%s", tt.expected, bar)
			}
			runeCount := utf8.RuneCountInString(bar)
			if runeCount != 20 {
				t.Errorf("expected bar rune count 20, got %d", runeCount)
			}
		})
	}
}

func TestCreateBar_EdgeCases(t *testing.T) {
	tests := []struct {
		name  string
		level int
	}{
		{"negative level", -10},
		{"over 100%", 150},
		{"1%", 1},
		{"99%", 99},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Should not panic on edge cases
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("createBar() panicked on level %d: %v", tt.level, r)
				}
			}()

			bar := createBar(tt.level)
			runeCount := utf8.RuneCountInString(bar)
			if runeCount != 20 {
				t.Errorf("expected bar rune count 20, got %d", runeCount)
			}
		})
	}
}

func BenchmarkInfo_ToJSON(b *testing.B) {
	info := &Info{
		Name:         "Test Printer",
		Model:        "Test Model",
		State:        "Idle",
		StateReasons: "none",
		InkLevels: []InkLevel{
			{Name: "Black", Level: 80, Color: "#000000"},
			{Name: "Cyan", Level: 60, Color: "#00FFFF"},
			{Name: "Magenta", Level: 70, Color: "#FF00FF"},
			{Name: "Yellow", Level: 85, Color: "#FFFF00"},
			{Name: "Gray", Level: 90, Color: "#808080"},
			{Name: "Photo Black", Level: 95, Color: "#000000"},
		},
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = info.ToJSON()
	}
}

func BenchmarkCreateBar(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		createBar(75)
	}
}
