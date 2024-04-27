package grpc_pickpoint

import (
	"context"

	pickpoint_pb "gitlab.ozon.dev/mer_marat/homework/internal/pkg/pb/homework/pickpoints/v1"
)

func (s *grpcService) Update(ctx context.Context, updateRequest *pickpoint_pb.UpdateRequest) (*pickpoint_pb.UpdateResponse, error) {
	point := pb2Model(updateRequest.Point)
	err := s.service.Update(ctx, point)
	return nil, err
}
