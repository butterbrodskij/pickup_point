//go:generate mockgen -source=./service.go -destination=./service_mocks_test.go -package=pickpoint
package pickpoint

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"
	"github.com/prometheus/client_golang/prometheus"
	"gitlab.ozon.dev/mer_marat/homework/internal/model"
	pickpoint_pb "gitlab.ozon.dev/mer_marat/homework/internal/pkg/pb/pickpoint"
	"google.golang.org/protobuf/types/known/emptypb"
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

type service struct {
	pickpoint_pb.UnimplementedPickPointsServer
	repo       storage
	cache      cache
	transactor transactor
	counter    prometheus.Counter
}

// New returns type Service associated with storage
func NewService(stor storage, cache cache, transactor transactor, counter prometheus.Counter) service {
	return service{
		repo:       stor,
		cache:      cache,
		transactor: transactor,
		counter:    counter,
	}
}

func pb2Model(point *pickpoint_pb.PickPoint) *model.PickPoint {
	return &model.PickPoint{
		ID:      point.Id,
		Name:    point.Name,
		Address: point.Address,
		Contact: point.Contact,
	}
}

func model2Pb(point *model.PickPoint) *pickpoint_pb.PickPoint {
	return &pickpoint_pb.PickPoint{
		Id:      point.ID,
		Name:    point.Name,
		Address: point.Address,
		Contact: point.Contact,
	}
}

func (s service) Read(ctx context.Context, idRequest *pickpoint_pb.IdRequest) (*pickpoint_pb.PickPoint, error) {
	defer s.counter.Add(1)
	id := idRequest.Id
	if !isValidID(id) {
		return nil, model.ErrorInvalidInput
	}
	point := new(model.PickPoint)
	err := s.cache.Get(ctx, fmt.Sprint(id), point)
	if err == nil {
		return model2Pb(point), nil
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
	return model2Pb(pPoint), nil
}

func (s service) Create(ctx context.Context, point *pickpoint_pb.PickPoint) (*pickpoint_pb.PickPoint, error) {
	id, err := s.repo.Add(ctx, pb2Model(point))
	if err != nil {
		return nil, err
	}
	point.Id = id
	return point, nil
}

func (s service) Update(ctx context.Context, point *pickpoint_pb.PickPoint) (*emptypb.Empty, error) {
	modelPoint := pb2Model(point)
	if !isValidPickPoint(modelPoint) {
		return nil, model.ErrorInvalidInput
	}
	return nil, s.transactor.RunSerializable(ctx, pgx.ReadWrite, func(ctxTX context.Context) error {
		err := s.repo.Update(ctx, modelPoint)
		if err != nil {
			return err
		}
		return s.cache.Delete(ctx, fmt.Sprint(modelPoint.ID))
	})
}

func (s service) Delete(ctx context.Context, idRequest *pickpoint_pb.IdRequest) (*emptypb.Empty, error) {
	id := idRequest.Id
	if !isValidID(id) {
		return nil, model.ErrorInvalidInput
	}
	return nil, s.transactor.RunSerializable(ctx, pgx.ReadWrite, func(ctxTX context.Context) error {
		err := s.repo.Delete(ctx, id)
		if err != nil {
			return err
		}
		return s.cache.Delete(ctx, fmt.Sprint(id))
	})
}

/*
// Create writes pick-up points information in storage

	func (s service) Create(ctx context.Context, point *model.PickPoint) (*model.PickPoint, error) {
		id, err := s.repo.Add(ctx, point)
		if err != nil {
			return nil, err
		}
		point.ID = id
		return point, nil
	}

// Read gets pick-up points information from storage by id

	func (s service) Read(ctx context.Context, id int64) (*model.PickPoint, error) {
		if !isValidID(id) {
			return nil, model.ErrorInvalidInput
		}
		point := new(model.PickPoint)
		err := s.cache.Get(ctx, fmt.Sprint(id), point)
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

	func (s service) Update(ctx context.Context, point *model.PickPoint) error {
		if !isValidPickPoint(point) {
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

	func (s service) Delete(ctx context.Context, id int64) error {
		if !isValidID(id) {
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
*/
func isValidPickPoint(point *model.PickPoint) bool {
	return point.ID > 0
}

func isValidID(id int64) bool {
	return id > 0
}
