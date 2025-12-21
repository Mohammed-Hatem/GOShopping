package models

type Book struct {
	Isbn           string  `db:"isbn" json:"isbn"`
	Title          string  `db:"title" json:"title"`
	Pub_year       int     `db:"publication_year" json:"publication_year"`
	Selling_price  float64 `db:"selling_price" json:"selling_price"`
	Category       string  `db:"category" json:"category"`
	Stock_quantity int     `db:"stock_quantity" json:"stock_quantity"`
	Threshold      int     `db:"threshold" json:"threshold"`
	Publisher_id   int     `db:"publisher_id" json:"publisher_id"`
}
