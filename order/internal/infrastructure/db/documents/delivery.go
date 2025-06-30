package documents

import (
	"time"
)

type Delivery struct {
	CourierID *string    `bson:"courier_id,omitempty"`
	Address   string     `bson:"address"`
	Arrived   *time.Time `bson:"arrived,omitempty"`
}
