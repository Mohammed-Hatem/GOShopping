package repository

import (
	"bookstore-project/internal/models"
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
)

type SalesRepo struct {
	db *sqlx.DB
}

func NewSalesRepo(db *sqlx.DB) *SalesRepo {
	return &SalesRepo{db: db}
}

// CreateSalesOrder creates a new sales order and returns the order ID
func (r *SalesRepo) CreateSalesOrder(username string, totalAmount float64, creditCardNo string, expiryDate time.Time) (int, error) {
	var orderID int
	query := `
		INSERT INTO sales_order (order_date, total_amount, credit_card_no, expiry_date, username)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING order_id
	`
	err := r.db.QueryRow(query, time.Now(), totalAmount, creditCardNo, expiryDate, username).Scan(&orderID)
	if err != nil {
		return 0, err
	}
	return orderID, nil
}

// CreateOrderItem creates an order item for a sales order
func (r *SalesRepo) CreateOrderItem(orderID int, isbn string, quantity int, unitPrice float64) error {
	query := `
		INSERT INTO order_item (order_id, isbn, quantity, unit_price)
		VALUES ($1, $2, $3, $4)
	`
	_, err := r.db.Exec(query, orderID, isbn, quantity, unitPrice)
	return err
}

// GetSalesOrder retrieves a sales order by ID
func (r *SalesRepo) GetSalesOrder(orderID int) (*models.SalesOrder, error) {
	var order models.SalesOrder
	err := r.db.Get(&order, "SELECT * FROM sales_order WHERE order_id = $1", orderID)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &order, nil
}

