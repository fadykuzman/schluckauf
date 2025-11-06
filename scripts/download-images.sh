#!/bin/bash

# Configuration
NUM_UNIQUE=5         # Number of unique images to download
COPIES_PER_IMAGE=3   # Number of duplicates per image

echo "Downloading $NUM_UNIQUE unique images with $COPIES_PER_IMAGE copies each..."

for i in $(seq 1 $NUM_UNIQUE); do
  # Use fixed image ID to ensure we can get the same image
  image_id=$((100 + i))
  width=$((800 + RANDOM % 400))
  height=$((600 + RANDOM % 400))

  original="image${i}_original.jpg"

  echo "[$i/$NUM_UNIQUE] Downloading unique image: $original (${width}x${height})..."
  curl -L -o "$original" "https://picsum.photos/id/${image_id}/${width}/${height}.jpg"

  if [ $? -ne 0 ]; then
    echo "  Error downloading $original, skipping..."
    continue
  fi

  # Create duplicates of this image
  for copy in $(seq 1 $COPIES_PER_IMAGE); do
    duplicate="image${i}_copy${copy}.jpg"
    cp "$original" "$duplicate"
    echo "  Created duplicate: $duplicate"
  done

  sleep 0.5
done

echo ""
echo "Done! Created $NUM_UNIQUE unique images with $COPIES_PER_IMAGE copies each."
echo "Total files: $((NUM_UNIQUE * (COPIES_PER_IMAGE + 1)))"
