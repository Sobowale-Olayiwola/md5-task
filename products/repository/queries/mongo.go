package queries

import (
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoQuery struct {
}

func (m MongoQuery) GetProductBySKU(sku string) interface{} {
	query := bson.D{{"sku", sku}}
	return query
}

func (m MongoQuery) ConsumeProductStock(sku string, amount int64) (filterQuery, updateQuery interface{}) {
	filter := bson.D{{"sku", sku}}
	update := bson.D{{"$inc", bson.D{{"stock_change", amount}}}, {"$set", bson.D{{"updatedAt", time.Now()}}}}
	return filter, update
}

func (m MongoQuery) BulkUpdateWithCSV(csvLines [][]string) (interface{}, []string) {
	var models []mongo.WriteModel
	var skus []string
	for _, lines := range csvLines {
		amount, _ := strconv.ParseInt(lines[3], 10, 64)
		skus = append(skus, lines[1])
		update := mongo.NewUpdateManyModel().SetFilter(bson.D{{"sku", lines[1]}}).SetUpdate(bson.D{{"$inc", bson.D{{"stock_change", amount}}}, {"$set", bson.D{{"updatedAt", time.Now()}}}, {"$setOnInsert", bson.D{{"name", lines[2]}, {"sku", lines[1]}, {"country", lines[0]}, {"createdAt", time.Now()}}}}).SetUpsert(true)
		models = append(models, update)
	}
	return models, skus
}
