package pickpoint

import (
	"testing"

	"github.com/golang/mock/gomock"
)

type pointServiceFixture struct {
	ctrl      *gomock.Controller
	serv      service
	mockCache *Mockcache
	mockRepo  *Mockstorage
}

func setUp(t *testing.T) pointServiceFixture {
	ctrl := gomock.NewController(t)
	mockRepo := NewMockstorage(ctrl)
	mockCache := NewMockcache(ctrl)
	serv := NewService(mockRepo, mockCache)
	return pointServiceFixture{
		ctrl:      ctrl,
		serv:      serv,
		mockCache: mockCache,
		mockRepo:  mockRepo,
	}
}

func (a *pointServiceFixture) tearDown() {
	a.ctrl.Finish()
}
