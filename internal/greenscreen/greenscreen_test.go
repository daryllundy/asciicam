package greenscreen

import (
	"image"
	"image/color"
	"os"
	"path/filepath"
	"testing"
)

func TestNewProcessor(t *testing.T) {
	processor := NewProcessor("test_samples", 0.15)

	if processor == nil {
		t.Fatal("NewProcessor() returned nil")
	}

	if processor.samplePath != "test_samples" {
		t.Errorf("Expected samplePath 'test_samples', got %s", processor.samplePath)
	}

	if processor.threshold != 0.15 {
		t.Errorf("Expected threshold 0.15, got %f", processor.threshold)
	}

	if processor.background != nil {
		t.Error("Expected background to be nil initially")
	}
}

func TestGetThreshold(t *testing.T) {
	processor := NewProcessor("test", 0.25)

	if processor.GetThreshold() != 0.25 {
		t.Errorf("Expected threshold 0.25, got %f", processor.GetThreshold())
	}
}

func TestSetThreshold(t *testing.T) {
	processor := NewProcessor("test", 0.1)

	processor.SetThreshold(0.3)

	if processor.GetThreshold() != 0.3 {
		t.Errorf("Expected threshold 0.3 after setting, got %f", processor.GetThreshold())
	}
}

func TestGetSamplePath(t *testing.T) {
	processor := NewProcessor("my_samples", 0.1)

	if processor.GetSamplePath() != "my_samples" {
		t.Errorf("Expected sample path 'my_samples', got %s", processor.GetSamplePath())
	}
}

func TestSetSamplePath(t *testing.T) {
	processor := NewProcessor("old_path", 0.1)

	processor.SetSamplePath("new_path")

	if processor.GetSamplePath() != "new_path" {
		t.Errorf("Expected sample path 'new_path' after setting, got %s", processor.GetSamplePath())
	}
}

func TestHasBackground(t *testing.T) {
	processor := NewProcessor("test", 0.1)

	if processor.HasBackground() {
		t.Error("Expected HasBackground() to return false initially")
	}

	// Set a background (simulate loading)
	processor.background = image.NewRGBA(image.Rect(0, 0, 10, 10))

	if !processor.HasBackground() {
		t.Error("Expected HasBackground() to return true after setting background")
	}
}

func TestApply_NoBackground(t *testing.T) {
	processor := NewProcessor("test", 0.1)

	// Create a test image
	img := image.NewRGBA(image.Rect(0, 0, 4, 4))
	originalColor := color.RGBA{255, 0, 0, 255}
	img.Set(0, 0, originalColor)

	// Apply should do nothing when no background is set
	processor.Apply(img)

	// Color should be unchanged
	resultColor := img.RGBAAt(0, 0)
	if resultColor != originalColor {
		t.Error("Apply() should not modify image when no background is set")
	}
}

func TestApply_WithBackground(t *testing.T) {
	processor := NewProcessor("test", 0.05) // Low threshold for easier testing

	// Create foreground and background images
	fg := image.NewRGBA(image.Rect(0, 0, 2, 2))
	bg := image.NewRGBA(image.Rect(0, 0, 2, 2))

	// Set similar colors (should be made transparent)
	similarColor := color.RGBA{100, 100, 100, 255}
	fg.Set(0, 0, similarColor)
	bg.Set(0, 0, similarColor)

	// Set different colors (should remain)
	fg.Set(1, 1, color.RGBA{255, 0, 0, 255})
	bg.Set(1, 1, color.RGBA{0, 255, 0, 255})

	processor.background = bg
	processor.Apply(fg)

	// Check that similar color was made transparent
	resultColor := fg.RGBAAt(0, 0)
	if resultColor.A != 0 {
		t.Error("Similar color should have been made transparent")
	}

	// Check that different color remains
	differentColor := fg.RGBAAt(1, 1)
	if differentColor.A == 0 {
		t.Error("Different color should not have been made transparent")
	}
}

func TestApply_WithHighThreshold(t *testing.T) {
	processor := NewProcessor("test", 0.8) // High threshold

	// Create images with moderately different colors
	fg := image.NewRGBA(image.Rect(0, 0, 2, 2))
	bg := image.NewRGBA(image.Rect(0, 0, 2, 2))

	fg.Set(0, 0, color.RGBA{100, 100, 100, 255})
	bg.Set(0, 0, color.RGBA{120, 120, 120, 255}) // Slightly different

	processor.background = bg
	processor.Apply(fg)

	// With high threshold, even moderately different colors should be made transparent
	resultColor := fg.RGBAAt(0, 0)
	if resultColor.A != 0 {
		t.Error("With high threshold, moderately different colors should be made transparent")
	}
}

