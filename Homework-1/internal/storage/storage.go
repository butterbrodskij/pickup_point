package storage

import (
	"bufio"
	"encoding/json"
	"homework1/pup/internal/model"
	"io"
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

func Order2DTO(order model.Order) OrderDTO {
	return OrderDTO{
		ID:         order.ID,
		Recipient:  order.Recipient,
		ExpireDate: order.ExpireDate,
	}
}

func (s *Storage) Get(order model.Order) error {
	all, err := s.listAll()
	if err != nil {
		return err
	}

	newOrder := Order2DTO(order)
	all = append(all, newOrder)

	return writeBytes(all)
}

func writeBytes(toDos []OrderDTO) error {
	rawBytes, err := json.Marshal(toDos)
	if err != nil {
		return err
	}

	return os.WriteFile(storageName, rawBytes, 0777)
}

func (s *Storage) listAll() ([]OrderDTO, error) {
	reader := bufio.NewReader(s.storage)
	rawBytes, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	orders := make([]OrderDTO, 0)
	if len(rawBytes) == 0 {
		return orders, nil
	}

	err = json.Unmarshal(rawBytes, &orders)
	if err != nil {
		return nil, err
	}

	return orders, nil
}
