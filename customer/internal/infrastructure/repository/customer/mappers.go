package customer

import (
	customerDomain "customer/internal/domain/customer"
	"customer/internal/infrastructure/db/tables"
)

func ToDomain(model *tables.Customer) *customerDomain.Customer {
	return &customerDomain.Customer{
		ID:       model.ID,
		Name:     model.Name,
		Phone:    model.Phone,
		Password: model.Password,
		Created:  model.Created,
	}
}

func ToModel(domain *customerDomain.Customer) *tables.Customer {
	return &tables.Customer{
		ID:       domain.ID,
		Name:     domain.Name,
		Phone:    domain.Phone,
		Password: domain.Password,
		Created:  domain.Created,
	}
}
