package handlers

import (
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
	id, _ := strconv.Atoi(idStr)

	result, err := h.service.GetPesertaByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Peserta tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": result})
}
