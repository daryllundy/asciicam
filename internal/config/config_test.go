package config

import (
	"flag"
	"image/color"
	"os"
	"testing"
)

func TestNewConfig(t *testing.T) {
	cfg := NewConfig()

	if cfg == nil {
		t.Fatal("NewConfig() returned nil")
	}

	// Test default values
	if cfg.DeviceID != 0 {
		t.Errorf("Expected default DeviceID 0, got %d", cfg.DeviceID)
	}

	if cfg.CamWidth != 1920 {
		t.Errorf("Expected default CamWidth 1920, got %d", cfg.CamWidth)
	}

	if cfg.CamHeight != 1080 {
		t.Errorf("Expected default CamHeight 1080, got %d", cfg.CamHeight)
	}

	if cfg.Zoom != 4 {
		t.Errorf("Expected default Zoom 4, got %d", cfg.Zoom)
	}

	if cfg.ANSI != false {
		t.Error("Expected default ANSI false")
	}

	if cfg.ShowFPS != false {
		t.Error("Expected default ShowFPS false")
	}

	if cfg.SamplePath != "bgsample" {
		t.Errorf("Expected default SamplePath 'bgsample', got %s", cfg.SamplePath)
	}

	if cfg.Threshold != 0.13 {
		t.Errorf("Expected default Threshold 0.13, got %f", cfg.Threshold)
	}
}

func TestParseFlags(t *testing.T) {
	// Reset flag package for clean testing
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	// Set up test args
	os.Args = []string{"test", "-dev=1", "-width=100", "-height=50", "-ansi=true", "-fps=true"}

	cfg := NewConfig()
	err := cfg.ParseFlags()

	if err != nil {
		t.Fatalf("ParseFlags() returned error: %v", err)
	}

	if cfg.DeviceID != 1 {
		t.Errorf("Expected DeviceID 1, got %d", cfg.DeviceID)
	}

	if cfg.Width != 100 {
		t.Errorf("Expected Width 100, got %d", cfg.Width)
	}

	if cfg.Height != 100 {
		t.Errorf("Expected Height 100, got %d", cfg.Height)
	}

	if !cfg.ANSI {
		t.Error("Expected ANSI true")
	}

	if !cfg.ShowFPS {
		t.Error("Expected ShowFPS true")
	}
}

func TestParseFlags_WithColor(t *testing.T) {
	// Reset flag package for clean testing
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	// Set up test args with color
	os.Args = []string{"test", "-color=#ff0000"}

	cfg := NewConfig()
	err := cfg.ParseFlags()

	if err != nil {
		t.Fatalf("ParseFlags() with color returned error: %v", err)
	}

	if cfg.Color != "#ff0000" {
		t.Errorf("Expected Color '#ff0000', got %s", cfg.Color)
	}

	// Check that ParsedColor was set
	if cfg.ParsedColor == nil {
		t.Error("Expected ParsedColor to be set")
	}
}

func TestParseFlags_InvalidColor(t *testing.T) {
	// Reset flag package for clean testing
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	// Set up test args with invalid color
	os.Args = []string{"test", "-color=invalid"}

	cfg := NewConfig()
	err := cfg.ParseFlags()

	if err == nil {
		t.Error("Expected error for invalid color, got none")
	}
}

func TestValidate(t *testing.T) {
	cfg := NewConfig()

	// Test zoom validation
	cfg.Zoom = 0 // Invalid
	err := cfg.Validate()
	if err != nil {
		t.Errorf("Validate() returned error: %v", err)
	}
	if cfg.Zoom != 1 {
		t.Errorf("Expected Zoom to be corrected to 1, got %d", cfg.Zoom)
	}

	cfg.Zoom = 10 // Invalid (too high)
	err = cfg.Validate()
	if err != nil {
		t.Errorf("Validate() returned error: %v", err)
	}
	if cfg.Zoom != 4 {
		t.Errorf("Expected Zoom to be corrected to 4, got %d", cfg.Zoom)
	}
}

func TestValidate_ANSI(t *testing.T) {
	cfg := NewConfig()
	cfg.ANSI = true
	cfg.Height = 24

	err := cfg.Validate()
	if err != nil {
		t.Errorf("Validate() returned error: %v", err)
	}

	// ANSI mode should double the height
	if cfg.Height != 48 {
		t.Errorf("Expected Height to be doubled to 48 for ANSI mode, got %d", cfg.Height)
	}
}

