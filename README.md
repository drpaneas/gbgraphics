# GBGraphics

I always wanted to re-create my favorite retro-games using modern Game-Engines, such as Godot or Unity.
To do that, I need the game's graphics, but I am no designer, nor I am interested in (illegally) downloading copyrighted material (such as ROMs and its tiles/sprites).
So I decided to extract the graphics from my own ROMs, and use them.
This tool, does exactly that, it extracts the graphics from a ROM, and saves them as PNGs.
But to do that, it needs to know _where_ the graphics are (memory address).
Since there are many ROMs, and many versions of the same game, I decided to make this tool which supports multiple ROMs, and multiple versions of the same game.
How?
By using a reference screenshot image.
You have to take a screenshot of a specific scene in the game, the one you are interested in ripping its graphic assets and use that as a reference.
The tool will analyze the image, and extract the graphics from the actual ROM.
In this way, I can (programmatically) extract the original 8x8 tiles used for a specific scene in any Gameboy game.
Having these 8x8 tiles, I can assemble them using [Aseprite](https://www.aseprite.org/) into a canonical game object (in case it consists of multiple tiles), and then export it as a PNG file.

## In a nutshell

1. Takes two inputs: a ROM file, and a reference screenshot image from the game.
2. Chops the reference screenshot into individual 8x8 images, decodes the DMG color palette and converts them back to 2BPP format.
3. Searches into the ROM file to find any of these 8x8 2BPP-format images.
4. Saves any of the findings as PNG formatted images.

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

<style>
    a { color: #E0F8D0 }
    b { color: #88C070 }
    c { color: #346856 }
    d { color: #081820 }
</style>

- <a>{0xE0, 0xF8, 0xD0}, // lightest (#E0F8D0)</a>
- <b>{0x88, 0xC0, 0x70}, // light (#88C070)</b>
- <c>{0x34, 0x68, 0x56}, // dark (#346856)</c>
- <d>{0x08, 0x18, 0x20}, // darkest (#081820)</d>

NOTE: It's very important to take a screenshot with these specifications, otherwise this tool won't work!
NOTE: Read [Gameboy 2BPP Graphics Format](https://www.huderlem.com/demos/gameboy2bpp.html) article by [Huderlem](https://www.huderlem.com/) for further details.

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

## Usage

```bash
GBGraphics - extract graphics from Gameboy ROM using a screenshot
git commit 6e5708bd3fe042b3f035d8182edbe1f7b61a8e14
Usage: gbgraphics --img <SCREENSHOT> [--output <FILE>] ROM

    Positional arguments:
    ROM                    Path to the ROM file

    Options:
    --img <SCREENSHOT>     path of in-game screenshot
    --output <FILE>        output file [default: out.png]
    --help, -h             display this help and exit
    --version              display version and exit
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

## How to build it

```bash
# Clone the repo
git clone https://github.com/drpaneas/gbgraphics; cd gbgraphics

# Build it (with versioning info if you have Go 1.18)
go build -buildvcs
```