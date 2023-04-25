package db

import (
	"context"
	"test0/internal/models"
)

type Repository interface {
	Close()
	InsertRow(ctx context.Context, order models.Order) error
	ListTable(ctx context.Context) ([]models.Order, error)
}

var impl Repository

func SetRepository(repository Repository) {
	impl = repository
}

func Close() {
	impl.Close()
}

func InsertRow(ctx context.Context, order models.Order) error {
	return impl.InsertRow(ctx, order)
}

func ListTabe(ctx context.Context) ([]models.Order, error) {
	return impl.ListTable(ctx)
}
