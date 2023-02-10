package database

import (
	"context"
)

/*

 */
type Collection[E interface{}] interface {
	InsertOne(ctx context.Context, document E) (string, error)
	FindOne(ctx context.Context, result *E, filter interface{}) error
	FindAll(ctx context.Context, filter interface{}) ([]E, error)
}