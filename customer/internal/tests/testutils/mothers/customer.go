package mothers

import (
	"time"

	customerDomain "customer/internal/domain/customer"
	"customer/internal/tests/testutils/builders"
)

func DefaultCustomer() *customerDomain.Customer {
	return builders.NewCustomerBuilder().Build()
}

func CustomerWithPassword(password string) *customerDomain.Customer {
	return builders.NewCustomerBuilder().
		WithPassword(password).
		Build()
}

func CustomerWithLogin(phone string, password string) *customerDomain.Customer {
	return builders.NewCustomerBuilder().
		WithPhone(phone).
		WithPassword(password).
		Build()
}

func LockedCustomer() *customerDomain.Customer {
	return builders.NewCustomerBuilder().
		LockedFor(30 * time.Minute).
		Build()
}

func LockedCustomerWithLogin(phone string, password string) *customerDomain.Customer {
	return builders.NewCustomerBuilder().
		WithPhone(phone).
		WithPassword(password).
		LockedFor(30 * time.Minute).
		Build()
}

func LockedUntilCustomer(t time.Time) *customerDomain.Customer {
	return builders.NewCustomerBuilder().
		LockedUntil(t).
		Build()
}

func CustomerWithEmail(email string) *customerDomain.Customer {
	return builders.NewCustomerBuilder().
		WithEmail(email).
		Build()
}

func CustomerWithPhone(phone string) *customerDomain.Customer {
	return builders.NewCustomerBuilder().
		WithPhone(phone).
		Build()
}
