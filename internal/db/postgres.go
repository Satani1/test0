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

func (pr *PostgresRepository) InsertRow(ctx context.Context, orderParams models.CreateOrder) error {
	query := `INSERT INTO test (order_uid, data) values ($1,$2)`

	_, err := pr.db.Exec(query, orderParams.OrderUID, orderParams.Data)
	if err != nil {
		return err
	}
	return nil
}

func (pr *PostgresRepository) GetOrder(id string) (*models.Order, error) {
	query := `select data from test where order_uid = $1`
	var order models.Order
	err := pr.db.QueryRow(query, id).Scan(&order)
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (pr *PostgresRepository) GetAllOrders() ([]models.CreateOrder, error) {
	query := `select order_uid,data from test`
	var orders []models.CreateOrder

	rows, err := pr.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var order models.CreateOrder
		err := rows.Scan(&order.OrderUID, &order.Data)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}
	return orders, nil
}
