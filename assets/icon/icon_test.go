package icon

import (
	"bytes"
	"testing"
)

// PNG magic number for validation
var pngHeader = []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}

func TestIconOpenIsValidPNG(t *testing.T) {
	if len(IconOpen) == 0 {
		t.Fatal("IconOpen is empty")
	}
	if !bytes.HasPrefix(IconOpen, pngHeader) {
		t.Fatal("IconOpen does not start with PNG header")
	}
}

func TestIconClosedIsValidPNG(t *testing.T) {
	if len(IconClosed) == 0 {
		t.Fatal("IconClosed is empty")
	}
	if !bytes.HasPrefix(IconClosed, pngHeader) {
		t.Fatal("IconClosed does not start with PNG header")
	}
}

func TestIconSizes(t *testing.T) {
	// Both icons should exist and be non-trivial in size
	t.Logf("IconOpen: %d bytes", len(IconOpen))
	t.Logf("IconClosed: %d bytes", len(IconClosed))

	if len(IconOpen) < 50 {
		t.Error("IconOpen too small to be a valid icon")
	}
	if len(IconClosed) < 50 {
		t.Error("IconClosed too small to be a valid icon")
	}
}

func TestIconsAreDifferent(t *testing.T) {
	// The open and closed icons must be different byte sequences
	if bytes.Equal(IconOpen, IconClosed) {
		t.Error("open and closed icons should be different")
	}
}
