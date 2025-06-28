package grpc_pickpoint

import (
	"context"

	"gitlab.ozon.dev/mer_marat/homework/internal/model"
	pickpoint_pb "gitlab.ozon.dev/mer_marat/homework/internal/pkg/pb/homework/pickpoints/v1"
	"google.golang.org/grpc"
)

func pb2Model(point *pickpoint_pb.PickPoint) *model.PickPoint {
	return &model.PickPoint{
		ID:      point.Id,
		Name:    point.Name,
		Address: point.Address,
		Contact: point.Contact,
	}
}

func model2Pb(point *model.PickPoint) *pickpoint_pb.PickPoint {
	return &pickpoint_pb.PickPoint{
		Id:      point.ID,
		Name:    point.Name,
		Address: point.Address,
		Contact: point.Contact,
	}
}

type service interface {
	Read(ctx context.Context, id int64) (res *model.PickPoint, err error)
	Create(ctx context.Context, point *model.PickPoint) (res *model.PickPoint, err error)
	Update(ctx context.Context, point *model.PickPoint) (err error)
	Delete(ctx context.Context, id int64) (err error)
}

type grpcService struct {
	pickpoint_pb.UnimplementedPickPointsServer
	service service
}

func NewGRPCPickpointService(s service) *grpcService {
	return &grpcService{service: s}
}

func (s *grpcService) RegisterGRPC(server grpc.ServiceRegistrar) {
	pickpoint_pb.RegisterPickPointsServer(server, s)
}
