package souko

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"souko/models"
)

func configureRouter() *gin.Engine {
	r := gin.Default()
	r.POST("/product", createProductHandler)
	return r
}

func createProductHandler(c *gin.Context) {
	var product *models.Product
	productDao := models.GetProductDao()

	err := c.ShouldBindJSON(&product)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	_, err = productDao.Get(product.Name)
	if err != nil && err != gorm.ErrRecordNotFound {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	err = productDao.Insert(product)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, product)
}
