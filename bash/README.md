# Epson ET-8550 Printing Scripts & Go Development

Documentation for bash printing scripts and future Go IPP implementation.

## Table of Contents

- [Overview](#overview)
- [Printing Scripts](#printing-scripts)
  - [1. Photo Batch Printing (`print-folder.sh`)](#1-photo-batch-printing-print-foldersh)
  - [2. PDF Printing (`print-pdf.sh`)](#2-pdf-printing-print-pdfsh)
- [Available Options Reference](#available-options-reference)
  - [Paper Sizes](#paper-sizes)
  - [Paper Trays (InputSlot)](#paper-trays-inputslot)
  - [Media Types](#media-types)
  - [Print Quality](#print-quality)
- [Typical Workflows](#typical-workflows)
  - [Photo Printing from Lightroom](#photo-printing-from-lightroom)
  - [Calendar/Poster Printing](#calendarposter-printing)
  - [Document Printing](#document-printing)
- [Quick Reference](#quick-reference)
  - [Scripts in This Repository](#scripts-in-this-repository)
  - [Most Common Commands](#most-common-commands)
  - [Check Printer Status](#check-printer-status)
- [Future Development](#future-development)

---

## Overview

Two bash scripts provide complete control over Epson ET-8550 printing via CUPS:
- **`print-folder.sh`** - Batch print photos from folders
- **`print-pdf.sh`** - Print PDFs with full control

**Future development:** See [GO-DEVELOPMENT.md](../GO-DEVELOPMENT.md) for plans to implement a Go version using IPP protocol.

---

## Printing Scripts

### 1. Photo Batch Printing (`print-folder.sh`)

**Purpose:** Print all photos in a folder (Lightroom export workflow)

**Syntax:**
```bash
# Positional parameters (backward compatible)
./print-folder.sh [folder] [paper-size] [tray] [media-type] [quality]

# Flag-based parameters (flexible - specify only what you need)
./print-folder.sh [folder] [-p paper] [-t tray] [-m media] [-q quality]
```

**Parameters:**
- `folder` - Path to folder containing images (default: current directory)
- `-p` paper-size - Paper size with optional .Borderless suffix (default: 4x6.Borderless)
- `-t` tray - Paper source: Auto, Photo, Main, Rear (default: Photo)
- `-m` media-type - Media type using friendly names or full CUPS names (default: glossy)
- `-q` quality - Print quality: 3 (draft), 4 (high), 5 (best) (default: 5)
- `-h` - Show help message

**Media Type Aliases (case-insensitive):**
- `ultra` / `ultra-glossy` / `high-gloss` → Ultra Glossy Photo Paper
- `glossy` → Standard Glossy (default)
- `semi` / `semi-glossy` / `luster` → Semi-gloss/Luster
- `matte` → Matte Photo Paper
- `photo` / `photographic` → Generic Photographic
- `plain` / `paper` → Plain Paper
- `coated` / `inkjet` → Coated Inkjet Paper
- `velvet` / `fine-art` → Fine Art Paper
- Or use full CUPS names: `PhotographicHighGloss`, `PhotographicGlossy`, etc.

**File Types:** JPG, JPEG, PNG, TIF

**Examples:**
```bash
# Use all defaults (4x6, Photo tray, glossy, quality 5)
./print-folder.sh ~/Photos

# Override only quality - draft prints with other defaults
./print-folder.sh ~/Photos -q 3

# Override specific settings with flags
./print-folder.sh ~/Photos -p A4.Borderless -m ultra -q 5

# A3 borderless matte, high quality (flags)
./print-folder.sh ~/Photos -p A3.Borderless -t Rear -m matte -q 4

# Positional parameters still work (backward compatible)
./print-folder.sh ~/Photos 4x6.Borderless Photo ultra 5

# Mix: change only media type to ultra, keep other defaults
./print-folder.sh ~/Photos -m ultra
```

**Output:**
```
Printing all images from: /path/to/photos
Paper size: 4x6.Borderless
Paper tray: Photo
Media type: ultra → PhotographicHighGloss
Quality: 5 (3=draft, 4=high, 5=best)

Printing: IMG_001.jpg
Printing: IMG_002.jpg
...
Sent 24 images to printer!
```

---

### 2. PDF Printing (`print-pdf.sh`)

**Purpose:** Print PDFs with complete control over all settings

**Syntax:**
```bash
./print-pdf.sh <pdf-file> [paper-size] [tray] [media-type] [quality] [pages]
```

**Parameters:**
- `pdf-file` - Path to PDF file (required)
- `paper-size` - Paper size (default: A4.Borderless)
- `tray` - Paper source (default: Auto)
- `media-type` - Media type (default: photographic-glossy)
- `quality` - Print quality 3-5 (default: 4)
- `pages` - Page range (default: all)

**Page Range Syntax (Go-style slices):**
| Syntax | Prints | Description |
|--------|--------|-------------|
| `1` | Page 1 | Single page |
| `1-5` | Pages 1-5 | Range notation |
| `1,3,5` | Pages 1,3,5 | Specific pages |
| `:5` | Pages 1-5 | **Go slice: first to 5** |
| `5:` | Pages 5-end | **Go slice: 5 to end** |
| `all` | All pages | All (default, can omit) |

**Examples:**
```bash
# Basic usage (defaults: A4 borderless, auto, glossy, quality 4, all pages)
./print-pdf.sh document.pdf

# A3+ matte, best quality, page 1 only (test print)
./print-pdf.sh calendar.pdf 13x19.Borderless Rear photographic-matte 5 1

# Same but full quality (final print)
./print-pdf.sh calendar.pdf 13x19.Borderless Rear photographic-matte 5

# Print pages 2-5
./print-pdf.sh document.pdf A4.Borderless Main stationery 4 2-5

# Go-style: Print first 5 pages
./print-pdf.sh document.pdf A4.Borderless Auto photographic-glossy 4 :5

# Go-style: Print from page 5 to end
./print-pdf.sh document.pdf A4.Borderless Auto photographic-glossy 4 5:

# Draft quality on plain paper
./print-pdf.sh document.pdf A4 Main stationery 3

# Pages 1 and 3 only
./print-pdf.sh document.pdf A4 Main stationery 4 1,3
```

**Output:**
```
=========================================
PDF PRINT SETTINGS
=========================================
File:        calendar.pdf
Paper size:  13x19.Borderless
Tray:        Rear
Media type:  photographic-matte
Quality:     5 (3=draft, 4=high, 5=best)
Pages:       :5 (converted to 1-5)
=========================================

request id is EPSON_ET-8550_Series-42 (1 file(s))
✓ Print job sent successfully!
```

---

## Available Options Reference

### Paper Sizes

**Common Sizes:**
- `4x6.Borderless` - 101.6x152.4mm (standard photo)
- `5x7.Borderless` - 127x178mm
- `8x10.Borderless` - 203x254mm
- `A4.Borderless` - 210x297mm
- `A3.Borderless` - 297x420mm
- `13x19.Borderless` - 329x483mm **(A3+)**
- `Letter.Borderless` - 8.5x11"

**Note:** Add `.Borderless` for edge-to-edge printing, omit for margins.

### Paper Trays (InputSlot)

- `Auto` - Automatically select based on paper size
- `Photo` - Front photo tray (4x6, 5x7)
- `Main` - Main cassette (A4 documents)
- `Rear` - Rear feed (A3, A3+, thick paper)
- `Manual` - Manual feed
- `Disc` - CD/DVD printing

**Recommendation:**
- Small photos (4x6, 5x7) → **Photo**
- Large/thick (A3, A3+) → **Rear**
- Documents (A4) → **Main**

### Media Types

- `photographic-glossy` - Glossy photo paper
- `photographic-matte` - Matte photo paper
- `stationery` - Plain copy paper
- `stationery-inkjet` - Inkjet paper

### Print Quality

- `3` - Draft (fast, test prints)
- `4` - High (normal prints)
- `5` - Best (final prints, slow)

---

## Typical Workflows

### Photo Printing from Lightroom
```bash
# 1. Export from Lightroom to folder
# 2. Test print one photo with draft quality
./print-folder.sh /path/to/exported-photos 4x6.Borderless Photo ultra 3

# 3. Batch print all photos with best quality
./print-folder.sh /path/to/exported-photos 4x6.Borderless Photo ultra 5
```

### Calendar/Poster Printing
```bash
# 1. Create in OnlyOffice/LibreOffice, export as PDF
# 2. Test print one page (draft quality)
./print-pdf.sh calendar.pdf 13x19.Borderless Rear photographic-matte 3 1

# 3. Final print all pages (best quality)
./print-pdf.sh calendar.pdf 13x19.Borderless Rear photographic-matte 5
```

### Document Printing
```bash
# Draft review
./print-pdf.sh document.pdf A4 Main stationery 3

# Final print, first 10 pages
./print-pdf.sh document.pdf A4 Main stationery 4 :10
```

---

## Quick Reference

### Scripts in This Repository
```
print-folder.sh          - Photo batch printing
print-pdf.sh            - PDF printing with full control
```

### Most Common Commands
```bash
# Batch print 4x6 photos (uses defaults: 4x6, Photo, glossy, quality 5)
./print-folder.sh ~/Photos

# Quick draft test - override only quality
./print-folder.sh ~/Photos -q 3

# Ultra glossy, best quality
./print-folder.sh ~/Photos -m ultra

# Test print A3+ matte (page 1, draft)
./print-pdf.sh calendar.pdf 13x19.Borderless Rear photographic-matte 3 1

# Final print A3+ matte (best quality)
./print-pdf.sh calendar.pdf 13x19.Borderless Rear photographic-matte 5

# Print document (A4 plain, draft)
./print-pdf.sh document.pdf A4 Main stationery 3

# Print first 5 pages (Go-style slice)
./print-pdf.sh document.pdf A4 Auto stationery 4 :5
```

### Check Printer Status

**Command Line:**
```bash
lpstat -p -d                                    # Printer status
lpoptions -p EPSON_ET-8550_Series -l           # All available options

# Check ink levels via IPP
ipptool -tv ipp://EPSONXXXXXX.local:631/ipp/print get-printer-attributes.test 2>&1 | grep -i "marker-levels"
```

**Web Interface:**
```
http://EPSONXXXXXX.local.:80/PRESENTATION/HTML/TOP/PRTINFO.HTML
```

**What you can monitor:**
- **Ink Tank Levels** - All 6 color tanks (MB, PB, C, Y, M, GY)
- **Waste Ink Maintenance Box** - Water droplet icon (far right)

**About the Waste Ink Maintenance Box:**

The maintenance box (water droplet symbol) collects waste ink from:
- Print head cleaning cycles
- Borderless printing overspray (ink sprayed beyond paper edges)
- Initial ink charging and maintenance operations

**Important Notes:**
- Borderless printing uses more waste ink than bordered prints
- Replacement part: **Epson T04D1** Maintenance Box
- Printer alerts you at ~90% full
- EcoTank boxes last thousands of prints (much larger than traditional Epson printers)
- Monitor occasionally if doing frequent borderless printing

**When to replace:** The printer will display a warning when the box is nearly full. Continue printing until prompted - no need to replace early.

---

## Future Development

Interested in a Go implementation with GUI, cross-platform support, and advanced features? See [GO-DEVELOPMENT.md](../GO-DEVELOPMENT.md) for the complete development plan.

---

**Last Updated:** December 2025
**Printer:** Epson EcoTank ET-8550
**System:** Linux (tested on Arch-based systems)
