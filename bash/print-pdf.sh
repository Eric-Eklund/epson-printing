#!/bin/bash
# Print PDF with full control over settings
# Usage: ./print-pdf.sh <pdf-file> [paper-size] [tray] [media-type] [quality] [pages]
# Examples:
#   ./print-pdf.sh document.pdf 13x19.Borderless Rear photographic-matte 5 1
#   ./print-pdf.sh document.pdf A4 Main stationery 3
#   ./print-pdf.sh document.pdf 4x6.Borderless Photo photographic-glossy 5 1-3
#   ./print-pdf.sh document.pdf A4.Borderless Auto photographic-glossy 4 :5    # First to page 5
#   ./print-pdf.sh document.pdf A4.Borderless Auto photographic-glossy 4 5:    # Page 5 to end

PDF_FILE="$1"
PAPER="${2:-A4.Borderless}"
TRAY="${3:-Auto}"
MEDIA="${4:-photographic-glossy}"
QUALITY="${5:-4}"
PAGES="${6:-all}"

# Check if PDF file exists
if [ ! -f "$PDF_FILE" ]; then
    echo "Error: PDF file not found: $PDF_FILE"
    echo ""
    echo "Usage: $0 <pdf-file> [paper-size] [tray] [media-type] [quality] [pages]"
    echo ""
    echo "Paper sizes: A4.Borderless, A3.Borderless, 13x19.Borderless, 4x6.Borderless, Letter.Borderless"
    echo "Trays: Auto, Photo, Main, Rear"
    echo "Media types: photographic-glossy, photographic-matte, stationery, stationery-inkjet"
    echo "Quality: 3 (draft), 4 (high), 5 (best)"
    echo "Pages: 1, 1-3, 1,3,5, :5 (first to 5), 5: (5 to end), or 'all' (default)"
    exit 1
fi

# Handle Go-style slice syntax for pages
ORIGINAL_PAGES="$PAGES"
if [ "$PAGES" != "all" ]; then
    # :5 means from first page to page 5 (1-5)
    if [[ "$PAGES" =~ ^:([0-9]+)$ ]]; then
        PAGES="1-${BASH_REMATCH[1]}"
    # 5: means from page 5 to end (5-)
    elif [[ "$PAGES" =~ ^([0-9]+):$ ]]; then
        PAGES="${BASH_REMATCH[1]}-"
    fi
fi

echo "========================================="
echo "PDF PRINT SETTINGS"
echo "========================================="
echo "File:        $PDF_FILE"
echo "Paper size:  $PAPER"
echo "Tray:        $TRAY"
echo "Media type:  $MEDIA"
echo "Quality:     $QUALITY (3=draft, 4=high, 5=best)"
if [ "$PAGES" != "$ORIGINAL_PAGES" ]; then
    echo "Pages:       $ORIGINAL_PAGES (converted to $PAGES)"
else
    echo "Pages:       $PAGES"
fi
echo "========================================="
echo ""

# Build the print command
CMD="lp -d EPSON_ET-8550_Series"
CMD="$CMD -o PageSize=$PAPER"
CMD="$CMD -o InputSlot=$TRAY"
CMD="$CMD -o media=$MEDIA"
CMD="$CMD -o print-quality=$QUALITY"
CMD="$CMD -o fit-to-page"

# Add page range if not 'all'
if [ "$PAGES" != "all" ]; then
    CMD="$CMD -o page-ranges=$PAGES"
fi

CMD="$CMD \"$PDF_FILE\""

# Execute
eval $CMD

if [ $? -eq 0 ]; then
    echo "✓ Print job sent successfully!"
else
    echo "✗ Print job failed!"
    exit 1
fi
