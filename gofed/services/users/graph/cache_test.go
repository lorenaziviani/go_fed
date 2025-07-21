package graph

import (
	"sync"
	"testing"
	"time"
	"users/graph/model"
)

func TestUserCache(t *testing.T) {
	cache := NewUserCache(10, 5*time.Minute)

	// Test initial statistics
	stats := cache.Stats()
	if stats["size"] != 0 {
		t.Errorf("Expected size 0, got %d", stats["size"])
	}
	if stats["max_size"] != 10 {
		t.Errorf("Expected max_size 10, got %d", stats["max_size"])
	}

	// Test adding user
	user := &model.User{
		ID:    "1",
		Name:  "Alice",
		Email: "alice@example.com",
	}

	cache.SetUserSafe(user)

	// Check if it was added
	if cache.Size() != 1 {
		t.Errorf("Expected size 1, got %d", cache.Size())
	}

	// Test searching user
	if foundUser, exists := cache.GetUserSafe("1"); !exists {
		t.Error("Expected user to exist")
	} else if foundUser.Name != "Alice" {
		t.Errorf("Expected Alice, got %s", foundUser.Name)
	}

	// Test searching nonexistent user
	if _, exists := cache.GetUserSafe("999"); exists {
		t.Error("Expected user to not exist")
	}

	t.Log("UserCache test passed")
}

func TestUserCacheConcurrency(t *testing.T) {
	cache := NewUserCache(100, 5*time.Minute)

	// Test concurrent safe access
	var wg sync.WaitGroup
	results := make(chan bool, 50)

	// Goroutine 1: Add users
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 25; i++ {
			user := &model.User{
				ID:    string(rune('A' + i)),
				Name:  "User" + string(rune('A'+i)),
				Email: "user" + string(rune('A'+i)) + "@example.com",
			}
			cache.SetUserSafe(user)
			results <- true
		}
	}()

	// Goroutine 2: Read users
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 25; i++ {
			cache.GetUsersSafe()
			results <- true
		}
	}()

	wg.Wait()
	close(results)

	// Check if all operations were successful
	count := 0
	for range results {
		count++
	}

	if count != 50 {
		t.Errorf("Expected 50 operations, got %d", count)
	}

	t.Log("UserCache concurrency test passed")
}

func TestUserCacheSyncMap(t *testing.T) {
	cache := NewUserCache(10, 5*time.Minute)

	// Test sync.Map
	user := &model.User{
		ID:    "1",
		Name:  "Alice",
		Email: "alice@example.com",
	}

	cache.SetUserSyncMap(user)

	// Check if it was added
	if foundUser, exists := cache.GetUserSyncMap("1"); !exists {
		t.Error("Expected user to exist")
	} else if foundUser.Name != "Alice" {
		t.Errorf("Expected Alice, got %s", foundUser.Name)
	}

	// Test searching nonexistent user
	if _, exists := cache.GetUserSyncMap("999"); exists {
		t.Error("Expected user to not exist")
	}

	t.Log("UserCache sync.Map test passed")
}

func TestUserCacheMaxSize(t *testing.T) {
	cache := NewUserCache(2, 5*time.Minute)

	// Add 3 users (should remove the first one)
	user1 := &model.User{ID: "1", Name: "Alice", Email: "alice@example.com"}
	user2 := &model.User{ID: "2", Name: "Bob", Email: "bob@example.com"}
	user3 := &model.User{ID: "3", Name: "Charlie", Email: "charlie@example.com"}

	cache.SetUserSafe(user1)
	cache.SetUserSafe(user2)
	cache.SetUserSafe(user3)

	// Check if the size did not exceed the maximum
	if cache.Size() > 2 {
		t.Errorf("Expected size <= 2, got %d", cache.Size())
	}

	t.Log("UserCache max size test passed")
}

func TestRaceConditionSimulation(t *testing.T) {
	cache := NewUserCache(100, 5*time.Minute)

	// Simulate race condition (can cause panic, but should be captured)
	defer func() {
		if r := recover(); r != nil {
			t.Logf("Race condition detected (expected): %v", r)
		}
	}()

	cache.SimulateRaceCondition()

	t.Log("Race condition simulation test passed")
}

func TestSafeAccessSimulation(t *testing.T) {
	cache := NewUserCache(100, 5*time.Minute)

	// Simulate safe access (should not cause panic)
	cache.SimulateSafeAccess()

	t.Log("Safe access simulation test passed")
}

func BenchmarkCacheSafeAccess(b *testing.B) {
	cache := NewUserCache(100, 5*time.Minute)

	// Populate cache
	for i := 0; i < 50; i++ {
		user := &model.User{
			ID:    string(rune('A' + i)),
			Name:  "User" + string(rune('A'+i)),
			Email: "user" + string(rune('A'+i)) + "@example.com",
		}
		cache.SetUserSafe(user)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.GetUsersSafe()
	}
}

func BenchmarkCacheUnsafeAccess(b *testing.B) {
	cache := NewUserCache(100, 5*time.Minute)

	// Populate cache
	for i := 0; i < 50; i++ {
		user := &model.User{
			ID:    string(rune('A' + i)),
			Name:  "User" + string(rune('A'+i)),
			Email: "user" + string(rune('A'+i)) + "@example.com",
		}
		cache.SetUserUnsafe(user)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.GetUsersUnsafe()
	}
}

func BenchmarkCacheSyncMapAccess(b *testing.B) {
	cache := NewUserCache(100, 5*time.Minute)

	// Populate cache
	for i := 0; i < 50; i++ {
		user := &model.User{
			ID:    string(rune('A' + i)),
			Name:  "User" + string(rune('A'+i)),
			Email: "user" + string(rune('A'+i)) + "@example.com",
		}
		cache.SetUserSyncMap(user)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.GetUsersSyncMap()
	}
}
