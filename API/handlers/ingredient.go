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

func (h *IngredientHandler) ListIngredients(c *gin.Context) {
	ingredients, err := h.service.ListIngredients(c.Request.Context())
	if err != nil {
		log.Printf("Error listing ingredients: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil ingredient"})
		return
	}
	c.JSON(http.StatusOK, ingredients)
}

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
