package main

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io"
	"log"
	"math"
	"os"
	"runtime/debug"
	"strconv"
	"strings"

	"github.com/alexflint/go-arg"
)

var Commit = func() string {
	if info, ok := debug.ReadBuildInfo(); ok {
		for _, setting := range info.Settings {
			if setting.Key == "vcs.revision" {
				return setting.Value
			}
		}
	}

	return ""
}()

type args struct {
	Rom        string `arg:"positional,required" help:"Path to the ROM file"`
	Screenshot string `arg:"required,--img" help:"path of in-game screenshot" placeholder:"<SCREENSHOT>"`
	Output     string `arg:"--output" help:"output file" default:"out.png" placeholder:"<FILE>"`
}

func (args) Description() string {
	return "GBGraphics - extract graphics from Gameboy ROM using a screenshot"
}

func (args) Version() string {
	return "git commit " + Commit
}

var rangeStartOffset int32

const (
	width       = 8
	bitDepth    = 2
	rangeLength = 16
)

func main() {
	var userInput args

	arg.MustParse(&userInput)

	outputFilename := userInput.Output
	screenshot := userInput.Screenshot

	path := userInput.Rom
	if path == "" {
		fmt.Println("No ROM specified!")
		os.Exit(1)
	}

	romBytes, errReadFile := os.ReadFile(path)
	if errReadFile != nil {
		fmt.Println(errReadFile)
		os.Exit(1)
	}

	locations := getTiles(screenshot, romBytes)
	for i, v := range locations {
		withoutPng := strings.ReplaceAll(outputFilename, ".png", "")
		newOutputFilename := fmt.Sprintf("%s_%d.png", withoutPng, i)

		rangeStartOffset = convertHexToInt32(v)

		if rangeStartOffset < 0 {
			fmt.Println("Invalid start offset specified!")
			os.Exit(1)
		}

		tile := romBytes[rangeStartOffset : rangeStartOffset+rangeLength]

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
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Printf("'%s' (Found at location %s) converted to '%s'\n", hexValue, v, newOutputFilename)
	}

	locations2 := getTiles2(screenshot, romBytes)

	for i, v := range locations2 {
		withoutPng := strings.ReplaceAll(outputFilename, ".png", "")
		newOutputFilename := fmt.Sprintf("second_%s_%d.png", withoutPng, i)

		rangeStartOffset = convertHexToInt32(v)

		if rangeStartOffset < 0 {
			fmt.Println("Invalid start offset specified!")
			os.Exit(1)
		}

		tile := romBytes[rangeStartOffset : rangeStartOffset+rangeLength]

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
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Printf("'%s' (Found at location %s) converted to '%s'\n", hexValue, v, newOutputFilename)
	}

	locations3 := getTiles3(screenshot, romBytes)

	for i, v := range locations3 {
		withoutPng := strings.ReplaceAll(outputFilename, ".png", "")
		newOutputFilename := fmt.Sprintf("third_%s_%d.png", withoutPng, i)

		rangeStartOffset = convertHexToInt32(v)

		if rangeStartOffset < 0 {
			fmt.Println("Invalid start offset specified!")
			os.Exit(1)
		}

		tile := romBytes[rangeStartOffset : rangeStartOffset+rangeLength]

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
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Printf("'%s' (Found at location %s) converted to '%s'\n", hexValue, v, newOutputFilename)
	}

	locations4 := getTiles4(screenshot, romBytes)

	for i, v := range locations4 {
		withoutPng := strings.ReplaceAll(outputFilename, ".png", "")
		newOutputFilename := fmt.Sprintf("forth_%s_%d.png", withoutPng, i)

		rangeStartOffset = convertHexToInt32(v)

		if rangeStartOffset < 0 {
			fmt.Println("Invalid start offset specified!")
			os.Exit(1)
		}

		tile := romBytes[rangeStartOffset : rangeStartOffset+rangeLength]

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
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Printf("'%s' (Found at location %s) converted to '%s'\n", hexValue, v, newOutputFilename)
	}

	locations5 := getTiles5(screenshot, romBytes)

	for i, v := range locations5 {
		withoutPng := strings.ReplaceAll(outputFilename, ".png", "")
		newOutputFilename := fmt.Sprintf("fifth_%s_%d.png", withoutPng, i)

		rangeStartOffset = convertHexToInt32(v)

		if rangeStartOffset < 0 {
			fmt.Println("Invalid start offset specified!")
			os.Exit(1)
		}

		tile := romBytes[rangeStartOffset : rangeStartOffset+rangeLength]

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
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Printf("'%s' (Found at location %s) converted to '%s'\n", hexValue, v, newOutputFilename)
	}

	locations6 := getTiles6(screenshot, romBytes)

	for i, v := range locations6 {
		withoutPng := strings.ReplaceAll(outputFilename, ".png", "")
		newOutputFilename := fmt.Sprintf("sixth_%s_%d.png", withoutPng, i)

		rangeStartOffset = convertHexToInt32(v)

		if rangeStartOffset < 0 {
			fmt.Println("Invalid start offset specified!")
			os.Exit(1)
		}

		tile := romBytes[rangeStartOffset : rangeStartOffset+rangeLength]

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
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Printf("'%s' (Found at location %s) converted to '%s'\n", hexValue, v, newOutputFilename)
	}

	locations7 := getTiles7(screenshot, romBytes)

	for i, v := range locations7 {
		withoutPng := strings.ReplaceAll(outputFilename, ".png", "")
		newOutputFilename := fmt.Sprintf("seventh_%s_%d.png", withoutPng, i)

		rangeStartOffset = convertHexToInt32(v)

		if rangeStartOffset < 0 {
			fmt.Println("Invalid start offset specified!")
			os.Exit(1)
		}

		tile := romBytes[rangeStartOffset : rangeStartOffset+rangeLength]

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
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Printf("'%s' (Found at location %s) converted to '%s'\n", hexValue, v, newOutputFilename)
	}

	locations8 := getTiles8(screenshot, romBytes)

	for i, v := range locations8 {
		withoutPng := strings.ReplaceAll(outputFilename, ".png", "")
		newOutputFilename := fmt.Sprintf("eighth_%s_%d.png", withoutPng, i)

		rangeStartOffset = convertHexToInt32(v)

		if rangeStartOffset < 0 {
			fmt.Println("Invalid start offset specified!")
			os.Exit(1)
		}

		tile := romBytes[rangeStartOffset : rangeStartOffset+rangeLength]

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
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Printf("'%s' (Found at location %s) converted to '%s'\n", hexValue, v, newOutputFilename)
	}

	// Compare all the images of the disk that start with:
	// 'out_*.png' and 'second_out_*.png' and third_out_*.png and forth_out_*.png and fifth_out_*.png and sixth_out_*.png and seventh_out_*.png and eighth_out_*.png
	// and delete the duplicates

	// Get all the files in the current directory
	files, err := os.ReadDir(".")
	if err != nil {
		log.Fatal(err)
	}

	// Create a slice of all the filenames
	var filenames []string
	for _, file := range files {
		filenames = append(filenames, file.Name())
	}

	// move all the files that start with 'out_' to a new slice
	var outFiles []string
	for _, file := range filenames {

		if strings.HasPrefix(file, "out_") {
			outFiles = append(outFiles, file)
		}
	}

	// move all the files that start with 'second_out_' to a new slice
	var secondOutFiles []string
	for _, file := range filenames {

		if strings.HasPrefix(file, "second_out_") {
			secondOutFiles = append(secondOutFiles, file)
		}
	}

	// move all the files that start with 'third_out_' to a new slice
	var thirdOutFiles []string
	for _, file := range filenames {

		if strings.HasPrefix(file, "third_out_") {
			thirdOutFiles = append(thirdOutFiles, file)
		}
	}

	// move all the files that start with 'forth_out_' to a new slice
	var forthOutFiles []string
	for _, file := range filenames {

		if strings.HasPrefix(file, "forth_out_") {
			forthOutFiles = append(forthOutFiles, file)
		}
	}

	// move all the files that start with 'fifth_out_' to a new slice
	var fifthOutFiles []string
	for _, file := range filenames {

		if strings.HasPrefix(file, "fifth_out_") {
			fifthOutFiles = append(fifthOutFiles, file)
		}
	}

	// move all the files that start with 'sixth_out_' to a new slice
	var sixthOutFiles []string
	for _, file := range filenames {

		if strings.HasPrefix(file, "sixth_out_") {
			sixthOutFiles = append(sixthOutFiles, file)
		}
	}

	// move all the files that start with 'seventh_out_' to a new slice
	var seventhOutFiles []string
	for _, file := range filenames {

		if strings.HasPrefix(file, "seventh_out_") {
			seventhOutFiles = append(seventhOutFiles, file)
		}
	}

	// move all the files that start with 'eighth_out_' to a new slice
	var eighthOutFiles []string
	for _, file := range filenames {

		if strings.HasPrefix(file, "eighth_out_") {
			eighthOutFiles = append(eighthOutFiles, file)
		}
	}

	// Create one slice with all the elements from all of the slices
	var allFiles []string
	allFiles = append(allFiles, outFiles...)
	allFiles = append(allFiles, secondOutFiles...)
	allFiles = append(allFiles, thirdOutFiles...)
	allFiles = append(allFiles, forthOutFiles...)
	allFiles = append(allFiles, fifthOutFiles...)
	allFiles = append(allFiles, sixthOutFiles...)
	allFiles = append(allFiles, seventhOutFiles...)
	allFiles = append(allFiles, eighthOutFiles...)

	// Create a folder 'final' if it doesn't exist
	// if it exists, delete it and create it again
	if _, err := os.Stat("final"); os.IsNotExist(err) {
		os.Mkdir("final", 0755)
	} else {
		os.RemoveAll("final")
		os.Mkdir("final", 0755)
	}

	// Move all the files to the 'final' folder
	for _, file := range allFiles {
		// Copy the file to the 'final' folder
		copyFile(file, "final/"+file)

		// Remove the file from the current directory
		os.Remove(file)
	}

	// Compare all the files in the 'final' folder and delete the duplicates
	// Get all the files in the current directory
	files, err = os.ReadDir("final")
	if err != nil {
		log.Fatal(err)
	}

	// create a slice for duplicate files
	var duplicateFiles []string

	// Compare all 'files' with each other (binary comparison) and delete the duplicates
	for _, file := range files {
		for _, file2 := range files {
			if file.Name() != file2.Name() {
				if compareFiles("final/"+file.Name(), "final/"+file2.Name()) {
					// Add the file to the slice of duplicate files
					duplicateFiles = append(duplicateFiles, file2.Name())
				}
			}
		}
	}

	// Delete all the files in the 'duplicateFiles' slice
	for _, file := range duplicateFiles {
		_ = os.Remove("final/" + file)
	}
}

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

