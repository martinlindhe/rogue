package rogue

import (
	"image"
	"image/color"
	"image/png"
	"math"
	"os"

	"github.com/martinlindhe/rogue/rollingparticle"
	"github.com/ojrac/opensimplex-go"
)

// Island ...
type Island struct {
	Width     int
	Height    int
	Seed      int64
	HeightMap [][]byte
	Spawns    []Npc
}

// GenerateIsland returns a new island
func GenerateIsland(seed int64, width int, height int) Island {

	var island Island

	island.Width = width
	island.Height = height
	island.Seed = seed

	particleLength := 6
	outerBlur := 0.95
	innerBlur := 0.70
	roller := rollingparticle.New(seed, island.Width, island.Height, particleLength, innerBlur, outerBlur)

	rollerImage := Slice2DAsImage(&roller, island.Width, island.Height)
	rollerImgFile, _ := os.Create("roller.png")
	png.Encode(rollerImgFile, rollerImage)

	m := make2DByteSlice(width, height)

	noise := opensimplex.NewWithSeed(seed)

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {

			// combine some sizes of noise
			fBig := noise.Eval2(float64(x)*0.005, float64(y)*0.005)
			fMid := noise.Eval2(float64(x)*0.01, float64(y)*0.01)
			fSmall := noise.Eval2(float64(x)*0.02, float64(y)*0.02)
			fMini := noise.Eval2(float64(x)*0.04, float64(y)*0.04)

			f := (fBig + fMid + fSmall + fMini) / 4

			// scale from -1.0-1.0 to 0.0-1.0
			f = (f + 1.0) / 2.0

			// scale from 0.0-1.0 to 0-255
			b := byte(0)
			if f == 1.0 {
				b = 255
			} else {
				b = byte(math.Floor(f * 256.0))
			}

			// combine with rolling particle
			opacity := 0.5
			b = byte((1-opacity)*float64(b) + opacity*float64(roller[y][x]))

			m[y][x] = b
		}
	}

	island.HeightMap = m
	return island
}

func (i *Island) Add(npc Npc) {
	i.Spawns = append(i.Spawns, npc)
}

func (i *Island) ColoredHeightMapAsImage() image.Image {

	img := image.NewRGBA(image.Rectangle{image.Point{0, 0}, image.Point{i.Width, i.Height}})

	for y := 0; y < i.Height; y++ {
		for x := 0; x < i.Width; x++ {
			b := i.HeightMap[y][x]

			var col color.RGBA
			switch {
			case b < 80:
				col = color.RGBA{0x26, 0x2f, 0x71, 0xff} // deep water

			case b < 90:
				col = color.RGBA{0x46, 0x4D, 0x85, 0xff} // shallow water

			case b < 100:
				col = color.RGBA{0xD4, 0xBC, 0x6A, 0xff} // beach

			case b < 150:
				col = color.RGBA{0x2D, 0x88, 0x2D, 0xff} // grass (green)

			case b < 230:
				col = color.RGBA{0x00, 0x4E, 0x00, 0xff} // forest (dark green)

			case b < 240:
				col = color.RGBA{0x4B, 0x2D, 0x12, 0xff} // hills (brown)

			default:
				col = color.RGBA{0xF2, 0xED, 0xE6, 0xff} // gray (mountains)
			}

			img.Set(x, y, col)
		}
	}

	return img

}

// HeightMapAsImage renders height map to an Image
func (i *Island) HeightMapAsImage() image.Image {

	return Slice2DAsImage(&i.HeightMap, i.Width, i.Height)
}
