package max7219

import (
	"bytes"
	"fmt"
	"time"
	//"strings"

	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
)

type Rotation int

const (
	RotateNone					Rotation = 0	// No rotation needed (normally used when individual LED modules make up matrix)
	RotateClockwise				Rotation = 1	// Used to rotate 4 in 1 LED matrix clockwise
	RotateAntiClockwise			Rotation = 2	// Used to rotate 4 in 1 LED matrix anti-clockwise
	RotateClockwiseInvert 		Rotation = 3	// Used to rotate 4 in 1 LED matrix clockwise (sometimes the modules are 180 degrees to the other 4 in 1 modules)
	RotateAntiClockwiseInvert	Rotation = 4	// Used to rotate 4 in 1 LED matrix anti-clockwise (sometimes the modules are 180 degrees to the other 4 in 1 modules)
 )

type Matrix struct {
	Device *Device
	Rotation Rotation
}

func NewMatrix(cascaded int, rotate Rotation) *Matrix {
	this := &Matrix{}
	this.Device = NewDevice(cascaded)
	this.Rotation = rotate
	return this
}

func (this *Matrix) Open(spibus int, spidevice int, brightness byte) error {
	return this.Device.Open(spibus, spidevice, brightness)
}

func (this *Matrix) Close() {
	this.Device.Close()
}

func (this *Matrix) Clear() {
	for count := 0; count < this.Device.GetCascadeCount(); count++ {
		this.OutputChar(count, FontCP437,' ', true)
	}
}

func getLineCondense(line byte) int {
	var condense int = 0
	var i uint
	for i = 0; i < 8; i++ {
		if line&(1<<i) > 0 {
			condense += 1
		}
	}
	return condense
}

func getLetterPatternLimits(pattern []byte) (start int, end int) {
	startIndex := -1
	endIndex := -1
	for i := 0; i < len(pattern); i++ {
		if pattern[i] != 0 && startIndex == -1 {
			startIndex = i
		}
	}
	if startIndex == -1 {
		startIndex = 0
	}
	for i := len(pattern) - 1; i >= 0; i-- {
		if pattern[i] != 0 && endIndex == -1 {
			endIndex = i
		}
	}
	if endIndex == -1 {
		endIndex = len(pattern) - 1
	}
	return startIndex, endIndex
}

func preparePatterns(text []byte, font [][]byte,
	condenseLetterPattern bool) []byte {
	var patrns [][]byte
	var limits [][]int
	totalWidth := 0
	for _, c := range text {
		pattern := font[c]
		start, end := getLetterPatternLimits(pattern)
		totalWidth += end - start + 1
		patrns = append(patrns, pattern)
		limits = append(limits, []int{start, end})
	}
	averageWidth := totalWidth / len(text)
	log.Debug("Average width: %d\n", averageWidth)
	var buf []byte
	for i := 0; i < len(patrns); i++ {
		if condenseLetterPattern {
			var startC = getLineCondense(patrns[i][limits[i][0]])
			var endC int = 0
			if i > 0 {
				endC = getLineCondense(patrns[i-1][limits[i-1][1]])
			}
			// In case of space char...
			if isEmpty(patrns[i]) {
				// ... specify average char width + extra line.
				limits[i][1] = averageWidth - 1 - 1
			}
			// ... + extra lines.
			if endC+startC == 0 {
			} else if endC+startC <= 2 || i == 0 {
				buf = append(buf, 0)
			} else if endC+startC <= 10 {
				buf = append(buf, 0, 0)
			} else {
				buf = append(buf, 0, 0, 0)
			}
			buf = append(buf, patrns[i][limits[i][0]:limits[i][1]+1]...)
		} else {
			buf = append(buf, patrns[i]...)
		}
	}
	if condenseLetterPattern {
		buf = append(buf, 0)
	}
	return buf
}

func repeat(b byte, count int) []byte {
	buf := make([]byte, count)
	for i := 0; i < len(buf); i++ {
		buf[i] = b
	}
	return buf
}

func isEmpty(pattern []byte) bool {
	for _, b := range pattern {
		if b != 0 {
			return false
		}
	}
	return true
}

// Output unicode char to the led matrix.
// Unicode char transforms to ascii code based on
// information taken from font.GetCodePage() call.
func (this *Matrix) OutputChar(cascadeId int, font Font,
	char rune, redraw bool) error {
	text := string(char)
	b := convertUnicodeToAscii(text, font.GetCodePage())
	if len(b) != 1 {
		return fmt.Errorf("One char expected: \"%s\"", text)
	}
	buf := preparePatterns(b, font.GetLetterPatterns(), false)
	buf = rotateCharacter(buf, this.Rotation)
	for i, value := range buf {
		//fmt.Printf("value: %#x\n", value)
		err := this.Device.SetBufferLine(cascadeId, i, value, redraw)
		if err != nil {
			return err
		}
	}
	return nil
}

// Output ascii code to the led matrix.
func (this *Matrix) OutputAsciiCode(cascadeId int, font Font,
	asciiCode int, redraw bool) error {
	patterns := font.GetLetterPatterns()
	buf := patterns[asciiCode]
	buf = rotateCharacter(buf, this.Rotation)

	for i, value := range buf {
		//fmt.Printf("value: %#x\n", value)
		err := this.Device.SetBufferLine(cascadeId, i, value, redraw)
		if err != nil {
			return err
		}
	}
	return nil
}

