package queries

type ProductQueries interface {
	GetProductBySKU(sku string) interface{}
	ConsumeProductStock(sku string, amount int64) (filterQuery, updateQuery interface{})
	BulkUpdateWithCSV(csvLines [][]string) (query interface{}, skus []string)
}
