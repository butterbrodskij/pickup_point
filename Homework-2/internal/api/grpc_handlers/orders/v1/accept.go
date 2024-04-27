package grpc_order

import (
	"context"

	"gitlab.ozon.dev/mer_marat/homework/internal/model"
	order_pb "gitlab.ozon.dev/mer_marat/homework/internal/pkg/pb/homework/orders/v1"
)

func (s *grpcService) AcceptFromCourier(ctx context.Context, req *order_pb.AcceptFromCourierRequest) (*order_pb.AcceptFromCourierResponse, error) {
	err := s.service.AcceptFromCourier(ctx, model.OrderInput{
		ID:           req.Id,
		RecipientID:  req.RecipientId,
		WeightGrams:  req.WeightGrams,
		PriceKopecks: req.PriceKopecks,
		Cover:        req.Cover,
		ExpireDate:   req.ExpireDate,
	})
	return nil, err
}
