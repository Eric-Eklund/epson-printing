package printer

import (
	"encoding/json"
	"fmt"
)

// Print displays the printer information to the console in a formatted way
func (p *Info) Print() {
	fmt.Println("--- PRINTER INFORMATION ---")
	fmt.Printf("Printer Info: %s\n", p.Name)
	fmt.Printf("Model: %s\n", p.Model)

	fmt.Println("\n--- PRINTER STATUS ---")
	fmt.Printf("State: %s\n", p.State)
	fmt.Printf("State Reasons: %s\n", p.StateReasons)
	if p.StateMessage != "" {
		fmt.Printf("Message: %s\n", p.StateMessage)
	}

	fmt.Println("\n--- INK LEVELS (6-COLOR SYSTEM) ---")
	for _, ink := range p.InkLevels {
		bar := createBar(ink.Level)
		fmt.Printf("%-20s [%s] %3d%% (%s)\n", ink.Name, bar, ink.Level, ink.Color)
	}
}

// ToJSON returns the printer information as a JSON string
func (p *Info) ToJSON() (string, error) {
	data, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return "", fmt.Errorf("marshaling to JSON: %w", err)
	}
	return string(data), nil
}

// createBar creates a visual progress bar for ink levels
func createBar(level int) string {
	barLength := 20
	filled := (level * barLength) / 100
	bar := ""

	for i := 0; i < barLength; i++ {
		if i < filled {
			bar += "█"
		} else {
			bar += "░"
		}
	}

	return bar
}
