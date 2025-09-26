package csvc

import (
	"bytes"
	stdcsv "encoding/csv"
	"io"
	"testing"
)

func benchmarkData() []byte {
	const row = "\"alpha\",beta,\"gamma\ninside\",delta,epsilon\r\n"
	buf := bytes.Repeat([]byte(row), 1024)
	return buf
}

func BenchmarkReader(b *testing.B) {
	data := benchmarkData()
	b.ReportAllocs()
	b.SetBytes(int64(len(data)))

	for i := 0; i < b.N; i++ {
		rdr := bytes.NewReader(data)
		var src io.Reader = rdr
		cr := NewReader(&src)
		cr.ReuseRecord = true

		for {
			if _, err := cr.Read(); err != nil {
				if err == io.EOF {
					break
				}
				b.Fatal(err)
			}
		}
	}
}

func BenchmarkEncodingCSV(b *testing.B) {
	data := benchmarkData()
	b.ReportAllocs()
	b.SetBytes(int64(len(data)))

	for i := 0; i < b.N; i++ {
		rdr := bytes.NewReader(data)
		cr := stdcsv.NewReader(rdr)

		for {
			if _, err := cr.Read(); err != nil {
				if err == io.EOF {
					break
				}
				b.Fatal(err)
			}
		}
	}
}
