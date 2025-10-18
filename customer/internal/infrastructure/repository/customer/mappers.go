package customer

import (
	customerDomain "customer/internal/domain/customer"
	"customer/internal/infrastructure/db/tables"
)

func ToDomain(model *tables.Customer) *customerDomain.Customer {
	return &customerDomain.Customer{
		ID:                 model.ID,
		Name:               model.Name,
		Phone:              model.Phone,
		Email:              model.Email,
		Password:           model.Password,
		Created:            model.Created,
		FailedCount:        model.FailedCount,
		LockedUntil:        model.LockedUntil,
		PasswordUpdated:    model.PasswordUpdated,
		MustChangePassword: model.MustChangePassword,
	}
}

func ToModel(domain *customerDomain.Customer) *tables.Customer {
	return &tables.Customer{
		ID:                 domain.ID,
		Name:               domain.Name,
		Phone:              domain.Phone,
		Email:              domain.Email,
		Password:           domain.Password,
		Created:            domain.Created,
		FailedCount:        domain.FailedCount,
		LockedUntil:        domain.LockedUntil,
		PasswordUpdated:    domain.PasswordUpdated,
		MustChangePassword: domain.MustChangePassword,
	}
}
