package main

import (
	"fmt"
	"image"
	"log"
	"os"
)

func readImageFromFilePath(path string) image.Image {
	// Load GB ROM
	infile, err := os.Open(path)
	if err != nil {
		// replace this with real error handling
		log.Fatal(err)
	}

	defer func(infile *os.File) {
		err := infile.Close()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}(infile)

	// Decode the image
	imData, imType, err := image.Decode(infile)
	if err != nil {
		fmt.Println(err)
	}

	// Make sure it's a PNG
	if imType != "png" {
		fmt.Println("Not a PNG")
	}

	return imData
}

func saveUniqueTiles(tiles []image.Image, tilesPerRow int, numRows int, filename string) ([]image.Image, error) {
	uniqueTiles := []image.Image{}

	for _, tile := range tiles {
		isUnique := true
		for j := 0; j < len(uniqueTiles); j++ {
			if areImagesEquivalent(tile, uniqueTiles[j]) {
				isUnique = false
				break
			}
		}

		if isUnique {
			uniqueTiles = append(uniqueTiles, tile)
		}
	}

	return uniqueTiles, nil
}
