//go:generate mockgen -source=./handler.go -destination=./handler_mocks_test.go -package=handler
package handler

import (
	"context"

	pickpoint_pb "gitlab.ozon.dev/mer_marat/homework/internal/pkg/pb/homework/pickpoints/v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

type service interface {
	Create(ctx context.Context, point *pickpoint_pb.PickPoint) (*pickpoint_pb.PickPoint, error)
	Read(ctx context.Context, idRequest *pickpoint_pb.IdRequest) (*pickpoint_pb.PickPoint, error)
	Update(ctx context.Context, point *pickpoint_pb.PickPoint) (*emptypb.Empty, error)
	Delete(ctx context.Context, idRequest *pickpoint_pb.IdRequest) (*emptypb.Empty, error)
}

type handler struct {
	service
}

func NewHandler(s service) *handler {
	return &handler{
		service: s,
	}
}
