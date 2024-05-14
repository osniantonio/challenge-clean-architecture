package database

import (
	"database/sql"

	"github.com/osniantonio/challenge-clean-architecture/internal/entity"
)

type OrderRepository struct {
	db *sql.DB
}

func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{db}
}

func (r *OrderRepository) Save(order *entity.Order) error {
	stmt, err := r.db.Prepare("insert into `orders` (id, price, tax, final_price) values (?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()
	if _, err := stmt.Exec(order.ID, order.Price, order.Tax, order.FinalPrice); err != nil {
		return err
	}
	return nil
}

func (r *OrderRepository) FindAll() ([]entity.Order, error) {
	stmt, err := r.db.Prepare("select id, price, tax, final_price from `orders`")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var orders []entity.Order
	for rows.Next() {
		var order entity.Order
		if err := rows.Scan(&order.ID, &order.Price, &order.Tax, &order.FinalPrice); err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return orders, nil
}
