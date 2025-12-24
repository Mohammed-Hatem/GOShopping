package models

type Book struct {
	Isbn          string  `db:"isbn" json:"isbn"`
	Title         string  `db:"title" json:"title"`
	PubYear       int     `db:"publication_year" json:"publication_year"`
	SellingPrice  float64 `db:"selling_price" json:"selling_price"`
	Category      string  `db:"category" json:"category"`
	AuthorName    string  `db:"author_name" json:"author_name"`
	StockQuantity int     `db:"stock_quantity" json:"stock_quantity"`
	Threshold     int     `db:"threshold" json:"threshold"`
	PublisherId   int     `db:"publisher_id" json:"publisher_id"`
}
