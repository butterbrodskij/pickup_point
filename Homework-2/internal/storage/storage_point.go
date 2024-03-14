package storage

import (
	"encoding/json"
	"homework2/pup/internal/model"
	"os"
	"sync"
)

type StoragePoints struct {
	storageName string
	content     []model.PickPoint
	mt          sync.Mutex
}

// New returns new storage associated with file storageName
func NewPoints(storageName string) (StoragePoints, error) {
	_, err := os.Stat(storageName)
	if os.IsNotExist(err) {
		f, err := os.Create(storageName)
		if err != nil {
			return StoragePoints{}, err
		}
		err = f.Close()
		return StoragePoints{}, err
	} else if err != nil {
		return StoragePoints{}, err
	}
	content, err := listAllPoints(storageName)
	if err != nil {
		return StoragePoints{}, err
	}
	return StoragePoints{
		storageName: storageName,
		content:     content,
	}, nil
}

// writeBytes writes orders in file in json
func (s *StoragePoints) writeBytes(orders []model.PickPoint) error {
	s.content = orders
	rawBytes, err := json.Marshal(orders)
	if err != nil {
		return err
	}

	return os.WriteFile(s.storageName, rawBytes, 0777)
}

// listAll returns all orders in storage
func listAllPoints(storageName string) ([]model.PickPoint, error) {
	rawBytes, err := os.ReadFile(storageName)
	if err != nil {
		return nil, err
	}

	orders := make([]model.PickPoint, 0)
	if len(rawBytes) == 0 {
		return orders, nil
	}

	err = json.Unmarshal(rawBytes, &orders)
	if err != nil {
		return nil, err
	}

	return orders, nil
}
