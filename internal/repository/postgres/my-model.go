package postgres

import (
	"context"
	"go-start-template/internal/domain"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
)

func NewMyModelStore(log *slog.Logger, pool *pgxpool.Pool) *myModelStore {
	return &myModelStore{
		log:  log,
		pool: pool,
	}
}

type myModelStore struct {
	log  *slog.Logger
	pool *pgxpool.Pool
}

func (store *myModelStore) Create(ctx context.Context, params domain.CreateMyModelParams) (int32, error) {
	const createMyModelQuery = `
		INSERT INTO "my_models" (
			name,
			age
		) VALUES (
			$1, $2
		) RETURNING id
	`

	row := store.pool.QueryRow(ctx, createMyModelQuery, params.Name, params.Age)
	var id int32
	err := row.Scan(&id)

	if err != nil {
		return 0, errors.WithStack(err)
	}

	return id, nil
}

func (store *myModelStore) FindOne(ctx context.Context, id int32) (domain.MyModel, error) {
	const getUserOrganization = `
		SELECT
			id, name, age
		FROM
			"my_models"
		WHERE
			id = $1
	`

	row := store.pool.QueryRow(ctx, getUserOrganization, id)
	var myModel domain.MyModel
	err := row.Scan(&myModel)

	if err != nil {
		return myModel, errors.WithStack(err)
	}

	return myModel, nil
}

func (store *myModelStore) SomeOtherMethod(ctx context.Context) error {
	return nil
}
