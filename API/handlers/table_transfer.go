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

type TableTransferHandler struct {
	service *services.TableTransferService
}

func NewTableTransferHandler(service *services.TableTransferService) *TableTransferHandler {
	return &TableTransferHandler{service}
}

type CreateTableTransferRequest struct {
	OrderID       int    `json:"order_id" binding:"required"`
	FromTableID   int    `json:"from_table_id" binding:"required"`
	ToTableID     int    `json:"to_table_id" binding:"required"`
	TransferredBy int    `json:"transferred_by" binding:"required"`
	Reason        string `json:"reason"`
}

func (h *TableTransferHandler) Create(c *gin.Context) {
	var req CreateTableTransferRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	data := &models.TableTransfer{
		OrderID:       req.OrderID,
		FromTableID:   req.FromTableID,
		ToTableID:     req.ToTableID,
		TransferredBy: req.TransferredBy,
		Reason:        sql.NullString{String: req.Reason, Valid: req.Reason != ""},
	}

	id, err := h.service.Create(c.Request.Context(), data)
	if err != nil {
		log.Printf("Gagal create table transfer: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menambahkan data"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": id})
}

func (h *TableTransferHandler) List(c *gin.Context) {
	data, err := h.service.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data"})
		return
	}
	c.JSON(http.StatusOK, data)
}

func (h *TableTransferHandler) GetByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	data, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Data tidak ditemukan"})
		return
	}
	c.JSON(http.StatusOK, data)
}

func (h *TableTransferHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	var req CreateTableTransferRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	data := &models.TableTransfer{
		ID:            id,
		OrderID:       req.OrderID,
		FromTableID:   req.FromTableID,
		ToTableID:     req.ToTableID,
		TransferredBy: req.TransferredBy,
		Reason:        sql.NullString{String: req.Reason, Valid: req.Reason != ""},
	}

	if err := h.service.Update(c.Request.Context(), data); err != nil {
		log.Printf("Gagal update table transfer ID %d: %v", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengupdate data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Data berhasil diperbarui"})
}

func (h *TableTransferHandler) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus data"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Berhasil dihapus"})
}
