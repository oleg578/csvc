# Project Overview

This project is a golang library (package) for parsing and writing CSV data.
It provides functionalities to read and write CSV files.
It supports custom delimiters.
It provides CSVReader and CSVWriter structs for handling CSV data.
CSVReader supports reading CSV data from files.
CSVReader provides methods for parsing and validating CSV data and return result as slices of strings.
CSVWriter supports writing CSV data to files.
It provide error handling for malformed CSV data.

## Folder Structure

- `/docs`: Contains the RFC files and design documents.

## Coding Standards

- Use coding standards from [Google's Go Style Guide](https://google.github.io/styleguide/go/)
- Ensure code is well-documented with comments and follows best practices.
- Write unit tests for all functions and methods.
- Write benchmarks for performance-critical code.
- Use Go's built-in testing package for unit tests.
- Use meaningful variable and function names for better readability.
- Use 'go fmt' to format the code consistently.
- Ensure proper error handling and logging.
- Follow semantic versioning for releases.
- Use Go modules for dependency management.
- Ensure compatibility with the latest stable version of Go.

## Links

- [RFC 4180 Specification](https://tools.ietf.org/html/rfc4180)
- [Reference C Implementation (libcsv)](https://github.com/rgamble/libcsv)
- [Go's Built-in CSV Package](https://golang.org/pkg/encoding/csv/)
- [My GitHub Profile](https://github.com/oleg578)

## License

This project is licensed under the MIT License.