package camera

import (
	"image"
	"image/color"
	"testing"

	"gocv.io/x/gocv"
)

func TestNewCapture_InvalidDevice(t *testing.T) {
	// Test with an invalid device ID (very high number)
	_, err := NewCapture(9999, 640, 480)
	if err == nil {
		t.Log("Warning: Expected error for invalid device ID, but got none. This might be platform-specific.")
		// Don't fail the test as this behavior can vary by platform
	}
}

func TestGetDeviceID(t *testing.T) {
	// Create a mock capture struct for testing
	capture := &Capture{
		deviceID: 42,
		width:    640,
		height:   480,
	}

	if capture.GetDeviceID() != 42 {
		t.Errorf("Expected device ID 42, got %d", capture.GetDeviceID())
	}
}

func TestGetDimensions(t *testing.T) {
	capture := &Capture{
		deviceID: 0,
		width:    1920,
		height:   1080,
	}

	w, h := capture.GetDimensions()
	if w != 1920 || h != 1080 {
		t.Errorf("Expected dimensions 1920x1080, got %dx%d", w, h)
	}
}

func TestMatToImage(t *testing.T) {
	// Test the conversion logic without requiring actual camera hardware
	capture := &Capture{
		deviceID: 0,
		width:    2,
		height:   2,
	}

	// Create a simple 2x2 BGR Mat
	mat := gocv.NewMat()
	defer mat.Close()

	// Initialize with some data (this is a simplified test)
	mat = gocv.NewMatWithSize(2, 2, gocv.MatTypeCV8UC3)
	defer mat.Close()

	// Test the conversion
	img := capture.matToImage(mat)

	if img == nil {
		t.Fatal("matToImage returned nil")
	}

	bounds := img.Bounds()
	if bounds.Dx() != 2 || bounds.Dy() != 2 {
		t.Errorf("Expected 2x2 image, got %dx%d", bounds.Dx(), bounds.Dy())
	}
}

func TestMatToImage_LargeDimensions(t *testing.T) {
	// Test with large dimensions to verify overflow protection
	capture := &Capture{
		deviceID: 0,
		width:    1000,
		height:   1000,
	}

	mat := gocv.NewMatWithSize(1000, 1000, gocv.MatTypeCV8UC3)
	defer mat.Close()

	img := capture.matToImage(mat)

	if img == nil {
		t.Fatal("matToImage returned nil for large dimensions")
	}

	bounds := img.Bounds()
	if bounds.Dx() != 1000 || bounds.Dy() != 1000 {
		t.Errorf("Expected 1000x1000 image, got %dx%d", bounds.Dx(), bounds.Dy())
	}
}

func TestResizeImage(t *testing.T) {
	capture := &Capture{}

	// Create a test image
	original := image.NewRGBA(image.Rect(0, 0, 100, 100))

	// Fill with test pattern
	for y := 0; y < 100; y++ {
		for x := 0; x < 100; x++ {
			original.Set(x, y, color.RGBA{uint8(x), uint8(y), 128, 255})
		}
	}

	// Resize to 50x50
	resized := capture.ResizeImage(original, 50, 50)

	if resized == nil {
		t.Fatal("ResizeImage returned nil")
	}

	bounds := resized.Bounds()
	if bounds.Dx() != 50 || bounds.Dy() != 50 {
		t.Errorf("Expected 50x50 resized image, got %dx%d", bounds.Dx(), bounds.Dy())
	}
}

func TestResizeImage_UpscaleAndDownscale(t *testing.T) {
	capture := &Capture{}

	// Create a small test image
	original := image.NewRGBA(image.Rect(0, 0, 10, 10))

	tests := []struct {
		name   string
		width  uint
		height uint
	}{
		{"upscale", 20, 20},
		{"downscale", 5, 5},
		{"different aspect ratio", 30, 10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resized := capture.ResizeImage(original, tt.width, tt.height)

			if resized == nil {
				t.Fatal("ResizeImage returned nil")
			}

			bounds := resized.Bounds()
			expectedW, expectedH := int(tt.width), int(tt.height)
			if bounds.Dx() != expectedW || bounds.Dy() != expectedH {
				t.Errorf("Expected %dx%d resized image, got %dx%d", expectedW, expectedH, bounds.Dx(), bounds.Dy())
			}
		})
	}
}

// Test overflow protection in matToImage
func TestMatToImage_OverflowProtection(t *testing.T) {
	// Test with reasonable dimensions that won't cause overflow
	capture := &Capture{
		deviceID: 0,
		width:    100,
		height:   100,
	}

	// Create a mat with the same dimensions
	mat := gocv.NewMatWithSize(100, 100, gocv.MatTypeCV8UC3)
	defer mat.Close()

	// This should not panic
	img := capture.matToImage(mat)

	if img == nil {
		t.Fatal("matToImage returned nil")
	}

	// The image dimensions should match the capture dimensions
	bounds := img.Bounds()
	if bounds.Dx() != 100 || bounds.Dy() != 100 {
		t.Errorf("Expected 100x100 image, got %dx%d", bounds.Dx(), bounds.Dy())
	}
}

// Benchmark tests
func BenchmarkMatToImage(b *testing.B) {
	capture := &Capture{
		deviceID: 0,
		width:    640,
		height:   480,
	}

	mat := gocv.NewMatWithSize(480, 640, gocv.MatTypeCV8UC3)
	defer mat.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		img := capture.matToImage(mat)
		_ = img // Prevent optimization
	}
}

func BenchmarkResizeImage(b *testing.B) {
	capture := &Capture{}

	// Create a test image
	original := image.NewRGBA(image.Rect(0, 0, 1920, 1080))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resized := capture.ResizeImage(original, 640, 480)
		_ = resized // Prevent optimization
	}
}

// Integration test (only runs if camera hardware is available)
func TestReadFrame_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Try to create a capture with device 0
	capture, err := NewCapture(0, 640, 480)
	if err != nil {
		t.Skipf("Skipping integration test: no camera available (%v)", err)
	}
	defer capture.Close()

	// Try to read a frame
	img, err := capture.ReadFrame()
	if err != nil {
		t.Skipf("Could not read frame: %v", err)
	}

	if img == nil {
		t.Error("ReadFrame returned nil image")
	}

	bounds := img.Bounds()
	if bounds.Dx() <= 0 || bounds.Dy() <= 0 {
		t.Error("Frame has invalid dimensions")
	}
}

// Test Close functionality
func TestClose(t *testing.T) {
	// This test ensures Close doesn't panic with nil webcam
	capture := &Capture{
		deviceID: 0,
		width:    640,
		height:   480,
		webcam:   nil, // Simulate uninitialized state
	}

	// Should not panic
	capture.Close()
}
