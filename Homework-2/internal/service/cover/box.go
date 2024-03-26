package cover

import "gitlab.ozon.dev/mer_marat/homework/internal/model"

type box struct {
	order model.Order
}

func newBox(order *model.Order) *box {
	return &box{order: *order}
}

func (b *box) OrderRequirements() bool {
	return b.order.Weight < 30
}

func (b *box) OrderChanges() *model.Order {
	b.order.Price += 20
	return &b.order
}
