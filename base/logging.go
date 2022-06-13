package base

import (
	"context"
	"events/model"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/transport/http"
	"github.com/gorilla/handlers"
)

//Middleware ...
type Middleware func(Service) Service

//NewLoggingMiddleware ...
func NewLoggingMiddleware(logger log.Logger) Middleware {
	return func(next Service) Service {
		return &loggingMiddleware{
			next:   next,
			logger: logger,
		}
	}
}

type loggingMiddleware struct {
	next   Service
	logger log.Logger
}

//NewPanicLogger implements the RecoveryHandler logger interface
func NewPanicLogger(logger log.Logger) handlers.RecoveryHandlerLogger {
	return panicLogger{
		logger,
	}
}

type panicLogger struct {
	log.Logger
}

//Println ....
func (pl panicLogger) Println(msgs ...interface{}) {
	for _, msg := range msgs {
		pl.Log("panic", msg)
	}
}

func (mw loggingMiddleware) Check(ctx context.Context) (res bool, err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "Check", "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.Check(ctx)
}

func cid(ctx context.Context) string {
	cid, _ := ctx.Value(http.ContextKeyRequestXRequestID).(string)
	return cid
}
func xff(ctx context.Context) string {
	xff, _ := ctx.Value(http.ContextKeyRequestXForwardedFor).(string)
	return xff
}

func (mw loggingMiddleware) GetEvents(ctx context.Context, username string) (resp model.Events, err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "GetEvents", "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.GetEvents(ctx, username)
}

func (mw loggingMiddleware) UpdateEvents(ctx context.Context, request model.UpdateEventRequest) (err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "UpdateEvents", "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.UpdateEvents(ctx, request)
}
