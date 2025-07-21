package graph

import (
	"sync"
	"time"
	"users/graph/model"
)

// UserCache simulates a cache with intentional race condition
type UserCache struct {
	// PROBLEM: Shared map without protection - RACE CONDITION!
	users map[string]*model.User

	// SOLUTION 1: Mutex for protection
	mu sync.RWMutex

	// SOLUTION 2: sync.Map for native thread-safety
	safeMap sync.Map

	// Configurations
	maxSize int
	ttl     time.Duration
}

// NewUserCache creates a new user cache
func NewUserCache(maxSize int, ttl time.Duration) *UserCache {
	return &UserCache{
		users:   make(map[string]*model.User),
		maxSize: maxSize,
		ttl:     ttl,
	}
}

// ===== IMPLEMENTATION WITH RACE CONDITION (PROBLEM) =====

// GetUserUnsafe - RACE CONDITION! Do not use in production
func (c *UserCache) GetUserUnsafe(id string) (*model.User, bool) {
	// PROBLEM: Concurrent access without protection
	user, exists := c.users[id]
	return user, exists
}

// SetUserUnsafe - RACE CONDITION! Do not use in production
func (c *UserCache) SetUserUnsafe(user *model.User) {
	// PROBLEM: Concurrent write without protection
	c.users[user.ID] = user
}

// GetUsersUnsafe - RACE CONDITION! Do not use in production
func (c *UserCache) GetUsersUnsafe() []*model.User {
	// PROBLEM: Concurrent read without protection
	users := make([]*model.User, 0, len(c.users))
	for _, user := range c.users {
		users = append(users, user)
	}
	return users
}

// ===== SOLUTION 1: MUTEX =====

// GetUserSafe - Thread-safe with mutex
func (c *UserCache) GetUserSafe(id string) (*model.User, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	user, exists := c.users[id]
	return user, exists
}

// SetUserSafe - Thread-safe with mutex
func (c *UserCache) SetUserSafe(user *model.User) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Check size limit
	if len(c.users) >= c.maxSize {
		// Remove oldest item (simple implementation)
		for key := range c.users {
			delete(c.users, key)
			break
		}
	}

	c.users[user.ID] = user
}

// GetUsersSafe - Thread-safe with mutex
func (c *UserCache) GetUsersSafe() []*model.User {
	c.mu.RLock()
	defer c.mu.RUnlock()

	users := make([]*model.User, 0, len(c.users))
	for _, user := range c.users {
		users = append(users, user)
	}
	return users
}

// ===== SOLUTION 2: SYNC.MAP =====

// GetUserSyncMap - Thread-safe with sync.Map
func (c *UserCache) GetUserSyncMap(id string) (*model.User, bool) {
	if value, ok := c.safeMap.Load(id); ok {
		if user, ok := value.(*model.User); ok {
			return user, true
		}
	}
	return nil, false
}

// SetUserSyncMap - Thread-safe with sync.Map
func (c *UserCache) SetUserSyncMap(user *model.User) {
	c.safeMap.Store(user.ID, user)
}

// GetUsersSyncMap - Thread-safe with sync.Map
func (c *UserCache) GetUsersSyncMap() []*model.User {
	users := make([]*model.User, 0)

	c.safeMap.Range(func(key, value interface{}) bool {
		if user, ok := value.(*model.User); ok {
			users = append(users, user)
		}
		return true
	})

	return users
}

// ===== UTILITY METHODS =====

// Clear clears the cache
func (c *UserCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.users = make(map[string]*model.User)
	c.safeMap = sync.Map{}
}

// Size returns the size of the cache
func (c *UserCache) Size() int {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return len(c.users)
}

// Stats returns statistics of the cache
func (c *UserCache) Stats() map[string]interface{} {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return map[string]interface{}{
		"size":     len(c.users),
		"max_size": c.maxSize,
		"ttl":      c.ttl.String(),
	}
}

// SimulateRaceCondition simulates a race condition for demonstration
func (c *UserCache) SimulateRaceCondition() {
	// Simulate multiple goroutines accessing the cache simultaneously
	var wg sync.WaitGroup

	// Goroutine 1: Read users
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 100; i++ {
			c.GetUsersUnsafe() // RACE CONDITION!
			time.Sleep(1 * time.Millisecond)
		}
	}()

	// Goroutine 2: Write users
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 100; i++ {
			user := &model.User{
				ID:    string(rune('A' + i)),
				Name:  "User" + string(rune('A'+i)),
				Email: "user" + string(rune('A'+i)) + "@example.com",
			}
			c.SetUserUnsafe(user) // RACE CONDITION!
			time.Sleep(1 * time.Millisecond)
		}
	}()

	// Goroutine 3: Read and write simultaneously
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 100; i++ {
			c.GetUserUnsafe("1") // RACE CONDITION!
			user := &model.User{
				ID:    "1",
				Name:  "UpdatedUser",
				Email: "updated@example.com",
			}
			c.SetUserUnsafe(user) // RACE CONDITION!
			time.Sleep(1 * time.Millisecond)
		}
	}()

	wg.Wait()
}

// SimulateSafeAccess simulates safe access for comparison
func (c *UserCache) SimulateSafeAccess() {
	var wg sync.WaitGroup

	// Goroutine 1: Read users safely
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 100; i++ {
			c.GetUsersSafe() // THREAD-SAFE
			time.Sleep(1 * time.Millisecond)
		}
	}()

	// Goroutine 2: Write users safely
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 100; i++ {
			user := &model.User{
				ID:    string(rune('A' + i)),
				Name:  "User" + string(rune('A'+i)),
				Email: "user" + string(rune('A'+i)) + "@example.com",
			}
			c.SetUserSafe(user) // THREAD-SAFE
			time.Sleep(1 * time.Millisecond)
		}
	}()

	// Goroutine 3: Read and write simultaneously safely
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 100; i++ {
			c.GetUserSafe("1") // THREAD-SAFE
			user := &model.User{
				ID:    "1",
				Name:  "UpdatedUser",
				Email: "updated@example.com",
			}
			c.SetUserSafe(user) // THREAD-SAFE
			time.Sleep(1 * time.Millisecond)
		}
	}()

	wg.Wait()
}
