# Go-Max7219

Library written in Go to allow controlling of the MAX7219LED module.  Source forked from https://github.com/d2r2/go-max7219

## Enhancements ##

* Works with pre-assembled 4-in-1 MAX7219 LED module matrices.
* Set the rotational direction of the MAX7219 LED modules.
    * Some of the pre-assembled 4-in-1 LED modules are connected "upsidedown" to the circuit board.  Using the <i>reversed</i> rotational directions allows for supporting of these MAX7219 4-in-1 LED module matrices.
* Sliding of text with blank padding before and after text

## Usage ##

1) go get github.com/adrianh-za/go-max7219
2) browse to $/go/src/github.com/adrianh-za/go-max7219/examples
3) sudo -E go run [filename].go
4) run to end (or ctrl-c to quit)

Examples filenames
* chars.go
* slide.go
* slide-4in1.go
* slide-4in1-reverse.go

## Compatibility ##

Tested on Raspberry PI 3 B+

## Acknowledgements ##

Thanks to <a href="https://github.com/d2r2" target="blank"><b>Denis Dyakov</b></a> for his excellent libraries

## Gits ##

https://github.com/d2r2/go-max7219
