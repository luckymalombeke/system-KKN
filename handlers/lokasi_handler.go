package handlers

import (
	"kkn-system/models/entity"
	"kkn-system/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type LokasiHandler struct {
	service services.LokasiService
}

func NewLokasiHandler(service services.LokasiService) *LokasiHandler {
	return &LokasiHandler{service}
}

func (h *LokasiHandler) Create(c *gin.Context) {
	var input entity.Lokasi
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.service.CreateLokasi(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Lokasi berhasil ditambahkan", "data": result})
}

func (h *LokasiHandler) GetAll(c *gin.Context) {
	result, err := h.service.GetAllLokasi()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": result})
}
