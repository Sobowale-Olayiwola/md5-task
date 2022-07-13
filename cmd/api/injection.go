package main

import (
	"jumia/cmd/socket"
	"jumia/internal/constants"
	_mongoProductsRepo "jumia/products/repository/mongo"
	_mongoProductQueries "jumia/products/repository/queries"
	_redisProductsRepo "jumia/products/repository/redis"
	"net/http"

	_productService "jumia/products/service"

	_productsHandler "jumia/products/handler/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	socketio "github.com/googollee/go-socket.io"
)

func inject(d *DataSources) (*gin.Engine, *socketio.Server) {
	SocketIOServer := socket.NewSocketIO().StartSocketIOServer()
	/*
	 * repository layer
	 */
	mongoProductRepo := _mongoProductsRepo.NewMongoProductRepository(d.DB, d.Client, constants.ProductsCollection)
	redisProductsRepo := _redisProductsRepo.NewRedisInMemoryDB(d.InMemoryDB)

	/*
	 * service layer
	 */
	productService := _productService.NewProductService(mongoProductRepo, redisProductsRepo, _mongoProductQueries.MongoQuery{})

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.Use(cors.Default())
	router.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"message": "Welcome to jumia"})
	})
	router.GET("/socket.io/*any", gin.WrapH(SocketIOServer))
	router.POST("/socket.io/*any", gin.WrapH(SocketIOServer))
	router.Use(static.Serve("/web", static.LocalFile("../web/public", true)))
	/*
	 * handler layer
	 */
	_productsHandler.NewProductHandler(router, productService, SocketIOServer)
	return router, SocketIOServer
}
