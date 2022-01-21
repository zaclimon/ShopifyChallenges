package souko

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
	"souko/models"
	"testing"
)

func TestCreateProduct(t *testing.T) {
	router, db := configureBaseComponents(t)
	defer db.Close()

	var productJson models.Product
	productJsonStr := `{"name": "PlayStation 5", "brand": "Sony", "description": "The contender to the most popular game console"}`
	json.NewDecoder(bytes.NewBuffer([]byte(productJsonStr))).Decode(&productJson)

	t.Run("HTTP request", func(t *testing.T) {
		res := executeHttpRequest(t, router, http.MethodPost, "/products", productJsonStr)

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
		res := executeHttpRequest(t, router, http.MethodPost, "/products", productJsonStr)
		if res.Code != http.StatusConflict {
			t.Errorf("Product conflict should have been detected, but code %v received instead", res.Code)
		}
	})
}

func TestReadProduct(t *testing.T) {
	router, db := configureBaseComponents(t)
	defer db.Close()

	t.Run("Search by ID", func(t *testing.T) {
		var tempProduct *models.Product
		res := executeHttpRequest(t, router, http.MethodGet, "/products/1", "")
		err := json.NewDecoder(res.Body).Decode(&tempProduct)

		if err != nil {
			t.Errorf("Error while decoding JSON")
		}

		if tempProduct.Name != getTestProducts()[0].Name {
			t.Errorf("Both products are not the same")
		}
	})
}

func TestModifyProduct(t *testing.T) {
	router, db := configureBaseComponents(t)
	defer db.Close()
	newProductJsonStr := `{"name": "iPhone 12 Pro Max", "brand": "Apple", "description": "The 'previous' most popular smartphone"}`

	var tempProduct *models.Product
	var newProduct *models.Product
	res := executeHttpRequest(t, router, http.MethodPut, "/products/1", newProductJsonStr)
	json.NewDecoder(bytes.NewReader([]byte(newProductJsonStr))).Decode(&tempProduct)
	json.NewDecoder(res.Body).Decode(&newProduct)

	if tempProduct.Name != newProduct.Name {
		t.Errorf("Updated product is not the same from memory and HTTP request")
	}

	t.Run("Modify one product to the same name as another one", func(t *testing.T) {
		res = executeHttpRequest(t, router, http.MethodPut, "/products/2", newProductJsonStr)
		if res.Code != http.StatusConflict {
			t.Errorf("Expecting product duplicate, while error code is: %v instead", res.Code)
		}
	})
}

func configureBaseComponents(t *testing.T) (*gin.Engine, *sql.DB) {
	t.Helper()
	router := configureRouter()
	db := ConfigureDatabase()
	products := getTestProducts()
	productDao := models.GetProductDao()
	dbObj, _ := db.DB()

	for _, product := range products {
		productDao.Insert(&product)
	}

	return router, dbObj
}

func executeHttpRequest(t *testing.T, router *gin.Engine, httpMethod string, url string, requestBody string) *httptest.ResponseRecorder {
	t.Helper()
	var req *http.Request

	if requestBody != "" {
		req, _ = http.NewRequest(httpMethod, url, bytes.NewReader([]byte(requestBody)))
	} else {
		req, _ = http.NewRequest(httpMethod, url, nil)
	}

	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()
	router.ServeHTTP(res, req)
	return res
}

func getTestProducts() []models.Product {
	return []models.Product{
		{Name: "Pixel 6 Pro", Brand: "Google", Description: "The most intelligent smartphone"},
		{Name: "iPhone 13 Pro Max", Brand: "Apple", Description: "The most popular smartphone"},
		{Name: "Switch", Brand: "Nintendo", Description: "The most popular portable game console"},
	}
}
