package models

type Author struct {
	Name string `db:"name" json:"name"`
	ISBN string `db:"isbn" json:"isbn"`
}
