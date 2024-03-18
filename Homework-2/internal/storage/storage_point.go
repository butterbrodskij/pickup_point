package storage

import (
	"encoding/json"
	"errors"
	"homework2/pup/internal/model"
	"os"
	"sync"
)

type StoragePoints struct {
	storageName string
	content     []model.PickPoint
	mt          sync.RWMutex
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

// Write adds new pick-up point to storage
func (s *StoragePoints) Write(point model.PickPoint) error {
	s.mt.Lock()
	defer s.mt.Unlock()
	all := s.content
	for _, p := range all {
		if p.ID == point.ID {
			return errors.New("can not write new pick-up point: trying to add existing point")
		}
	}
	all = append(all, point)
	return s.writeBytes(all)
}

// Get returns pick-up point by its id
func (s *StoragePoints) Get(id int64) (model.PickPoint, bool) {
	s.mt.RLock()
	defer s.mt.RUnlock()
	all := s.content
	for _, point := range all {
		if point.ID == id {
			return point, true
		}
	}
	return model.PickPoint{}, false
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
