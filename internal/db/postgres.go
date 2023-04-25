package db

import (
	"context"
	"database/sql"
	"fmt"
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

type testData struct {
	id    int
	uid   string
	order models.Order
}

func (pr *PostgresRepository) ListTable() ([]testData, error) {
	query := `select * from "modelDB".test2`

	rows, err := pr.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	data := []testData{}

	for rows.Next() {
		d := testData{}
		err := rows.Scan(&d.id, &d.uid, &d.order)
		if err != nil {
			fmt.Println("error", err)
			continue
		}
		data = append(data, d)
	}

	for _, d := range data {
		fmt.Println(d.id, d.uid, d.order)
	}

	return nil, nil
}
