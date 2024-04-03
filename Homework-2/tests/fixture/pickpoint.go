package fixture

import (
	"gitlab.ozon.dev/mer_marat/homework/internal/model"
	"gitlab.ozon.dev/mer_marat/homework/tests/states"
)

type PickpointBuilder struct {
	point *model.PickPoint
}

func PickPoint() *PickpointBuilder {
	return &PickpointBuilder{point: &model.PickPoint{}}
}

func (b *PickpointBuilder) ID(id int64) *PickpointBuilder {
	b.point.ID = id
	return b
}

func (b *PickpointBuilder) Name(name string) *PickpointBuilder {
	b.point.Name = name
	return b
}

func (b *PickpointBuilder) Address(addr string) *PickpointBuilder {
	b.point.Address = addr
	return b
}

func (b *PickpointBuilder) Contact(cont string) *PickpointBuilder {
	b.point.Contact = cont
	return b
}

func (b *PickpointBuilder) P() *model.PickPoint {
	return b.point
}

func (b *PickpointBuilder) V() model.PickPoint {
	return *b.point
}

func (b *PickpointBuilder) Valid1() *PickpointBuilder {
	return b.ID(states.ValidID1).Name(states.ValidName1).Address(states.ValidAddress1).Contact(states.ValidContact1)
}
