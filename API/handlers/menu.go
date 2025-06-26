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

func (h *MenuItemHandler) ListMenuItemsById(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID harus berupa angka"})
		return
	}

	items, err := h.service.ListMenuItemsById(c.Request.Context(), id)
	if err != nil {
		log.Printf("[ListMenuItemsById] Gagal mengambil data menu: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data menu"})
		return
	}
	c.JSON(http.StatusOK, items)
}

func (h *MenuItemHandler) ListMenuItems(c *gin.Context) {
	items, err := h.service.ListMenuItems(c.Request.Context())
	if err != nil {
		log.Printf("[ListMenuItems] Gagal mengambil data menu: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data menu"})
		return
	}
	c.JSON(http.StatusOK, items)
}

func (h *MenuItemHandler) ListActiveMenuItems(c *gin.Context) {
	items, err := h.service.ListActiveMenuItems(c.Request.Context())
	if err != nil {
		log.Printf("[ListActiveMenuItems] Gagal mengambil data aktif: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data aktif"})
		return
	}
	c.JSON(http.StatusOK, items)
}

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
