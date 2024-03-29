package souko

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/mattn/go-sqlite3"
	"gorm.io/gorm"
	"net/http"
	"souko/models"
	"strconv"
)

type MultiPageProductResponse struct {
	Products      []models.Product `json:"products"`
	NextPageToken int              `json:"nextPageToken"`
}

func configureRouter() *gin.Engine {
	r := gin.Default()
	v1 := r.Group("/v1")
	{
		v1.POST("/products", createProductHandler)
		v1.GET("/products", getProductsHandler)
		v1.GET("/products/:id", getProductHandler)
		v1.PUT("/products/:id", modifyProductHandler)
		v1.DELETE("/products/:id", deleteProductHandler)
	}

	return r
}

func createProductHandler(c *gin.Context) {
	var product *models.Product
	productDao := models.GetProductDao()

	err := c.ShouldBindJSON(&product)
	if validateError(c, http.StatusBadRequest, err) {
		return
	}

	_, err = productDao.GetByName(product.Name)

	if err != nil && err != gorm.ErrRecordNotFound {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	err = productDao.Insert(product)

	if validateError(c, http.StatusInternalServerError, err) {
		return
	}

	c.JSON(http.StatusOK, product)
}

func getProductHandler(c *gin.Context) {
	productDao := models.GetProductDao()
	id, err := strconv.Atoi(c.Param("id"))

	if validateError(c, http.StatusBadRequest, err) {
		return
	}

	product, err := productDao.GetById(id)

	if validateError(c, http.StatusBadRequest, err) {
		return
	}

	c.JSON(http.StatusOK, product)
}

func getProductsHandler(c *gin.Context) {
	productDao := models.GetProductDao()
	token, err := strconv.Atoi(c.DefaultQuery("token", "1"))

	if validateError(c, http.StatusBadRequest, err) {
		return
	}

	size, err := strconv.Atoi(c.DefaultQuery("size", "5"))

	if validateError(c, http.StatusBadRequest, err) {
		return
	}

	if token < 0 || size < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "An invalid value for size or token has been set"})
		return
	}

	products, nextPageId, err := productDao.GetPage(token, size)

	if validateError(c, http.StatusInternalServerError, err) {
		return
	}

	c.JSON(http.StatusOK, MultiPageProductResponse{products, nextPageId})
}

func modifyProductHandler(c *gin.Context) {
	productDao := models.GetProductDao()
	id, err := strconv.Atoi(c.Param("id"))

	if validateError(c, http.StatusBadRequest, err) {
		return
	}

	_, err = productDao.GetById(id)

	if err != nil && err == gorm.ErrRecordNotFound {
		errorStr := fmt.Sprintf("Product with id %v doesn't exist", id)
		c.JSON(http.StatusNotFound, gin.H{"error": errorStr})
		return
	}

	var tempProduct *models.Product
	err = c.ShouldBindJSON(&tempProduct)

	if validateError(c, http.StatusBadRequest, err) {
		return
	}

	updatedProduct, err := productDao.Update(id, tempProduct)
	if validateError(c, http.StatusInternalServerError, err) {
		return
	}

	c.JSON(http.StatusOK, updatedProduct)
}

func deleteProductHandler(c *gin.Context) {
	productDao := models.GetProductDao()
	id, err := strconv.Atoi(c.Param("id"))

	if validateError(c, http.StatusBadRequest, err) {
		return
	}

	err = productDao.Delete(id)

	if err != nil && err != gorm.ErrRecordNotFound {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
	} else if validateError(c, http.StatusInternalServerError, err) {
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "The product has been deleted"})
}

func validateError(c *gin.Context, statusCode int, err error) bool {
	if err != nil {
		if sqliteErr, ok := err.(sqlite3.Error); ok {
			if sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
				c.JSON(http.StatusConflict, gin.H{"error": "This product already exists"})
			}
		} else {
			c.JSON(statusCode, gin.H{"error": err})
		}
		return true
	}
	return false
}