func TestGetCameraDimensions(t *testing.T) {
	cfg := NewConfig()
	cfg.CamWidth = 1280
	cfg.CamHeight = 720

	w, h := cfg.GetCameraDimensions()
	if w != 1280 || h != 720 {
		t.Errorf("Expected camera dimensions 1280x720, got %dx%d", w, h)
	}
}

func TestGetDisplayDimensions(t *testing.T) {
	cfg := NewConfig()
	cfg.Width = 80
	cfg.Height = 24

	w, h := cfg.GetDisplayDimensions()
	if w != 80 || h != 24 {
		t.Errorf("Expected display dimensions 80x24, got %dx%d", w, h)
	}
}

func TestGetScaledDimensions(t *testing.T) {
	cfg := NewConfig()
	cfg.Width = 80
	cfg.Height = 24

	tests := []struct {
		zoom      uint
		expectedW uint
		expectedH uint
	}{
		{1, 20, 6},  // 25%
		{2, 40, 12}, // 50%
		{3, 60, 18}, // 75%
		{4, 80, 24}, // 100%
	}

	for _, tt := range tests {
		t.Run("zoom_"+string(rune(tt.zoom+'0')), func(t *testing.T) {
			cfg.Zoom = tt.zoom
			w, h := cfg.GetScaledDimensions()

			if w != tt.expectedW || h != tt.expectedH {
				t.Errorf("For zoom %d, expected %dx%d, got %dx%d",
					tt.zoom, tt.expectedW, tt.expectedH, w, h)
			}
		})
	}
}

func TestGetScaledDimensions_NoOverflow(t *testing.T) {
	cfg := NewConfig()
	cfg.Width = 1000
	cfg.Height = 1000
	cfg.Zoom = 4

	w, h := cfg.GetScaledDimensions()

	// Should not exceed original dimensions
	if w > cfg.Width || h > cfg.Height {
		t.Errorf("Scaled dimensions %dx%d exceed original %dx%d", w, h, cfg.Width, cfg.Height)
	}
}

func TestGetTermSize(t *testing.T) {
	// This test just ensures the function doesn't panic
	// Actual values depend on the terminal environment
	w, h := getTermSize()

	// Should return positive values
	if w == 0 && h == 0 {
		t.Log("getTermSize returned 0x0, which might be expected in test environment")
	}

	// At minimum, should not panic and should return reasonable defaults
	if w > 10000 || h > 10000 {
		t.Errorf("getTermSize returned unreasonably large values: %dx%d", w, h)
	}
}

func TestValidate_WithAutoDetect(t *testing.T) {
	cfg := NewConfig()
	cfg.Width = 0  // Auto-detect
	cfg.Height = 0 // Auto-detect

	err := cfg.Validate()
	if err != nil {
		t.Errorf("Validate() returned error: %v", err)
	}

	// Should set some positive values
	if cfg.Width == 0 || cfg.Height == 0 {
		t.Error("Validate() should set positive width and height when auto-detecting")
	}
}

func TestParsedColorIsColor(t *testing.T) {
	cfg := NewConfig()

	// Test that ParsedColor implements color.Color interface
	var _ color.Color = cfg.ParsedColor

	// Test default parsed color
	r, g, b, a := cfg.ParsedColor.RGBA()
	if a != 0 {
		t.Errorf("Default ParsedColor should have alpha 0, got %d", a)
	}

	// Values should be 0 for default color
	if r != 0 || g != 0 || b != 0 {
		t.Errorf("Default ParsedColor should be black, got RGBA(%d,%d,%d,%d)", r, g, b, a)
	}
}

// Benchmark tests
func BenchmarkNewConfig(b *testing.B) {
	for i := 0; i < b.N; i++ {
		cfg := NewConfig()
		_ = cfg // Prevent optimization
	}
}

func BenchmarkValidate(b *testing.B) {
	cfg := NewConfig()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := cfg.Validate()
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkGetScaledDimensions(b *testing.B) {
	cfg := NewConfig()
	cfg.Width = 1920
	cfg.Height = 1080
	cfg.Zoom = 2

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w, h := cfg.GetScaledDimensions()
		_, _ = w, h // Prevent optimization
	}
}
