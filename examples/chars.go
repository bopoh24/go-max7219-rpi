package main

import (
	"log"
	"time"

	"github.com/adrianh-za/go-max7219"
)

func main() {

	// Output a sequence of ascii codes in a loop
	mtx := max7219.NewMatrix(1, max7219.RotateNone)
	err := mtx.Open(0, 0, 1)
	if err != nil {
		log.Fatal(err)
	}
	
	for i := 0; i < len(max7219.FontCP437.GetLetterPatterns()); i++ {
		mtx.OutputAsciiCode(0, max7219.FontCP437, i, true)
		time.Sleep(200 * time.Millisecond)
	}
	time.Sleep(1 * time.Second)
	mtx.Clear()
	mtx.Close()
}