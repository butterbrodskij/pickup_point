package service

type storage interface {
}

type Service struct {
	s storage
}

func New(stor storage) Service {
	return Service{s: stor}
}
