package ascii

import (
	"image"
	"image/color"
	"strings"
	"testing"

	"github.com/muesli/termenv"
)

func TestNewConverter(t *testing.T) {
	converter := NewConverter()
	if converter == nil {
		t.Fatal("NewConverter() returned nil")
	}

	if len(converter.pixels) == 0 {
		t.Error("Converter should have default pixels")
	}

	expectedPixels := []rune{' ', '.', ',', ':', ';', 'i', '1', 't', 'f', 'L', 'C', 'G', '0', '8', '@'}
	if len(converter.pixels) != len(expectedPixels) {
		t.Errorf("Expected %d pixels, got %d", len(expectedPixels), len(converter.pixels))
	}
}

func TestSetGlobalColor(t *testing.T) {
	converter := NewConverter()
	testColor := color.RGBA{255, 0, 0, 255} // Red

	converter.SetGlobalColor(testColor)

	// Test that the color was set (we can't directly access it, but we can test behavior)
	// This is tested indirectly through the ImageToASCII function
}

func TestPixelToASCII(t *testing.T) {
	converter := NewConverter()

	tests := []struct {
		name     string
		pixel    color.Color
		expected bool // whether we expect a valid ASCII character
	}{
		{
			name:     "black pixel",
			pixel:    color.RGBA{0, 0, 0, 255},
			expected: true,
		},
		{
			name:     "white pixel",
			pixel:    color.RGBA{255, 255, 255, 255},
			expected: true,
		},
		{
			name:     "transparent pixel",
			pixel:    color.RGBA{128, 128, 128, 0},
			expected: true,
		},
		{
			name:     "red pixel",
			pixel:    color.RGBA{255, 0, 0, 255},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := converter.pixelToASCII(tt.pixel)

			if tt.expected {
				// Check that result is one of our expected ASCII characters
				found := false
				for _, p := range converter.pixels {
					if result == p {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("pixelToASCII returned unexpected character: %c", result)
				}
			}
		})
	}
}

func TestImageToASCII(t *testing.T) {
	converter := NewConverter()

	// Create a simple 2x2 test image
	img := image.NewRGBA(image.Rect(0, 0, 2, 2))
	img.Set(0, 0, color.RGBA{0, 0, 0, 255})       // Black
	img.Set(1, 0, color.RGBA{255, 255, 255, 255}) // White
	img.Set(0, 1, color.RGBA{128, 128, 128, 255}) // Gray
	img.Set(1, 1, color.RGBA{255, 0, 0, 255})     // Red

	profile := termenv.ANSI
	result := converter.ImageToASCII(2, 2, profile, img)

	// Check that we got a string result
	if result == "" {
		t.Error("ImageToASCII returned empty string")
	}

	// Check that result contains newlines (should have 2 lines)
	lines := strings.Split(strings.TrimRight(result, "\n"), "\n")
	if len(lines) != 2 {
		t.Errorf("Expected 2 lines, got %d", len(lines))
	}

	// Each line should have some content (ASCII chars + ANSI codes)
	for i, line := range lines {
		if len(line) == 0 {
			t.Errorf("Line %d is empty", i)
		}
	}
}

func TestImageToANSI(t *testing.T) {
	converter := NewConverter()

	// Create a simple 2x4 test image (even height for ANSI blocks)
	img := image.NewRGBA(image.Rect(0, 0, 2, 4))
	for y := 0; y < 4; y++ {
		for x := 0; x < 2; x++ {
			// Create a gradient
			intensity := uint8(255 * y / 3)
			img.Set(x, y, color.RGBA{intensity, intensity, intensity, 255})
		}
	}

	profile := termenv.ANSI
	result := converter.ImageToANSI(profile, img)

	// Check that we got a string result
	if result == "" {
		t.Error("ImageToANSI returned empty string")
	}

	// Check that result contains the upper half block character (▀)
	if !strings.Contains(result, "▀") {
		t.Error("ImageToANSI should contain upper half block character (▀)")
	}

	// Should have 2 lines (4 pixel rows / 2 = 2 block rows)
	lines := strings.Split(strings.TrimRight(result, "\n"), "\n")
	if len(lines) != 2 {
		t.Errorf("Expected 2 lines, got %d", len(lines))
	}
}

func TestImageToASCII_WithGlobalColor(t *testing.T) {
	converter := NewConverter()
	globalColor := color.RGBA{255, 0, 0, 255} // Red with full alpha
	converter.SetGlobalColor(globalColor)

	// Create a simple test image
	img := image.NewRGBA(image.Rect(0, 0, 2, 2))
	img.Set(0, 0, color.RGBA{0, 0, 0, 255})
	img.Set(1, 0, color.RGBA{255, 255, 255, 255})
	img.Set(0, 1, color.RGBA{128, 128, 128, 255})
	img.Set(1, 1, color.RGBA{0, 255, 0, 255})

	profile := termenv.ANSI
	result := converter.ImageToASCII(2, 2, profile, img)

	// Should still produce output
	if result == "" {
		t.Error("ImageToASCII with global color returned empty string")
	}

	lines := strings.Split(strings.TrimRight(result, "\n"), "\n")
	if len(lines) != 2 {
		t.Errorf("Expected 2 lines, got %d", len(lines))
	}
}

func TestImageToASCII_LargeDimensions(t *testing.T) {
	converter := NewConverter()

	// Test with large dimensions to check overflow protection
	img := image.NewRGBA(image.Rect(0, 0, 100, 100))
	for y := 0; y < 100; y++ {
		for x := 0; x < 100; x++ {
			img.Set(x, y, color.RGBA{uint8(x), uint8(y), 128, 255})
		}
	}

	profile := termenv.ANSI
	result := converter.ImageToASCII(100, 100, profile, img)

	if result == "" {
		t.Error("ImageToASCII with large dimensions returned empty string")
	}

	lines := strings.Split(strings.TrimRight(result, "\n"), "\n")
	if len(lines) != 100 {
		t.Errorf("Expected 100 lines, got %d", len(lines))
	}
}

// Benchmark tests
func BenchmarkPixelToASCII(b *testing.B) {
	converter := NewConverter()
	testPixel := color.RGBA{128, 128, 128, 255}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		converter.pixelToASCII(testPixel)
	}
}

func BenchmarkImageToASCII(b *testing.B) {
	converter := NewConverter()
	img := image.NewRGBA(image.Rect(0, 0, 80, 24))

	// Fill with test data
	for y := 0; y < 24; y++ {
		for x := 0; x < 80; x++ {
			img.Set(x, y, color.RGBA{uint8(x), uint8(y), 128, 255})
		}
	}

	profile := termenv.ANSI

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		converter.ImageToASCII(80, 24, profile, img)
	}
}

func BenchmarkImageToANSI(b *testing.B) {
	converter := NewConverter()
	img := image.NewRGBA(image.Rect(0, 0, 80, 48)) // Even height for ANSI

	// Fill with test data
	for y := 0; y < 48; y++ {
		for x := 0; x < 80; x++ {
			img.Set(x, y, color.RGBA{uint8(x), uint8(y), 128, 255})
		}
	}

	profile := termenv.ANSI

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		converter.ImageToANSI(profile, img)
	}
}
