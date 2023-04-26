package db

import (
	"context"
	"database/sql"
	_ "github.com/lib/pq"
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
	if err := db.Ping(); err != nil {
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

func (pr *PostgresRepository) GetOrder(id string) (*models.Order, error) {
	query := `select data from "modelDB".test where order_uid = $1`
	var order models.Order
	err := pr.db.QueryRow(query, id).Scan(&order)
	if err != nil {
		return nil, err
	}
	return &order, nil
}
