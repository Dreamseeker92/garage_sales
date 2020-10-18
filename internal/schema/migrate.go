package schema

import (
	"garagesale/internal/product"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(&product.Product{})
}
