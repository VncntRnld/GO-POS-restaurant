package handlers

import (
	"database/sql"
	"log"
	"net/http"
	"pos-restaurant/models"
	"pos-restaurant/services"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type ReservationHandler struct {
	service *services.ReservationService
}

func NewReservationHandler(service *services.ReservationService) *ReservationHandler {
	return &ReservationHandler{service: service}
}

type CreateReservationRequest struct {
	CustomerID      int    `json:"customer_id" binding:"required"`
	ReservationTime string `json:"reservation_time" binding:"required"`
	Pax             int    `json:"pax" binding:"required"`
	TableID         int    `json:"table_id" binding:"required"`
	Status          string `json:"status" binding:"required"`
	SpecialRequest  string `json:"special_request"`
}

func (h *ReservationHandler) Create(c *gin.Context) {
	var req CreateReservationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println("Create reservation bind error:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	timeParsed, err := time.Parse(time.RFC3339, req.ReservationTime)
	if err != nil {
		log.Println("Invalid time format:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format waktu tidak valid (harus RFC3339)"})
		return
	}

	res := &models.Reservation{
		CustomerID:      req.CustomerID,
		ReservationTime: timeParsed,
		Pax:             req.Pax,
		TableID:         req.TableID,
		Status:          req.Status,
		SpecialRequest:  sql.NullString{String: req.SpecialRequest, Valid: req.SpecialRequest != ""},
	}

	id, err := h.service.Create(c.Request.Context(), res)
	if err != nil {

		if strings.Contains(err.Error(), "meja sudah dipesan") {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}

		log.Println("Gagal create reservation:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat reservasi"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": id})
}

func (h *ReservationHandler) List(c *gin.Context) {
	sortBy := c.DefaultQuery("sort", "reservation_time")

	reservations, err := h.service.List(c.Request.Context(), sortBy)
	if err != nil {
		log.Printf("Gagal mengambil data reservasi: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data reservasi"})
		return
	}

	c.JSON(http.StatusOK, reservations)
}

func (h *ReservationHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Println("Invalid ID:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	res, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		log.Println("Gagal ambil reservation by ID:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Reservasi tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *ReservationHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Println("Invalid ID:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	var req CreateReservationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println("Bind error update:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	timeParsed, err := time.Parse(time.RFC3339, req.ReservationTime)
	if err != nil {
		log.Println("Invalid time format:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format waktu tidak valid (harus RFC3339)"})
		return
	}

	res := &models.Reservation{
		ID:              id,
		CustomerID:      req.CustomerID,
		ReservationTime: timeParsed,
		Pax:             req.Pax,
		TableID:         req.TableID,
		Status:          req.Status,
		SpecialRequest:  sql.NullString{String: req.SpecialRequest, Valid: req.SpecialRequest != ""},
	}

	if err := h.service.Update(c.Request.Context(), res); err != nil {
		log.Println("Gagal update reservation:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal update reservasi"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Reservasi berhasil diupdate"})
}

func (h *ReservationHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Println("Invalid ID:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		log.Println("Gagal delete reservation:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus reservasi"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Reservasi berhasil dihapus"})
}
