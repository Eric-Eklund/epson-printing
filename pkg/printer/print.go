package printer

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

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
// Supports: "1" (single page), "1-5" (range), ":5" (first 5), "5:" (from 5 to end)
// Note: Comma-separated pages "1,3,5" will print only the first page listed
func convertPageRange(pageRange string) goipp.Range {
	pageRange = strings.TrimSpace(pageRange)

	// Single page: "1"
	if !strings.ContainsAny(pageRange, "-:,") {
		if page, err := strconv.Atoi(pageRange); err == nil && page > 0 {
			return goipp.Range{Lower: page, Upper: page}
		}
	}

	// Range with dash: "1-5"
	if strings.Contains(pageRange, "-") {
		parts := strings.Split(pageRange, "-")
		if len(parts) == 2 {
			lower, err1 := strconv.Atoi(strings.TrimSpace(parts[0]))
			upper, err2 := strconv.Atoi(strings.TrimSpace(parts[1]))
			if err1 == nil && err2 == nil && lower > 0 && upper > 0 {
				return goipp.Range{Lower: lower, Upper: upper}
			}
		}
	}

	// Range with colon: ":5" or "5:"
	if strings.Contains(pageRange, ":") {
		parts := strings.Split(pageRange, ":")
		if len(parts) == 2 {
			// ":5" - first 5 pages (1-5)
			if parts[0] == "" {
				if upper, err := strconv.Atoi(strings.TrimSpace(parts[1])); err == nil && upper > 0 {
					return goipp.Range{Lower: 1, Upper: upper}
				}
			}
			// "5:" - from page 5 to end
			if parts[1] == "" {
				if lower, err := strconv.Atoi(strings.TrimSpace(parts[0])); err == nil && lower > 0 {
					return goipp.Range{Lower: lower, Upper: 999}
				}
			}
		}
	}

	// Comma-separated: "1,3,5" - print only first page
	// (IPP ranges don't support non-contiguous pages well)
	if strings.Contains(pageRange, ",") {
		parts := strings.Split(pageRange, ",")
		if len(parts) > 0 {
			if page, err := strconv.Atoi(strings.TrimSpace(parts[0])); err == nil && page > 0 {
				return goipp.Range{Lower: page, Upper: page}
			}
		}
	}

	// Default: all pages
	return goipp.Range{Lower: 1, Upper: 999}
}
