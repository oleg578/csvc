package csvc

import (
	"errors"
	"fmt"
	"io"
	"unsafe"
)

const defaultBufferSize = 64 * 1024

var (
	// ErrBareQuote is returned when an unexpected quote is found in an unquoted field.
	ErrBareQuote = errors.New("csvc: bare quote in non-quoted field")
	// ErrUnterminatedQuote is returned when a quoted field is not closed before EOF or record end.
	ErrUnterminatedQuote = errors.New("csvc: unterminated quoted field")
)

// ParseError contains location information for CSV parsing errors.
type ParseError struct {
	Line   int
	Column int
	Err    error
}

func (e *ParseError) Error() string {
	if e == nil {
		return ""
	}
	return fmt.Sprintf("csvc: parse error on line %d, column %d: %v", e.Line, e.Column, e.Err)
}

// Unwrap returns the underlying error.
func (e *ParseError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Err
}

// Reader provides high-performance CSV parsing with support for customizable delimiters.
type Reader struct {
	src io.Reader

	// Comma is the field delimiter. Default is ','.
	Comma byte
	// Quote is the quote character. Default is '"'.
	Quote byte
	// ReuseRecord indicates whether Read should reuse the backing array of the returned slice.
	ReuseRecord bool

	buf    []byte
	bufPos int
	bufLen int
	bufErr error

	record      []string
	dataBuf     []byte
	fieldBounds []int
	finished    bool
	line        int
}

// NewReader creates a new Reader with an internal buffered reader tuned for bulk throughput.
func NewReader(r *io.Reader) *Reader {
	if r == nil || *r == nil {
		panic("csvc: reader source cannot be nil")
	}

	return &Reader{
		src:         *r,
		Comma:       ',',
		Quote:       '"',
		buf:         make([]byte, defaultBufferSize),
		record:      make([]string, 0, 16),
		dataBuf:     make([]byte, 0, 512),
		fieldBounds: make([]int, 0, 32),
		line:        1,
	}
}

// Read parses the next CSV record. It reuses allocations when possible to provide
// higher throughput than encoding/csv. When no more records are available it returns io.EOF.
func (r *Reader) Read() (dst []string, err error) {
	if r == nil || r.src == nil {
		return nil, io.EOF
	}
	if r.finished {
		return nil, io.EOF
	}

	comma := r.Comma
	if comma == 0 {
		comma = ','
	}
	quote := r.Quote
	if quote == 0 {
		quote = '"'
	}

	if r.ReuseRecord {
		r.record = r.record[:0]
	} else {
		r.record = nil
	}
	r.dataBuf = r.dataBuf[:0]
	r.fieldBounds = r.fieldBounds[:0]

	inQuotes := false
	sawQuotedField := false
	column := 1
	fieldStart := 0

	for {
		b, readErr := r.readByte()
		curColumn := column
		if readErr != nil {
			if readErr == io.EOF {
				if inQuotes {
					r.finished = true
					return nil, r.wrapError(curColumn, ErrUnterminatedQuote)
				}
				if len(r.fieldBounds) > 0 || len(r.dataBuf) > 0 || sawQuotedField {
					r.fieldBounds = append(r.fieldBounds, fieldStart, len(r.dataBuf))
					r.finished = true
					return r.buildRecord(), nil
				}
				r.finished = true
				return nil, io.EOF
			}
			return nil, readErr
		}

		if inQuotes {
			if b == quote {
				next, err := r.peekByte()
				if err == nil && next == quote {
					r.bufPos++
					r.dataBuf = append(r.dataBuf, quote)
					column = curColumn + 2
					continue
				}
				if err != nil && err != io.EOF {
					return nil, err
				}
				inQuotes = false
				column = curColumn + 1
				continue
			}
			if b == '\n' {
				r.line++
				column = 1
			} else {
				column = curColumn + 1
			}
			r.dataBuf = append(r.dataBuf, b)
			continue
		}

		switch b {
		case comma:
			r.fieldBounds = append(r.fieldBounds, fieldStart, len(r.dataBuf))
			fieldStart = len(r.dataBuf)
			sawQuotedField = false
			column = curColumn + 1
		case '\n':
			r.fieldBounds = append(r.fieldBounds, fieldStart, len(r.dataBuf))
			fieldStart = len(r.dataBuf)
			sawQuotedField = false
			r.line++
			column = 1
			return r.buildRecord(), nil
		case '\r':
			next, err := r.peekByte()
			if err == nil && next == '\n' {
				r.bufPos++
			}
			if err != nil && err != io.EOF {
				return nil, err
			}
			r.fieldBounds = append(r.fieldBounds, fieldStart, len(r.dataBuf))
			fieldStart = len(r.dataBuf)
			sawQuotedField = false
			r.line++
			column = 1
			return r.buildRecord(), nil
		case quote:
			if len(r.dataBuf) == fieldStart && !sawQuotedField {
				inQuotes = true
				sawQuotedField = true
				column = curColumn + 1
				continue
			}
			return nil, r.wrapError(curColumn, ErrBareQuote)
		default:
			r.dataBuf = append(r.dataBuf, b)
			column = curColumn + 1
		}
	}
}

func (r *Reader) buildRecord() []string {
	fieldCount := len(r.fieldBounds) / 2

	var recordStr string
	if r.ReuseRecord {
		if len(r.dataBuf) == 0 {
			recordStr = ""
		} else {
			recordStr = unsafe.String(unsafe.SliceData(r.dataBuf), len(r.dataBuf))
		}
		if cap(r.record) < fieldCount {
			r.record = make([]string, fieldCount)
		}
		r.record = r.record[:fieldCount]
	} else {
		recordStr = string(r.dataBuf)
		r.record = make([]string, fieldCount)
	}

	for i := 0; i < fieldCount; i++ {
		start := r.fieldBounds[2*i]
		end := r.fieldBounds[2*i+1]
		r.record[i] = recordStr[start:end]
	}

	return r.record
}

func (r *Reader) wrapError(column int, err error) error {
	return &ParseError{Line: r.line, Column: column, Err: err}
}

func (r *Reader) readByte() (byte, error) {
	for {
		if r.bufPos < r.bufLen {
			b := r.buf[r.bufPos]
			r.bufPos++
			return b, nil
		}
		if r.bufErr != nil {
			err := r.bufErr
			r.bufErr = nil
			return 0, err
		}

		n, err := r.src.Read(r.buf)
		if n == 0 && err != nil {
			return 0, err
		}
		if n == 0 {
			continue
		}
		r.bufPos = 0
		r.bufLen = n
		r.bufErr = err
	}
}

func (r *Reader) peekByte() (byte, error) {
	for {
		if r.bufPos < r.bufLen {
			return r.buf[r.bufPos], nil
		}
		if r.bufErr != nil {
			return 0, r.bufErr
		}

		n, err := r.src.Read(r.buf)
		if n == 0 && err != nil {
			return 0, err
		}
		if n == 0 {
			continue
		}
		r.bufPos = 0
		r.bufLen = n
		r.bufErr = err
	}
}
