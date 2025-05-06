package sales

import "errors"

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
		return nil, errors.New("Sale not found")
	}
	return sale, nil
}

func (st *SaleStorage) PutSale(sale *Sale) error {
	if sale.Id == "" {
		errors.New("Sale ID is empty")
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
	if len(sales) == 0 {
		return nil, errors.New("No sales found")
	}
	return &sales, nil
}
