package handlers

import (
	"database/sql"
	"log"
	"net/http"
	"pos-restaurant/models"
	"pos-restaurant/services"
	"strconv"

	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	service *services.OrderService
}

func NewOrderHandler(service *services.OrderService) *OrderHandler {
	return &OrderHandler{service: service}
}

type NewOrderRequest struct {
	ID         int                     `json:"id"`
	TableID    int                     `json:"table_id"`
	CustomerID int                     `json:"customer_id"`
	HotelRoom  string                  `json:"hotel_room"`
	WaiterID   int                     `json:"waiter_id"`
	OutletID   int                     `json:"outlet_id"`
	Status     string                  `json:"status"`
	OrderType  string                  `json:"order_type"`
	Items      []models.OrderItemInput `json:"items"`
}

func (h *OrderHandler) Create(c *gin.Context) {
	var req NewOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Bind JSON error (Create Order): %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	order := &models.OrderRequest{
		TableID:    req.TableID,
		CustomerID: req.CustomerID,
		HotelRoom:  sql.NullString{String: req.HotelRoom, Valid: req.HotelRoom != ""},
		WaiterID:   req.WaiterID,
		OutletID:   req.OutletID,
		Status:     req.Status,
		OrderType:  req.OrderType,
		Items:      req.Items,
	}

	id, err := h.service.Create(c.Request.Context(), order)
	if err != nil {
		log.Printf("Create Order error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat order"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": id})
}

func (h *OrderHandler) List(c *gin.Context) {
	orders, err := h.service.List(c.Request.Context())
	if err != nil {
		log.Printf("List Orders error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get orders"})
		return
	}
	c.JSON(http.StatusOK, orders)
}

func (h *OrderHandler) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Printf("Invalid ID (GetByID): %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	order, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		log.Printf("Order not found (ID %d): %v", id, err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	c.JSON(http.StatusOK, order)
}

func (h *OrderHandler) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Printf("Invalid ID (Update): %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var req NewOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Bind JSON error (Update Order): %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req.ID = id
	order := &models.Order{
		ID:         req.ID,
		TableID:    req.TableID,
		CustomerID: req.CustomerID,
		HotelRoom:  sql.NullString{String: req.HotelRoom, Valid: req.HotelRoom != ""},
		WaiterID:   req.WaiterID,
		OutletID:   req.OutletID,
		Status:     req.Status,
		OrderType:  req.OrderType,
	}

	if err := h.service.Update(c.Request.Context(), order); err != nil {
		log.Printf("Update Order error (ID %d): %v", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update order"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order updated"})
}

func (h *OrderHandler) AddItem(c *gin.Context) {
	orderID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Printf("Invalid Order ID (AddItem): %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}

	var req models.AddOrderItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Bind JSON error (AddItem): %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.AddItem(c.Request.Context(), orderID, &req); err != nil {
		log.Printf("Add Item to Order error (Order ID %d): %v", orderID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menambahkan item"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Item berhasil ditambahkan ke order"})
}

func (h *OrderHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Printf("Invalid ID (Delete): %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := h.service.SoftDelete(c.Request.Context(), id); err != nil {
		log.Printf("Delete Order error (ID %d): %v", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete order"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order deleted"})
}
