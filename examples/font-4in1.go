package main

import (
	"log"
	"time"

	"github.com/adrianh-za/go-max7219"
)

func main() {

	// No rotate - Show 4 letters
	mtx := max7219.NewMatrix(4, max7219.RotateNone)
	err := mtx.Open(0, 0, 1)
	if err != nil {
		log.Fatal(err)
	}

	mtx.OutputChar(0, max7219.FontZXSpectrumRus, 'Я', true)
	mtx.OutputChar(1, max7219.FontZXSpectrumRus, 'Я', true)
	mtx.OutputChar(2, max7219.FontZXSpectrumRus, 'Я', true)
	mtx.OutputChar(3, max7219.FontZXSpectrumRus, 'Я', true)
	time.Sleep(2 * time.Second)
	mtx.Clear()
	mtx.Close()

	// Clockwise - Show 4 letters
	mtx = max7219.NewMatrix(4, max7219.RotateClockwise)
	err = mtx.Open(0, 0, 1)
	if err != nil {
		log.Fatal(err)
	}

	mtx.OutputChar(0, max7219.FontZXSpectrumRus, 'Я', true)
	mtx.OutputChar(1, max7219.FontZXSpectrumRus, 'Я', true)
	mtx.OutputChar(2, max7219.FontZXSpectrumRus, 'Я', true)
	mtx.OutputChar(3, max7219.FontZXSpectrumRus, 'Я', true)
	time.Sleep(2 * time.Second)
	mtx.Clear()
	mtx.Close()

	// AntiClockwise - Show 4 letters
	mtx = max7219.NewMatrix(4, max7219.RotateAntiClockwise)
	err = mtx.Open(0, 0, 1)
	if err != nil {
		log.Fatal(err)
	}

	mtx.OutputChar(0, max7219.FontZXSpectrumRus, 'Я', true)
	mtx.OutputChar(1, max7219.FontZXSpectrumRus, 'Я', true)
	mtx.OutputChar(2, max7219.FontZXSpectrumRus, 'Я', true)
	mtx.OutputChar(3, max7219.FontZXSpectrumRus, 'Я', true)
	time.Sleep(2 * time.Second)
	mtx.Clear()
	mtx.Close()
}