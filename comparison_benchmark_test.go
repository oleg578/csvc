package csvc

import (
	"encoding/csv"
	"fmt"
	"io"
	"strings"
	"testing"
)

// generateCSVData creates test CSV data with specified rows and columns
func generateCSVData(rows, cols int, quoted bool) string {
	var builder strings.Builder

	// Write header
	for i := 0; i < cols; i++ {
		if i > 0 {
			builder.WriteString(",")
		}
		if quoted {
			builder.WriteString(fmt.Sprintf("\"col%d\"", i+1))
		} else {
			builder.WriteString(fmt.Sprintf("col%d", i+1))
		}
	}
	builder.WriteString("\n")

	// Write data rows
	for row := 0; row < rows; row++ {
		for col := 0; col < cols; col++ {
			if col > 0 {
				builder.WriteString(",")
			}
			if quoted {
				builder.WriteString(fmt.Sprintf("\"data%d_%d\"", row+1, col+1))
			} else {
				builder.WriteString(fmt.Sprintf("data%d_%d", row+1, col+1))
			}
		}
		builder.WriteString("\n")
	}

	return builder.String()
}

// generateComplexCSVData creates CSV with mixed quoted/unquoted fields and special characters
func generateComplexCSVData(rows int) string {
	var builder strings.Builder

	// Header
	builder.WriteString("id,name,description,price,tags\n")

	// Data rows with various complexities
	for i := 0; i < rows; i++ {
		builder.WriteString(fmt.Sprintf("%d,", i+1))
		builder.WriteString(fmt.Sprintf("\"Product %d\",", i+1))
		builder.WriteString(fmt.Sprintf("\"Description with, commas and \"\"quotes\"\" for product %d\",", i+1))
		builder.WriteString(fmt.Sprintf("%.2f,", float64(i+1)*1.23))
		builder.WriteString("\"tag1,tag2,tag3\"")
		builder.WriteString("\n")
	}

	return builder.String()
}

// Benchmark Comparison: Small Simple Dataset (100 rows, 5 columns)
func BenchmarkComparison_SmallSimple_CSVC(b *testing.B) {
	data := generateCSVData(100, 5, false)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		reader := NewReader(strings.NewReader(data))

		for {
			_, err := reader.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				b.Fatal(err)
			}
		}
	}
}

func BenchmarkComparison_SmallSimple_GoBuiltin(b *testing.B) {
	data := generateCSVData(100, 5, false)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		reader := csv.NewReader(strings.NewReader(data))

		for {
			_, err := reader.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				b.Fatal(err)
			}
		}
	}
}

// Benchmark Comparison: Small Quoted Dataset (100 rows, 5 columns)
func BenchmarkComparison_SmallQuoted_CSVC(b *testing.B) {
	data := generateCSVData(100, 5, true)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		reader := NewReader(strings.NewReader(data))

		for {
			_, err := reader.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				b.Fatal(err)
			}
		}
	}
}

func BenchmarkComparison_SmallQuoted_GoBuiltin(b *testing.B) {
	data := generateCSVData(100, 5, true)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		reader := csv.NewReader(strings.NewReader(data))

		for {
			_, err := reader.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				b.Fatal(err)
			}
		}
	}
}

// Benchmark Comparison: Medium Simple Dataset (1000 rows, 10 columns)
func BenchmarkComparison_MediumSimple_CSVC(b *testing.B) {
	data := generateCSVData(1000, 10, false)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		reader := NewReader(strings.NewReader(data))

		for {
			_, err := reader.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				b.Fatal(err)
			}
		}
	}
}

func BenchmarkComparison_MediumSimple_GoBuiltin(b *testing.B) {
	data := generateCSVData(1000, 10, false)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		reader := csv.NewReader(strings.NewReader(data))

		for {
			_, err := reader.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				b.Fatal(err)
			}
		}
	}
}

// Benchmark Comparison: Medium Quoted Dataset (1000 rows, 10 columns)
func BenchmarkComparison_MediumQuoted_CSVC(b *testing.B) {
	data := generateCSVData(1000, 10, true)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		reader := NewReader(strings.NewReader(data))

		for {
			_, err := reader.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				b.Fatal(err)
			}
		}
	}
}

func BenchmarkComparison_MediumQuoted_GoBuiltin(b *testing.B) {
	data := generateCSVData(1000, 10, true)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		reader := csv.NewReader(strings.NewReader(data))

		for {
			_, err := reader.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				b.Fatal(err)
			}
		}
	}
}

// Benchmark Comparison: Large Simple Dataset (10000 rows, 20 columns)
func BenchmarkComparison_LargeSimple_CSVC(b *testing.B) {
	data := generateCSVData(10000, 20, false)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		reader := NewReader(strings.NewReader(data))

		for {
			_, err := reader.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				b.Fatal(err)
			}
		}
	}
}

func BenchmarkComparison_LargeSimple_GoBuiltin(b *testing.B) {
	data := generateCSVData(10000, 20, false)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		reader := csv.NewReader(strings.NewReader(data))

		for {
			_, err := reader.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				b.Fatal(err)
			}
		}
	}
}

// Benchmark Comparison: Complex Fields with Escaped Quotes and Commas
func BenchmarkComparison_ComplexFields_CSVC(b *testing.B) {
	data := generateComplexCSVData(1000)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		reader := NewReader(strings.NewReader(data))

		for {
			_, err := reader.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				b.Fatal(err)
			}
		}
	}
}

