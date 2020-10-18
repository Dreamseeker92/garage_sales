package product

import (
	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

// Product is an item we sell.
type Product struct {
	gorm.Model
	ID       uuid.UUID `gorm:"type=uuid;primaryKey"`
	Name     string
	Cost     int
	Quantity int
}

func (p *Product) BeforeCreate(tx *gorm.DB) (err error) {
	p.ID, err = uuid.NewV4()
	return
}