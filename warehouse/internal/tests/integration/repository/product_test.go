//go:build integration

package repository

import (
	"context"
	"testing"
	productDomain "warehouse/internal/domain/product"
	"warehouse/internal/infrastructure/db/migrations"
	productRepository "warehouse/internal/infrastructure/repository/product"
	"warehouse/internal/tests/testutils"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type ProductRepositoryTestSuite struct {
	suite.Suite
	ctx    context.Context
	testDB *testutils.TestDB
}

func (s *ProductRepositoryTestSuite) SetupSuite() {
	config, err := migrations.NewConfig()
	require.NoError(s.T(), err)

	s.ctx = context.Background()

	s.testDB, err = testutils.NewTestDB(s.ctx, config)
	require.NoError(s.T(), err)
}

func (s *ProductRepositoryTestSuite) TearDownSuite() {
	if s.testDB != nil {
		err := s.testDB.Close(s.ctx)
		require.NoError(s.T(), err)
	}
}

func (s *ProductRepositoryTestSuite) getRepo() productDomain.Repository {
	return productRepository.New(s.testDB.DB)
}

func (s *ProductRepositoryTestSuite) createTestProduct() *productDomain.Product {
	product, _, err := productDomain.Create("Test Product", decimal.NewFromInt(100))
	require.NoError(s.T(), err)
	return product
}

func (s *ProductRepositoryTestSuite) TestCreate() {
	tests := []struct {
		name          string
		setup         func(repo productDomain.Repository) *productDomain.Product
		expectedError error
	}{
		{
			name: "Success",
			setup: func(repo productDomain.Repository) *productDomain.Product {
				return s.createTestProduct()
			},
			expectedError: nil,
		},
		{
			name: "Failure: Product already exists",
			setup: func(repo productDomain.Repository) *productDomain.Product {
				product := s.createTestProduct()
				err := repo.Create(s.ctx, product)
				require.NoError(s.T(), err)
				return product
			},
			expectedError: productRepository.ErrProductAlreadyExists,
		},
	}

	repo := s.getRepo()
	for _, tc := range tests {
		s.Run(tc.name, func() {
			product := tc.setup(repo)
			err := repo.Create(s.ctx, product)

			if tc.expectedError != nil {
				require.Error(s.T(), err)
				require.Equal(s.T(), tc.expectedError, err)
			} else {
				require.NoError(s.T(), err)

				createdProduct, err := repo.GetByID(s.ctx, product.ID)
				require.NoError(s.T(), err)
				require.NotNil(s.T(), createdProduct)
				require.Equal(s.T(), product.ID, createdProduct.ID)
			}
		})
	}
}

func (s *ProductRepositoryTestSuite) TestGetByID() {
	tests := []struct {
		name          string
		setup         func(repo productDomain.Repository) uuid.UUID
		expectedError error
	}{
		{
			name: "Success",
			setup: func(repo productDomain.Repository) uuid.UUID {
				product := s.createTestProduct()
				err := repo.Create(s.ctx, product)
				require.NoError(s.T(), err)
				return product.ID
			},
			expectedError: nil,
		},
		{
			name: "Failure: Product not found",
			setup: func(repo productDomain.Repository) uuid.UUID {
				return uuid.New()
			},
			expectedError: productRepository.ErrProductNotFound,
		},
	}

	repo := s.getRepo()
	for _, tc := range tests {
		s.Run(tc.name, func() {
			productID := tc.setup(repo)

			product, err := repo.GetByID(s.ctx, productID)

			if tc.expectedError != nil {
				require.Error(s.T(), err)
				require.Equal(s.T(), tc.expectedError, err)
			} else {
				require.NoError(s.T(), err)
				require.NotNil(s.T(), product)
				require.Equal(s.T(), productID, product.ID)
			}
		})
	}
}

func TestProductRepository(t *testing.T) {
	suite.Run(t, new(ProductRepositoryTestSuite))
}
