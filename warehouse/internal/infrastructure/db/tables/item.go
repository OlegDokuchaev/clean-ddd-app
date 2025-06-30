package tables

import "github.com/google/uuid"

type Item struct {
	ID        uuid.UUID `gorm:"primaryKey"`
	ProductID uuid.UUID
	Count     int
	Version   uuid.UUID
	Product   Product `gorm:"foreignKey:ProductID;"`
}
