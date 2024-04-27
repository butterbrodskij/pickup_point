package grpc_pickpoint

import (
	"context"

	pickpoint_pb "gitlab.ozon.dev/mer_marat/homework/internal/pkg/pb/homework/pickpoints/v1"
)

func (s *grpcService) Read(ctx context.Context, readRequest *pickpoint_pb.ReadRequest) (*pickpoint_pb.ReadResponse, error) {
	id := readRequest.Id
	res, err := s.service.Read(ctx, id)
	if err != nil {
		return nil, err
	}
	pb := model2Pb(res)
	return &pickpoint_pb.ReadResponse{Point: pb}, nil
}
