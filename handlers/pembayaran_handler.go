package handlers

import (
	"kkn-system/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PembayaranHandler struct {
	service services.PembayaranService
}

func NewPembayaranHandler(service services.PembayaranService) *PembayaranHandler {
	return &PembayaranHandler{service}
}

func (h *PembayaranHandler) CreateInvoice(c *gin.Context) {
	var input struct {
		PesertaID uint  `json:"peserta_id" binding:"required"`
		Amount    int64 `json:"amount" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.service.CreateInvoice(input.PesertaID, input.Amount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Invoice berhasil dibuat",
		"data":    result,
	})
}

func (h *PembayaranHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)

	result, err := h.service.GetStatus(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Data pembayaran tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": result})
}
