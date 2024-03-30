package cover

import "gitlab.ozon.dev/mer_marat/homework/internal/model"

type film struct {
	order *model.Order
}

func newFilm(order *model.Order) film {
	return film{order: order}
}

func (b film) validateOrder() error {
	return nil
}

func (b film) getPackagingPrice() int64 {
	return 1 * model.KopecksInRuble
}
