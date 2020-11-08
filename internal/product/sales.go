package product

import (
	"context"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

// List gets all Sale instances from the database.
func ListSales(ctx context.Context, db *gorm.DB) ([]Sale, error) {
	var sales []Sale
	if db.WithContext(ctx).Find(&sales); db.Error != nil {
		return nil, db.Error
	}

	return sales, nil
}

// Fetches a Sale by a given id.
func FetchSale(ctx context.Context, db *gorm.DB, id int) (*Sale, error) {
	sale := new(Sale)
	if err := db.WithContext(ctx).Where("id = ?", id).First(sale).Error; err != nil {
		return nil, err
	}

	return sale, nil
}

// Persist given attributes of a Sale to the database.
func PersistSale(ctx context.Context, db *gorm.DB, newSale *Sale) error {
	if db.WithContext(ctx).Create(newSale); db.Error != nil {
		return errors.Wrapf(db.Error, "Persisting a sale %v", newSale)
	}

	return nil
}
