package graph

import (
	"context"
	"testing"
	"time"
)

func TestSemaphoreStats(t *testing.T) {
	resolver := NewResolver()

	// Test initial statistics
	stats, err := resolver.SemaphoreStats(context.Background())
	if err != nil {
		t.Errorf("Failed to get semaphore stats: %v", err)
	}

	if stats.Max != 3 {
		t.Errorf("Expected max 3, got %d", stats.Max)
	}
	if stats.Current != 0 {
		t.Errorf("Expected current 0, got %d", stats.Current)
	}
	if stats.Available != 3 {
		t.Errorf("Expected available 3, got %d", stats.Available)
	}
	if stats.Usage != 0 {
		t.Errorf("Expected usage 0, got %d", stats.Usage)
	}

	t.Log("SemaphoreStats test passed")
}

func TestProductsWithSemaphore(t *testing.T) {
	resolver := NewResolver()

	// Test concurrent resolution with semaphore
	ids := []string{"1", "2", "3", "4", "5"}

	start := time.Now()
	products, err := resolver.ProductsWithSemaphore(context.Background(), ids)
	duration := time.Since(start)

	if err != nil {
		t.Errorf("Failed to get products with semaphore: %v", err)
	}

	if len(products) != 5 {
		t.Errorf("Expected 5 products, got %d", len(products))
	}

	// Check if the duration is greater than expected (due to backpressure)
	expectedMinDuration := 200 * time.Millisecond * 2 // At least 2 products in parallel
	if duration < expectedMinDuration {
		t.Errorf("Expected duration >= %v (backpressure), got %v", expectedMinDuration, duration)
	}

	t.Logf("ProductsWithSemaphore test passed - Duration: %v", duration)
}

func TestProductsByCategory(t *testing.T) {
	resolver := NewResolver()

	products, err := resolver.ProductsByCategory(context.Background(), "Electronics")
	if err != nil {
		t.Errorf("Failed to get products by category: %v", err)
	}

	if len(products) == 0 {
		t.Error("Expected electronics products, got none")
	}

	for _, product := range products {
		if product.Category != "Electronics" {
			t.Errorf("Expected Electronics category, got %s", product.Category)
		}
	}

	t.Logf("ProductsByCategory test passed - Found %d electronics products", len(products))
}

func TestProductsByIds(t *testing.T) {
	resolver := NewResolver()

	ids := []string{"1", "2", "3"}

	start := time.Now()
	products, err := resolver.ProductsByIds(context.Background(), ids)
	duration := time.Since(start)

	if err != nil {
		t.Errorf("Failed to get products by IDs: %v", err)
	}

	if len(products) != 3 {
		t.Errorf("Expected 3 products, got %d", len(products))
	}

	// Check if the duration is less than expected (no backpressure)
	expectedMaxDuration := 200 * time.Millisecond * 3 // Maximum 3 products in parallel
	if duration > expectedMaxDuration {
		t.Errorf("Expected duration <= %v (no backpressure), got %v", expectedMaxDuration, duration)
	}

	t.Logf("ProductsByIds test passed - Duration: %v", duration)
}
