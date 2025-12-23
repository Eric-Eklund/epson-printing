#!/bin/bash
# Batch print all images in a folder
# Usage: ./print-folder.sh [folder] [paper-size] [tray] [media-type] [quality]
# Defaults: 4x6.Borderless, Photo tray, glossy, quality 5
# Example: ./print-folder.sh ~/Pictures/Exports

FOLDER="${1:-.}"
PAPER="${2:-4x6.Borderless}"
TRAY="${3:-Photo}"
MEDIA_INPUT="${4:-glossy}"
QUALITY="${5:-5}"

# Map friendly names to CUPS MediaType values
case "${MEDIA_INPUT,,}" in  # Convert to lowercase for case-insensitive matching
    ultra|ultra-glossy|high-gloss)
        MEDIA="PhotographicHighGloss"
        ;;
    glossy)
        MEDIA="PhotographicGlossy"
        ;;
    semi-glossy|semi|luster)
        MEDIA="PhotographicSemiGloss"
        ;;
    matte)
        MEDIA="PhotographicMatte"
        ;;
    photo|photographic)
        MEDIA="Photographic"
        ;;
    plain|paper)
        MEDIA="Stationery"
        ;;
    coated|inkjet)
        MEDIA="StationeryCoated"
        ;;
    velvet|fine-art)
        MEDIA="Com.epsonVelvetFineArt"
        ;;
    *)
        # If not a friendly name, assume it's a full CUPS name
        MEDIA="$MEDIA_INPUT"
        ;;
esac

echo "Printing all images from: $FOLDER"
echo "Paper size: $PAPER"
echo "Paper tray: $TRAY"
echo "Media type: $MEDIA_INPUT → $MEDIA"
echo "Quality: $QUALITY (3=draft, 4=high, 5=best)"
echo ""

count=0
for img in "$FOLDER"/*.{jpg,JPG,jpeg,JPEG,png,PNG,tif,TIF}; do
    [ -f "$img" ] || continue
    echo "Printing: $(basename "$img")"
    lp -d EPSON_ET-8550_Series \
       -o PageSize="$PAPER" \
       -o InputSlot="$TRAY" \
       -o MediaType="$MEDIA" \
       -o print-quality="$QUALITY" \
       "$img"
    ((count++))
done

echo ""
echo "Sent $count images to printer!"
