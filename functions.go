package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"strconv"
	"strings"
)

// convert2BPPToPNG converts a 2BPP tile to a PNG image
func convert2BPPToPNG(height int, tile []byte, img *image.RGBA, xPos int, yPos int) {
	for i := 0; i < width*height; i += 8 * bitDepth {
		highBit := 0
		lowBit := 0
		colorVal := 0

		for x := 0; x < 8; x++ {
			for y := 0; y < 8; y++ {
				if i+2*y+1 >= len(tile) {
					break
				}

				// Given that bit depth is 2:
				highBit = int(tile[i+2*y+1]>>(7-x)) & 0x01
				lowBit = int(tile[i+2*y]>>(7-x)) & 0x01
				value := (highBit << 1) | lowBit
				colorVal = int(255 * (float32(3-value) / 3))

				var c color.Color = color.RGBA{R: uint8(colorVal), G: uint8(colorVal), B: uint8(colorVal), A: 255}

				img.Set(x+(xPos*8), y+(yPos*8), c)
			}
		}

		xPos++
		if xPos >= (width / 8) {
			xPos = 0
			yPos++
		}
	}
}

// convertHexToInt32 converts a hex string to an int32
func convertHexToInt32(v string) int32 {
	startOffsetString := v
	// Remove the 0x prefix and keep only the actual number (as a string)
	if !strings.Contains(startOffsetString, "0x") {
		fmt.Println("Invalid start offset specified! Please specify a HEX value (e.g. 0x3f)")
		os.Exit(1)
	}
	// Remove the 0x prefix (it's not needed anymore)
	startOffsetString = strings.ReplaceAll(startOffsetString, "0x", "")
	// and convert the strings (which represent a hex value) to a decimal int32
	a, _ := strconv.ParseInt(startOffsetString, 16, 32)

	return int32(a)
}

func saveToDisk(outputFilename string, img image.Image) error {
	f, err := os.Create(outputFilename)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}

	defer func(f *os.File) {
		if err := f.Close(); err != nil {
			return
		}
	}(f)

	err = png.Encode(f, img)
	if err != nil {
		return fmt.Errorf("failed to encode image: %w", err)
	}

	return nil
}

// Remove duplicates [][]byte
func removeDuplicateByte(byteSlice [][]byte) [][]byte {
	allKeys := make(map[string]bool)

	var list [][]byte

	for _, item := range byteSlice {
		if _, value := allKeys[string(item)]; !value {
			allKeys[string(item)] = true

			list = append(list, item)
		}
	}

	return list
}

func getHexCodes(tiles []image.Image) [][]byte {
	origCodeTiles := make([][]byte, 0)
	for _, tile := range tiles {
		origCodeTiles = append(origCodeTiles, pngTo2BPP(tile))
	}

	return origCodeTiles
}

func compare(tile []byte, code []byte) bool {
	for i, b := range tile {
		if b != code[i] {
			return false
		}
	}

	return true
}

func areImagesEqual(img1, img2 image.Image) bool {
	bounds1 := img1.Bounds()
	bounds2 := img2.Bounds()

	if bounds1.Dx() != bounds2.Dx() || bounds1.Dy() != bounds2.Dy() {
		return false
	}

	for y := bounds1.Min.Y; y < bounds1.Max.Y; y++ {
		for x := bounds1.Min.X; x < bounds1.Max.X; x++ {
			if img1.At(x, y) != img2.At(x, y) {
				return false
			}
		}
	}

	return true
}

func flipImageHorizontally(img image.Image) image.Image {
	bounds := img.Bounds()
	flipped := image.NewRGBA(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			flipped.Set(bounds.Max.X-1-x, y, img.At(x, y))
		}
	}

	return flipped
}

func flipImageVertically(img image.Image) image.Image {
	bounds := img.Bounds()
	flipped := image.NewRGBA(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			flipped.Set(x, bounds.Max.Y-1-y, img.At(x, y))
		}
	}

	return flipped
}

func areImagesEquivalent(img1, img2 image.Image) bool {
	if areImagesEqual(img1, img2) {
		return true
	}

	flippedHorizontally := flipImageHorizontally(img1)
	if areImagesEqual(flippedHorizontally, img2) {
		return true
	}

	flippedVertically := flipImageVertically(img1)
	if areImagesEqual(flippedVertically, img2) {
		return true
	}

	flippedBoth := flipImageVertically(flippedHorizontally)
	return areImagesEqual(flippedBoth, img2)
}

func removeDuplicateString(slice []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range slice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}