// Convert unicode text to ASCII text
// using specific codepage mapping.
func convertUnicodeToAscii(text string,
	codepage encoding.Encoding) []byte {
	b := []byte(text)
	// fmt.Printf("Text length: %d\n", len(b))
	var buf bytes.Buffer
	if codepage == nil {
		codepage = charmap.Windows1252
	}
	w := transform.NewWriter(&buf, codepage.NewEncoder())
	defer w.Close()
	w.Write(b)
	// fmt.Printf("Buffer length: %d\n", len(buf.Bytes()))
	return buf.Bytes()
}

// Show message sliding it by led matrix from the right to left.
func (this *Matrix) SlideMessage(text string, font Font, condensePattern bool, pixelDelay time.Duration) error {
	
	// No rotation configured, do normal scrolling
	if (this.Rotation == RotateNone) {

		b := convertUnicodeToAscii(text, font.GetCodePage())
		buffer := preparePatterns(b, font.GetLetterPatterns(), condensePattern)

		for _, b := range buffer {
			time.Sleep(pixelDelay)
			err := this.Device.ScrollLeft(true)
			if err != nil {
				return err
			}
			err = this.Device.SetBufferLine(
				this.Device.GetCascadeCount()-1,
				this.Device.GetLedLineCount()-1, b, true)
			if err != nil {
				return err
			}
		}
	
	// Rotation configured, do "special" scrolling
	} else {
		
		b := convertUnicodeToAscii(text, font.GetCodePage())
		buffer := preparePatterns(b, font.GetLetterPatterns(), condensePattern)
		
		// Add empty chars onto start of buffer so it can start sliding "off-screen" (8 bytes per LED module)
		buffer = append(make([]byte, this.Device.GetCascadeCount() * 8), buffer...)  
		
		//Start sliding
		var shiftCount = 0
		for { 
			this.Device.Flush()
			
			// End of buffer reached
			if shiftCount > len(buffer) {
				break
			}

			// Build up the buffers to be displayed
			var shiftBuffer = append(buffer[shiftCount:])		// Shift the buffer 1 vertical line of LEDs to left
			var paddingBuffer = make([]byte, shiftCount + 1)	// Add empty char onto end of buffer so it can continue sliding "off-screen"
			shiftBuffer = append(shiftBuffer, paddingBuffer...) // Concatenate the 2 buffers
			
			// Rotate the letters in the buffer
			var rotatedBuffer = rotateCharacters(shiftBuffer, this.Rotation)

			// Set the start LED block (cascade)
			var cascadeCount = 0
			if (this.Rotation == RotateAntiClockwise) || (this.Rotation == RotateAntiClockwiseInvert) {
				cascadeCount = this.Device.GetCascadeCount() - 1
			}

			// Render each LED line accross all the LED blocks (cascades)
			for lineCount, byteValue := range rotatedBuffer {
				err := this.Device.SetBufferLine(
					cascadeCount,
					(lineCount % 8),
					byteValue,
					true)
				if err != nil {
					return err
				}

				// Move onto next LED block (cascade)
				if (lineCount % 8 == 7) {
					if (this.Rotation == RotateClockwise) || (this.Rotation == RotateClockwiseInvert) {
						cascadeCount++
					} else if (this.Rotation == RotateAntiClockwise) || (this.Rotation == RotateAntiClockwiseInvert){
						cascadeCount--
					}
				}

				// All LED lines for all LED blocks (cascades) rendered, we can break
				if (lineCount == (this.Device.GetCascadeCount() * this.Device.GetLedLineCount()) - 1) {
					break
				}
			}
			
			// Delay for the specified time
			time.Sleep(pixelDelay)
			shiftCount++;
		}
	}

	return nil
}

// Rotate the set of characters based on the Rotation specified, ultimately calls rotateCharacter()
func rotateCharacters(character []byte, rotate Rotation) ([]byte) {
	var totalCharacters = len(character) / 8
	var rotatedCharacters = make([]byte, len(character))

	for count := 0; count < totalCharacters; count++ {
		start := (count * 8)
		end := (count * 8) + 8
		characterBytes := rotateCharacter(character[start:end], rotate)
		for i, byteValue := range characterBytes {
			rotatedCharacters[(count * 8) + i] = byteValue
		}
	}

	return rotatedCharacters
}

// Rotate the specified character based on the Rotation specified
func rotateCharacter(character []byte, rotate Rotation) ([]byte) {
	var tempByte = []byte{0,0,0,0,0,0,0,0}

	//No rotation, just return passed in character bit matrix
	if (rotate == RotateNone) {
		return character
	}

	// Rotate the character bit matrix anti-clockwise (minus 90 degrees)
	if (rotate == RotateAntiClockwise) || (rotate == RotateClockwiseInvert) {
		
		var counter = 7
		for outer := byte(0); outer < 8; outer++ {
			var byteValue = character[outer]
			for inner := byte(0); inner < 8; inner++ {
				var bitValue = (byteValue >> inner) & 1
				if (bitValue == 1) {
					tempByte[inner] = setBit(tempByte[inner], byte(counter))
				}
			}
			counter--
		}

		return tempByte
	}

	// Rotate the character bit matrix clockwise (plus 90 degrees)
	if (rotate == RotateClockwise) || (rotate == RotateAntiClockwiseInvert) {
		
		for outer := byte(0); outer < 8; outer++ {
			var byteValue = character[outer]
			var counter = 7
			for inner := byte(0); inner < 8; inner++ {
				var bitValue = (byteValue >> inner) & 1
				if (bitValue == 1) {
					tempByte[counter] = setBit(tempByte[counter], outer)
				}
				counter--
			}
		}

		return tempByte
	}

	return nil
}

// Sets the bit at position in the byte value.
func setBit(value byte, position byte) byte {
    value |= (1 << position)
    return value
}