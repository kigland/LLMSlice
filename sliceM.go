package LLMSlice

import "sync"

// SliceM 多供应商切片
type SliceM[T any] struct {
	DefaultKey string
	items      map[string]*Slice[T]
	rw         sync.RWMutex
}

func (s *SliceM[T]) Append(key string, item T) {
	s.rw.Lock()
	defer s.rw.Unlock()
	if s.items == nil {
		s.items = make(map[string]*Slice[T])
	}
	if key == "" {
		key = s.DefaultKey
	}
	if s.items[key] == nil {
		s.items[key] = &Slice[T]{}
	}
	s.items[key].Append(item)
}

func (s *SliceM[T]) Get(key string) T {
	if key == "" {
		key = s.DefaultKey
	}
	s.rw.RLock()
	defer s.rw.RUnlock()
	if s.items[key] == nil {
		var v T
		return v
	}
	return s.items[key].Get()
}
