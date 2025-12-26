package models

type ShoppingCart struct {
    ID       int    `db:"cart_id" json:"cart_id"`
    Username string `db:"username" json:"username"`
}

type CartItem struct {
    CartID   int    `db:"cart_id"`
    ISBN     string `db:"isbn"`
    Quantity int    `db:"quantity"`
}