const (
	// PaletteBGB is the default palette used by BGB

	// PaletteGreyscale is the default greyscale gameboy colour palette.
	// PaletteGreyscale = byte(iota)
	// PaletteOriginal is more authentic looking green tinted gameboy
	// colour palette  as it would have been on the GameBoy
	// PaletteOriginal

	PaletteBGB = 2
)

const (
	lightest = byte(iota)
	light
	dark
	darkest
)

// Palettes is a mapping from colour palettes to their colour values
// to be used by the emulator.
var Palettes = [][][]byte{
	// PaletteGreyscale
	{
		{0xFF, 0xFF, 0xFF},
		{0xCC, 0xCC, 0xCC},
		{0x77, 0x77, 0x77},
		{0x00, 0x00, 0x00},
	},
	// PaletteOriginal
	{
		{0x9B, 0xBC, 0x0F},
		{0x8B, 0xAC, 0x0F},
		{0x30, 0x62, 0x30},
		{0x0F, 0x38, 0x0F},
	},
	// PaletteBGB
	{
		{0xE0, 0xF8, 0xD0}, // lightest (#E0F8D0)
		{0x88, 0xC0, 0x70}, // light (#88C070)
		{0x34, 0x68, 0x56}, // dark (#346856)
		{0x08, 0x18, 0x20}, // darkest (#081820)
	},
}

