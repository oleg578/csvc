package csvc

import (
	"bufio"
	"bytes"
	"io"
)

// Constants for CSV parsing
const (
	ASCII_COMMA           byte = 0x2C // Comma character
	ASCII_DOUBLE_QUOTE    byte = 0x22 // Double quote character
	ASCII_LINE_FEED       byte = 0x0A // Line feed character
	ASCII_CARRIAGE_RETURN byte = 0x0D // Carriage return character
)

// Reader represents a CSV reader
type Reader struct {
	Comma      byte
	FieldCount int // number of fields per record
	r          *bufio.Reader
	buf        []byte
	fields     []string
	out        []byte
	starts     []int
	ends       []int
}

func NewReader(r io.Reader) *Reader {
	return &Reader{
		Comma:  ASCII_COMMA,
		r:      bufio.NewReaderSize(r, 8192), // Larger buffer for better performance
		buf:    make([]byte, 0, 8192),
		fields: make([]string, 0, 32),
		out:    make([]byte, 0, 2048),
		starts: make([]int, 0, 32),
		ends:   make([]int, 0, 32),
	}
}

func (b *Reader) Read() (dst []string, err error) {
	// Reset reusable slices but keep capacity
	b.buf = b.buf[:0]
	b.fields = b.fields[:0]
	b.out = b.out[:0]
	b.starts = b.starts[:0]
	b.ends = b.ends[:0]

	if err = b.readRecord(); err != nil {
		return
	}
	// Empty line should return one empty field, not EOF
	// Only return EOF if there's no data at all
	if len(b.buf) == 0 && err == io.EOF {
		return nil, io.EOF
	}

	// Pre-size output buffer to current line size to reduce reallocations in quoted path
	if cap(b.out) < len(b.buf) {
		b.out = make([]byte, 0, len(b.buf))
	}

	fieldCount := 0

	// Fast path: no quotes present, split directly from input buffer
	end := len(b.buf)
	noQuotes := true
	for i := 0; i < end; i++ {
		if b.buf[i] == ASCII_DOUBLE_QUOTE {
			noQuotes = false
			break
		}
	}
	if noQuotes {
		// Collect ranges, build one string, then slice
		b.starts = b.starts[:0]
		b.ends = b.ends[:0]
		start := 0
		for i := 0; i < end; i++ {
			if b.buf[i] == b.Comma {
				b.starts = append(b.starts, start)
				b.ends = append(b.ends, i)
				fieldCount++
				start = i + 1
			}
		}
		// append last field (handles trailing comma as empty)
		b.starts = append(b.starts, start)
		b.ends = append(b.ends, end)
		fieldCount++

		lineStr := string(b.buf[:end])
		for i := 0; i < len(b.starts); i++ {
			s, e := b.starts[i], b.ends[i]
			b.fields = append(b.fields, lineStr[s:e])
		}

		if b.FieldCount == 0 {
			b.FieldCount = fieldCount
		}
		if fieldCount < b.FieldCount {
			for i := fieldCount; i < b.FieldCount; i++ {
				b.fields = append(b.fields, "")
			}
		}

		return b.fields, nil
	}

	// RFC 4180 State Machine for CSV parsing
	const (
		STATE_START_FIELD     = iota // Beginning of field
		STATE_UNQUOTED               // Inside unquoted field
		STATE_QUOTED                 // Inside quoted field
		STATE_QUOTE_IN_QUOTED        // Quote character within quoted field
		STATE_END_FIELD              // End of field
	)

	state := STATE_START_FIELD
	lastWasComma := false // track if last processed character was a comma
	fieldStart := -1      // start index in out for the current field

	for i := 0; i < len(b.buf); i++ {
		ch := b.buf[i]

		lastWasComma = false // reset unless we see a comma

		switch state {
		case STATE_START_FIELD:
			switch ch {
			case ASCII_DOUBLE_QUOTE:
				// Start of quoted field
				fieldStart = len(b.out)
				state = STATE_QUOTED
			case b.Comma:
				// Empty field
				pos := len(b.out)
				b.starts = append(b.starts, pos)
				b.ends = append(b.ends, pos)
				fieldCount++
				lastWasComma = true
				state = STATE_START_FIELD
			default:
				// Start of unquoted field
				fieldStart = len(b.out)
				b.out = append(b.out, ch)
				state = STATE_UNQUOTED
			}

		case STATE_UNQUOTED:
			if ch == b.Comma {
				// End of unquoted field
				b.starts = append(b.starts, fieldStart)
				b.ends = append(b.ends, len(b.out))
				fieldCount++
				fieldStart = -1
				lastWasComma = true
				state = STATE_START_FIELD
			} else {
				// Continue unquoted field
				b.out = append(b.out, ch)
			}

		case STATE_QUOTED:
			// Jump to next quote to avoid per-byte loop
			if ch != ASCII_DOUBLE_QUOTE {
				if idx := bytes.IndexByte(b.buf[i:], ASCII_DOUBLE_QUOTE); idx >= 0 {
					// copy run up to quote
					b.out = append(b.out, b.buf[i:i+idx]...)
					i += idx
					ch = b.buf[i]
					// fallthrough to handle quote below
				} else {
					// no more quotes: append rest and finish later
					b.out = append(b.out, b.buf[i:]...)
					i = len(b.buf) - 1
					break
				}
			}
			// ch is a quote here
			state = STATE_QUOTE_IN_QUOTED

		case STATE_QUOTE_IN_QUOTED:
			switch ch {
			case ASCII_DOUBLE_QUOTE:
				// Escaped quote (two consecutive quotes = one quote)
				b.out = append(b.out, ASCII_DOUBLE_QUOTE)
				state = STATE_QUOTED
			case b.Comma:
				// End of quoted field
				b.starts = append(b.starts, fieldStart)
				b.ends = append(b.ends, len(b.out))
				fieldCount++
				fieldStart = -1 // reset for next field
				lastWasComma = true
				state = STATE_START_FIELD
			default:
				// End quoted field and start a new unquoted field with this character
				b.starts = append(b.starts, fieldStart)
				b.ends = append(b.ends, len(b.out))
				fieldCount++
				fieldStart = len(b.out)
				b.out = append(b.out, ch)
				state = STATE_UNQUOTED
			}
		}
	}

	// Handle the last field
	if fieldStart != -1 && (state == STATE_QUOTED || state == STATE_QUOTE_IN_QUOTED || state == STATE_UNQUOTED) {
		// finalize last field
		b.starts = append(b.starts, fieldStart)
		b.ends = append(b.ends, len(b.out))
		fieldCount++
	} else if lastWasComma {
		// If line ended with a comma, add empty field
		pos := len(b.out)
		b.starts = append(b.starts, pos)
		b.ends = append(b.ends, pos)
		fieldCount++
	}

	// Build final string once and slice for each field
	if len(b.starts) > 0 {
		lineStr := string(b.out)
		for i := 0; i < len(b.starts); i++ {
			s, e := b.starts[i], b.ends[i]
			b.fields = append(b.fields, lineStr[s:e])
		}
	}
	// Don't add empty field at end for START_FIELD state unless lastWasComma

	if b.FieldCount == 0 {
		b.FieldCount = fieldCount
	}
	if fieldCount < b.FieldCount {
		for i := fieldCount; i < b.FieldCount; i++ {
			b.fields = append(b.fields, "")
		}
	}

	return b.fields, nil
}

