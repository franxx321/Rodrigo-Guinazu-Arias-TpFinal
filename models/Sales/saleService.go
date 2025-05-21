package Sales

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"math/rand/v2"
	"time"
)

var ErrSaleNotPending = errors.New("Sale is not Pending")
var ErrInvalidTransition = errors.New("Invalid Transition")
var ErrInvalidStatus = errors.New("Invalid Status")

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
	fmt.Printf("holaa")
	rand := rand.IntN(3)
	if rand == 0 {
		status = Pending
	} else if rand == 1 {
		status = Aproved
	} else if rand == 2 {
		status = Rejected
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

func (s *SaleService) Update(id string, status string) error {
	sale, err := s.saleStorage.GetSale(id)
	if err != nil {
		return err
	}
	if status != Pending && status != Aproved && status != Rejected {
		return ErrInvalidStatus
	}
	if status == Pending {
		return ErrInvalidTransition
	}
	if sale.Status != Pending {
		return ErrSaleNotPending
	}
	sale.Status = status
	sale.UpdatedAt = time.Now()
	sale.Version++
	s.saleStorage.PutSale(sale)
	return nil
}

func (s *SaleService) GetByUserStatus(userId string, status string) (*[]Sale, error) {

	if status != "" && status != Pending && status != Aproved && status != Rejected {
		return nil, ErrInvalidStatus
	}
	sales, err := s.saleStorage.GetByUserStatus(userId, status)
	if err != nil {
		return nil, err
	}
	return sales, nil
}
