package repos

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go-mux/models"
	"gorm.io/gorm"
	"testing"
)

var GormDB *gorm.DB

type RepoMockRepository struct {
	mock.Mock
}

func (mock *RepoMockRepository) CreateNewStore(db *gorm.DB, store *models.Store) error {
	args := mock.Called()
	res := args.Error(0)
	return res
}

func (mock *RepoMockRepository) GetProduct(db *gorm.DB, pId int) ([]models.Product, error) {
	args := mock.Called()
	result := args.Get(0)
	return result.([]models.Product), args.Error(1)
}

func (mock *RepoMockRepository) FetchProductsByIDList(db *gorm.DB, pIds []int) ([]models.Product, error) {
	args := mock.Called()
	result := args.Get(0)
	return result.([]models.Product), args.Error(1)
}

func (mock *RepoMockRepository) GetStoreProducts(db *gorm.DB, sId int) ([]models.StoreProducts, error) {
	args := mock.Called()
	result := args.Get(0)
	return result.([]models.StoreProducts), args.Error(1)
}

func (mock *RepoMockRepository) CheckStoreAvailableInDB(db *gorm.DB, sId int) ([]models.Store, error) {
	args := mock.Called()
	result := args.Get(0)
	return result.([]models.Store), args.Error(1)
}

func (mock *RepoMockRepository) AddProductsToStore(db *gorm.DB, StoreId int, s []*models.StoreProducts) error {
	args := mock.Called()
	res := args.Error(0)
	return res
}

func (mock *RepoMockRepository) GetAllProductsFromStore(db *gorm.DB, StoreId int) ([]models.Product, error) {
	args := mock.Called()
	result := args.Get(0)
	return result.([]models.Product), args.Error(1)
}

func TestCreateNewStore(t *testing.T) {
	mockRepo := new(RepoMockRepository)
	var store = models.Store{StoreId: 1, StoreName: "Store1"}
	mockRepo.On("CreateNewStore").Return(nil)
	err := mockRepo.CreateNewStore(GormDB, &store)
	mockRepo.AssertExpectations(t)
	assert.Nil(t, err)
}

func TestGetProduct(t *testing.T) {
	mockRepo := new(RepoMockRepository)
	var prod = models.Product{ID: 1, Name: "Prod1", Price: 10.0}
	mockRepo.On("GetProduct").Return([]models.Product{prod}, nil)
	res, err := mockRepo.GetProduct(GormDB, 10)
	mockRepo.AssertExpectations(t)
	assert.Nil(t, err)
	assert.Equal(t, 1, res[0].ID)
	assert.Equal(t, "Prod1", res[0].Name)
	assert.Equal(t, 10.0, res[0].Price)
}

func TestFetchProductsByIDList(t *testing.T) {
	mockRepo := new(RepoMockRepository)
	var prods1 = models.Product{ID: 1, Name: "Prod1", Price: 10.0}
	//var prods2 = models.Product{ID: 2, Name: "Prod2", Price: 20}
	mockRepo.On("FetchProductsByIDList").Return([]models.Product{prods1}, nil)
	//mockRepo.On("FetchProductsByIDList").Return([]models.Product{prods2}, nil)
	pIds := make([]int, 10)
	res, err := mockRepo.FetchProductsByIDList(GormDB, pIds)
	mockRepo.AssertExpectations(t)
	assert.Nil(t, err)
	assert.Equal(t, 1, res[0].ID)
	assert.Equal(t, "Prod1", res[0].Name)
	assert.Equal(t, 10.0, res[0].Price)
	//assert.Nil(t, err)
	//assert.Equal(t, 2, res[0].ID)
	//assert.Equal(t, "Prod2", res[0].Name)
	//assert.Equal(t, 20, res[0].Price)
}

func TestGetStoreProducts(t *testing.T) {
	mockRepo := new(RepoMockRepository)
	var storeProduct = models.StoreProducts{StoreId: 1, ProductId: 1, IsAvailable: true}
	mockRepo.On("GetStoreProducts").Return([]models.StoreProducts{storeProduct}, nil)
	res, err := mockRepo.GetStoreProducts(GormDB, 201)
	mockRepo.AssertExpectations(t)
	assert.Nil(t, err)
	assert.Equal(t, 1, res[0].StoreId)
	assert.Equal(t, 1, res[0].ProductId)
	assert.Equal(t, true, res[0].IsAvailable)
}

func TestCheckStoreAvailableInDB(t *testing.T) {
	mockRepo := new(RepoMockRepository)
	var store = models.Store{StoreId: 1, StoreName: "Store1"}
	mockRepo.On("CheckStoreAvailableInDB").Return([]models.Store{store}, nil)
	res, err := mockRepo.CheckStoreAvailableInDB(GormDB, 201)
	mockRepo.AssertExpectations(t)
	assert.Nil(t, err)
	assert.Equal(t, 1, res[0].StoreId)
	assert.Equal(t, "Store1", res[0].StoreName)
}

//func TestAddProductsToStore(t *testing.T) {
//	mockRepo := new(RepoMockRepository)
//	var storeProducts = models.StoreProducts{StoreId: 1, ProductId: 1, IsAvailable: true}
//	StoreId := 1
//	mockRepo.On("AddProductsToStore").Return(nil)
//	err := mockRepo.AddProductsToStore(GormDB, StoreId, storeProducts)
//	mockRepo.AssertExpectations(t)
//	assert.Nil(t, err)
//}
