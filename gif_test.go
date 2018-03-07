package gifconv

import (
	"os"
	"testing"
)

func TestSplit(t *testing.T) {
	file, _ := os.Open("test.gif")
	defer file.Close()

	FramePath = "./frames/"
	framePaths, err := Split(file)
	if err != nil {
		t.Fatal("Return Split err")
	}

	if len(framePaths) != 2 {
		t.Fatal("Return invalid total frame path count")
	}
}

func TestEncodeAll(t *testing.T) {
	file, _ := os.Open("test.gif")
	defer file.Close()

	FramePath = "./frames/"
	framePaths, _ := Split(file)

	err := EncodeAll(framePaths, 10)
	if err != nil {
		t.Fatal("EncodeAll return err")
	}

	if _, err = os.Stat(FramePath + "output.gif"); os.IsNotExist(err) {
		t.Fatal("File not created")
	}
}
