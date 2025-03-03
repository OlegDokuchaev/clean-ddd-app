package domain

import (
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"testing"
	itemDomain "warehouse/internal/domain/item"
	productDomain "warehouse/internal/domain/product"
)

type ItemDomainTestSuite struct {
	suite.Suite
}

func (i *ItemDomainTestSuite) TestCreate() {
	tests := []struct {
		name        string
		setup       func() *productDomain.Product
		count       int
		expectedErr error
	}{
		{
			name: "Success",
			setup: func() *productDomain.Product {
				product, _, err := productDomain.Create("test", decimal.NewFromInt(1))
				require.NoError(i.T(), err)
				return product
			},
			count:       10,
			expectedErr: nil,
		},
		{
			name: "Failure: Zero count",
			setup: func() *productDomain.Product {
				product, _, err := productDomain.Create("test", decimal.NewFromInt(1))
				require.NoError(i.T(), err)
				return product
			},
			count:       0,
			expectedErr: itemDomain.ErrInvalidItemCount,
		},
		{
			name: "Failure: Negative count",
			setup: func() *productDomain.Product {
				product, _, err := productDomain.Create("test", decimal.NewFromInt(1))
				require.NoError(i.T(), err)
				return product
			},
			count:       -10,
			expectedErr: itemDomain.ErrInvalidItemCount,
		},
	}

	for _, tc := range tests {
		tc := tc
		i.Run(tc.name, func() {
			i.T().Parallel()
			product := tc.setup()

			item, err := itemDomain.Create(product, tc.count)

			if tc.expectedErr != nil {
				require.Error(i.T(), err)
				require.ErrorIs(i.T(), err, tc.expectedErr)
			} else {
				require.NoError(i.T(), err)
				require.NotNil(i.T(), item)
			}
		})
	}
}

func (i *ItemDomainTestSuite) TestReserve() {
	tests := []struct {
		name          string
		setup         func() *itemDomain.Item
		reserveCount  int
		expectedCount int
		expectedErr   error
	}{
		{
			name: "Success",
			setup: func() *itemDomain.Item {
				product, _, err := productDomain.Create("test", decimal.NewFromInt(1))
				require.NoError(i.T(), err)
				item, err := itemDomain.Create(product, 10)
				require.NoError(i.T(), err)
				return item
			},
			reserveCount:  1,
			expectedCount: 9,
			expectedErr:   nil,
		},
		{
			name: "Failure: Zero reserve count",
			setup: func() *itemDomain.Item {
				product, _, err := productDomain.Create("test", decimal.NewFromInt(1))
				require.NoError(i.T(), err)
				item, err := itemDomain.Create(product, 10)
				require.NoError(i.T(), err)
				return item
			},
			reserveCount:  0,
			expectedCount: 10,
			expectedErr:   itemDomain.ErrInvalidItemCount,
		},
		{
			name: "Failure: Negative reserve count",
			setup: func() *itemDomain.Item {
				product, _, err := productDomain.Create("test", decimal.NewFromInt(1))
				require.NoError(i.T(), err)
				item, err := itemDomain.Create(product, 10)
				require.NoError(i.T(), err)
				return item
			},
			reserveCount:  -1,
			expectedCount: 10,
			expectedErr:   itemDomain.ErrInvalidItemCount,
		},
		{
			name: "Failure: Zero item count",
			setup: func() *itemDomain.Item {
				product, _, err := productDomain.Create("test", decimal.NewFromInt(1))
				require.NoError(i.T(), err)
				item, err := itemDomain.Create(product, 10)
				require.NoError(i.T(), err)
				return item
			},
			reserveCount:  10,
			expectedCount: 10,
			expectedErr:   itemDomain.ErrInvalidItemCount,
		},
		{
			name: "Failure: Negative item count",
			setup: func() *itemDomain.Item {
				product, _, err := productDomain.Create("test", decimal.NewFromInt(1))
				require.NoError(i.T(), err)
				item, err := itemDomain.Create(product, 10)
				require.NoError(i.T(), err)
				return item
			},
			reserveCount:  100,
			expectedCount: 10,
			expectedErr:   itemDomain.ErrInvalidItemCount,
		},
	}

	for _, tc := range tests {
		tc := tc
		i.Run(tc.name, func() {
			i.T().Parallel()
			item := tc.setup()

			err := item.Reserve(tc.reserveCount)

			if tc.expectedErr != nil {
				require.Error(i.T(), err)
				require.ErrorIs(i.T(), err, tc.expectedErr)
			} else {
				require.NoError(i.T(), err)
			}
			require.Equal(i.T(), tc.expectedCount, item.Count)
		})
	}
}

func (i *ItemDomainTestSuite) TestRelease() {
	tests := []struct {
		name          string
		setup         func() *itemDomain.Item
		releaseCount  int
		expectedCount int
		expectedErr   error
	}{
		{
			name: "Success",
			setup: func() *itemDomain.Item {
				product, _, err := productDomain.Create("test", decimal.NewFromInt(1))
				require.NoError(i.T(), err)
				item, err := itemDomain.Create(product, 10)
				require.NoError(i.T(), err)
				return item
			},
			releaseCount:  1,
			expectedCount: 11,
			expectedErr:   nil,
		},
		{
			name: "Failure: Zero reserve count",
			setup: func() *itemDomain.Item {
				product, _, err := productDomain.Create("test", decimal.NewFromInt(1))
				require.NoError(i.T(), err)
				item, err := itemDomain.Create(product, 10)
				require.NoError(i.T(), err)
				return item
			},
			releaseCount:  0,
			expectedCount: 10,
			expectedErr:   itemDomain.ErrInvalidItemCount,
		},
		{
			name: "Failure: Negative reserve count",
			setup: func() *itemDomain.Item {
				product, _, err := productDomain.Create("test", decimal.NewFromInt(1))
				require.NoError(i.T(), err)
				item, err := itemDomain.Create(product, 10)
				require.NoError(i.T(), err)
				return item
			},
			releaseCount:  -1,
			expectedCount: 10,
			expectedErr:   itemDomain.ErrInvalidItemCount,
		},
	}

	for _, tc := range tests {
		tc := tc
		i.Run(tc.name, func() {
			i.T().Parallel()
			item := tc.setup()

			err := item.Release(tc.releaseCount)

			if tc.expectedErr != nil {
				require.Error(i.T(), err)
				require.ErrorIs(i.T(), err, tc.expectedErr)
			} else {
				require.NoError(i.T(), err)
			}
			require.Equal(i.T(), tc.expectedCount, item.Count)
		})
	}
}

func TestItemDomainTestSuite(t *testing.T) {
	suite.Run(t, new(ItemDomainTestSuite))
}