// GetPaletteColour returns the colour based on the colour index and the currently
// selected palette.
func GetPaletteColour(index byte, palette byte) (uint8, uint8, uint8) {
	col := Palettes[palette][index]
	r, g, b := col[0], col[1], col[2]

	return r, g, b
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

func getTiles2(screenshot string, romBytes []byte) []string {
	// Step 1: Load the screenshot from the disk
	img := readImageFromFilePath(screenshot)

	// Print the dimensions of the image
	fmt.Println("Image dimensions:", img.Bounds().Max.X, img.Bounds().Max.Y)

	// Step 2: Create a new image with the same dimensions as the original image
	// but without the first column
	newImg := image.NewRGBA(image.Rect(0, 0, img.Bounds().Max.X-1, img.Bounds().Max.Y))

	// Step 3: Copy the pixels from the original image to the new image
	for y := 0; y < img.Bounds().Max.Y; y++ {
		for x := 1; x < img.Bounds().Max.X; x++ {
			newImg.Set(x-1, y, img.At(x, y))
		}
	}

	// Print the dimensions of the new image
	fmt.Println("New image dimensions:", newImg.Bounds().Max.X, newImg.Bounds().Max.Y)

	// Step 4: Save the new image to disk
	saveToDisk("newImg2.png", newImg)

	// ---------------------------------------------

	// Step 1: Split the image into 7x8 tiles (7x8 because we removed the first column)
	tiles := split7x8(newImg)

	// Step 2: Print the number of tiles
	fmt.Println("Number of tiles new image:", len(tiles))

	// Step 3: Create a new image from the tiles
	tiledImage := createImageFromTiles(tiles, 19, 18)

	// Step 4: Save the tiled image to disk
	saveToDisk("tiledImage2.png", tiledImage)

	// Step 5: Save unique tiles to disk and the final image
	uniqueTiles, err := saveUniqueTiles(tiles, 19, 18, "2")
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	// Step 6: Create TiledArray
	tiledArray := createTiledArray(tiles, uniqueTiles)

	// Step 7: Print the tiledArray to the console in readable format, like 19x18
	for i, v := range tiledArray {
		if i%19 == 0 {
			fmt.Println()
		}
		fmt.Printf("%02d ", v)
	}
	fmt.Println()

	// Step 8: Change this line later
	tiles = uniqueTiles

	// Step 9: These tiles are in RGBA format, so we need to convert them to 2BPP
	// before we can compare them to the original gameboy tileset
	origCodeTiles := getHexCodes(tiles)
	uniqCodeTiles := removeDuplicateByte(origCodeTiles)

	var addr []string

	// Step 10: Search for each tile in the screenshot
	for i, tile := range uniqCodeTiles {
		// Search for the tile in the rom
		for j := 0; j < len(romBytes); j++ {
			if j+16 > len(romBytes) { // check for out-of-bounds
				break
			}

			if compare(tile, romBytes[j:j+16]) {
				tmp := fmt.Sprintf("0x%X", j)
				addr = append(addr, tmp)
				_ = i

				break // no reason to keep searching for this tile if you have already found it
			}
		}
	}

	return addr
}

func getTiles3(screenshot string, romBytes []byte) []string {
	// Step 1: Load the screenshot from the disk
	img := readImageFromFilePath(screenshot)

	// Print the dimensions of the image
	fmt.Println("Image dimensions:", img.Bounds().Max.X, img.Bounds().Max.Y)

	// Step 2: Create a new image with the same dimensions as the original image
	// but without the first 2 columns
	newImg := image.NewRGBA(image.Rect(0, 0, img.Bounds().Max.X-2, img.Bounds().Max.Y))

	// Step 3: Copy the pixels from the original image to the new image (x-2)
	for y := 0; y < img.Bounds().Max.Y; y++ {
		for x := 2; x < img.Bounds().Max.X; x++ {
			newImg.Set(x-2, y, img.At(x, y))
		}
	}

	// Print the dimensions of the new image
	fmt.Println("New image dimensions:", newImg.Bounds().Max.X, newImg.Bounds().Max.Y)

	// Step 4: Save the new image to disk
	saveToDisk("new3Img.png", newImg)

	// ---------------------------------------------

	// Step 1: Split the image into 7x8 tiles (7x8 because we removed the first column)
	tiles := split7x8(newImg)

	// Step 2: Print the number of tiles
	fmt.Println("Number of tiles new image:", len(tiles))

	// Step 3: Create a new image from the tiles
	tiledImage := createImageFromTiles(tiles, 19, 18)

	// Step 4: Save the tiled image to disk
	saveToDisk("tiledImage3.png", tiledImage)

	// Step 5: Save unique tiles to disk and the final image
	uniqueTiles, err := saveUniqueTiles(tiles, 19, 18, "3")
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	// Step 6: Create TiledArray
	tiledArray := createTiledArray(tiles, uniqueTiles)

	// Step 7: Print the tiledArray to the console in readable format, like 19x18
	for i, v := range tiledArray {
		if i%19 == 0 {
			fmt.Println()
		}
		fmt.Printf("%02d ", v)
	}
	fmt.Println()

	// Step 8: Change this line later
	tiles = uniqueTiles

	// Step 9: These tiles are in RGBA format, so we need to convert them to 2BPP
	// before we can compare them to the original gameboy tileset
	origCodeTiles := getHexCodes(tiles)
	uniqCodeTiles := removeDuplicateByte(origCodeTiles)

	var addr []string

	// Step 10: Search for each tile in the screenshot
	for i, tile := range uniqCodeTiles {
		// Search for the tile in the rom
		for j := 0; j < len(romBytes); j++ {
			if j+16 > len(romBytes) { // check for out-of-bounds
				break
			}

			if compare(tile, romBytes[j:j+16]) {
				tmp := fmt.Sprintf("0x%X", j)
				addr = append(addr, tmp)
				_ = i

				break // no reason to keep searching for this tile if you have already found it
			}
		}
	}

	return addr
}

func getTiles4(screenshot string, romBytes []byte) []string {
	// Step 1: Load the screenshot from the disk
	img := readImageFromFilePath(screenshot)

	// Print the dimensions of the image
	fmt.Println("Image dimensions:", img.Bounds().Max.X, img.Bounds().Max.Y)

	// Step 2: Create a new image with the same dimensions as the original image
	// but without the first 3 columns
	newImg := image.NewRGBA(image.Rect(0, 0, img.Bounds().Max.X-3, img.Bounds().Max.Y))

	// Step 3: Copy the pixels from the original image to the new image
	for y := 0; y < img.Bounds().Max.Y; y++ {
		for x := 3; x < img.Bounds().Max.X; x++ {
			newImg.Set(x-3, y, img.At(x, y))
		}
	}

	// Print the dimensions of the new image
	fmt.Println("New image dimensions:", newImg.Bounds().Max.X, newImg.Bounds().Max.Y)

	// Step 4: Save the new image to disk
	saveToDisk("new4Img.png", newImg)

	// ---------------------------------------------

	// Step 1: Split the image into 7x8 tiles (7x8 because we removed the first column)
	tiles := split7x8(newImg)

	// Step 2: Print the number of tiles
	fmt.Println("Number of tiles new image:", len(tiles))

	// Step 3: Create a new image from the tiles
	tiledImage := createImageFromTiles(tiles, 19, 18)

	// Step 4: Save the tiled image to disk
	saveToDisk("tiledImage4.png", tiledImage)

	// Step 5: Save unique tiles to disk and the final image
	uniqueTiles, err := saveUniqueTiles(tiles, 19, 18, "4")
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	// Step 6: Create TiledArray
	tiledArray := createTiledArray(tiles, uniqueTiles)

	// Step 7: Print the tiledArray to the console in readable format, like 19x18
	for i, v := range tiledArray {
		if i%19 == 0 {
			fmt.Println()
		}
		fmt.Printf("%02d ", v)
	}
	fmt.Println()

	// Step 8: Change this line later
	tiles = uniqueTiles

	// Step 9: These tiles are in RGBA format, so we need to convert them to 2BPP
	// before we can compare them to the original gameboy tileset
	origCodeTiles := getHexCodes(tiles)
	uniqCodeTiles := removeDuplicateByte(origCodeTiles)

	var addr []string

	// Step 10: Search for each tile in the screenshot
	for i, tile := range uniqCodeTiles {
		// Search for the tile in the rom
		for j := 0; j < len(romBytes); j++ {
			if j+16 > len(romBytes) { // check for out-of-bounds
				break
			}

			if compare(tile, romBytes[j:j+16]) {
				tmp := fmt.Sprintf("0x%X", j)
				addr = append(addr, tmp)
				_ = i

				break // no reason to keep searching for this tile if you have already found it
			}
		}
	}

	return addr
}

func getTiles5(screenshot string, romBytes []byte) []string {
	// Step 1: Load the screenshot from the disk
	img := readImageFromFilePath(screenshot)

	// Print the dimensions of the image
	fmt.Println("Image dimensions:", img.Bounds().Max.X, img.Bounds().Max.Y)

	// Step 2: Create a new image with the same dimensions as the original image
	// but without the first 3 columns
	newImg := image.NewRGBA(image.Rect(0, 0, img.Bounds().Max.X-4, img.Bounds().Max.Y))

	// Step 3: Copy the pixels from the original image to the new image
	for y := 0; y < img.Bounds().Max.Y; y++ {
		for x := 4; x < img.Bounds().Max.X; x++ {
			newImg.Set(x-4, y, img.At(x, y))
		}
	}

	// Print the dimensions of the new image
	fmt.Println("New image dimensions:", newImg.Bounds().Max.X, newImg.Bounds().Max.Y)

	// Step 4: Save the new image to disk
	saveToDisk("new5Img.png", newImg)

	// ---------------------------------------------

	// Step 1: Split the image into 7x8 tiles (7x8 because we removed the first column)
	tiles := split7x8(newImg)

	// Step 2: Print the number of tiles
	fmt.Println("Number of tiles new image:", len(tiles))

	// Step 3: Create a new image from the tiles
	tiledImage := createImageFromTiles(tiles, 19, 18)

	// Step 4: Save the tiled image to disk
	saveToDisk("tiledImage5.png", tiledImage)

	// Step 5: Save unique tiles to disk and the final image
	uniqueTiles, err := saveUniqueTiles(tiles, 19, 18, "5")
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	// Step 6: Create TiledArray
	tiledArray := createTiledArray(tiles, uniqueTiles)

	// Step 7: Print the tiledArray to the console in readable format, like 19x18
	for i, v := range tiledArray {
		if i%19 == 0 {
			fmt.Println()
		}
		fmt.Printf("%02d ", v)
	}
	fmt.Println()

	// Step 8: Change this line later
	tiles = uniqueTiles

	// Step 9: These tiles are in RGBA format, so we need to convert them to 2BPP
	// before we can compare them to the original gameboy tileset
	origCodeTiles := getHexCodes(tiles)
	uniqCodeTiles := removeDuplicateByte(origCodeTiles)

	var addr []string

	// Step 10: Search for each tile in the screenshot
	for i, tile := range uniqCodeTiles {
		// Search for the tile in the rom
		for j := 0; j < len(romBytes); j++ {
			if j+16 > len(romBytes) { // check for out-of-bounds
				break
			}

			if compare(tile, romBytes[j:j+16]) {
				tmp := fmt.Sprintf("0x%X", j)
				addr = append(addr, tmp)
				_ = i

				break // no reason to keep searching for this tile if you have already found it
			}
		}
	}

	return addr
}

func getTiles6(screenshot string, romBytes []byte) []string {
	// Step 1: Load the screenshot from the disk
	img := readImageFromFilePath(screenshot)

	// Print the dimensions of the image
	fmt.Println("Image dimensions:", img.Bounds().Max.X, img.Bounds().Max.Y)

	// Step 2: Create a new image with the same dimensions as the original image
	// but without the first 3 columns
	newImg := image.NewRGBA(image.Rect(0, 0, img.Bounds().Max.X-5, img.Bounds().Max.Y))

	// Step 3: Copy the pixels from the original image to the new image
	for y := 0; y < img.Bounds().Max.Y; y++ {
		for x := 5; x < img.Bounds().Max.X; x++ {
			newImg.Set(x-5, y, img.At(x, y))
		}
	}

	// Print the dimensions of the new image
	fmt.Println("New image dimensions:", newImg.Bounds().Max.X, newImg.Bounds().Max.Y)

	// Step 4: Save the new image to disk
	saveToDisk("new6Img.png", newImg)

	// ---------------------------------------------

	// Step 1: Split the image into 7x8 tiles (7x8 because we removed the first column)
	tiles := split7x8(newImg)

	// Step 2: Print the number of tiles
	fmt.Println("Number of tiles new image:", len(tiles))

	// Step 3: Create a new image from the tiles
	tiledImage := createImageFromTiles(tiles, 19, 18)

	// Step 4: Save the tiled image to disk
	saveToDisk("tiledImage6.png", tiledImage)

	// Step 5: Save unique tiles to disk and the final image
	uniqueTiles, err := saveUniqueTiles(tiles, 19, 18, "6")
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	// Step 6: Create TiledArray
	tiledArray := createTiledArray(tiles, uniqueTiles)

	// Step 7: Print the tiledArray to the console in readable format, like 19x18
	for i, v := range tiledArray {
		if i%19 == 0 {
			fmt.Println()
		}
		fmt.Printf("%02d ", v)
	}
	fmt.Println()

	// Step 8: Change this line later
	tiles = uniqueTiles

	// Step 9: These tiles are in RGBA format, so we need to convert them to 2BPP
	// before we can compare them to the original gameboy tileset
	origCodeTiles := getHexCodes(tiles)
	uniqCodeTiles := removeDuplicateByte(origCodeTiles)

	var addr []string

	// Step 10: Search for each tile in the screenshot
	for i, tile := range uniqCodeTiles {
		// Search for the tile in the rom
		for j := 0; j < len(romBytes); j++ {
			if j+16 > len(romBytes) { // check for out-of-bounds
				break
			}

			if compare(tile, romBytes[j:j+16]) {
				tmp := fmt.Sprintf("0x%X", j)
				addr = append(addr, tmp)
				_ = i

				break // no reason to keep searching for this tile if you have already found it
			}
		}
	}

	return addr
}

func getTiles7(screenshot string, romBytes []byte) []string {
	// Step 1: Load the screenshot from the disk
	img := readImageFromFilePath(screenshot)

	// Print the dimensions of the image
	fmt.Println("Image dimensions:", img.Bounds().Max.X, img.Bounds().Max.Y)

	// Step 2: Create a new image with the same dimensions as the original image
	// but without the first 3 columns
	newImg := image.NewRGBA(image.Rect(0, 0, img.Bounds().Max.X-6, img.Bounds().Max.Y))

	// Step 3: Copy the pixels from the original image to the new image
	for y := 0; y < img.Bounds().Max.Y; y++ {
		for x := 6; x < img.Bounds().Max.X; x++ {
			newImg.Set(x-6, y, img.At(x, y))
		}
	}

	// Print the dimensions of the new image
	fmt.Println("New image dimensions:", newImg.Bounds().Max.X, newImg.Bounds().Max.Y)

	// Step 4: Save the new image to disk
	saveToDisk("new7Img.png", newImg)

	// ---------------------------------------------

	// Step 1: Split the image into 7x8 tiles (7x8 because we removed the first column)
	tiles := split7x8(newImg)

	// Step 2: Print the number of tiles
	fmt.Println("Number of tiles new image:", len(tiles))

	// Step 3: Create a new image from the tiles
	tiledImage := createImageFromTiles(tiles, 19, 18)

	// Step 4: Save the tiled image to disk
	saveToDisk("tiledImage7.png", tiledImage)

	// Step 5: Save unique tiles to disk and the final image
	uniqueTiles, err := saveUniqueTiles(tiles, 19, 18, "7")
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	// Step 6: Create TiledArray
	tiledArray := createTiledArray(tiles, uniqueTiles)

	// Step 7: Print the tiledArray to the console in readable format, like 19x18
	for i, v := range tiledArray {
		if i%19 == 0 {
			fmt.Println()
		}
		fmt.Printf("%02d ", v)
	}
	fmt.Println()

	// Step 8: Change this line later
	tiles = uniqueTiles

	// Step 9: These tiles are in RGBA format, so we need to convert them to 2BPP
	// before we can compare them to the original gameboy tileset
	origCodeTiles := getHexCodes(tiles)
	uniqCodeTiles := removeDuplicateByte(origCodeTiles)

	var addr []string

	// Step 10: Search for each tile in the screenshot
	for i, tile := range uniqCodeTiles {
		// Search for the tile in the rom
		for j := 0; j < len(romBytes); j++ {
			if j+16 > len(romBytes) { // check for out-of-bounds
				break
			}

			if compare(tile, romBytes[j:j+16]) {
				tmp := fmt.Sprintf("0x%X", j)
				addr = append(addr, tmp)
				_ = i

				break // no reason to keep searching for this tile if you have already found it
			}
		}
	}

	return addr
}

func getTiles8(screenshot string, romBytes []byte) []string {
	// Step 1: Load the screenshot from the disk
	img := readImageFromFilePath(screenshot)

	// Print the dimensions of the image
	fmt.Println("Image dimensions:", img.Bounds().Max.X, img.Bounds().Max.Y)

	// Step 2: Create a new image with the same dimensions as the original image
	// but without the first 3 columns
	newImg := image.NewRGBA(image.Rect(0, 0, img.Bounds().Max.X-7, img.Bounds().Max.Y))

	// Step 3: Copy the pixels from the original image to the new image
	for y := 0; y < img.Bounds().Max.Y; y++ {
		for x := 7; x < img.Bounds().Max.X; x++ {
			newImg.Set(x-7, y, img.At(x, y))
		}
	}

	// Print the dimensions of the new image
	fmt.Println("New image dimensions:", newImg.Bounds().Max.X, newImg.Bounds().Max.Y)

	// Step 4: Save the new image to disk
	saveToDisk("new8Img.png", newImg)

	// ---------------------------------------------

	// Step 1: Split the image into 7x8 tiles (7x8 because we removed the first column)
	tiles := split7x8(newImg)

	// Step 2: Print the number of tiles
	fmt.Println("Number of tiles new image:", len(tiles))

	// Step 3: Create a new image from the tiles
	tiledImage := createImageFromTiles(tiles, 19, 18)

	// Step 4: Save the tiled image to disk
	saveToDisk("tiledImage8.png", tiledImage)

	// Step 5: Save unique tiles to disk and the final image
	uniqueTiles, err := saveUniqueTiles(tiles, 19, 18, "8")
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	// Step 6: Create TiledArray
	tiledArray := createTiledArray(tiles, uniqueTiles)

	// Step 7: Print the tiledArray to the console in readable format, like 19x18
	for i, v := range tiledArray {
		if i%19 == 0 {
			fmt.Println()
		}
		fmt.Printf("%02d ", v)
	}
	fmt.Println()

	// Step 8: Change this line later
	tiles = uniqueTiles

	// Step 9: These tiles are in RGBA format, so we need to convert them to 2BPP
	// before we can compare them to the original gameboy tileset
	origCodeTiles := getHexCodes(tiles)
	uniqCodeTiles := removeDuplicateByte(origCodeTiles)

	var addr []string

	// Step 10: Search for each tile in the screenshot
	for i, tile := range uniqCodeTiles {
		// Search for the tile in the rom
		for j := 0; j < len(romBytes); j++ {
			if j+16 > len(romBytes) { // check for out-of-bounds
				break
			}

			if compare(tile, romBytes[j:j+16]) {
				tmp := fmt.Sprintf("0x%X", j)
				addr = append(addr, tmp)
				_ = i

				break // no reason to keep searching for this tile if you have already found it
			}
		}
	}

	return addr
}

func getTiles(screenshot string, romBytes []byte) []string {
	// Load a screenshot from the emulator and split it into 8x8 tiles
	img := readImageFromFilePath(screenshot)
	tiles := split8x8(img)

	if len(tiles) != 23040/64 { // 23040 is the number of pixels in the screenshot (e.g. 160x144)
		fmt.Println("Not 23040/64 tiles") // 64 is the number of pixels in a 8x8 tile, 8*8=64
		os.Exit(1)                        // so we expect 360 tiles
	}

	tiledImage := createImageFromTiles(tiles, 20, 18)

	// Save the tiled image to disk with "tiledImage.png"
	saveToDisk("tiledImage.png", tiledImage)

	// Save unique tiles to disk and the final image
	uniqueTiles, err := saveUniqueTiles(tiles, 20, 18, "1")
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	// Create TiledArray
	tiledArray := createTiledArray(tiles, uniqueTiles)

	// Print the tiledArray to the console in readable format, like 20x18
	for i, v := range tiledArray {
		if i%20 == 0 {
			fmt.Println()
		}
		fmt.Printf("%02d ", v)
	}
	fmt.Println()

	// Change this line later
	tiles = uniqueTiles

	// These tiles are in RGBA format, so we need to convert them to 2BPP
	// before we can compare them to the original gameboy tileset
	origCodeTiles := getHexCodes(tiles)
	uniqCodeTiles := removeDuplicateByte(origCodeTiles)

	var addr []string

	// Search for each tile in the screenshot
	for i, tile := range uniqCodeTiles {
		// Search for the tile in the rom
		for j := 0; j < len(romBytes); j++ {
			if j+16 > len(romBytes) { // check for out-of-bounds
				break
			}

			if compare(tile, romBytes[j:j+16]) {
				tmp := fmt.Sprintf("0x%X", j)
				addr = append(addr, tmp)
				_ = i

				break // no reason to keep searching for this tile if you have already found it
			}
		}
	}

	return addr
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

// Takes a 8x8 PNG RGBA images and converts it to 2BPP
// and returns the original byte array (what the GB rom would contain, and you could see with a hex editor)
// And also it saves the image as dest.2bpp
func pngTo2BPP(imData image.Image) []byte {
	// Make sure it is 8x8
	if imData.Bounds().Max.X != 8 || imData.Bounds().Max.Y != 8 {
		fmt.Println("Not 8x8")
		fmt.Println("It is: ", imData.Bounds().Max.X, imData.Bounds().Max.Y)
		os.Exit(1)
	}

	checkColor(imData)

	// print image pixel color r,g,b,a into uint32
	var binCode []byte

	for y := 0; y < imData.Bounds().Max.Y; y++ {
		var binLow, binHigh uint8

		for x := 0; x < imData.Bounds().Max.X; x++ {
			var lowBit, highBit uint8

			pixelColor := imData.At(x, y)

			// type assertion must be checked
			col, ok := pixelColor.(color.RGBA)
			if !ok {
				fmt.Println("Not RGBA")
				os.Exit(1)
			}

			r := col.R
			g := col.G
			b := col.B

			if r == g && g == b {
				switch r {
				case 0: // black, r = 3
					highBit = 1
					lowBit = 1
				case 85: // dark gray, r = 2
					highBit = 1
					lowBit = 0
				case 170: // light gray, r = 1
					highBit = 0
					lowBit = 1
				case 255: // white, r = 0
					highBit = 0
					lowBit = 0
				}
			} else {
				// Find the closest palette colour
				if red, green, blue := GetPaletteColour(darkest, PaletteBGB); r == red && g == green && b == blue {
					highBit = 1
					lowBit = 1
				} else if red, green, blue := GetPaletteColour(dark, PaletteBGB); r == red && g == green && b == blue {
					highBit = 1
					lowBit = 0
				} else if red, green, blue := GetPaletteColour(light, PaletteBGB); r == red && g == green && b == blue {
					highBit = 0
					lowBit = 1
				} else if red, green, blue := GetPaletteColour(lightest, PaletteBGB); r == red && g == green && b == blue {
					highBit = 0
					lowBit = 0
				} else {
					panic("Unknown colour")
				}
			}

			binLow += lowBit * uint8(math.Pow(2, float64(7-x)))
			binHigh += highBit * uint8(math.Pow(2, float64(7-x)))
		}

		binCode = append(binCode, binLow, binHigh)
	}

	return binCode
}

func checkColor(imData image.Image) {
	// Make sure it is 32-bit RGBA color, each R,G,B, A component requires 8-bits
	if imData.ColorModel() != color.RGBAModel {
		fmt.Println("Not RGBA")
		fmt.Println("Color model:", imData.ColorModel())

		switch imData.ColorModel() {
		case color.RGBAModel: // 32-bit RGBA color, each R,G,B, A component requires 8-bits
			fmt.Println("RGBA")
		case color.GrayModel: // 8-bit grayscale
			fmt.Println("Gray")
		case color.NRGBAModel: // 32-bit non-alpha-premultiplied RGB color, each R,G,B component requires 8-bits
			fmt.Println("NRGBA")
		case imData.ColorModel(): // 32-bit non-alpha-premultiplied YCbCr color, each Y,Cb,Cr component requires 8-bits
			fmt.Println("NYCbCrA")
		case color.YCbCrModel: // 24-bit YCbCr color, each Y,Cb,Cr component requires 8-bits
			fmt.Println("YCbCr")
		default:
			fmt.Println("Unknown")
		}

		os.Exit(1)
	}
}

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

	// Save unique tiles to individual files
	for i, tile := range uniqueTiles {
		filename := fmt.Sprintf("tile%d.png", i)
		file, err := os.Create(filename)
		if err != nil {
			return nil, fmt.Errorf("failed to create file %s: %w", filename, err)
		}
		defer file.Close()

		if err := png.Encode(file, tile); err != nil {
			return nil, fmt.Errorf("failed to save tile %s: %w", filename, err)
		}
	}

	// Save the final image with unique tiles
	uniqueTilesImage := createImageFromTiles(uniqueTiles, tilesPerRow, numRows)
	file, err := os.Create(fmt.Sprintf("unique_tiles-%s.png", filename))
	if err != nil {
		return nil, fmt.Errorf("failed to create file unique_tiles.png: %w", err)
	}
	defer file.Close()

	if err := png.Encode(file, uniqueTilesImage); err != nil {
		return nil, fmt.Errorf("failed to save unique_tiles.png: %w", err)
	}

	return uniqueTiles, nil
}
