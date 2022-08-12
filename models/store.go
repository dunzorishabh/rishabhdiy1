package models

import "gorm.io/gorm"

type Store struct {
	gorm.Model
	StoreId   int    `json:"store_id"`
	StoreName string `json:"name"`
}
