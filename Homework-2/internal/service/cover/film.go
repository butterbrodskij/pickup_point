package cover

import "gitlab.ozon.dev/mer_marat/homework/internal/model"

type film struct {
	order model.Order
}

func newFilm(order *model.Order) *film {
	return &film{order: *order}
}

func (b *film) OrderRequirements() bool {
	return true
}

func (b *film) OrderChanges() *model.Order {
	b.order.PriceKopecks += 1 * model.KopecksInRuble
	return &b.order
}
