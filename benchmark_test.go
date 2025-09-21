package csvc

import (
	"bufio"
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

// BenchmarkReader_Read_SmallSimple benchmarks reading small simple CSV
func BenchmarkReader_Read_SmallSimple(b *testing.B) {
	data := generateCSVData(100, 5, false)


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

// BenchmarkReader_Read_SmallQuoted benchmarks reading small quoted CSV
func BenchmarkReader_Read_SmallQuoted(b *testing.B) {
	data := generateCSVData(100, 5, true)


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

// BenchmarkReader_Read_MediumSimple benchmarks reading medium simple CSV
func BenchmarkReader_Read_MediumSimple(b *testing.B) {
	data := generateCSVData(1000, 10, false)


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

// BenchmarkReader_Read_MediumQuoted benchmarks reading medium quoted CSV
func BenchmarkReader_Read_MediumQuoted(b *testing.B) {
	data := generateCSVData(1000, 10, true)


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

// BenchmarkReader_Read_LargeSimple benchmarks reading large simple CSV
func BenchmarkReader_Read_LargeSimple(b *testing.B) {
	data := generateCSVData(10000, 20, false)


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

// BenchmarkReader_Read_LargeQuoted benchmarks reading large quoted CSV
func BenchmarkReader_Read_LargeQuoted(b *testing.B) {
	data := generateCSVData(10000, 20, true)


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

// BenchmarkReader_Read_ComplexFields benchmarks reading CSV with complex quoted fields
func BenchmarkReader_Read_ComplexFields(b *testing.B) {
	data := generateComplexCSVData(1000)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
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

// BenchmarkReader_Read_SingleRecord benchmarks reading a single CSV record
func BenchmarkReader_Read_SingleRecord(b *testing.B) {
	data := "field1,field2,field3,field4,field5\n"


	for b.Loop() {
		reader := NewReader(bufio.NewReader(strings.NewReader(data)))
		_, err := reader.Read()
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkReader_Read_SingleQuotedRecord benchmarks reading a single quoted CSV record
func BenchmarkReader_Read_SingleQuotedRecord(b *testing.B) {
	data := "\"field1\",\"field2\",\"field3\",\"field4\",\"field5\"\n"


	for b.Loop() {
		reader := NewReader(bufio.NewReader(strings.NewReader(data)))
		_, err := reader.Read()
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkReader_Read_EscapedQuotes benchmarks reading CSV with escaped quotes
func BenchmarkReader_Read_EscapedQuotes(b *testing.B) {
	data := "\"field with \"\"quotes\"\"\",\"another \"\"quoted\"\" field\",normal\n"


	for b.Loop() {
		reader := NewReader(bufio.NewReader(strings.NewReader(data)))
		_, err := reader.Read()
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkReader_Read_ManyColumns benchmarks reading CSV with many columns
func BenchmarkReader_Read_ManyColumns(b *testing.B) {
	data := generateCSVData(100, 50, false)


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

// BenchmarkReader_Read_EmptyFields benchmarks reading CSV with many empty fields
func BenchmarkReader_Read_EmptyFields(b *testing.B) {
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

// BenchmarkReader_Read_CustomDelimiter benchmarks reading CSV with custom delimiter
func BenchmarkReader_Read_CustomDelimiter(b *testing.B) {
	data := generateCSVData(1000, 10, false)
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

// BenchmarkReader_Read_LongFields benchmarks reading CSV with very long fields
func BenchmarkReader_Read_LongFields(b *testing.B) {
	longField := strings.Repeat("a", 1000)
	data := fmt.Sprintf("\"%s\",\"%s\",\"%s\"\n", longField, longField, longField)
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

// BenchmarkReader_Read_MultilineFields benchmarks reading CSV with multiline fields
func BenchmarkReader_Read_MultilineFields(b *testing.B) {
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
