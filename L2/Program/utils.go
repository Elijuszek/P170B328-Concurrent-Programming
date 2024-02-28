package main

import (
	"encoding/json"
	"fmt"
	"hash/fnv"
	"io"
	"os"
)

func readJson(fileName string) []City {
	// Open the JSON file
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	// Read the file content
	byteResult, _ := io.ReadAll(file)
	var cities Cities
	err2 := json.Unmarshal(byteResult, &cities)

	if err2 != nil {
		fmt.Println("Error while decoding the data")
	}
	return cities.Cities
}

// for printing
func printTable(cities []ComputedCity) {
	// Print table header
	fmt.Printf("%-15s%-15s%-15s%-15s\n", "Name", "Population", "Area", "Hash")

	// Iterate over the cities and print data in a tabular format
	for _, comp := range cities {
		fmt.Printf("%-15s%-15d%-15.2f%-15d\n", comp.City.Name, comp.City.Population, comp.City.Area, comp.Hash)
	}
}

func writeTableToTxtFile(cities []ComputedCity, filename string) error {
	// Open the file for writing
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write table header to the file
	fmt.Fprintf(file, "%-15s%-15s%-15s%-15s\n", "Name", "Population", "Area", "Hash")

	// Iterate over the cities and write data in a tabular format to the file
	for _, comp := range cities {
		fmt.Fprintf(file, "%-15s%-15d%-15.2f%-15d\n", comp.City.Name, comp.City.Population, comp.City.Area, comp.Hash)
	}

	return nil
}

func generateHashCode(data City) uint64 {
	// Create a new FNV-1a hash
	h := fnv.New64a()

	// Running fibonacci to make it more complex
	_ = FibonacciRecursion(40)

	// Encode the struct fields into the hash
	h.Write([]byte(fmt.Sprintf("%d%f%s", data.Population, data.Area, data.Name)))

	// Return the hash code as a uint64
	return h.Sum64()
}

func FibonacciRecursion(n int) int {
	if n <= 1 {
		return n
	}
	return FibonacciRecursion(n-1) + FibonacciRecursion(n-2)
}
