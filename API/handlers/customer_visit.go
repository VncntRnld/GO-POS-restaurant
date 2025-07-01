package handlers

import (
	"database/sql"
	"log"
	"net/http"
	"pos-restaurant/models"
	"pos-restaurant/services"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type CustomerVisitHandler struct {
	service *services.CustomerVisitService
}

func NewCustomerVisitHandler(service *services.CustomerVisitService) *CustomerVisitHandler {
	return &CustomerVisitHandler{service: service}
}

type CustomerVisitRequest struct {
	CustomerID    int       `json:"customer_id"`
	VisitType     string    `json:"visit_type"`
	VisitDate     time.Time `json:"visit_date"`
	RoomNumber    string    `json:"room_number"`
	ReservationID int64     `json:"reservation_id"`
	OutletID      int       `json:"outlet_id"`
	TotalSpent    float64   `json:"total_spent"`
	Pax           int       `json:"pax"`
}

func (h *CustomerVisitHandler) Create(c *gin.Context) {
	var req CustomerVisitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Bind error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	visit := &models.CustomerVisit{
		CustomerID:    req.CustomerID,
		VisitType:     req.VisitType,
		VisitDate:     req.VisitDate,
		RoomNumber:    sql.NullString{String: req.RoomNumber, Valid: req.RoomNumber != ""},
		ReservationID: sql.NullInt64{Int64: req.ReservationID, Valid: req.ReservationID != 0},
		OutletID:      req.OutletID,
		TotalSpent:    sql.NullFloat64{Float64: req.TotalSpent, Valid: req.TotalSpent != 0},
		Pax:           req.Pax,
	}

	id, err := h.service.Create(c.Request.Context(), visit)
	if err != nil {
		log.Printf("Create error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menambahkan data"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": id})
}

func (h *CustomerVisitHandler) List(c *gin.Context) {
	data, err := h.service.List(c.Request.Context())
	if err != nil {
		log.Printf("List error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data"})
		return
	}
	c.JSON(http.StatusOK, data)
}

func (h *CustomerVisitHandler) GetByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	data, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		log.Printf("GetByID error: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Data tidak ditemukan"})
		return
	}
	c.JSON(http.StatusOK, data)
}

func (h *CustomerVisitHandler) Update(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var req CustomerVisitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Bind error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	visit := &models.CustomerVisit{
		ID:            id,
		CustomerID:    req.CustomerID,
		VisitType:     req.VisitType,
		VisitDate:     req.VisitDate,
		RoomNumber:    sql.NullString{String: req.RoomNumber, Valid: req.RoomNumber != ""},
		ReservationID: sql.NullInt64{Int64: req.ReservationID},
		OutletID:      req.OutletID,
		TotalSpent:    sql.NullFloat64{Float64: req.TotalSpent},
		Pax:           req.Pax,
	}

	if err := h.service.Update(c.Request.Context(), visit); err != nil {
		log.Printf("Update error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal update data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Data berhasil diupdate"})
}

func (h *CustomerVisitHandler) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		log.Printf("Delete error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus data"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Data berhasil dihapus"})
}
