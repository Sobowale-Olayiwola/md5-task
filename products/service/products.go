package service

import (
	"context"
	"jumia/domain"
	"jumia/products/repository/queries"
)

type productService struct {
	productRepository domain.ProductRepository
	dbQueries         queries.ProductQueries
	inMemoryDB        domain.ProductInMemoryDB
}

func NewProductService(j domain.ProductRepository, r domain.ProductInMemoryDB, q queries.ProductQueries) domain.ProductService {
	return &productService{productRepository: j, dbQueries: q, inMemoryDB: r}
}

func (p *productService) GetProductBySKU(ctx context.Context, sku string) ([]domain.Products, error) {
	products, err := p.inMemoryDB.Get(ctx, sku)
	if err == domain.ErrKeyNotFound {
		query := p.dbQueries.GetProductBySKU(sku)
		products, err = p.productRepository.GetProductBySKU(ctx, query)
		if err != nil {
			return nil, err
		}
		if len(products) == 0 {
			return nil, domain.ErrRecordNotFound
		}
		p.inMemoryDB.Set(ctx, sku, products)
	}
	return products, nil
}

func (p *productService) ConsumeProductStock(ctx context.Context, sku string, amount int64) error {
	filterQuery, updateQuery := p.dbQueries.ConsumeProductStock(sku, amount)
	err := p.productRepository.ConsumeProductStock(ctx, filterQuery, updateQuery)
	if err != nil {
		return err
	}
	p.inMemoryDB.Delete(ctx, sku)
	return nil
}

func (p *productService) BulkUpdateWithCSV(ctx context.Context, csvLines [][]string) error {
	query, skus := p.dbQueries.BulkUpdateWithCSV(csvLines)
	err := p.productRepository.BulkUpdateWithCSV(ctx, query)
	if err != nil {
		return err
	}
	p.inMemoryDB.DeleteMany(ctx, skus)
	return nil
}
