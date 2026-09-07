package controllers

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"web_sarpras/config"
	"web_sarpras/models"

	"github.com/gin-gonic/gin"
)

func generateKodeAset(jenis string) string {
	prefix := "SRN"
	if strings.ToLower(jenis) == "prasarana" {
		prefix = "PRS"
	}
	return fmt.Sprintf("%s-%d", prefix, time.Now().UnixNano()/1e6)
}

func GetAllAset(c *gin.Context) {
	var asetList []models.Aset
	config.DB.Find(&asetList)

	c.JSON(http.StatusOK, gin.H{
		"data": asetList,
	})
}

func GetAsetByID(c *gin.Context) {
	id := c.Param("id")
	var aset models.Aset

	if err := config.DB.First(&aset, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Aset tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": aset,
	})
}

func CreateAset(c *gin.Context) {
	var input models.Aset
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if strings.TrimSpace(input.KodeAset) == "" {
		input.KodeAset = generateKodeAset(input.Jenis)
	}

	if err := config.DB.Create(&input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menambah aset"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Aset berhasil ditambahkan",
		"data":    input,
	})
}

func UpdateAset(c *gin.Context) {
	id := c.Param("id")
	var aset models.Aset

	if err := config.DB.First(&aset, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Aset tidak ditemukan"})
		return
	}

	var input models.Aset
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	config.DB.Model(&aset).Updates(models.Aset{
		KodeAset: input.KodeAset,
		NamaAset: input.NamaAset,
		Jenis:    input.Jenis,
		Kategori: input.Kategori,
		Jumlah:   input.Jumlah,
		Kondisi:  input.Kondisi,
		Lokasi:   input.Lokasi,
	})

	c.JSON(http.StatusOK, gin.H{
		"message": "Data aset berhasil diperbarui",
		"data":    aset,
	})
}

func DeleteAset(c *gin.Context) {
	id := c.Param("id")
	var aset models.Aset

	if err := config.DB.First(&aset, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Aset tidak ditemukan"})
		return
	}

	config.DB.Delete(&aset)

	c.JSON(http.StatusOK, gin.H{
		"message": "Aset berhasil dihapus",
	})
}