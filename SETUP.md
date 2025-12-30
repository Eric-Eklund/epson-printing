# Setup Guide - Go Tools

Quick setup guide for the Go-based printer tools.

## Prerequisites

- Go 1.18 or later
- Epson ET-8550 printer configured in CUPS (for local printing)
- Or network-accessible ET-8550

## Setup

### 1. Configure Printer URI

The Go tools require the `PRINTER_URI` environment variable to be set.

**Option A: Using .env file (recommended for permanent setup)**

```bash
# Copy the example file
cp .env.example .env

# Edit .env and set your printer URI
# For local CUPS:
PRINTER_URI=http://localhost:631/printers/EPSON_ET-8550_Series

# Or for network printer:
PRINTER_URI=ipp://your-printer.local:631/ipp/print
```

Then load it:
```bash
source .env
# Or use: export $(cat .env | xargs)
```

**Option B: Export directly (for one-time use)**

```bash
# For local CUPS printer:
export PRINTER_URI="http://localhost:631/printers/EPSON_ET-8550_Series"

# Or for network printer:
export PRINTER_URI="ipp://your-printer.local:631/ipp/print"
```

### 2. Find Your Printer URI

If you don't know your printer URI:

```bash
# List all configured printers
lpstat -v

# Output example:
# device for EPSON_ET-8550_Series: ipp://localhost/printers/EPSON_ET-8550_Series
```

Use the URI from the output, replacing `ipp://localhost` with `http://localhost:631`.

### 3. Build the Tools

```bash
# Build all tools
go build -o cmd/test-ipp/test-ipp ./cmd/test-ipp
go build -o cmd/test-print/test-print ./cmd/test-print
go build -o cmd/print-info/print-info ./cmd/print-info
```

## Usage

### Check Printer Status

```bash
./cmd/test-ipp/test-ipp
```

Shows:
- Printer name and model
- Current state (Idle, Processing, etc.)
- Ink levels for all 6 tanks with visual bars

### Test Print

```bash
./cmd/test-print/test-print
```

Prints `testprint_gopher.pdf` on A4 plain paper from Main tray (draft quality).

### Print Status Report

```bash
./cmd/print-info/print-info
```

Generates and prints a PDF status report containing:
- Printer information (name, model, status)
- Ink levels for all 6 tanks with visual bars
- Program information (Go version, libraries used)
- Timestamp and printer URI
- Saves PDF as `printer-status-YYYYMMDD-HHMMSS.pdf` for reference

## Troubleshooting

### "PRINTER_URI environment variable not set"

You forgot to set the environment variable. See step 1 above.

### "Error: connection refused"

- Check that CUPS is running: `systemctl status cups`
- Verify printer is turned on and connected
- Check printer URI is correct: `lpstat -v`

### "Error: Test PDF not found"

Run the command from the repository root where `testprint_gopher.pdf` is located.

## Shell Configuration (Optional)

To make PRINTER_URI permanent, add to your `~/.bashrc` or `~/.zshrc`:

```bash
export PRINTER_URI="http://localhost:631/printers/EPSON_ET-8550_Series"
```

Then reload: `source ~/.bashrc`
