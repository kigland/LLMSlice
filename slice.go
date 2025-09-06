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
	if e == nil {
		return false
	}
	return !e.failure && e.holdUntil.Before(time.Now())
}

func (e *Elem[T]) MarkAsFailure() {
	if e == nil {
		return
	}
	e.failure = true
}

func (e *Elem[T]) MarkAsHold(duration time.Duration) {
	if e == nil {
		return
	}
	e.holdUntil = time.Now().Add(duration)
}

func (e *Elem[T]) MarkAsAvailable() {
	if e == nil {
		return
	}
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

func (s *Slice[T]) availableCount() int {
	count := 0
	for _, item := range s.items {
		if item.Available() {
			count++
		}
	}
	return count
}

func (s *Slice[T]) AvailableCount() int {
	s.rw.RLock()
	defer s.rw.RUnlock()
	return s.availableCount()
}

func (s *Slice[T]) Append(items ...T) *Slice[T] {
	s.rw.Lock()
	defer s.rw.Unlock()
	for _, item := range items {
		s.items = append(s.items, &Elem[T]{val: item})
	}
	return s
}

func (s *Slice[T]) Get() TContext[T] {
	s.rw.RLock()
	defer s.rw.RUnlock()
	if len(s.items) == 0 {
		return nil
	}
	for {
		if s.availableCount() == 0 {
			return nil
		}
		if e := s.getOrNil(); e != nil {
			return e
		}
	}
}
func (s *Slice[T]) getOrNil() *Elem[T] {
	if len(s.items) == 0 {
		return nil
	}
	idx := s.cur.Add(1) % int64(len(s.items))
	e := s.items[idx]
	if e.Available() {
		return e
	}
	return nil
}

func (s *Slice[T]) MustGet(sleep time.Duration) TContext[T] {
	s.rw.RLock()
	defer s.rw.RUnlock()
	if len(s.items) == 0 {
		return nil
	}
	for {
		if s.availableCount() == 0 {
			time.Sleep(sleep)
			continue
		}
		e := s.getOrNil()
		if e != nil {
			return e
		}
	}
}

func (s *Slice[T]) Len() int {
	return len(s.items)
}
