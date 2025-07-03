// Package main provides functionality for converting webcam images to ASCII/ANSI art.
// This file contains image processing functions for the ASCII webcam application.
package main

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io/ioutil"
	"math"
	"strings"

	"github.com/lucasb-eyer/go-colorful"
	"github.com/muesli/termenv"
	"github.com/nfnt/resize"
	"gocv.io/x/gocv"
)

// matToImage converts a gocv Mat to Go's native image.RGBA format.
// This is necessary because gocv uses OpenCV's BGR color format, while
// Go's standard image library uses RGBA format.
//
// Parameters:
//   - mat: The OpenCV Mat containing the image data
//   - width: The desired width of the output image
//   - height: The desired height of the output image
//
// Returns:
//   - A pointer to an image.RGBA containing the converted image
func matToImage(mat gocv.Mat, width, height uint) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, int(width), int(height)))

	// Create a copy of the Mat to avoid modifying the original
	bgrMat := mat.Clone()
	defer bgrMat.Close() // Ensure the Mat is properly closed to avoid memory leaks

	// Copy the data from mat to img, converting BGR to RGBA
	for y := 0; y < int(height); y++ {
		for x := 0; x < int(width); x++ {
			pixel := bgrMat.GetVecbAt(y, x)
			img.SetRGBA(x, y, color.RGBA{
				B: pixel[0], // OpenCV stores colors as BGR
				G: pixel[1],
				R: pixel[2],
				A: 255, // Set full opacity
			})
		}
	}

	return img
}

// pixelToASCII converts a color pixel to an ASCII character based on its intensity.
// Darker pixels are represented by characters with less "ink" (like spaces or dots),
// while brighter pixels use more "ink-heavy" characters (like @ or 8).
//
// Parameters:
//   - pixel: The color.Color to convert
//
// Returns:
//   - A rune (character) representing the pixel's brightness
func pixelToASCII(pixel color.Color) rune {
	// Get RGBA values (each in range 0-65535)
	r2, g2, b2, a2 := pixel.RGBA()

	// Convert to 0-255 range
	r := uint(r2 / 256)
	g := uint(g2 / 256)
	b := uint(b2 / 256)
	a := uint(a2 / 256)

	// Calculate intensity, taking alpha into account
	intensity := (r + g + b) * a / 255

	// Calculate precision based on number of available ASCII characters
	precision := float64(255 * 3 / (len(pixels) - 1))

	// Map intensity to an index in the pixels array
	v := int(math.Floor(float64(intensity)/precision + 0.5))
	return pixels[v]
}

// imageToASCII converts an image to ASCII art with color.
// Each pixel is represented by an ASCII character with the appropriate color.
//
// Parameters:
//   - width: The width of the output in characters
//   - height: The height of the output in characters
//   - p: The terminal color profile to use for rendering
//   - img: The source image to convert
//
// Returns:
//   - A string containing the ASCII representation of the image
func imageToASCII(width, height uint, p termenv.Profile, img image.Image) string {
	str := strings.Builder{}

	for i := 0; i < int(height); i++ {
		for j := 0; j < int(width); j++ {
			// Get pixel and convert to ASCII character
			pixel := color.NRGBAModel.Convert(img.At(j, i))
			s := termenv.String(string(pixelToASCII(pixel)))

			// Apply color - either the global color (if set) or the pixel's color
			_, _, _, a := col.RGBA()
			if a > 0 {
				// Use global color if it has been set
				s = s.Foreground(p.FromColor(col))
			} else {
				// Otherwise use the pixel's color
				s = s.Foreground(p.FromColor(pixel))
			}
			str.WriteString(s.String())
		}
		str.WriteString("\n") // End of row
	}

	return str.String()
}

// imageToANSI converts an image to colored ANSI blocks.
// It uses the upper half block character (▀) with foreground and background
// colors to represent two pixels vertically in a single character position.
// This provides higher vertical resolution than ASCII art.
//
// Parameters:
//   - _: Unused width parameter (kept for interface consistency with imageToASCII)
//   - _: Unused height parameter (kept for interface consistency with imageToASCII)
//   - p: The terminal color profile to use for rendering
//   - img: The source image to convert
//
// Returns:
//   - A string containing the ANSI representation of the image
func imageToANSI(_, _ uint, p termenv.Profile, img image.Image) string {
	b := img.Bounds()

	str := strings.Builder{}
	for y := 0; y < b.Max.Y; y += 2 {
		for x := 0; x < b.Max.X; x++ {
			// Use the upper half block character (▀)
			// The foreground color is the top pixel
			// The background color is the bottom pixel
			str.WriteString(termenv.String("▀").
				Foreground(p.FromColor(img.At(x, y))).
				Background(p.FromColor(img.At(x, y+1))).
				String())
		}
		str.WriteString("\n") // End of row
	}

	return str.String()
}

// greenscreen applies a virtual greenscreen effect to an image.
// It compares each pixel in the image to the corresponding pixel in a background
// image. If they are similar enough (within the distance threshold), the pixel
// is made transparent.
//
// Parameters:
//   - img: The foreground image to process (modified in-place)
//   - bg: The background image to compare against
//   - dist: The color distance threshold - lower values make fewer pixels transparent
//
// Note: This function assumes that img and bg have the same dimensions.
func greenscreen(img *image.RGBA, bg image.Image, dist float64) {
	if bg == nil {
		return
	}

	// Commented size check - could be uncommented for safety
	/*
		if img.Bounds().Size().Y != v.Bounds().Size().Y {
			panic(nil)
		}
		if img.Bounds().Size().X != v.Bounds().Size().X {
			panic(nil)
		}
	*/

	for y := 0; y < img.Bounds().Size().Y; y++ {
		for x := 0; x < img.Bounds().Size().X; x++ {
			// Convert to colorful.Color for better color distance calculation
			c1, _ := colorful.MakeColor(img.At(x, y))
			c2, _ := colorful.MakeColor(bg.At(x, y))

			// Potential future enhancement: face detection to avoid making faces transparent
			/*
				add face detection?
				if (x > 42 && x < 78) && (y > 5 && y < 40) {
					continue
				}
			*/

			// If colors are similar (within threshold), make pixel transparent
			if c1.DistanceLab(c2) < dist {
				img.Set(x, y, image.Transparent)
			}
		}
	}
}

// loadBgSamples loads a background sample image for use with the greenscreen effect.
// Currently, it loads a single sample image, but could be extended to average multiple samples.
//
// Parameters:
//   - path: The directory path where background samples are stored
//   - width: The desired width of the output image
//   - height: The desired height of the output image
//
// Returns:
//   - The loaded and resized background image
//   - An error if loading fails
//
// TODO: Enhance to take average of multiple sample images for better background detection
func loadBgSamples(path string, width, height uint) (image.Image, error) {
	//TODO: take average of sample set
	// for i := 40; i < 41; i++ {
	i := 40 // Currently only using a single sample
	b, err := ioutil.ReadFile(fmt.Sprintf("%s/%d.png", path, i))
	if err != nil {
		return nil, err
	}

	img, err := png.Decode(bytes.NewReader(b))
	if err != nil {
		return nil, err
	}

	// Resize the background image to match the terminal dimensions
	return resize.Resize(width, height, img, resize.Bilinear).(*image.RGBA), nil
	// }
}
