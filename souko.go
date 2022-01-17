package souko

import (
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"net/http"
)

type Product struct {
	gorm.Model
	Name        string `json:"name" binding:"required"`
	Brand       string `json:"brand" binding:"required"`
	Description string `json:"description" binding:"required"`
}

type ProductDaoImpl struct {
	dbObj *gorm.DB
}

type ProductDao interface {
	Insert(product *Product) error
	Get(name string) (*Product, error)
	GetAll() ([]Product, error)
}

var productDao = configureDatabase()

func configureRouter() *gin.Engine {
	r := gin.Default()
	r.POST("/product", createProductHandler)
	return r
}

func configureDatabase() ProductDao {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("Could not connect to the database")
	}

	err = db.AutoMigrate(&Product{})
	if err != nil {
		panic("Could not migrate database objects")
	}

	return &ProductDaoImpl{db}
}

func createProductHandler(c *gin.Context) {
	var product *Product

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

func (pd *ProductDaoImpl) Insert(product *Product) error {
	result := pd.dbObj.Create(&product)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (pd *ProductDaoImpl) Get(name string) (*Product, error) {
	var product Product
	result := pd.dbObj.Where("name = ?", name).First(&product)

	if result.Error != nil {
		return nil, result.Error
	}

	return &product, nil
}

func (pd *ProductDaoImpl) GetAll() ([]Product, error) {
	var products []Product
	result := pd.dbObj.Find(&products)

	if result.Error != nil {
		return nil, result.Error
	}

	return products, nil
}
