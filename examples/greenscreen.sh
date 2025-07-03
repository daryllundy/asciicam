#!/bin/bash
# Greenscreen example with background sample generation

echo "Step 1: Generate background samples..."
echo "Please step out of the camera view and press Enter"
read
./asciicam -gen=true -sample=bgdata

echo "Step 2: Using greenscreen effect..."
echo "Step back into the camera view and press Enter"
read
./asciicam -greenscreen=true -sample=bgdata -threshold=0.12 -ansi=true