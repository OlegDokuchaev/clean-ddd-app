package order

import (
	"context"
	"order/internal/infrastructure/db/documents"

	orderDomain "order/internal/domain/order"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type RepositoryImpl struct {
	collection *mongo.Collection
}

func New(collection *mongo.Collection) *RepositoryImpl {
	return &RepositoryImpl{collection: collection}
}

func (r *RepositoryImpl) Create(ctx context.Context, order *orderDomain.Order) error {
	doc := toDoc(order)
	_, err := r.collection.InsertOne(ctx, doc)
	return ParseError(err)
}

func (r *RepositoryImpl) Update(ctx context.Context, order *orderDomain.Order) error {
	oldVersion := order.Version
	newVersion := uuid.New()
	order.Version = newVersion
	doc := toDoc(order)

	filter := bson.M{"_id": order.ID.String(), "version": oldVersion.String()}
	result, err := r.collection.ReplaceOne(ctx, filter, doc)
	if err != nil {
		return ParseError(err)
	}
	if result.MatchedCount == 0 {
		return ErrOrderNotFound
	}

	return nil
}

func (r *RepositoryImpl) GetByID(ctx context.Context, orderID uuid.UUID) (*orderDomain.Order, error) {
	filter := bson.M{"_id": orderID.String()}
	var doc documents.Order
	if err := r.collection.FindOne(ctx, filter).Decode(&doc); err != nil {
		return nil, ParseError(err)
	}
	return toDomain(&doc)
}

func (r *RepositoryImpl) GetAllByCustomer(ctx context.Context, customerID uuid.UUID) ([]*orderDomain.Order, error) {
	filter := bson.M{"customer_id": customerID.String()}
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, ParseError(err)
	}
	defer func() { _ = cursor.Close(ctx) }()

	var docs []documents.Order
	if err := cursor.All(ctx, &docs); err != nil {
		return nil, ParseError(err)
	}
	return toDomains(docs)
}

func (r *RepositoryImpl) GetCurrentByCourier(ctx context.Context, courierID uuid.UUID) ([]*orderDomain.Order, error) {
	filter := bson.M{"delivery.courier_id": courierID.String()}
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, ParseError(err)
	}
	defer func() { _ = cursor.Close(ctx) }()

	var docs []documents.Order
	if err := cursor.All(ctx, &docs); err != nil {
		return nil, ParseError(err)
	}
	return toDomains(docs)
}

var _ orderDomain.Repository = (*RepositoryImpl)(nil)
