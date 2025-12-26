package repository

import (
	"time"

	"github.com/jmoiron/sqlx"
)

type ReportRepo struct {
	db *sqlx.DB
}

func NewReportRepo(db *sqlx.DB) *ReportRepo {
	return &ReportRepo{db: db}
}

// MonthlySales represents total sales for a month
type MonthlySales struct {
	TotalSales float64 `db:"total_sales" json:"total_sales"`
	Month      int     `db:"month" json:"month"`
	Year       int     `db:"year" json:"year"`
}

// GetMonthlySales gets total sales for the previous month
func (r *ReportRepo) GetMonthlySales() (*MonthlySales, error) {
	var result MonthlySales
	query := `
		SELECT 
			COALESCE(SUM(total_amount), 0) as total_sales,
			EXTRACT(MONTH FROM CURRENT_DATE - INTERVAL '1 month')::INTEGER as month,
			EXTRACT(YEAR FROM CURRENT_DATE - INTERVAL '1 month')::INTEGER as year
		FROM sales_order
		WHERE EXTRACT(YEAR FROM order_date) = EXTRACT(YEAR FROM CURRENT_DATE - INTERVAL '1 month')
		  AND EXTRACT(MONTH FROM order_date) = EXTRACT(MONTH FROM CURRENT_DATE - INTERVAL '1 month')
	`
	err := r.db.Get(&result, query)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// DailySales represents total sales for a specific date
type DailySales struct {
	TotalSales float64 `db:"total_sales" json:"total_sales"`
	Date       string  `db:"date" json:"date"`
}

// GetDailySales gets total sales for a specific date
func (r *ReportRepo) GetDailySales(date time.Time) (*DailySales, error) {
	var result DailySales
	query := `
		SELECT 
			COALESCE(SUM(total_amount), 0) as total_sales,
			$1::DATE::TEXT as date
		FROM sales_order
		WHERE order_date = $1
	`
	err := r.db.Get(&result, query, date)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// TopCustomer represents a customer with their total purchase amount
type TopCustomer struct {
	Username      string  `db:"username" json:"username"`
	FirstName     string  `db:"first_name" json:"first_name"`
	LastName      string  `db:"last_name" json:"last_name"`
	TotalPurchase float64 `db:"total_purchase" json:"total_purchase"`
}

// GetTopCustomers gets top 5 customers by total purchase amount for the last N months
func (r *ReportRepo) GetTopCustomers(months int) ([]TopCustomer, error) {
	var customers []TopCustomer
	query := `
		SELECT 
			c.username,
			c.first_name,
			c.last_name,
			COALESCE(SUM(so.total_amount), 0) as total_purchase
		FROM customer c
		LEFT JOIN sales_order so ON c.username = so.username
			AND so.order_date >= CURRENT_DATE - INTERVAL '1 month' * $1
		GROUP BY c.username, c.first_name, c.last_name
		ORDER BY total_purchase DESC
		LIMIT 5
	`
	err := r.db.Select(&customers, query, months)
	if err != nil {
		return nil, err
	}
	return customers, nil
}

// TopSellingBook represents a book with its total sales quantity
type TopSellingBook struct {
	ISBN        string  `db:"isbn" json:"isbn"`
	Title       string  `db:"title" json:"title"`
	AuthorName  string  `db:"author_name" json:"author_name"`
	TotalSold   int     `db:"total_sold" json:"total_sold"`
	TotalAmount float64 `db:"total_amount" json:"total_amount"`
}

// GetTopSellingBooks gets top 10 selling books by quantity for the last N months
func (r *ReportRepo) GetTopSellingBooks(months int) ([]TopSellingBook, error) {
	var books []TopSellingBook
	query := `
		SELECT 
			b.isbn,
			b.title,
			b.author_name,
			COALESCE(SUM(oi.quantity), 0)::INTEGER as total_sold,
			COALESCE(SUM(oi.quantity * oi.unit_price), 0) as total_amount
		FROM book b
		LEFT JOIN order_item oi ON b.isbn = oi.isbn
		LEFT JOIN sales_order so ON oi.order_id = so.order_id
			AND so.order_date >= CURRENT_DATE - INTERVAL '1 month' * $1
		GROUP BY b.isbn, b.title, b.author_name
		ORDER BY total_sold DESC
		LIMIT 10
	`
	err := r.db.Select(&books, query, months)
	if err != nil {
		return nil, err
	}
	return books, nil
}

// BookOrderCount represents the count of times a book has been ordered
type BookOrderCount struct {
	ISBN       string `db:"isbn" json:"isbn"`
	Title      string `db:"title" json:"title"`
	OrderCount int    `db:"order_count" json:"order_count"`
}

// GetBookOrderCount gets the total number of times a specific book has been ordered
func (r *ReportRepo) GetBookOrderCount(isbn string) (*BookOrderCount, error) {
	var result BookOrderCount
	query := `
		SELECT 
			b.isbn,
			b.title,
			COUNT(po.rep_order_id)::INTEGER as order_count
		FROM book b
		LEFT JOIN publisher_order po ON b.isbn = po.isbn
		WHERE b.isbn = $1
		GROUP BY b.isbn, b.title
	`
	err := r.db.Get(&result, query, isbn)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

