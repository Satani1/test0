package db

import (
	"context"
	"test0/internal/models"
)

type Repository interface {
	Close()
	InsertRow(ctx context.Context, order models.Order) error
	ListTable() ([]testData, error)
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

func ListTable() ([]testData, error) {
	return impl.ListTable()
}
