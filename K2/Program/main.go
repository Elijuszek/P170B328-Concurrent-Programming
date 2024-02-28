package main

import (
	"fmt"
)

func main() {
	data := make(chan int)
	even := make(chan int)
	odd := make(chan int)
	done := make(chan bool)
	// Start senders and receiver
	go Sender(data, 0, 10)
	go Sender(data, 11, 21)
	go Receiver(data, even, odd, 2)

	// Start printers
	go Printer(even, done)
	go Printer(odd, done)

	<-done
	<-done
}

func Sender(data chan<- int, start int, finish int) {
	for i := start; i <= finish; i++ {
		data <- i
	}
	data <- -1
}

func Receiver(data <-chan int, even chan<- int, odd chan<- int, senders int) {
	count := 0
	for count < senders {
		number := <-data
		if number == -1 {
			count++
			continue
		} else if number%2 == 0 {
			even <- number
		} else {
			odd <- number
		}
	}
	close(even)
	close(odd)
}

func Printer(channel <-chan int, done chan<- bool) {
	numbers := make([]int, 11)
	count := 0
	for number := range channel {
		numbers[count] = number
		count++
	}
	for i := 0; i <= count-1; i++ {
		fmt.Println(numbers[i])
	}
	done <- true
}
