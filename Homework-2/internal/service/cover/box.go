package cover

import "gitlab.ozon.dev/mer_marat/homework/internal/model"

type box struct {
	order model.Order
}

func newBox(order *model.Order) *box {
	return &box{order: *order}
}

func (b *box) OrderRequirements() bool {
	return true
}

func (b *box) OrderChanges() *model.Order {
	return &b.order
}
