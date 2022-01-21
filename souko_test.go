package souko

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"souko/models"
	"testing"
)

func TestCreateProduct(t *testing.T) {
	productJsonStr := []byte(`{"name": "Pixel 6 Pro", "brand": "Google", "description": "The most intelligent smartphone"}`)
	router := configureRouter()
	db := ConfigureDatabase()
	dbObj, _ := db.DB()
	ts := httptest.NewServer(router)

	defer dbObj.Close()
	defer ts.Close()

	req, _ := http.NewRequest(http.MethodPost, "/products", bytes.NewBuffer(productJsonStr))
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
			t.Errorf("Error while comparing products in request")
		}
	})

	t.Run("Validate object in database", func(t *testing.T) {
		productDao := models.GetProductDao()
		dbProduct, _ := productDao.GetByName(productJson.Name)

		if dbProduct == nil {
			t.Errorf("The product created in the server and the one in the database aren't the same")
		}
	})

	t.Run("Same product should not be created twice", func(t *testing.T) {
		req, _ = http.NewRequest(http.MethodPost, "/products", bytes.NewBuffer(productJsonStr))
		res = httptest.NewRecorder()
		router.ServeHTTP(res, req)
		if res.Code != http.StatusConflict {
			t.Errorf("Product conflict should have been detected, but code %v received instead", res.Code)
		}
	})
}

func TestReadProduct(t *testing.T) {
	router := configureRouter()
	db := ConfigureDatabase()
	dbObj, _ := db.DB()
	ts := httptest.NewServer(router)
	defer dbObj.Close()
	defer ts.Close()

	productDao := models.GetProductDao()
	product := &models.Product{
		Name:        "Pixel 6 Pro",
		Brand:       "Google",
		Description: "The most intelligent smartphone",
	}

	productDao.Insert(product)

	t.Run("Search by ID", func(t *testing.T) {
		var tempProduct *models.Product
		req, _ := http.NewRequest(http.MethodGet, "/products/1", nil)
		res := httptest.NewRecorder()
		router.ServeHTTP(res, req)
		err := json.NewDecoder(res.Body).Decode(&tempProduct)
		if err != nil {
			t.Errorf("Error while decoding JSON")
		}
		if tempProduct.Name != product.Name {
			t.Errorf("Both products are not the same")
		}
	})
}

func TestModifyProduct(t *testing.T) {
	router := configureRouter()
	db := ConfigureDatabase()
	dbObj, _ := db.DB()
	ts := httptest.NewServer(router)
	defer dbObj.Close()
	defer ts.Close()

	productDao := models.GetProductDao()
	product := &models.Product{
		Name:        "Pixel 6 Pro",
		Brand:       "Google",
		Description: "The most intelligent smartphone",
	}

	productDao.Insert(product)
	newProductJsonStr := []byte(`{"name": "iPhone 12 Pro Max", "brand": "Apple", "description": "The most popular smartphone"}`)

	var tempProduct *models.Product
	var newProduct *models.Product
	req, _ := http.NewRequest(http.MethodPut, "/products/1", bytes.NewReader(newProductJsonStr))
	res := httptest.NewRecorder()
	router.ServeHTTP(res, req)
	_ = json.NewDecoder(bytes.NewReader(newProductJsonStr)).Decode(&tempProduct)
	err := json.NewDecoder(res.Body).Decode(&newProduct)

	if err != nil {
		t.Errorf("Error while decoding JSON")
	}

	if tempProduct.Name != newProduct.Name {
		t.Errorf("Updated product is not the same from memory and HTTP request")
	}
}
