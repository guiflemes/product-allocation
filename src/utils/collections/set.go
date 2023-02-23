package collections

import (
	"math/rand"
	"sync"
)

type Set[T comparable] struct {
	m     sync.RWMutex
	items map[T]bool
	keys  []T
}

func (s *Set[T]) Add(item T) {
	s.m.Lock()
	defer s.m.Unlock()
	if _, ok := s.items[item]; !ok {
		s.keys = append(s.keys, item)
		s.items[item] = true
	}
}

func (s *Set[T]) Iter() <-chan T {
	ch := make(chan T)
	s.m.RLock()

	go func() {
		defer close(ch)
		defer s.m.RUnlock()
		for _, item := range s.keys {
			ch <- item
		}
	}()

	return ch
}

func (s *Set[T]) Pop() T {
	var r T
	s.m.Lock()
	defer s.m.Unlock()

	if len(s.keys) == 0 {
		return r
	}

	index := rand.Intn(len(s.keys))
	item := s.keys[index]
	delete(s.items, item)
	s.keys[index] = s.keys[len(s.keys)-1]
	s.keys = s.keys[:len(s.keys)-1]
	return item

}

func (s *Set[T]) Len() int {
	s.m.RLock()
	defer s.m.RUnlock()
	return len(s.items)
}

func NewSet[T comparable]() *Set[T] {
	return &Set[T]{
		items: make(map[T]bool),
	}
}
