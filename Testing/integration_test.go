package Testing

import (
	"Rodrigo-Guinazu-Arias-TpFinal/api"
	"Rodrigo-Guinazu-Arias-TpFinal/models/Sales"
	"Rodrigo-Guinazu-Arias-TpFinal/models/users"
	"Rodrigo-Guinazu-Arias-TpFinal/utils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"resty.dev/v3"
	"testing"
)

func TestSaleServiceIntegration(t *testing.T) {

	gin.SetMode(gin.TestMode)

	saleService := Sales.NewSaleService(Sales.NewSaleStorage())
	userService := users.NewUserService(users.NewUserStorage())

	userList, saleList := utils.InitSystem(saleService, userService)
	assert.NotEmpty(t, userList)
	assert.NotEmpty(t, saleList)

	r := gin.Default()
	api.InitRoutes(r, userService, saleService)

	t.Run("Create Sale", func(t *testing.T) {

		reqBody := map[string]interface{}{
			"user_id": userList[0].ID,
			"amount":  150.75,
		}

		server := httptest.NewServer(r)
		defer server.Close()

		client := resty.New()

		var sale Sales.Sale

		resp, err := client.R().
			SetHeader("Content-Type", "application/json").
			SetBody(reqBody).
			SetResult(&sale).
			Post(server.URL + "/sales")

		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode())

		assert.NotEmpty(t, sale.Id)
		assert.Equal(t, userList[0].ID, sale.UserId)
		assert.Equal(t, float32(150.75), sale.Amount)
	})

	t.Run("Update Sale", func(t *testing.T) {
		// Find a pending sale to update
		var pendingSale *Sales.Sale
		for _, sale := range saleList {
			if sale.Status == Sales.Pending {
				pendingSale = sale
				break
			}
		}

		if pendingSale == nil {
			t.Skip("No pending sale found to update")
		}

		reqBody := map[string]interface{}{
			"status": Sales.Aproved,
		}

		server := httptest.NewServer(r)
		defer server.Close()

		client := resty.New()

		resp, err := client.R().
			SetHeader("Content-Type", "application/json").
			SetBody(reqBody).
			Patch(server.URL + "/sales/" + pendingSale.Id)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode())

		var getResponse struct {
			Metadata struct {
				Quantity    int     `json:"quantity"`
				Aproved     int     `json:"aproved"`
				Rejected    int     `json:"rejected"`
				Pending     int     `json:"pending"`
				TotalAmount float32 `json:"total_amount"`
			} `json:"metadata"`
			Sales []Sales.Sale `json:"results"`
		}

		getResp, err := client.R().
			SetResult(&getResponse).
			Get(server.URL + "/sales?user_id=" + pendingSale.UserId + "&status=" + Sales.Aproved)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, getResp.StatusCode())

		var found bool
		for _, sale := range getResponse.Sales {
			if sale.Id == pendingSale.Id {
				assert.Equal(t, Sales.Aproved, sale.Status)
				found = true
				break
			}
		}
		assert.True(t, found, "Updated sale not found in GET response")
	})

	t.Run("Get Sales By User", func(t *testing.T) {

		server := httptest.NewServer(r)
		defer server.Close()

		client := resty.New()

		var response struct {
			Metadata struct {
				Quantity    int     `json:"quantity"`
				Aproved     int     `json:"aproved"`
				Rejected    int     `json:"rejected"`
				Pending     int     `json:"pending"`
				TotalAmount float32 `json:"total_amount"`
			} `json:"metadata"`
			Sales []Sales.Sale `json:"results"`
		}

		resp, err := client.R().
			SetResult(&response).
			Get(server.URL + "/sales?user_id=" + userList[0].ID)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode())

		for _, sale := range response.Sales {
			assert.Equal(t, userList[0].ID, sale.UserId)
		}
	})
}
