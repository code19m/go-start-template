// mongostore package implements the repository pattern for the MongoDB database.
// Located in the internal/repository/mongostore directory.
// It is called "mongostore" because "mongo" is the name of the official MongoDB Go driver.
package mongostore

import (
	"context"
	"fmt"
	"go-start-template/internal/config"
	"time"

	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	// Collection names
	userLogsCollName = "user_logs"
)

type mongoStore struct {
	client *mongo.Client
	db     *mongo.Database

	// Collections
	userLogsColl *mgm.Collection
}

// New creates a new instance of the mongoStore.
// It returns a pointer to the mongoStore and an error if any.
// Connection timeout is set to 5 seconds.
func New(cfg *config.Mongo) (*mongoStore, error) {
	const op = "mongostore.New"

	ctx, cancel := context.WithTimeout(context.Background(), 7*time.Second)
	defer cancel()

	uri := fmt.Sprintf("mongodb://%s:%s@%s:%d/%s?authSource=admin",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Db)

	mongoClient, err := mongo.Connect(ctx, options.Client().
		ApplyURI(uri).
		SetConnectTimeout(time.Second*5))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	err = mongoClient.Ping(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	db := mongoClient.Database(cfg.Db)

	mongoStore := &mongoStore{
		client: mongoClient,
		db:     db,

		// Initialize collections
		userLogsColl: mgm.NewCollection(db, userLogsCollName),
	}

	err = mongoStore.applyIndexes(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return mongoStore, nil
}

// Close closes the mongoStore client.
func (m *mongoStore) Close() error {
	return m.client.Disconnect(context.Background())
}

// applyIndexes applies indexes to the collections.
// It is idempotent so it won't return create errors if indexes already exist.
func (m *mongoStore) applyIndexes(ctx context.Context) error {
	const op = "mongostore.applyIndexes"

	// This indexes covers most of the queries for the user_logs collection.
	// ("created_at": -1)
	// ("user_id": 1, "created_at": -1)
	// ("user_pinfl": 1, "created_at": -1)
	// ("organization_id": 1, "created_at": -1)
	// ("organization_tin": 1, "created_at": -1)
	// ("action": 1, "created_at": -1)
	_, err := m.userLogsColl.Indexes().
		CreateMany(ctx, []mongo.IndexModel{
			{
				Keys: bson.D{{Key: "created_at", Value: -1}},
			},
			{
				Keys: bson.D{{Key: "user_id", Value: 1}, {Key: "created_at", Value: -1}},
			},
			{
				Keys: bson.D{{Key: "user_pinfl", Value: 1}, {Key: "created_at", Value: -1}},
			},
			{
				Keys: bson.D{{Key: "organization_id", Value: 1}, {Key: "created_at", Value: -1}},
			},
			{
				Keys: bson.D{{Key: "organization_tin", Value: 1}, {Key: "created_at", Value: -1}},
			},
			{
				Keys: bson.D{{Key: "action", Value: 1}, {Key: "created_at", Value: -1}},
			},
		})
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
