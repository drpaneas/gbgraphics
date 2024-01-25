package main

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"io"
	"log"
	"os"
)

// compareFiles compares two files and returns true if they are equal
func compareFiles(file1, file2 string) bool {
	f1, err := os.ReadFile(file1)
	if err != nil {
		log.Fatal(err)
	}
	f2, err := os.ReadFile(file2)
	if err != nil {
		log.Fatal(err)
	}
	return bytes.Equal(f1, f2)
}

// copyFile copies a file from src to dst
func copyFile(src, dst string) error {

	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destination.Close()
	_, err = io.Copy(destination, source)
	return err
}

// createTiledArray For every tile, create an array with the size of the tiledImage
// and fill it with the tile index if the tile is found in the tiledImage
func createTiledArray(allTiles []image.Image, uniqueTiles []image.Image) []int {

	// Create an array if integers with the size of the tiledImageTiles
	tiledArray := make([]int, len(allTiles))

	for i, tile := range uniqueTiles {
		for j, imageTile := range allTiles {
			if areImagesEquivalent(tile, imageTile) {
				tiledArray[j] = i
			}
		}
	}

	return tiledArray

}

// createImageFromTiles creates a new image from the given tiles
func createImageFromTiles(tiles []image.Image, tilesPerRow, numRows int) image.Image {

	tileWidth, tileHeight := tiles[0].Bounds().Dx(), tiles[0].Bounds().Dy()
	imgWidth, imgHeight := tilesPerRow*(tileWidth+2), numRows*(tileHeight+2)

	output := image.NewRGBA(image.Rect(0, 0, imgWidth, imgHeight))

	// Assign a unique color for each unique tile
	colors := make([]color.Color, len(tiles))
	for i := range tiles {
		colors[i] = color.RGBA{uint8(i * 29 % 256), uint8(i * 73 % 256), uint8(i * 151 % 256), 255}

		for j := 0; j < i; j++ {
			if areImagesEquivalent(tiles[i], tiles[j]) {
				colors[i] = colors[j]
				break
			}
		}
	}

	for i, tile := range tiles {
		if i >= tilesPerRow*numRows {
			log.Println("Warning: More tiles provided than space in the output image.")
			break
		}

		x := (i % tilesPerRow) * (tileWidth + 2)
		y := (i / tilesPerRow) * (tileHeight + 2)

		r := image.Rect(x+1, y+1, x+1+tileWidth, y+1+tileHeight)
		draw.Draw(output, r, tile, image.Point{}, draw.Src)

		// Draw the perimeter with the assigned color
		perimeterColor := colors[i]
		for px := x; px <= x+tileWidth+1; px++ {
			output.Set(px, y, perimeterColor)              // Top border
			output.Set(px, y+tileHeight+1, perimeterColor) // Bottom border
		}
		for py := y; py <= y+tileHeight+1; py++ {
			output.Set(x, py, perimeterColor)             // Left border
			output.Set(x+tileWidth+1, py, perimeterColor) // Right border
		}
	}

	return output
}

// validateTileDimensions validates the dimensions of the tiles
func validateTileDimensions(tiles int, gbScreenXRes int, gbScreenYRes int) error {
	pixelsPerScreen := 23040 // the number of pixels in the screenshot (e.g. 160x144)
	pixelsPerTile := 64      // the number of pixels in a 8x8 tile, 8*8=64

	if tiles != pixelsPerScreen/pixelsPerTile { // we expect 360 tiles
		return fmt.Errorf("Not 23040/64 tiles")
	}

	if tilesPerRow != 20 {
		return fmt.Errorf("Not 20 tiles per row. It is: %d", tilesPerRow)
	}

	if tilesPerCol != 18 {
		return fmt.Errorf("Not 18 tiles per column. It is: %d", tilesPerCol)
	}

	return nil
}
