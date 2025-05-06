package sales

import (
	"errors"
	"github.com/google/uuid"
	"math/rand/v2"
	"time"
)

type SaleService struct {
	saleStorage *SaleStorage
}

func NewSaleService(saleStorage *SaleStorage) *SaleService {
	return &SaleService{
		saleStorage: saleStorage,
	}
}

func (s *SaleService) Create(userId string, amount float32) (*Sale, error) {

	//TODO: chequear que el usuario existe(falta que el joako termine el user Service)
	var status string
	rand := rand.IntN(3)
	if rand == 0 {
		status = pending
	} else if rand == 1 {
		status = aproved
	} else if rand == 2 {
		status = rejected
	}

	sale := &Sale{
		Id:        uuid.NewString(),
		UserId:    userId,
		Amount:    amount,
		Status:    status,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Version:   1,
	}

	err := s.saleStorage.PutSale(sale)
	if err != nil {
		return nil, err
	}

	return sale, nil
}

func (s *SaleService) update(id string, status string) error {
	sale, err := s.saleStorage.GetSale(id)
	if err != nil {
		return err
	}
	if sale.Status != pending {
		return errors.New("Sale is not pending")
	}
	if status != aproved && status != rejected {
		return errors.New("Invalid new status")
	}
	sale.Status = status
	sale.UpdatedAt = time.Now()
	sale.Version++
	s.saleStorage.PutSale(sale)
	return nil
}

func (s *SaleService) GetByUserStatus(userId string, status string) (*[]Sale, error) {

	if status != "" && status != pending && status != aproved && status != rejected {
		return nil, errors.New("Invalid status")
	}
	sales, err := s.saleStorage.GetByUserStatus(userId, status)
	if err != nil {
		return nil, err
	}
	return sales, nil
}
