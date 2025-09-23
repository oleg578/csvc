# CSVC Optimization Report

Date: 2025-09-23
Environment: Windows 11, Go (goos=windows, goarch=amd64)
CPU: 11th Gen Intel(R) Core(TM) i5-1135G7 @ 2.40GHz
Branch: v2

## Summary

We implemented several performance optimizations to the CSV reader while preserving RFC 4180 behavior and not trimming whitespace:

- Single output buffer with indexed fields to minimize per-field allocations
- Unquoted fast path that splits directly from the input buffer and constructs a single string per record
- Quoted path optimized with bytes.IndexByte to jump between quotes, reducing branchy per-byte loops
- Fixed empty-field and CR-only line ending handling across paths

All unit tests pass after the changes. Benchmarks below compare our library (CSVC) against Go's built-in `encoding/csv` on the same datasets.

## Benchmark Comparison

Numbers are from `go test -bench=Comparison -benchmem` variants, run on the machine above. Lower is better for ns/op and B/op.

| Benchmark | Library | ns/op | B/op | allocs/op |
|---|---|---:|---:|---:|
| SmallSimple (100x5, unquoted) | CSVC | 17,804 | 24,424 | 109 |
| SmallSimple (100x5, unquoted) | Go builtin | 20,553 | 17,608 | 217 |
| SmallQuoted (100x5, quoted) | CSVC | 18,319 | 24,408 | 109 |
| SmallQuoted (100x5, quoted) | Go builtin | 21,072 | 17,608 | 217 |
| MediumSimple (1000x10, unquoted) | CSVC | 216,067 | 129,922 | 1,009 |
| MediumSimple (1000x10, unquoted) | Go builtin | 328,336 | 261,296 | 2,020 |
| MediumQuoted (1000x10, quoted) | CSVC | 230,456 | 115,507 | 1,009 |
| MediumQuoted (1000x10, quoted) | Go builtin | 328,843 | 261,296 | 2,020 |

Notes:

- Our implementation now consistently outperforms the builtin in ns/op across these datasets.
- We also reduce allocations per operation significantly. B/op is higher than builtin on smaller datasets due to our reader buffer strategy, but much lower on medium datasets.

## Key Changes (Technical)

- Reader struct now includes:
  - `out []byte`, `starts []int`, `ends []int` to accumulate one record and reference fields by indices
- Unquoted fast path:
  - Detects absence of '"'; scans once collecting field boundaries
  - Builds a single string from the line and produces sub-slices for fields
  - Trims trailing `\r` for CR-only endings
- Quoted parsing:
  - State machine remains RFC 4180-compliant
  - In QUOTED state, uses `bytes.IndexByte` to jump to next '"' and append intervening spans in one shot
  - Correctly handles doubled quotes (""") and commas/newlines inside quoted fields
- Buffer management:
  - Output buffer pre-sized to line length when needed to reduce growth
  - Improved handling of trailing commas for empty fields

## Future Work

- Reader buffer sizing: allocate based on observed line length rather than fixed caps in NewReader to reduce B/op on small datasets when many readers are created.
- Consider `bufio.Reader.ReadSlice('\n')` to avoid the extra copy into `b.buf` (ensure CR-only handling remains correct).
- Optional micro-optimizations to further reduce branches in the hot loops.

## How to Reproduce

Run the specific benchmark groups:

```powershell
# Small datasets
go test -bench=Comparison_SmallSimple -benchmem
go test -bench=Comparison_SmallQuoted -benchmem

# Medium datasets
go test -bench=Comparison_MediumSimple -benchmem
go test -bench=Comparison_MediumQuoted -benchmem
```

All unit tests:

```powershell
go test ./...
```
