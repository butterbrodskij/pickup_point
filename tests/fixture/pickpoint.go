package fixture

import (
	"gitlab.ozon.dev/mer_marat/homework/internal/model"
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

func (b *PickpointBuilder) ValidInput1() *PickpointBuilder {
	return b.Name(ValidName1).Address(ValidAddress1).Contact(ValidContact1)
}

func (b *PickpointBuilder) Valid1() *PickpointBuilder {
	return b.ID(ValidID1).Name(ValidName1).Address(ValidAddress1).Contact(ValidContact1)
}

func (b *PickpointBuilder) InValid1() *PickpointBuilder {
	return b.ID(InValidID1).Name(ValidName1).Address(ValidAddress1).Contact(ValidContact1)
}
