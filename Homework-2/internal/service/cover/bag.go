package cover

import "gitlab.ozon.dev/mer_marat/homework/internal/model"

type bag struct {
	order model.Order
}

func newBag(order *model.Order) *bag {
	return &bag{order: *order}
}

func (b *bag) OrderRequirements() bool {
	return true
}

func (b *bag) OrderChanges() *model.Order {
	return &b.order
}
