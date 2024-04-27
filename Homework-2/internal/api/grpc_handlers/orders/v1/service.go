package grpc_order

import (
	"context"

	"gitlab.ozon.dev/mer_marat/homework/internal/model"
	order_pb "gitlab.ozon.dev/mer_marat/homework/internal/pkg/pb/homework/orders/v1"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type service interface {
	AcceptFromCourier(ctx context.Context, input model.OrderInput) error
	Remove(ctx context.Context, id int64) error
	Give(ctx context.Context, ids []int64) error
	List(ctx context.Context, recipient int64, n int, onlyNotGivenOrders bool) ([]model.Order, error)
	Return(ctx context.Context, id, recipient int64) error
	ListReturn(ctx context.Context, pageNum, ordersPerPage int) ([]model.Order, error)
}

type grpcService struct {
	order_pb.UnimplementedOrdersServer
	service service
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

func NewGRPCOrderService(s service) *grpcService {
	return &grpcService{service: s}
}

func (s *grpcService) RegisterGRPC(server grpc.ServiceRegistrar) {
	order_pb.RegisterOrdersServer(server, s)
}
