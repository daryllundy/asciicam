package errors

import (
	"errors"
	"testing"
)

func TestCameraError(t *testing.T) {
	originalErr := errors.New("original error")
	cameraErr := NewCameraError(1, "read", originalErr)

	expected := "camera error (device 1, op: read): original error"
	if cameraErr.Error() != expected {
		t.Errorf("Expected %s, got %s", expected, cameraErr.Error())
	}

	if !errors.Is(cameraErr, originalErr) {
		t.Error("Expected camera error to wrap original error")
	}
}

func TestConfigError(t *testing.T) {
	originalErr := errors.New("invalid value")
	configErr := NewConfigError("width", 0, originalErr)

	expected := "config error (field: width, value: 0): invalid value"
	if configErr.Error() != expected {
		t.Errorf("Expected %s, got %s", expected, configErr.Error())
	}

	if !errors.Is(configErr, originalErr) {
		t.Error("Expected config error to wrap original error")
	}
}

func TestFileError(t *testing.T) {
	originalErr := errors.New("permission denied")
	fileErr := NewFileError("/path/to/file", "read", originalErr)

	expected := "file error (path: /path/to/file, op: read): permission denied"
	if fileErr.Error() != expected {
		t.Errorf("Expected %s, got %s", expected, fileErr.Error())
	}

	if !errors.Is(fileErr, originalErr) {
		t.Error("Expected file error to wrap original error")
	}
}

func TestImageError(t *testing.T) {
	originalErr := errors.New("decode failed")
	imageErr := NewImageError("decode", "1920x1080", originalErr)

	expected := "image error (op: decode, dim: 1920x1080): decode failed"
	if imageErr.Error() != expected {
		t.Errorf("Expected %s, got %s", expected, imageErr.Error())
	}

	if !errors.Is(imageErr, originalErr) {
		t.Error("Expected image error to wrap original error")
	}
}

func TestIsRetryable(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{"nil error", nil, false},
		{"camera read failed", ErrCameraReadFailed, true},
		{"file read failed", ErrFileReadFailed, true},
		{"image process failed", ErrImageProcessFailed, true},
		{"camera not found", ErrCameraNotFound, false},
		{"other error", errors.New("other"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if IsRetryable(tt.err) != tt.expected {
				t.Errorf("IsRetryable(%v) = %v, expected %v", tt.err, IsRetryable(tt.err), tt.expected)
			}
		})
	}
}

func TestIsFatal(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{"nil error", nil, false},
		{"camera not found", ErrCameraNotFound, true},
		{"camera unsupported", ErrCameraUnsupported, true},
		{"invalid config", ErrInvalidConfig, true},
		{"terminal not tty", ErrTerminalNotTTY, true},
		{"camera read failed", ErrCameraReadFailed, false},
		{"other error", errors.New("other"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if IsFatal(tt.err) != tt.expected {
				t.Errorf("IsFatal(%v) = %v, expected %v", tt.err, IsFatal(tt.err), tt.expected)
			}
		})
	}
}

func TestErrorUnwrapping(t *testing.T) {
	originalErr := errors.New("root cause")
	cameraErr := NewCameraError(1, "test", originalErr)

	// Test direct unwrapping
	if errors.Unwrap(cameraErr) != originalErr {
		t.Error("Expected to unwrap to original error")
	}

	// Test errors.Is with wrapped error
	if !errors.Is(cameraErr, originalErr) {
		t.Error("Expected errors.Is to find wrapped error")
	}

	// Test with nested wrapping
	configErr := NewConfigError("field", "value", cameraErr)
	if !errors.Is(configErr, originalErr) {
		t.Error("Expected errors.Is to find nested wrapped error")
	}
}
