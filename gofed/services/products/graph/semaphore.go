package graph

import (
	"context"
	"sync"
	"time"
)

type Semaphore struct {
	permits chan struct{}
	mu      sync.RWMutex
	max     int
	current int
}

// NewSemaphore creates a new semaphore with the maximum number of permits
func NewSemaphore(maxConcurrent int) *Semaphore {
	if maxConcurrent <= 0 {
		maxConcurrent = 1
	}

	return &Semaphore{
		permits: make(chan struct{}, maxConcurrent),
		max:     maxConcurrent,
		current: 0,
	}
}

// Acquire acquires a permit from the semaphore
func (s *Semaphore) Acquire(ctx context.Context) error {
	select {
	case s.permits <- struct{}{}:
		s.mu.Lock()
		s.current++
		s.mu.Unlock()
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// Release releases a permit from the semaphore
func (s *Semaphore) Release() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.current > 0 {
		<-s.permits
		s.current--
	}
}

// CurrentCount returns the current number of permits in use
func (s *Semaphore) CurrentCount() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.current
}

// MaxCount returns the maximum number of permits
func (s *Semaphore) MaxCount() int {
	return s.max
}

// AvailableCount returns the number of available permits
func (s *Semaphore) AvailableCount() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.max - s.current
}

// WaitForAvailable waits until a permit is available
func (s *Semaphore) WaitForAvailable(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			if s.AvailableCount() > 0 {
				return nil
			}
			time.Sleep(10 * time.Millisecond)
		}
	}
}

// Stats returns statistics of the semaphore
func (s *Semaphore) Stats() map[string]int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return map[string]int{
		"max":       s.max,
		"current":   s.current,
		"available": s.max - s.current,
		"usage":     int(float64(s.current) / float64(s.max) * 100),
	}
}
