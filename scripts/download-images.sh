#!/bin/bash

# Configuration
NUM_UNIQUE=70      # Number of unique images to download
COPIES_PER_IMAGE=6 # Number of duplicates per image

echo "Downloading $NUM_UNIQUE unique images with $COPIES_PER_IMAGE copies each..."

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

  # Create duplicates of this image
  for copy in $(seq 1 $COPIES_PER_IMAGE); do
    duplicate="image${i}_copy${copy}.jpg"
    cp "$original" "$duplicate"
    echo "  Created duplicate: $duplicate"
  done

  successful=$((successful + 1))
  sleep 0.5
done

echo ""
echo "Done! Successfully created $successful unique images with $COPIES_PER_IMAGE copies each."
echo "Total files: $((successful * (COPIES_PER_IMAGE + 1)))"
echo "Failed: $((NUM_UNIQUE - successful))"
