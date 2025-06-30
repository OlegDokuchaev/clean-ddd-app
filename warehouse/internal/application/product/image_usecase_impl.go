package product

import (
	"context"
	"io"
	"warehouse/internal/domain/uow"

	"github.com/google/uuid"
)

type ImageUseCaseImpl struct {
	uow          uow.UoW
	imageService ImageService
}

func NewImageUseCase(uow uow.UoW, imageService ImageService) *ImageUseCaseImpl {
	return &ImageUseCaseImpl{
		uow:          uow,
		imageService: imageService,
	}
}

func (u *ImageUseCaseImpl) GetByID(ctx context.Context, productID uuid.UUID) (io.ReadCloser, string, error) {
	product, err := u.uow.Product().GetByID(ctx, productID)
	if err != nil {
		return nil, "", err
	}
	return u.imageService.Get(ctx, product.Image.Path)
}

func (u *ImageUseCaseImpl) UpdateByID(
	ctx context.Context,
	productID uuid.UUID,
	fileReader io.Reader,
	contentType string,
) error {
	product, err := u.uow.Product().GetByID(ctx, productID)
	if err != nil {
		return err
	}
	return u.imageService.Update(ctx, product.Image.Path, fileReader, contentType)
}

var _ ImageUseCase = (*ImageUseCaseImpl)(nil)
