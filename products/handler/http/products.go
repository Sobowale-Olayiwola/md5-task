package http

import (
	"context"
	"errors"
	"fmt"
	"jumia/domain"
	"jumia/internal/helpers"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	socketio "github.com/googollee/go-socket.io"
)

type ProductHandler struct {
	ProductService domain.ProductService
	Socket         *socketio.Server
}

func NewProductHandler(router *gin.Engine, j domain.ProductService, s *socketio.Server) {
	handler := &ProductHandler{
		ProductService: j,
		Socket:         s,
	}
	api := router.Group("/api/v1")
	api.GET("/products/sku/:sku", handler.GetProductBySKU)
	api.PUT("/products/stocks", handler.ConsumeProductStock)
	api.PUT("/products/csv/uploads", handler.BulkUpdateWithCSV)
}

func (j *ProductHandler) GetProductBySKU(c *gin.Context) {
	sku, ok := c.Params.Get("sku")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "sku is required"})
	}
	ctx := context.TODO()
	products, err := j.ProductService.GetProductBySKU(ctx, sku)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrRecordNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}
	c.JSON(http.StatusFound, gin.H{"success": true, "payload": products, "message": "Successfully found product"})
}

func (j *ProductHandler) ConsumeProductStock(c *gin.Context) {
	sku := c.Query("sku")
	if sku == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "sku is required"})
	}
	amount, err := strconv.ParseInt(c.Query("amount"), 10, 64)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}
	ctx := context.TODO()
	err = j.ProductService.ConsumeProductStock(ctx, sku, amount)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrRecordNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Successfully consumed stock"})
}

func (j *ProductHandler) BulkUpdateWithCSV(c *gin.Context) {
	fileName := c.Query("csv")
	if fileName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "file name is required"})
		return
	}
	csvLines, err := helpers.LoadCSV(fmt.Sprintf("/home/layitheinfotechguru/mds_challenge/internal/csv/%s.csv", fileName))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "csv file not found"})
		return
	}
	ctx := context.TODO()
	helpers.Background(func() {
		err = j.ProductService.BulkUpdateWithCSV(ctx, csvLines)
		if err != nil {
			switch {
			case errors.Is(err, domain.ErrRecordNotFound):
				j.Socket.BroadcastToNamespace("/csv", "bulk-update-status", err.Error())
				return
			default:
				j.Socket.BroadcastToNamespace("/csv", "bulk-update-status", "Internal Server: "+err.Error())
				return
			}
		} else {
			fmt.Println("Bulk upload succeeded")
			j.Socket.BroadcastToNamespace("/csv", "bulk-update-status", "Bulk update successfully")
			return
		}
	})
	c.JSON(http.StatusAccepted, gin.H{"message": "Acknowledged"})
}
