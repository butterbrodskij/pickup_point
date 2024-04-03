package postgres

import (
	"testing"

	"github.com/golang/mock/gomock"
	mock_repo "gitlab.ozon.dev/mer_marat/homework/internal/storage/postgres/mocks"
)

type pointRepoFixture struct {
	ctrl   *gomock.Controller
	repo   *PickpointRepo
	mockDB *mock_repo.Mockdatabase
}

func setUp(t *testing.T) pointRepoFixture {
	ctrl := gomock.NewController(t)
	mockDB := mock_repo.NewMockdatabase(ctrl)
	repo := NewRepo(mockDB)
	return pointRepoFixture{
		ctrl:   ctrl,
		repo:   repo,
		mockDB: mockDB,
	}
}

func (a *pointRepoFixture) tearDown() {
	a.ctrl.Finish()
}
