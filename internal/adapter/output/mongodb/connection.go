package mongodb

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Config contém as configurações para conexão com MongoDB.
type Config struct {
	URI      string
	Database string
	Timeout  time.Duration
}

// Connection representa a conexão com MongoDB.
type Connection struct {
	client   *mongo.Client
	database *mongo.Database
}

// NewConnection cria uma nova conexão com MongoDB.
func NewConnection(config Config) (*Connection, error) {
	ctx, cancel := context.WithTimeout(context.Background(), config.Timeout)
	defer cancel()

	clientOptions := options.Client().
		ApplyURI(config.URI).
		SetServerSelectionTimeout(config.Timeout)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to mongodb: %w", err)
	}

	// Verifica a conexão
	if err = client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("failed to ping mongodb: %w", err)
	}

	return &Connection{
		client:   client,
		database: client.Database(config.Database),
	}, nil
}

// Database retorna a instância do banco de dados.
func (c *Connection) Database() *mongo.Database {
	return c.database
}

// Close fecha a conexão com o MongoDB.
func (c *Connection) Close(ctx context.Context) error {
	if err := c.client.Disconnect(ctx); err != nil {
		return fmt.Errorf("failed to disconnect from mongodb: %w", err)
	}
	return nil
}

// Ping verifica se a conexão está ativa.
func (c *Connection) Ping(ctx context.Context) error {
	if err := c.client.Ping(ctx, nil); err != nil {
		return fmt.Errorf("mongodb ping failed: %w", err)
	}
	return nil
}
