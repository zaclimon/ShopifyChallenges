package souko

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"souko/models"
	"strconv"
)

func configureRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/products/:id", getProductHandler)
	r.POST("/products", createProductHandler)
	return r
}

func createProductHandler(c *gin.Context) {
	var product *models.Product
	productDao := models.GetProductDao()

	err := c.ShouldBindJSON(&product)
	if validateError(c, http.StatusBadRequest, err) {
		return
	}

	product, err = productDao.GetByName(product.Name)
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

func validateError(c *gin.Context, statusCode int, err error) bool {
	if err != nil {
		c.JSON(statusCode, gin.H{"error": err})
		return true
	}
	return false
}
