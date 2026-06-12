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

type BengkelDTO struct {
	Nama         string  `json:"nama" binding:"required"`
	JenisService string  `json:"jenis_service" binding:"required"`
	Tipe         string  `json:"tipe" binding:"required"`
	NoTelp       *string `json:"no_telp"`
	Provinsi     string  `json:"provinsi" binding:"required"`
	KotaKab      string  `json:"kota_kab" binding:"required"`
	Kecamatan    string  `json:"kecamatan" binding:"required"`
	Kelurahan    string  `json:"kelurahan" binding:"required"`
	Alamat       string  `json:"alamat" binding:"required"`
}

// CreateBengkel godoc
// @Summary Tambah Bengkel Baru
// @Tags Master Bengkel
// @Accept json
// @Produce json
// @Param bengkel body controllers.BengkelDTO true "Data Bengkel"
// @Success 200 {object} models.Bengkel
// @Router /api/bengkel [post]
func CreateBengkel(c *gin.Context) {
	var input BengkelDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	utils.MultiToUpper(&input.Nama, &input.JenisService, &input.Tipe, &input.Provinsi, &input.KotaKab, &input.Kecamatan, &input.Kelurahan, &input.Alamat)
	bengkel := models.Bengkel{
		Nama:         input.Nama,
		JenisService: input.JenisService,
		Tipe:         input.Tipe,
		NoTelp:       input.NoTelp,
		Provinsi:     input.Provinsi,
		KotaKab:      input.KotaKab,
		Kecamatan:    input.Kecamatan,
		Kelurahan:    input.Kelurahan,
		Alamat:       input.Alamat,
	}

	if err := config.DB.Create(&bengkel).Error; err != nil {
		code, message := utils.HandleDBError(err)
		c.JSON(code, gin.H{"error": message})
		return
	}
	c.JSON(http.StatusOK, bengkel)
}

// GetBengkel godoc
// @Summary List Semua Bengkel dengan pagination fleksibel
// @Description Mengambil data bengkel. Contoh: ?page=1&limit=5 atau ?page=2&limit=15
// @Tags Master Bengkel
// @Produce json
// @Param aktif query string false "Filter Aktif (true, false, all)"
// @Param tipe query string false "Tipe Bengkel (mobil bbm, motor bbm, mobil ev, listrik ev, mobil bbm & ev, motor bbm & ev)"
// @Param page query int false "Halaman ke berapa (default 1)" default(1)
// @Param limit query int false "Jumlah data per halaman (bisa 5, 10, 15, dll)" default(10)
// @Param sort_by query string false "Kolom untuk sorting (nama, jenis_service, tipe, created_at)" default(nama)
// @Param order query string false "Arah urutan (asc, desc)" default(asc)
// @Param search query string false "Cari berdasarkan nama bengkel"
// @Param provinsi query string false "Filter nama Provinsi"
// @Param kota_kab query string false "Filter nama Kota/Kabupaten"
// @Param kecamatan query string false "Filter nama Kecamatan"
// @Param kelurahan query string false "Filter nama Kelurahan"
// @Success 200 {object} utils.PaginationResult{data=[]models.Bengkel} "Berhasil mengambil data"
// @Router /api/bengkel [get]
func GetBengkel(c *gin.Context) {
	var bengkel []models.Bengkel
	statusFilter := c.DefaultQuery("aktif", "all")
	sortBy := c.DefaultQuery("sort_by", "nama")
	order := c.DefaultQuery("order", "asc")

	searchName := c.Query("search")
	filterProv := c.Query("provinsi")
	filterKota := c.Query("kota_kab")
	filterKec := c.Query("kecamatan")
	filterKel := c.Query("kelurahan")
	filterTipe := c.Query("tipe")

	query := config.DB

	switch statusFilter {
	case "true":
		query = query.Where("aktif = ?", true)
	case "false":
		query = query.Where("aktif = ?", false)
	}

	if filterTipe != "" {
		query = query.Where("tipe = ?", strings.ToUpper(filterTipe))
	}

	if searchName != "" {
		// Ubah ke huruf besar karena saat simpan data kita sudah paksa UPPERCASE
		searchNameUpper := strings.ToUpper(searchName)
		query = query.Where("nama LIKE ?", "%"+searchNameUpper+"%")
	}

	if filterProv != "" {
		query = query.Where("provinsi = ?", strings.ToUpper(filterProv))
	}
	if filterKota != "" {
		query = query.Where("kota_kab = ?", strings.ToUpper(filterKota))
	}
	if filterKec != "" {
		query = query.Where("kecamatan = ?", strings.ToUpper(filterKec))
	}
	if filterKel != "" {
		query = query.Where("kelurahan = ?", strings.ToUpper(filterKel))
	}

	allowedSortColumns := map[string]bool{
		"nama":          true,
		"jenis_service": true,
		"tipe":          true,
		"created_at":    true,
	}

	if !allowedSortColumns[sortBy] {
		sortBy = "nama"
	}
	if order != "asc" && order != "desc" {
		order = "asc"
	}
	orderClause := fmt.Sprintf("%s %s", sortBy, order)
	query = query.Order(orderClause)
	result, err := utils.ExecuteWithPagination(c, query, &models.Bengkel{}, &bengkel)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data bengkel"})
		return
	}

	c.JSON(http.StatusOK, result)

}

