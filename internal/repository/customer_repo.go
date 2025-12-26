package repository

import (
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"bookstore-project/internal/models"
)

type CustomerRepo struct {
	db *sqlx.DB
}

func NewCustomerRepo(db *sqlx.DB) *CustomerRepo {
	return &CustomerRepo{db: db}
}

func (r *CustomerRepo) CreateCustomer(username string, password string, firstName string, lastName string, email string, phone string, address string) error {
	_, err := r.db.Exec("INSERT INTO customer (username,password,first_name,last_name,email,phone_number,shipping_address) VALUES ($1, $2, $3, $4, $5, $6, $7)", username, password, firstName, lastName, email, phone, address)

	if err != nil { //postgres unique violation error code
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			switch pqErr.Constraint {
			case "customer_pkey":
				return errors.New("username already exists")
			case "customer_email_key":
				return errors.New("email already exists")
			}
		}
		return err
	}
	return nil
}

func (r *CustomerRepo) GetCustomerByUsername(username string) (*models.Customer, error) {
	var customer models.Customer

	err := r.db.Get(&customer, "SELECT * FROM customer WHERE username = $1", username)

	if err != nil {
		return nil, err
	}
	return &customer, nil
}

func (r *CustomerRepo) GetCustomerByEmail(email string) (*models.Customer, error) {

	var customer models.Customer
	err := r.db.Get(&customer, "SELECT * FROM customer WHERE email = $1", email)

	if err != nil {
		return nil, err
	}
	return &customer, nil
}

func (r *CustomerRepo) GetOrdersByCustomer(username string) ([]models.SalesOrder, error) {
	var orders []models.SalesOrder

	err := r.db.Select(&orders, "SELECT * FROM sales_order WHERE username = $1", username)
	if err != nil {
		return nil, err
	}
	return orders, nil
}

func (r *CustomerRepo) GetOrderItemsByOrderID(orderID int) ([]models.OrderItem, error) {
	var items []models.OrderItem

	err := r.db.Select(&items, "SELECT * FROM order_item WHERE order_id = $1", orderID)
	if err != nil {
		return nil, err
	}
	return items, nil
}
