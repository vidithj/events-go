package base

import (
	"context"
	"events/db"
	"events/model"

	"github.com/go-kit/kit/log"
)

//Service ...
type Service interface {
	Check(ctx context.Context) (bool, error)
	GetEvents(ctx context.Context, username string) (model.Events, error)
	UpdateEvents(ctx context.Context, request model.UpdateEventRequest) (err error)
}

type baseService struct {
	logger log.Logger
	dbs    *db.RoundRobin
}

//NewService ...
func NewService(l log.Logger, dbs *db.RoundRobin) Service {
	return baseService{
		logger: l,
		dbs:    dbs,
	}
}

//Check ...
func (s baseService) Check(ctx context.Context) (bool, error) {
	return true, nil
}

func (s baseService) GetEvents(ctx context.Context, username string) (model.Events, error) {
	db, err := s.dbs.DB()
	if err != nil {
		return model.Events{}, err
	}
	return db.GetEvents(username)
}

func (s baseService) UpdateEvents(ctx context.Context, request model.UpdateEventRequest) (err error) {
	db, err := s.dbs.DB()
	if err != nil {
		return err
	}
	return db.UpdateEvents(request)
}
