package printer

import (
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/go-pdf/fpdf"
)

// GenerateStatusReport creates a PDF report with printer information
func GenerateStatusReport(info *Info, printerURI, outputPath string) error {
	pdf := fpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	// Header with colored background
	pdf.SetFillColor(41, 128, 185)  // Blue background
	pdf.SetTextColor(255, 255, 255) // White text
	pdf.SetFont("Helvetica", "B", 18)
	pdf.CellFormat(0, 12, "EPSON ET-8550 STATUS REPORT", "1", 1, "C", true, 0, "")
	pdf.SetTextColor(0, 0, 0) // Reset to black

	// Date/Time with light background
	pdf.SetFillColor(236, 240, 241)
	pdf.SetFont("Helvetica", "", 10)
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	pdf.CellFormat(0, 8, fmt.Sprintf("Generated: %s", timestamp), "1", 1, "C", true, 0, "")
	pdf.Ln(5)

	// Printer Information Section
	drawSectionHeader(pdf, "PRINTER", 52, 152, 219)
	pdf.SetFont("Helvetica", "", 11)
	drawInfoRow(pdf, "Name:", info.Name)
	drawInfoRow(pdf, "Model:", info.Model)
	drawInfoRow(pdf, "Status:", info.State)
	if info.StateReasons != "" && info.StateReasons != "none" {
		drawInfoRow(pdf, "Reason:", info.StateReasons)
	}
	pdf.Ln(2)

	// Ink Levels Section with graphical bars
	drawSectionHeader(pdf, "INK LEVELS (6-COLOR SYSTEM)", 46, 204, 113)
	pdf.Ln(1)

	for _, ink := range info.InkLevels {
		drawInkLevelBar(pdf, ink.Name, ink.Level, ink.Color)
	}
	pdf.Ln(2)

	// Program Information Section
	drawSectionHeader(pdf, "PROGRAM INFORMATION", 155, 89, 182)
	pdf.SetFont("Helvetica", "", 10)
	drawInfoRow(pdf, "Command:", "print-info")
	drawInfoRow(pdf, "Go version:", runtime.Version())
	drawInfoRow(pdf, "IPP library:", "github.com/OpenPrinting/goipp v1.2.0")
	drawInfoRow(pdf, "PDF library:", "github.com/go-pdf/fpdf v0.9.0")
	pdf.Ln(2)

	// Printer URI Section
	drawSectionHeader(pdf, "PRINTER_URI", 241, 196, 15)
	pdf.SetFont("Courier", "", 9)
	pdf.SetFillColor(255, 255, 255)
	pdf.CellFormat(0, 8, printerURI, "1", 1, "L", true, 0, "")

	// Push footer to bottom if there's space
	currentY := pdf.GetY()
	footerHeight := 15.0 // Height needed for footer (line + 2 text lines + spacing)
	pageHeight := 297.0  // A4 height in mm
	bottomMargin := 20.0 // Margin from bottom
	targetY := pageHeight - bottomMargin - footerHeight

	// If we have room, push footer to bottom, otherwise place it right after content
	if currentY < targetY {
		pdf.SetY(targetY)
	} else {
		pdf.Ln(4)
	}

	// Footer with border
	pdf.SetDrawColor(189, 195, 199)
	pdf.Line(20, pdf.GetY(), 190, pdf.GetY())
	pdf.Ln(2)
	pdf.SetFont("Helvetica", "I", 9)
	pdf.SetTextColor(127, 140, 141)
	pdf.CellFormat(0, 4, "Generated with epson-printing Go tools", "", 1, "C", false, 0, "")
	pdf.SetFont("Courier", "I", 8)
	pdf.CellFormat(0, 4, "github.com/Eric-Eklund/epson-printing", "", 1, "C", false, 0, "")

	// Save PDF
	return pdf.OutputFileAndClose(outputPath)
}

// drawSectionHeader draws a colored section header
func drawSectionHeader(pdf *fpdf.Fpdf, title string, r, g, b int) {
	pdf.SetFillColor(r, g, b)
	pdf.SetTextColor(255, 255, 255)
	pdf.SetFont("Helvetica", "B", 12)
	pdf.CellFormat(0, 7, " "+title, "1", 1, "L", true, 0, "")
	pdf.SetTextColor(0, 0, 0)
	pdf.Ln(0.5)
}

// drawInfoRow draws a key-value information row
func drawInfoRow(pdf *fpdf.Fpdf, key, value string) {
	pdf.SetFont("Helvetica", "B", 10)
	pdf.CellFormat(45, 6, "  "+key, "", 0, "L", false, 0, "")
	pdf.SetFont("Helvetica", "", 10)
	pdf.CellFormat(0, 6, value, "", 1, "L", false, 0, "")
}

// drawInkLevelBar draws a graphical ink level bar
func drawInkLevelBar(pdf *fpdf.Fpdf, name string, level int, color string) {
	// Ink name
	pdf.SetFont("Helvetica", "B", 10)
	pdf.CellFormat(40, 8, "  "+name, "", 0, "L", false, 0, "")

	// Bar background
	x := pdf.GetX()
	y := pdf.GetY()
	barWidth := 100.0
	barHeight := 6.0

	// Draw background (light gray)
	pdf.SetFillColor(220, 220, 220)
	pdf.Rect(x, y+1, barWidth, barHeight, "F")

	// Draw filled portion (green gradient based on level)
	if level > 0 {
		fillWidth := (float64(level) / 100.0) * barWidth
		// Color based on level: red if low, yellow if medium, green if high
		if level < 20 {
			pdf.SetFillColor(231, 76, 60) // Red
		} else if level < 50 {
			pdf.SetFillColor(241, 196, 15) // Yellow
		} else {
			pdf.SetFillColor(46, 204, 113) // Green
		}
		pdf.Rect(x, y+1, fillWidth, barHeight, "F")
	}

	// Draw border
	pdf.SetDrawColor(100, 100, 100)
	pdf.Rect(x, y+1, barWidth, barHeight, "D")

	// Percentage text
	pdf.SetX(x + barWidth + 3)
	pdf.SetFont("Helvetica", "", 10)
	pdf.CellFormat(20, 8, fmt.Sprintf("%3d%%", level), "", 1, "R", false, 0, "")
}

// createTextBar creates a simple ASCII progress bar for PDF
func createTextBar(level int) string {
	barLength := 20
	filled := (level * barLength) / 100
	bar := ""

	for i := 0; i < barLength; i++ {
		if i < filled {
			bar += "#"
		} else {
			bar += "-"
		}
	}

	return "[" + bar + "]"
}

// PrintStatusReport generates and prints a status report
func PrintStatusReport(printerURI string) (string, int, error) {
	// Get printer information
	info, err := GetPrinterInfo(printerURI)
	if err != nil {
		return "", 0, fmt.Errorf("getting printer info: %w", err)
	}

	// Generate PDF
	pdfPath := fmt.Sprintf("printer-status-%s.pdf", time.Now().Format("20060102-150405"))
	err = GenerateStatusReport(info, printerURI, pdfPath)
	if err != nil {
		return "", 0, fmt.Errorf("generating PDF: %w", err)
	}

	// Print the PDF
	opts := PrintOptions{
		PaperSize: "A4",
		Tray:      "Main",
		MediaType: "stationery",
		Quality:   4, // High quality for reports
		PageRange: "all",
		Copies:    1,
	}

	jobID, err := PrintPDF(printerURI, pdfPath, opts)
	if err != nil {
		// Clean up PDF on print error
		os.Remove(pdfPath)
		return "", 0, fmt.Errorf("printing PDF: %w", err)
	}

	return pdfPath, jobID, nil
}
