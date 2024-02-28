package main

import (
	"fmt"
	"sync"
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
	readCities := readJson("IFF1-1_ZekonisElijus_L1a_dat_2.json")
	threadCount := 10
	finishedReading := false
	var dataMonitor DataMonitor
	var resultMonitor ResultMonitor

	dataMonitor.AllocateDataMonitor(50)

	var wg sync.WaitGroup
	wg.Add(threadCount)

	// Start threads
	for i := 0; i < threadCount; i++ {
		go ParallelFunction(&dataMonitor, &resultMonitor, &wg, &finishedReading)
	}

	start := time.Now()

	for _, city := range readCities {
		dataMonitor.addItem(city)
	}

	finishedReading = true

	wg.Wait()
	duration := time.Since(start)

	// print
	printResultMonitor(resultMonitor)
	printResultMonitorJson("IFF-1-1_ZekonisElijus_L1a_rez.json", resultMonitor)
	fmt.Println(duration)
}
