package gif

import (
	"fmt"
	"image"
	"image/draw"
	"image/gif"
	"io"
	"io/ioutil"
	"os"
)

var (
	// FramePath is frame path of gif image
	FramePath = "./frames/"
)

// Split decode reads and analyzes the given reader as a GIF image
func Split(reader io.Reader) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("Error while decoding: %s", r)
		}
	}()

	gifbin, err := gif.DecodeAll(reader)

	if err != nil {
		return err
	}

	imgWidth, imgHeight := getGifDimensions(gifbin)

	overpaintImage := image.NewRGBA(image.Rect(0, 0, imgWidth, imgHeight))

	if _, err := os.Stat(FramePath); os.IsNotExist(err) {
		os.Mkdir(FramePath, 0777)
	}

	for i, srcImg := range gifbin.Image {
		draw.Draw(overpaintImage, overpaintImage.Bounds(), srcImg, image.ZP, draw.Over)
		file, err := os.Create(fmt.Sprintf("%s%d%s", FramePath, i, ".gif"))
		if err != nil {
			return err
		}

		opts := &gif.Options{}
		opts.NumColors = 256
		err = gif.Encode(file, overpaintImage, opts)

		if err != nil {
			return err
		}

		file.Close()
	}

	return nil
}

// EncodeAll creates gif animation from frames
func EncodeAll() (err error) {
	files, err := ioutil.ReadDir(FramePath)
	if err != nil {
		return err
	}

	outGIF := &gif.GIF{}
	for _, file := range files {
		f, err := os.Open(FramePath + file.Name())
		if err != nil {
			fmt.Println("Error opening file")
		}
		inGIF, _ := gif.Decode(f)
		f.Close()

		outGIF.Image = append(outGIF.Image, inGIF.(*image.Paletted))
		outGIF.Delay = append(outGIF.Delay, 0)
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
