package LLMSlice

import (
	"sync"
	"sync/atomic"
)

// Slice 单供应商切片
type Slice[T any] struct {
	cur   atomic.Int64
	items []T
	rw    sync.RWMutex
}

func (s *Slice[T]) Append(item T) {
	s.rw.Lock()
	defer s.rw.Unlock()
	s.items = append(s.items, item)
}

func (s *Slice[T]) Get() T {
	s.rw.RLock()
	defer s.rw.RUnlock()
	val := s.cur.Add(1)
	return s.items[val]
}

func (s *Slice[T]) Len() int {
	return len(s.items)
}
