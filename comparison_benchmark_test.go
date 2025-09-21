package csvc

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"strings"
	"testing"
)

// generateCSVDataForComparison creates test CSV data with specified rows and columns
func generateCSVDataForComparison(rows, cols int, quoted bool) string {
	var builder strings.Builder

	// Write header
	for i := range cols {
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
	for row := range rows {
		for col := range cols {
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

// generateComplexCSVDataForComparison creates CSV with mixed quoted/unquoted fields
func generateComplexCSVDataForComparison(rows int) string {
	var builder strings.Builder

	// Header
	builder.WriteString("id,name,description,price,tags\n")

	// Data rows with various complexities
	for i := range rows {
		builder.WriteString(fmt.Sprintf("%d,", i+1))
		builder.WriteString(fmt.Sprintf("\"Product %d\",", i+1))
		builder.WriteString(fmt.Sprintf("\"Description with, commas and \"\"quotes\"\" for product %d\",", i+1))
		builder.WriteString(fmt.Sprintf("%.2f,", float64(i+1)*1.23))
		builder.WriteString("\"tag1,tag2,tag3\"")
		builder.WriteString("\n")
	}

	return builder.String()
}

// Benchmark: Small Simple CSV - CSVC vs Go Built-in
func BenchmarkComparison_SmallSimple_CSVC(b *testing.B) {
	data := generateCSVDataForComparison(100, 5, false)


	for b.Loop() {
		reader := NewReader(bufio.NewReader(strings.NewReader(data)))

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
	data := generateCSVDataForComparison(100, 5, false)


	for b.Loop() {
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

// Benchmark: Small Quoted CSV - CSVC vs Go Built-in
func BenchmarkComparison_SmallQuoted_CSVC(b *testing.B) {
	data := generateCSVDataForComparison(100, 5, true)


	for b.Loop() {
		reader := NewReader(bufio.NewReader(strings.NewReader(data)))

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
	data := generateCSVDataForComparison(100, 5, true)


	for b.Loop() {
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

// Benchmark: Medium Simple CSV - CSVC vs Go Built-in
func BenchmarkComparison_MediumSimple_CSVC(b *testing.B) {
	data := generateCSVDataForComparison(1000, 10, false)


	for b.Loop() {
		reader := NewReader(bufio.NewReader(strings.NewReader(data)))

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
	data := generateCSVDataForComparison(1000, 10, false)


	for b.Loop() {
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

// Benchmark: Medium Quoted CSV - CSVC vs Go Built-in
func BenchmarkComparison_MediumQuoted_CSVC(b *testing.B) {
	data := generateCSVDataForComparison(1000, 10, true)


	for b.Loop() {
		reader := NewReader(bufio.NewReader(strings.NewReader(data)))

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
	data := generateCSVDataForComparison(1000, 10, true)


	for b.Loop() {
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

// Benchmark: Large Simple CSV - CSVC vs Go Built-in
func BenchmarkComparison_LargeSimple_CSVC(b *testing.B) {
	data := generateCSVDataForComparison(5000, 20, false)


	for b.Loop() {
		reader := NewReader(bufio.NewReader(strings.NewReader(data)))

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
	data := generateCSVDataForComparison(5000, 20, false)

	
	for b.Loop() {
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

// Benchmark: Complex Fields - CSVC vs Go Built-in
func BenchmarkComparison_ComplexFields_CSVC(b *testing.B) {
	data := generateComplexCSVDataForComparison(1000)


	for b.Loop() {
		reader := NewReader(bufio.NewReader(strings.NewReader(data)))

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
	data := generateComplexCSVDataForComparison(1000)


	for b.Loop() {
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

// Benchmark: Single Record - CSVC vs Go Built-in
func BenchmarkComparison_SingleRecord_CSVC(b *testing.B) {
	data := "field1,field2,field3,field4,field5\n"


	for b.Loop() {
		reader := NewReader(bufio.NewReader(strings.NewReader(data)))
		_, err := reader.Read()
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkComparison_SingleRecord_GoBuiltin(b *testing.B) {
	data := "field1,field2,field3,field4,field5\n"


	for b.Loop() {
		reader := csv.NewReader(strings.NewReader(data))
		_, err := reader.Read()
		if err != nil {
			b.Fatal(err)
		}
	}
}

// Benchmark: Single Quoted Record - CSVC vs Go Built-in
func BenchmarkComparison_SingleQuotedRecord_CSVC(b *testing.B) {
	data := "\"field1\",\"field2\",\"field3\",\"field4\",\"field5\"\n"


	for b.Loop() {
		reader := NewReader(bufio.NewReader(strings.NewReader(data)))
		_, err := reader.Read()
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkComparison_SingleQuotedRecord_GoBuiltin(b *testing.B) {
	data := "\"field1\",\"field2\",\"field3\",\"field4\",\"field5\"\n"


	for b.Loop() {
		reader := csv.NewReader(strings.NewReader(data))
		_, err := reader.Read()
		if err != nil {
			b.Fatal(err)
		}
	}
}

// Benchmark: Escaped Quotes - CSVC vs Go Built-in
func BenchmarkComparison_EscapedQuotes_CSVC(b *testing.B) {
	data := "\"field with \"\"quotes\"\"\",\"another \"\"quoted\"\" field\",normal\n"


	for b.Loop() {
		reader := NewReader(bufio.NewReader(strings.NewReader(data)))
		_, err := reader.Read()
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkComparison_EscapedQuotes_GoBuiltin(b *testing.B) {
	data := "\"field with \"\"quotes\"\"\",\"another \"\"quoted\"\" field\",normal\n"


	for b.Loop() {
		reader := csv.NewReader(strings.NewReader(data))
		_, err := reader.Read()
		if err != nil {
			b.Fatal(err)
		}
	}
}

// Benchmark: Empty Fields - CSVC vs Go Built-in
func BenchmarkComparison_EmptyFields_CSVC(b *testing.B) {
	data := strings.Repeat(",,,,\n", 1000)


	for b.Loop() {
		reader := NewReader(bufio.NewReader(strings.NewReader(data)))

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


	for b.Loop() {
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

// Benchmark: Custom Delimiter - CSVC vs Go Built-in
func BenchmarkComparison_CustomDelimiter_CSVC(b *testing.B) {
	data := generateCSVDataForComparison(1000, 10, false)
	data = strings.ReplaceAll(data, ",", ";") // Replace commas with semicolons


	for b.Loop() {
		reader := NewReader(bufio.NewReader(strings.NewReader(data)))
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
	data := generateCSVDataForComparison(1000, 10, false)
	data = strings.ReplaceAll(data, ",", ";") // Replace commas with semicolons


	for b.Loop() {
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

// Benchmark: Many Columns - CSVC vs Go Built-in
func BenchmarkComparison_ManyColumns_CSVC(b *testing.B) {
	data := generateCSVDataForComparison(100, 50, false)


	for b.Loop() {
		reader := NewReader(bufio.NewReader(strings.NewReader(data)))

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
	data := generateCSVDataForComparison(100, 50, false)


	for b.Loop() {
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

// Benchmark: Multiline Fields - CSVC vs Go Built-in
func BenchmarkComparison_MultilineFields_CSVC(b *testing.B) {
	multilineField := "line1\nline2\nline3"
	data := fmt.Sprintf("field1,\"%s\",field3\n", multilineField)
	data = strings.Repeat(data, 100)


	for b.Loop() {
		reader := NewReader(bufio.NewReader(strings.NewReader(data)))

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

func BenchmarkComparison_MultilineFields_GoBuiltin(b *testing.B) {
	multilineField := "line1\nline2\nline3"
	data := fmt.Sprintf("field1,\"%s\",field3\n", multilineField)
	data = strings.Repeat(data, 100)


	for b.Loop() {
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
