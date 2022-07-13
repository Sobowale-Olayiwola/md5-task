package service

import (
	"context"
	"errors"
	"jumia/domain"
	inmemorydb "jumia/domain/mocks/inmemory"
	"jumia/domain/mocks/repository"
	"jumia/products/repository/queries"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetProductBySKU(t *testing.T) {
	as := assert.New(t)
	inmem := &inmemorydb.InMemoryDBMock{}
	productRepo := &repository.ProductRepositoryMock{}
	t.Run("happy path: Successfully gets product by sku", func(t *testing.T) {
		sku := "9befa247cd11"
		inmem.On("Get", context.Background(), sku).Return([]domain.Products{}, domain.ErrKeyNotFound).Once()
		inmem.On("Set", context.Background(), sku, mock.Anything).Return(nil).Once()
		productRepo.On("GetProductBySKU", context.Background(), mock.Anything).Return([]domain.Products{
			{
				Country:     "ke",
				SKU:         "9befa247cd11",
				StockChange: 100,
			},
		}, nil).Once()
		service := NewProductService(productRepo, inmem, queries.MongoQuery{})
		products, err := service.GetProductBySKU(context.Background(), sku)
		as.NoError(err)
		as.Equal(len(products), 1)
		as.Equal(products[0].Country, "ke")
		inmem.AssertExpectations(t)
		productRepo.AssertExpectations(t)
	})
}

func TestConsumeProductStock(t *testing.T) {
	as := assert.New(t)
	inmem := &inmemorydb.InMemoryDBMock{}
	productRepo := &repository.ProductRepositoryMock{}
	t.Run("happy path: Successfully consumes product stock", func(t *testing.T) {
		sku := "9befa247cd11"
		amount := int64(-20)
		inmem.On("Delete", context.Background(), sku, mock.Anything).Return(nil).Once()
		productRepo.On("ConsumeProductStock", context.Background(), mock.Anything, mock.Anything).Return(nil).Once()
		service := NewProductService(productRepo, inmem, queries.MongoQuery{})
		err := service.ConsumeProductStock(context.Background(), sku, amount)
		as.NoError(err)
		inmem.AssertExpectations(t)
		productRepo.AssertExpectations(t)
	})
	t.Run("system error: Database failed", func(t *testing.T) {
		sku := "9befa247cd11"
		amount := int64(-20)
		productRepo.On("ConsumeProductStock", context.Background(), mock.Anything, mock.Anything).Return(errors.New("an error occured")).Once()
		service := NewProductService(productRepo, inmem, queries.MongoQuery{})
		err := service.ConsumeProductStock(context.Background(), sku, amount)
		as.Error(err)
		inmem.AssertExpectations(t)
		productRepo.AssertExpectations(t)
	})
}

func TestBulkUpdateWithCSV(t *testing.T) {
	as := assert.New(t)
	inmem := &inmemorydb.InMemoryDBMock{}
	productRepo := &repository.ProductRepositoryMock{}
	t.Run("happy path: Successfully consumes product stock", func(t *testing.T) {
		csvLines := [][]string{{"ke", "9befa247cd11", "New jumia", "4"}, {"dz", "04ea8ee8dccc", "White Group Shoes", "-4"}}
		inmem.On("DeleteMany", context.Background(), mock.Anything, mock.Anything).Return(nil).Once()
		productRepo.On("BulkUpdateWithCSV", context.Background(), mock.Anything).Return(nil).Once()
		service := NewProductService(productRepo, inmem, queries.MongoQuery{})
		err := service.BulkUpdateWithCSV(context.Background(), csvLines)
		as.NoError(err)
		inmem.AssertExpectations(t)
		productRepo.AssertExpectations(t)
	})
	t.Run("system error: Database failed", func(t *testing.T) {
		csvLines := [][]string{{"ke", "9befa247cd11", "New jumia", "4"}, {"dz", "04ea8ee8dccc", "White Group Shoes", "-4"}}
		productRepo.On("BulkUpdateWithCSV", context.Background(), mock.Anything, mock.Anything).Return(errors.New("an error occured")).Once()
		service := NewProductService(productRepo, inmem, queries.MongoQuery{})
		err := service.BulkUpdateWithCSV(context.Background(), csvLines)
		as.Error(err)
		inmem.AssertExpectations(t)
		productRepo.AssertExpectations(t)
	})
}
