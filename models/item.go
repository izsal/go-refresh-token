package models

type Item struct {
	ID    uint    `json:"id" gorm:"primary_key"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}
