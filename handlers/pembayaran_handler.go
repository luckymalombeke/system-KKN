package handlers

import (
	"encoding/json"
	"io"
	"kkn-system/middleware"
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

func (h *PembayaranHandler) CreateMyInvoice(c *gin.Context) {
	userID, err := middleware.GetAuthUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	var input struct {
		Amount int64 `json:"amount" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.service.CreateInvoice(userID, input.Amount)
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

func (h *PembayaranHandler) GetMyPembayaran(c *gin.Context) {
	userID, err := middleware.GetAuthUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	result, err := h.service.GetByPesertaID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tagihan belum dibuat"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": result})
}

func (h *PembayaranHandler) HandleNotification(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "gagal membaca body notifikasi"})
		return
	}

	var notification services.MidtransNotification
	if err := json.Unmarshal(body, &notification); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "format notifikasi tidak valid"})
		return
	}

	err = h.service.HandleMidtransNotification(notification)
	if err != nil {
		if err.Error() == "signature_key tidak valid" || err.Error() == "signature_key wajib ada" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Status pembayaran berhasil diupdate"})
}
