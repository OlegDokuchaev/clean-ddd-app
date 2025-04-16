package product

import (
	"context"
	"github.com/google/uuid"
	"io"
	"warehouse/internal/domain/uow"
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
