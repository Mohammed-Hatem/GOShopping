package models


type Publisher struct {
	PublisherID int    `db:"publisher_id" json:"publisher_id"`
	Name        string `db:"name" json:"name"`
	Address     string `db:"address" json:"address"`
	Phone       string `db:"phone" json:"phone"`
}
