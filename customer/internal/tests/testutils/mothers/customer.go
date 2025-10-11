package mothers

import (
	"time"

	customerDomain "customer/internal/domain/customer"
	"customer/internal/tests/testutils/builders"
	"github.com/google/uuid"
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

func CustomerWithFailedAttempts(n int) *customerDomain.Customer {
	return builders.NewCustomerBuilder().
		WithFailedCount(n).
		Build()
}

func CustomerWithEmail(email string) *customerDomain.Customer {
	return builders.NewCustomerBuilder().
		WithEmail(email).
		Build()
}

func ListOfCustomers(n int) []*customerDomain.Customer {
	out := make([]*customerDomain.Customer, 0, n)
	for i := 0; i < n; i++ {
		out = append(out, DefaultCustomer())
	}
	return out
}

func TechnicalCustomer() *customerDomain.Customer {
	return builders.NewCustomerBuilder().
		WithName("Tech User").
		WithEmail("tech@example.com").
		WithPhone("+79990000000").
		WithPassword("TechP@ssw0rd!").
		Unlocked().
		Build()
}

func CustomerWithId(id uuid.UUID) *customerDomain.Customer {
	return builders.NewCustomerBuilder().
		WithID(id).
		Build()
}
