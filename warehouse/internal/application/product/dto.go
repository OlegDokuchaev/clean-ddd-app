package product

import "github.com/shopspring/decimal"

type CreateDto struct {
	Name  string
	Price decimal.Decimal
}
