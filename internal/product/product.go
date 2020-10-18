package product

import (
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

// List gets all Products from the database.
func List(db *gorm.DB) ([]Product, error) {
	var products []Product
	if db.Find(&products); db.Error != nil {
		return nil, db.Error
	}

	return products, nil
}

// Fetches a product by a given id
func Fetch(db *gorm.DB, id string) (*Product, error) {
	product := new(Product)
	if err := db.Where("id = ?", id).First(product).Error; err != nil {
		return nil, err
	}

	return product, nil
}

func Persist(db *gorm.DB, newProduct *Product) (*Product, error) {
	if db.Create(newProduct); db.Error != nil {
		return nil, errors.Wrapf(db.Error, "Persisting a product %v", newProduct)
	}
	
	return newProduct, nil
}