package pipes

import "sync"

type Cleanup struct {
	E    chan struct{}
	wait *sync.WaitGroup
}

func NewCleanup() *Cleanup {
	return &Cleanup{E: make(chan struct{}), wait: &sync.WaitGroup{}}
}

func (c *Cleanup) Cleanup() {
	close(c.E)
	c.wait.Wait()
}

func (c *Cleanup) Add() {
	c.wait.Add(1)
}

func (c *Cleanup) Done() {
	c.wait.Done()
}
