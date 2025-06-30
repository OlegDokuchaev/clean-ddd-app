package documents

type OrderItem struct {
	ProductID string `bson:"product_id"`
	Price     string `bson:"price"`
	Count     int    `bson:"count"`
}
