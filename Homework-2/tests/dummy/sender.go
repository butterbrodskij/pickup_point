package dummy

import "gitlab.ozon.dev/mer_marat/homework/internal/model"

type dummySender struct{}

func NewDummySender() dummySender {
	return dummySender{}
}

func (dummySender) SendMessage(message model.RequestMessage) error {
	return nil
}
