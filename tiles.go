package main

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"math"
	"os"
	"strings"
)

func getTiles(screenshot string, romBytes []byte, numColumns int) []string {
	// 1. Load a screenshot from the disk
	img := readImageFromFilePath(screenshot)

	// Step 2: Create a new image with the same dimensions as the original image but without the numColumns first columns
	newImg := image.NewRGBA(image.Rect(0, 0, img.Bounds().Max.X-numColumns, img.Bounds().Max.Y))

	// Step 3: Copy the pixels from the original image to the new image
	for y := 0; y < img.Bounds().Max.Y; y++ {
		for x := numColumns; x < img.Bounds().Max.X; x++ {
			newImg.Set(x-numColumns, y, img.At(x, y))
		}
	}

	// Print the dimensions of the new image
	// fmt.Println("New image dimensions:", newImg.Bounds().Max.X, newImg.Bounds().Max.Y)

	// Step 4: Split the image into sub-images
	tiles := split7x8(img) // Split the image into sub-images of 7x8 pixels (tiles) and store them in a slice
	if numColumns == 0 {
		tiles = split8x8(img) // Split the image into sub-images of 8x8 pixels (tiles) and store them in a slice
	} else if numColumns >= 1 && numColumns <= 7 {
		tiles = split7x8(img) // Split the image into sub-images of 7x8 pixels (tiles) and store them in a slice
	}

	// Step 5: Checks if the number of tiles is correct
	// err := validateTileDimensions(len(tiles), gbScreenXRes, gbScreenYRes)
	// if err != nil {
	// 	fmt.Println(err)
	// 	os.Exit(1)
	// }

	// Step 6: From the all the tiles, remove the duplicates
	uniqueTiles, err := saveUniqueTiles(tiles, gbScreenXRes, tilesPerCol, fmt.Sprintf("%d", numColumns+1))
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	// Step 7:  These tiles are in RGBA format, so we need to convert them to 2BPP
	// 			before we can compare them to the original gameboy tileset
	origCodeTiles := getHexCodes(uniqueTiles)
	uniqCodeTiles := removeDuplicateByte(origCodeTiles)

	// Step 8: Search for each tile in the screenshot and return the addresses from the ROM
	return findTileAddresses(uniqCodeTiles, romBytes)
}

// processTile Processes receives the addresses of tiles and converts them to PNG
func processTile(i int, v string, outputFilename string, romBytes []byte, rangeLength int, width int, bitDepth int) error {
	withoutPng := strings.ReplaceAll(outputFilename, ".png", "")
	newOutputFilename := fmt.Sprintf("%s_%d.png", withoutPng, i)

	rangeStartOffset := convertHexToInt32(v)

	if rangeStartOffset < 0 {
		return errors.New("invalid start offset specified")
	}

	rangeLengthInt32 := int32(rangeLength) // Convert rangeLength to int32

	tile := romBytes[rangeStartOffset : rangeStartOffset+rangeLengthInt32] // Use rangeLengthInt32

	// Calculate the height of the img in bytes
	height := 8 * int(math.Ceil(float64(len(tile))/float64(8*8*bitDepth)))
	hexValue := fmt.Sprintf("% X", tile)
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			img.Set(x, y, color.White)
		}
	}

	xPos, yPos := 0, 0

	// Modify the img
	convert2BPPToPNG(height, tile, img, xPos, yPos)

	if err := saveToDisk(newOutputFilename, img); err != nil {
		return err
	}

	fmt.Printf("'%s' (Found at location %s) converted to '%s'\n", hexValue, v, newOutputFilename)

	return nil
}

func split8x8(src image.Image) []image.Image {
	// Check if image resolution is 160x144
	if src.Bounds().Max.X != 160 || src.Bounds().Max.Y != 144 {
		fmt.Println("Not 160x144")
		fmt.Println("It is: ", src.Bounds().Max.X, src.Bounds().Max.Y)
		os.Exit(1)
	}

	checkColor(src)

	// Iterate over the image pixels and split it into 8x8 sub-images
	var tiles []image.Image

	for y := 0; y < src.Bounds().Max.Y; y += 8 {
		for x := 0; x < src.Bounds().Max.X; x += 8 {
			tile := image.NewRGBA(image.Rect(0, 0, 8, 8))

			for i := 0; i < 8; i++ {
				for j := 0; j < 8; j++ {
					tile.Set(i, j, src.At(x+i, y+j))
				}
			}

			tiles = append(tiles, tile)
		}
	}

	return tiles
}

func split7x8(src image.Image) []image.Image {
	// // Check if image resolution is 159x144
	// if src.Bounds().Max.X != 159 && src.Bounds().Max.X != 158 || src.Bounds().Max.Y != 144 {
	// 	fmt.Println("Not 159x144 or not 158x144")
	// 	fmt.Println("Error: It is: ", src.Bounds().Max.X, src.Bounds().Max.Y)
	// 	os.Exit(1)
	// }

	checkColor(src)

	// Iterate over the image pixels and split it into 7x8 sub-images
	var tiles []image.Image

	// Skip the last tile of every row
	for y := 0; y < src.Bounds().Max.Y; y += 8 {
		for x := 0; x < src.Bounds().Max.X-8; x += 8 {
			tile := image.NewRGBA(image.Rect(0, 0, 8, 8))

			for i := 0; i < 8; i++ {
				for j := 0; j < 8; j++ {
					tile.Set(i, j, src.At(x+i, y+j))
				}
			}

			tiles = append(tiles, tile)
		}
	}

	return tiles
}

// Function to find the addresses of tiles in the ROM
func findTileAddresses(uniqCodeTiles [][]byte, romBytes []byte) []string {
	var addr []string

	// Search for each tile in the screenshot
	for _, tile := range uniqCodeTiles {
		// Calculate the maximum index in romBytes where a match with tile can start
		// This is to prevent an "out of bounds" error when accessing the slice romBytes[j:j+len(tile)]
		maxIndex := len(romBytes) - len(tile) + 1

		// Search for the tile in the rom
		for j := 0; j < maxIndex; j++ {
			// Get a slice of romBytes that has the same length as tile
			// This slice will be compared with tile to check for a match
			romSlice := romBytes[j : j+len(tile)]

			if compare(tile, romSlice) {
				// If a match is found, convert the index of the match to a hexadecimal string
				// and append it to the addr slice
				addr = append(addr, fmt.Sprintf("0x%X", j))

				// No reason to keep searching for this tile if you have already found it
				break
			}
		}
	}

	return addr
}
