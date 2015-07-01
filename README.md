## Foreword

This project is mainly a fork of respective functionality originally written by Richard Hull in python: <https://github.com/rm-hull/max7219>. Newetheless it differs in some parts: refuse some functionality (works only with matrix led), include extra functionality (extra fonts, support of national languages).

## MAX7219 driver

This library intended to output text messages to 8x8 LED display (pdf reference) interacting with MAX7219 driver chip (pdf reference):

Today there are many master devices which can drive MAX7219 chip, but this lib intended to work on Raspberry PI and its clones (tested with Raspberry PI and Banana PI). It may works with any Raspberry PI clone, which support Kernel SPI bus, and you should carry out all necessary preparations to make SPI bus device present in /dev/ list.

## Golang usage

```go
func main() {
  dev.SlideMessage("Hello world !!!", font, true, 50*time.Millisecond)
}
```

## Dependencies

Import and use package [github.com/fulr/spidev](http://github.com/fulr/spidev) to interact with RPi via Linux SPI device.

## Documentation

GoDoc [documentation](http://godoc.org/github.com/d2r2/go-max7219/max7219)

## Installation

```bash
$ go get github.com/d2r2/go-max7219/max7219
```

## Quick Start

To output a single letter to LED matrix by specifing ascii code use OutputAsciiCode call:
```go
	// Output a sequence of ascii codes in a loop
	font = max7219.FontCP437
	for i := 0; i <= len(font.GetLetterPatterns()); i++ {
		mtx.OutputAsciiCode(0, font, i, true)
		time.Sleep(500 * time.Millisecond)
	}
```
To output a single national letter either unicode letter (rune) to LED matrix use OutputChar call:
```go
	// Output non-latin national letter (russian for example).
	// You must be sure, that your national letter should match code page of font used.
	mtx.OutputChar(0, max7219.FontZXSpectrumRus, 'Я', true)
```

## FAQ

## License

Go-max7219 is licensed inder MIT License.
