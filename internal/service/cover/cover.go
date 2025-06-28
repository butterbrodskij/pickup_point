package cover

import "gitlab.ozon.dev/mer_marat/homework/internal/model"

type service struct{}

type cover interface {
	validateOrder() error
	getPackagingPrice() int64
}

func NewService() service {
	return service{}
}

func (s service) ValidateOrder(order model.Order) error {
	cov := getCover(&order)
	if cov == nil {
		return model.ErrorInvalidInput
	}
	return cov.validateOrder()
}

func (s service) GetPackagingPrice(order model.Order) int64 {
	cov := getCover(&order)
	if cov == nil {
		return order.PriceKopecks
	}
	return cov.getPackagingPrice()
}

func getCover(order *model.Order) cover {
	switch order.Cover {
	case model.BagCover:
		return newBag(order)
	case model.BoxCover:
		return newBox(order)
	case model.FilmCover:
		return newFilm(order)
	default:
		return nil
	}
}
