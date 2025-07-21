package graph

import (
	"context"
	"testing"
	"time"
)

func TestSemaphore(t *testing.T) {
	sem := NewSemaphore(3)

	// Check initial statistics
	stats := sem.Stats()
	if stats["max"] != 3 {
		t.Errorf("Expected max 3, got %d", stats["max"])
	}
	if stats["current"] != 0 {
		t.Errorf("Expected current 0, got %d", stats["current"])
	}
	if stats["available"] != 3 {
		t.Errorf("Expected available 3, got %d", stats["available"])
	}

	ctx := context.Background()
	err := sem.Acquire(ctx)
	if err != nil {
		t.Errorf("Failed to acquire semaphore: %v", err)
	}

	// Check statistics after acquisition
	stats = sem.Stats()
	if stats["current"] != 1 {
		t.Errorf("Expected current 1, got %d", stats["current"])
	}
	if stats["available"] != 2 {
		t.Errorf("Expected available 2, got %d", stats["available"])
	}

	// Release permit
	sem.Release()

	// Check statistics after release
	stats = sem.Stats()
	if stats["current"] != 0 {
		t.Errorf("Expected current 0, got %d", stats["current"])
	}
	if stats["available"] != 3 {
		t.Errorf("Expected available 3, got %d", stats["available"])
	}

	t.Log("Semaphore test passed")
}

func TestSemaphoreConcurrency(t *testing.T) {
	sem := NewSemaphore(2)

	results := make(chan int, 5)

	// Function to test concurrency
	testFunc := func(id int) {
		ctx := context.Background()

		err := sem.Acquire(ctx)
		if err != nil {
			t.Errorf("Failed to acquire semaphore: %v", err)
			return
		}
		defer sem.Release()

		time.Sleep(100 * time.Millisecond)

		// Check that there are no more than 2 permits in use
		stats := sem.Stats()
		if stats["current"] > 2 {
			t.Errorf("Expected max 2 concurrent, got %d", stats["current"])
		}

		results <- id
	}

	for i := 0; i < 5; i++ {
		go testFunc(i)
	}

	for i := 0; i < 5; i++ {
		<-results
	}

	t.Log("Semaphore concurrency test passed")
}

func TestSemaphoreTimeout(t *testing.T) {
	sem := NewSemaphore(1)

	ctx := context.Background()
	err := sem.Acquire(ctx)
	if err != nil {
		t.Errorf("Failed to acquire semaphore: %v", err)
	}

	// Try to acquire another permit with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	err = sem.Acquire(ctx)
	if err == nil {
		t.Error("Expected timeout error, got nil")
	}

	sem.Release()

	t.Log("Semaphore timeout test passed")
}
