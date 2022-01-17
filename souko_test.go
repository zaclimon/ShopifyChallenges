package souko

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestProducts(t *testing.T) {
	productJsonStr := []byte(`{"name": "Pixel 6 Pro", "brand": "Google", "description": "The most intelligent smartphone"}`)
	router := configureRouter()
	ts := httptest.NewServer(router)
	defer ts.Close()

	req, _ := http.NewRequest(http.MethodPost, "/product", bytes.NewBuffer(productJsonStr))
	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()

	var productJson Product
	json.NewDecoder(bytes.NewBuffer(productJsonStr)).Decode(&productJson)

	t.Run("HTTP request", func(t *testing.T) {
		router.ServeHTTP(res, req)

		if res.Code != http.StatusOK {
			t.Errorf("Couldn't create product successfully")
		}

		var decodedProduct Product
		err := json.NewDecoder(res.Body).Decode(&decodedProduct)

		if err != nil || decodedProduct.Name != productJson.Name {
			t.Errorf("Error while comparing product")
		}
	})

	t.Run("Validate object in database", func(t *testing.T) {
		dbProduct, _ := productDao.Get(productJson.Name)

		if dbProduct == nil {
			t.Errorf("The product created in the server and the one in the database aren't the same")
		}
	})
}
