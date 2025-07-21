package graph

import "users/graph/model"

type Resolver struct{}

// Mock data for users
var users = []*model.User{
	{
		ID:    "1",
		Name:  "Alice",
		Email: "alice@example.com",
	},
	{
		ID:    "2",
		Name:  "Bob",
		Email: "bob@example.com",
	},
	{
		ID:    "3",
		Name:  "Charlie",
		Email: "charlie@example.com",
	},
	{
		ID:    "4",
		Name:  "Diana",
		Email: "diana@example.com",
	},
	{
		ID:    "5",
		Name:  "Eve",
		Email: "eve@example.com",
	},
	{
		ID:    "6",
		Name:  "Frank",
		Email: "frank@example.com",
	},
	{
		ID:    "7",
		Name:  "Grace",
		Email: "grace@example.com",
	},
	{
		ID:    "8",
		Name:  "Henry",
		Email: "henry@example.com",
	},
}
