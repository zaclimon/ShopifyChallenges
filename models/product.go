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
	GetPage(token int, size int) ([]Product, int, error)
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

// Inspired from: https://cloud.google.com/apis/design/design_patterns#list_pagination
func (pd *productDaoImpl) GetPage(token int, size int) ([]Product, int, error) {
	var products []Product
	var productCount int64
	var nextTokenId int
	result := pd.dbObj.Order("id").Where("id >= ?", token).Limit(size).Find(&products)
	pd.dbObj.Model(&Product{}).Count(&productCount)

	if result.Error != nil {
		return nil, -1, result.Error
	}

	// We assume that the token will always be positive since it's auto-incrementing
	if int64(token+size) > productCount {
		nextTokenId = -1
	} else {
		nextTokenId = token + size
	}

	return products, nextTokenId, nil
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
