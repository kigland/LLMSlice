package LLMSlice

import (
	"sync"
	"sync/atomic"
	"time"
)

// Slice 单供应商切片
type Slice[T any] struct {
	cur   atomic.Int64
	items []*Elem[T]
	rw    sync.RWMutex
}

type Elem[T any] struct {
	val       T
	failure   bool
	holdUntil time.Time
}

func (e *Elem[T]) Available() bool {
	return !e.failure && e.holdUntil.Before(time.Now())
}

func (e *Elem[T]) MarkAsFailure() {
	e.failure = true
}

func (e *Elem[T]) MarkAsHold(duration time.Duration) {
	e.holdUntil = time.Now().Add(duration)
}

func (e *Elem[T]) MarkAsAvailable() {
	e.failure = false
	e.holdUntil = time.Time{}
}

func (e *Elem[T]) Val() T {
	if e == nil {
		var v T
		return v
	}
	return e.val
}

func (s *Slice[T]) Append(item T) {
	s.rw.Lock()
	defer s.rw.Unlock()
	s.items = append(s.items, &Elem[T]{val: item})
}

func (s *Slice[T]) Get() TContext[T] {
	s.rw.RLock()
	defer s.rw.RUnlock()
	for {
		idx := s.cur.Add(1) % int64(len(s.items))
		if s.items[idx].Available() {
			return s.items[idx]
		}
	}
}

func (s *Slice[T]) Len() int {
	return len(s.items)
}
