package grpc_order

import (
	"context"

	order_pb "gitlab.ozon.dev/mer_marat/homework/internal/pkg/pb/homework/orders/v1"
)

func (s *grpcService) Give(ctx context.Context, req *order_pb.GiveRequest) (*order_pb.GiveResponse, error) {
	err := s.service.Give(ctx, req.Ids)
	return nil, err
}
