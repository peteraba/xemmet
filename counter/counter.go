package counter

import (
	"sync"
)

var (
	// nolint:gochecknoglobals
	globalTabStopCounter = 0
	// nolint:gochecknoglobals
	globalTabStopMutex = &sync.Mutex{}
)

func ResetGlobalTabStopCounter() {
	globalTabStopMutex.Lock()
	defer globalTabStopMutex.Unlock()

	globalTabStopCounter = 0
}

func GetGlobalTabStopCounter() int {
	globalTabStopMutex.Lock()
	defer globalTabStopMutex.Unlock()

	counter := globalTabStopCounter

	globalTabStopCounter++

	return counter
}
