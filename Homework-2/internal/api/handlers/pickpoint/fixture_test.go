package handler

import (
	"testing"

	"github.com/golang/mock/gomock"
)

type pointHandlerFixture struct {
	ctrl     *gomock.Controller
	handl    *handler
	mockServ *Mockservice
}

func setUp(t *testing.T) pointHandlerFixture {
	ctrl := gomock.NewController(t)
	mockServ := NewMockservice(ctrl)
	handl := NewHandler(mockServ)
	return pointHandlerFixture{
		ctrl:     ctrl,
		handl:    handl,
		mockServ: mockServ,
	}
}

func (a *pointHandlerFixture) tearDown() {
	a.ctrl.Finish()
}
