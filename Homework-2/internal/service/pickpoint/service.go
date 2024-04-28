//go:generate mockgen -source=./service.go -destination=./service_mocks_test.go -package=pickpoint
package pickpoint

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v4"
	"gitlab.ozon.dev/mer_marat/homework/internal/metrics"
	"gitlab.ozon.dev/mer_marat/homework/internal/model"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type storage interface {
	Add(context.Context, *model.PickPoint) (int64, error)
	GetByID(context.Context, int64) (*model.PickPoint, error)
	Update(context.Context, *model.PickPoint) error
	Delete(context.Context, int64) error
}

type cache interface {
	Set(ctx context.Context, key string, value interface{}) error
	Get(ctx context.Context, key string, value interface{}) error
	Delete(ctx context.Context, keys ...string) error
}

type transactor interface {
	RunSerializable(ctx context.Context, role pgx.TxAccessMode, f func(ctxTX context.Context) error) error
}

type counter interface {
	CounterInc()
}

type histogram interface {
	Observe(start time.Time)
}

type service struct {
	repo            storage
	cache           cache
	transactor      transactor
	counter         counter
	requestHandling histogram
}

// New returns type Service associated with storage
func NewService(stor storage, cache cache, transactor transactor) service {
	return service{
		repo:       stor,
		cache:      cache,
		transactor: transactor,
		counter:    &metrics.UnImplementedCounter{},
	}
}

func (s *service) AddRequestHistogram(hist histogram) {
	s.requestHandling = hist
}

func (s *service) AddCounterMetric(counter counter) {
	s.counter = counter
}

func (s service) Read(ctx context.Context, id int64) (res *model.PickPoint, err error) {
	start := time.Now()
	defer s.requestHandling.Observe(start)
	defer s.counter.CounterInc()
	ctx, span := otel.GetTracerProvider().Tracer("pickpoint").Start(ctx, "Read", trace.WithAttributes(
		attribute.String("request", fmt.Sprint(id)),
	))
	defer func() {
		if err != nil {
			span.SetAttributes(attribute.Bool("error", true))
			span.RecordError(err)
		} else {
			span.AddEvent("put done")
		}
		span.End()
	}()
	if !isValidID(ctx, id) {
		return nil, model.ErrorInvalidInput
	}
	point := new(model.PickPoint)
	err = s.cache.Get(ctx, fmt.Sprint(id), point)
	if err == nil {
		return point, nil
	}
	var pPoint *model.PickPoint
	if err := s.transactor.RunSerializable(ctx, pgx.ReadOnly, func(ctxTX context.Context) error {
		pPoint, err = s.repo.GetByID(ctx, id)
		if err != nil {
			return err
		}
		return s.cache.Set(ctx, fmt.Sprint(id), *pPoint)
	}); err != nil {
		return nil, err
	}
	return pPoint, nil
}

func (s service) Create(ctx context.Context, point *model.PickPoint) (res *model.PickPoint, err error) {
	start := time.Now()
	defer s.requestHandling.Observe(start)
	defer s.counter.CounterInc()
	ctx, span := otel.GetTracerProvider().Tracer("pickpoint").Start(ctx, "Create", trace.WithAttributes(
		attribute.String("request", fmt.Sprint(point)),
	))
	defer func() {
		if err != nil {
			span.SetAttributes(attribute.Bool("error", true))
			span.RecordError(err)
		} else {
			span.AddEvent("put done")
		}
		span.End()
	}()
	id, err := s.repo.Add(ctx, point)
	if err != nil {
		return nil, err
	}
	point.ID = id
	return point, nil
}

func (s service) Update(ctx context.Context, point *model.PickPoint) (err error) {
	start := time.Now()
	defer s.requestHandling.Observe(start)
	defer s.counter.CounterInc()
	ctx, span := otel.GetTracerProvider().Tracer("pickpoint").Start(ctx, "Update", trace.WithAttributes(
		attribute.String("request", fmt.Sprint(point)),
	))
	defer func() {
		if err != nil {
			span.SetAttributes(attribute.Bool("error", true))
			span.RecordError(err)
		} else {
			span.AddEvent("put done")
		}
		span.End()
	}()
	if !isValidPickPoint(ctx, point) {
		return model.ErrorInvalidInput
	}
	return s.transactor.RunSerializable(ctx, pgx.ReadWrite, func(ctxTX context.Context) error {
		err := s.repo.Update(ctx, point)
		if err != nil {
			return err
		}
		return s.cache.Delete(ctx, fmt.Sprint(point.ID))
	})
}

func (s service) Delete(ctx context.Context, id int64) (err error) {
	start := time.Now()
	defer s.requestHandling.Observe(start)
	defer s.counter.CounterInc()
	ctx, span := otel.GetTracerProvider().Tracer("pickpoint").Start(ctx, "Deelete", trace.WithAttributes(
		attribute.String("request", fmt.Sprint(id)),
	))
	defer func() {
		if err != nil {
			span.SetAttributes(attribute.Bool("error", true))
			span.RecordError(err)
		} else {
			span.AddEvent("put done")
		}
		span.End()
	}()
	if !isValidID(ctx, id) {
		return model.ErrorInvalidInput
	}
	return s.transactor.RunSerializable(ctx, pgx.ReadWrite, func(ctxTX context.Context) error {
		err := s.repo.Delete(ctx, id)
		if err != nil {
			return err
		}
		return s.cache.Delete(ctx, fmt.Sprint(id))
	})
}

func isValidPickPoint(ctx context.Context, point *model.PickPoint) bool {
	_, span := otel.GetTracerProvider().Tracer("pickpoint").Start(ctx, "isValidPickPoint")
	defer func() {
		span.AddEvent("put done")
		span.End()
	}()
	return point.ID > 0
}

func isValidID(ctx context.Context, id int64) bool {
	_, span := otel.GetTracerProvider().Tracer("pickpoint").Start(ctx, "isValidPickPoint")
	defer func() {
		span.AddEvent("put done")
		span.End()
	}()
	return id > 0
}
