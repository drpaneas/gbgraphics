package main

import (
	"fmt"
	"image"
	"image/color"
	"os"
)

const (
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
