package csvc

import (
	"encoding/csv"
	"errors"
	"io"
	"reflect"
	"strings"
	"testing"
)

func newTestReader(s string) *Reader {
	var src io.Reader = strings.NewReader(s)
	return NewReader(&src)
}

func TestReaderBasic(t *testing.T) {
	r := newTestReader("a,b,c\n1,2,3\n")

	rec, err := r.Read()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !reflect.DeepEqual(rec, []string{"a", "b", "c"}) {
		t.Fatalf("unexpected record: %#v", rec)
	}

	rec, err = r.Read()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !reflect.DeepEqual(rec, []string{"1", "2", "3"}) {
		t.Fatalf("unexpected record: %#v", rec)
	}

	rec, err = r.Read()
	if err != io.EOF {
		t.Fatalf("expected EOF, got %v with record %#v", err, rec)
	}
}

func TestReaderQuotedFields(t *testing.T) {
	r := newTestReader("\"a\",\"b\",\"c\"\n\"x,y\",z,\"1\"\"2\"\n")

	rec, err := r.Read()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expect := []string{"a", "b", "c"}
	if !reflect.DeepEqual(rec, expect) {
		t.Fatalf("unexpected record: %#v", rec)
	}

	rec, err = r.Read()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expect = []string{"x,y", "z", "1\"2"}
	if !reflect.DeepEqual(rec, expect) {
		t.Fatalf("unexpected record: %#v", rec)
	}
}

func TestReaderHandlesCRLF(t *testing.T) {
	r := newTestReader("a,b,c\r\n1,2,3\r\n")

	rec, err := r.Read()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !reflect.DeepEqual(rec, []string{"a", "b", "c"}) {
		t.Fatalf("unexpected record: %#v", rec)
	}

	rec, err = r.Read()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !reflect.DeepEqual(rec, []string{"1", "2", "3"}) {
		t.Fatalf("unexpected record: %#v", rec)
	}
}

func TestReaderBareQuoteError(t *testing.T) {
	r := newTestReader("a\"b,c\n")
	_, err := r.Read()
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if !errors.Is(err, ErrBareQuote) {
		t.Fatalf("expected ErrBareQuote, got %v", err)
	}
	var perr *ParseError
	if !errors.As(err, &perr) {
		t.Fatalf("expected ParseError, got %T", err)
	}
	if perr.Line != 1 {
		t.Fatalf("expected line 1, got %d", perr.Line)
	}
}

func TestReaderUnterminatedQuote(t *testing.T) {
	r := newTestReader("\"unterminated\n")
	_, err := r.Read()
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if !errors.Is(err, ErrUnterminatedQuote) {
		t.Fatalf("expected ErrUnterminatedQuote, got %v", err)
	}
}

func TestReaderReuseRecord(t *testing.T) {
	r := newTestReader("a,b\n1,2\n")
	r.ReuseRecord = true

	first, err := r.Read()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(first) == 0 {
		t.Fatalf("expected data in first record")
	}
	firstPtr := &first[0]

	second, err := r.Read()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(second) == 0 {
		t.Fatalf("expected data in second record")
	}
	secondPtr := &second[0]

	if firstPtr != secondPtr {
		t.Fatalf("expected slice reuse, but backing arrays differ")
	}
}

func TestReaderCustomComma(t *testing.T) {
	r := newTestReader("a;b;c\n1;2;3\n")
	r.Comma = ';'

	rec, err := r.Read()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !reflect.DeepEqual(rec, []string{"a", "b", "c"}) {
		t.Fatalf("unexpected record: %#v", rec)
	}
}

func TestReaderParityWithEncoding(t *testing.T) {
	data := strings.Repeat("\"field,with,comma\",plain,\"multi\nline\"\r\n", 8)

	var src1 io.Reader = strings.NewReader(data)
	var src2 io.Reader = strings.NewReader(data)

	fast := NewReader(&src1)
	std := csv.NewReader(src2)

	for {
		fastRec, fastErr := fast.Read()
		stdRec, stdErr := std.Read()

		if fastErr == io.EOF && stdErr == io.EOF {
			return
		}
		if fastErr == io.EOF || stdErr == io.EOF {
			t.Fatalf("mismatched EOF: fast=%v std=%v", fastErr, stdErr)
		}
		if fastErr != nil || stdErr != nil {
			t.Fatalf("unexpected errors: fast=%v std=%v", fastErr, stdErr)
		}
		if !reflect.DeepEqual(fastRec, stdRec) {
			t.Fatalf("records differ: fast=%#v std=%#v", fastRec, stdRec)
		}
	}
}
