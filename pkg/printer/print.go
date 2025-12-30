package printer

import (
	"bytes"
	"fmt"
	"net/http"
	"os"

	"github.com/OpenPrinting/goipp"
)

// PrintPDF sends a PDF file to the printer via IPP
func PrintPDF(printerURI, pdfPath string, opts PrintOptions) (int, error) {
	// Read PDF file
	pdfData, err := os.ReadFile(pdfPath)
	if err != nil {
		return 0, fmt.Errorf("reading PDF file: %w", err)
	}

	// Get filename for job name
	filename := pdfPath
	if stat, err := os.Stat(pdfPath); err == nil {
		filename = stat.Name()
	}

	// Build IPP Print-Job request
	msg := goipp.NewRequest(goipp.DefaultVersion, goipp.OpPrintJob, 1)

	// Operation attributes
	msg.Operation.Add(goipp.MakeAttr("attributes-charset",
		goipp.TagCharset, goipp.String("utf-8")))
	msg.Operation.Add(goipp.MakeAttr("attributes-natural-language",
		goipp.TagLanguage, goipp.String("en-US")))
	msg.Operation.Add(goipp.MakeAttr("printer-uri",
		goipp.TagURI, goipp.String(printerURI)))
	msg.Operation.Add(goipp.MakeAttr("requesting-user-name",
		goipp.TagName, goipp.String(os.Getenv("USER"))))
	msg.Operation.Add(goipp.MakeAttr("job-name",
		goipp.TagName, goipp.String(filename)))
	msg.Operation.Add(goipp.MakeAttr("document-format",
		goipp.TagMimeType, goipp.String("application/pdf")))

	// Job attributes - print settings
	// Note: Using CUPS-style attribute names for compatibility
	msg.Job.Add(goipp.MakeAttr("PageSize",
		goipp.TagKeyword, goipp.String(opts.PaperSize)))
	msg.Job.Add(goipp.MakeAttr("InputSlot",
		goipp.TagKeyword, goipp.String(opts.Tray)))
	msg.Job.Add(goipp.MakeAttr("media",
		goipp.TagKeyword, goipp.String(opts.MediaType)))
	msg.Job.Add(goipp.MakeAttr("print-quality",
		goipp.TagInteger, goipp.Integer(opts.Quality)))
	msg.Job.Add(goipp.MakeAttr("copies",
		goipp.TagInteger, goipp.Integer(opts.Copies)))
	msg.Job.Add(goipp.MakeAttr("fit-to-page",
		goipp.TagBoolean, goipp.Boolean(true)))

	// Add page range if not "all"
	if opts.PageRange != "" && opts.PageRange != "all" {
		pageRangeIPP := convertPageRange(opts.PageRange)
		msg.Job.Add(goipp.MakeAttr("page-ranges",
			goipp.TagRange, pageRangeIPP))
	}

	// Encode request
	request, err := msg.EncodeBytes()
	if err != nil {
		return 0, fmt.Errorf("encoding IPP request: %w", err)
	}

	// Append PDF data to request
	requestWithPDF := append(request, pdfData...)

	// Send HTTP request
	resp, err := http.Post(printerURI, goipp.ContentType, bytes.NewBuffer(requestWithPDF))
	if err != nil {
		return 0, fmt.Errorf("sending print job: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	// Decode response
	var respMsg goipp.Message
	err = respMsg.Decode(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("decoding response: %w", err)
	}

	// Extract job ID
	jobID := 0
	if attr := getAttribute(&respMsg, "job-id"); attr != nil && len(attr.Values) > 0 {
		if idVal, ok := attr.Values[0].V.(goipp.Integer); ok {
			jobID = int(idVal)
		}
	}

	// Check for errors
	status := goipp.Status(respMsg.Code)
	if status != goipp.StatusOk && status != goipp.StatusOkIgnoredOrSubstituted {
		return jobID, fmt.Errorf("printer returned error: %s", status.String())
	}

	return jobID, nil
}

// convertPageRange converts page range notation to IPP range format
// Supports: "1-5", "1,3,5", ":5", "5:", "all"
func convertPageRange(pageRange string) goipp.Range {
	// For now, return a simple range
	// TODO: Implement full page range parsing
	return goipp.Range{Lower: 1, Upper: 999}
}

// PrintTestPDF prints the test PDF (testprint_gopher.pdf) on A4 plain paper
// This is a convenience function for testing the printer
func PrintTestPDF(printerURI, testPDFPath string) (int, error) {
	opts := TestPrintOptions()

	fmt.Println("========================================")
	fmt.Println("PRINTING TEST PDF")
	fmt.Println("========================================")
	fmt.Printf("File:        %s\n", testPDFPath)
	fmt.Printf("Paper size:  %s\n", opts.PaperSize)
	fmt.Printf("Tray:        %s\n", opts.Tray)
	fmt.Printf("Media type:  %s (plain paper)\n", opts.MediaType)
	fmt.Printf("Quality:     %d (draft)\n", opts.Quality)
	fmt.Println("========================================")

	return PrintPDF(printerURI, testPDFPath, opts)
}
