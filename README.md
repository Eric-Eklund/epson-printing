<div align="center">
  <img src="gopher-printer.png" alt="Gopher Printer" width="300">
</div>

# Epson ET-8550 Printing Tools (Go)

Go-based tools for controlling Epson ET-8550 EcoTank printer via IPP (Internet Printing Protocol).

[![Go Version](https://img.shields.io/badge/Go-1.18+-00ADD8?style=flat&logo=go)](https://go.dev)
[![Tests](https://img.shields.io/badge/tests-passing-brightgreen)](./pkg/printer)

---

## Features

- 🖨️ **Printer Status Monitoring** - Get printer state, model, and detailed information
- 📊 **Ink Level Monitoring** - Read all 6 ink tank levels with visual displays
- 📄 **PDF Printing** - Print PDFs with complete control over all settings
- 📋 **Status Reports** - Generate styled PDF reports with printer information
- ✅ **Comprehensive Testing** - 34 tests covering all major functionality
- 🔌 **Direct IPP Communication** - No external dependencies beyond Go libraries

---

## Quick Start

### Prerequisites

- Go 1.18 or later
- Epson ET-8550 printer configured in CUPS or accessible via network

### Installation

```bash
# Clone the repository
git clone https://github.com/Eric-Eklund/epson-printing.git
cd epson-printing

# Install dependencies
go mod download

# Build all tools
go build -o bin/test-ipp ./cmd/test-ipp
go build -o bin/test-print ./cmd/test-print
go build -o bin/print-info ./cmd/print-info
```

### Configuration

Set the `PRINTER_URI` environment variable:

```bash
# For local CUPS printer (recommended)
export PRINTER_URI="http://localhost:631/printers/EPSON_ET-8550_Series"

# Or for network printer
export PRINTER_URI="http://EPSONXXXXXX.local:631/ipp/print"
```

**Find your printer URI:**
```bash
lpstat -v
# Output: device for EPSON_ET-8550_Series: ipp://EPSONXXXXXX.local:631/ipp/print
```

### Usage

**Check printer status:**
```bash
./bin/test-ipp
```
Shows printer name, model, state, and ink levels for all 6 tanks.

**Test print:**
```bash
./bin/test-print
```
Prints a test PDF (requires `testprint_gopher.pdf` in current directory).

**Generate status report:**
```bash
./bin/print-info
```
Generates a styled PDF report with printer information and prints it.

---

## Available Commands

### `test-ipp` - Printer Status Check

Connects to the printer via IPP and displays:
- Printer name and model
- Current state (Idle, Processing, Stopped)
- State reasons and messages
- Ink levels for all 6 tanks with visual bars:
  - **MB** - Matte Black
  - **PB** - Photo Black
  - **C** - Cyan
  - **Y** - Yellow
  - **M** - Magenta
  - **GY** - Gray

**Example output:**
```
===========================================
Epson ET-8550 IPP Connection Test
===========================================

Connecting to: http://localhost:631/printers/EPSON_ET-8550_Series

--- PRINTER INFORMATION ---
Printer Info: EPSON ET-8550 Series
Model: EPSON ET-8550

--- PRINTER STATUS ---
State: Idle
State Reasons: none

--- INK LEVELS (6-COLOR SYSTEM) ---
Matte Black          [████████████████░░░░]  80% (#000000)
Photo Black          [██████████████████░░]  90% (#000000)
Cyan                 [███████████████░░░░░]  75% (#00FFFF)
Yellow               [██████████░░░░░░░░░░]  50% (#FFFF00)
Magenta              [████░░░░░░░░░░░░░░░░]  20% (#FF00FF)
Gray                 [█████████████████░░░]  85% (#808080)

===========================================
✓ IPP Communication successful!
===========================================
```

### `test-print` - Test PDF Printing

Prints `testprint_gopher.pdf` with default settings:
- Paper: A4
- Tray: Main
- Media: Plain paper (stationery)
- Quality: 3 (draft)

### `print-info` - Status Report Generation

Generates and prints a styled PDF report containing:
- Printer information (name, model, status)
- Graphical ink level bars with color coding:
  - 🟢 Green: >50% (healthy)
  - 🟡 Yellow: 20-50% (medium)
  - 🔴 Red: <20% (low)
- Program information (Go version, libraries used)
- Printer URI
- Timestamp
- Footer with project link

The PDF is saved as `printer-status-YYYYMMDD-HHMMSS.pdf` for reference.

---

## Project Structure

```
.
├── cmd/
│   ├── test-ipp/          # Printer status checker
│   │   └── main.go
│   ├── test-print/        # Test print tool
│   │   └── main.go
│   └── print-info/        # Status report generator
│       └── main.go
├── pkg/
│   └── printer/           # Core printer library
│       ├── info.go        # Printer information and status
│       ├── info_test.go   # Info tests
│       ├── print.go       # PDF printing functionality
│       ├── print_test.go  # Print tests
│       ├── report.go      # PDF report generation
│       ├── report_test.go # Report tests
│       ├── options.go     # Print options and defaults
│       ├── options_test.go# Options tests
│       ├── format.go      # Output formatting
│       └── format_test.go # Format tests
├── bin/                   # Compiled binaries (gitignored)
├── go.mod                 # Go module definition
├── go.sum                 # Dependency checksums
├── SETUP.md              # Detailed setup guide
└── README.md             # This file
```

---

## Library Usage

### Import

```go
import "github.com/Eric-Eklund/epson-printing/pkg/printer"
```

### Get Printer Information

```go
printerURI := "http://localhost:631/printers/EPSON_ET-8550_Series"

// Get all printer information
info, err := printer.GetPrinterInfo(printerURI)
if err != nil {
    log.Fatal(err)
}

// Access printer details
fmt.Println("Printer:", info.Name)
fmt.Println("Model:", info.Model)
fmt.Println("State:", info.State)

// Print formatted output
info.Print()

// Export to JSON
jsonData, _ := info.ToJSON()
fmt.Println(string(jsonData))
```

### Print PDF

```go
// Configure print options
opts := printer.PrintOptions{
    PaperSize: "A4.Borderless",
    Tray:      "Main",
    MediaType: "photographic-glossy",
    Quality:   5, // 3=draft, 4=high, 5=best
    PageRange: "all",
    Copies:    1,
}

// Print PDF
jobID, err := printer.PrintPDF(printerURI, "document.pdf", opts)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Print job sent! Job ID: %d\n", jobID)
```

### Generate Status Report

```go
// Generate PDF report
info, _ := printer.GetPrinterInfo(printerURI)
err := printer.GenerateStatusReport(info, printerURI, "status.pdf")
if err != nil {
    log.Fatal(err)
}

// Or generate and print in one step
pdfPath, jobID, err := printer.PrintStatusReport(printerURI)
```

---

## Testing

### Run All Tests

```bash
go test -v ./...
```

### Run Specific Package Tests

```bash
go test -v ./pkg/printer
```

### Run With Coverage

```bash
go test -cover ./pkg/printer
```

### Integration Tests

Integration tests that require a physical printer are skipped by default. To run them:

```bash
INTEGRATION_TEST=1 go test -v ./pkg/printer
```

**Test Results:**
- ✅ 32 unit tests passing
- ✅ 2 integration tests (skipped by default)
- ✅ Full coverage of report generation
- ✅ Edge case handling verified

---

## Print Options Reference

### Paper Sizes

- `4x6.Borderless` - 4x6" photo (101.6x152.4mm)
- `5x7.Borderless` - 5x7" photo (127x178mm)
- `8x10.Borderless` - 8x10" photo (203x254mm)
- `A4.Borderless` - A4 borderless (210x297mm)
- `A3.Borderless` - A3 borderless (297x420mm)
- `13x19.Borderless` - A3+ borderless (329x483mm)
- `Letter.Borderless` - Letter borderless (8.5x11")

Add `.Borderless` suffix for edge-to-edge printing, omit for margins.

### Paper Trays

- `Auto` - Auto-select based on paper size
- `Photo` - Front photo tray (4x6, 5x7)
- `Main` - Main cassette (A4 documents)
- `Rear` - Rear feed (A3, A3+, thick paper)
- `Manual` - Manual feed
- `Disc` - CD/DVD printing

### Media Types

- `photographic-glossy` - Glossy photo paper
- `photographic-matte` - Matte photo paper
- `PhotographicHighGloss` - Ultra glossy photo paper
- `stationery` - Plain copy paper
- `stationery-inkjet` - Inkjet paper

### Print Quality

- `3` - Draft (fast, test prints)
- `4` - High (normal prints)
- `5` - Best (final prints, slow)

---

## IPP Protocol

This project uses IPP (Internet Printing Protocol) for direct printer communication:

**What is IPP?**
- Standard protocol for network printing (RFC 8010)
- Based on HTTP/HTTPS
- Supported by CUPS and modern printers
- No proprietary drivers needed

**IPP Operations Used:**
- `Get-Printer-Attributes` - Retrieve printer status and capabilities
- `Print-Job` - Submit print jobs with full option control

**Libraries:**
- [OpenPrinting/goipp](https://github.com/OpenPrinting/goipp) - IPP protocol implementation
- [go-pdf/fpdf](https://github.com/go-pdf/fpdf) - PDF generation for reports

---

## Future Development

### Planned Features 🚀

**Enhanced Monitoring:**
- Waste ink maintenance box level monitoring
- Print job queue monitoring
- Historical ink usage tracking and analytics
- Email/notification alerts for low ink
- Web dashboard for real-time status

**Advanced Printing:**
- Batch photo printing tool (Go equivalent of `print-folder.sh`)
- Print preview functionality
- Saved print presets (e.g., "4x6 glossy high quality")
- Print cost calculator (estimate ink usage per job)
- Color profile management and ICC support

**GUI Applications:**
- **Desktop GUI** (Fyne or Gio framework)
  - Drag-and-drop PDF printing
  - Visual ink level display
  - Print queue management
  - Settings presets
- **Web Interface**
  - Remote printing from any device
  - Mobile-responsive design
  - Multi-user support
- **Mobile Companion App**
  - Status monitoring on phone
  - Quick print from mobile photos

**Developer Features:**
- Plugin system for custom workflows
- REST API server for integrations
- Webhook support for printer events
- Lightroom/Darktable integration

### Long-term Vision 💡

**Cross-Platform Distribution:**
- Single-binary releases for Linux, macOS, Windows
- Package manager support (apt, brew, chocolatey, AUR)
- Docker container for headless server printing
- Snap/Flatpak packaging

**Community Features:**
- ICC color profile sharing repository
- Print presets community library
- Plugin marketplace
- Support for other Epson EcoTank models (ET-2800, ET-4800, etc.)

**Enterprise Features:**
- Print job scheduling
- Automatic maintenance reminders
- Ink refill cost tracking
- Print statistics dashboard
- Fleet management for multiple printers
- Cloud printing integration

---

## Why Go?

**Advantages:**
- ✅ Direct IPP communication (minimal dependencies)
- ✅ Cross-platform compatibility (Linux, macOS, Windows)
- ✅ Single static binary (easy distribution)
- ✅ Type-safe API with excellent error handling
- ✅ Built-in testing framework
- ✅ Fast compilation and execution
- ✅ Great concurrency for batch operations
- ✅ Easy to build GUI applications

**Bash vs Go:**

The original bash scripts are now in a separate directory and still work great for quick, personal use. The Go implementation provides:
- Better error handling and validation
- Cross-platform compatibility
- Foundation for GUI applications
- Easier to test and maintain
- Professional distribution to other users
- Type safety and modern tooling

**Recommendation:** Use Go tools for regular use and as a foundation for future development. Bash scripts remain available for quick operations if needed.

---

## Troubleshooting

### "PRINTER_URI environment variable not set"

Set the environment variable:
```bash
export PRINTER_URI="http://localhost:631/printers/EPSON_ET-8550_Series"
```

Or create a `.env` file (see `.env.example`).

### "connection refused"

- Check CUPS is running: `systemctl status cups`
- Verify printer is on and connected
- Verify URI is correct: `lpstat -v`

### "Message truncated" or IPP errors

- Use CUPS URI format: `http://localhost:631/printers/PRINTER_NAME`
- Network URIs should use `http://` not `ipp://` scheme

### Test files not found

Run commands from repository root where test files are located.

---

## Resources

**Documentation:**
- [SETUP.md](./SETUP.md) - Detailed setup guide
- [pkg/printer/README.md](./pkg/printer/README.md) - Library API documentation

**IPP Protocol:**
- [RFC 8010 - IPP/1.1](https://datatracker.ietf.org/doc/html/rfc8010)
- [CUPS IPP Documentation](https://www.cups.org/doc/spec-ipp.html)

**Epson ET-8550:**
- [Product Page](https://epson.com/ecotank-et-8550)
- 6-color ink system (MB, PB, C, Y, M, GY)
- Maximum paper size: A3+ (13x19")
- Network and USB connectivity

---

## Contributing

Contributions and ideas welcome! Areas of interest:
- Additional Epson printer model support
- GUI application development
- Cross-platform testing
- Feature requests and bug reports

---

## License

Personal project - check individual library licenses for dependencies.

---

## Acknowledgments

- [OpenPrinting](https://openprinting.github.io/) for the excellent goipp library
- [go-pdf](https://github.com/go-pdf) for fpdf PDF generation
- CUPS project for IPP protocol implementation

---

**Project Status:** ✅ Production Ready
**Last Updated:** December 30, 2025
**Printer:** Epson EcoTank ET-8550
**Go Version:** 1.18+
**Author:** Eric Eklund
