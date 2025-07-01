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

type CustomerHandler struct {
	service *services.CustomerService
}

func NewCustomerHandler(service *services.CustomerService) *CustomerHandler {
	return &CustomerHandler{service: service}
}

type CustomerRequest struct {
	HotelGuestID string  `json:"hotel_guest_id"`
	Type         string  `json:"type" binding:"required"`
	Name         string  `json:"name" binding:"required"`
	Phone        string  `json:"phone"`
	VisitCount   int     `json:"visit_count"`
	LastVisit    *string `json:"last_visit"` // ISO string expected
}

func (h *CustomerHandler) Create(c *gin.Context) {
	var req CustomerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Invalid request: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	customer := &models.Customer{
		HotelGuestID: sql.NullString{String: req.HotelGuestID, Valid: req.HotelGuestID != ""},
		Type:         req.Type,
		Name:         req.Name,
		Phone:        sql.NullString{String: req.Phone, Valid: req.Phone != ""},
		VisitCount:   req.VisitCount,
	}

	id, err := h.service.CreateCustomer(c.Request.Context(), customer)
	if err != nil {
		log.Printf("Failed to create customer: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat customer"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": id, "name": customer.Name})
}

func (h *CustomerHandler) List(c *gin.Context) {
	data, err := h.service.GetAllCustomers(c.Request.Context())
	if err != nil {
		log.Printf("Failed to get customers: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil customer"})
		return
	}
	c.JSON(http.StatusOK, data)
}

func (h *CustomerHandler) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Printf("Invalid customer ID: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	customer, err := h.service.GetCustomerByID(c.Request.Context(), id)
	if err != nil {
		log.Printf("Failed to get customer by ID: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data"})
		return
	}

	c.JSON(http.StatusOK, customer)
}

func (h *CustomerHandler) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Printf("Invalid customer ID: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	var req CustomerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Invalid update request: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	customer := &models.Customer{
		CustID:       id,
		HotelGuestID: sql.NullString{String: req.HotelGuestID, Valid: req.HotelGuestID != ""},
		Type:         req.Type,
		Name:         req.Name,
		Phone:        sql.NullString{String: req.Phone, Valid: req.Phone != ""},
		VisitCount:   req.VisitCount,
	}

	if err := h.service.UpdateCustomer(c.Request.Context(), customer); err != nil {
		log.Printf("Failed to update customer: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui customer"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Customer berhasil diperbarui"})
}

func (h *CustomerHandler) SoftDelete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Printf("Invalid customer ID: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	if err := h.service.SoftDeleteCustomer(c.Request.Context(), id); err != nil {
		log.Printf("Failed to delete customer: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus customer"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Customer berhasil dihapus"})
}
