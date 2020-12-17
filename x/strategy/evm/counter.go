package evm

import "sync"

type Counter struct {
	cur int
	sum int
	lock *sync.RWMutex
}

func NewCounter(num int) *Counter {
	return &Counter{
		cur: 0,
		sum: num,
		lock: new(sync.RWMutex),
	}
}

func (c *Counter) Add() int {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.cur++
	cur := c.cur
	return cur
}

func (c *Counter) IsOver() bool {
	if c.cur >= c.sum {
		return true
	} else {
		return false
	}
}

func (c *Counter) GetCurrentNum() int {
	return c.cur
}