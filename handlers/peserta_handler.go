package handlers

import (
	"kkn-system/middleware"
	"kkn-system/models/entity"
	"kkn-system/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PesertaHandler struct {
	service services.PesertaService
}

func NewPesertaHandler(service services.PesertaService) *PesertaHandler {
	return &PesertaHandler{service}
}

func (h *PesertaHandler) Register(c *gin.Context) {
	var input entity.Peserta
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.service.DaftarKKN(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Pendaftaran berhasil",
		"data":    result,
	})
}

func (h *PesertaHandler) GetAll(c *gin.Context) {
	result, err := h.service.GetAllPeserta()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": result})
}

func (h *PesertaHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID peserta tidak valid"})
		return
	}

	result, err := h.service.GetPesertaByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Peserta tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": result})
}

func (h *PesertaHandler) GetMyProfile(c *gin.Context) {
	userID, err := middleware.GetAuthUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	result, err := h.service.GetPesertaByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Peserta tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": result})
}

func (h *PesertaHandler) UpdateStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)

	var input struct {
		Status string `json:"status" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.service.UpdateStatus(uint(id), input.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Status berhasil diperbarui"})
}

func (h *PesertaHandler) AssignLocation(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)

	var input struct {
		LokasiID uint `json:"lokasi_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.service.AssignLocation(uint(id), input.LokasiID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Lokasi peserta berhasil diperbarui"})
}
