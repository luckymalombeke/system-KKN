package handlers

import (
	"kkn-system/middleware"
	"kkn-system/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type NotifikasiHandler struct {
	service services.NotifikasiService
}

func NewNotifikasiHandler(service services.NotifikasiService) *NotifikasiHandler {
	return &NotifikasiHandler{service}
}

func (h *NotifikasiHandler) Send(c *gin.Context) {
	var input struct {
		PesertaID uint   `json:"peserta_id" binding:"required"`
		Title     string `json:"title" binding:"required"`
		Message   string `json:"message" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.service.Send(input.PesertaID, input.Title, input.Message)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Notifikasi terkirim", "data": result})
}

func (h *NotifikasiHandler) GetMyNotifikasi(c *gin.Context) {
	userID, err := middleware.GetAuthUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	result, err := h.service.GetByPeserta(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": result})
}

func (h *NotifikasiHandler) MarkMyNotifikasiAsRead(c *gin.Context) {
	userID, err := middleware.GetAuthUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID notifikasi tidak valid"})
		return
	}

	err = h.service.ReadForPeserta(uint(id), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Notifikasi telah dibaca"})
}
