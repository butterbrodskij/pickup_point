package pickpoint

import (
	"testing"

	"github.com/golang/mock/gomock"
)

type pointServiceFixture struct {
	ctrl     *gomock.Controller
	serv     service
	mockRepo *Mockstorage
}

func setUp(t *testing.T) pointServiceFixture {
	ctrl := gomock.NewController(t)
	mockRepo := NewMockstorage(ctrl)
	serv := NewService(mockRepo)
	return pointServiceFixture{
		ctrl:     ctrl,
		serv:     serv,
		mockRepo: mockRepo,
	}
}

func (a *pointServiceFixture) tearDown() {
	a.ctrl.Finish()
}
