package mongodb

import (
	"context"
	"errors"
	"fmt"

	"github.com/diiineeei/neoway-backend/internal/domain/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const ordersCollection = "orders"

// OrderDocument representa o documento do pedido no MongoDB.
type OrderDocument struct {
	OrderID   string             `bson:"order_id"`
	Product   string             `bson:"product"`
	Quantity  int                `bson:"quantity"`
	Status    entity.OrderStatus `bson:"status"`
	CreatedAt int64              `bson:"created_at"` // Unix timestamp
}

// OrderRepositoryAdapter é o adaptador de saída para MongoDB.
// Implementa port.OrderRepository.
type OrderRepositoryAdapter struct {
	collection *mongo.Collection
}

// NewOrderRepositoryAdapter cria um novo adapter de repository a partir de um database.
func NewOrderRepositoryAdapter(db *mongo.Database) *OrderRepositoryAdapter {
	return &OrderRepositoryAdapter{
		collection: db.Collection(ordersCollection),
	}
}

// NewOrderRepositoryAdapterFromCollection cria um novo adapter de repository a partir de uma coleção.
func NewOrderRepositoryAdapterFromCollection(collection *mongo.Collection) *OrderRepositoryAdapter {
	return &OrderRepositoryAdapter{
		collection: collection,
	}
}

// Save persiste um pedido no MongoDB.
func (r *OrderRepositoryAdapter) Save(ctx context.Context, order *entity.Order) error {
	doc := OrderDocument{
		OrderID:   order.ID(),
		Product:   order.Product(),
		Quantity:  order.Quantity(),
		Status:    order.Status(),
		CreatedAt: order.CreatedAt().Unix(),
	}

	_, err := r.collection.InsertOne(ctx, doc)
	if err != nil {
		return fmt.Errorf("failed to insert order: %w", err)
	}

	return nil
}

// FindByID busca um pedido pelo ID.
func (r *OrderRepositoryAdapter) FindByID(ctx context.Context, orderID string) (*entity.Order, error) {
	filter := bson.M{"order_id": orderID}

	var doc OrderDocument
	err := r.collection.FindOne(ctx, filter).Decode(&doc)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, fmt.Errorf("order not found: %s", orderID)
		}
		return nil, fmt.Errorf("failed to find order: %w", err)
	}

	// Reconstrói a entidade de domínio
	order := entity.Reconstruct(
		doc.OrderID,
		doc.Product,
		doc.Quantity,
		doc.Status,
		doc.CreatedAt,
	)

	return order, nil
}

// UpdateStatus atualiza o status de um pedido.
func (r *OrderRepositoryAdapter) UpdateStatus(ctx context.Context, orderID string, status entity.OrderStatus) error {
	filter := bson.M{"order_id": orderID}
	update := bson.M{"$set": bson.M{"status": status}}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update order status: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("order not found: %s", orderID)
	}

	return nil
}
