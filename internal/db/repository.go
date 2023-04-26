package db

import (
	"context"
	"test0/internal/models"
)

type Repository interface {
	Close()
	InsertRow(ctx context.Context, order models.Order) error
	GetOrder(id string) (*models.Order, error)
}

var impl Repository

func SetRepository(repository *PostgresRepository) {
	impl = repository
}

func Close() {
	impl.Close()
}

func InsertRow(ctx context.Context, order models.Order) error {
	return impl.InsertRow(ctx, order)
}

func GetOrder(id string) (*models.Order, error) {
	return impl.GetOrder(id)
}
