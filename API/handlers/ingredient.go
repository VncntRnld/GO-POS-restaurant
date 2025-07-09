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

type IngredientHandler struct {
	service *services.IngredientService
}

func NewIngredientHandler(service *services.IngredientService) *IngredientHandler {
	return &IngredientHandler{service: service}
}

type CreateIngredientRequest struct {
	Name        string  `json:"name" binding:"required"`
	Qty         float64 `json:"qty"`
	Unit        string  `json:"unit"`
	IsAllergen  bool    `json:"is_allergen"`
	IsActive    bool    `json:"is_active"`
	Description string  `json:"description"`
}

// CreateIngredient godoc
// @Summary Tambah ingredient baru
// @Tags Ingredient
// @Accept json
// @Produce json
// @Param request body CreateIngredientRequest true "Data ingredient baru"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /ingredients [post]
func (h *IngredientHandler) CreateIngredient(c *gin.Context) {
	var req CreateIngredientRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Error parsing request body: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ingredient := &models.Ingredient{
		Name:        req.Name,
		Qty:         req.Qty,
		Unit:        req.Unit,
		IsAllergen:  req.IsAllergen,
		IsActive:    req.IsActive,
		Description: sql.NullString{String: req.Description, Valid: req.Description != ""},
	}

	id, err := h.service.CreateIngredient(c.Request.Context(), ingredient)
	if err != nil {
		log.Printf("Error creating ingredient: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat ingredient"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": id, "name": ingredient.Name})
}

// ListIngredients godoc
// @Summary Tampilkan semua ingredient
// @Tags Ingredient
// @Produce json
// @Success 200 {array} models.Ingredient
// @Failure 500 {object} map[string]string
// @Router /ingredients [get]
func (h *IngredientHandler) ListIngredients(c *gin.Context) {
	ingredients, err := h.service.ListIngredients(c.Request.Context())
	if err != nil {
		log.Printf("Error listing ingredients: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil ingredient"})
		return
	}
	c.JSON(http.StatusOK, ingredients)
}

// GetIngredientByID godoc
// @Summary Ambil ingredient berdasarkan ID
// @Tags Ingredient
// @Produce json
// @Param id path int true "ID Ingredient"
// @Success 200 {object} models.Ingredient
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /ingredients/{id} [get]
func (h *IngredientHandler) GetIngredientByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("ID tidak valid: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	ingredient, err := h.service.GetIngredientByID(c.Request.Context(), id)
	if err != nil {
		log.Printf("Gagal mengambil ingredient ID %d: %v", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data ingredient"})
		return
	}

	c.JSON(http.StatusOK, ingredient)
}

// GetIngredientByID godoc
// @Summary Ambil ingredient berdasarkan ID
// @Tags Ingredient
// @Produce json
// @Param id path int true "ID Ingredient"
// @Success 200 {object} models.Ingredient
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /ingredients/{id} [get]
func (h *IngredientHandler) UpdateIngredient(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("ID tidak valid untuk update: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	var req CreateIngredientRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Error binding update request: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ingredient := &models.Ingredient{
		ID:          id,
		Name:        req.Name,
		Qty:         req.Qty,
		Unit:        req.Unit,
		IsAllergen:  req.IsAllergen,
		IsActive:    req.IsActive,
		Description: sql.NullString{String: req.Description, Valid: req.Description != ""},
	}

	if err := h.service.UpdateIngredient(c.Request.Context(), ingredient); err != nil {
		log.Printf("Gagal update ingredient ID %d: %v", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengupdate ingredient"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Ingredient berhasil diupdate"})
}

// DeleteIngredient godoc
// @Summary Hapus (soft delete) ingredient berdasarkan ID
// @Tags Ingredient
// @Produce json
// @Param id path int true "ID Ingredient"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /ingredients/{id} [delete]
func (h *IngredientHandler) DeleteIngredient(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("ID ingredient tidak valid untuk delete: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	err = h.service.DeleteIngredient(c.Request.Context(), id)
	if err != nil {
		log.Printf("Gagal menghapus (soft delete) ingredient ID %d: %v", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus ingredient"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Ingredient berhasil dihapus (soft delete)"})
}
