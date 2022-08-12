//go:generate mockgen -source repo.go -destination mock/repo_mock.go -package mock
package repos

import (
	"fmt"
	"go-mux/models"
	"gorm.io/gorm"
)

type Repos struct {
}

type RepoInterface interface {
	CreateNewStore(db *gorm.DB, store *models.Store)
	GetProduct(db *gorm.DB, pId int) ([]*models.Product, error)
	FetchProductsByIDList(db *gorm.DB, pIds []int) ([]*models.Product, error)
	GetStoreProducts(db *gorm.DB, sId int) ([]*models.StoreProducts, error)
	checkStoreAvailableInDB(db *gorm.DB, sId int) ([]*models.Store, error)
	AddProductsToStore(db *gorm.DB, StoreId int, storeProducts []*models.StoreProducts) error
	GetAllProductsFromStore(db *gorm.DB, StoreId int) ([]*models.Product, error)
}

// create new store

func (repo *Repos) CreateNewStore(db *gorm.DB, store *models.Store) error {
	return db.Create(store).Error
}

// product retrieval based on id

func (repo *Repos) GetProduct(db *gorm.DB, pId int) ([]*models.Product, error) {
	var prods []*models.Product
	response := db.Find(&models.Product{ID: pId}).First(&prods)
	err := response.Error
	if err != nil {
		return nil, err
	}
	return prods, nil
}

// get all the products based on IDs

func (repo *Repos) FetchProductsByIDList(db *gorm.DB, pIds []int) ([]*models.Product, error) {
	var prods []*models.Product
	response := db.Where(pIds).Find(&prods)
	err := response.Error
	if err != nil {
		return nil, err
	}
	return prods, nil
}

// retrieval of all stores products in store based on storeID

func (repo *Repos) GetStoreProducts(db *gorm.DB, sId int) ([]*models.StoreProducts, error) {
	var storeProduct []*models.StoreProducts
	response := db.Where(&models.StoreProducts{StoreId: sId}).Find(&storeProduct)
	fmt.Println("Getting all store products")
	err := response.Error
	if err != nil {
		return nil, err
	}
	return storeProduct, nil
}

// check store availability

func (repo *Repos) checkStoreAvailableInDB(db *gorm.DB, sId int) ([]*models.Store, error) {
	var storeCheck []*models.Store
	response := db.Where(&models.Store{StoreId: sId}).First(&storeCheck)
	err := response.Error
	fmt.Println(err)
	if err != nil {
		return nil, err
	}
	return storeCheck, nil
}

// add products service

func (repo *Repos) AddProductsToStore(db *gorm.DB, StoreId int, s []*models.StoreProducts) error {
	for _, sProd := range s {
		_, err := repo.GetProduct(db, sProd.ProductId)
		if err != nil {
			return err
		}
		sProd.StoreId = StoreId
	}
	return db.Create(&s).Error
}

// get products service

func (repo *Repos) GetAllProductsFromStore(db *gorm.DB, StoreId int) ([]*models.Product, error) {
	sProds, err := repo.GetStoreProducts(db, StoreId)
	if err != nil {
		return nil, err
	}
	var prodIds []int
	for _, prods := range sProds {
		prodIds = append(prodIds, prods.ProductId)
	}
	products, err := repo.FetchProductsByIDList(db, prodIds)
	if err != nil {
		return nil, err
	}
	return products, nil
}
