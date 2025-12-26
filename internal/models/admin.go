package models

import "time"

type Admin struct {
	AdminID  int    `db:"admin_id" json:"admin_id"`
	Username string `db:"username" json:"username"`
	Password string `db:"password" json:"-"`
}

type PublisherOrder struct {
	ID        int       `db:"rep_order_id" json:"rep_order_id"`
	OrderDate time.Time `db:"order_date" json:"order_date"`
	Status    string    `db:"status" json:"status"`
	Quantity  int       `db:"quantity" json:"quantity"`
	ISBN      string    `db:"isbn" json:"isbn"`
	AdminID   int       `db:"admin_id" json:"admin_id"`
}
