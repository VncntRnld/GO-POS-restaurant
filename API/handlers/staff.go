package handlers

import (
	"log"
	"net/http"
	"pos-restaurant/models"
	"pos-restaurant/services"
	"strconv"

	"github.com/gin-gonic/gin"
)

type StaffHandler struct {
	service *services.StaffService
}

func NewStaffHandler(service *services.StaffService) *StaffHandler {
	return &StaffHandler{service: service}
}

type StaffRequest struct {
	Name     string `json:"name" binding:"required"`
	Role     string `json:"role" binding:"required"`
	PinCode  string `json:"pin_code" binding:"required"`
	IsActive bool   `json:"is_active"`
}

// Create godoc
// @Summary Tambah staff baru
// @Tags Staff
// @Accept json
// @Produce json
// @Param request body StaffRequest true "Data staff"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /staff [post]
func (h *StaffHandler) Create(c *gin.Context) {
	var req StaffRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	staff := &models.Staff{
		Name:     req.Name,
		Role:     req.Role,
		PinCode:  req.PinCode,
		IsActive: req.IsActive,
	}
	id, err := h.service.CreateStaff(c.Request.Context(), staff)
	if err != nil {
		log.Printf("Gagal tambah staff: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat staff"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": id, "message": "Staff berhasil ditambahkan"})
}

// List godoc
// @Summary Lihat semua staff
// @Tags Staff
// @Produce json
// @Success 200 {array} models.Staff
// @Failure 500 {object} map[string]string
// @Router /staff [get]
func (h *StaffHandler) List(c *gin.Context) {
	staff, err := h.service.ListStaff(c.Request.Context())
	if err != nil {
		log.Printf("Gagal ambil staff: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data staff"})
		return
	}
	c.JSON(http.StatusOK, staff)
}

// GetByID godoc
// @Summary Lihat staff berdasarkan ID
// @Tags Staff
// @Produce json
// @Param id path int true "ID staff"
// @Success 200 {object} models.Staff
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /staff/{id} [get]
func (h *StaffHandler) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}
	staff, err := h.service.GetStaffByID(c.Request.Context(), id)
	if err != nil {
		log.Printf("Gagal ambil staff id %d: %v", id, err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Staff tidak ditemukan"})
		return
	}
	c.JSON(http.StatusOK, staff)
}

// Update godoc
// @Summary Perbarui data staff
// @Tags Staff
// @Accept json
// @Produce json
// @Param id path int true "ID staff"
// @Param request body StaffRequest true "Data staff yang diperbarui"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /staff/{id} [put]
func (h *StaffHandler) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}
	var req StaffRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	staff := &models.Staff{
		ID:       id,
		Name:     req.Name,
		Role:     req.Role,
		PinCode:  req.PinCode,
		IsActive: req.IsActive,
	}
	if err := h.service.UpdateStaff(c.Request.Context(), staff); err != nil {
		log.Printf("Gagal update staff: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui staff"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Staff berhasil diupdate"})
}

// SoftDelete godoc
// @Summary Hapus (soft delete) staff
// @Tags Staff
// @Produce json
// @Param id path int true "ID staff"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /staff/{id} [delete]
func (h *StaffHandler) SoftDelete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}
	if err := h.service.SoftDeleteStaff(c.Request.Context(), id); err != nil {
		log.Printf("Gagal hapus staff: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus staff"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Staff berhasil dihapus"})
}
