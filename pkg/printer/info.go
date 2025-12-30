package printer

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/OpenPrinting/goipp"
)

// InkLevel represents a single ink tank level
type InkLevel struct {
	Name  string `json:"name"`
	Level int    `json:"level"`
	Color string `json:"color"`
}

// Info contains all printer information and status
type Info struct {
	Name         string     `json:"name"`
	Model        string     `json:"model"`
	State        string     `json:"state"`
	StateReasons string     `json:"state_reasons"`
	StateMessage string     `json:"state_message,omitempty"`
	InkLevels    []InkLevel `json:"ink_levels"`
}

// GetPrinterInfo retrieves all printer information and status via IPP
func GetPrinterInfo(printerURI string) (*Info, error) {
	msg, err := queryPrinter(printerURI)
	if err != nil {
		return nil, err
	}

	info := &Info{
		State:     getPrinterState(msg),
		InkLevels: getInkLevels(msg),
	}

	// Get optional string attributes
	if attr := getAttribute(msg, "printer-info"); attr != nil && len(attr.Values) > 0 {
		info.Name = fmt.Sprintf("%v", attr.Values[0].V)
	}
	if attr := getAttribute(msg, "printer-make-and-model"); attr != nil && len(attr.Values) > 0 {
		info.Model = fmt.Sprintf("%v", attr.Values[0].V)
	}
	if attr := getAttribute(msg, "printer-state-reasons"); attr != nil && len(attr.Values) > 0 {
		info.StateReasons = fmt.Sprintf("%v", attr.Values[0].V)
	}
	if attr := getAttribute(msg, "printer-state-message"); attr != nil && len(attr.Values) > 0 {
		info.StateMessage = fmt.Sprintf("%v", attr.Values[0].V)
	}

	return info, nil
}

// queryPrinter sends an IPP request to get printer attributes
func queryPrinter(printerURI string) (*goipp.Message, error) {
	// Build IPP Get-Printer-Attributes request
	msg := goipp.NewRequest(goipp.DefaultVersion, goipp.OpGetPrinterAttributes, 1)
	msg.Operation.Add(goipp.MakeAttr("attributes-charset",
		goipp.TagCharset, goipp.String("utf-8")))
	msg.Operation.Add(goipp.MakeAttr("attributes-natural-language",
		goipp.TagLanguage, goipp.String("en-US")))
	msg.Operation.Add(goipp.MakeAttr("printer-uri",
		goipp.TagURI, goipp.String(printerURI)))
	msg.Operation.Add(goipp.MakeAttr("requested-attributes",
		goipp.TagKeyword, goipp.String("all")))

	// Encode request
	request, err := msg.EncodeBytes()
	if err != nil {
		return nil, fmt.Errorf("encoding request: %w", err)
	}

	// Send HTTP request
	resp, err := http.Post(printerURI, goipp.ContentType, bytes.NewBuffer(request))
	if err != nil {
		return nil, fmt.Errorf("sending HTTP request: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	// Decode response
	var respMsg goipp.Message
	err = respMsg.Decode(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	return &respMsg, nil
}

// getAttribute retrieves a specific attribute from the printer response
func getAttribute(msg *goipp.Message, name string) *goipp.Attribute {
	for _, attr := range msg.Printer {
		if attr.Name == name {
			return &attr
		}
	}
	return nil
}

// getPrinterState returns the printer state as a human-readable string
func getPrinterState(msg *goipp.Message) string {
	attr := getAttribute(msg, "printer-state")
	if attr != nil && len(attr.Values) > 0 {
		if stateVal, ok := attr.Values[0].V.(goipp.Integer); ok {
			stateInt := int(stateVal)
			switch stateInt {
			case 3:
				return "Idle"
			case 4:
				return "Processing"
			case 5:
				return "Stopped"
			default:
				return fmt.Sprintf("Unknown (%d)", stateInt)
			}
		}
	}
	return "Unknown"
}

// getInkLevels retrieves all ink tank levels from the printer
func getInkLevels(msg *goipp.Message) []InkLevel {
	var levels []InkLevel

	namesAttr := getAttribute(msg, "marker-names")
	levelsAttr := getAttribute(msg, "marker-levels")
	colorsAttr := getAttribute(msg, "marker-colors")

	if namesAttr != nil && levelsAttr != nil && colorsAttr != nil {
		for i := range namesAttr.Values {
			if i < len(levelsAttr.Values) && i < len(colorsAttr.Values) {
				name := fmt.Sprintf("%v", namesAttr.Values[i].V)
				var level int
				if levelVal, ok := levelsAttr.Values[i].V.(goipp.Integer); ok {
					level = int(levelVal)
				}
				color := fmt.Sprintf("%v", colorsAttr.Values[i].V)

				levels = append(levels, InkLevel{
					Name:  name,
					Level: level,
					Color: color,
				})
			}
		}
	}

	return levels
}
