package main

import (
	"context"
	"fmt"
	"image"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/muesli/asciicam/internal/ascii"
	"github.com/muesli/asciicam/internal/camera"
	"github.com/muesli/asciicam/internal/config"
	"github.com/muesli/asciicam/internal/greenscreen"
	"github.com/muesli/termenv"
)

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
	// Initialize configuration
	cfg := config.NewConfig()
	if err := cfg.ParseFlags(); err != nil {
		return fmt.Errorf("error parsing flags: %w", err)
	}

	// Initialize camera capture
	camWidth, camHeight := cfg.GetCameraDimensions()
	capture, err := camera.NewCapture(cfg.DeviceID, camWidth, camHeight)
	if err != nil {
		return fmt.Errorf("error initializing camera: %w", err)
	}
	defer capture.Close()

	// Initialize ASCII converter
	converter := ascii.NewConverter()
	if cfg.ParsedColor != nil {
		converter.SetGlobalColor(cfg.ParsedColor)
	}

	// Initialize greenscreen processor if needed
	var gsProcessor *greenscreen.Processor
	if cfg.UseGreenscreen || cfg.GenerateSamples {
		gsProcessor = greenscreen.NewProcessor(cfg.SamplePath, cfg.Threshold)
		if cfg.UseGreenscreen {
			termWidth, termHeight := cfg.GetDisplayDimensions()
			if err := gsProcessor.LoadBackground(termWidth, termHeight); err != nil {
				return fmt.Errorf("error loading background samples: %w", err)
			}
		}
	}

	// Get display dimensions
	termWidth, termHeight := cfg.GetDisplayDimensions()
	scaledWidth, scaledHeight := cfg.GetScaledDimensions()

	// Set up terminal
	output := termenv.NewOutput(os.Stdout)
	p := output.ColorProfile()
	output.HideCursor()
	defer output.ShowCursor()
	output.AltScreen()
	defer output.ExitAltScreen()

	// Clear screen at the beginning
	fmt.Print("\033[2J") // Clear entire screen
	fmt.Print("\033[H")  // Move cursor to the top-left corner

	// FPS tracking
	var fps []float64
	for i := 0; i < 10; i++ {
		fps = append(fps, 0)
	}

	frameCount := 0
	for {
		if ctx.Err() != nil {
			return nil
		}

		// Read frame from camera
		img, err := capture.ReadFrame()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading frame: %v\n", err)
			time.Sleep(100 * time.Millisecond)
			continue
		}

		// Handle background sample generation
		if cfg.GenerateSamples {
			if err := gsProcessor.GenerateSamples(img, frameCount); err != nil {
				return fmt.Errorf("error generating background sample: %w", err)
			}
			frameCount++
			if frameCount > 100 {
				fmt.Println("Generated 100 background samples. Exiting.")
				return nil
			}
			continue
		}

		// Resize image based on calculated dimensions
		resizedImg := capture.ResizeImage(img, scaledWidth, scaledHeight)

		// Apply greenscreen effect if enabled
		if cfg.UseGreenscreen && gsProcessor != nil {
			if rgbaImg, ok := resizedImg.(*image.RGBA); ok {
				gsProcessor.Apply(rgbaImg)
				resizedImg = rgbaImg
			}
		}

		// Convert to ASCII/ANSI
		now := time.Now()
		var output string
		if cfg.ANSI {
			output = converter.ImageToANSI(p, resizedImg)
		} else {
			output = converter.ImageToASCII(termWidth, termHeight, p, resizedImg)
		}

		// Render output
		fmt.Print("\033[H") // Move cursor to top-left (home)
		fmt.Print("\033[J") // Clear screen from cursor to end of screen
		fmt.Print(output)   // Print the ASCII/ANSI image

		// Update and display FPS if requested
		if cfg.ShowFPS {
			for i := len(fps) - 1; i > 0; i-- {
				fps[i] = fps[i-1]
			}
			fps[0] = float64(time.Second / time.Since(now))

			var fpsa float64
			for _, f := range fps {
				fpsa += f
			}

			// Calculate position for FPS display
			ansiHeightMultiplier := uint(1)
			if cfg.ANSI {
				ansiHeightMultiplier = 2
			}

			// Move cursor to the bottom and print FPS
			// Safe conversion with bounds checking
			const maxInt = int(^uint(0) >> 1)
			heightDiv := termHeight / ansiHeightMultiplier
			var cursorLine int
			if heightDiv > uint(maxInt-1) {
				cursorLine = maxInt // Cap at max int to avoid overflow
			} else {
				cursorLine = int(heightDiv) + 1
			}
			// Ensure cursor line is positive
			if cursorLine < 1 {
				cursorLine = 1
			}
			fmt.Printf("\033[%d;0H", cursorLine)
			fmt.Printf("FPS: %.0f", fpsa/float64(len(fps)))
		}
	}
}
