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
	mockTrans *Mocktransactor
}

func setUp(t *testing.T) pointServiceFixture {
	ctrl := gomock.NewController(t)
	mockRepo := NewMockstorage(ctrl)
	mockCache := NewMockcache(ctrl)
	mockTrans := NewMocktransactor(ctrl)
	serv := NewService(mockRepo, mockCache, mockTrans)
	return pointServiceFixture{
		ctrl:      ctrl,
		serv:      serv,
		mockCache: mockCache,
		mockRepo:  mockRepo,
		mockTrans: mockTrans,
	}
}

func (a *pointServiceFixture) tearDown() {
	a.ctrl.Finish()
}
