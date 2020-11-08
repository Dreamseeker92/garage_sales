package schema

import (
	"garagesale/internal/product"
	"gorm.io/gorm"
)

var seeds = []product.Product{
	{
		Name:     "Comic Books",
		Cost:     50,
		Quantity: 42,
		Sales:    []product.Sale{{Quantity: 10, Paid: 100}},
	},
	{
		Name:     "McDonalds Toys",
		Cost:     76,
		Quantity: 100,
		Sales:    []product.Sale{{Quantity: 20, Paid: 200}},
	},
}

// Seed runs the set of seed-data queries against db. The queries are ran in a
// transaction and rolled back if any fail.
func Seed(db *gorm.DB) error {
	transaction := db.Begin()
	defer func() {
		if failure := recover(); failure != nil {
			transaction.Rollback()
		}
	}()

	transaction.Create(&seeds)
	if errTransaction := transaction.Error; errTransaction != nil {
		transaction.Rollback()
		return errTransaction
	}

	return transaction.Commit().Error
}
