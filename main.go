package main

import (
	"fmt"
	"github.com/alexflint/go-arg"
	"image"
	"image/color"
	"image/png"
	"log"
	"math"
	"os"
	"runtime/debug"
	"strconv"
	"strings"
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
	Rom    string `arg:"positional,required"`
	Bpp    int    `arg:"--bpp" help:"Bits Per Pixel" default:"2" placeholder:"<BITS>"`
	Width  int
	Offset int
	Length int
	Output string
}

func (args) Description() string {
	desc := fmt.Sprintf("GBGraphics v%s ( https://github.com/drpaneas/gbgraphics ) - converts a portion of binary file to a PNG image", Commit)
	return desc
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
	//img := Bitmap{}
	useFileRange := false
	outputFilename := ""

	if len(os.Args) == 1 {
		// No arguments were specified!
		fmt.Println("Usage: gbgraphics <img path> -b -w -r -o")
		fmt.Println()
		fmt.Println("Optional arguments:")
		fmt.Println("-b: bits per pixel(1 or 2)(e.g. -b2)")
		fmt.Println("-w: width(e.g. -w32)")
		fmt.Println("-r: file range in HEX (start offset, length in bytes)(e.g. -r 0x3f 0x200/512)")
		fmt.Println("-o: Output file name (e.g. -o font.png)")
		os.Exit(1)
	} else {
		// Parse arguments
		for i := 1; i < len(os.Args); i++ {
			switch os.Args[i] {
			case "-b":
				// Bits per pixel
				bitDepth, _ = strconv.Atoi(os.Args[i+1])
			case "-w":
				width, _ = strconv.Atoi(os.Args[i+1])
			case "-r":
				useFileRange = true

				// We expect two arguments after this one
				// both arguments should be in HEX (the user specifies sth like: 0x3f)
				startOffsetString := os.Args[i+1]
				lengthString := os.Args[i+2]

				// Remove the 0x prefix and keep only the actual number (as a string)
				if !strings.Contains(startOffsetString, "0x") {
					fmt.Println("Invalid start offset specified! Please specify a HEX value (e.g. 0x3f)")
					os.Exit(1)
				}

				if !strings.Contains(lengthString, "0x") {
					fmt.Println("Invalid length specified! Please specify a HEX value (e.g. 0x200/512)")
					os.Exit(1)
				}

				// Remove the 0x prefix (it's not needed anymore)
				startOffsetString = strings.Replace(startOffsetString, "0x", "", -1)
				lengthString = strings.Replace(lengthString, "0x", "", -1)

				// and convert the strings (which represent a hex value) to a decimal int32
				a, _ := strconv.ParseInt(startOffsetString, 16, 32)
				rangeStartOffset = int32(a)
				a, _ = strconv.ParseInt(lengthString, 16, 32)
				rangeLength = int32(a)
				i += 2
			case "-o":
				outputFilename = os.Args[i+1]
				i++
			default:
				path = os.Args[1] // Path to img is always the first argument
			}
		}
	}
	// Debug: info for arguments
	fmt.Printf("path: %v\nflag '-b': %v\nflag '-w': %v\nflag '-r': %v %v\nflag '-o': %v\nuseFileRange: %v\n", path, bitDepth, width, rangeStartOffset, rangeLength, outputFilename, useFileRange)

	// Load GB ROM
	imageData, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	// If the user specified a file range, we only want to use the specified range
	if useFileRange {
		// e.g. -r 0x640A0 0x10 (it's 16 bytes long)
		imageData = imageData[rangeStartOffset : rangeStartOffset+rangeLength]
	}

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
	fmt.Println(height)
	fmt.Printf("% X\n", imageData)
	fmt.Printf("%08b\n", imageData)
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
		log.Fatal(err)
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(f)
	err = png.Encode(f, img)
	if err != nil {
		log.Fatal(err)
	}

}
