package gifconv

import "testing"

func TestInfo(t *testing.T) {
	err := Info("test.gif")
	if err != nil {
		t.Fatal("Info return err")
	}
}
