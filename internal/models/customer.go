package models



type Customer struct {
	Username string `db:"username" json:"username"`
	Password string `db:"password" json:"-"`
	FName     string `db:"first_name" json:"fname"`
	LName     string `db:"last_name" json:"lname"`
	Email    string `db:"email" json:"email"`
	Phone    string `db:"phone_number" json:"phone"`
	Address  string `db:"shipping_address" json:"address"`

}


