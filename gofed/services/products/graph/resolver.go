package graph

import "products/graph/model"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct{}

// Mock data para produtos com owners
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
}
