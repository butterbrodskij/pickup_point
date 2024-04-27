package grpc_pickpoint

import (
	"context"

	pickpoint_pb "gitlab.ozon.dev/mer_marat/homework/internal/pkg/pb/homework/pickpoints/v1"
)

func (s *grpcService) Delete(ctx context.Context, deleteRequest *pickpoint_pb.DeleteRequest) (*pickpoint_pb.DeleteResponse, error) {
	err := s.service.Delete(ctx, deleteRequest.Id)
	return nil, err
}
