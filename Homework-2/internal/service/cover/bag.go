package cover

import "gitlab.ozon.dev/mer_marat/homework/internal/model"

type bag struct {
	order *model.Order
}

func newBag(order *model.Order) bag {
	return bag{order: order}
}

func (b bag) validateOrder() error {
	if b.order.WeightGrams >= 10*model.GramsInKilo {
		return model.ErrorExcessWeight
	}
	return nil
}

func (b bag) getPackagingPrice() int64 {
	return b.order.PriceKopecks + 5*model.KopecksInRuble
}
