package csvc

import (
	"bufio"
	"io"
)

// Constants for CSV parsing
const (
	ASCII_COMMA           byte = 0x2C // Comma character
	ASCII_DOUBLE_QUOTE    byte = 0x22 // Double quote character
	ASCII_LINE_FEED       byte = 0x0A // Line feed character
	ASCII_CARRIAGE_RETURN byte = 0x0D // Carriage return character
	ASCII_TAB             byte = 0x09 // Tab character
)

// Reader represents a CSV reader
type Reader struct {
	Comma      byte
	FieldCount int // number of fields per record
	r          *bufio.Reader
	buf        []byte // buffer to hold data]
}

func NewReader(r io.Reader) *Reader {
	return &Reader{
		Comma: ASCII_COMMA,
		r:     bufio.NewReaderSize(r, 4096),
		buf:   make([]byte, 0, 4096), // bytes buffer with capacity of 4096 bytes
	}
}

func (b *Reader) Read() (dst []string, err error) {
	defer func() { b.buf = b.buf[:0] }() // reset buffer after reading
	if err = b.readLine(); err != nil {
		return
	}
	// Empty line should return one empty field, not EOF
	// Only return EOF if there's no data at all
	if len(b.buf) == 0 && err == io.EOF {
		return nil, io.EOF
	}

	fieldCount := 0
	field := make([]byte, 0, 256)

	// RFC 4180 State Machine for CSV parsing
	const (
		STATE_START_FIELD     = iota // Beginning of field
		STATE_UNQUOTED               // Inside unquoted field
		STATE_QUOTED                 // Inside quoted field
		STATE_QUOTE_IN_QUOTED        // Quote character within quoted field
		STATE_END_FIELD              // End of field
	)

	state := STATE_START_FIELD

	for i := 0; i < len(b.buf); i++ {
		ch := b.buf[i]

		// Skip line ending characters that might be included in buffer
		if ch == ASCII_CARRIAGE_RETURN || ch == ASCII_LINE_FEED {
			continue
		}

		switch state {
		case STATE_START_FIELD:
			switch ch {
			case ASCII_DOUBLE_QUOTE:
				// Start of quoted field
				state = STATE_QUOTED
			case b.Comma:
				// Empty field
				dst = append(dst, "")
				fieldCount++
				field = field[:0] // reset field buffer
				state = STATE_START_FIELD
			default:
				// Start of unquoted field
				field = append(field, ch)
				state = STATE_UNQUOTED
			}

		case STATE_UNQUOTED:
			if ch == b.Comma {
				// End of unquoted field
				dst = append(dst, string(field))
				fieldCount++
				field = field[:0] // reset field buffer
				state = STATE_START_FIELD
			} else {
				// Continue unquoted field
				field = append(field, ch)
			}

		case STATE_QUOTED:
			if ch == ASCII_DOUBLE_QUOTE {
				// Potential end of quoted field or escaped quote
				state = STATE_QUOTE_IN_QUOTED
			} else {
				// Continue quoted field (including commas, newlines, etc.)
				field = append(field, ch)
			}

		case STATE_QUOTE_IN_QUOTED:
			switch ch {
			case ASCII_DOUBLE_QUOTE:
				// Escaped quote (two consecutive quotes = one quote)
				field = append(field, ASCII_DOUBLE_QUOTE)
				state = STATE_QUOTED
			case b.Comma:
				// End of quoted field
				dst = append(dst, string(field))
				fieldCount++
				field = field[:0] // reset field buffer
				state = STATE_START_FIELD
			default:
				// This should not happen in valid CSV, but we'll treat it as end of quoted field
				// and start processing the character as a new field
				dst = append(dst, string(field))
				fieldCount++
				field = field[:0] // reset field buffer
				field = append(field, ch)
				state = STATE_UNQUOTED
			}
		}
	}

	// Handle the last field
	if state == STATE_QUOTED || state == STATE_QUOTE_IN_QUOTED {
		// Unterminated quoted field - this is technically an error in RFC 4180
		// but we'll be lenient and include the field
		dst = append(dst, string(field))
		fieldCount++
	} else if len(field) > 0 || state == STATE_START_FIELD {
		// Last field or empty field at end
		dst = append(dst, string(field))
		fieldCount++
	}
	if b.FieldCount == 0 {
		b.FieldCount = fieldCount
	}
	if fieldCount < b.FieldCount {
		for i := fieldCount; i < b.FieldCount; i++ {
			dst = append(dst, "")
		}
	}
	return
}

// readLine reads a single line from the CSV input into b.buf.
func (b *Reader) readLine() (err error) {
	var isPrefix bool
	line := make([]byte, 0, 4096)
	for {
		line, isPrefix, err = b.r.ReadLine()
		if err != nil && err.Error() != "EOF" {
			return err
		}
		b.buf = append(b.buf, line...)
		if !isPrefix { // if we have read a complete line
			break
		}
	}
	return
}
