package models

import "gorm.io/gorm"

type Product struct {
	gorm.Model
	Name        string `json:"name" binding:"required"`
	Brand       string `json:"brand" binding:"required"`
	Description string `json:"description" binding:"required"`
}

type productDaoImpl struct {
	dbObj *gorm.DB
}

type ProductDao interface {
	Insert(product *Product) error
	Get(name string) (*Product, error)
	GetAll() ([]Product, error)
}

var pDao ProductDao

func (pd *productDaoImpl) Insert(product *Product) error {
	result := pd.dbObj.Create(&product)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (pd *productDaoImpl) Get(name string) (*Product, error) {
	var product Product
	result := pd.dbObj.Where("name = ?", name).First(&product)

	if result.Error != nil {
		return nil, result.Error
	}

	return &product, nil
}

func (pd *productDaoImpl) GetAll() ([]Product, error) {
	var products []Product
	result := pd.dbObj.Find(&products)

	if result.Error != nil {
		return nil, result.Error
	}

	return products, nil
}

func ConfigureProductDao(db *gorm.DB) {
	pDao = &productDaoImpl{db}
}

func GetProductDao() ProductDao {
	return pDao
}
