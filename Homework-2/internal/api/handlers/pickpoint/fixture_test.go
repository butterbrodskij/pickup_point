package handler

import (
	"testing"

	"github.com/golang/mock/gomock"
)

type pointHandlerFixture struct {
	ctrl      *gomock.Controller
	handl     *handler
	mockServ  *Mockservice
	mockCache *Mockcache
}

func setUp(t *testing.T) pointHandlerFixture {
	ctrl := gomock.NewController(t)
	mockServ := NewMockservice(ctrl)
	mockCache := NewMockcache(ctrl)
	handl := NewHandler(mockServ, mockCache)
	return pointHandlerFixture{
		ctrl:      ctrl,
		handl:     handl,
		mockServ:  mockServ,
		mockCache: mockCache,
	}
}

func (a *pointHandlerFixture) tearDown() {
	a.ctrl.Finish()
}
