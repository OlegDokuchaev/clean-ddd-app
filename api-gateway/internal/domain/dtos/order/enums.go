package order

type (
	Status string
)

const (
	Created                 Status = "created"
	CanceledCourierNotFound Status = "canceled_courier_not_found"
	CanceledOutOfStock      Status = "canceled_out_of_stock"
	Delivering              Status = "delivering"
	Delivered               Status = "delivered"
	CustomerCanceled        Status = "customer_canceled"
)
