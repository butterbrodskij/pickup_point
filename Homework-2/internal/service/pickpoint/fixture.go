package pickpoint

import (
	"testing"

	"github.com/golang/mock/gomock"
	mock_service "gitlab.ozon.dev/mer_marat/homework/internal/service/pickpoint/mocks"
)

type pointServiceFixture struct {
	ctrl     *gomock.Controller
	serv     service
	mockRepo *mock_service.Mockstorage
}

func setUp(t *testing.T) pointServiceFixture {
	ctrl := gomock.NewController(t)
	mockRepo := mock_service.NewMockstorage(ctrl)
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
