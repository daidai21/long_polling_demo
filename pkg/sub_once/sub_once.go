package sub_once

import "sync"

type SubOnce struct {
	ch     []chan struct{}
	mu     sync.Mutex
	closed bool
}

func New() *SubOnce {
	return &SubOnce{
		ch:     make([]chan struct{}, 0),
		mu:     sync.Mutex{},
		closed: false,
	}
}

func (s *SubOnce) Sub() chan struct{} {
	defer s.mu.Unlock()
	s.mu.Lock()
	if s.closed {
		return nil
	}

	c := make(chan struct{}, 1)
	s.ch = append(s.ch, c)
	return c
}

func (s *SubOnce) Pub() {
	defer s.mu.Unlock()
	s.mu.Lock()
	if s.closed {
		return
	}

	for _, c := range s.ch {
		c <- struct{}{}
	}
	s.ch = nil
}

func (s *SubOnce) Close() {
	defer s.mu.Unlock()
	s.mu.Lock()
	if s.closed {
		return
	}

	for _, c := range s.ch {
		close(c)
	}
	s.ch = nil
}