// GetBengkelByID godoc
// @Summary Ambil Data Bengkel berdasarkan ID
// @Description Mengambil data spesifik bengkel berdasarkan ID. Contoh: /api/bengkel/1
// @Tags Master Bengkel
// @Produce json
// @Param id path int true "ID Bengkel"
// @Success 200 {object} models.Bengkel "Berhasil mengambil data"
// @Failure 404 {object} map[string]string "Data tidak ditemukan"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /api/bengkel/{id} [get]
func GetBengkelByID(c *gin.Context) {
	id := c.Param("id")
	var bengkel models.Bengkel

	if err := config.DB.First(&bengkel, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Bengkel dengan ID " + id + " tidak ditemukan"})
		return
	}
	c.JSON(http.StatusOK, bengkel)
}

// GetBengkelAktif godoc
// @Summary Ambil Data Bengkel berdasarkan status
// @Description Mengambil data Bengkel berdasarkan status. Contoh: /api/bengkel/getBengkelByStatus/
// @Tags Master Bengkel
// @Produce json
// @Param status path string true "Status"
// @Success 200 {array} models.Bengkel "Berhasil mengambil data"
// @Failure 404 {object} map[string]string "Data tidak ditemukan"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /api/bengkel/getBengkelByStatus/{status} [get]
func GetBengkelByStatus(c *gin.Context) {
	status := c.Param("status")
	var bengkels []models.Bengkel

	if err := config.DB.Find(&bengkels, "aktif = ?", status).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memproses data bengkel"})
		return
	}

	if len(bengkels) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tidak ada bengkel dengan status tersebut"})
		return
	}

	c.JSON(http.StatusOK, bengkels)
}

// UpdateBengkel godoc
// @Summary Update Data Bengkel
// @Description Mengubah data bengkel
// @Tags Master Bengkel
// @Accept json
// @Produce json
// @Param id path string true "ID Bengkel"
// @Param bengkel body models.Bengkel true "Data Bengkel"
// @Success 200 {object} models.Bengkel
// @Router /api/bengkel/{id} [put]
func UpdateBengkel(c *gin.Context) {
	id := c.Param("id")
	var bengkel models.Bengkel

	if err := config.DB.First(&bengkel, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Bengkel tidak ditemukan"})
		return
	}

	if err := c.ShouldBindJSON(&bengkel); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format data tidak valid: " + err.Error()})
		return
	}

	utils.MultiToUpper(
		&bengkel.Nama,
		&bengkel.JenisService,
		&bengkel.Provinsi,
		&bengkel.KotaKab,
		&bengkel.Kecamatan,
		&bengkel.Kelurahan,
		&bengkel.Alamat,
		&bengkel.Tipe,
	)

	if err := config.DB.Save(&bengkel).Error; err != nil {
		code, message := utils.HandleDBError(err)
		c.JSON(code, gin.H{"error": message})
		return
	}
	c.JSON(http.StatusOK, bengkel)
}
