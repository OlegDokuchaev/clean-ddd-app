package courier

import (
	"context"
	courierDomain "courier/internal/domain/courier"
	"github.com/google/uuid"
	"math/rand"
)

type UseCaseImpl struct {
	repo courierDomain.Repository
}

func (u *UseCaseImpl) AssignOrder(ctx context.Context, _ uuid.UUID) (uuid.UUID, error) {
	orders, err := u.repo.GetAll(ctx)
	if err != nil {
		return uuid.Nil, err
	}
	if len(orders) == 0 {
		return uuid.Nil, ErrAvailableCourierNotFound
	}

	orderIndex := rand.Intn(len(orders))
	selectedOrder := orders[orderIndex]

	return selectedOrder.ID, nil
}

var _ UseCase = (*UseCaseImpl)(nil)
