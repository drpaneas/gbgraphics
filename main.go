package main

import (
	"bufio"
	"fmt"
	"github.com/alexflint/go-arg"
	"golang.org/x/image/draw"
	"image"
	"image/color"
	"image/png"
	"log"
	"math"
	"os"
	"runtime/debug"
	"strconv"
	"strings"
	"time"
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
	//Bpp    int    `arg:"--bpp" help:"bits per pixel" default:"2" placeholder:"<BITS>"`
	//Width  int    `arg:"--width" help:"width of the image" default:"8" placeholder:"<WIDTH>"`
	//Offset string `arg:"required,--offset" help:"beginning of the range" placeholder:"<OFFSET>"`
	//Length string `arg:"--length" help:"length of the range in hex" default:"0x10" placeholder:"<LENGTH>"`
	Output string `arg:"--output" help:"output file" default:"out.png" placeholder:"<FILE>"`
}

func (args) Description() string {
	desc := fmt.Sprintf("GBGraphics - extract graphics from Gameboy ROM using a screenshot")
	return desc
}

func (args) Version() string {
	return "git commit " + Commit
}

var rangeStartOffset, rangeLength int32

func main() {
	var args args
	arg.MustParse(&args)

	// There are 64 total pixels in a single tile (8x8 pixels).
	// Therefore, exactly 128 bits, or 16 bytes,
	// are required to fully represent a single tile.
	width := 128 // 16 bytes per tile
	bitDepth := 0
	path := ""
	outputFilename := ""

	// bitDepth = args.Bpp
	bitDepth = 2
	// width = args.Width
	width = 8
	outputFilename = args.Output
	path = args.Rom

	// --------------------------------
	screenshot := args.Screenshot
	romBytes, err := os.ReadFile(args.Rom)
	if err != nil {
		log.Fatal(err)
	}
	locations := getTiles(screenshot, romBytes)
	for i, v := range locations {

		outputFilename = args.Output
		outputFilename = strings.Replace(outputFilename, ".png", "", -1)
		outputFilename = fmt.Sprintf("%s_%d.png", outputFilename, i)

		// --------------------------------

		// startOffsetString := args.Offset
		startOffsetString := v
		// Remove the 0x prefix and keep only the actual number (as a string)
		if !strings.Contains(startOffsetString, "0x") {
			fmt.Println("Invalid start offset specified! Please specify a HEX value (e.g. 0x3f)")
			os.Exit(1)
		}
		// Remove the 0x prefix (it's not needed anymore)
		startOffsetString = strings.Replace(startOffsetString, "0x", "", -1)
		// and convert the strings (which represent a hex value) to a decimal int32
		a, _ := strconv.ParseInt(startOffsetString, 16, 32)
		rangeStartOffset = int32(a)

		// lengthString := args.Length
		lengthString := "0x10"
		// Remove the 0x prefix and keep only the actual number (as a string)
		if !strings.Contains(lengthString, "0x") {
			fmt.Println("Invalid length specified! Please specify a HEX value (e.g. 0x200/512)")
			os.Exit(1)
		}
		// Remove the 0x prefix (it's not needed anymore)
		lengthString = strings.Replace(lengthString, "0x", "", -1)
		// and convert the strings (which represent a hex value) to a decimal int32
		a, _ = strconv.ParseInt(lengthString, 16, 32)
		rangeLength = int32(a)

		if path == "" {
			fmt.Println("No ROM specified!")
			os.Exit(1)
		}

		if bitDepth == 0 {
			fmt.Println("No bit depth specified!")
			os.Exit(1)
		}

		if bitDepth != 1 && bitDepth != 2 {
			fmt.Println("Invalid bit depth specified!")
			os.Exit(1)
		}

		if width == 0 {
			fmt.Println("No width specified!")
			os.Exit(1)
		}

		if width%8 != 0 {
			fmt.Println("Invalid width specified!")
			os.Exit(1)
		}

		if rangeStartOffset < 0 {
			fmt.Println("Invalid start offset specified!")
			os.Exit(1)
		}

		if rangeLength < 0 {
			fmt.Println("Invalid length specified!")
			os.Exit(1)
		}

		// Load GB ROM
		imageData, err := os.ReadFile(path)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// If the user specified a file range, we only want to use the specified range
		// e.g. -r 0x640A0 0x10 (it's 16 bytes long)
		imageData = imageData[rangeStartOffset : rangeStartOffset+rangeLength]

		// Calculate the height of the img in bytes
		//
		// (e.g. if the img is 8x8 pixels, the height is 8/8 = 1 byte)
		// (e.g. if the img is 8x16 pixels, the height is 16/8 = 2 bytes)
		//
		// The GameBoy displays its graphics using 8x8-pixel tiles.
		// As the name 2BPP implies, it takes exactly two bits to store the information about a single pixel.
		// There are 64 total pixels in a single tile (8x8 pixels)
		// Therefore, exactly 128 bits, or 16 bytes, are required to fully represent a single tile.
		// As a result, any uncompressed graphics data present in a GameBoy ROM file is represented using exactly 16 bytes.
		tileBits := 8 * 8 * bitDepth // 8x8 pixels, 2 bits per pixel, 16 bytes per tile

		// Calculate the height of the img in bytes
		height := 8 * int(math.Ceil(float64(len(imageData))/float64(tileBits)))
		hexValue := fmt.Sprintf("% X", imageData)

		// fmt.Printf("%08b\n", imageData)
		img := image.NewRGBA(image.Rect(0, 0, width, height))
		for x := 0; x < width; x++ {
			for y := 0; y < height; y++ {
				img.Set(x, y, color.White)
			}
		}
		var xPos = 0
		var yPos = 0
		for i := 0; i < width*height; i += 8 * bitDepth {
			highBit := 0
			lowBit := 0
			colorVal := 0
			for x := 0; x < 8; x++ {
				for y := 0; y < 8; y++ {
					if (i+y >= len(imageData) && bitDepth == 1) || ((i+2*y+1 >= len(imageData)) && bitDepth == 2) {
						break
					}
					if bitDepth == 2 {
						highBit = int(imageData[i+2*y+1]>>(7-x)) & 0x01
						lowBit = int(imageData[i+2*y]>>(7-x)) & 0x01
						value := (highBit << 1) | lowBit
						colorVal = int(255 * (float32(3-value) / 3))
					} else {
						value := (imageData[i+y] >> (7 - x)) & 0x01
						colorVal = int(1-value) * 255
					}

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
		f, err := os.Create(outputFilename)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		defer func(f *os.File) {
			err := f.Close()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}(f)
		err = png.Encode(f, img)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Printf("'%s' (Found at location %s) converted to '%s'\n", hexValue, v, outputFilename)
	}
}

const (
	// PaletteGreyscale is the default greyscale gameboy colour palette.
	PaletteGreyscale = byte(iota)
	// PaletteOriginal is more authentic looking green tinted gameboy
	// colour palette  as it would have been on the GameBoy
	PaletteOriginal
	// PaletteBGB used by default in the BGB emulator.
	PaletteBGB
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
func GetPaletteColour(index byte, palette byte) (r, g, b uint8) {
	col := Palettes[palette][index]
	r, g, b = col[0], col[1], col[2]
	return r, g, b
}

func removeDuplicateInt(intSlice []int) []int {
	allKeys := make(map[int]bool)
	list := []int{}
	for _, item := range intSlice {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}

// Remove duplicates [][]byte
func removeDuplicateByte(byteSlice [][]byte) [][]byte {
	allKeys := make(map[string]bool)
	list := [][]byte{}
	for _, item := range byteSlice {
		if _, value := allKeys[string(item)]; !value {
			allKeys[string(item)] = true
			list = append(list, item)
		}
	}
	return list
}

func getTiles(screenshot string, romBytes []byte) []string {

	// Load a screenshot from the emulator and split it into 8x8 tiles
	img := readImageFromFilePath(screenshot)
	tiles := split8x8(img)
	if len(tiles) != 23040/64 {
		fmt.Println("Not 23040/64 tiles")
		os.Exit(1)
	} else {
		// fmt.Println("Screenshot is splitted into: ", 23040/64, " 8x8 tiles")
	}

	// These tiles are in RGBA format, so we need to convert them to 2BPP
	// before we can compare them to the original gameboy tileset
	origCodeTiles := getHexCodes(tiles)
	uniqCodeTiles := removeDuplicateByte(origCodeTiles)

	// Print the 2BPP tiles to the console
	//for _, tile := range uniqCodeTiles {
	//	// fmt.Printf("Tile %d: % 02X\n", i, tile)
	//	fmt.Printf("% 02X\n", tile)
	//}

	//// Case 1: If you already have a tile and want to see if it exists in the screenshot
	//// do this:
	//// Look at the original code tiles if a specific tile exists
	//specificImg := readImageFromFilePath(tileToLookFor)
	//specificImgCode := pngTo2BPP(specificImg, "specific_tile")
	//foundMsg := "We couldn't file the tile you were looking for"
	//for i, tile := range origCodeTiles {
	//	if compare(tile, specificImgCode) {
	//		foundMsg = fmt.Sprintf("Found tile %d\n", i)
	//	}
	//}
	//fmt.Println(foundMsg)

	// Case 2: Load the GB rom and search to find each tile in the screenshot
	// If you find it, print the address of the tile in the rom
	// do this:
	// Load GB ROM
	// Read the rom into a byte array

	// Search for each tile in the screenshot
	var addr []string
	for i, tile := range uniqCodeTiles {
		//// if tile is full of only 0 (plain white) skip it
		//if bytes.Compare(tile, []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}) == 0 {
		//	continue
		//}

		// Search for the tile in the rom
		for j := 0; j < len(romBytes); j++ {
			// check for out-of-bounds
			if j+16 > len(romBytes) {
				break
			}
			if compare(tile, romBytes[j:j+16]) {
				// fmt.Printf("Found tile %03d [% 02X] at address 0x%X\n", i, tile, j)
				tmp := fmt.Sprintf("0x%X", j)
				addr = append(addr, tmp)
				_ = i

				break // no reason to keep searching for this tile if you have already found it
			}
		}
	}

	return addr
}

func getHexCodes(tiles []image.Image) [][]byte {
	origCodeTiles := make([][]byte, 0)
	for i, tile := range tiles {
		destFilename := fmt.Sprintf("tile_%d", i)
		origCodeTiles = append(origCodeTiles, pngTo2BPP(tile, destFilename))
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
func pngTo2BPP(imData image.Image, dest string) []byte {
	// Make sure it is 8x8
	if imData.Bounds().Max.X != 8 || imData.Bounds().Max.Y != 8 {
		fmt.Println("Not 8x8")
		fmt.Println("It is: ", imData.Bounds().Max.X, imData.Bounds().Max.Y)
		os.Exit(1)
	}

	// Make sure it is 32-bit RGBA color, each R,G,B, A component requires 8-bits
	if imData.ColorModel() != color.RGBAModel {
		fmt.Println("Not RGBA")
		fmt.Println("Color model:", imData.ColorModel())
		if imData.ColorModel() == color.RGBAModel {
			//32-bit RGBA color, each R,G,B, A component requires 8-bits
			fmt.Println("RGBA")
		} else if imData.ColorModel() == color.GrayModel {
			//8-bit grayscale
			fmt.Println("Gray")
		} else if imData.ColorModel() == color.NRGBAModel {
			//32-bit non-alpha-premultiplied RGB color, each R,G,B component requires 8-bits
			fmt.Println("NRGBA")
		} else if imData.ColorModel() == color.NYCbCrAModel {
			//32-bit non-alpha-premultiplied YCbCr color, each Y,Cb,Cr component requires 8-bits
			fmt.Println("NYCbCrA")
		} else if imData.ColorModel() == color.YCbCrModel {
			//24-bit YCbCr color, each Y,Cb,Cr component requires 8-bits
			fmt.Println("YCbCr")
		} else {
			fmt.Println("Unknown")
		}
		os.Exit(1)
	}

	// print image pixel color r,g,b,a into uint32
	var binCode []byte
	for y := 0; y < imData.Bounds().Max.Y; y++ {
		var binLow uint8
		var binHigh uint8
		for x := 0; x < imData.Bounds().Max.X; x++ {
			pixelColor := imData.At(x, y)
			var highBit uint8
			var lowBit uint8

			r := pixelColor.(color.RGBA).R
			g := pixelColor.(color.RGBA).G
			b := pixelColor.(color.RGBA).B
			// fmt.Printf("0x%02X, 0x%02X, 0x%2X\n", r, g, b)
			if r == g && g == b {
				switch r {
				case 0:
					r = 3 // black
					highBit = 1
					lowBit = 1
				case 85:
					r = 2 // dark gray
					highBit = 1
					lowBit = 0
				case 170:
					r = 1 // light gray
					highBit = 0
					lowBit = 1
				case 255:
					r = 0 // white
					highBit = 0
					lowBit = 0
				}
			} else {
				// Find the closest palette colour
				if red, green, blue := GetPaletteColour(darkest, PaletteBGB); r == red && g == green && b == blue {
					r = 3 // black
					highBit = 1
					lowBit = 1
				} else if red, green, blue := GetPaletteColour(dark, PaletteBGB); r == red && g == green && b == blue {
					r = 2 // dark gray
					highBit = 1
					lowBit = 0
				} else if red, green, blue := GetPaletteColour(light, PaletteBGB); r == red && g == green && b == blue {
					r = 1 // light gray
					highBit = 0
					lowBit = 1
				} else if red, green, blue := GetPaletteColour(lightest, PaletteBGB); r == red && g == green && b == blue {
					r = 0 // white
					highBit = 0
					lowBit = 0
				} else {
					panic("Unknown colour")
				}
			}

			binLow += lowBit * uint8(math.Pow(2, float64(7-x)))
			binHigh += highBit * uint8(math.Pow(2, float64(7-x)))
			//fmt.Printf("%v ", r) // Now in range 0..255
		}
		//fmt.Printf("%02X %02X ", binLow, binHigh) // Now in range 0..255
		binCode = append(binCode, binLow, binHigh)
		//fmt.Print("\n") // Change line
	}

	// fmt.Printf("% 02X\n", binCode)
	//fmt.Println()

	//// Open a new file for writing only
	//file, err := os.OpenFile(
	//	fmt.Sprintf("%s.2bpp", dest),
	//	os.O_WRONLY|os.O_TRUNC|os.O_CREATE,
	//	0666,
	//)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//defer func(file *os.File) {
	//	err := file.Close()
	//	if err != nil {
	//		fmt.Println(err)
	//	}
	//}(file)
	//
	//// Write bytes to file
	//bytesWritten, err := file.Write(binCode)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//log.Printf("Wrote %d bytes.\n", bytesWritten)

	return binCode
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
		os.Exit(1)
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

	// Check if image is 32-bit RGBA color, each R,G,B, A component requires 8-bits
	if src.ColorModel() != color.RGBAModel {
		fmt.Println("Not RGBA")
		fmt.Println("Color model:", src.ColorModel())
		if src.ColorModel() == color.RGBAModel {
			//32-bit RGBA color, each R,G,B, A component requires 8-bits
			fmt.Println("RGBA")
		} else if src.ColorModel() == color.GrayModel {
			//8-bit grayscale
			fmt.Println("Gray")
		} else if src.ColorModel() == color.NRGBAModel {
			//32-bit non-alpha-premultiplied RGB color, each R,G,B component requires 8-bits
			fmt.Println("NRGBA")
		} else if src.ColorModel() == color.NYCbCrAModel {
			//32-bit non-alpha-premultiplied YCbCr color, each Y,Cb,Cr component requires 8-bits
			fmt.Println("NYCbCrA")
		} else if src.ColorModel() == color.YCbCrModel {
			//24-bit YCbCr color, each Y,Cb,Cr component requires 8-bits
			fmt.Println("YCbCr")
		} else {
			fmt.Println("Unknown")
		}
		os.Exit(1)
	}

	// iterate over the image and split it into 8x8 tiles
	// var tiles []image.Image

	//for y := 0; y < src.Bounds().Max.Y; y += 8 {
	//	for x := 0; x < src.Bounds().Max.X; x += 8 {
	//		tile := src.(*image.RGBA).SubImage(image.Rect(x, y, x+8, y+8))
	//		tiles = append(tiles, tile)
	//	}
	//}

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

	//var images []image.Image
	//for y := 0; y < src.Bounds().Max.Y; y += 8 {
	//	for x := 0; x < src.Bounds().Max.X; x += 8 {
	//		rect := image.Rect(x, y, x+8, y+8)
	//		images = append(images, src.(interface {
	//			SubImage(r image.Rectangle) image.Image
	//		}).SubImage(rect))
	//	}
	//}
	//return images
}

func printSize(path string) {
	// Open file.
	inputFile, _ := os.Open(path)
	defer func(inputFile *os.File) {
		err := inputFile.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(inputFile)
	reader := bufio.NewReader(inputFile)

	// Decode image and get its settings.
	config, _, _ := image.DecodeConfig(reader)

	// Print config.
	fmt.Printf("IMAGE: width=%v height=%v\n", config.Width, config.Height)
}

func resizeImage(path, saveAs string) {
	inputFile2, _ := os.Open(path)
	defer func(inputFile *os.File) {
		err := inputFile.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(inputFile2)

	// Decode image and get the image object.
	src, _, err := image.Decode(inputFile2)
	if err != nil {
		log.Fatal("Error ", err)
	}

	// Resize image
	// new size of image
	dr := image.Rect(0, 0, src.Bounds().Max.X/2, src.Bounds().Max.Y/2)
	// resize using given scaler
	var res image.Image
	{ // show time to resize
		tp := time.Now()
		// perform resizing
		res = scaleTo(src, dr)
		// report time to scaling to console
		log.Printf("scaling using %q takes %v time",
			"NearestNeighbor", time.Now().Sub(tp))
	}
	// open file to save
	dstFile, err := os.Create(saveAs + ".png")
	if err != nil {
		log.Fatal(err)
	}
	// encode as .png to the file
	err = png.Encode(dstFile, res)
	// close the file
	err = dstFile.Close()
	if err != nil {
		return
	}

	if err != nil {
		log.Fatal(err)
	}
}

func scaleTo(src image.Image, rect image.Rectangle) image.Image {
	var scale draw.Scaler = draw.NearestNeighbor
	dst := image.NewRGBA(rect)
	scale.Scale(dst, rect, src, src.Bounds(), draw.Over, nil)
	return dst
}