func TestGenerateSamples(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()

	processor := NewProcessor(tempDir, 0.1)

	// Create a test image
	img := image.NewRGBA(image.Rect(0, 0, 10, 10))
	for y := 0; y < 10; y++ {
		for x := 0; x < 10; x++ {
			img.Set(x, y, color.RGBA{uint8(x * 25), uint8(y * 25), 128, 255})
		}
	}

	// Generate a sample
	err := processor.GenerateSamples(img, 42)
	if err != nil {
		t.Fatalf("GenerateSamples() returned error: %v", err)
	}

	// Check that the file was created
	expectedPath := filepath.Join(tempDir, "42.png")
	if _, err := os.Stat(expectedPath); os.IsNotExist(err) {
		t.Errorf("Expected sample file %s was not created", expectedPath)
	}
}

func TestGenerateSamples_InvalidPath(t *testing.T) {
	// Use a path that can't be created (inside a non-existent parent)
	invalidPath := "/nonexistent/path/that/cannot/be/created"
	processor := NewProcessor(invalidPath, 0.1)

	img := image.NewRGBA(image.Rect(0, 0, 2, 2))

	err := processor.GenerateSamples(img, 1)
	if err == nil {
		t.Error("Expected error for invalid path, got none")
	}
}

func TestLoadBackground_NonexistentFile(t *testing.T) {
	// Use a path that doesn't exist
	processor := NewProcessor("nonexistent_path", 0.1)

	err := processor.LoadBackground(10, 10)
	if err == nil {
		t.Error("Expected error for nonexistent background file, got none")
	}
}

func TestLoadBackground_WithValidFile(t *testing.T) {
	// Create a temporary directory and sample file
	tempDir := t.TempDir()

	processor := NewProcessor(tempDir, 0.1)

	// First generate a sample file
	img := image.NewRGBA(image.Rect(0, 0, 20, 20))
	err := processor.GenerateSamples(img, 40) // Creates 40.png
	if err != nil {
		t.Fatalf("Failed to generate sample: %v", err)
	}

	// Now try to load it
	err = processor.LoadBackground(20, 20)
	if err != nil {
		t.Fatalf("LoadBackground() returned error: %v", err)
	}

	if !processor.HasBackground() {
		t.Error("Expected background to be loaded")
	}
}

// Benchmark tests
func BenchmarkApply(b *testing.B) {
	processor := NewProcessor("test", 0.1)

	// Create test images
	fg := image.NewRGBA(image.Rect(0, 0, 100, 100))
	bg := image.NewRGBA(image.Rect(0, 0, 100, 100))

	// Fill with test data
	for y := 0; y < 100; y++ {
		for x := 0; x < 100; x++ {
			fg.Set(x, y, color.RGBA{uint8(x), uint8(y), 128, 255})
			bg.Set(x, y, color.RGBA{uint8(x + 10), uint8(y + 10), 128, 255})
		}
	}

	processor.background = bg

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Create a fresh copy for each iteration
		testImg := image.NewRGBA(fg.Bounds())
		for y := 0; y < 100; y++ {
			for x := 0; x < 100; x++ {
				testImg.Set(x, y, fg.At(x, y))
			}
		}
		processor.Apply(testImg)
	}
}

func BenchmarkGenerateSamples(b *testing.B) {
	tempDir := b.TempDir()
	processor := NewProcessor(tempDir, 0.1)

	// Create a test image
	img := image.NewRGBA(image.Rect(0, 0, 50, 50))
	for y := 0; y < 50; y++ {
		for x := 0; x < 50; x++ {
			img.Set(x, y, color.RGBA{uint8(x), uint8(y), 128, 255})
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := processor.GenerateSamples(img, i)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func TestApply_EdgeCases(t *testing.T) {
	processor := NewProcessor("test", 0.1)

	// Test with nil image (should not panic)
	t.Run("nil_image", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Apply() panicked with nil background: %v", r)
			}
		}()
		processor.Apply(nil) // Should not panic
	})

	// Test with empty image
	t.Run("empty_image", func(t *testing.T) {
		emptyImg := image.NewRGBA(image.Rect(0, 0, 0, 0))
		processor.Apply(emptyImg) // Should not panic
	})

	// Test with mismatched sizes
	t.Run("mismatched_sizes", func(t *testing.T) {
		fg := image.NewRGBA(image.Rect(0, 0, 10, 10))
		bg := image.NewRGBA(image.Rect(0, 0, 5, 5)) // Different size

		processor.background = bg

		// This might panic or produce undefined behavior, but shouldn't crash the test
		defer func() {
			if r := recover(); r != nil {
				t.Logf("Apply() panicked with mismatched sizes (expected): %v", r)
			}
		}()
		processor.Apply(fg)
	})
}
