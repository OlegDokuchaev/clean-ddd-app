package domain

import (
	"testing"
	productDomain "warehouse/internal/domain/product"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type ProductDomainTestSuite struct {
	suite.Suite
}

func (p *ProductDomainTestSuite) TestCreate() {
	tests := []struct {
		name        string
		pName       string
		price       decimal.Decimal
		expectedErr error
	}{
		{
			name:        "Success",
			pName:       "test",
			price:       decimal.NewFromInt(1),
			expectedErr: nil,
		},
		{
			name:        "Failure: Empty name",
			pName:       "",
			price:       decimal.NewFromInt(1),
			expectedErr: productDomain.ErrInvalidProductName,
		},
		{
			name:        "Failure: Name consists only of whitespaces",
			pName:       "   ",
			price:       decimal.NewFromInt(1),
			expectedErr: productDomain.ErrInvalidProductName,
		},
		{
			name:        "Failure: Price is zero",
			pName:       "test",
			price:       decimal.Zero,
			expectedErr: productDomain.ErrInvalidProductPrice,
		},
		{
			name:        "Failure: Price less then zero",
			pName:       "test",
			price:       decimal.NewFromInt(-1),
			expectedErr: productDomain.ErrInvalidProductPrice,
		},
	}

	for _, tc := range tests {
		tc := tc
		p.Run(tc.name, func() {
			p.T().Parallel()

			product, events, err := productDomain.Create(tc.pName, tc.price)

			if tc.expectedErr != nil {
				require.Error(p.T(), err)
				require.ErrorIs(p.T(), err, tc.expectedErr)
			} else {
				require.NoError(p.T(), err)
				require.NotNil(p.T(), product)
				require.Len(p.T(), events, 1)
				require.IsType(p.T(), productDomain.CreateEvent{}, events[0])
			}
		})
	}
}

func TestProductDomainTestSuite(t *testing.T) {
	suite.Run(t, new(ProductDomainTestSuite))
}
