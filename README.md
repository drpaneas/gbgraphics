<img src="./gbgraphics-logo.png?sanitize=true" alt="GBGraphics" width="240">

A tool that extracts the graphics from any Gameboy ROM (into PNGs) for a given screenshot of the game.

In case a sprite consists of multiple tiles, you can combine them (e.g. using [Aseprite](https://www.aseprite.org/))
to create the canonical game asset and then export it as a PNG file.

## Usage

```bash
GBGraphics - extract graphics from Gameboy ROM using a screenshot
git commit 6e5708bd3fe042b3f035d8182edbe1f7b61a8e14
Usage: gbgraphics --img SCREENSHOT [--output FILE] ROM

Positional arguments:
ROM                    Path to the ROM file

Options:
--img SCREENSHOT     path of in-game screenshot
--output FILE        output file [default: out.png]
--help, -h           display this help and exit
--version            display version and exit
```

### Example
```bash
$ ./gbgraphics --img screen.png pokemon.gb
```

Output:

```
'00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00' (Found at location 0xBE) converted to 'out_0.png'
[... omitted for brevity ...]
'00 00 60 60 0F 0E 6D 6D 6D 6D 6D 6D 00 00 00 00' (Found at location 0x121D8) converted to 'out_97.png'
'00 00 00 00 78 38 60 60 63 63 7B 3B 00 00 00 00' (Found at location 0x121E8) converted to 'out_98.png'
```

Screenshot (`screen.png`) used in this example, taken from my modified GoBoy emulator:

![screen.png](screen.png)

The unique tiles extracted from this one:

![tiles.png](tiles.png)

## For Developers

```bash
# Clone the repo
git clone https://github.com/drpaneas/gbgraphics; cd gbgraphics

# Build it (with versioning info if you have Go 1.18)
go build -buildvcs
```

It works like this:

1. Takes two inputs: a ROM file, and a reference screenshot from the game.
2. Chops the reference screenshot into individual 8x8 images, decodes the DMG color palette and converts them back to 2BPP format.
3. Searches into the ROM file to find any of these 8x8 2BPP-format images.
4. Saves any of the findings as PNG formatted images.

NOTE: Read [Gameboy 2BPP Graphics Format](https://www.huderlem.com/demos/gameboy2bpp.html) article by [Huderlem](https://www.huderlem.com/) for further details.


## For Users

### Input 1: Get a ROM file

First step is to dump a Rom file (e.g. `game.gb`) taken from your GameBoy physical cartridge you already own.
You will need a cartridge reader, such as [Sanni's Cardreader](https://github.com/sanni/cartreader), which dumps pretty much every well-know retro-console.
There are other card-readers, dedicated to dumping Gameboy cards only, such as [GBxCart RW v1.4 Pro](https://retrogamerepairshop.com/products/gbxcart-rw-gameboy-gbc-gba-cart-reader-writer-flasher).

### Input 2: Take a screenshot

Launch your favorite Gameboy emulator and start playing the game until you find the scene you are interested in ripping its graphics.
The requirements for a proper screenshot are the following:

1. image type: ARGB
2. resolution: 160×144 pixels (it's the native GameBoy res)
3. color palette: DMG (GameBoy) palette

The gameboy palette is a 4-color palette, which is used to render the 2BPP images, is the following:

| Color | Hex | RGB | RGB (in HEX) |
| --- | --- | --- | --- |
| Lightest | #E0F8D0 | 224, 248, 208 | 0xE0, 0xF8, 0xD0 |
| Light | #88C070 | 136, 192, 112 | 0x88, 0xC0, 0x70 |
| Dark  | #346856 | 52, 104, 86 | 0x34, 0x68, 0x56 |
| Darkest | #081820 | 8, 24, 32 | 0x08, 0x18, 0x20 |

**NOTE**: It's very important to take a screenshot with these specifications, otherwise this tool won't work!

#### Optional step (recommended): using GoBoy emulator

In my case, I am using my modified version of [GoBoy](https://github.com/drpaneas/goboy) emulator.
Here's the instructions for you, to follow my example:

```bash
# Clone my fork
git clone https://github.com/drpaneas/goboy; cd goboy

# Build it (you need to have Go installed)
go build -o goboy cmd/goboy/main.go

# Run it against your ROM using the DMG color palette
./goboy -dmg pokemon.gb
```

The game will start playing in a tiny borderless window!
Press key `t` to take a screenshot. It will be saved locally with `screenshot.png` filename.
