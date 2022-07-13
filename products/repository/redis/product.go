package redis

import (
	"context"
	"encoding/json"
	"jumia/domain"

	"github.com/go-redis/redis/v8"
)

type redisInMemoryDB struct {
	db *redis.Client
}

func NewRedisInMemoryDB(redisClient *redis.Client) domain.ProductInMemoryDB {
	return &redisInMemoryDB{redisClient}
}

func (r *redisInMemoryDB) Get(ctx context.Context, sku string) ([]domain.Products, error) {
	var products []domain.Products
	val, err := r.db.Get(ctx, sku).Result()
	if err == redis.Nil {
		return nil, domain.ErrKeyNotFound
	} else if err != nil {
		return nil, err
	}
	json.Unmarshal([]byte(val), &products)
	return products, nil
}

func (r *redisInMemoryDB) Set(ctx context.Context, sku string, products []domain.Products) error {
	data, err := json.Marshal(products)
	if err != nil {
		return err
	}
	err = r.db.Set(ctx, sku, string(data), 0).Err()
	return err
}

func (r *redisInMemoryDB) Delete(ctx context.Context, sku string) error {
	err := r.db.Del(ctx, sku).Err()
	return err
}

func (r *redisInMemoryDB) DeleteMany(ctx context.Context, skus []string) error {
	err := r.db.Del(ctx, skus...).Err()
	return err
}
