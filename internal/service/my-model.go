package service

import (
	"context"
	"go-start-template/internal/domain"
	"log/slog"
)

type myModelRepo interface {
	Create(ctx context.Context, params domain.CreateMyModelParams) (int32, error)
	FindOne(ctx context.Context, id int32) (domain.MyModel, error)
}

type myModelSrv struct {
	log  *slog.Logger
	repo myModelRepo
}

func NewMyModelSrv(log *slog.Logger, repo myModelRepo) *myModelSrv {
	return &myModelSrv{
		log:  log,
		repo: repo,
	}
}

func (srv *myModelSrv) Create(ctx context.Context, params domain.CreateMyModelParams) (int32, error) {
	// Some other business logic

	id, err := srv.repo.Create(ctx, params)
	return id, err
}

func (srv *myModelSrv) FindOne(ctx context.Context, id int32) (domain.MyModel, error) {
	// Some other business logic

	return srv.repo.FindOne(ctx, id)
}
