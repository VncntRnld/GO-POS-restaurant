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

type BillHandler struct {
	service *services.BillService
}

func NewBillHandler(service *services.BillService) *BillHandler {
	return &BillHandler{service: service}
}

type CreateBillRequest struct {
	OrderID        int     `json:"order_id" binding:"required"`
	DiscountAmount float64 `json:"discount_amount"`
}

// Create godoc
// @Summary Buat tagihan untuk sebuah order
// @Tags Bills
// @Accept json
// @Produce json
// @Param request body CreateBillRequest true "Data tagihan"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /bills [post]
func (h *BillHandler) Create(c *gin.Context) {
	var req CreateBillRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Invalid request body: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Permintaan tidak valid"})
		return
	}

	billID, err := h.service.Create(c.Request.Context(), req.OrderID, req.DiscountAmount)
	if err != nil {
		log.Printf("Gagal membuat tagihan untuk order %d: %v", req.OrderID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat tagihan"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Tagihan berhasil dibuat", "bill_id": billID})
}

// CreateSplit godoc
// @Summary Buat tagihan split dari satu order
// @Tags Bills
// @Accept json
// @Produce json
// @Param request body models.SplitBillRequest true "Data split bill"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /bills/split [post]
func (h *BillHandler) CreateSplit(c *gin.Context) {
	var req models.SplitBillRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Request tidak valid: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Permintaan tidak valid"})
		return
	}

	billIDs, err := h.service.CreateSplit(c.Request.Context(), req)
	if err != nil {
		log.Printf("Gagal membuat split bill: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat split bill"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Split bill berhasil dibuat", "bill_ids": billIDs})
}

// List godoc
// @Summary Ambil semua tagihan
// @Tags Bills
// @Produce json
// @Success 200 {array} models.Bill
// @Failure 500 {object} map[string]string
// @Router /bills [get]
func (h *BillHandler) List(c *gin.Context) {
	bills, err := h.service.List(c.Request.Context())
	if err != nil {
		log.Printf("Gagal mengambil daftar tagihan: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil daftar tagihan"})
		return
	}

	c.JSON(http.StatusOK, bills)
}

// GetByID godoc
// @Summary Ambil tagihan berdasarkan ID
// @Tags Bills
// @Produce json
// @Param id path int true "ID tagihan"
// @Success 200 {object} models.Bill
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /bills/{id} [get]
func (h *BillHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("ID tidak valid: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	bill, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		log.Printf("Gagal mengambil bill ID %d: %v", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil bill"})
		return
	}
	if bill == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Bill tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, bill)
}

// Delete godoc
// @Summary Soft delete tagihan
// @Tags Bills
// @Produce json
// @Param id path int true "ID tagihan"
// @Success 200 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /bills/{id} [delete]
func (h *BillHandler) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := h.service.SoftDelete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete bill"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Bill deleted"})
}

type BillPaymentRequest struct {
	BillID               int     `json:"bill_id" binding:"required"`
	PaymentMethod        string  `json:"payment_method" binding:"required"`
	Amount               float64 `json:"amount" binding:"required"`
	ReferenceNumber      string  `json:"reference_number"`
	RoomChargeApprovedBy int     `json:"room_charge_approved_by"`
}

// Pay godoc
// @Summary Proses pembayaran tagihan
// @Tags Bills
// @Accept json
// @Produce json
// @Param request body BillPaymentRequest true "Data pembayaran"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /bills/pay [post]
func (h *BillHandler) Pay(c *gin.Context) {
	var req BillPaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	payment := &models.BillPayment{
		BillID:               req.BillID,
		PaymentMethod:        req.PaymentMethod,
		Amount:               req.Amount,
		ReferenceNumber:      sql.NullString{String: req.ReferenceNumber, Valid: req.ReferenceNumber != ""},
		RoomChargeApprovedBy: sql.NullInt64{Int64: int64(req.RoomChargeApprovedBy), Valid: req.RoomChargeApprovedBy != 0},
	}

	err := h.service.Pay(c.Request.Context(), payment)
	if err != nil {
		log.Printf("Gagal memproses pembayaran: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memproses pembayaran"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Pembayaran berhasil diproses"})
}
