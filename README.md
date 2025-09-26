# 📄 CSVC - CSV Parser Library for Go

A high-performance, RFC 4180 compliant CSV parsing library written in Go. CSVC focuses on predictable throughput and minimal allocations, making it a practical replacement for the standard `encoding/csv` reader when processing large datasets or streaming CSV payloads.

## Features

- Handles RFC 4180 constructs: quoted fields, embedded delimiters, CRLF endings, and multi-line values.
- Supports custom field delimiters (`Reader.Comma`) and quote characters (`Reader.Quote`).
- Minimises allocations via reusable record buffers (`Reader.ReuseRecord`) and a zero-copy string builder.
- Surfaces structured parse errors with line and column numbers (`ParseError`).
- Provides benchmarks and extensive tests to guard correctness and performance.

## chatGPT Codex

The time it took for the Code to produce the final result was about 45 minutes.

## API Overview

### `func NewReader(src *io.Reader) *Reader`

Creates a new reader backed by an aggressively sized internal buffer. Passing a pointer keeps the underlying source swappable while avoiding unnecessary interface boxing.

### `func (r *Reader) Read() ([]string, error)`

Parses the next CSV record and returns it as a slice of strings. When the input is exhausted, it returns `io.EOF`. If `ReuseRecord` is enabled, the returned slice is reused on subsequent calls—copy data you need before calling `Read` again.

### `type ParseError`

```go
type ParseError struct {
    Line   int // 1-based line number where parsing failed
    Column int // 1-based column within the current record
    Err    error // one of ErrBareQuote, ErrUnterminatedQuote, or an I/O error
}
```

`ParseError` implements `error` and supports `errors.Is` / `errors.As`.

### Sentinel Errors

- `ErrBareQuote`: encountered an unexpected quote in an unquoted field.
- `ErrUnterminatedQuote`: hit EOF while a quoted field remained open.

## Usage Example

```go
package main

import (
    "fmt"
    "io"
    "strings"

    "csvc"
)

func main() {
    var src io.Reader = strings.NewReader("name,price\nWidget,12.50\n")
    reader := csvc.NewReader(&src)

    for {
     record, err := reader.Read()
     if err == io.EOF {
      break
     }
     if err != nil {
      panic(err)
     }
     fmt.Println(record)
    }
}
```

For larger workloads, enable `reader.ReuseRecord = true` to reduce allocations—remember to copy fields before the next `Read` call if you need them later.

## Benchmarks

Run the bundled benchmarks to compare CSVC with the standard library reader:

```bash
go test -run=^$ -bench=. -benchmem
```

Sample output on an 11th Gen Intel i5 (Windows):

```text
BenchmarkReader-8        259594 ns/op   169.62 MB/s   66622 B/op       5 allocs/op
BenchmarkEncodingCSV-8   232113 ns/op   189.70 MB/s  135801 B/op    2063 allocs/op
```

The standard reader remains slightly faster in raw throughput, while CSVC reduces heap allocations by several orders of magnitude—helpful for GC-bound workloads.

## Testing

All behaviour is covered by table-driven unit tests. Execute the suite with:

```bash
go test ./...
```

## Development Notes

This code base is generated and maintained using ChatGPT Codex tooling. Contributions and manual edits are welcome—please include tests for new features.

## License

Released under the MIT License. See the full text below:

```text
MIT License

Copyright (c) 2025 Oleg

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
```
