package getimage

import (
	"bytes"
	"image"
	"os"
	"testing"
)

func TestGetImageHandler(t *testing.T) {

}

func TestResizeImage(t *testing.T) {
	// 1. Set up input
	originalImageBytes, err := os.ReadFile("img.jpg")
	width, height := 100, 50
	// 2. Call the function
	resizedImageBytes, err := resizeImage(originalImageBytes, width, height)
	// 3. Assert the results
	if err != nil {
		t.Errorf("Error resizing image: %v", err)
		return
	}
	// Check if the resized image has the correct dimensions
	img, _, err := image.Decode(bytes.NewReader(resizedImageBytes))
	if err != nil {
		t.Errorf("Error decoding resized image: %v", err)
		return
	}
	if img.Bounds().Dx() != width || img.Bounds().Dy() != height {
		t.Errorf("Resized image has incorrect dimensions. Expected (%d, %d), got (%d, %d)", width, height, img.Bounds().Dx(), img.Bounds().Dy())
	}
}
