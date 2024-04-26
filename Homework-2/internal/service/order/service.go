package order

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"gitlab.ozon.dev/mer_marat/homework/internal/model"
	order_pb "gitlab.ozon.dev/mer_marat/homework/internal/pkg/pb/order"
	storage "gitlab.ozon.dev/mer_marat/homework/internal/storage/file"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type coverService interface {
	ValidateOrder(model.Order) error
	GetPackagingPrice(model.Order) int64
}

type storageInterface interface {
	AcceptFromCourier(model.Order) error
	Remove(int64) error
	Give([]int64) error
	ListAll(int64) ([]model.Order, error)
	ListNotGiven(int64) ([]model.Order, error)
	Return(int64) error
	ListReturn() ([]model.Order, error)
	GetByID(int64) (storage.OrderDTO, bool)
}

type service struct {
	order_pb.UnimplementedOrdersServer
	s                  storageInterface
	cov                coverService
	givenOrdersCounter prometheus.Gauge
}

func pb2Order(input *order_pb.OrderInput) (model.Order, error) {
	if input.Id <= 0 {
		return model.Order{}, errors.New("id should be positive")
	}
	if input.RecipientId <= 0 {
		return model.Order{}, errors.New("recipient id should be positive")
	}
	if input.WeightGrams <= 0 {
		return model.Order{}, errors.New("order weight should be positive")
	}
	if input.PriceKopecks <= 0 {
		return model.Order{}, errors.New("order price should be positive")
	}

	t, err := time.Parse("2.1.2006", input.ExpireDate)
	if err != nil {
		return model.Order{}, errors.New("wrong date format")
	}

	return model.Order{
		ID:           input.Id,
		RecipientID:  input.RecipientId,
		WeightGrams:  input.WeightGrams,
		PriceKopecks: input.PriceKopecks,
		Cover:        input.Cover,
		ExpireDate:   t,
	}, nil
}

func order2PbSlice(input []model.Order) []*order_pb.Order {
	res := make([]*order_pb.Order, len(input))
	for i, order := range input {
		res[i] = &order_pb.Order{
			Id:           order.ID,
			RecipientId:  order.RecipientID,
			WeightGrams:  order.WeightGrams,
			PriceKopecks: order.PriceKopecks,
			Cover:        order.Cover,
			ExpireDate:   timestamppb.New(order.ExpireDate),
		}
	}
	return res
}

// New returns type Service associated with storage
func NewService(stor storageInterface, cov coverService) service {
	return service{s: stor, cov: cov}
}

func (s *service) AddGivenOrdersMetrics(metrics prometheus.Gauge) {
	s.givenOrdersCounter = metrics
}

// Get checks validity of given data and adds new order to storage
func (s service) AcceptFromCourier(ctx context.Context, in *order_pb.OrderInput) (*emptypb.Empty, error) {
	order, err := pb2Order(in)
	if err != nil {
		return nil, err
	}
	if order.ExpireDate.Before(time.Now()) {
		return nil, errors.New("can not get order: trying to get expired order")
	}
	if err = s.cov.ValidateOrder(order); err != nil {
		return nil, err
	}
	order.PriceKopecks += s.cov.GetPackagingPrice(order)
	return nil, s.s.AcceptFromCourier(order)
}

// Remove checks validity of given id and deletes an order from storage
func (s service) Remove(ctx context.Context, idRequest *order_pb.IdRequest) (*emptypb.Empty, error) {
	id := idRequest.Id
	if id <= 0 {
		return nil, errors.New("id should be positive")
	}
	if order, ok := s.s.GetByID(id); ok && order.ExpireDate.After(time.Now()) || order.IsGiven {
		return nil, errors.New("order can not be removed: trying to remove order that is given or not expired")
	}
	return nil, s.s.Remove(id)
}

