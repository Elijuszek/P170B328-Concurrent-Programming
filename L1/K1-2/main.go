package main

import (
	"fmt"
	"sync"
	"time"
)

type DataMonitor struct {
	c         int
	d         int
	readCount int
	printed   bool
	mu        sync.Mutex
	cond      *sync.Cond
}

func NewDataMonitor() *DataMonitor {

	dm := &DataMonitor{
		c:         10,
		d:         100,
		readCount: 0,
		printed:   false,
	}
	dm.cond = sync.NewCond(&dm.mu)
	return dm
}

func (dm *DataMonitor) Change() {
	dm.mu.Lock()
	for dm.readCount < 2 {
		dm.cond.Wait()
	}
	defer dm.mu.Unlock()
	dm.c++
	dm.d--
	dm.readCount = 0
	dm.printed = false
}

func (dm *DataMonitor) Print() {
	dm.mu.Lock()
	defer dm.mu.Unlock()
	dm.readCount++
	if !dm.printed {
		fmt.Printf("c: %d d: %d\n", dm.c, dm.d)
		dm.printed = true
	} else {
		fmt.Print("read\n")
	}
	dm.cond.Broadcast()
}

func functionPrint(dm *DataMonitor) {
	for true {
		dm.Print()
		time.Sleep(time.Millisecond * 1000)
	}
}

func functionWrite(dm *DataMonitor) {
	for true {
		dm.Change()
	}
}

func main() {
	dm := NewDataMonitor()
	writeThreads := 2
	readThreads := 3
	var wg sync.WaitGroup
	wg.Add(writeThreads + readThreads)
	for i := 0; i < readThreads; i++ {
		go functionPrint(dm)
	}
	for i := 0; i < writeThreads; i++ {
		go functionWrite(dm)
	}
	wg.Wait()
}
