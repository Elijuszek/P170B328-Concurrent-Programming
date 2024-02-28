package main

import (
	"fmt"
	"math"
)

func main() {
	data := make(chan float64)
	avgLow := make(chan float64)
	avgMid := make(chan float64)
	avgHigh := make(chan float64)
	done := make(chan bool)

	for i := 0; i < 10; i++ {
		go Sender(data, i)
	}
	go Receiver(data, avgLow, avgMid, avgHigh)
	go Printer(avgLow, done)
	go Printer(avgMid, done)
	go Printer(avgHigh, done)

	<-done
	<-done
	<-done
}

func Sender(data chan<- float64, i int) {
	data <- float64(i)
	for {
		i = int(math.Pow(float64(i), 2) - 4*float64(i) + 1) // Assign the new value to the outer 'i'
		data <- float64(i)
	}
}

func Receiver(data <-chan float64, avgLow chan<- float64, avgMid chan<- float64, avgHigh chan<- float64) {
	array := [15]float64{}
	index := 0
	average := 0.0
	for average < 200 {
		average = 0.0
		array[index] = <-data
		index = (index + 1) % 15
		for i := 0; i < 15; i++ {
			average += array[i]
		}

		average /= 15
		if average < 10 {
			avgLow <- average
		} else if average <= 100 && average >= 0 {
			avgMid <- average
		} else if average <= 200 && average >= 75 {
			avgHigh <- average
		}

	}
	close(avgLow)
	close(avgMid)
	close(avgHigh)
}

func Printer(channel <-chan float64, done chan<- bool) {
	for value := range channel {
		fmt.Println(value)
	}
	done <- true
}
