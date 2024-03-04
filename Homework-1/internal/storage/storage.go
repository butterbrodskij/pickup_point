package storage

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"homework1/pup/internal/model"
	"io"
	"os"
	"time"
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

	for _, ord := range all {
		if ord.ID == newOrder.ID {
			return errors.New("trying to get existing order")
		}
	}

	all = append(all, newOrder)

	return writeBytes(all)
}

func (s *Storage) Remove(id int) error {
	all, err := s.listAll()
	if err != nil {
		return err
	}

	for i, ord := range all {
		if ord.ID == id {
			if ord.ExpireDate.Before(time.Now()) && !ord.IsGiven {
				all = append(all[:i], all[i+1:]...)
				return writeBytes(all)
			} else {
				return errors.New("trying to remove order that is given or not expired")
			}
		}
	}

	return errors.New("order not found")
}

func (s *Storage) Give(ids []int) error {
	all, err := s.listAll()
	if err != nil {
		return err
	}
	var recipient int
	toModify := make([]int, len(ids))

	for i, id := range ids {
		idx, ok := searchId(all, id)
		if !ok {
			return fmt.Errorf("order %d is not in the storage", id)
		}
		if recipient != 0 && all[idx].Recipient != recipient {
			return errors.New("orders belong to different recipients")
		}
		if all[idx].IsGiven {
			return fmt.Errorf("order %d is already given", id)
		}
		if all[idx].ExpireDate.Before(time.Now()) {
			return fmt.Errorf("order %d is expired", id)
		}
		recipient = all[idx].Recipient
		toModify[i] = idx
	}

	for _, i := range toModify {
		all[i].IsGiven = true
		all[i].GivenTime = time.Now()
	}

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

func searchId(all []OrderDTO, id int) (int, bool) {
	for i, order := range all {
		if order.ID == id {
			return i, true
		}
	}
	return 0, false
}
