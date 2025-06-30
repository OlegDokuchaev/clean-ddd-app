package tables

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"time"
)

type Product struct {
	ID      uuid.UUID `gorm:"primaryKey"`
	Name    string
	Price   decimal.Decimal `gorm:"type:numeric(10, 2)"`
	Created time.Time
	Image   ProductImage `gorm:"foreignKey:ProductID;"`
}
