package main

import (
	"log"
	"time"

	"github.com/adrianh-za/go-max7219"
)

func main() {

	mtx := max7219.NewMatrix(4, max7219.RotateNone)
	err := mtx.Open(0, 0, 1)
	if err != nil {
		log.Fatal(err)
	}

	mtx.SlideMessage("NO ROTATION slide text!", max7219.FontCP437, true, 50 * time.Millisecond)
	time.Sleep(2 * time.Second)
	mtx.Clear()
	mtx.Close()
}
