package main

import (
	"fmt"
	"image"
	"image/png"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

func decodeFiles(args []string) []image.Image {
	var images []image.Image

	for _, filename := range args {
		f, err := os.Open(filename)
		if err != nil {
			fmt.Printf("Could't open the file \"%s\" -- %s\n", filename, err)
			continue
		}
		defer f.Close()

		img, err := png.Decode(f)
		if err != nil {
			fmt.Printf("Couldn't decode the file \"%s\" -- %s\n", filename, err)
			continue
		}

		images = append(images, img)
	}

	return images
}

func imagesToASCII(images []image.Image) []string {
	const symbols = "@%#*+=-:. "

	asciis := make([]string, len(images))
	for i, image := range images {
		rect := image.Bounds()

		var currentString strings.Builder
		for y := rect.Min.Y; y < rect.Max.Y; y++ {
			for x := rect.Min.X; x < rect.Max.X; x++ {
				r, g, b, a := image.At(x, y).RGBA()

				luminance := (.2126*float64(r) + .7152*float64(g) + .0722*float64(b)) / 65535

				index := len(symbols) - 1
				if a != 0 {
					index = int(luminance * float64(len(symbols)-1))
				}

				currentString.WriteByte(symbols[index])
			}

			currentString.WriteByte('\n')
		}

		asciis[i] = currentString.String()
	}

	return asciis
}

func saveASCIIs(asciis []string) {
	for i, ascii := range asciis {
		filename := strconv.Itoa(i) + ".txt"
		err := ioutil.WriteFile(filename, []byte(ascii), 0644)
		if err != nil {
			fmt.Printf("Couldn't save ASCII -- %s", err)
			continue
		}
	}
}

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		fmt.Println("Usage: ./png2ascii FILE1 [FILE2 FILE3...]")
		return
	}

	images := decodeFiles(args)
	fmt.Println("Decoded the images")

	asciis := imagesToASCII(images)
	fmt.Println("Converted the images to ASCII")

	saveASCIIs(asciis)
	fmt.Println("Saved the images")

	fmt.Println("Done!")
}
