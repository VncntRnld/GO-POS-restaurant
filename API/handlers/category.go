package handlers

import (
	"log"
	"net/http"
	"pos-restaurant/models"
	"pos-restaurant/services"
	"strconv"

	"github.com/gin-gonic/gin"
)

type MenuCategoryHandler struct {
	service *services.MenuCategoryService
}

func NewMenuCategoryHandler(service *services.MenuCategoryService) *MenuCategoryHandler {
	return &MenuCategoryHandler{service: service}
}

type CreateCategoryRequest struct {
	Name string `json:"name" binding:"required"`
}

// CreateCategory godoc
// @Summary Tambah kategori menu baru
// @Tags Category
// @Accept json
// @Produce json
// @Param request body CreateCategoryRequest true "Data kategori baru"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /menu/category [post]
func (h *MenuCategoryHandler) CreateCategory(c *gin.Context) {
	var req CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	category := &models.MenuCategory{Name: req.Name}
	id, err := h.service.CreateCategory(c.Request.Context(), category)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat kategori"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": id, "name": category.Name})
}

// ListCategories godoc
// @Summary Tampilkan semua kategori menu
// @Tags Category
// @Produce json
// @Success 200 {array} models.MenuCategory
// @Failure 500 {object} map[string]string
// @Router /menu/category [get]
func (h *MenuCategoryHandler) ListCategories(c *gin.Context) {
	categories, err := h.service.ListCategories(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil kategori"})
		return
	}
	c.JSON(http.StatusOK, categories)
}

// DeleteCategory godoc
// @Summary Hapus (soft delete) kategori menu
// @Tags Category
// @Produce json
// @Param id path int true "ID Kategori"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /menu/category/{id} [delete]
func (h *MenuCategoryHandler) DeleteCategory(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("ID kategori tidak valid: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	err = h.service.DeleteCategory(c.Request.Context(), id)
	if err != nil {
		log.Printf("Gagal soft delete kategori: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus kategori"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Kategori berhasil dihapus (soft delete)"})
}
