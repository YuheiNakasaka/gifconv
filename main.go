package main

import (
	"fmt"
	"os"

	"github.com/YuheiNakasaka/gifconv/gif"
)

func main() {
	file, err := os.Open("test.gif")
	if err != nil {
		fmt.Println("Error opening file")
	}
	defer file.Close()
	gif.FramePath = "./frames/"
	gif.Split(file)
	gif.EncodeAll()
}
