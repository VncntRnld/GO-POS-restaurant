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

type MenuItemHandler struct {
	service *services.MenuService
}

func NewMenuItemHandler(service *services.MenuService) *MenuItemHandler {
	return &MenuItemHandler{service: service}
}

type CreateMenuItemRequest struct {
	CategoryID      int      `json:"category_id" binding:"required"`
	SKU             string   `json:"sku" binding:"required,max=50"`
	Name            string   `json:"name" binding:"required,max=255"`
	Description     string   `json:"description"`
	Price           float64  `json:"price" binding:"required,gt=0"`
	Cost            float64  `json:"cost" binding:"gte=0"`
	IsActive        bool     `json:"is_active"`
	PreparationTime *int     `json:"preparation_time"`
	Tags            []string `json:"tags"`
}

// CreateMenuItem godoc
// @Summary Tambah menu baru
// @Tags Menu-Items
// @Accept json
// @Produce json
// @Param request body CreateMenuItemRequest true "Data menu baru"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /menu/menu-items [post]
func (h *MenuItemHandler) CreateMenuItem(c *gin.Context) {
	var req CreateMenuItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("[CreateMenuItem] Bad request: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	menuItem := &models.MenuItem{
		CategoryID:  req.CategoryID,
		SKU:         req.SKU,
		Name:        req.Name,
		Description: sql.NullString{String: req.Description, Valid: req.Description != ""},
		Price:       req.Price,
		Cost:        req.Cost,
		IsActive:    req.IsActive,
		Tags:        req.Tags,
	}
	if req.PreparationTime != nil {
		menuItem.PreparationTime = sql.NullInt64{Int64: int64(*req.PreparationTime), Valid: true}
	}

	id, err := h.service.CreateMenuItem(c.Request.Context(), menuItem)
	if err != nil {
		log.Printf("[CreateMenuItem] Gagal membuat menu: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat menu"})
		return
	}

	log.Printf("[CreateMenuItem] Menu item berhasil dibuat dengan ID: %d", id)
	c.JSON(http.StatusCreated, gin.H{"id": id, "name": menuItem.Name})
}

// GetMenuItemsByCategory godoc
// @Summary Dapatkan menu berdasarkan kategori
// @Tags Menu
// @Produce json
// @Param category_id query int true "ID Kategori"
// @Success 200 {array} models.MenuItem
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /menu/menu-items/category [get]
func (h *MenuItemHandler) GetMenuItemsByCategory(c *gin.Context) {
	categoryIDStr := c.Query("category_id")
	categoryID, err := strconv.Atoi(categoryIDStr)
	if err != nil {
		log.Printf("Parameter category_id tidak valid: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parameter category_id tidak valid"})
		return
	}

	items, err := h.service.GetMenuItemsByCategory(c.Request.Context(), categoryID)
	if err != nil {
		log.Printf("Gagal mengambil menu berdasarkan kategori: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data menu"})
		return
	}

	c.JSON(http.StatusOK, items)
}

// ListMenuItems godoc
// @Summary List semua menu
// @Tags Menu
// @Produce json
// @Success 200 {array} models.MenuItem
// @Failure 500 {object} map[string]string
// @Router /menu/menu-items [get]
func (h *MenuItemHandler) ListMenuItems(c *gin.Context) {
	items, err := h.service.ListMenuItems(c.Request.Context())
	if err != nil {
		log.Printf("[ListMenuItems] Gagal mengambil data menu: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data menu"})
		return
	}
	c.JSON(http.StatusOK, items)
}

// ListActiveMenuItems godoc
// @Summary List menu aktif
// @Tags Menu
// @Produce json
// @Success 200 {array} models.MenuItem
// @Failure 500 {object} map[string]string
// @Router /menu/menu-items-active [get]
func (h *MenuItemHandler) ListActiveMenuItems(c *gin.Context) {
	items, err := h.service.ListActiveMenuItems(c.Request.Context())
	if err != nil {
		log.Printf("[ListActiveMenuItems] Gagal mengambil data aktif: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data aktif"})
		return
	}
	c.JSON(http.StatusOK, items)
}

// SearchMenuItems godoc
// @Summary Cari menu berdasarkan keyword
// @Tags Menu
// @Produce json
// @Param search query string true "Kata kunci pencarian"
// @Success 200 {array} models.MenuItem
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /menu/menu-items/search [get]
func (h *MenuItemHandler) SearchMenuItems(c *gin.Context) {
	keyword := c.Query("search")
	if keyword == "" {
		log.Println("[SearchMenuItems] parameter 'search' kosong")
		c.JSON(http.StatusBadRequest, gin.H{"error": "parameter 'search' diperlukan"})
		return
	}

	items, err := h.service.SearchMenuItems(c.Request.Context(), keyword)
	if err != nil {
		log.Printf("[SearchMenuItems] Gagal mencari menu dengan keyword '%s': %v", keyword, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mencari menu"})
		return
	}

	c.JSON(http.StatusOK, items)
}

// UpdateMenuItem godoc
// @Summary Perbarui menu
// @Tags Menu
// @Accept json
// @Produce json
// @Param id path int true "ID Menu"
// @Param request body CreateMenuItemRequest true "Data menu yang diperbarui"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /menu/menu-items/{id} [put]
func (h *MenuItemHandler) UpdateMenuItem(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("[UpdateMenuItem] ID tidak valid: %s", idStr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	var req CreateMenuItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("[UpdateMenuItem] Bad request: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	menuItem := &models.MenuItem{
		ID:          id,
		CategoryID:  req.CategoryID,
		SKU:         req.SKU,
		Name:        req.Name,
		Description: sql.NullString{String: req.Description, Valid: req.Description != ""},
		Price:       req.Price,
		Cost:        req.Cost,
		IsActive:    req.IsActive,
		Tags:        req.Tags,
	}
	if req.PreparationTime != nil {
		menuItem.PreparationTime = sql.NullInt64{Int64: int64(*req.PreparationTime), Valid: true}
	}

	if err := h.service.UpdateMenuItem(c.Request.Context(), menuItem); err != nil {
		log.Printf("[UpdateMenuItem] Gagal update ID %d: %v", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui menu item"})
		return
	}

	log.Printf("[UpdateMenuItem] Berhasil update menu item ID %d", id)
	c.JSON(http.StatusOK, gin.H{"message": "Menu item berhasil diupdate"})
}

// DeleteMenuItem godoc
// @Summary Hapus (soft delete) menu
// @Tags Menu
// @Produce json
// @Param id path int true "ID Menu"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /menu/menu-items/{id} [delete]
func (h *MenuItemHandler) DeleteMenuItem(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	err = h.service.DeleteMenuItem(c.Request.Context(), id)
	if err != nil {
		log.Printf("[DeleteMenuItem] Gagal soft delete: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus menu"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Menu berhasil dihapus (soft delete)"})
}

// GetMenuDetail godoc
// @Summary Dapatkan detail menu dengan bahan
// @Tags Menu
// @Produce json
// @Param id path int true "ID Menu"
// @Success 200 {object} models.MenuItemWithIngredients
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /menu/menu-items/detail/{id} [get]
func (h *MenuItemHandler) GetMenuDetail(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("ID menu tidak valid: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	menu, err := h.service.GetMenuWithIngredients(c.Request.Context(), id)
	if err != nil {
		log.Printf("Gagal mengambil detail menu: %v", err)
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Menu tidak ditemukan"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data menu"})
		return
	}

	c.JSON(http.StatusOK, menu)
}
