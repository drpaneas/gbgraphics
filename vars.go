package main

var rangeStartOffset int32

const (
	width           = 8
	bitDepth        = 2
	rangeLength     = 16
	gbScreenXRes    = 160
	gbScreenYRes    = 144
	pixelsPerScreen = gbScreenXRes * gbScreenYRes
	pixelsPerTile   = 8 * 8
	tilesPerRow     = gbScreenXRes / 8
	tilesPerCol     = gbScreenYRes / 8
)
