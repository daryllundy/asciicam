// Package camera provides functionality for webcam capture and image processing.
package camera

import (
	"fmt"
	"image"
	"image/color"

	"github.com/nfnt/resize"
	"gocv.io/x/gocv"
)

// Capture handles webcam capture operations.
type Capture struct {
	webcam   *gocv.VideoCapture
	deviceID int
	width    uint
	height   uint
}

// NewCapture creates a new camera capture instance.
func NewCapture(deviceID int, width, height uint) (*Capture, error) {
	webcam, err := gocv.OpenVideoCapture(deviceID)
	if err != nil {
		return nil, fmt.Errorf("error opening video capture device: %v", err)
	}

	// Set camera properties
	webcam.Set(gocv.VideoCaptureFrameWidth, float64(width))
	webcam.Set(gocv.VideoCaptureFrameHeight, float64(height))

	// Check if camera opened correctly
	if !webcam.IsOpened() {
		webcam.Close()
		return nil, fmt.Errorf("camera device %d failed to open", deviceID)
	}

	return &Capture{
		webcam:   webcam,
		deviceID: deviceID,
		width:    width,
		height:   height,
	}, nil
}

// Close closes the camera capture.
func (c *Capture) Close() {
	if c.webcam != nil {
		c.webcam.Close()
	}
}

// ReadFrame reads a frame from the webcam and returns it as an image.
func (c *Capture) ReadFrame() (image.Image, error) {
	frame := gocv.NewMat()
	defer frame.Close()

	// Read a new frame from the webcam
	if ok := c.webcam.Read(&frame); !ok {
		return nil, fmt.Errorf("error reading frame from camera")
	}

	// Skip empty frames
	if frame.Empty() {
		return nil, fmt.Errorf("empty frame received")
	}

	// Convert gocv Mat to Go image
	img := c.matToImage(frame)
	return img, nil
}

// matToImage converts a gocv Mat to Go's native image.RGBA format.
// This is necessary because gocv uses OpenCV's BGR color format, while
// Go's standard image library uses RGBA format.
func (c *Capture) matToImage(mat gocv.Mat) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, int(c.width), int(c.height)))

	// Create a copy of the Mat to avoid modifying the original
	bgrMat := mat.Clone()
	defer bgrMat.Close() // Ensure the Mat is properly closed to avoid memory leaks

	// Copy the data from mat to img, converting BGR to RGBA
	for y := 0; y < int(c.height); y++ {
		for x := 0; x < int(c.width); x++ {
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

// ResizeImage resizes an image to the specified dimensions.
func (c *Capture) ResizeImage(img image.Image, width, height uint) image.Image {
	return resize.Resize(width, height, img, resize.Bilinear)
}

// GetDeviceID returns the device ID of the camera.
func (c *Capture) GetDeviceID() int {
	return c.deviceID
}

// GetDimensions returns the width and height of the camera.
func (c *Capture) GetDimensions() (uint, uint) {
	return c.width, c.height
}