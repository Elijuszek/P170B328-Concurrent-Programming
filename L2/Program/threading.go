package main

func DataFunction(dataChannel <-chan City, workerChannel chan<- City) {
	var dataArray [10]City
	length := 0
	defer close(workerChannel)
	for {
		if length != 10 && length != 0 {
			select {
			case data := <-dataChannel:
				dataArray[length] = data
				length++
			case workerChannel <- dataArray[length-1]:
				length--
			}
		} else if length == 10 {
			workerChannel <- dataArray[length-1]
			length--
		} else {
			data, ok := <-dataChannel
			if !ok {
				break
			}
			dataArray[length] = data
			length++
		}
	}
}

func WorkerFunction(workerChannel <-chan City, resultChannel chan<- ComputedCity) {
	finishValue := ComputedCity{City{"", 0, 0.0}, 0}
	for c := range workerChannel {
		if c.Population >= 20000 {
			resultChannel <- ComputedCity{c, generateHashCode(c)}
		}
	}
	resultChannel <- finishValue
}

func ResultFunction(resultChannel <-chan ComputedCity, returnChan chan<- []ComputedCity, activeThreads int) {
	var resultArray []ComputedCity
	for activeThreads != 0 {
		c := <-resultChannel
		if c.City.Name == "" {
			activeThreads--
			continue
		}
		resultArray = append(resultArray, c)
	}
	returnChan <- resultArray
}
