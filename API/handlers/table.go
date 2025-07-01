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

type TableHandler struct {
	service *services.TableService
}

func NewTableHandler(service *services.TableService) *TableHandler {
	return &TableHandler{service: service}
}

type newTableRequest struct {
	OutletID     int    `json:"outlet_id"`
	TableNumber  string `json:"table_number"`
	Capacity     int    `json:"capacity"`
	LocationType string `json:"location_type"`
	Status       string `json:"status"`
}

func (h *TableHandler) Create(c *gin.Context) {
	var req newTableRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println("Invalid request:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	table := &models.Table{
		OutletID:     req.OutletID,
		TableNumber:  req.TableNumber,
		Capacity:     req.Capacity,
		LocationType: sql.NullString{String: req.LocationType, Valid: req.LocationType != ""},
		Status:       req.Status,
	}

	id, err := h.service.CreateTable(c.Request.Context(), table)
	if err != nil {
		log.Println("Failed to create table:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat meja"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": id, "message": "Table berhasil ditambahkan"})
}

func (h *TableHandler) List(c *gin.Context) {
	tables, err := h.service.ListTables(c.Request.Context())
	if err != nil {
		log.Println("Failed to get tables:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data meja"})
		return
	}
	c.JSON(http.StatusOK, tables)
}

func (h *TableHandler) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Println("Invalid ID:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}
	table, err := h.service.GetTableByID(c.Request.Context(), id)
	if err != nil {
		log.Println("Failed to get table by ID:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Meja tidak ditemukan"})
		return
	}
	c.JSON(http.StatusOK, table)
}

func (h *TableHandler) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Println("Invalid ID:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}
	var req newTableRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println("Invalid update payload:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	table := &models.Table{
		ID:           id,
		OutletID:     req.OutletID,
		TableNumber:  req.TableNumber,
		Capacity:     req.Capacity,
		LocationType: sql.NullString{String: req.LocationType, Valid: req.LocationType != ""},
		Status:       req.Status,
	}

	if err := h.service.UpdateTable(c.Request.Context(), table); err != nil {
		log.Println("Failed to update table:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal update meja"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Meja berhasil diupdate"})
}

func (h *TableHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Println("Invalid ID:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}
	if err := h.service.SoftDeleteTable(c.Request.Context(), id); err != nil {
		log.Println("Failed to delete table:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal hapus meja"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Meja berhasil dihapus"})
}
