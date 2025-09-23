package csvc

import (
	"bufio"
	"io"
	"reflect"
	"strings"
	"testing"
)

func TestNewReader(t *testing.T) {
	input := "test,data"
	reader := NewReader(bufio.NewReader(strings.NewReader(input)))

	if reader == nil {
		t.Fatal("NewReader() returned nil")
	}

	if reader.Comma != ',' {
		t.Errorf("Expected default comma to be ',', got %c", reader.Comma)
	}

	if reader.r == nil {
		t.Fatal("Reader.r is nil")
	}
}

func TestReader_Read_BasicFields(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "simple fields",
			input:    "field1,field2,field3\n",
			expected: []string{"field1", "field2", "field3"},
		},
		{
			name:     "single field",
			input:    "singlefield\n",
			expected: []string{"singlefield"},
		},
		{
			name:     "empty fields",
			input:    ",,\n",
			expected: []string{"", "", ""},
		},
		{
			name:     "mixed empty and filled",
			input:    "field1,,field3\n",
			expected: []string{"field1", "", "field3"},
		},
		{
			name:     "trailing empty field",
			input:    "field1,field2,\n",
			expected: []string{"field1", "field2", ""},
		},
		{
			name:     "leading empty field",
			input:    ",field2,field3\n",
			expected: []string{"", "field2", "field3"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := NewReader(bufio.NewReader(strings.NewReader(tt.input)))
			result, err := reader.Read()

			if err != nil {
				t.Fatalf("Read() error = %v", err)
			}

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Read() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestReader_Read_QuotedFields(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "quoted field with comma",
			input:    "\"field, with comma\",normal\n",
			expected: []string{"field, with comma", "normal"},
		},
		{
			name:     "simple quoted field",
			input:    "\"quoted field\",normal\n",
			expected: []string{"quoted field", "normal"},
		},
		{
			name:     "all quoted fields",
			input:    "\"field1\",\"field2\",\"field3\"\n",
			expected: []string{"field1", "field2", "field3"},
		},
		{
			name:     "quoted field with spaces",
			input:    "\"  spaced field  \",normal\n",
			expected: []string{"  spaced field  ", "normal"},
		},
		{
			name:     "empty quoted field",
			input:    "\"\",normal\n",
			expected: []string{"", "normal"},
		},
		{
			name:     "quoted field at end",
			input:    "normal,\"quoted field\"\n",
			expected: []string{"normal", "quoted field"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := NewReader(bufio.NewReader(strings.NewReader(tt.input)))
			result, err := reader.Read()

			if err != nil {
				t.Fatalf("Read() error = %v", err)
			}

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Read() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestReader_Read_EscapedQuotes(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "escaped quotes in middle",
			input:    "\"field with \"\"quotes\"\"\",normal\n",
			expected: []string{"field with \"quotes\"", "normal"},
		},
		{
			name:     "multiple escaped quotes",
			input:    "\"\"\"quoted\"\" and \"\"more\"\"\",normal\n",
			expected: []string{"\"quoted\" and \"more\"", "normal"},
		},
		{
			name:     "escaped quote at start",
			input:    "\"\"\"starts with quote\",normal\n",
			expected: []string{"\"starts with quote", "normal"},
		},
		{
			name:     "escaped quote at end",
			input:    "\"ends with quote\"\"\",normal\n",
			expected: []string{"ends with quote\"", "normal"},
		},
		{
			name:     "only escaped quotes",
			input:    "\"\"\"\"\n",
			expected: []string{"\""},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := NewReader(bufio.NewReader(strings.NewReader(tt.input)))
			result, err := reader.Read()

			if err != nil {
				t.Fatalf("Read() error = %v", err)
			}

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Read() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestReader_Read_LineEndings(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "LF ending",
			input:    "field1,field2\n",
			expected: []string{"field1", "field2"},
		},
		{
			name:     "CRLF ending",
			input:    "field1,field2\r\n",
			expected: []string{"field1", "field2"},
		},
		{
			name:     "CR ending",
			input:    "field1,field2\r",
			expected: []string{"field1", "field2"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := NewReader(bufio.NewReader(strings.NewReader(tt.input)))
			result, err := reader.Read()

			if err != nil {
				t.Fatalf("Read() error = %v", err)
			}

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Read() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestReader_Read_CustomDelimiter(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		delimiter byte
		expected  []string
	}{
		{
			name:      "semicolon delimiter",
			input:     "field1;field2;field3\n",
			delimiter: ';',
			expected:  []string{"field1", "field2", "field3"},
		},
		{
			name:      "tab delimiter",
			input:     "field1\tfield2\tfield3\n",
			delimiter: '\t',
			expected:  []string{"field1", "field2", "field3"},
		},
		{
			name:      "pipe delimiter",
			input:     "field1|field2|field3\n",
			delimiter: '|',
			expected:  []string{"field1", "field2", "field3"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := NewReader(bufio.NewReader(strings.NewReader(tt.input)))
			reader.Comma = tt.delimiter
			result, err := reader.Read()

			if err != nil {
				t.Fatalf("Read() error = %v", err)
			}

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Read() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestReader_Read_ComplexCases(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "quoted field with comma and quotes",
			input:    "normal,\"field, with \"\"comma\"\" and quotes\",last\n",
			expected: []string{"normal", "field, with \"comma\" and quotes", "last"},
		},
		{
			name:     "numbers and text mixed",
			input:    "123,\"text field\",456.78\n",
			expected: []string{"123", "text field", "456.78"},
		},
		{
			name:     "special characters",
			input:    "field1,\"!@#$%^&*()\",field3\n",
			expected: []string{"field1", "!@#$%^&*()", "field3"},
		},
		{
			name:     "unicode characters",
			input:    "field1,\"héllo wörld\",field3\n",
			expected: []string{"field1", "héllo wörld", "field3"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := NewReader(bufio.NewReader(strings.NewReader(tt.input)))
			result, err := reader.Read()

			if err != nil {
				t.Fatalf("Read() error = %v", err)
			}

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Read() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestReader_Read_EOF(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectError bool
		errorType   error
	}{
		{
			name:        "empty input",
			input:       "",
			expectError: true,
			errorType:   io.EOF,
		},
		{
			name:        "only newline",
			input:       "\n",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := NewReader(bufio.NewReader(strings.NewReader(tt.input)))
			result, err := reader.Read()

			if tt.expectError {
				if err == nil {
					t.Fatalf("Expected error, got nil")
				}
				if err != tt.errorType {
					t.Errorf("Expected error %v, got %v", tt.errorType, err)
				}
			} else {
				if err != nil {
					t.Fatalf("Unexpected error: %v", err)
				}
				if result == nil {
					t.Error("Expected non-nil result")
				}
			}
		})
	}
}

func TestReader_Read_MultipleLines(t *testing.T) {
	input := "line1,field2\nline2,field4\nline3,field6\n"
	reader := NewReader(bufio.NewReader(strings.NewReader(input)))

	// Read first line
	result1, err := reader.Read()
	if err != nil {
		t.Fatalf("First Read() error = %v", err)
	}
	expected1 := []string{"line1", "field2"}
	if !reflect.DeepEqual(result1, expected1) {
		t.Errorf("First Read() = %v, want %v", result1, expected1)
	}

	// Read second line
	result2, err := reader.Read()
	if err != nil {
		t.Fatalf("Second Read() error = %v", err)
	}
	expected2 := []string{"line2", "field4"}
	if !reflect.DeepEqual(result2, expected2) {
		t.Errorf("Second Read() = %v, want %v", result2, expected2)
	}

	// Read third line
	result3, err := reader.Read()
	if err != nil {
		t.Fatalf("Third Read() error = %v", err)
	}
	expected3 := []string{"line3", "field6"}
	if !reflect.DeepEqual(result3, expected3) {
		t.Errorf("Third Read() = %v, want %v", result3, expected3)
	}

	// Try to read beyond EOF
	_, err = reader.Read()
	if err != io.EOF {
		t.Errorf("Expected EOF error, got %v", err)
	}
}

func TestReader_Read_EdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "line with only commas",
			input:    ",,,\n",
			expected: []string{"", "", "", ""},
		},
		{
			name:     "line with only quotes",
			input:    "\"\"\n",
			expected: []string{""},
		},
		{
			name:     "single comma",
			input:    ",\n",
			expected: []string{"", ""},
		},
		{
			name:     "trailing comma",
			input:    "field1,\n",
			expected: []string{"field1", ""},
		},
		{
			name:     "leading comma",
			input:    ",field2\n",
			expected: []string{"", "field2"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := NewReader(bufio.NewReader(strings.NewReader(tt.input)))
			result, err := reader.Read()

			if err != nil {
				t.Fatalf("Read() error = %v", err)
			}

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Read() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// Benchmark tests for performance measurement
func BenchmarkReader_Read_SimpleFields(b *testing.B) {
	input := "field1,field2,field3,field4,field5\n"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		reader := NewReader(bufio.NewReader(strings.NewReader(input)))
		_, err := reader.Read()
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkReader_Read_QuotedFields(b *testing.B) {
	input := "\"field1\",\"field2\",\"field3\",\"field4\",\"field5\"\n"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		reader := NewReader(bufio.NewReader(strings.NewReader(input)))
		_, err := reader.Read()
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkReader_Read_ComplexFields(b *testing.B) {
	input := "field1,\"field with, comma\",\"field with \"\"quotes\"\"\",field4\n"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		reader := NewReader(bufio.NewReader(strings.NewReader(input)))
		_, err := reader.Read()
		if err != nil {
			b.Fatal(err)
		}
	}
}
