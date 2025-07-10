package handlers

import (
	"log"
	"net/http"
	"pos-restaurant/models"
	"pos-restaurant/services"
	"strconv"

	"github.com/gin-gonic/gin"
)

type MenuIngredientHandler struct {
	service *services.MenuIngredientService
}

func NewMenuIngredientHandler(service *services.MenuIngredientService) *MenuIngredientHandler {
	return &MenuIngredientHandler{service: service}
}

type CreateMenuIngredientRequest struct {
	MenuItemID   int     `json:"menu_item_id" binding:"required"`
	IngredientID int     `json:"ingredient_id" binding:"required"`
	Qty          float64 `json:"qty" binding:"required"`
	IsRemovable  bool    `json:"is_removable"`
	IsDefault    bool    `json:"is_default"`
}

// Create godoc
// @Summary Tambah menu ingredient baru
// @Tags MenuIngredient
// @Accept json
// @Produce json
// @Param request body CreateMenuIngredientRequest true "Data menu ingredient"
// @Success 201 {object} map[string]int
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /menu-ingredients [post]
func (h *MenuIngredientHandler) Create(c *gin.Context) {
	var req CreateMenuIngredientRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Error parsing request: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	m := &models.MenuIngredient{
		MenuItemID:   req.MenuItemID,
		IngredientID: req.IngredientID,
		Qty:          req.Qty,
		IsRemovable:  req.IsRemovable,
		IsDefault:    req.IsDefault,
	}

	id, err := h.service.CreateMenuIngredient(c.Request.Context(), m)
	if err != nil {
		log.Printf("Failed to create menu ingredient: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menambahkan ingredient ke menu"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": id})
}

// ListByMenuItem godoc
// @Summary Tampilkan semua ingredient berdasarkan menu item
// @Tags MenuIngredient
// @Produce json
// @Param menu_item_id path int true "ID Menu Item"
// @Success 200 {array} models.MenuIngredient
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /menu-ingredients/{menu_item_id} [get]
func (h *MenuIngredientHandler) ListByMenuItem(c *gin.Context) {
	idStr := c.Param("menu_item_id")
	menuItemID, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("Invalid menu_item_id: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	items, err := h.service.GetIngredientsByMenuItem(c.Request.Context(), menuItemID)
	if err != nil {
		log.Printf("Failed to get menu ingredients: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data ingredients untuk menu"})
		return
	}

	c.JSON(http.StatusOK, items)
}

// UpdateMenuIngredient godoc
// @Summary Perbarui data menu ingredient
// @Tags MenuIngredient
// @Accept json
// @Produce json
// @Param id path int true "ID menu ingredient"
// @Param request body CreateMenuIngredientRequest true "Data yang diperbarui"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /menu-ingredients/{id} [put]
func (h *MenuIngredientHandler) UpdateMenuIngredient(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("Invalid menu_ingredient ID: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	var req CreateMenuIngredientRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Error binding JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	m := &models.MenuIngredient{
		ID:           id,
		MenuItemID:   req.MenuItemID,
		IngredientID: req.IngredientID,
		Qty:          req.Qty,
		IsRemovable:  req.IsRemovable,
		IsDefault:    req.IsDefault,
	}

	if err := h.service.UpdateMenuIngredient(c.Request.Context(), m); err != nil {
		log.Printf("Failed to update menu ingredient: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui menu ingredient"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Menu ingredient berhasil diperbarui"})
}

// DeleteMenuIngredient godoc
// @Summary Hapus menu ingredient berdasarkan ID
// @Tags MenuIngredient
// @Produce json
// @Param id path int true "ID menu ingredient"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /menu-ingredients/{id} [delete]
func (h *MenuIngredientHandler) DeleteMenuIngredient(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("ID menu_ingredient tidak valid untuk delete: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	err = h.service.DeleteMenuIngredient(c.Request.Context(), id)
	if err != nil {
		log.Printf("Gagal menghapus Menu_ingredient %d: %v", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus menu ingredient"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Menu ingredient berhasil dihapus"})
}
