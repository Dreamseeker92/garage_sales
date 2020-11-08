package product

import (
	"context"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

// List gets all Products from the database.
func List(ctx context.Context, db *gorm.DB) ([]Product, error) {
	var products []Product

	db.WithContext(ctx).
		Table("products").
		Joins("inner join sales on products.id = sales.product_id").
		Select("products.*, sum(sales.quantity) as sold, sum(sales.paid) as revenue").
		Group("products.id").
		Find(&products)

	if db.Error != nil {
		return nil, db.Error
	}

	return products, nil
}

// Fetches a product by a given id
func Fetch(ctx context.Context, db *gorm.DB, id string) (*Product, error) {
	product := new(Product)
	if err := db.WithContext(ctx).Where("id = ?", id).First(product).Error; err != nil {
		return nil, err
	}

	return product, nil
}

func Persist(ctx context.Context, db *gorm.DB, newProduct *Product) (*Product, error) {
	if db.WithContext(ctx).Create(newProduct); db.Error != nil {
		return nil, errors.Wrapf(db.Error, "Persisting a product %v", newProduct)
	}

	return newProduct, nil
}

func Update(ctx context.Context, db *gorm.DB, id string, updateParams Product) error {
	var product Product
	fetchResult := db.WithContext(ctx).Where("id = ?", id).First(&product)
	if fetchResult.Error != nil {
		return errors.Wrapf(fetchResult.Error, "Fetching product with id %s", id)
	}

	if updateResult := db.Model(&product).Updates(updateParams); updateResult.Error != nil {
		return errors.Wrapf(updateResult.Error, "Updating a product with id %s", id)
	}

	return nil

}
