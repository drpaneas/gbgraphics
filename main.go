package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
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

var rangeStartOffset, rangeLength int32

func main() {
	var (
		args           args
		outputFilename string
	)

	arg.MustParse(&args)

	// There are 64 total pixels in a single tile (8x8 pixels).
	// Therefore, exactly 128 bits, or 16 bytes,
	// are required to fully represent a single tile.
	width := 8
	bitDepth := 2
	path := ""
	path = args.Rom
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
		tileBits := 8 * 8 * bitDepth // 8x8 pixels, 2 bits per pixel, 16 bytes per tile

		// Calculate the height of the img in bytes
		height := 8 * int(math.Ceil(float64(len(imageData))/float64(tileBits)))
		hexValue := fmt.Sprintf("% X", imageData)
		img := image.NewRGBA(image.Rect(0, 0, width, height))

		for x := 0; x < width; x++ {
			for y := 0; y < height; y++ {
				img.Set(x, y, color.White)
			}
		}

		var xPos, yPos = 0, 0

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

func getTiles(screenshot string, romBytes []byte) []string {
	// Load a screenshot from the emulator and split it into 8x8 tiles
	img := readImageFromFilePath(screenshot)
	tiles := split8x8(img)

	if len(tiles) != 23040/64 {
		fmt.Println("Not 23040/64 tiles")
		os.Exit(1)
	}

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
				if red, green, blue := GetPaletteColour(darkest, PaletteBGB); r == red && g == green && b == blue { // black, r = 3
					highBit = 1
					lowBit = 1
				} else if red, green, blue := GetPaletteColour(dark, PaletteBGB); r == red && g == green && b == blue { // dark gray, r = 2
					highBit = 1
					lowBit = 0
				} else if red, green, blue := GetPaletteColour(light, PaletteBGB); r == red && g == green && b == blue { // light gray, r = 1
					highBit = 0
					lowBit = 1
				} else if red, green, blue := GetPaletteColour(lightest, PaletteBGB); r == red && g == green && b == blue { // white, r = 0
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
