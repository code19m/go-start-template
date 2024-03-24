//go:build integration

package mongostore_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"go-start-template/internal/config"
	"go-start-template/internal/repository/mongostore"

	tcmongo "github.com/romnn/testcontainers/mongo"
)

func setupMongoDBContainer(t *testing.T) (*config.Mongo, func()) {
	ctx := context.Background()

	container, err := tcmongo.Start(ctx, tcmongo.Options{
		ImageTag: "7.0",
		User:     "test",
		Password: "testpwd",
	})
	if err != nil {
		t.Fatal(err)
	}

	cfg := &config.Mongo{
		Host:     container.Host,
		Port:     int32(container.Port),
		User:     container.User,
		Password: container.Password,
		Db:       "testdatabase",
	}

	createRoleCmd := fmt.Sprintf(
		`db.createUser({user: '%s', pwd: '%s', roles: [{role: 'readWrite', db: '%s'}]})`,
		cfg.User, cfg.Password, cfg.Db,
	)

	// Create a test database and user
	_, _, err = container.Container.Exec(ctx, []string{"mongo", "--eval",
		createRoleCmd,
	})
	if err != nil {
		t.Fatal(err)
	}

	return cfg, func() { container.Terminate(ctx) } //nolint: errcheck
}

func TestMongoStoreNew(t *testing.T) {
	cfg, cleanup := setupMongoDBContainer(t)
	defer cleanup()

	mongoStore, err := mongostore.New(cfg)
	require.NoError(t, err)
	require.NotNil(t, mongoStore)
	require.NotEmpty(t, mongoStore)
}
