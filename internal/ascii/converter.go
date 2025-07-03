// Package ascii provides functionality for converting images to ASCII and ANSI art.
package ascii

import (
	"image"
	"image/color"
	"math"
	"strings"

	"github.com/muesli/termenv"
)

// Converter handles the conversion of images to ASCII/ANSI art.
type Converter struct {
	// pixels defines the ASCII characters used for different intensity levels
	pixels []rune
	// globalColor is the global color to use for ASCII output (if set)
	globalColor color.Color
}

// NewConverter creates a new ASCII converter with default settings.
func NewConverter() *Converter {
	return &Converter{
		pixels:      []rune{' ', '.', ',', ':', ';', 'i', '1', 't', 'f', 'L', 'C', 'G', '0', '8', '@'},
		globalColor: color.Color(color.RGBA{0, 0, 0, 0}), // alpha 0 means use truecolor
	}
}

// SetGlobalColor sets a global color for ASCII output.
func (c *Converter) SetGlobalColor(col color.Color) {
	c.globalColor = col
}

// pixelToASCII converts a color pixel to an ASCII character based on its intensity.
// Darker pixels are represented by characters with less "ink" (like spaces or dots),
// while brighter pixels use more "ink-heavy" characters (like @ or 8).
func (c *Converter) pixelToASCII(pixel color.Color) rune {
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
	precision := float64(255 * 3 / (len(c.pixels) - 1))

	// Map intensity to an index in the pixels array
	v := int(math.Floor(float64(intensity)/precision + 0.5))
	return c.pixels[v]
}

// ImageToASCII converts an image to ASCII art with color.
// Each pixel is represented by an ASCII character with the appropriate color.
func (c *Converter) ImageToASCII(width, height uint, p termenv.Profile, img image.Image) string {
	str := strings.Builder{}

	for i := 0; i < int(height); i++ {
		for j := 0; j < int(width); j++ {
			// Get pixel and convert to ASCII character
			pixel := color.NRGBAModel.Convert(img.At(j, i))
			s := termenv.String(string(c.pixelToASCII(pixel)))

			// Apply color - either the global color (if set) or the pixel's color
			_, _, _, a := c.globalColor.RGBA()
			if a > 0 {
				// Use global color if it has been set
				s = s.Foreground(p.FromColor(c.globalColor))
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

// ImageToANSI converts an image to colored ANSI blocks.
// It uses the upper half block character (▀) with foreground and background
// colors to represent two pixels vertically in a single character position.
// This provides higher vertical resolution than ASCII art.
func (c *Converter) ImageToANSI(p termenv.Profile, img image.Image) string {
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
