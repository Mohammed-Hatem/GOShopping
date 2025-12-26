package models

type SalesOrder struct {
	OrderID      int     `db:"order_id" json:"order_id"`
	OrderDate    string  `db:"order_date" json:"order_date"`
	TotalAmt     float64 `db:"total_amount" json:"total_amount"`
	CreditCardNo string  `db:"credit_card_no" json:"credit_card_no"`
	ExpiryDate   string  `db:"expiry_date" json:"expiry_date"`
	Username     string  `db:"username" json:"username"`
}

type OrderItem struct {
	OrderID int     `db:"order_id" json:"order_id"`
	ISBN    string  `db:"isbn" json:"isbn"`
	Quantity int    `db:"quantity" json:"quantity"`
	Price   float64 `db:"unit_price" json:"unit_price"`
}

