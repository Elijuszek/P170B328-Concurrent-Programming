package main

import (
	"sync"
)

// TODO: all lock unlock cond funbctionality must be implemented here:
type DataMonitor struct {
	Chars string
	mutex sync.Mutex
	cond  *sync.Cond
}

func (dm *DataMonitor) AllocateDataMonitor() {
	dm.Chars = "*"
	dm.cond = sync.NewCond(&dm.mutex)
}

func (dm *DataMonitor) addItem(ch rune) {
	dm.mutex.Lock()
	defer dm.mutex.Unlock()
	dm.Chars = dm.Chars + string(ch)
}
