# CSV Parser Comparison Report

Date: 2025-09-23
Environment: Windows 11, Go (goos=windows, goarch=amd64)
CPU: 11th Gen Intel(R) Core(TM) i5-1135G7 @ 2.40GHz
Branch: v2

This report compares performance of this library (CSVC) vs Go's built-in `encoding/csv` using the in-repo benchmarks.

## Benchmarks

Command used:

```powershell
go test -bench=Comparison -benchmem
# For MediumQuoted specifically
go test -bench=Comparison_MediumQuoted -benchmem
```

Datasets:

- SmallSimple: 100 rows × 5 columns, unquoted
- SmallQuoted: 100 rows × 5 columns, quoted
- MediumSimple: 1000 rows × 10 columns, unquoted
- MediumQuoted: 1000 rows × 10 columns, quoted

## Results

| Benchmark | Library | ns/op | B/op | allocs/op |
|---|---|---:|---:|---:|
| SmallSimple | CSVC | 21,600 | 24,424 | 109 |
| SmallSimple | Go builtin | 21,271 | 17,608 | 217 |
| SmallQuoted | CSVC | 23,150 | 24,408 | 109 |
| SmallQuoted | Go builtin | 22,319 | 17,608 | 217 |
| MediumSimple | CSVC | 275,020 | 129,920 | 1,009 |
| MediumSimple | Go builtin | 345,628 | 261,296 | 2,020 |
| MediumQuoted | CSVC | 336,299 | 115,505 | 1,009 |
| MediumQuoted | Go builtin | 327,792 | 261,296 | 2,020 |

## Summary

- SmallSimple: Go slightly faster (~1.5%), CSVC uses fewer allocations.
- SmallQuoted: Go slightly faster (~3.7%), CSVC uses fewer allocations.
- MediumSimple: CSVC faster (~20.4%), with ~50% fewer allocations and lower B/op.
- MediumQuoted: Go slightly faster (~2.6%), CSVC still much lower allocations and B/op.

## Notes

- CSVC now supports multi-line quoted fields (CR/LF within quoted fields) and very large fields (> 8KB).
- CSVC preserves whitespace (no trimming) and follows RFC 4180 quoting rules (including escaped quotes).
