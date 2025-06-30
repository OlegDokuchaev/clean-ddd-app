package documents

import (
	orderDomain "order/internal/domain/order"
	"time"
)

type Order struct {
	ID         string             `bson:"_id"`
	CustomerID string             `bson:"customer_id"`
	Status     orderDomain.Status `bson:"status"`
	Created    time.Time          `bson:"created"`
	Version    string             `bson:"version"`
	Delivery   Delivery           `bson:"delivery"`
	Items      []OrderItem        `bson:"items"`
}
