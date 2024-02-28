package main

import (
	"sync"
)

func ParallelFunction(dataMonitor *DataMonitor, resultMonitor *ResultMonitor, wg *sync.WaitGroup, finishedReading *bool) {
	defer wg.Done()
	for {
		var city = dataMonitor.removeItem()
		if city.Name == "" {
			if *finishedReading {
				return
			}
		} else {
			var computedCity ComputedCity
			computedCity.City = city
			computedCity.Hash = generateHashCode(city)
			resultMonitor.addItem(computedCity)
		}
	}
}

func SingleThreadFunction(dataMonitor *DataMonitor, resultMonitor *ResultMonitor) {
	for len(dataMonitor.Cities) > 0 {
		var city = dataMonitor.removeItem()
		var computedCity ComputedCity
		computedCity.City = city
		computedCity.Hash = generateHashCode(city)
		resultMonitor.addItem(computedCity)
	}
}
