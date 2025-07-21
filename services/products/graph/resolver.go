package graph

import "products/graph/model"

type Resolver struct {
	semaphore *Semaphore
}

// NewResolver creates a new resolver with the semaphore configured
func NewResolver() *Resolver {
	return &Resolver{
		semaphore: NewSemaphore(3), // Maximum 3 concurrent resolutions
	}
}

// Semaphore returns the semaphore of the resolver
func (r *Resolver) Semaphore() *Semaphore {
	return r.semaphore
}

// Mock data for products
var products = []*model.Product{
	{
		ID:          "1",
		Name:        "iPhone 15 Pro",
		Description: "Smartphone Apple com chip A17 Pro",
		Price:       999.99,
		Category:    "Electronics",
		Owner: &model.User{
			ID:    "1",
			Name:  "Alice",
			Email: "alice@example.com",
		},
	},
	{
		ID:          "2",
		Name:        "MacBook Air M2",
		Description: "Notebook Apple com chip M2",
		Price:       1199.99,
		Category:    "Electronics",
		Owner: &model.User{
			ID:    "2",
			Name:  "Bob",
			Email: "bob@example.com",
		},
	},
	{
		ID:          "3",
		Name:        "Nike Air Max",
		Description: "Tênis esportivo Nike",
		Price:       129.99,
		Category:    "Sports",
		Owner: &model.User{
			ID:    "1",
			Name:  "Alice",
			Email: "alice@example.com",
		},
	},
	{
		ID:          "4",
		Name:        "Coffee Maker",
		Description: "Máquina de café automática",
		Price:       89.99,
		Category:    "Home",
		Owner: &model.User{
			ID:    "3",
			Name:  "Charlie",
			Email: "charlie@example.com",
		},
	},
	{
		ID:          "5",
		Name:        "Gaming Mouse",
		Description: "Mouse gamer com RGB",
		Price:       79.99,
		Category:    "Electronics",
		Owner: &model.User{
			ID:    "4",
			Name:  "Diana",
			Email: "diana@example.com",
		},
	},
	{
		ID:          "6",
		Name:        "Yoga Mat",
		Description: "Tapete de yoga premium",
		Price:       45.99,
		Category:    "Sports",
		Owner: &model.User{
			ID:    "5",
			Name:  "Eve",
			Email: "eve@example.com",
		},
	},
	{
		ID:          "7",
		Name:        "Bluetooth Speaker",
		Description: "Alto-falante Bluetooth portátil",
		Price:       129.99,
		Category:    "Electronics",
		Owner: &model.User{
			ID:    "6",
			Name:  "Frank",
			Email: "frank@example.com",
		},
	},
	{
		ID:          "8",
		Name:        "Smart Watch",
		Description: "Relógio inteligente com monitor cardíaco",
		Price:       299.99,
		Category:    "Electronics",
		Owner: &model.User{
			ID:    "7",
			Name:  "Grace",
			Email: "grace@example.com",
		},
	},
}
