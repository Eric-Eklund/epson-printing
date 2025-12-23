# Go Development Plan

Future development plan for implementing Epson ET-8550 printing in Go using the IPP protocol.

---

## Why Go for Printing?

**Current Limitations:**
- GUI apps (OnlyOffice, darktable, even Epson's own apps) don't expose all printer settings
- Settings get ignored or simplified
- True on Linux, Android, iOS, Windows

**Go Advantages:**
- Direct IPP communication (no CUPS dependency)
- Single compiled binary
- Cross-platform
- Type-safe API
- Can build GUI (Fyne, GTK, web interface)
- Better error handling
- Same functionality as bash scripts, more extensible

---

## IPP (Internet Printing Protocol)

**What is IPP?**
- Standard protocol for network printing
- Based on HTTP/HTTPS
- Supported by CUPS and modern printers
- ET-8550 fully supports IPP

**ET-8550 IPP Connection:**
```
ipp://localhost/printers/EPSON_ET-8550_Series
```
Network discovery:
```
dnssd://EPSON ET-8550 Series._ipp._tcp.local/
```

**IPP Features:**
- Submit print jobs
- Query printer status
- Get ink levels
- Cancel jobs
- Set all print options (paper size, tray, quality, etc.)
- Same features as CUPS command-line

---

## Go IPP Library

**Primary Library:** `github.com/phin1x/go-ipp`

**Installation:**
```bash
go get github.com/phin1x/go-ipp
```

**Basic Example:**
```go
package main

import (
    "fmt"
    "os"

    "github.com/phin1x/go-ipp"
)

func main() {
    // Connect to printer
    client := ipp.NewClient("http://localhost:631/printers/EPSON_ET-8550_Series")

    // Get printer attributes
    req := ipp.NewRequest(ipp.OperationGetPrinterAttributes, 1)
    req.OperationAttributes[ipp.AttributePrinterURI] = "http://localhost:631/printers/EPSON_ET-8550_Series"

    resp, err := client.SendRequest(req)
    if err != nil {
        panic(err)
    }

    fmt.Println("Printer Status:", resp.PrinterAttributes["printer-state"])
}
```

**Print Job Example:**
```go
func printPDF(filename, paperSize, tray, media string, quality int) error {
    client := ipp.NewClient("http://localhost:631/printers/EPSON_ET-8550_Series")

    // Read PDF file
    data, err := os.ReadFile(filename)
    if err != nil {
        return err
    }

    // Create print job
    req := ipp.NewRequest(ipp.OperationPrintJob, 1)
    req.OperationAttributes[ipp.AttributePrinterURI] = "http://localhost:631/printers/EPSON_ET-8550_Series"
    req.OperationAttributes[ipp.AttributeRequestingUserName] = "username"
    req.OperationAttributes[ipp.AttributeJobName] = filename
    req.OperationAttributes[ipp.AttributeDocumentFormat] = "application/pdf"

    // Job attributes (print settings)
    req.JobAttributes[ipp.AttributeMediaSize] = paperSize              // "13x19.Borderless"
    req.JobAttributes[ipp.AttributeMediaSource] = tray                 // "Rear"
    req.JobAttributes[ipp.AttributeMediaType] = media                  // "photographic-matte"
    req.JobAttributes[ipp.AttributePrintQuality] = quality             // 3, 4, or 5

    req.File = data

    resp, err := client.SendRequest(req)
    if err != nil {
        return err
    }

    fmt.Printf("Job ID: %d\\n", resp.JobAttributes[ipp.AttributeJobID])
    return nil
}
```

---

## Project Structure

**Proposed Go Project:**
```
epson-printing-go/
├── cmd/
│   ├── print-pdf/          # CLI tool for PDF printing
│   │   └── main.go
│   ├── print-photos/       # CLI tool for photo batch printing
│   │   └── main.go
│   └── print-gui/          # GUI application (Fyne)
│       └── main.go
├── pkg/
│   ├── printer/            # Printer abstraction
│   │   ├── printer.go      # Main printer interface
│   │   ├── ipp.go          # IPP implementation
│   │   └── options.go      # Print options (paper size, tray, etc.)
│   └── jobs/               # Job management
│       ├── pdf.go          # PDF printing
│       ├── photos.go       # Photo batch printing
│       └── queue.go        # Job queue
├── go.mod
├── go.sum
└── README.md
```

---

## Getting Started with Go Version

**1. Initialize Go Project:**
```bash
# In the repository directory
mkdir go-version && cd go-version
go mod init github.com/yourusername/epson-printing
go get github.com/phin1x/go-ipp
```

**2. Create Basic Printer Package:**
```go
// pkg/printer/printer.go
package printer

type Printer struct {
    Name string
    URI  string
}

type PrintOptions struct {
    PaperSize  string
    Tray       string
    MediaType  string
    Quality    int
    PageRanges string
}

func NewPrinter(name, uri string) *Printer {
    return &Printer{Name: name, URI: uri}
}

func (p *Printer) PrintPDF(filename string, opts PrintOptions) error {
    // Implementation using go-ipp
    return nil
}

func (p *Printer) GetStatus() (string, error) {
    // Implementation using go-ipp
    return "", nil
}
```

**3. Create CLI Tool:**
```go
// cmd/print-pdf/main.go
package main

import (
    "flag"
    "fmt"
    "os"

    "github.com/yourusername/epson-printing/pkg/printer"
)

func main() {
    paperSize := flag.String("paper", "A4.Borderless", "Paper size")
    tray := flag.String("tray", "Auto", "Paper tray")
    media := flag.String("media", "photographic-glossy", "Media type")
    quality := flag.Int("quality", 4, "Print quality (3-5)")
    pages := flag.String("pages", "all", "Page ranges")

    flag.Parse()

    if flag.NArg() < 1 {
        fmt.Println("Usage: print-pdf [options] <file.pdf>")
        os.Exit(1)
    }

    filename := flag.Arg(0)

    p := printer.NewPrinter(
        "EPSON_ET-8550_Series",
        "http://localhost:631/printers/EPSON_ET-8550_Series",
    )

    opts := printer.PrintOptions{
        PaperSize:  *paperSize,
        Tray:       *tray,
        MediaType:  *media,
        Quality:    *quality,
        PageRanges: *pages,
    }

    err := p.PrintPDF(filename, opts)
    if err != nil {
        fmt.Printf("Error: %v\\n", err)
        os.Exit(1)
    }

    fmt.Println("✓ Print job sent successfully!")
}
```

**4. Build and Run:**
```bash
# Build
go build -o print-pdf ./cmd/print-pdf

# Run
./print-pdf --paper 13x19.Borderless --tray Rear --media photographic-matte --quality 5 calendar.pdf
```

---

## Additional Resources

**IPP Protocol Documentation:**
- RFC 8010: https://datatracker.ietf.org/doc/html/rfc8010
- CUPS Documentation: https://www.cups.org/doc/spec-ipp.html

**Go IPP Libraries:**
- `github.com/phin1x/go-ipp` - Primary, well-maintained
- `github.com/OpenPrinting/goipp` - Alternative

**GUI Frameworks:**
- Fyne: https://fyne.io/ (native-looking, cross-platform)
- Gio: https://gioui.org/ (immediate mode)
- GTK: via `gotk3` (Linux-focused)

**Build Tools:**
- Task: https://taskfile.dev/ (Makefile alternative)
- GoReleaser: https://goreleaser.com/ (release automation)

---

## Comparison: Bash vs Go

### Current Bash Scripts

**Pros:**
- ✅ Simple, readable
- ✅ Direct CUPS integration
- ✅ No compilation needed
- ✅ Works perfectly for current use case
- ✅ Easy to modify

**Cons:**
- ❌ CUPS dependency (though not really a con on Linux)
- ❌ Limited error handling
- ❌ No type safety
- ❌ Hard to build GUI
- ❌ Shell-specific syntax

### Future Go Implementation

**Pros:**
- ✅ Direct IPP communication (no CUPS required)
- ✅ Cross-platform (Windows, Mac, Linux)
- ✅ Single binary (easy distribution)
- ✅ Type-safe API
- ✅ Better error handling
- ✅ Can build GUI easily
- ✅ Extensible (web interface, mobile app, etc.)
- ✅ Better for large projects

**Cons:**
- ❌ More code to write
- ❌ Requires compilation
- ❌ Overkill for simple printing scripts

**Recommendation:** Keep bash scripts for personal use. Build Go version if you want to:
- Learn Go
- Build a GUI application
- Distribute to other users
- Add advanced features (print preview, queue management, etc.)

---

## Development Roadmap

**Phase 1: Basic CLI Tools**
1. Learn IPP protocol basics
2. Experiment with `go-ipp` library
3. Build simple CLI tool (print-pdf equivalent)
4. Add features (status checking, ink levels)

**Phase 2: GUI Application**
5. Build GUI with Fyne
6. Package as single binary
7. Add drag-and-drop support

**Phase 3: Advanced Features**
- Print preview
- Queue management
- Saved presets
- Web interface
- Mobile companion app
- Print cost calculator (ink usage)

---

**Last Updated:** December 2025
**Printer:** Epson EcoTank ET-8550
**System:** Linux (tested on Arch-based systems)
