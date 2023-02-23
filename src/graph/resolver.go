package graph

import (
	"sync"

	"github.com/takeru-a/self-introduction-app-backend/graph/model"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Keys struct{
    token string
    userId string
}

type Resolver struct {
    subscribers map[Keys]chan<- *model.Room
    mutex       sync.Mutex
}

func NewResolver() *Resolver {
    return &Resolver{
        subscribers: map[Keys]chan<- *model.Room{},
        mutex:       sync.Mutex{},
    }
}