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

// Create godoc
// @Summary Tambah meja baru
// @Description Menambahkan data meja baru ke dalam sistem
// @Tags Table
// @Accept json
// @Produce json
// @Param request body newTableRequest true "Data meja"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /tables [post]
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

// List godoc
// @Summary Lihat semua data meja
// @Tags Table
// @Produce json
// @Success 200 {array} models.Table
// @Failure 500 {object} map[string]string
// @Router /tables [get]
func (h *TableHandler) List(c *gin.Context) {
	tables, err := h.service.ListTables(c.Request.Context())
	if err != nil {
		log.Println("Failed to get tables:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data meja"})
		return
	}
	c.JSON(http.StatusOK, tables)
}

// GetByID godoc
// @Summary Lihat detail meja berdasarkan ID
// @Tags Table
// @Produce json
// @Param id path int true "ID meja"
// @Success 200 {object} models.Table
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /tables/{id} [get]
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

// Update godoc
// @Summary Update data meja
// @Tags Table
// @Accept json
// @Produce json
// @Param id path int true "ID meja"
// @Param request body newTableRequest true "Data meja yang diperbarui"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /tables/{id} [put]
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

// Delete godoc
// @Summary Hapus (soft delete) data meja
// @Tags Table
// @Produce json
// @Param id path int true "ID meja"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /tables/{id} [delete]
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
