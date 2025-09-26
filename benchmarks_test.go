package csvc

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"strings"
	"testing"
)

// Note: We reuse generateCSV from csvc_reader_test.go for general cases.
// This file adds additional data generators and table-driven benchmarks comparing
// csvc vs the standard library encoding/csv.
// generateQuotedHeavyCSV creates n records with m fields emphasizing quotes and commas.
func generateQuotedHeavyCSV(n, m int) []byte {
	var b bytes.Buffer
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			val := fmt.Sprintf("value %d, %d says \"hello, \"\"world\"\"\"", i, j)
			if j > 0 {
				b.WriteByte(',')
			}
			// Always quote and double quotes inside
			b.WriteByte('"')
			b.WriteString(strings.ReplaceAll(val, "\"", "\"\""))
			b.WriteByte('"')
		}
		b.WriteByte('\n')
	}
	return b.Bytes()
}

// generateSemicolonCSV creates n records with m fields using ';' as delimiter.
func generateSemicolonCSV(n, m int) []byte {
	var b bytes.Buffer
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			val := fmt.Sprintf("field_%d_%d", i, j)
			if j > 0 {
				b.WriteByte(';')
			}
			b.WriteString(val)
		}
		b.WriteByte('\n')
	}
	return b.Bytes()
}

func BenchmarkCompare_Read_Sizes(b *testing.B) {
	sizes := []struct{
		name string
		n, m int
	}{
		{"small", 200, 10},
		{"medium", 2_000, 10},
		{"large", 20_000, 10},
	}

	for _, sz := range sizes {
		b.Run(sz.name, func(b *testing.B) {
			data := generateCSV(sz.n, sz.m)
			b.SetBytes(int64(len(data)))

			b.Run("csvc", func(b *testing.B) {
				b.ReportAllocs()
				for i := 0; i < b.N; i++ {
					r := NewReader(bytes.NewReader(data))
					for {
						_, err := r.Read()
						if err != nil {
							if err == io.EOF {
								break
							}
							b.Fatalf("csvc read error: %v", err)
						}
					}
				}
			})

			b.Run("stdcsv", func(b *testing.B) {
				b.ReportAllocs()
				for i := 0; i < b.N; i++ {
					r := csv.NewReader(bytes.NewReader(data))
					for {
						_, err := r.Read()
						if err != nil {
							if err == io.EOF {
								break
							}
							b.Fatalf("std csv read error: %v", err)
						}
					}
				}
			})
		})
	}
}

func BenchmarkCompare_Read_QuotedHeavy(b *testing.B) {
	data := generateQuotedHeavyCSV(5_000, 8)
	b.SetBytes(int64(len(data)))

	b.Run("csvc", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			r := NewReader(bytes.NewReader(data))
			for {
				_, err := r.Read()
				if err != nil {
					if err == io.EOF {
						break
					}
					b.Fatalf("csvc read error: %v", err)
				}
			}
		}
	})

	b.Run("stdcsv", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			r := csv.NewReader(bytes.NewReader(data))
			for {
				_, err := r.Read()
				if err != nil {
					if err == io.EOF {
						break
					}
					b.Fatalf("std csv read error: %v", err)
				}
			}
		}
	})
}

func BenchmarkCompare_Read_Semicolon(b *testing.B) {
	data := generateSemicolonCSV(10_000, 8)
	b.SetBytes(int64(len(data)))

	b.Run("csvc", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			r := NewReader(bytes.NewReader(data))
			r.Comma = ';'
			for {
				_, err := r.Read()
				if err != nil {
					if err == io.EOF {
						break
					}
					b.Fatalf("csvc read error: %v", err)
				}
			}
		}
	})

	b.Run("stdcsv", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			r := csv.NewReader(bytes.NewReader(data))
			r.Comma = ';'
			for {
				_, err := r.Read()
				if err != nil {
					if err == io.EOF {
						break
					}
					b.Fatalf("std csv read error: %v", err)
				}
			}
		}
	})
}
