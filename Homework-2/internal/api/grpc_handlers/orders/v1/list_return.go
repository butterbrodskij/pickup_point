package grpc_order

import (
	"context"

	order_pb "gitlab.ozon.dev/mer_marat/homework/internal/pkg/pb/homework/orders/v1"
)

func (s *grpcService) ListReturn(ctx context.Context, req *order_pb.ListReturnRequest) (*order_pb.ListReturnResponse, error) {
	res, err := s.service.ListReturn(ctx, int(req.PageNum), int(req.OrdersPerPage))
	return &order_pb.ListReturnResponse{Orders: &order_pb.OrderList{Orders: order2PbSlice(res)}}, err
}
