package storage

import (
	"homework1/pup/internal/model"
	"os"
)

const storageName = "storage"

type Storage struct {
	storage *os.File
}

// New returns new storage associated with file storageName
func New() (Storage, error) {
	file, err := os.OpenFile(storageName, os.O_CREATE, 0777)
	if err != nil {
		return Storage{}, err
	}
	return Storage{storage: file}, nil
}

func Into2DTO(order model.OrderInput) OrderDTO {
	return OrderDTO{
		ID:         order.ID,
		Recipient:  order.Recipient,
		ExpireDate: order.ExpireDate,
	}
}
