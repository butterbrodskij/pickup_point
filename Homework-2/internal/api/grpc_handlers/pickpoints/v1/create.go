package grpc_pickpoint

import (
	"context"

	pickpoint_pb "gitlab.ozon.dev/mer_marat/homework/internal/pkg/pb/homework/pickpoints/v1"
)

func (s *grpcService) Create(ctx context.Context, createRequest *pickpoint_pb.CreateRequest) (*pickpoint_pb.CreateResponse, error) {
	point := pb2Model(createRequest.Point)
	res, err := s.service.Create(ctx, point)
	if err != nil {
		return nil, err
	}
	pb := model2Pb(res)
	return &pickpoint_pb.CreateResponse{Point: pb}, nil
}
