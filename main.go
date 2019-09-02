package main

import (
	"errors"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

type col struct {
	r, g, b uint32
}

func (c col) RGBA() (uint32, uint32, uint32, uint32) {
	return c.r, c.g, c.b, 0
}

func main() {
	rootCmd := &cobra.Command{
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			filename := args[0]
			img, err := readImage(filename)
			if err != nil {
				panic(err)
			}

			bounds := img.Bounds()

			red := image.NewRGBA(image.Rectangle{
				image.Point{0, 0},
				image.Point{bounds.Max.X, bounds.Max.Y},
			})

			green := image.NewRGBA(image.Rectangle{
				image.Point{0, 0},
				image.Point{bounds.Max.X, bounds.Max.Y},
			})

			blue := image.NewRGBA(image.Rectangle{
				image.Point{0, 0},
				image.Point{bounds.Max.X, bounds.Max.Y},
			})

			alpha := image.NewRGBA(image.Rectangle{
				image.Point{0, 0},
				image.Point{bounds.Max.X, bounds.Max.Y},
			})

			width, height := bounds.Dx(), bounds.Dy()
			for y := 0; y < height; y++ {
				for x := 0; x < width; x++ {
					pixel := img.At(x, y)
					r, g, b, a := pixel.RGBA()

					red.Set(x, y, col{r, r, r})
					green.Set(x, y, col{g, g, g})
					blue.Set(x, y, col{b, b, b})
					alpha.Set(x, y, col{a, a, a})
				}
			}

			extension := filepath.Ext(filename)
			name := filename[0 : len(filename)-len(extension)]
			redname := name + "-red.jpg"
			greenname := name + "-green.jpg"
			bluename := name + "-blue.jpg"
			alphaname := name + "-alpha.jpg"

			err = writeImage(redname, red)
			if err != nil {
				panic(err)
			}

			err = writeImage(greenname, green)
			if err != nil {
				panic(err)
			}

			err = writeImage(bluename, blue)
			if err != nil {
				panic(err)
			}

			err = writeImage(alphaname, alpha)
			if err != nil {
				panic(err)
			}
		},
	}

	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}

func readImage(filename string) (image.Image, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	ext := strings.ToLower(path.Ext(filename))
	switch ext {
	case ".jpg":
		fallthrough
	case ".jpeg":
		return jpeg.Decode(file)
	case ".png":
		return png.Decode(file)
	default:
		return nil, errors.New("unsupported output type: " + filename)
	}
}

func writeImage(filename string, img image.Image) error {
	if img == nil {
		return errors.New("nil image")
	}

	ext := path.Ext(filename)
	switch ext {
	case ".jpg":
		f, err := os.Create(filename)
		if err != nil {
			return err
		}
		return jpeg.Encode(f, img, &jpeg.Options{
			Quality: 100,
		})
	case ".png":
		f, err := os.Create(filename)
		if err != nil {
			return err
		}
		return png.Encode(f, img)
	default:
		return errors.New("unsupported output type")
	}
}
