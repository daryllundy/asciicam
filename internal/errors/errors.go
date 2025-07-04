// Package errors provides custom error types for the asciicam application.
package errors

import (
	"errors"
	"fmt"
)

// Error types for different failure modes
var (
	// Camera errors
	ErrCameraNotFound    = errors.New("camera device not found")
	ErrCameraInitFailed  = errors.New("failed to initialize camera")
	ErrCameraReadFailed  = errors.New("failed to read frame from camera")
	ErrCameraUnsupported = errors.New("camera device not supported")

	// Configuration errors
	ErrInvalidConfig     = errors.New("invalid configuration")
	ErrConfigParseFailed = errors.New("failed to parse configuration")
	ErrInvalidColorCode  = errors.New("invalid color code")
	ErrInvalidDimensions = errors.New("invalid dimensions")

	// File operation errors
	ErrFileNotFound    = errors.New("file not found")
	ErrFileReadFailed  = errors.New("failed to read file")
	ErrFileWriteFailed = errors.New("failed to write file")
	ErrDirCreateFailed = errors.New("failed to create directory")

	// Image processing errors
	ErrImageProcessFailed = errors.New("failed to process image")
	ErrImageResizeFailed  = errors.New("failed to resize image")
	ErrImageDecodeFailed  = errors.New("failed to decode image")
	ErrImageEncodeFailed  = errors.New("failed to encode image")

	// Greenscreen errors
	ErrGreenscreenLoadFailed  = errors.New("failed to load greenscreen background")
	ErrGreenscreenApplyFailed = errors.New("failed to apply greenscreen effect")
	ErrSampleGenerateFailed   = errors.New("failed to generate background sample")

	// Terminal errors
	ErrTerminalSizeFailed = errors.New("failed to get terminal size")
	ErrTerminalNotTTY     = errors.New("not running in a terminal")
)

// CameraError represents camera-related errors with additional context
type CameraError struct {
	DeviceID int
	Op       string
	Err      error
}

func (e *CameraError) Error() string {
	return fmt.Sprintf("camera error (device %d, op: %s): %v", e.DeviceID, e.Op, e.Err)
}

func (e *CameraError) Unwrap() error {
	return e.Err
}

// NewCameraError creates a new camera error with context
func NewCameraError(deviceID int, operation string, err error) *CameraError {
	return &CameraError{
		DeviceID: deviceID,
		Op:       operation,
		Err:      err,
	}
}

// ConfigError represents configuration-related errors
type ConfigError struct {
	Field string
	Value interface{}
	Err   error
}

func (e *ConfigError) Error() string {
	return fmt.Sprintf("config error (field: %s, value: %v): %v", e.Field, e.Value, e.Err)
}

func (e *ConfigError) Unwrap() error {
	return e.Err
}

// NewConfigError creates a new configuration error with context
func NewConfigError(field string, value interface{}, err error) *ConfigError {
	return &ConfigError{
		Field: field,
		Value: value,
		Err:   err,
	}
}

// FileError represents file operation errors
type FileError struct {
	Path string
	Op   string
	Err  error
}

func (e *FileError) Error() string {
	return fmt.Sprintf("file error (path: %s, op: %s): %v", e.Path, e.Op, e.Err)
}

func (e *FileError) Unwrap() error {
	return e.Err
}

// NewFileError creates a new file error with context
func NewFileError(path, operation string, err error) *FileError {
	return &FileError{
		Path: path,
		Op:   operation,
		Err:  err,
	}
}

// ImageError represents image processing errors
type ImageError struct {
	Operation string
	Dimension string
	Err       error
}

func (e *ImageError) Error() string {
	return fmt.Sprintf("image error (op: %s, dim: %s): %v", e.Operation, e.Dimension, e.Err)
}

func (e *ImageError) Unwrap() error {
	return e.Err
}

// NewImageError creates a new image processing error with context
func NewImageError(operation, dimension string, err error) *ImageError {
	return &ImageError{
		Operation: operation,
		Dimension: dimension,
		Err:       err,
	}
}

// IsRetryable determines if an error is potentially retryable
func IsRetryable(err error) bool {
	if err == nil {
		return false
	}

	// Check for specific retryable errors
	switch {
	case errors.Is(err, ErrCameraReadFailed):
		return true
	case errors.Is(err, ErrFileReadFailed):
		return true
	case errors.Is(err, ErrImageProcessFailed):
		return true
	default:
		return false
	}
}

// IsFatal determines if an error is fatal and should stop execution
func IsFatal(err error) bool {
	if err == nil {
		return false
	}

	// Check for fatal errors
	switch {
	case errors.Is(err, ErrCameraNotFound):
		return true
	case errors.Is(err, ErrCameraUnsupported):
		return true
	case errors.Is(err, ErrInvalidConfig):
		return true
	case errors.Is(err, ErrTerminalNotTTY):
		return true
	default:
		return false
	}
}