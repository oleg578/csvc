# CSV Reader Benchmarks

This file contains comprehensive benchmarks for the CSV Reader's `Read` function.

## Running Benchmarks

To run all benchmarks:

```bash
go test -bench=. -benchmem
```

To run specific benchmark categories:

```bash
# Single record benchmarks
go test -bench=BenchmarkReader_Read_Single -benchmem

# Small dataset benchmarks
go test -bench=BenchmarkReader_Read_Small -benchmem

# Medium dataset benchmarks
go test -bench=BenchmarkReader_Read_Medium -benchmem

# Large dataset benchmarks
go test -bench=BenchmarkReader_Read_Large -benchmem

# Complex field benchmarks
go test -bench=BenchmarkReader_Read_Complex -benchmem
```

## Benchmark Categories

### Dataset Size Benchmarks

- **Small**: 100 rows, 5 columns
- **Medium**: 1,000 rows, 10 columns
- **Large**: 10,000 rows, 20 columns

### Field Type Benchmarks

- **Simple**: Basic unquoted fields
- **Quoted**: All fields enclosed in quotes
- **Complex**: Mixed quoted/unquoted with special characters
- **Escaped**: Fields with escaped quotes
- **Multiline**: Fields containing newlines

### Special Case Benchmarks

- **SingleRecord**: Single CSV row performance
- **EmptyFields**: CSV with many empty fields
- **LongFields**: CSV with very long field content
- **ManyColumns**: CSV with 50 columns
- **CustomDelimiter**: CSV with semicolon delimiter

## Performance Interpretation

Benchmark output format:

```bash
BenchmarkName-8    iterations    ns/operation    bytes/op    allocs/op
```

- **iterations**: Number of times the benchmark ran
- **ns/operation**: Nanoseconds per operation (lower is better)
- **bytes/op**: Bytes allocated per operation (lower is better)
- **allocs/op**: Number of allocations per operation (lower is better)

## Optimization Notes

The benchmarks help identify:

1. Memory allocation patterns
2. Performance differences between quoted vs unquoted fields
3. Scalability with dataset size
4. Impact of field complexity on performance
5. Custom delimiter performance overhead
