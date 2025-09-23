package main

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"csvc"
)

func main() {
	// Open the CSV file
	file, err := os.Open("dummy.csv")
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		return
	}
	defer file.Close()

	// Create a buffered reader
	bufReader := bufio.NewReader(file)

	// Create a CSV reader
	csvReader := csvc.NewReader(bufReader)


	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Printf("Error reading record: %v\n", err)
			break
		}

		fmt.Printf("Record: %v with %d fields and last field: %v\n", record, len(record), []byte(record[len(record)-1]))
	}

}