// Give checks validity of given ids and gives orders to recipient
func (s service) Give(ctx context.Context, idsRequest *order_pb.Ids) (*emptypb.Empty, error) {
	var recipient int64
	ids := idsRequest.Ids

	for _, id := range ids {
		order, ok := s.s.GetByID(id)
		switch {
		case !ok:
			return nil, fmt.Errorf("can not give orders: order %d is not in the storage", id)
		case recipient != 0 && order.RecipientID != recipient:
			return nil, errors.New("can not give orders: orders belong to different recipients")
		case order.IsGiven:
			return nil, fmt.Errorf("can not give orders: order %d is already given", id)
		case order.IsReturned:
			return nil, fmt.Errorf("can not give orders: order %d is already returned by recipient", id)
		case order.ExpireDate.Before(time.Now()):
			return nil, fmt.Errorf("can not give orders: order %d is expired", id)
		case recipient == 0:
			recipient = order.RecipientID
		}
	}
	err := s.s.Give(ids)
	if err != nil {
		return nil, err
	}
	if s.givenOrdersCounter != nil {
		s.givenOrdersCounter.Add(float64(len(ids)))
	}

	return nil, nil
}

// List checks validity of given recipient id and n and returns slice of all his orders (last n)
func (s service) List(ctx context.Context, req *order_pb.ListRequest) (*order_pb.OrderList, error) {
	recipient := req.Recipient
	n := int(req.N)
	if recipient <= 0 {
		return nil, errors.New("recipient id should be positive")
	}
	if n < 0 {
		return nil, errors.New("n should not be negative")
	}
	var (
		all []model.Order
		err error
	)
	if !req.OnlyNotGivenOrders {
		all, err = s.s.ListAll(recipient)
	} else {
		all, err = s.s.ListNotGiven(recipient)
	}
	if err != nil {
		return nil, err
	}
	if n == 0 || len(all) <= n {
		return &order_pb.OrderList{Orders: order2PbSlice(all)}, nil
	}

	return &order_pb.OrderList{Orders: order2PbSlice(all[len(all)-n:])}, err
}

// Return checks validity of given order id and recipient id and gets order back from recipient
func (s service) Return(ctx context.Context, returnRequest *order_pb.ReturnRequest) (*emptypb.Empty, error) {
	id := returnRequest.Id
	recipient := returnRequest.Recipient
	if id <= 0 {
		return nil, errors.New("id should be positive")
	}
	if recipient <= 0 {
		return nil, errors.New("recipient id should be positive")
	}
	order, ok := s.s.GetByID(id)
	switch {
	case !ok:
		return nil, errors.New("order can not be returned: order not found")
	case order.RecipientID != recipient:
		return nil, errors.New("order can not be returned: order belongs to different recipient")
	case order.IsReturned:
		return nil, errors.New("order can not be returned: order is already returned")
	case !order.IsGiven:
		return nil, errors.New("order can not be returned: order is not given yet")
	case order.GivenTime.AddDate(0, 0, 2).Before(time.Now()):
		return nil, errors.New("order can not be returned: more than 2 days passed")
	}
	err := s.s.Return(id)
	if err != nil {
		return nil, err
	}
	if s.givenOrdersCounter != nil {
		s.givenOrdersCounter.Sub(1)
	}

	return nil, nil
}

// ListReturn checks validity of given args and returns k returned orders on nth page
func (s service) ListReturn(ctx context.Context, request *order_pb.ListReturnRequest) (*order_pb.OrderList, error) {
	pageNum := int(request.PageNum)
	ordersPerPage := int(request.OrdersPerPage)
	if pageNum < 0 {
		return nil, errors.New("pageNum should not be negative")
	}
	if ordersPerPage <= 0 {
		return nil, errors.New("ordersPerPage should be positive")
	}
	all, err := s.s.ListReturn()
	if err != nil || pageNum == 0 {
		return &order_pb.OrderList{Orders: order2PbSlice(all)}, err
	}
	firstPos := (pageNum - 1) * ordersPerPage
	if len(all) == 0 || len(all) <= firstPos {
		return &order_pb.OrderList{Orders: order2PbSlice(all)}, errors.New("empty list")
	}
	newLen := ordersPerPage
	if len(all) < pageNum*ordersPerPage {
		newLen = len(all) % ordersPerPage
	}
	return &order_pb.OrderList{Orders: order2PbSlice(all[firstPos : firstPos+newLen])}, nil
}
