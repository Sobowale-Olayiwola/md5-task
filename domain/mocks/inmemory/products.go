package inmemorydb

import (
	"context"
	"jumia/domain"

	"github.com/stretchr/testify/mock"
)

type InMemoryDBMock struct {
	mock.Mock
}

func (w *InMemoryDBMock) Get(ctx context.Context, sku string) ([]domain.Products, error) {
	output := w.Mock.Called(ctx, sku)
	products := output.Get(0)
	err := output.Error(1)
	return products.([]domain.Products), err
}

func (w *InMemoryDBMock) Set(ctx context.Context, sku string, products []domain.Products) error {
	output := w.Mock.Called(ctx, sku, products)
	err := output.Error(0)
	return err
}

func (w *InMemoryDBMock) Delete(ctx context.Context, sku string) error {
	output := w.Mock.Called(ctx, sku)
	err := output.Error(0)
	return err
}

func (w *InMemoryDBMock) DeleteMany(ctx context.Context, sku []string) error {
	output := w.Mock.Called(ctx, sku)
	err := output.Error(0)
	return err
}
