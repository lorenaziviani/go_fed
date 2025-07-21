package graph

import (
	"context"
	"testing"
	"time"
)

// Benchmark for sequential user resolution (simulated)
func BenchmarkSequentialUserResolution(b *testing.B) {
	ids := []string{"1", "2", "3", "4", "5"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, id := range ids {
			time.Sleep(100 * time.Millisecond)

			for _, user := range users {
				if user.ID == id {
					break
				}
			}
		}
	}
}

// Benchmark for concurrent user resolution
func BenchmarkConcurrentUserResolution(b *testing.B) {
	ids := []string{"1", "2", "3", "4", "5"}
	resolver := &Resolver{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ctx := context.Background()
		_, _ = resolver.UsersByIds(ctx, ids)
	}
}

// Benchmark for different array sizes
func BenchmarkConcurrentUserResolution_Small(b *testing.B) {
	ids := []string{"1", "2"}
	resolver := &Resolver{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ctx := context.Background()
		_, _ = resolver.UsersByIds(ctx, ids)
	}
}

func BenchmarkConcurrentUserResolution_Medium(b *testing.B) {
	ids := []string{"1", "2", "3", "4", "5", "6", "7", "8"}
	resolver := &Resolver{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ctx := context.Background()
		_, _ = resolver.UsersByIds(ctx, ids)
	}
}

func BenchmarkConcurrentUserResolution_Large(b *testing.B) {
	ids := make([]string, 20)
	for i := 0; i < 20; i++ {
		ids[i] = string(rune('1' + i))
	}
	resolver := &Resolver{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ctx := context.Background()
		_, _ = resolver.UsersByIds(ctx, ids)
	}
}

// Test for race condition
func TestConcurrentUserResolution_RaceCondition(t *testing.T) {
	ids := []string{"1", "2", "3", "4", "5"}
	resolver := &Resolver{}

	// Execute multiple goroutines simultaneously
	for i := 0; i < 10; i++ {
		go func() {
			ctx := context.Background()
			_, _ = resolver.UsersByIds(ctx, ids)
		}()
	}

	// Wait a bit for all goroutines to finish
	time.Sleep(1 * time.Second)
}

// Test for timeout
func TestConcurrentUserResolution_Timeout(t *testing.T) {
	ids := []string{"1", "2", "3", "4", "5"}
	resolver := &Resolver{}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer cancel()

	_, err := resolver.UsersByIds(ctx, ids)

	if err == nil {
		t.Error("Expected timeout error, got nil")
	}
}

// Test for cancellation
func TestConcurrentUserResolution_Cancellation(t *testing.T) {
	ids := []string{"1", "2", "3", "4", "5"}
	resolver := &Resolver{}

	ctx, cancel := context.WithCancel(context.Background())

	cancel()

	_, err := resolver.UsersByIds(ctx, ids)

	if err == nil {
		t.Error("Expected cancellation error, got nil")
	}
}

// Test for valid data
func TestConcurrentUserResolution_ValidData(t *testing.T) {
	ids := []string{"1", "2", "3", "4", "5"}
	resolver := &Resolver{}

	ctx := context.Background()
	result, err := resolver.UsersByIds(ctx, ids)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(result) != 5 {
		t.Errorf("Expected 5 users, got %d", len(result))
	}

	expectedIDs := map[string]bool{"1": true, "2": true, "3": true, "4": true, "5": true}
	for _, user := range result {
		if !expectedIDs[user.ID] {
			t.Errorf("Unexpected user ID: %s", user.ID)
		}
	}
}

// Test for invalid data
func TestConcurrentUserResolution_InvalidData(t *testing.T) {
	ids := []string{"999", "998", "997"}
	resolver := &Resolver{}

	ctx := context.Background()
	result, err := resolver.UsersByIds(ctx, ids)

	if err != nil {
		t.Errorf("Expected no error for invalid IDs, got %v", err)
	}

	if len(result) != 0 {
		t.Errorf("Expected 0 users for invalid IDs, got %d", len(result))
	}
}

// Test for stress with multiple concurrent executions
func TestConcurrentUserResolution_Stress(t *testing.T) {
	ids := []string{"1", "2", "3", "4", "5"}
	resolver := &Resolver{}

	done := make(chan bool, 50)

	for i := 0; i < 50; i++ {
		go func() {
			ctx := context.Background()
			_, _ = resolver.UsersByIds(ctx, ids)
			done <- true
		}()
	}

	for i := 0; i < 50; i++ {
		<-done
	}
}
