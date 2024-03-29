package cover

import "gitlab.ozon.dev/mer_marat/homework/internal/model"

type Cover interface {
	OrderRequirements() bool
	OrderChanges() *model.Order
}

func CoveredOrder(order *model.Order) (Cover, error) {
	switch order.Cover {
	case model.BagCover:
		return newBag(order), nil
	case model.BoxCover:
		return newBox(order), nil
	case model.FilmCover:
		return newFilm(order), nil
	default:
		return nil, model.ErrorInvalidInput
	}
}
