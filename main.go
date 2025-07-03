package main

import (
	"context"
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/lucasb-eyer/go-colorful"
	"github.com/muesli/termenv"
	"github.com/nfnt/resize"
	"gocv.io/x/gocv"
	"golang.org/x/term"
)

var (
	col    = color.Color(color.RGBA{0, 0, 0, 0}) // if alpha is 0, use truecolor
	pixels = []rune{' ', '.', ',', ':', ';', 'i', '1', 't', 'f', 'L', 'C', 'G', '0', '8', '@'}
)

// getTermSize returns the current terminal dimensions
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

func main() {
	// graceful shutdown on SIGINT, SIGTERM
	ctx, cancel := context.WithCancel(context.Background())
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs

		fmt.Println("\nShutting down...")
		cancel()
	}()

	if err := run(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	deviceID := flag.Int("dev", 0, "camera device ID (default: 0)")
	sample := flag.String("sample", "bgsample", "Where to find/store the sample data")
	gen := flag.Bool("gen", false, "Generate a new background")
	screen := flag.Bool("greenscreen", false, "Use greenscreen")
	screenDist := flag.Float64("threshold", 0.13, "Greenscreen threshold")
	ansi := flag.Bool("ansi", false, "Use ANSI")
	usecol := flag.String("color", "", "Use single color")
	w := flag.Uint("width", 0, "output width")
	h := flag.Uint("height", 0, "output height")
	camWidth := flag.Uint("camWidth", 1920, "cam input width")
	camHeight := flag.Uint("camHeight", 1080, "cam input height")
	zoom := flag.Uint("zoom", 4, "image zoom level (1-4, where 1=25%, 2=50%, 3=75%, 4=100%)")
	showFPS := flag.Bool("fps", false, "Show FPS")

	flag.Parse()
	if *usecol != "" {
		c, err := colorful.Hex(*usecol)
		if err != nil {
			return fmt.Errorf("invalid color: %v", err)
		}

		col = c
	}
	// Initialize width/height from command line
	termWidth := *w  // width of the terminal output
	termHeight := *h // height of the terminal output

	// If not explicitly set via command line, auto-detect from terminal
	if termWidth == 0 || termHeight == 0 {
		autoWidth, autoHeight := getTermSize()
		if termWidth == 0 {
			termWidth = autoWidth
		}
		if termHeight == 0 {
			termHeight = autoHeight
		}
	}

	// Set reasonable defaults if detection failed
	if termWidth == 0 {
		termWidth = 125
	}
	if termHeight == 0 {
		termHeight = 50
	}

	// ANSI rendering uses half-height blocks - adjust height
	ansiHeightMultiplier := uint(1)
	if *ansi {
		ansiHeightMultiplier = 2
		termHeight *= 2
	}

	// Store the last known terminal size to detect changes
	lastTermWidth, lastTermHeight := termWidth, termHeight

	// Open webcam with gocv
	webcam, err := gocv.OpenVideoCapture(*deviceID)
	if err != nil {
		return fmt.Errorf("error opening video capture device: %v", err)
	}
	defer webcam.Close()

	// Set camera properties
	webcam.Set(gocv.VideoCaptureFrameWidth, float64(*camWidth))
	webcam.Set(gocv.VideoCaptureFrameHeight, float64(*camHeight))

	// Check if camera opened correctly
	if !webcam.IsOpened() {
		return fmt.Errorf("camera device %d failed to open", *deviceID)
	}

	// Create a Mat to store frames
	frame := gocv.NewMat()
	defer frame.Close()

	var bg image.Image
	if !*gen && *screen {
		bg, err = loadBgSamples(*sample, termWidth, termHeight)
		if err != nil {
			return fmt.Errorf("could not load background samples: %w", err)
		}
	}

	p := termenv.EnvColorProfile()

	// Set up terminal
	termenv.HideCursor()
	defer termenv.ShowCursor()
	termenv.AltScreen()
	defer termenv.ExitAltScreen()

	// Clear screen at the beginning
	fmt.Print("\033[2J") // Clear entire screen
	fmt.Print("\033[H")  // Move cursor to the top-left corner

	// seed fps counter
	var fps []float64
	for i := 0; i < 10; i++ {
		fps = append(fps, 0)
	}

	i := 0
	for {
		if ctx.Err() != nil {
			return nil //nolint:nilerr
		}

		// Read a new frame from the webcam
		if ok := webcam.Read(&frame); !ok {
			fmt.Fprintln(os.Stderr, "Error reading frame")
			time.Sleep(100 * time.Millisecond)
			continue
		}

		// Skip empty frames
		if frame.Empty() {
			continue
		}

		// Convert gocv Mat to Go image
		img := matToImage(frame, *camWidth, *camHeight)

		// generate background sample data
		if *gen {
			// Make sure sample directory exists
			if err := os.MkdirAll(*sample, 0755); err != nil {
				return fmt.Errorf("failed to create sample directory: %w", err)
			}

			f, err := os.Create(fmt.Sprintf("%s/%d.png", *sample, i))
			if err != nil {
				return fmt.Errorf("failed to create sample file: %w", err)
			}
			if err := png.Encode(f, img); err != nil {
				return fmt.Errorf("failed to encode sample frame: %w", err)
			}
			_ = f.Close()

			i++
			if i > 100 {
				os.Exit(0)
			}
		}

		// Check for terminal resize
		if term.IsTerminal(int(os.Stdout.Fd())) {
			newWidth, newHeight := getTermSize()

			// If terminal size changed, update our dimensions
			if newWidth != lastTermWidth || newHeight != lastTermHeight {
				termWidth = newWidth
				termHeight = newHeight
				if *ansi {
					termHeight *= ansiHeightMultiplier
				}

				// Update the stored values
				lastTermWidth, lastTermHeight = newWidth, newHeight

				// Clear entire screen when size changes
				fmt.Print("\033[2J")
				fmt.Print("\033[H")
			}
		}

		// Apply zoom/scale factor and resize for display
		// Validate zoom level (1-4)
		if *zoom < 1 {
			*zoom = 1
		} else if *zoom > 4 {
			*zoom = 4
		}

		// Calculate dimensions to fill the entire terminal window
		// Use zoom to adjust the fill percentage
		scaleFactor := float64(*zoom) / 4.0 // Convert zoom 1-4 to 0.25-1.0 range

		// Always use the full terminal width (adjusted by zoom)
		scaledWidth := uint(float64(termWidth) * scaleFactor)

		// Calculate height - just use full terminal height adjusted by zoom factor
		scaledHeight := uint(float64(termHeight) * scaleFactor)

		// Make sure we don't exceed the terminal dimensions
		if scaledWidth > termWidth {
			scaledWidth = termWidth
		}
		if scaledHeight > termHeight {
			scaledHeight = termHeight
		}

		// Resize image based on calculated dimensions
		img = resize.Resize(scaledWidth, scaledHeight, img, resize.Bilinear).(*image.RGBA)

		// virtual green screen
		if !*gen && *screen {
			greenscreen(img, bg, *screenDist)
		}

		now := time.Now()
		// convert frame to ascii/ansi
		var s string
		if *ansi {
			s = imageToANSI(termWidth, termHeight, p, img)
		} else {
			s = imageToASCII(termWidth, termHeight, p, img)
		}

		// render - first clear the screen and then reposition cursor
		fmt.Print("\033[H")      // Move cursor to top-left (home)
		fmt.Print("\033[J")      // Clear screen from cursor to end of screen (more efficient than clearing whole screen)
		fmt.Fprint(os.Stdout, s) // Print the ASCII/ANSI image

		// Update FPS counter
		if *showFPS {
			for i := len(fps) - 1; i > 0; i-- {
				fps[i] = fps[i-1]
			}
			fps[0] = float64(time.Second / time.Since(now))

			var fpsa float64
			for _, f := range fps {
				fpsa += f
			}

			// Move cursor to the bottom and print FPS
			fmt.Printf("\033[%d;0H", int(termHeight/ansiHeightMultiplier)+1) // Move to the line after the image
			fmt.Printf("FPS: %.0f", fpsa/float64(len(fps)))
		}
	}
}
