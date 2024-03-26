package cover

import "gitlab.ozon.dev/mer_marat/homework/internal/model"

type Cover interface {
	OrderRequirements() bool
	OrderChanges() *model.Order
}

func CoveredOrder(order *model.Order) (Cover, error) {
	switch order.Cover {
	case "bag":
		return newBag(order), nil
	case "box":
		return newBox(order), nil
	case "film":
		return newFilm(order), nil
	default:
		return nil, model.ErrorInvalidInput
	}
}
