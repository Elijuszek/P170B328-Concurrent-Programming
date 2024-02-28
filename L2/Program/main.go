package main

import (
	"fmt"
	"time"
)

type Cities struct {
	Cities []City `json:"cities"`
}

type City struct {
	Name       string  `json:"name"`
	Population int     `json:"population"`
	Area       float64 `json:"area"`
}

type ComputedCity struct {
	City City
	Hash uint64
}

func main() {
	readCities := readJson("IFF-1-1_ZekonisE_L1_dat_2.json")
	workerThreads := 10
	dataChan := make(chan City)
	workerChan := make(chan City, len(readCities))
	resultChan := make(chan ComputedCity)
	returnChan := make(chan []ComputedCity)
	timeStart := time.Now()

	go DataFunction(dataChan, workerChan)
	go ResultFunction(resultChan, returnChan, workerThreads)
	for i := 0; i < workerThreads; i++ {
		go WorkerFunction(workerChan, resultChan)
	}

	for i := 0; i < len(readCities); i++ {
		dataChan <- readCities[i]
	}
	close(dataChan)

	computedCities := <-returnChan

	// Printing results
	totalTime := time.Since(timeStart)
	printTable(computedCities)
	writeTableToTxtFile(computedCities, "IFF-1-1_ZekonisE_L1_rez.txt")
	fmt.Printf("\nResult count: %d\n", len(computedCities))
	fmt.Println(totalTime)
}
