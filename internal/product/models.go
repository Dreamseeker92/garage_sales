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
	Sold     int
	Revenue  int
	Sales    []Sale `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
}

// Provides uuid before persistence to the storage.
func (p *Product) BeforeCreate(tx *gorm.DB) (err error) {
	p.ID, err = uuid.NewV4()
	return
}

// Sale represents an item of a transaction where some amount of a product was sold.
// Note: due to haggling the Paid value might not equal Quantity sold.
type Sale struct {
	gorm.Model
	Quantity  int       // Number of units sold
	Paid      int       // Total price
	ProductID uuid.UUID `gorm:"type:uuid;not null" json:"product_id"`
}
