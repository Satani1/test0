package db

import (
	"context"
	"database/sql"
	"test0/internal/models"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgres(url string) (*PostgresRepository, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}
	return &PostgresRepository{
		db,
	}, nil
}

func (pr *PostgresRepository) Close() {
	pr.db.Close()
}

func (pr *PostgresRepository) InsertRow(ctx context.Context, order models.Order) error {
	return nil
}

func (pr *PostgresRepository) ListTable(ctx context.Context) ([]models.Order, error) {
	return nil, nil
}
