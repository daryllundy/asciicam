// Package config provides configuration management for the asciicam application.
package config

import (
	"flag"
	"fmt"
	"image/color"
	"os"

	"github.com/lucasb-eyer/go-colorful"
	"golang.org/x/term"
)

// Config holds all configuration options for the application.
type Config struct {
	// Camera settings
	DeviceID  int
	CamWidth  uint
	CamHeight uint

	// Display settings
	Width  uint
	Height uint
	Zoom   uint

	// Rendering settings
	ANSI     bool
	Color    string
	ShowFPS  bool

	// Greenscreen settings
	GenerateSamples bool
	UseGreenscreen  bool
	SamplePath      string
	Threshold       float64

	// Parsed color (internal use)
	ParsedColor color.Color
}

// NewConfig creates a new configuration with default values.
func NewConfig() *Config {
	return &Config{
		DeviceID:        0,
		CamWidth:        1920,
		CamHeight:       1080,
		Width:           0, // Auto-detect
		Height:          0, // Auto-detect
		Zoom:            4,
		ANSI:            false,
		Color:           "",
		ShowFPS:         false,
		GenerateSamples: false,
		UseGreenscreen:  false,
		SamplePath:      "bgsample",
		Threshold:       0.13,
		ParsedColor:     color.RGBA{0, 0, 0, 0}, // Alpha 0 means use truecolor
	}
}

// ParseFlags parses command line flags and updates the configuration.
func (c *Config) ParseFlags() error {
	deviceID := flag.Int("dev", c.DeviceID, "camera device ID (default: 0)")
	sample := flag.String("sample", c.SamplePath, "Where to find/store the sample data")
	gen := flag.Bool("gen", c.GenerateSamples, "Generate a new background")
	screen := flag.Bool("greenscreen", c.UseGreenscreen, "Use greenscreen")
	screenDist := flag.Float64("threshold", c.Threshold, "Greenscreen threshold")
	ansi := flag.Bool("ansi", c.ANSI, "Use ANSI")
	usecol := flag.String("color", c.Color, "Use single color")
	w := flag.Uint("width", c.Width, "output width")
	h := flag.Uint("height", c.Height, "output height")
	camWidth := flag.Uint("camWidth", c.CamWidth, "cam input width")
	camHeight := flag.Uint("camHeight", c.CamHeight, "cam input height")
	zoom := flag.Uint("zoom", c.Zoom, "image zoom level (1-4, where 1=25%, 2=50%, 3=75%, 4=100%)")
	showFPS := flag.Bool("fps", c.ShowFPS, "Show FPS")

	flag.Parse()

	// Update config with parsed values
	c.DeviceID = *deviceID
	c.SamplePath = *sample
	c.GenerateSamples = *gen
	c.UseGreenscreen = *screen
	c.Threshold = *screenDist
	c.ANSI = *ansi
	c.Color = *usecol
	c.Width = *w
	c.Height = *h
	c.CamWidth = *camWidth
	c.CamHeight = *camHeight
	c.Zoom = *zoom
	c.ShowFPS = *showFPS

	// Parse color if provided
	if c.Color != "" {
		col, err := colorful.Hex(c.Color)
		if err != nil {
			return fmt.Errorf("invalid color: %v", err)
		}
		c.ParsedColor = col
	}

	return c.Validate()
}

// Validate validates the configuration and sets reasonable defaults.
func (c *Config) Validate() error {
	// Validate zoom level (1-4)
	if c.Zoom < 1 {
		c.Zoom = 1
	} else if c.Zoom > 4 {
		c.Zoom = 4
	}

	// Auto-detect terminal size if not explicitly set
	if c.Width == 0 || c.Height == 0 {
		autoWidth, autoHeight := getTermSize()
		if c.Width == 0 {
			c.Width = autoWidth
		}
		if c.Height == 0 {
			c.Height = autoHeight
		}
	}

	// Set reasonable defaults if detection failed
	if c.Width == 0 {
		c.Width = 125
	}
	if c.Height == 0 {
		c.Height = 50
	}

	// ANSI rendering uses half-height blocks - adjust height
	if c.ANSI {
		c.Height *= 2
	}

	return nil
}

// getTermSize returns the current terminal dimensions.
func getTermSize() (width, height uint) {
	w, h := 0, 0
	var err error

	if term.IsTerminal(int(os.Stdout.Fd())) {
		w, h, err = term.GetSize(int(os.Stdout.Fd()))
		if err != nil {
			// Default fallback values if we can't get the terminal size
			w, h = 80, 24
		}
	} else {
		// Default values for non-terminal output
		w, h = 80, 24
	}

	return uint(w), uint(h)
}

// GetDisplayDimensions returns the calculated display dimensions.
func (c *Config) GetDisplayDimensions() (uint, uint) {
	return c.Width, c.Height
}

// GetCameraDimensions returns the camera capture dimensions.
func (c *Config) GetCameraDimensions() (uint, uint) {
	return c.CamWidth, c.CamHeight
}

// GetScaledDimensions returns the dimensions adjusted for zoom level.
func (c *Config) GetScaledDimensions() (uint, uint) {
	scaleFactor := float64(c.Zoom) / 4.0 // Convert zoom 1-4 to 0.25-1.0 range
	
	scaledWidth := uint(float64(c.Width) * scaleFactor)
	scaledHeight := uint(float64(c.Height) * scaleFactor)
	
	// Make sure we don't exceed the terminal dimensions
	if scaledWidth > c.Width {
		scaledWidth = c.Width
	}
	if scaledHeight > c.Height {
		scaledHeight = c.Height
	}
	
	return scaledWidth, scaledHeight
}