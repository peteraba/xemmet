package main

import (
	"sync"
)

type Counter struct {
	counter int
	lock    *sync.Mutex
}

func NewCounter() *Counter {
	return &Counter{
		counter: 1,
		lock:    &sync.Mutex{},
	}
}

func (c *Counter) Reset() {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.counter = 1
}

func (c *Counter) Get() int {
	c.lock.Lock()
	defer c.lock.Unlock()

	counter := c.counter

	c.counter++

	return counter
}
