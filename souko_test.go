package souko

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"souko/models"
	"testing"
)

func TestProducts(t *testing.T) {
	productJsonStr := []byte(`{"name": "Pixel 6 Pro", "brand": "Google", "description": "The most intelligent smartphone"}`)
	router := configureRouter()
	db := ConfigureDatabase()
	dbObj, _ := db.DB()
	ts := httptest.NewServer(router)

	defer dbObj.Close()
	defer ts.Close()

	req, _ := http.NewRequest(http.MethodPost, "/product", bytes.NewBuffer(productJsonStr))
	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()

	var productJson models.Product
	json.NewDecoder(bytes.NewBuffer(productJsonStr)).Decode(&productJson)

	t.Run("HTTP request", func(t *testing.T) {
		router.ServeHTTP(res, req)

		if res.Code != http.StatusOK {
			t.Errorf("Couldn't create product successfully")
		}

		var decodedProduct models.Product
		err := json.NewDecoder(res.Body).Decode(&decodedProduct)

		if err != nil || decodedProduct.Name != productJson.Name {
			t.Errorf("Error while comparing product")
		}
	})

	t.Run("Validate object in database", func(t *testing.T) {
		productDao := models.GetProductDao()
		dbProduct, _ := productDao.Get(productJson.Name)

		if dbProduct == nil {
			t.Errorf("The product created in the server and the one in the database aren't the same")
		}
	})
}
