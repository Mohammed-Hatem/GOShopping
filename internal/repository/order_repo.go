package repository

import (
	"bookstore-project/internal/models"
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
)

type OrderRepo struct {
	db *sqlx.DB
}

func NewOrderRepo(db *sqlx.DB) *OrderRepo {
	return &OrderRepo{db: db}
}

// GetAllPublisherOrders retrieves all publisher orders
func (r *OrderRepo) GetAllPublisherOrders() ([]models.PublisherOrder, error) {
	var orders []models.PublisherOrder
	err := r.db.Select(&orders, "SELECT * FROM publisher_order ORDER BY order_date DESC")
	if err != nil {
		return nil, err
	}
	return orders, nil
}

// GetPublisherOrder retrieves a publisher order by ID
func (r *OrderRepo) GetPublisherOrder(id int) (*models.PublisherOrder, error) {
	var order models.PublisherOrder
	err := r.db.Get(&order, "SELECT * FROM publisher_order WHERE rep_order_id = $1", id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &order, nil
}

// PlacePublisherOrder creates a new publisher order
func (r *OrderRepo) PlacePublisherOrder(order models.PublisherOrder) error {
	query := `
		INSERT INTO publisher_order (order_date, status, quantity, isbn, admin_id)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING rep_order_id`

	err := r.db.QueryRow(query, time.Now(), "Pending", order.Quantity, order.ISBN, order.AdminID).Scan(&order.ID)
	return err
}

// ConfirmPublisherOrder confirms a publisher order (changes status to 'Confirmed')
func (r *OrderRepo) ConfirmPublisherOrder(id int) error {
	query := `UPDATE publisher_order SET status = 'Confirmed' WHERE rep_order_id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

// GetPendingPublisherOrders retrieves all pending publisher orders
func (r *OrderRepo) GetPendingPublisherOrders() ([]models.PublisherOrder, error) {
	var orders []models.PublisherOrder
	err := r.db.Select(&orders, "SELECT * FROM publisher_order WHERE status = 'Pending' ORDER BY order_date ASC")
	if err != nil {
		return nil, err
	}
	return orders, nil
}
