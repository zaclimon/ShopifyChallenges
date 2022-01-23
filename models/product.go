package models

import "gorm.io/gorm"

type Product struct {
	gorm.Model
	Name        string `json:"name" binding:"required" gorm:"unique"`
	Brand       string `json:"brand" binding:"required"`
	Description string `json:"description" binding:"required"`
}

type productDaoImpl struct {
	dbObj *gorm.DB
}

type ProductDao interface {
	Insert(product *Product) error
	GetById(id int) (*Product, error)
	GetByName(name string) (*Product, error)
	GetAll() ([]Product, error)
	Update(id int, newProduct *Product) (*Product, error)
	Delete(id int) error
}

var pDao ProductDao

func (pd *productDaoImpl) Insert(product *Product) error {
	result := pd.dbObj.Create(&product)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (pd *productDaoImpl) GetById(id int) (*Product, error) {
	var product Product
	result := pd.dbObj.First(&product, id)

	if result.Error != nil {
		return nil, result.Error
	}
	return &product, nil
}

func (pd *productDaoImpl) GetByName(name string) (*Product, error) {
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

func (pd *productDaoImpl) Update(id int, newProduct *Product) (*Product, error) {
	var currentProduct Product
	pd.dbObj.Find(&currentProduct, id)

	if newProduct.Name != "" {
		currentProduct.Name = newProduct.Name
	}
	if newProduct.Brand != "" {
		currentProduct.Brand = newProduct.Brand
	}
	if newProduct.Description != "" {
		currentProduct.Description = newProduct.Description
	}

	result := pd.dbObj.Save(&currentProduct)

	if result.Error != nil {
		return nil, result.Error
	}
	return &currentProduct, nil
}

func (pd *productDaoImpl) Delete(id int) error {
	result := pd.dbObj.Delete(&Product{}, id)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func ConfigureProductDao(db *gorm.DB) {
	pDao = &productDaoImpl{db}
}

func GetProductDao() ProductDao {
	return pDao
}
