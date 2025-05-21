package Sales

import (
	"Rodrigo-Guinazu-Arias-TpFinal/models/users"
	"github.com/stretchr/testify/require"
	"testing"
)

type MockSaleService struct {
	saleStorage *SaleStorage
	mockCreate  func(userId string, amount float32) (*Sale, error)
}

func (m *MockSaleService) Create(userId string, amount float32) (*Sale, error) {
	if m.mockCreate != nil {
		return m.mockCreate(userId, amount)
	}
	return nil, nil
}

func TestSaleService_Create(t *testing.T) {

	type testCase struct {
		name     string
		mockFn   func(userId string, amount float32) (*Sale, error)
		userId   string
		amount   float32
		wantErr  func(t *testing.T, err error)
		wantSale func(t *testing.T, sale *Sale)
	}

	tests := []testCase{
		{
			name: "user not found",
			mockFn: func(userId string, amount float32) (*Sale, error) {
				return nil, users.ErrNotFound
			},
			userId: "non-existent-user-id",
			amount: 100.0,
			wantErr: func(t *testing.T, err error) {
				require.NotNil(t, err)
				require.Equal(t, users.ErrNotFound, err)
			},
			wantSale: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mockService := &MockSaleService{
				saleStorage: NewSaleStorage(),
				mockCreate:  tt.mockFn,
			}

			sale, err := mockService.Create(tt.userId, tt.amount)

			if tt.wantErr != nil {
				tt.wantErr(t, err)
			}

			if tt.wantSale != nil {
				tt.wantSale(t, sale)
			}
		})
	}
}
