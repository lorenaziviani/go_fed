package graph

import (
	"time"
	"users/graph/model"
)

type Resolver struct {
	cache *UserCache
}

// NewResolver cria um novo resolver com cache configurado
func NewResolver() *Resolver {
	return &Resolver{
		cache: NewUserCache(100, 5*time.Minute), // Cache com 100 itens, TTL 5min
	}
}

// Cache retorna o cache do resolver
func (r *Resolver) Cache() *UserCache {
	return r.cache
}

// Mock data para usuários
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
