package models

import "gorm.io/gorm"

type StoreProducts struct {
	gorm.Model
	StoreId     int  `json:"store_id"`
	ProductId   int  `json:"product_id"`
	IsAvailable bool `json:"is_available"`
}
