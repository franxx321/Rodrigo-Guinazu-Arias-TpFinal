package Sales

import "errors"

var ErrNotFound = errors.New("sale not found")
var ErrEmptyID = errors.New("empty sale ID")

type SaleStorage struct {
	m map[string]*Sale
}

func NewSaleStorage() *SaleStorage {
	return &SaleStorage{
		m: map[string]*Sale{},
	}
}

func (st *SaleStorage) GetSale(saleID string) (*Sale, error) {
	sale, ok := st.m[saleID]
	if !ok {
		return nil, ErrNotFound
	}
	return sale, nil
}

func (st *SaleStorage) PutSale(sale *Sale) error {
	if sale.Id == "" {
		return ErrEmptyID
	}
	st.m[sale.Id] = sale
	return nil
}

func (st *SaleStorage) GetByUserStatus(userId string, status string) (*[]Sale, error) {
	sales := []Sale{}
	if status == "" {
		for _, sale := range st.m {
			if sale.UserId == userId {
				sales = append(sales, *sale)
			}
		}
	} else {
		for _, sale := range st.m {
			if sale.UserId == userId && sale.Status == status {
				sales = append(sales, *sale)
			}
		}
	}
	return &sales, nil
}
