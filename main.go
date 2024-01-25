package main

import (
	"fmt"
	"os"
	"runtime/debug"

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
	return "Version (git commit):" + Commit
}

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

	// loop from 0 to 7
	var locations [][]string
	for i := 0; i < 8; i++ {
		locations = append(locations, getTiles(screenshot, romBytes, i))
	}

	// Create a slice of all the addresses
	var allAddresses []string
	for _, location := range locations {
		allAddresses = append(allAddresses, location...)
	}

	uniqueAddresses := removeDuplicateString(allAddresses)

	for i, address := range uniqueAddresses {
		// for every address, get the tile and save it to disk
		//tile := romBytes[convertHexToInt32(address) : convertHexToInt32(address)+16]
		//fmt.Printf("Address: %s, Tile: %v\n", address, hex.EncodeToString(tile))
		if err := processTile(i, address, outputFilename, romBytes, rangeLength, width, bitDepth); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
}
