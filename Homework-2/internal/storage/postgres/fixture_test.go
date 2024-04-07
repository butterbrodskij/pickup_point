package postgres

import (
	"testing"

	"github.com/golang/mock/gomock"
)

type pointRepoFixture struct {
	ctrl   *gomock.Controller
	repo   *PickpointRepo
	mockDB *Mockdatabase
}

func setUp(t *testing.T) pointRepoFixture {
	ctrl := gomock.NewController(t)
	mockDB := NewMockdatabase(ctrl)
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
