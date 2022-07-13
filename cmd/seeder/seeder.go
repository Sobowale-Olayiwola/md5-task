package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"jumia/domain"
	"jumia/internal/constants"
	"jumia/internal/helpers"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func init() {
	godotenv.Load()
}

func SeedDataDynamically() {
	csvFileNames := make([]string, 0, 10)
	files, err := ioutil.ReadDir(".")
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		if strings.Contains(file.Name(), ".csv") {
			csvFileNames = append(csvFileNames, file.Name())
		}
	}
	conc := len(csvFileNames)
	csvLines := make(chan [][]string, conc)
	startReadingCSV := time.Now()
	for _, fileName := range csvFileNames {
		go func(name string) {
			loadCSV, _ := helpers.LoadCSV(name)
			csvLines <- loadCSV
		}(fileName)
	}
	var output [][]string
	for i := 0; i < conc; i++ {
		out := <-csvLines
		output = append(output, out...)
	}
	timeElapsedReadingCSV := time.Since(startReadingCSV)
	fmt.Printf("Time reading CSV of length %v  of %v documents took \n", len(output), timeElapsedReadingCSV)
	var dataToSeed []interface{}
	err = DecodeCSV(output, &dataToSeed)
	if err != nil {
		fmt.Println("Inside decode error")
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGO_URI")))
	if err != nil {
		fmt.Println("Inside exit of Mongo")
		log.Fatal(err)
	}

	//Verify MongoDB connection
	if err = client.Ping(context.TODO(), readpref.Primary()); err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to MongoDB")
	database := client.Database(os.Getenv("MONGO_DATABASE"))
	database.Collection(constants.ProductsCollection).Drop(context.TODO())
	opts := options.InsertMany().SetOrdered(false)
	start := time.Now()
	res, err := database.Collection(constants.ProductsCollection).InsertMany(context.TODO(), dataToSeed, opts)
	if err != nil {
		log.Fatal(err)
	}
	elapsed := time.Since(start)
	fmt.Printf("Inserted  of %v documents took %s\n", len(res.InsertedIDs), elapsed)
}

// DecodeCSV is used to decode the csvlines into an array of Products struct
func DecodeCSV(csvLines [][]string, dataToDecode *[]interface{}) error {
	for _, line := range csvLines {
		stockChange, err := strconv.ParseInt(line[3], 10, 64)
		if err != nil {
			return err
		}
		product := domain.Products{
			Country:     line[0],
			SKU:         line[1],
			Name:        line[2],
			StockChange: stockChange,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
		*dataToDecode = append(*dataToDecode, product)
	}
	return nil
}
