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

type OutletHandler struct {
	service *services.OutletService
}

func NewOutletHandler(service *services.OutletService) *OutletHandler {
	return &OutletHandler{service: service}
}

type CreateOutletRequest struct {
	Name                 string  `json:"name" binding:"required"`
	Location             string  `json:"location"`
	ServiceChargePercent float64 `json:"service_charge_percentage"`
	TaxPercentage        float64 `json:"tax_percentage"`
	IsActive             bool    `json:"is_active"`
}

// Create godoc
// @Summary Tambah outlet baru
// @Tags Outlet
// @Accept json
// @Produce json
// @Param request body CreateOutletRequest true "Data outlet"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /outlets [post]
func (h *OutletHandler) Create(c *gin.Context) {
	var req CreateOutletRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Validation error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	outlet := &models.Outlet{
		Name:                 req.Name,
		Location:             sql.NullString{String: req.Location, Valid: req.Location != ""},
		ServiceChargePercent: req.ServiceChargePercent,
		TaxPercentage:        req.TaxPercentage,
		IsActive:             req.IsActive,
	}

	id, err := h.service.CreateOutlet(c.Request.Context(), outlet)
	if err != nil {
		log.Printf("DB error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat outlet"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": id, "message": "Outlet berhasil ditambahkan"})
}

// List godoc
// @Summary Tampilkan semua outlet
// @Tags Outlet
// @Produce json
// @Success 200 {array} models.Outlet
// @Failure 500 {object} map[string]string
// @Router /outlets [get]
func (h *OutletHandler) List(c *gin.Context) {
	outlets, err := h.service.ListOutlets(c.Request.Context())
	if err != nil {
		log.Printf("DB error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil outlet"})
		return
	}
	c.JSON(http.StatusOK, outlets)
}

// GetByID godoc
// @Summary Ambil data outlet berdasarkan ID
// @Tags Outlet
// @Produce json
// @Param id path int true "ID outlet"
// @Success 200 {object} models.Outlet
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /outlets/{id} [get]
func (h *OutletHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	outlet, err := h.service.GetOutletByID(c.Request.Context(), id)
	if err != nil {
		log.Printf("DB error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil outlet"})
		return
	}
	if outlet == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Outlet tidak ditemukan"})
		return
	}
	c.JSON(http.StatusOK, outlet)
}

// Update godoc
// @Summary Perbarui outlet
// @Tags Outlet
// @Accept json
// @Produce json
// @Param id path int true "ID outlet"
// @Param request body CreateOutletRequest true "Data outlet yang diperbarui"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /outlets/{id} [put]
func (h *OutletHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	var req CreateOutletRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	outlet := &models.Outlet{
		ID:                   id,
		Name:                 req.Name,
		Location:             sql.NullString{String: req.Location, Valid: req.Location != ""},
		ServiceChargePercent: req.ServiceChargePercent,
		TaxPercentage:        req.TaxPercentage,
		IsActive:             req.IsActive,
	}

	if err := h.service.UpdateOutlet(c.Request.Context(), outlet); err != nil {
		log.Printf("DB error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui outlet"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Berhasil update outlet"})
}

// Delete godoc
// @Summary Hapus outlet (soft delete)
// @Tags Outlet
// @Produce json
// @Param id path int true "ID outlet"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /outlets/{id} [delete]
func (h *OutletHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	if err := h.service.SoftDeleteOutlet(c.Request.Context(), id); err != nil {
		log.Printf("DB error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus outlet"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Outlet berhasil dihapus (soft delete)"})
}
