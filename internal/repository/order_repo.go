package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

type OrderRepo struct {
	db *sqlx.DB
}

func NewOrderRepo(db *sqlx.DB) *OrderRepo {
	return &OrderRepo{db: db}
}

type CheckoutResult struct {
	OrderID int     `json:"order_id"`
	Total   float64 `json:"total_amount"`
}

func (r *OrderRepo) CheckoutCart(username, creditCardNo string, expiryDate time.Time) (*CheckoutResult, error) {
	if username == "" {
		return nil, errors.New("username is required")
	}
	if creditCardNo == "" {
		return nil, errors.New("credit card number is required")
	}
	if expiryDate.IsZero() {
		return nil, errors.New("expiry date is required")
	}

	tx, err := r.db.BeginTxx(context.Background(), nil)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	var cartID int
	err = tx.Get(&cartID, "SELECT cart_id FROM shopping_cart WHERE username = $1", username)
	if err == sql.ErrNoRows {
		return nil, errors.New("cart is empty")
	}
	if err != nil {
		return nil, err
	}

	items := make([]struct {
		ISBN     string `db:"isbn"`
		Quantity int    `db:"quantity"`
	}, 0)
	err = tx.Select(&items, "SELECT isbn, quantity FROM cart_item WHERE cart_id = $1", cartID)
	if err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return nil, errors.New("cart is empty")
	}

	total := 0.0
	prices := make(map[string]float64, len(items))

	for _, it := range items {
		if it.Quantity <= 0 {
			return nil, errors.New("invalid cart quantity")
		}

		var stock int
		var unitPrice float64
		err = tx.QueryRowx(
			"SELECT stock_quantity, selling_price FROM book WHERE isbn = $1 FOR UPDATE",
			it.ISBN,
		).Scan(&stock, &unitPrice)
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("book not found: %s", it.ISBN)
		}
		if err != nil {
			return nil, err
		}

		if stock < it.Quantity {
			return nil, fmt.Errorf("insufficient stock for isbn %s", it.ISBN)
		}

		prices[it.ISBN] = unitPrice
		total += float64(it.Quantity) * unitPrice
	}

	var orderID int
	err = tx.QueryRowx(
		`INSERT INTO sales_order (order_date, total_amount, credit_card_no, expiry_date, username)
		 VALUES (CURRENT_DATE, $1, $2, $3, $4)
		 RETURNING order_id`,
		total, creditCardNo, expiryDate, username,
	).Scan(&orderID)
	if err != nil {
		return nil, err
	}

	for _, it := range items {
		unitPrice := prices[it.ISBN]
		_, err = tx.Exec(
			"INSERT INTO order_item (order_id, isbn, quantity, unit_price) VALUES ($1, $2, $3, $4)",
			orderID, it.ISBN, it.Quantity, unitPrice,
		)
		if err != nil {
			return nil, err
		}

		res, err := tx.Exec(
			"UPDATE book SET stock_quantity = stock_quantity - $1 WHERE isbn = $2 AND stock_quantity >= $1",
			it.Quantity, it.ISBN,
		)
		if err != nil {
			return nil, err
		}
		affected, err := res.RowsAffected()
		if err != nil {
			return nil, err
		}
		if affected == 0 {
			return nil, fmt.Errorf("insufficient stock for isbn %s", it.ISBN)
		}
	}

	_, err = tx.Exec("DELETE FROM cart_item WHERE cart_id = $1", cartID)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &CheckoutResult{OrderID: orderID, Total: total}, nil
}

