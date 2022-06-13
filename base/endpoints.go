package base

import (
	"context"
	"errors"
	"events/model"

	"github.com/go-kit/kit/endpoint"
)

//Endpoints ...
type Endpoints struct {
	Check       endpoint.Endpoint
	GetEvents   endpoint.Endpoint
	UpdateEvent endpoint.Endpoint
}

//MakeServerEndpoints ...
func MakeServerEndpoints(s Service) Endpoints {
	return Endpoints{
		Check:       MakeCheck(s),
		GetEvents:   MakeGetEvents(s),
		UpdateEvent: MakeUpdateEvent(s),
	}
}

//MakeCheck ...
func MakeCheck(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return s.Check(ctx)
	}
}
func MakeGetEvents(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		username, ok := request.(string)
		if !ok {
			return nil, errors.New("Bad Request")
		}
		return s.GetEvents(ctx, username)
	}
}

func MakeUpdateEvent(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		eventreq, ok := request.(model.UpdateEventRequest)
		if !ok {
			return nil, errors.New("Bad Request")
		}
		return "", s.UpdateEvents(ctx, eventreq)
	}
}
