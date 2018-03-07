/*
Package gifconv is simple gif manipurator.
*/
package gifconv

import (
	"fmt"
	"image/gif"
	"os"
)

// Info returns gif infomation
func Info(filePath string) error {
	file, _ := os.Open(filePath)
	defer file.Close()

	gifBin, err := gif.DecodeAll(file)
	if err != nil {
		fmt.Println("File is not gif")
	}

	fmt.Printf("logical screen %dx%d\n", gifBin.Config.Width, gifBin.Config.Height)
	fmt.Printf("%s is %d images\n", filePath, len(gifBin.Image))

	return err
}
