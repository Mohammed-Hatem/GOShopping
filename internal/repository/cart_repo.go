package repository

import (
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
)

type CartRepo struct {
	db *sqlx.DB
}

func NewCartRepo(db *sqlx.DB) *CartRepo {
	return &CartRepo{db: db}
}

type CartLineItem struct {
	ISBN      string  `db:"isbn" json:"isbn"`
	Title     string  `db:"title" json:"title"`
	Quantity  int     `db:"quantity" json:"quantity"`
	UnitPrice float64 `db:"unit_price" json:"unit_price"`
	LineTotal float64 `db:"line_total" json:"line_total"`
}

type CartSummary struct {
	Username string         `json:"username"`
	Items    []CartLineItem `json:"items"`
	Total    float64        `json:"total"`
}

func (r *CartRepo) EnsureCart(username string) (int, error) {
	if username == "" {
		return 0, errors.New("username is required")
	}

	var cartID int
	err := r.db.Get(&cartID, `
		INSERT INTO shopping_cart (username)
		VALUES ($1)
		ON CONFLICT (username)
		DO UPDATE SET username = EXCLUDED.username
		RETURNING cart_id
	`, username)
	if err != nil {
		return 0, err
	}
	return cartID, nil
}

func (r *CartRepo) AddToCart(username, isbn string, quantity int) error {
	if username == "" {
		return errors.New("username is required")
	}
	if isbn == "" {
		return errors.New("isbn is required")
	}
	if quantity <= 0 {
		return errors.New("quantity must be positive")
	}

	var exists int
	err := r.db.Get(&exists, "SELECT 1 FROM book WHERE isbn = $1", isbn)
	if err == sql.ErrNoRows {
		return errors.New("book not found")
	}
	if err != nil {
		return err
	}

	cartID, err := r.EnsureCart(username)
	if err != nil {
		return err
	}

	_, err = r.db.Exec(`
		INSERT INTO cart_item (cart_id, isbn, quantity)
		VALUES ($1, $2, $3)
		ON CONFLICT (cart_id, isbn)
		DO UPDATE SET quantity = cart_item.quantity + EXCLUDED.quantity
	`, cartID, isbn, quantity)
	return err
}

func (r *CartRepo) RemoveFromCart(username, isbn string) error {
	if username == "" {
		return errors.New("username is required")
	}
	if isbn == "" {
		return errors.New("isbn is required")
	}

	var cartID int
	err := r.db.Get(&cartID, "SELECT cart_id FROM shopping_cart WHERE username = $1", username)
	if err == sql.ErrNoRows {
		return nil
	}
	if err != nil {
		return err
	}

	_, err = r.db.Exec("DELETE FROM cart_item WHERE cart_id = $1 AND isbn = $2", cartID, isbn)
	return err
}

func (r *CartRepo) ClearCart(username string) error {
	if username == "" {
		return errors.New("username is required")
	}

	var cartID int
	err := r.db.Get(&cartID, "SELECT cart_id FROM shopping_cart WHERE username = $1", username)
	if err == sql.ErrNoRows {
		return nil
	}
	if err != nil {
		return err
	}

	_, err = r.db.Exec("DELETE FROM cart_item WHERE cart_id = $1", cartID)
	return err
}

func (r *CartRepo) GetCartSummary(username string) (*CartSummary, error) {
	if username == "" {
		return nil, errors.New("username is required")
	}

	var cartID int
	err := r.db.Get(&cartID, "SELECT cart_id FROM shopping_cart WHERE username = $1", username)
	if err == sql.ErrNoRows {
		return &CartSummary{Username: username, Items: []CartLineItem{}, Total: 0}, nil
	}
	if err != nil {
		return nil, err
	}

	items := make([]CartLineItem, 0)
	err = r.db.Select(&items, `
		SELECT
			ci.isbn,
			b.title,
			ci.quantity,
			b.selling_price AS unit_price,
			(ci.quantity * b.selling_price) AS line_total
		FROM cart_item ci
		JOIN book b ON b.isbn = ci.isbn
		WHERE ci.cart_id = $1
		ORDER BY b.title
	`, cartID)
	if err != nil {
		return nil, err
	}

	total := 0.0
	for _, it := range items {
		total += it.LineTotal
	}

	return &CartSummary{Username: username, Items: items, Total: total}, nil
}

