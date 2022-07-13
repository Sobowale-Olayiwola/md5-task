package repository

import (
	"context"
	"jumia/domain"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

var wc = writeconcern.New(writeconcern.WMajority())
var rc = readconcern.Snapshot()
var txnOpts = options.Transaction().SetWriteConcern(wc).SetReadConcern(rc)

type mongoProducts struct {
	db         *mongo.Database
	client     *mongo.Client
	collection string
}

func NewMongoProductRepository(db *mongo.Database, client *mongo.Client, collection string) domain.ProductRepository {
	return &mongoProducts{db: db, collection: collection, client: client}
}

func (m *mongoProducts) GetProductBySKU(ctx context.Context, query interface{}) ([]domain.Products, error) {
	callback := func(sessionContext mongo.SessionContext) (interface{}, error) {
		products := []domain.Products{}
		cursor, err := m.db.Collection(m.collection).Find(ctx, query)
		if err != nil {
			defer cursor.Close(ctx)
			return nil, err
		}
		if err := cursor.All(ctx, &products); err != nil {
			return nil, err
		}
		return products, nil
	}
	session, err := m.client.StartSession()
	if err != nil {
		panic(err)
	}
	defer session.EndSession(context.Background())
	products, err := session.WithTransaction(context.Background(), callback, txnOpts)
	if err != nil {
		return nil, err
	}
	return products.([]domain.Products), nil
}

func (m *mongoProducts) ConsumeProductStock(ctx context.Context, filterQuery, updateQuery interface{}) error {
	callback := func(sessionContext mongo.SessionContext) (interface{}, error) {
		result, err := m.db.Collection(m.collection).UpdateMany(context.TODO(), filterQuery, updateQuery)
		if result.MatchedCount == 0 {
			return nil, domain.ErrRecordNotFound
		}
		if err != nil {
			return nil, err
		}
		return result, nil
	}
	session, err := m.client.StartSession()
	if err != nil {
		panic(err)
	}
	defer session.EndSession(context.Background())
	_, err = session.WithTransaction(context.Background(), callback, txnOpts)
	if err != nil {
		return err
	}
	return nil

}

func (m *mongoProducts) BulkUpdateWithCSV(ctx context.Context, query interface{}) error {
	opts := options.BulkWrite().SetOrdered(false)
	callback := func(sessionContext mongo.SessionContext) (interface{}, error) {
		res, err := m.db.Collection(m.collection).BulkWrite(context.TODO(), query.([]mongo.WriteModel), opts)
		if err != nil {
			return nil, err
		}
		return res, nil
	}
	session, err := m.client.StartSession()
	if err != nil {
		panic(err)
	}
	defer session.EndSession(context.Background())
	_, err = session.WithTransaction(context.Background(), callback, txnOpts)
	if err != nil {
		return err
	}
	return nil
}
