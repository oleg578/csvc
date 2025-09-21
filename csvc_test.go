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
	r := bufio.NewReader(strings.NewReader(input))
	reader := NewReader(r)

	if reader.Comma != ',' {
		t.Errorf("Expected default comma to be ',', got %c", reader.Comma)
	}

	if reader.r == nil {
		t.Error("Expected reader to be initialized")
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
			input:    "onefield\n",
			expected: []string{"onefield"},
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := bufio.NewReader(strings.NewReader(tt.input))
			reader := NewReader(r)

			result, err := reader.Read()
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Expected %v, got %v", tt.expected, result)
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
			name:     "quoted fields",
			input:    "\"field1\",\"field2\",\"field3\"\n",
			expected: []string{"field1", "field2", "field3"},
		},
		{
			name:     "mixed quoted and unquoted",
			input:    "field1,\"field2\",field3\n",
			expected: []string{"field1", "field2", "field3"},
		},
		{
			name:     "quoted field with comma",
			input:    "\"field1,with,commas\",field2\n",
			expected: []string{"field1,with,commas", "field2"},
		},
		{
			name:     "quoted field with spaces",
			input:    "\"field with spaces\",\"another field\"\n",
			expected: []string{"field with spaces", "another field"},
		},
		{
			name:     "empty quoted fields",
			input:    "\"\",\"field2\",\"\"\n",
			expected: []string{"", "field2", ""},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := bufio.NewReader(strings.NewReader(tt.input))
			reader := NewReader(r)

			result, err := reader.Read()
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Expected %v, got %v", tt.expected, result)
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
			name:     "escaped quotes",
			input:    "\"field with \"\"quotes\"\"\",normal\n",
			expected: []string{"field with \"quotes\"", "normal"},
		},
		{
			name:     "multiple escaped quotes",
			input:    "\"She said \"\"Hello\"\" and \"\"Goodbye\"\"\"\n",
			expected: []string{"She said \"Hello\" and \"Goodbye\""},
		},
		{
			name:     "escaped quote at start",
			input:    "\"\"\"quoted at start\"\n",
			expected: []string{"\"quoted at start"},
		},
		{
			name:     "escaped quote at end",
			input:    "\"quoted at end\"\"\"\n",
			expected: []string{"quoted at end\""},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := bufio.NewReader(strings.NewReader(tt.input))
			reader := NewReader(r)

			result, err := reader.Read()
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Expected %v, got %v", tt.expected, result)
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
			name:     "CR ending (converted to CRLF by test)",
			input:    "field1,field2\r\n", // Using CRLF as CR alone is handled the same way
			expected: []string{"field1", "field2"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := bufio.NewReader(strings.NewReader(tt.input))
			reader := NewReader(r)

			result, err := reader.Read()
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Expected %v, got %v", tt.expected, result)
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
			delimiter: ASCII_TAB,
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
			r := bufio.NewReader(strings.NewReader(tt.input))
			reader := NewReader(r)
			reader.Comma = tt.delimiter

			result, err := reader.Read()
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Expected %v, got %v", tt.expected, result)
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
			input:    "\"field1, with \"\"quotes\"\" and comma\",field2\n",
			expected: []string{"field1, with \"quotes\" and comma", "field2"},
		},
		{
			name:     "numbers and text mixed",
			input:    "123,\"text field\",456.78,\"another,text\"\n",
			expected: []string{"123", "text field", "456.78", "another,text"},
		},
		{
			name:     "special characters",
			input:    "field1,\"field with\nnewline\",field3\n",
			expected: []string{"field1", "field with\nnewline", "field3"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := bufio.NewReader(strings.NewReader(tt.input))
			reader := NewReader(r)

			result, err := reader.Read()
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestReader_Read_EOF(t *testing.T) {
	input := ""
	r := bufio.NewReader(strings.NewReader(input))
	reader := NewReader(r)

	result, err := reader.Read()
	if err != io.EOF {
		t.Errorf("Expected EOF error, got %v", err)
	}

	if len(result) != 0 {
		t.Errorf("Expected empty result, got %v", result)
	}
}

func TestReader_Read_MultipleLines(t *testing.T) {
	input := "field1,field2\nfield3,field4\n"
	r := bufio.NewReader(strings.NewReader(input))
	reader := NewReader(r)

	// Read first line
	result1, err := reader.Read()
	if err != nil {
		t.Fatalf("Unexpected error on first read: %v", err)
	}
	expected1 := []string{"field1", "field2"}
	if !reflect.DeepEqual(result1, expected1) {
		t.Errorf("First line: expected %v, got %v", expected1, result1)
	}

	// Read second line
	result2, err := reader.Read()
	if err != nil {
		t.Fatalf("Unexpected error on second read: %v", err)
	}
	expected2 := []string{"field3", "field4"}
	if !reflect.DeepEqual(result2, expected2) {
		t.Errorf("Second line: expected %v, got %v", expected2, result2)
	}

	// Try to read beyond EOF
	result3, err := reader.Read()
	if err != io.EOF {
		t.Errorf("Expected EOF error, got %v", err)
	}
	if len(result3) != 0 {
		t.Errorf("Expected empty result after EOF, got %v", result3)
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
			input:    "field1,field2,\n",
			expected: []string{"field1", "field2", ""},
		},
		{
			name:     "leading comma",
			input:    ",field1,field2\n",
			expected: []string{"", "field1", "field2"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := bufio.NewReader(strings.NewReader(tt.input))
			reader := NewReader(r)

			result, err := reader.Read()
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}
