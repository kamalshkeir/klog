package klog

import (
	"fmt"
	"sync"
)

type LimitedSlice[T any] struct {
	Slice []T
	max   int
	mu    sync.RWMutex
}

func NewLimitedSlice[T any](max int) *LimitedSlice[T] {
	return &LimitedSlice[T]{
		max: max,
	}
}

func (ls *LimitedSlice[T]) Add(element T) {
	ls.mu.Lock()
	defer ls.mu.Unlock()

	if len(ls.Slice) == ls.max {
		ls.Slice = ls.Slice[1:]
	}
	ls.Slice = append(ls.Slice, element)
}

func (ls *LimitedSlice[T]) Get(index int) (T, error) {
	ls.mu.RLock()
	defer ls.mu.RUnlock()

	if index >= 0 && index < len(ls.Slice) {
		return ls.Slice[index], nil
	}
	return *new(T), fmt.Errorf("index out of range")
}

func (ls *LimitedSlice[T]) Delete(index int) error {
	ls.mu.Lock()
	defer ls.mu.Unlock()
	if index >= 0 && index < len(ls.Slice) {
		ls.Slice = append(ls.Slice[:index], ls.Slice[index+1:]...)
		return nil
	}
	return fmt.Errorf("index out of range")
}

func (ls *LimitedSlice[T]) Range(fn func(T, int) bool) {
	ls.mu.RLock()
	defer ls.mu.RUnlock()

	for i, element := range ls.Slice {
		if !fn(element, i) {
			break
		}
	}
}
