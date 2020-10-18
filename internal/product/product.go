package product

import (
	"gorm.io/gorm"
)

// List gets all Products from the database.
func List(db *gorm.DB) ([]Product, error) {
	var products []Product

	db.Find(&products)

	return products, nil
}
