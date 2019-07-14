package main

import (
	"log"
	"time"

	"github.com/adrianh-za/go-max7219"
)

func main() {

	mtx := max7219.NewMatrix(4, max7219.RotateAntiClockwiseReverse)
	err := mtx.Open(0, 0, 1)
	if err != nil {
		log.Fatal(err)
	}

	mtx.SlideMessage("ANTI-CLOCKWISE slide text!", max7219.FontCP437, true, 50 * time.Millisecond)
	mtx.Close()

	time.Sleep(1 * time.Second)

	mtx = max7219.NewMatrix(4, max7219.RotateClockwiseReverse)
	err = mtx.Open(0, 0, 1)
	if err != nil {
		log.Fatal(err)
	}

	mtx.SlideMessage("CLOCKWISE slide text!", max7219.FontCP437, true, 50 * time.Millisecond)
	mtx.Close()

}