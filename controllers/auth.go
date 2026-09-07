package controllers

import (
	"net/http"
	"time"

	"web_sarpras/config"
	"web_sarpras/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret = []byte("secret_key_rahasia_sarpras_123")

type RegisterInput struct {
	Nama               string `json:"nama" binding:"required"`
	Username           string `json:"username" binding:"required"`
	Password           string `json:"password" binding:"required,min=6"`
	KonfirmasiPassword string `json:"konfirmasi_password" binding:"required"`
}

type LoginInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Struct khusus untuk input penambahan user oleh Admin
type CreateUserInput struct {
	Nama     string `json:"nama" binding:"required"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
	Role     string `json:"role" binding:"required"`
}

func Register(c *gin.Context) {
	var input RegisterInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Semua field wajib diisi (minimal password 6 karakter)"})
		return
	}

	if input.Password != input.KonfirmasiPassword {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password dan Konfirmasi Password tidak cocok"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memproses password"})
		return
	}

	user := models.User{
		Nama:     input.Nama,
		Username: input.Username,
		Password: string(hashedPassword),
		Role:     "user",
	}

	if err := config.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username sudah terdaftar!"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Registrasi berhasil!",
		"data": gin.H{
			"id":       user.ID,
			"nama":     user.Nama,
			"username": user.Username,
			"role":     user.Role,
		},
	})
}

func Login(c *gin.Context) {
	var input LoginInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username dan password wajib diisi"})
		return
	}

	var user models.User
	if err := config.DB.Where("username = ?", input.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Username atau password salah"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Username atau password salah"})
		return
	}

	claims := jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"role":     user.Role,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat token autentikasi"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login berhasil",
		"token":   tokenString,
		"user": gin.H{
			"id":       user.ID,
			"nama":     user.Nama,
			"username": user.Username,
			"role":     user.Role,
		},
	})
}

func CreateUserByAdmin(c *gin.Context) {
	var input CreateUserInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Semua field wajib diisi (role: admin/teknisi/user)"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memproses password"})
		return
	}

	user := models.User{
		Nama:     input.Nama,
		Username: input.Username,
		Password: string(hashedPassword),
		Role:     input.Role,
	}

	if err := config.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username sudah terdaftar atau query gagal!"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User berhasil dibuat oleh Admin!",
		"data": gin.H{
			"id":       user.ID,
			"nama":     user.Nama,
			"username": user.Username,
			"role":     user.Role,
		},
	})
}