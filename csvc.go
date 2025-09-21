package csvc

import (
	"bufio"
)

// Constants for CSV parsing
const (
	ASCII_COMMA = ','  // Comma character
	ASCII_DQ    = '"'  // Double quote character
	ASCII_LF    = '\n' // Line feed character
	ASCII_CR    = '\r' // Carriage return character
	ASCII_TAB   = '\t' // Tab character
)

// Reader represents a CSV reader
type Reader struct {
	Comma byte
	r        *bufio.Reader
}

func NewReader(r *bufio.Reader) *Reader {
	return &Reader{
		Comma:    ',',
		r:        r,
	}
}

func (b *Reader) Read() (dst []string, err error) {
	return dst, err
}
