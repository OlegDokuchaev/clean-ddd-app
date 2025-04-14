package tables

import "github.com/google/uuid"

type ProductImage struct {
	ID        uuid.UUID `gorm:"primaryKey"`
	ProductID uuid.UUID
	Path      string
}
