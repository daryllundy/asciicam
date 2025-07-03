// Package greenscreen provides virtual greenscreen functionality for the asciicam application.
package greenscreen

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"io/ioutil"
	"os"

	"github.com/lucasb-eyer/go-colorful"
	"github.com/nfnt/resize"
)

// Processor handles greenscreen operations.
type Processor struct {
	samplePath string
	threshold  float64
	background image.Image
}

// NewProcessor creates a new greenscreen processor.
func NewProcessor(samplePath string, threshold float64) *Processor {
	return &Processor{
		samplePath: samplePath,
		threshold:  threshold,
	}
}

// LoadBackground loads the background sample image for greenscreen processing.
func (p *Processor) LoadBackground(width, height uint) error {
	bg, err := p.loadBgSamples(width, height)
	if err != nil {
		return err
	}
	p.background = bg
	return nil
}

// Apply applies the greenscreen effect to an image.
// It compares each pixel in the image to the corresponding pixel in the background
// image. If they are similar enough (within the distance threshold), the pixel
// is made transparent.
func (p *Processor) Apply(img *image.RGBA) {
	if p.background == nil {
		return
	}

	for y := 0; y < img.Bounds().Size().Y; y++ {
		for x := 0; x < img.Bounds().Size().X; x++ {
			// Convert to colorful.Color for better color distance calculation
			c1, _ := colorful.MakeColor(img.At(x, y))
			c2, _ := colorful.MakeColor(p.background.At(x, y))

			// If colors are similar (within threshold), make pixel transparent
			if c1.DistanceLab(c2) < p.threshold {
				img.Set(x, y, image.Transparent)
			}
		}
	}
}

// GenerateSamples generates background sample images for greenscreen processing.
func (p *Processor) GenerateSamples(img image.Image, frameNumber int) error {
	// Make sure sample directory exists
	if err := os.MkdirAll(p.samplePath, 0755); err != nil {
		return fmt.Errorf("failed to create sample directory: %w", err)
	}

	f, err := os.Create(fmt.Sprintf("%s/%d.png", p.samplePath, frameNumber))
	if err != nil {
		return fmt.Errorf("failed to create sample file: %w", err)
	}
	defer f.Close()

	if err := png.Encode(f, img); err != nil {
		return fmt.Errorf("failed to encode sample frame: %w", err)
	}

	return nil
}

// loadBgSamples loads a background sample image for use with the greenscreen effect.
// Currently, it loads a single sample image, but could be extended to average multiple samples.
func (p *Processor) loadBgSamples(width, height uint) (image.Image, error) {
	// TODO: take average of sample set
	// Currently only using a single sample
	i := 40
	b, err := ioutil.ReadFile(fmt.Sprintf("%s/%d.png", p.samplePath, i))
	if err != nil {
		return nil, fmt.Errorf("failed to read background sample: %w", err)
	}

	img, err := png.Decode(bytes.NewReader(b))
	if err != nil {
		return nil, fmt.Errorf("failed to decode background sample: %w", err)
	}

	// Resize the background image to match the terminal dimensions
	return resize.Resize(width, height, img, resize.Bilinear).(*image.RGBA), nil
}

// GetThreshold returns the current threshold value.
func (p *Processor) GetThreshold() float64 {
	return p.threshold
}

// SetThreshold sets the threshold value for greenscreen processing.
func (p *Processor) SetThreshold(threshold float64) {
	p.threshold = threshold
}

// GetSamplePath returns the path where background samples are stored.
func (p *Processor) GetSamplePath() string {
	return p.samplePath
}

// SetSamplePath sets the path where background samples are stored.
func (p *Processor) SetSamplePath(path string) {
	p.samplePath = path
}

// HasBackground returns true if a background image has been loaded.
func (p *Processor) HasBackground() bool {
	return p.background != nil
}