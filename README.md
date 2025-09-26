# 📄 CSVC - CSV Parser Library for Go

A high-performance, RFC 4180 compliant CSV parsing library written in Go. CSVC focuses on predictable throughput and minimal allocations, making it a practical replacement for the standard `encoding/csv` reader when processing large datasets or streaming CSV payloads.

## Features

- Handles RFC 4180 constructs: quoted fields, embedded delimiters, CRLF endings, and multi-line values.
- Supports custom field delimiters (`Reader.Comma`) and quote characters (`Reader.Quote`).
- Minimises allocations via reusable record buffers (`Reader.ReuseRecord`) and a zero-copy string builder.
- Surfaces structured parse errors with line and column numbers (`ParseError`).
- Provides benchmarks and extensive tests to guard correctness and performance.

## Windsurf timing

