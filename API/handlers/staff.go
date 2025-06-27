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

func (h *StaffHandler) List(c *gin.Context) {
	staff, err := h.service.ListStaff(c.Request.Context())
	if err != nil {
		log.Printf("Gagal ambil staff: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data staff"})
		return
	}
	c.JSON(http.StatusOK, staff)
}

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
