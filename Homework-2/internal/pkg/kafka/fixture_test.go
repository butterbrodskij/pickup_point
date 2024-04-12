package kafka

import (
	"testing"

	"github.com/golang/mock/gomock"
)

type senderFixture struct {
	ctrl     *gomock.Controller
	sender   *KafkaSender
	mockProd *Mockproducer
}

func setUp(t *testing.T, topic string) senderFixture {
	ctrl := gomock.NewController(t)
	mockProd := NewMockproducer(ctrl)
	sender := NewKafkaSender(mockProd, topic)
	return senderFixture{
		ctrl:     ctrl,
		sender:   sender,
		mockProd: mockProd,
	}
}

func (a *senderFixture) tearDown() {
	a.ctrl.Finish()
}
