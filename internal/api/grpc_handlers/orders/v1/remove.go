package grpc_order

import (
	"context"

	order_pb "gitlab.ozon.dev/mer_marat/homework/internal/pkg/pb/homework/orders/v1"
)

func (s *grpcService) Remove(ctx context.Context, req *order_pb.RemoveRequest) (*order_pb.RemoveResponse, error) {
	err := s.service.Remove(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &order_pb.RemoveResponse{}, nil
}
