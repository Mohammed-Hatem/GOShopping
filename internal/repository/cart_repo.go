package repository

import (
	"database/sql"
	"bookstore-project/internal/models"
	"github.com/jmoiron/sqlx"
)

type CartRepo struct {
	db *sqlx.DB
}

func NewCartRepo(db *sqlx.DB) *CartRepo {
	return &CartRepo{db: db}
}

// CartItemWithBook represents a cart item with book details
type CartItemWithBook struct {
	ISBN         string  `db:"isbn" json:"isbn"`
	Title        string  `db:"title" json:"title"`
	AuthorName   string  `db:"author_name" json:"author_name"`
	SellingPrice float64 `db:"selling_price" json:"selling_price"`
	Quantity     int     `db:"quantity" json:"quantity"`
	StockQuantity int    `db:"stock_quantity" json:"stock_quantity"`
}

// GetCartByUsername gets or creates a cart for a user
func (r *CartRepo) GetCartByUsername(username string) (*models.ShoppingCart, error) {
	var cart models.ShoppingCart
	
	// Try to get existing cart
	err := r.db.Get(&cart, "SELECT * FROM shopping_cart WHERE username = $1", username)
	if err == sql.ErrNoRows {
		// Cart doesn't exist, create one
		query := `INSERT INTO shopping_cart (username) VALUES ($1) RETURNING cart_id`
		err = r.db.QueryRow(query, username).Scan(&cart.ID)
		if err != nil {
			return nil, err
		}
		cart.Username = username
		return &cart, nil
	}
	if err != nil {
		return nil, err
	}
	
	return &cart, nil
}

// AddItemToCart adds or updates an item in the cart
func (r *CartRepo) AddItemToCart(cartID int, isbn string, quantity int) error {
	query := `
		INSERT INTO cart_item (cart_id, isbn, quantity)
		VALUES ($1, $2, $3)
		ON CONFLICT (cart_id, isbn)
		DO UPDATE SET quantity = cart_item.quantity + $3
	`
	_, err := r.db.Exec(query, cartID, isbn, quantity)
	return err
}

// GetCartItems retrieves all cart items with book details
func (r *CartRepo) GetCartItems(cartID int) ([]CartItemWithBook, error) {
	var items []CartItemWithBook
	query := `
		SELECT ci.isbn, b.title, b.author_name, b.selling_price, ci.quantity, b.stock_quantity
		FROM cart_item ci
		JOIN book b ON ci.isbn = b.isbn
		WHERE ci.cart_id = $1
	`
	err := r.db.Select(&items, query, cartID)
	if err != nil {
		return nil, err
	}
	return items, nil
}

// RemoveItemFromCart removes an item from the cart
func (r *CartRepo) RemoveItemFromCart(cartID int, isbn string) error {
	query := `DELETE FROM cart_item WHERE cart_id = $1 AND isbn = $2`
	_, err := r.db.Exec(query, cartID, isbn)
	return err
}

// ClearCart removes all items from the cart
func (r *CartRepo) ClearCart(cartID int) error {
	query := `DELETE FROM cart_item WHERE cart_id = $1`
	_, err := r.db.Exec(query, cartID)
	return err
}

// GetCartTotal calculates the total price of all items in the cart
func (r *CartRepo) GetCartTotal(cartID int) (float64, error) {
	var total float64
	query := `
		SELECT COALESCE(SUM(b.selling_price * ci.quantity), 0)
		FROM cart_item ci
		JOIN book b ON ci.isbn = b.isbn
		WHERE ci.cart_id = $1
	`
	err := r.db.Get(&total, query, cartID)
	if err != nil {
		return 0, err
	}
	return total, nil
}

// UpdateCartItemQuantity updates the quantity of a cart item
func (r *CartRepo) UpdateCartItemQuantity(cartID int, isbn string, quantity int) error {
	if quantity <= 0 {
		// If quantity is 0 or less, remove the item
		return r.RemoveItemFromCart(cartID, isbn)
	}
	query := `UPDATE cart_item SET quantity = $3 WHERE cart_id = $1 AND isbn = $2`
	_, err := r.db.Exec(query, cartID, isbn, quantity)
	return err
}

