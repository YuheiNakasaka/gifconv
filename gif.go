/*
Package gifconv is simple gif manipurator.
*/
package gifconv

import (
	"fmt"
	"image"
	"image/draw"
	"image/gif"
	"io"
	"os"
)

// FramePath is output path of gif image frames
var FramePath = "./frames/"

// Split decode reads and analyzes the given reader as a GIF image
func Split(reader io.Reader) ([]string, error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Error while decoding: %s", r)
		}
	}()

	gifbin, err := gif.DecodeAll(reader)

	if err != nil {
		fmt.Println("Decode error")
	}

	imgWidth, imgHeight := getGifDimensions(gifbin)

	overpaintImage := image.NewRGBA(image.Rect(0, 0, imgWidth, imgHeight))

	if _, err = os.Stat(FramePath); os.IsNotExist(err) {
		os.Mkdir(FramePath, 0777)
	}

	framePaths := make([]string, 0, 0)
	for i, srcImg := range gifbin.Image {
		draw.Draw(overpaintImage, overpaintImage.Bounds(), srcImg, image.ZP, draw.Over)
		filePath := fmt.Sprintf("%s%04d%s", FramePath, i, ".gif")
		file, _ := os.Create(filePath)
		framePaths = append(framePaths, filePath)

		opts := &gif.Options{}
		opts.NumColors = 256
		err = gif.Encode(file, overpaintImage, opts)

		if err != nil {
			fmt.Println(err)
		}

		file.Close()
	}

	return framePaths, err
}

// EncodeAll creates gif animation from frames
func EncodeAll(files []string, delay int) (err error) {
	outGIF := &gif.GIF{}
	for _, file := range files {
		f, err := os.Open(file)
		if err != nil {
			fmt.Println("Error opening file")
		}
		inGIF, _ := gif.Decode(f)
		f.Close()

		outGIF.Image = append(outGIF.Image, inGIF.(*image.Paletted))
		outGIF.Delay = append(outGIF.Delay, delay)
	}

	f, _ := os.OpenFile(fmt.Sprintf("%s%s", FramePath, "output.gif"), os.O_WRONLY|os.O_CREATE, 0666)
	defer f.Close()
	gif.EncodeAll(f, outGIF)

	return nil
}

// get logical screen size
func getGifDimensions(gif *gif.GIF) (x, y int) {
	var lowestX int
	var lowestY int
	var highestX int
	var highestY int

	for _, img := range gif.Image {
		if img.Rect.Min.X < lowestX {
			lowestX = img.Rect.Min.X
		}
		if img.Rect.Min.Y < lowestY {
			lowestY = img.Rect.Min.Y
		}
		if img.Rect.Max.X > highestX {
			highestX = img.Rect.Max.X
		}
		if img.Rect.Max.Y > highestY {
			highestY = img.Rect.Max.Y
		}
	}

	return highestX - lowestX, highestY - lowestY
}
