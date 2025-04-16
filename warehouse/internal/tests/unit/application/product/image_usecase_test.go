package product

import (
	"bytes"
	"context"
	"errors"
	"io"
	"testing"
	productApplication "warehouse/internal/application/product"
	productDomain "warehouse/internal/domain/product"
	"warehouse/internal/mocks"
	productMock "warehouse/internal/mocks/product"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type ProductImageUseCaseTestSuite struct {
	suite.Suite
	ctx context.Context
}

func (s *ProductImageUseCaseTestSuite) SetupTest() {
	s.ctx = context.Background()
}

func (s *ProductImageUseCaseTestSuite) createTestProduct() *productDomain.Product {
	product, _, err := productDomain.Create("Test Product", decimal.NewFromInt(10), "product.png")
	require.NoError(s.T(), err)
	return product
}

func (s *ProductImageUseCaseTestSuite) TestUpdateByID() {
	tests := []struct {
		name        string
		fileReader  io.Reader
		contentType string
		setup       func(uowMock *mocks.UoWMock, imageServiceMock *productMock.ImageServiceMock) uuid.UUID
		expectedErr error
	}{
		{
			name:        "Success",
			fileReader:  bytes.NewReader([]byte("test data")),
			contentType: "image/png",
			setup: func(uow *mocks.UoWMock, imageServiceMock *productMock.ImageServiceMock) uuid.UUID {
				product := s.createTestProduct()

				uow.ProductMock.On("GetByID", s.ctx, mock.Anything).
					Return(product, nil).Once()
				imageServiceMock.On("Update", s.ctx, product.Image.Path, mock.Anything, "image/png").
					Return(nil).Once()

				return product.ID
			},
			expectedErr: nil,
		},
		{
			name:        "Failure: Product not found",
			fileReader:  bytes.NewReader([]byte("test data")),
			contentType: "image/png",
			setup: func(uow *mocks.UoWMock, imageServiceMock *productMock.ImageServiceMock) uuid.UUID {
				uow.ProductMock.On("GetByID", s.ctx, mock.Anything).
					Return((*productDomain.Product)(nil), errors.New("product not found")).Once()
				return uuid.Nil
			},
			expectedErr: errors.New("product not found"),
		},
		{
			name:        "Failure: Image service update error",
			fileReader:  bytes.NewReader([]byte("test data")),
			contentType: "image/png",
			setup: func(uow *mocks.UoWMock, imageServiceMock *productMock.ImageServiceMock) uuid.UUID {
				product := s.createTestProduct()

				uow.ProductMock.On("GetByID", s.ctx, mock.AnythingOfType("uuid.UUID")).
					Return(product, nil).Once()
				imageServiceMock.On("Update", s.ctx, product.Image.Path, mock.Anything, "image/png").
					Return(errors.New("image update failed")).Once()

				return product.ID
			},
			expectedErr: errors.New("image update failed"),
		},
	}

	for _, tc := range tests {
		tc := tc
		s.Run(tc.name, func() {
			s.T().Parallel()
			uow := mocks.NewUowMock()
			imageServiceMock := &productMock.ImageServiceMock{}
			productID := tc.setup(uow, imageServiceMock)
			usecase := productApplication.NewImageUseCase(uow, imageServiceMock)

			err := usecase.UpdateByID(s.ctx, productID, tc.fileReader, tc.contentType)

			if tc.expectedErr == nil {
				require.NoError(s.T(), err)
			} else {
				require.Error(s.T(), err)
				require.EqualError(s.T(), err, tc.expectedErr.Error())
			}

			uow.AssertExpectations(s.T())
			imageServiceMock.AssertExpectations(s.T())
		})
	}
}

func (s *ProductImageUseCaseTestSuite) TestGetByID() {
	tests := []struct {
		name        string
		setup       func(uowMock *mocks.UoWMock, imageServiceMock *productMock.ImageServiceMock) (uuid.UUID, string)
		expectedErr error
	}{
		{
			name: "Success",
			setup: func(uow *mocks.UoWMock, imageServiceMock *productMock.ImageServiceMock) (uuid.UUID, string) {
				product := s.createTestProduct()
				expectedContentType := "image/png"
				readCloser := io.NopCloser(bytes.NewReader([]byte("dummy data")))

				uow.ProductMock.On("GetByID", s.ctx, product.ID).
					Return(product, nil).Once()
				imageServiceMock.On("Get", s.ctx, product.Image.Path).
					Return(readCloser, expectedContentType, nil).Once()

				return product.ID, expectedContentType
			},
			expectedErr: nil,
		},
		{
			name: "Failure: Product not found",
			setup: func(uow *mocks.UoWMock, imageServiceMock *productMock.ImageServiceMock) (uuid.UUID, string) {
				productID := uuid.New()
				uow.ProductMock.On("GetByID", s.ctx, productID).
					Return((*productDomain.Product)(nil), errors.New("product not found")).Once()
				return productID, ""
			},
			expectedErr: errors.New("product not found"),
		},
		{
			name: "Failure: Image service get error",
			setup: func(uow *mocks.UoWMock, imageServiceMock *productMock.ImageServiceMock) (uuid.UUID, string) {
				product := s.createTestProduct()

				uow.ProductMock.On("GetByID", s.ctx, product.ID).
					Return(product, nil).Once()
				imageServiceMock.On("Get", s.ctx, product.Image.Path).
					Return(io.NopCloser(nil), "", errors.New("image get failed")).Once()

				return product.ID, ""
			},
			expectedErr: errors.New("image get failed"),
		},
	}

	for _, tc := range tests {
		tc := tc
		s.Run(tc.name, func() {
			s.T().Parallel()
			uow := mocks.NewUowMock()
			imageServiceMock := &productMock.ImageServiceMock{}
			productID, expectedContentType := tc.setup(uow, imageServiceMock)
			usecase := productApplication.NewImageUseCase(uow, imageServiceMock)

			actualReadCloser, actualContentType, err := usecase.GetByID(s.ctx, productID)

			if tc.expectedErr != nil {
				require.Error(s.T(), err)
				require.EqualError(s.T(), err, tc.expectedErr.Error())
			} else {
				defer actualReadCloser.Close()

				require.NoError(s.T(), err)
				require.NotNil(s.T(), actualReadCloser)
				require.Equal(s.T(), expectedContentType, actualContentType)
			}

			uow.AssertExpectations(s.T())
			imageServiceMock.AssertExpectations(s.T())
		})
	}
}

func TestProductImageUseCaseTestSuite(t *testing.T) {
	suite.Run(t, new(ProductImageUseCaseTestSuite))
}
