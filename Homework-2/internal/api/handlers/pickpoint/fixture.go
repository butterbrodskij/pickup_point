package handler

import (
	"testing"

	"github.com/golang/mock/gomock"
	mock_handler "gitlab.ozon.dev/mer_marat/homework/internal/api/handlers/pickpoint/mocks"
)

type pointHandlerFixture struct {
	ctrl     *gomock.Controller
	handl    handler
	mockServ *mock_handler.Mockservice
}

func setUp(t *testing.T) pointHandlerFixture {
	ctrl := gomock.NewController(t)
	mockServ := mock_handler.NewMockservice(ctrl)
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
