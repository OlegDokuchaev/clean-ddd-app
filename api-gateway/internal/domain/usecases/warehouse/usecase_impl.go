package warehouse

import (
	warehouseDto "api-gateway/internal/domain/dtos/warehouse"
	"api-gateway/internal/port/output/auth/admin"
	warehouseClient "api-gateway/internal/port/output/clients/warehouse"
	"context"

	"github.com/google/uuid"
)

type UseCaseImpl struct {
	adminAuth       admin.Auth
	warehouseClient warehouseClient.Client
}

func NewUseCase(
	adminAuth admin.Auth,
	warehouseClient warehouseClient.Client,
) UseCase {
	return &UseCaseImpl{
		adminAuth:       adminAuth,
		warehouseClient: warehouseClient,
	}
}

func (u *UseCaseImpl) ReserveItems(ctx context.Context, items []warehouseDto.ItemInfoDto, adminToken string) error {
	if !u.adminAuth.Validate(adminToken) {
		return ErrUnauthorized
	}

	err := u.warehouseClient.ReserveItems(ctx, items)
	if err != nil {
		return err
	}

	return nil
}

func (u *UseCaseImpl) ReleaseItems(ctx context.Context, items []warehouseDto.ItemInfoDto, adminToken string) error {
	if !u.adminAuth.Validate(adminToken) {
		return ErrUnauthorized
	}

	err := u.warehouseClient.ReleaseItems(ctx, items)
	if err != nil {
		return err
	}

	return nil
}

func (u *UseCaseImpl) CreateProduct(ctx context.Context, data warehouseDto.CreateProductDto, adminToken string) (uuid.UUID, error) {
	if !u.adminAuth.Validate(adminToken) {
		return uuid.Nil, ErrUnauthorized
	}

	productID, err := u.warehouseClient.CreateProduct(ctx, data)
	if err != nil {
		return uuid.Nil, err
	}

	return productID, nil
}

func (u *UseCaseImpl) GetAllItems(ctx context.Context) ([]*warehouseDto.ItemDto, error) {
	items, err := u.warehouseClient.GetAllItems(ctx)
	if err != nil {
		return nil, err
	}
	return items, nil
}

var _ UseCase = (*UseCaseImpl)(nil)
