package handlers

import (
	"kkn-system/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	service services.AuthService
}

func NewAuthHandler(service services.AuthService) *AuthHandler {
	return &AuthHandler{service}
}

func (h *AuthHandler) RequestOTP(c *gin.Context) {
	var input struct {
		NIM string `json:"nim" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "NIM wajib diisi"})
		return
	}

	if err := h.service.RequestOTP(input.NIM); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "OTP berhasil dikirim ke email terdaftar",
	})
}

func (h *AuthHandler) VerifyOTP(c *gin.Context) {
	var input struct {
		NIM string `json:"nim" binding:"required"`
		OTP string `json:"otp" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "NIM dan OTP wajib diisi"})
		return
	}

	token, peserta, err := h.service.VerifyOTP(input.NIM, input.OTP)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "Login berhasil",
		"token":      token,
		"peserta_id": peserta.ID,
	})
}

func (h *AuthHandler) AdminLogin(c *gin.Context) {
	var input struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format email salah atau password kosong"})
		return
	}

	token, admin, err := h.service.AdminLogin(input.Email, input.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Admin login berhasil",
		"token":   token,
		"admin": gin.H{
			"id":    admin.ID,
			"name":  admin.Name,
			"email": admin.Email,
			"role":  admin.Role,
		},
	})
}
