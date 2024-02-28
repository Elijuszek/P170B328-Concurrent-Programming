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

func printResultMonitor(resultMonitor ResultMonitor) {
	fmt.Println("---------------------------------------------------")
	for _, compCity := range resultMonitor.getItems() {
		fmt.Printf("%s %d %g %d\n", compCity.City.Name, compCity.City.Population, compCity.City.Area, compCity.Hash)
	}
	fmt.Println("---------------------------------------------------")
}

func printResultMonitorJson(fileName string, resultMonitor ResultMonitor) {
	file, err := os.Create(fileName)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer file.Close()

	cities := resultMonitor.computedCities

	byteValue, err := json.MarshalIndent(cities, "", "  ")
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = file.Write(byteValue)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func generateHashCode(data City) uint64 {
	// Create a new FNV-1a hash
	h := fnv.New64a()

	// Running fibonacci to make it more complex
	FibonacciRecursion(30)

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
