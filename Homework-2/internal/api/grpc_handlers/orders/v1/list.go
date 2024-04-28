package grpc_order

import (
	"context"

	order_pb "gitlab.ozon.dev/mer_marat/homework/internal/pkg/pb/homework/orders/v1"
)

func (s *grpcService) List(ctx context.Context, req *order_pb.ListRequest) (*order_pb.ListResponse, error) {
	res, err := s.service.List(ctx, req.Recipient, int(req.N), req.OnlyNotGivenOrders)
	if err != nil {
		return nil, err
	}
	return &order_pb.ListResponse{Orders: &order_pb.OrderList{Orders: order2PbSlice(res)}}, nil
}
