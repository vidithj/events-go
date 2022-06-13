package base

import (
	"context"
	"events/model"
	"time"

	"github.com/go-kit/kit/metrics"
)

type instrumentingService struct {
	labelNames     []string
	requestCount   metrics.Counter
	errCount       metrics.Counter
	requestLatency metrics.Histogram
	next           Service
}

//NewInstrumentingService ...
func NewInstrumentingService(labelNames []string, counter metrics.Counter, errCounter metrics.Counter, latency metrics.Histogram,
	s Service) Service {
	return instrumentingService{
		labelNames:     labelNames,
		requestCount:   counter,
		errCount:       errCounter,
		requestLatency: latency,
		next:           s,
	}
}

func (s instrumentingService) Check(ctx context.Context) (res bool, err error) {
	defer func(begin time.Time) {
		s.instrument(begin, "Check", err)
	}(time.Now())
	return s.next.Check(ctx)
}

func (s instrumentingService) instrument(begin time.Time, methodName string, err error) {
	if len(s.labelNames) > 0 {
		s.requestCount.With(s.labelNames[0], methodName).Add(1)
		s.requestLatency.With(s.labelNames[0], methodName).Observe(time.Since(begin).Seconds())
		if err != nil {
			s.errCount.With(s.labelNames[0], methodName).Add(1)
		}
	}
}

func (s instrumentingService) GetEvents(ctx context.Context, username string) (resp model.Events, err error) {
	defer func(begin time.Time) {
		s.instrument(begin, "GetEvents", err)
	}(time.Now())
	return s.next.GetEvents(ctx, username)
}

func (s instrumentingService) UpdateEvents(ctx context.Context, request model.UpdateEventRequest) (err error) {
	defer func(begin time.Time) {
		s.instrument(begin, "UpdateEvents", err)
	}(time.Now())
	return s.next.UpdateEvents(ctx, request)
}
