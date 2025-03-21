package item

import (
	"github.com/google/uuid"
)

type CreateDto struct {
	ProductID uuid.UUID
	Count     int
}

type ReserveDto struct {
	Items []ItemDto
}

type ReleaseDto struct {
	Items []ItemDto
}

type ItemDto struct {
	ProductID uuid.UUID
	Count     int
}
