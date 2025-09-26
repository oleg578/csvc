// Package csvc provides a high-performance CSV reader and writer implementation
// that is faster than the standard library's encoding/csv package.
package csvc

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"strings"
)

const (
	defaultDelimiter     = ','
	defaultComment       = '#'
	defaultLazyQuotes   = false
	defaultFieldsPerRec  = 0 // 0 means variable number of fields
)

var (
	// ErrQuote indicates a parsing error due to an unbalanced quote
	ErrQuote = errors.New("extraneous or missing " + string('"') + " in quoted field")
	// ErrBareQuote indicates a bare quote in a non-quoted field
	ErrBareQuote = errors.New(`bare '"' in non-quoted-field`)
	// ErrFieldCount indicates a wrong number of fields in a record
	ErrFieldCount = errors.New("wrong number of fields")
	// ErrTrailingComma indicates an extra delimiter at the end of a line
	ErrTrailingComma = errors.New("extra delimiter at end of line")
)

// Reader reads records from a CSV-encoded file.
// It provides better performance than the standard library's csv.Reader.
type Reader struct {
	// Comma is the field delimiter (default is ',').
	Comma rune

	// Comment, if not 0, is the comment character. Lines beginning with the
	// Comment character without preceding whitespace are ignored.
	Comment rune

	// FieldsPerRecord is the number of expected fields per record.
	// If FieldsPerRecord is positive, Read requires each record to
	// have the given number of fields. If FieldsPerRecord is 0, Read sets it to
	// the number of fields in the first record. If FieldsPerRecord is negative,
	// no check is made and records may have a variable number of fields.
	FieldsPerRecord int

	// If LazyQuotes is true, a quote may appear in an unquoted field and a
	// non-doubled quote may appear in a quoted field.
	LazyQuotes bool

	r *bufio.Reader

	// internal state
	numLine      int          // current line number
	reuseRecord  bool         // whether to reuse the record buffer
	fieldIndexes []int        // slice to store field boundaries for reuse
	lineBuffer   bytes.Buffer // reusable buffer for line processing
}

// NewReader returns a new Reader that reads from r.
func NewReader(r io.Reader) *Reader {
	return &Reader{
		Comma:        defaultDelimiter,
		Comment:      defaultComment,
		LazyQuotes:   defaultLazyQuotes,
		r:            bufio.NewReader(r),
		fieldIndexes: make([]int, 0, 10), // Pre-allocate some space for field indexes
	}
}

// Read reads one record from r. The record is a slice of strings with each
// string representing one field.
func (r *Reader) Read() (record []string, err error) {
	line, err := r.readLine()
	if err != nil {
		return nil, err
	}

	// Skip empty lines or comment lines
	line = strings.TrimSpace(line)
	if line == "" || (r.Comment != 0 && line[0] == byte(r.Comment)) {
		return r.Read()
	}

	// Reset field indexes
	r.fieldIndexes = r.fieldIndexes[:0]
	inQuotes := false
	fieldStart := 0

	// Parse the line to find field boundaries
	for i := 0; i < len(line); i++ {
		c := line[i]
		switch c {
		case '"':
			if inQuotes {
				// Inside quotes: handle escaped quotes "" by skipping the next quote
				if i+1 < len(line) && line[i+1] == '"' {
					i++ // consume the escaped quote
				} else {
					inQuotes = false
				}
			} else {
				// entering quoted field
				inQuotes = true
			}
		case byte(r.Comma):
			if !inQuotes {
				r.fieldIndexes = append(r.fieldIndexes, fieldStart, i)
				fieldStart = i + 1
			}
		}
	}

	// Add the last field
	if fieldStart <= len(line) {
		r.fieldIndexes = append(r.fieldIndexes, fieldStart, len(line))
	}

	// Handle the case where the line ends with a comma
	if len(line) > 0 && line[len(line)-1] == byte(r.Comma) {
		r.fieldIndexes = append(r.fieldIndexes, len(line), len(line))
	}

	// Allocate record with exact capacity needed
	record = make([]string, 0, len(r.fieldIndexes)/2)

	// Extract fields
	for i := 0; i < len(r.fieldIndexes); i += 2 {
		start, end := r.fieldIndexes[i], r.fieldIndexes[i+1]
		field := line[start:end]
		
		// Remove surrounding quotes if present
		if len(field) >= 2 && field[0] == '"' && field[len(field)-1] == '"' {
			field = field[1 : len(field)-1]
			// Replace double quotes with single quotes
			field = strings.ReplaceAll(field, "\"\"", "\"")
		}
		record = append(record, field)
	}

	// Validate field count if needed
	if r.FieldsPerRecord > 0 {
		if len(record) != r.FieldsPerRecord {
			return nil, ErrFieldCount
		}
	} else if r.FieldsPerRecord == 0 {
		r.FieldsPerRecord = len(record)
	}

	return record, nil
}

// readLine reads the next line from the input.
func (r *Reader) readLine() (string, error) {
	var line []byte
	for {
		l, isPrefix, err := r.r.ReadLine()
		if err != nil {
			return "", err
		}
		line = append(line, l...)
		if !isPrefix {
			break
		}
	}
	r.numLine++
	return string(line), nil
}

// ReadAll reads all the remaining records from r.
// Each record is a slice of fields.
// A successful call returns err == nil, not err == io.EOF.
func (r *Reader) ReadAll() ([][]string, error) {
	var records [][]string

	for {
		record, err := r.Read()
		if err == io.EOF {
			return records, nil
		}
		if err != nil {
			return nil, err
		}
		records = append(records, record)
	}
}
