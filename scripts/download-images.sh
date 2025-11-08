#!/bin/bash

# Default configuration
NUM_UNIQUE=70      # Number of unique images to download
COPIES_PER_IMAGE=6 # Number of duplicates per image

# Function to show usage
show_help() {
  cat << EOF
Usage: $(basename "$0") [OPTIONS]

Download sample duplicate images for testing Schluckauf.

Options:
  --unique NUM    Number of unique images to download (default: 70)
  --copies NUM    Number of copies per image (default: 6)
  -h, --help      Show this help message

Examples:
  $(basename "$0")
  $(basename "$0") --unique=50 --copies=3
  $(basename "$0") --unique 100 --copies 10

Prerequisites:
  - ImageMagick must be installed for image validation
  - curl for downloading images

EOF
}

# Parse command-line arguments
while [[ $# -gt 0 ]]; do
  case $1 in
    --unique=*)
      NUM_UNIQUE="${1#*=}"
      shift
      ;;
    --unique)
      NUM_UNIQUE="$2"
      shift 2
      ;;
    --copies=*)
      COPIES_PER_IMAGE="${1#*=}"
      shift
      ;;
    --copies)
      COPIES_PER_IMAGE="$2"
      shift 2
      ;;
    -h|--help)
      show_help
      exit 0
      ;;
    *)
      echo "Error: Unknown option: $1"
      echo "Use --help for usage information"
      exit 1
      ;;
  esac
done

# Validate arguments
if ! [[ "$NUM_UNIQUE" =~ ^[0-9]+$ ]] || [ "$NUM_UNIQUE" -lt 1 ]; then
  echo "Error: --unique must be a positive integer"
  exit 1
fi

if ! [[ "$COPIES_PER_IMAGE" =~ ^[0-9]+$ ]]; then
  echo "Error: --copies must be a non-negative integer"
  exit 1
fi

if [ "$COPIES_PER_IMAGE" -eq 0 ]; then
  echo "Downloading $NUM_UNIQUE unique images (no copies)..."
else
  echo "Downloading $NUM_UNIQUE unique images with $COPIES_PER_IMAGE copies each..."
fi

successful=0
for i in $(seq 1 $NUM_UNIQUE); do
  # Use fixed image ID to ensure we can get the same image
  image_id=$((100 + i))
  width=$((800 + RANDOM % 400))
  height=$((600 + RANDOM % 400))

  original="image${i}_original.jpg"

  echo "[$i/$NUM_UNIQUE] Downloading unique image: $original (${width}x${height})..."

  # Try up to 3 times
  retry=0
  max_retries=3
  download_success=false

  while [ $retry -lt $max_retries ]; do
    curl -L -s -o "$original" "https://picsum.photos/id/${image_id}/${width}/${height}.jpg"

    # Verify file exists, not empty, and is valid image using identify
    if [ -s "$original" ] && identify -format '%m' "$original" &>/dev/null; then
      download_success=true
      break
    fi

    # Invalid or corrupted file, clean up
    rm -f "$original"

    retry=$((retry + 1))
    if [ $retry -lt $max_retries ]; then
      echo "  Retry $retry/$max_retries..."
      sleep 1
    fi
  done

  if [ "$download_success" = false ]; then
    echo "  ✗ Failed after $max_retries attempts, skipping..."
    rm -f "$original"
    continue
  fi

  echo "  ✓ Downloaded successfully"

  # Create duplicates of this image (only if copies > 0)
  if [ "$COPIES_PER_IMAGE" -gt 0 ]; then
    for copy in $(seq 1 $COPIES_PER_IMAGE); do
      duplicate="image${i}_copy${copy}.jpg"
      cp "$original" "$duplicate"
      echo "  Created duplicate: $duplicate"
    done
  fi

  successful=$((successful + 1))
  sleep 0.5
done

echo ""
if [ "$COPIES_PER_IMAGE" -eq 0 ]; then
  echo "Done! Successfully created $successful unique images (no copies)."
  echo "Total files: $successful"
else
  echo "Done! Successfully created $successful unique images with $COPIES_PER_IMAGE copies each."
  echo "Total files: $((successful * (COPIES_PER_IMAGE + 1)))"
fi
echo "Failed: $((NUM_UNIQUE - successful))"
