#!/bin/bash
# Batch print all images in a folder
# Usage: ./print-folder.sh [folder] [paper-size] [tray]
# Example: ./print-folder.sh ~/Pictures/Exports 4x6.Borderless Photo

FOLDER="${1:-.}"
PAPER="${2:-A4.Borderless}"
TRAY="${3:-Auto}"

echo "Printing all images from: $FOLDER"
echo "Paper size: $PAPER"
echo "Paper tray: $TRAY"
echo ""

count=0
for img in "$FOLDER"/*.{jpg,JPG,jpeg,JPEG,png,PNG,tif,TIF}; do
    [ -f "$img" ] || continue
    echo "Printing: $(basename "$img")"
    lp -d EPSON_ET-8550_Series \
       -o PageSize="$PAPER" \
       -o InputSlot="$TRAY" \
       -o media=photographic-glossy \
       -o print-quality=5 \
       "$img"
    ((count++))
done

echo ""
echo "Sent $count images to printer!"
