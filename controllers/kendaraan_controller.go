package controllers

import (
	"fmt"
	"net/http"
	"service-record/config"
	"service-record/models"
	"service-record/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

type KendaraanDTO struct {
	NoPolisi   string `json:"no_polisi" binding:"required"`
	Nama       string `json:"nama" binding:"required"`
	Manufaktur string `json:"manufaktur" binding:"required"`
	Jenis      string `json:"jenis" binding:"required" enums:"Mobil,Motor"`
	Tahun      string `json:"tahun" binding:"required"`
	Warna      string `json:"warna" binding:"required"`
}

// CreateKendaraan godoc
// @Summary Tambah Kendaraan Baru
// @Tags Master Kendaraan
// @Accept json
// @Produce json
// @Param kendaraan body controllers.KendaraanDTO true "Data Kendaraan (Enum Jenis: Mobil, Motor)"
// @Success 200 {object} models.Kendaraan
// @Router /api/kendaraan [post]
func CreateKendaraan(c *gin.Context) {
	var input KendaraanDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	utils.MultiToUpper(&input.NoPolisi, &input.Manufaktur, &input.Nama, &input.Jenis)
	kendaraan := models.Kendaraan{
		NoPolisi:   input.NoPolisi,
		Nama:       input.Nama,
		Manufaktur: input.Manufaktur,
		Jenis:      input.Jenis,
		Tahun:      input.Tahun,
		Warna:      input.Warna,
	}

	if err := config.DB.Create(&kendaraan).Error; err != nil {
		code, message := utils.HandleDBError(err)
		c.JSON(code, gin.H{"error": message})
		return
	}
	c.JSON(http.StatusOK, kendaraan)
}

// GetKendaraans godoc
// @Summary Ambil list kendaraan dengan pagination fleksibel
// @Description Mengambil data kendaraan. Contoh: ?page=1&limit=5 atau ?page=2&limit=15
// @Tags Master Kendaraan
// @Produce json
// @Param page query int false "Halaman ke berapa (default 1)" default(1)
// @Param limit query int false "Jumlah data per halaman (bisa 5, 10, 15, dll)" default(10)
// @Param aktif query string false "Filter Aktif (true, false, all)"
// @Param sort_by query string false "Kolom untuk sorting (nama, manufaktur, created_at)" default(nama)
// @Param order query string false "Arah urutan (asc, desc)" default(asc)
// @Param search query string false "Cari berdasarkan nama kendaraan"
// @Success 200 {object} utils.PaginationResult{data=[]models.Kendaraan} "Berhasil mengambil data"
// @Router /api/kendaraan [get]
func GetKendaraan(c *gin.Context) {
	var kendaraans []models.Kendaraan
	statusFilter := c.DefaultQuery("aktif", "all")
	sortBy := c.DefaultQuery("sort_by", "nama")
	order := c.DefaultQuery("order", "asc")
	searchName := c.Query("search")

	query := config.DB

	switch statusFilter {
	case "true":
		query = query.Where("aktif = ?", true)
	case "false":
		query = query.Where("aktif = ?", false)
	}

	if searchName != "" {
		// Ubah ke huruf besar karena saat simpan data kita sudah paksa UPPERCASE
		searchNameUpper := strings.ToUpper(searchName)
		query = query.Where("nama LIKE ?", "%"+searchNameUpper+"%")
	}

	allowedSortColumns := map[string]bool{
		"no_polisi":  true,
		"nama":       true,
		"manufaktur": true,
		"created_at": true,
	}

	if !allowedSortColumns[sortBy] {
		sortBy = "nama"
	}
	if order != "asc" && order != "desc" {
		order = "asc"
	}
	orderClause := fmt.Sprintf("%s %s", sortBy, order)
	query = query.Order(orderClause)

	result, err := utils.ExecuteWithPagination(c, query, &models.Kendaraan{}, &kendaraans)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data kendaraan"})
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetKendaraansByID godoc
// @Summary Ambil Data Kendaraan berdasarkan ID
// @Description Mengambil data spesifik kendaraan berdasarkan ID. Contoh: /api/kendaraan/G1234
// @Tags Master Kendaraan
// @Produce json
// @Param no_polisi path string true "No Polisi Kendaraan"
// @Success 200 {array} models.Kendaraan "Berhasil mengambil data"
// @Failure 404 {object} map[string]string "Data tidak ditemukan"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /api/kendaraan/{no_polisi} [get]
func GetKendaraanByID(c *gin.Context) {
	noPolisi := c.Param("no_polisi")
	var kendaraan []models.Kendaraan

	if err := config.DB.First(&kendaraan, "no_polisi = ?", noPolisi).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Kendaraan tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, kendaraan)
}

// GetKendaraansAktif godoc
// @Summary Ambil Data Kendaraan berdasarkan status
// @Description Mengambil data Kendaraan berdasarkan status. Contoh: /api/kendaraan/getKendaraanByStatus/
// @Tags Master Kendaraan
// @Produce json
// @Param status path string true "Status"
// @Success 200 {array} models.Kendaraan "Berhasil mengambil data"
// @Failure 404 {object} map[string]string "Data tidak ditemukan"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /api/kendaraan/getKendaraanByStatus/{status} [get]
func GetKendaraanByStatus(c *gin.Context) {
	status := c.Param("status")
	var kendaraans []models.Kendaraan

	if err := config.DB.Find(&kendaraans, "aktif = ?", status).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memproses data kendaraan"})
		return
	}

	if len(kendaraans) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tidak ada kendaraan dengan status tersebut"})
		return
	}

	c.JSON(http.StatusOK, kendaraans)
}

// UpdateKendaraan godoc
// @Summary Update Data Kendaraan
// @Description Mengubah data Kendaraan
// @Accept json
// @Produce json
// @Tags Master Kendaraan
// @Param no_polisi path string true "No Polisi"
// @Param kendaraan body models.Kendaraan true "Data Kendaraan"
// @Success 200 {object} models.Kendaraan
// @Router /api/kendaraan/{no_polisi} [put]
func UpdateKendaraan(c *gin.Context) {
	noPolisi := c.Param("no_polisi")
	var kendaraan models.Kendaraan

	if err := config.DB.First(&kendaraan, "no_polisi = ?", noPolisi).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Kendaraan tidak ditemukan"})
		return
	}

	if err := c.ShouldBindJSON(&kendaraan); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format data tidak valid: " + err.Error()})
		return
	}

	utils.MultiToUpper(&kendaraan.NoPolisi, &kendaraan.Manufaktur, &kendaraan.Nama, &kendaraan.Jenis, &kendaraan.Warna)

	if err := config.DB.Save(&kendaraan).Error; err != nil {
		code, message := utils.HandleDBError(err)
		c.JSON(code, gin.H{"error": message})
		return
	}
	c.JSON(http.StatusOK, kendaraan)
}
