package http

import (
	"beaver/inventory/core/domain"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type StockHandler struct {
	stockService *domain.StockService
}

func NewV1Handler(router *gin.RouterGroup, stockService *domain.StockService) {
	handler := &StockHandler{stockService: stockService}
	router.POST("/stock/change", handler.stockChange)
	router.GET("/products", handler.GetAllProducts)
}

type StockChangeTO struct {
	Batch    BatchTO    `json:"batch"`
	Location LocationTO `json:"location"`
	Quantity QuantityTO `json:"quantity"`
}

type QuantityTO struct {
	Amount float64     `json:"amount"`
	Unit   domain.Unit `json:"unit"`
}

type LocationTO struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type BatchTO struct {
	Id           string    `json:"id"`
	Product      ProductTO `json:"product"`
	ProducedAt   time.Time `json:"producedAt"`
	SellLatestAt time.Time `json:"sellLatestAt"`
}

type ProductTO struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

func (h *StockHandler) GetAllProducts(c *gin.Context) {
	c.JSON(http.StatusOK, h.stockService.GetAllProducts())
}

func (h *StockHandler) stockChange(c *gin.Context) {
	var stockChange StockChangeTO
	if err := c.ShouldBindJSON(&stockChange); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.stockService.StockChange(domain.StockChangeEvent{
		Batch: domain.Batch{
			Id:           stockChange.Batch.Id,
			Product:      domain.Product(stockChange.Batch.Product),
			ProducedAt:   stockChange.Batch.ProducedAt,
			SellLatestAt: stockChange.Batch.SellLatestAt,
		},
		Location: domain.Location(stockChange.Location),
		Quantity: domain.Quantity(stockChange.Quantity),
	}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}