func BenchmarkComparison_ComplexFields_GoBuiltin(b *testing.B) {
	data := generateComplexCSVData(1000)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		reader := csv.NewReader(strings.NewReader(data))

		for {
			_, err := reader.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				b.Fatal(err)
			}
		}
	}
}

// Benchmark Comparison: Single Record Performance
func BenchmarkComparison_SingleRecord_CSVC(b *testing.B) {
	data := "field1,field2,field3,field4,field5\n"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		reader := NewReader(strings.NewReader(data))
		_, err := reader.Read()
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkComparison_SingleRecord_GoBuiltin(b *testing.B) {
	data := "field1,field2,field3,field4,field5\n"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		reader := csv.NewReader(strings.NewReader(data))
		_, err := reader.Read()
		if err != nil {
			b.Fatal(err)
		}
	}
}

// Benchmark Comparison: Single Quoted Record Performance
func BenchmarkComparison_SingleQuotedRecord_CSVC(b *testing.B) {
	data := "\"field1\",\"field2\",\"field3\",\"field4\",\"field5\"\n"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		reader := NewReader(strings.NewReader(data))
		_, err := reader.Read()
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkComparison_SingleQuotedRecord_GoBuiltin(b *testing.B) {
	data := "\"field1\",\"field2\",\"field3\",\"field4\",\"field5\"\n"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		reader := csv.NewReader(strings.NewReader(data))
		_, err := reader.Read()
		if err != nil {
			b.Fatal(err)
		}
	}
}

// Benchmark Comparison: Escaped Quotes Performance
func BenchmarkComparison_EscapedQuotes_CSVC(b *testing.B) {
	data := "\"field with \"\"quotes\"\"\",\"another \"\"quoted\"\" field\",normal\n"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		reader := NewReader(strings.NewReader(data))
		_, err := reader.Read()
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkComparison_EscapedQuotes_GoBuiltin(b *testing.B) {
	data := "\"field with \"\"quotes\"\"\",\"another \"\"quoted\"\" field\",normal\n"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		reader := csv.NewReader(strings.NewReader(data))
		_, err := reader.Read()
		if err != nil {
			b.Fatal(err)
		}
	}
}

// Benchmark Comparison: Empty Fields Performance
func BenchmarkComparison_EmptyFields_CSVC(b *testing.B) {
	data := strings.Repeat(",,,,\n", 1000)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		reader := NewReader(strings.NewReader(data))

		for {
			_, err := reader.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				b.Fatal(err)
			}
		}
	}
}

func BenchmarkComparison_EmptyFields_GoBuiltin(b *testing.B) {
	data := strings.Repeat(",,,,\n", 1000)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		reader := csv.NewReader(strings.NewReader(data))

		for {
			_, err := reader.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				b.Fatal(err)
			}
		}
	}
}

// Benchmark Comparison: Custom Delimiter Performance
func BenchmarkComparison_CustomDelimiter_CSVC(b *testing.B) {
	data := generateCSVData(1000, 10, false)
	data = strings.ReplaceAll(data, ",", ";") // Replace commas with semicolons

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		reader := NewReader(strings.NewReader(data))
		reader.Comma = ';'

		for {
			_, err := reader.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				b.Fatal(err)
			}
		}
	}
}

func BenchmarkComparison_CustomDelimiter_GoBuiltin(b *testing.B) {
	data := generateCSVData(1000, 10, false)
	data = strings.ReplaceAll(data, ",", ";") // Replace commas with semicolons

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		reader := csv.NewReader(strings.NewReader(data))
		reader.Comma = ';'

		for {
			_, err := reader.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				b.Fatal(err)
			}
		}
	}
}

// Benchmark Comparison: Many Columns Performance
func BenchmarkComparison_ManyColumns_CSVC(b *testing.B) {
	data := generateCSVData(100, 50, false)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		reader := NewReader(strings.NewReader(data))

		for {
			_, err := reader.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				b.Fatal(err)
			}
		}
	}
}

func BenchmarkComparison_ManyColumns_GoBuiltin(b *testing.B) {
	data := generateCSVData(100, 50, false)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		reader := csv.NewReader(strings.NewReader(data))

		for {
			_, err := reader.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				b.Fatal(err)
			}
		}
	}
}

// Benchmark Comparison: Long Fields Performance
func BenchmarkComparison_LongFields_CSVC(b *testing.B) {
	longField := strings.Repeat("a", 1000)
	data := fmt.Sprintf("\"%s\",\"%s\",\"%s\"\n", longField, longField, longField)
	data = strings.Repeat(data, 100)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		reader := NewReader(strings.NewReader(data))

		for {
			_, err := reader.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				b.Fatal(err)
			}
		}
	}
}

func BenchmarkComparison_LongFields_GoBuiltin(b *testing.B) {
	longField := strings.Repeat("a", 1000)
	data := fmt.Sprintf("\"%s\",\"%s\",\"%s\"\n", longField, longField, longField)
	data = strings.Repeat(data, 100)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		reader := csv.NewReader(strings.NewReader(data))

		for {
			_, err := reader.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				b.Fatal(err)
			}
		}
	}
}
