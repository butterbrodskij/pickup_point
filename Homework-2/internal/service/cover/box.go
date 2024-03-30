package cover

import "gitlab.ozon.dev/mer_marat/homework/internal/model"

type box struct {
	order *model.Order
}

func newBox(order *model.Order) box {
	return box{order: order}
}

func (b box) validateOrder() error {
	if b.order.WeightGrams >= 30*model.GramsInKilo {
		return model.ErrorExcessWeight
	}
	return nil
}

func (b box) getPackagingPrice() int64 {
	return 20 * model.KopecksInRuble
}
