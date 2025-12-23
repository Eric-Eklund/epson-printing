#!/bin/bash
# Batch print all images in a folder
#
# Usage (positional): ./print-folder.sh [folder] [paper-size] [tray] [media-type] [quality]
# Usage (flags):      ./print-folder.sh [folder] [-p paper] [-t tray] [-m media] [-q quality]
#
# Defaults: 4x6.Borderless, Photo tray, glossy, quality 5
#
# Examples:
#   ./print-folder.sh ~/Photos                           # Use all defaults
#   ./print-folder.sh ~/Photos -q 3                      # Draft quality, other defaults
#   ./print-folder.sh ~/Photos -p A4.Borderless -m ultra # A4 ultra glossy
#   ./print-folder.sh ~/Photos 4x6 Photo glossy 3        # Positional parameters

# Set defaults
FOLDER=""
PAPER="4x6.Borderless"
TRAY="Photo"
MEDIA_INPUT="glossy"
QUALITY="5"

# Store first non-flag argument as folder
if [[ -n "$1" && ! "$1" =~ ^- ]]; then
    FOLDER="$1"
    shift
else
    FOLDER="."
fi

# Parse positional parameters (backward compatibility)
# Only if no flags are present in remaining args
if [[ "$#" -gt 0 && ! "$*" =~ "-" ]]; then
    [[ -n "$1" ]] && PAPER="$1"
    [[ -n "$2" ]] && TRAY="$2"
    [[ -n "$3" ]] && MEDIA_INPUT="$3"
    [[ -n "$4" ]] && QUALITY="$4"
else
    # Parse flags (override defaults or positional)
    while getopts "p:t:m:q:h" opt; do
        case $opt in
            p) PAPER="$OPTARG" ;;
            t) TRAY="$OPTARG" ;;
            m) MEDIA_INPUT="$OPTARG" ;;
            q) QUALITY="$OPTARG" ;;
            h)
                echo "Usage: $0 [folder] [-p paper] [-t tray] [-m media] [-q quality]"
                echo ""
                echo "Flags:"
                echo "  -p  Paper size (default: 4x6.Borderless)"
                echo "  -t  Paper tray: Auto, Photo, Main, Rear (default: Photo)"
                echo "  -m  Media type: ultra, glossy, matte, semi, plain (default: glossy)"
                echo "  -q  Quality: 3 (draft), 4 (high), 5 (best) (default: 5)"
                echo ""
                echo "Examples:"
                echo "  $0 ~/Photos -q 3"
                echo "  $0 ~/Photos -p A4.Borderless -m ultra -q 5"
                exit 0
                ;;
            *)
                echo "Use -h for help"
                exit 1
                ;;
        esac
    done
fi

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