// readLine removed in favor of readRecord which supports multi-line quoted fields

// readRecord reads a full CSV record possibly spanning multiple physical lines.
// It accumulates data until a line terminator is encountered while not inside a quoted field.
func (b *Reader) readRecord() error {
	b.buf = b.buf[:0]
	insideQuoted := false
	for {
		chunk, err := b.r.ReadSlice('\n')
		if len(chunk) > 0 {
			prevLen := len(b.buf)
			b.buf = append(b.buf, chunk...)
			// Scan from prevLen-1 to catch escaped quote pairs crossing chunk boundary
			start := prevLen
			if start > 0 {
				start = start - 1
			}
			for i := start; i < len(b.buf); i++ {
				c := b.buf[i]
				if c == ASCII_DOUBLE_QUOTE {
					if insideQuoted {
						if i+1 < len(b.buf) && b.buf[i+1] == ASCII_DOUBLE_QUOTE {
							i++ // skip escaped quote
						} else {
							insideQuoted = false
						}
					} else {
						insideQuoted = true
					}
				}
			}
		}
		if err == nil {
			// Newline encountered; if not inside quoted, record ends here
			if !insideQuoted {
				break
			}
			// Else continue to accumulate next line
			continue
		}
		if err == bufio.ErrBufferFull {
			// Continue reading next chunk; no newline encountered yet.
			continue
		}
		if err == io.EOF {
			// EOF
			if len(b.buf) == 0 {
				// no data at all
				return io.EOF
			}
			// return partial record as-is (may be unterminated quoted field)
			break
		}
		// Any other error
		return err
	}
	// Trim trailing record delimiter: CRLF, LF, or CR
	end := len(b.buf)
	if end > 0 && b.buf[end-1] == ASCII_LINE_FEED {
		end--
		if end > 0 && b.buf[end-1] == ASCII_CARRIAGE_RETURN {
			end--
		}
	} else if end > 0 && b.buf[end-1] == ASCII_CARRIAGE_RETURN {
		end--
	}
	b.buf = b.buf[:end]
	return nil
}
