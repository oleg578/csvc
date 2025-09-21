package csvc

import (
	"bufio"
	"io"
)

const (
	ASCII_LF    = 0x0A // '\n'
	ASCII_CR    = 0x0D // '\r'
	ASCII_DQ    = 0x22 // '"'
	ASCII_COMMA = 0x2C // ','
	ASCII_TAB   = 0x09 // '\t'
)

type Reader struct {
	Comma byte

	r        *bufio.Reader
	fieldBuf []byte // reusable buffer for building fields
}

func NewReader(r *bufio.Reader) *Reader {
	return &Reader{
		Comma:    ',',
		r:        r,
		fieldBuf: make([]byte, 0, 256), // initial capacity for field buffer
	}
}

func (b *Reader) Read() (dst []string, err error) {
	var fields []string
	var inQuotes bool

	// Pre-allocate slice with reasonable capacity to reduce reallocations
	if cap(fields) < 8 {
		fields = make([]string, 0, 8)
	}

	// Reset field buffer but keep capacity
	b.fieldBuf = b.fieldBuf[:0]

	for {
		ch, err := b.r.ReadByte()
		if err != nil {
			if err == io.EOF && len(b.fieldBuf) > 0 {
				// Handle last field if we have content
				fields = append(fields, string(b.fieldBuf))
			}
			return fields, err
		}

		switch ch {
		case ASCII_DQ: // Double quote
			if inQuotes {
				// Check if this is an escaped quote (double quote)
				nextCh, err := b.r.ReadByte()
				if err == nil && nextCh == ASCII_DQ {
					// Escaped quote - add single quote to field
					b.fieldBuf = append(b.fieldBuf, ASCII_DQ)
				} else {
					// End of quoted field - put back the character if not EOF
					if err == nil {
						b.r.UnreadByte()
					}
					inQuotes = false
				}
			} else {
				// Start of quoted field
				inQuotes = true
			}

		case b.Comma: // Field separator
			if inQuotes {
				// Comma inside quotes is part of the field
				b.fieldBuf = append(b.fieldBuf, ch)
			} else {
				// End of field - use fieldBuf directly to avoid string copying
				fields = append(fields, string(b.fieldBuf))
				b.fieldBuf = b.fieldBuf[:0] // reset but keep capacity
			}

		case ASCII_LF: // Line feed
			if inQuotes {
				// LF inside quotes is part of the field
				b.fieldBuf = append(b.fieldBuf, ch)
			} else {
				// End of record - add the last field and return
				fields = append(fields, string(b.fieldBuf))
				return fields, nil
			}

		case ASCII_CR: // Carriage return
			if inQuotes {
				// CR inside quotes is part of the field
				b.fieldBuf = append(b.fieldBuf, ch)
			} else {
				// Check if followed by LF for CRLF sequence
				nextCh, err := b.r.ReadByte()
				if err == nil && nextCh == ASCII_LF {
					// CRLF - end of record
					fields = append(fields, string(b.fieldBuf))
					return fields, nil
				} else {
					// Just CR - treat as regular character
					b.fieldBuf = append(b.fieldBuf, ch)
					if err == nil {
						b.r.UnreadByte()
					}
				}
			}

		default:
			// Regular character - add to current field
			b.fieldBuf = append(b.fieldBuf, ch)
		}
	}
}
