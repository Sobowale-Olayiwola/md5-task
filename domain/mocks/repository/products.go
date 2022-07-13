package repository

import (
	"context"
	"jumia/domain"

	"github.com/stretchr/testify/mock"
)

type ProductRepositoryMock struct {
	mock.Mock
}

func (j *ProductRepositoryMock) GetProductBySKU(ctx context.Context, query interface{}) ([]domain.Products, error) {
	output := j.Mock.Called(ctx, query)
	products := output.Get(0)
	err := output.Error(1)
	return products.([]domain.Products), err
}

func (j *ProductRepositoryMock) ConsumeProductStock(ctx context.Context, filterQuery, updateQuery interface{}) error {
	output := j.Mock.Called(ctx, filterQuery, updateQuery)
	err := output.Error(0)
	return err
}

func (j *ProductRepositoryMock) BulkUpdateWithCSV(ctx context.Context, query interface{}) error {
	output := j.Mock.Called(ctx, query)
	err := output.Error(0)
	return err
}
