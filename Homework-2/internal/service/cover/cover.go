package cover

import "gitlab.ozon.dev/mer_marat/homework/internal/model"

type cover interface {
	validateOrder() error
	getPackagingPrice() int64
}

type coveredOrder struct {
	cover
}

func NewCoveredOrder(order *model.Order) (coveredOrder, error) {
	switch order.Cover {
	case model.BagCover:
		return coveredOrder{newBag(order)}, nil
	case model.BoxCover:
		return coveredOrder{newBox(order)}, nil
	case model.FilmCover:
		return coveredOrder{newFilm(order)}, nil
	default:
		return coveredOrder{}, model.ErrorInvalidInput
	}
}

func (c coveredOrder) ValidateOrder() error {
	return c.validateOrder()
}

func (c coveredOrder) GetPackagingPrice() int64 {
	return c.getPackagingPrice()
}
