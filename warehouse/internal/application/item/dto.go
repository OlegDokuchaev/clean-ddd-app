package item

import (
	"github.com/google/uuid"
)

type CreateDto struct {
	ProductID uuid.UUID
	Count     int
}

type ReserveDto struct {
	ItemID uuid.UUID
	Count  int
}

type ReleaseDto struct {
	ItemID uuid.UUID
	Count  int
}
