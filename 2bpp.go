package main

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"os"
)

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
