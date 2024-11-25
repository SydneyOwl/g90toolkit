package tools

import (
	"bytes"
	"fmt"
	"github.com/sydneyowl/g90toolkit/firmware_data"
	"image"
	_ "image/png"
)

func PatchBootLogo(logo []byte, src []byte) error {
	// Search for original logo...
	index := bytes.Index(src, firmware_data.OriginalBootImage)
	if index == -1 {
		return fmt.Errorf("Original logo not found in the firmware. Have you already replaced it?")
	}
	// Check image
	img, _, err := image.Decode(bytes.NewReader(logo))
	if err != nil {
		return fmt.Errorf("Error decoding image: %v\n", err)
	}
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y
	if width != 48 || height != 48 {
		return fmt.Errorf("New bootlogo image width and height must be 48x48 (was %dx%d)\n", width, height)
	}
	fmt.Printf("Logo found at %d. Trying to replace it...\n", index)

	imgBytes := make([]byte, 0)
	for y := 0; y < height; y++ {
		for x := 0; x < width/8; x++ {
			b := byte(0)
			for i := 0; i < 8; i++ {
				px := img.At(x*8+i, y)
				// Convert pixel to black/white (1/0)
				gray, _, _, _ := px.RGBA()
				if gray > 0x7FFF {
					b |= 1 << (7 - i) // Set the corresponding bit to 1 for white
				}
			}
			imgBytes = append(imgBytes, b)
		}
	}
	var result []byte
	result = append(result, src[:index]...)
	result = append(result, imgBytes...)
	result = append(result, src[index+len(imgBytes):]...)
	copy(src, result)
	return nil
}

func PatchBootText(text []byte, src []byte) error {
	index := bytes.Index(src, firmware_data.OriginalBootText)
	if index == -1 {
		return fmt.Errorf("Original text not found in the firmware. Have you already replaced it?")
	}
	var result []byte
	result = append(result, src[:index]...)
	result = append(result, text...)
	result = append(result, src[index+len(text):]...)
	copy(src, result)
	return nil
}
