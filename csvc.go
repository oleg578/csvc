package csvc

import (
	"bufio"
	"bytes"
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

	r *bufio.Reader
}

func NewReader(r *bufio.Reader) *Reader {
	return &Reader{
		Comma: ',',
		r:     bufio.NewReader(r),
	}
}

func (b *Reader) Read() (dst []string, err error) {
	var fields []string
	var field bytes.Buffer
	var inQuotes bool

	for {
		ch, err := b.r.ReadByte()
		if err != nil {
			if err == io.EOF && field.Len() > 0 {
				// Handle last field if we have content
				fields = append(fields, field.String())
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
					field.WriteByte(ASCII_DQ)
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
				field.WriteByte(ch)
			} else {
				// End of field
				fields = append(fields, field.String())
				field.Reset()
			}

		case ASCII_LF: // Line feed
			if inQuotes {
				// LF inside quotes is part of the field
				field.WriteByte(ch)
			} else {
				// End of record - add the last field and return
				fields = append(fields, field.String())
				return fields, nil
			}

		case ASCII_CR: // Carriage return
			if inQuotes {
				// CR inside quotes is part of the field
				field.WriteByte(ch)
			} else {
				// Check if followed by LF for CRLF sequence
				nextCh, err := b.r.ReadByte()
				if err == nil && nextCh == ASCII_LF {
					// CRLF - end of record
					fields = append(fields, field.String())
					return fields, nil
				} else {
					// Just CR - treat as regular character
					field.WriteByte(ch)
					if err == nil {
						b.r.UnreadByte()
					}
				}
			}

		default:
			// Regular character - add to current field
			field.WriteByte(ch)
		}
	}
}
