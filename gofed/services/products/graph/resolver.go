package graph

import "products/graph/model"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct{}

var products = []*model.Product{
	{
		ID:          "1",
		Name:        "iPhone 15 Pro",
		Description: "Latest iPhone with advanced camera system",
		Price:       999.99,
		Category:    "Electronics",
	},
	{
		ID:          "2",
		Name:        "MacBook Air M2",
		Description: "Ultra-thin laptop with M2 chip",
		Price:       1199.99,
		Category:    "Electronics",
	},
	{
		ID:          "3",
		Name:        "Nike Air Max",
		Description: "Comfortable running shoes",
		Price:       129.99,
		Category:    "Sports",
	},
	{
		ID:          "4",
		Name:        "Coffee Maker",
		Description: "Automatic coffee machine",
		Price:       89.99,
		Category:    "Home & Kitchen",
	},
}
