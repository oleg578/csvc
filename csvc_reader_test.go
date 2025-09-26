package csvc

import (
	"bytes"
	"encoding/csv"
	"io"
	"math/rand"
	"strings"
	"testing"
	"time"
)

func TestReader_Basic(t *testing.T) {
	data := "a,b,c\n1,2,3\n"
	r := NewReader(strings.NewReader(data))

	rec, err := r.Read()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got, want := len(rec), 3; got != want {
		t.Fatalf("fields len=%d want %d", got, want)
	}
	if rec[0] != "a" || rec[1] != "b" || rec[2] != "c" {
		t.Fatalf("unexpected first record: %#v", rec)
	}

	rec, err = r.Read()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rec[0] != "1" || rec[1] != "2" || rec[2] != "3" {
		t.Fatalf("unexpected second record: %#v", rec)
	}

	_, err = r.Read()
	if err != io.EOF {
		t.Fatalf("expected EOF, got %v", err)
	}
}

func TestReader_Quotes(t *testing.T) {
	data := "name,comment\n\"Doe, John\",\"He said \"\"Hello\"\"\"\n"
	r := NewReader(strings.NewReader(data))

	rec, err := r.Read()
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if got, want := rec[0], "name"; got != want {
		t.Fatalf("got %q want %q", got, want)
	}

	rec, err = r.Read()
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if got, want := rec[0], "Doe, John"; got != want {
		t.Fatalf("got %q want %q", got, want)
	}
	if got, want := rec[1], "He said \"Hello\""; got != want {
		t.Fatalf("got %q want %q", got, want)
	}
}

func TestReader_CustomDelimiterAndComments(t *testing.T) {
	data := "# header to ignore\na;b;c\n1;2;3\n"
	r := NewReader(strings.NewReader(data))
	r.Comma = ';'
	r.Comment = '#'

	rec, err := r.Read()
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if rec[0] != "a" || rec[1] != "b" || rec[2] != "c" {
		t.Fatalf("unexpected first rec: %#v", rec)
	}

	rec, err = r.Read()
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if rec[0] != "1" || rec[1] != "2" || rec[2] != "3" {
		t.Fatalf("unexpected second rec: %#v", rec)
	}
}

func TestReader_FieldsPerRecord(t *testing.T) {
	data := "a,b\n1,2,3\n"
	r := NewReader(strings.NewReader(data))
	r.FieldsPerRecord = 2

	_, err := r.Read() // sets baseline ok
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	_, err = r.Read() // wrong number of fields
	if err == nil {
		t.Fatal("expected error for wrong field count")
	}
	if err != ErrFieldCount {
		t.Fatalf("expected ErrFieldCount, got %v", err)
	}
}

func TestReader_ReadAll(t *testing.T) {
	data := "a,b,c\n1,2,3\n4,5,6\n"
	r := NewReader(strings.NewReader(data))

	recs, err := r.ReadAll()
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if got, want := len(recs), 3; got != want {
		t.Fatalf("got %d recs want %d", got, want)
	}
	if recs[2][2] != "6" {
		t.Fatalf("unexpected last field: %q", recs[2][2])
	}
}

// generateCSV generates n records with m fields each, with occasional quotes and commas
func generateCSV(n, m int) []byte {
	var b bytes.Buffer
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			val := randomField(j)
			if j > 0 {
				b.WriteByte(',')
			}
			if strings.ContainsAny(val, ",\" ") {
				b.WriteByte('"')
				b.WriteString(strings.ReplaceAll(val, "\"", "\"\""))
				b.WriteByte('"')
			} else {
				b.WriteString(val)
			}
		}
		b.WriteByte('\n')
	}
	return b.Bytes()
}

func randomField(seed int) string {
	rnd := rand.Intn(5)
	switch rnd {
	case 0:
		return "foo"
	case 1:
		return "bar"
	case 2:
		return "baz"
	case 3:
		return "hello,world"
	default:
		return "say \"hi\""
	}
}

func BenchmarkReader_Read(b *testing.B) {
	rand.Seed(time.Now().UnixNano())
	data := generateCSV(2000, 10)
	b.SetBytes(int64(len(data)))
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		r := NewReader(bytes.NewReader(data))
		for {
			_, err := r.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				b.Fatalf("unexpected err: %v", err)
			}
		}
	}
}

func BenchmarkStdCSV_Read(b *testing.B) {
	rand.Seed(time.Now().UnixNano())
	data := generateCSV(2000, 10)
	b.SetBytes(int64(len(data)))
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		r := csv.NewReader(bytes.NewReader(data))
		for {
			_, err := r.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				b.Fatalf("unexpected err: %v", err)
			}
		}
	}
}
