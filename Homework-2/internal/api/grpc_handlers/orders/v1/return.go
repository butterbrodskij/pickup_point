package grpc_order

import (
	"context"

	order_pb "gitlab.ozon.dev/mer_marat/homework/internal/pkg/pb/homework/orders/v1"
)

func (s *grpcService) Return(ctx context.Context, req *order_pb.ReturnRequest) (*order_pb.ReturnResponse, error) {
	err := s.service.Return(ctx, req.Id, req.Recipient)
	if err != nil {
		return nil, err
	}
	return &order_pb.ReturnResponse{}, nil
}
