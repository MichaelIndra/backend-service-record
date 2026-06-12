package controllers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"service-record/config"
	"service-record/models"
	"service-record/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// CreateTransaction godoc
// @Summary Simpan Transaksi Lengkap
// @Description Simpan header transaksi beserta banyak detail service dan sparepart sekaligus.
// @Tags Transaction
// @Accept multipart/form-data
// @Produce json
// @Param no_polisi formData string true "No Polisi Kendaraan"
// @Param bengkel_id formData int true "ID Bengkel"
// @Param biaya_bruto formData number true "Total Biaya (Bruto)"
// @Param diskon formData number true "Total Diskon"
// @Param biaya_netto formData number true "Total Bayar (Netto)"
// @Param tanggal_transaksi formData string true "Format: 2006-01-02"
// @Param tanggal_masuk formData string true "Format: 2006-01-02"
// @Param tanggal_selesai formData string true "Format: 2006-01-02"
// @Param km formData string false "Kilometer Kendaraan"
// @Param no_nota formData string false "Nomor Nota Manual"
// @Param foto_nota[] formData []file false "Foto Nota Utama (Bisa upload lebih dari 1 foto)" collectionFormat(multi)
// @Param service_name[] formData []string false "Daftar Nama Service" collectionFormat(multi)
// @Param service_price[] formData []number false "Daftar Harga Service" collectionFormat(multi)
// @Param foto_service[] formData []file false "Daftar Foto Service" collectionFormat(multi)
// @Param sparepart_name[] formData []string false "Daftar Nama Sparepart" collectionFormat(multi)
// @Param sparepart_code[] formData []string false "Daftar Kode Sparepart" collectionFormat(multi)
// @Param sparepart_price[] formData []number false "Daftar Harga Sparepart" collectionFormat(multi)
// @Param sparepart_batas_km[] formData []string false "Daftar Batas KM Sparepart" collectionFormat(multi)
// @Param sparepart_batas_waktu[] formData []string false "Daftar Batas Waktu Sparepart" collectionFormat(multi)
// @Param sparepart_remind[] formData []boolean false "Daftar Remind Sparepart" collectionFormat(multi)
// @Param sparepart_qty[] formData []number false "Daftar Qty Sparepart" collectionFormat(multi)
// @Param foto_sparepart[] formData []file false "Daftar Foto Sparepart" collectionFormat(multi)
// @Success 200 {object} models.Transaction
// @Router /api/transaction [post]
func CreateTransaction(c *gin.Context) {

	// fmt.Println("========== DEBUG DATA MASUK ==========")
	// if err := c.Request.ParseMultipartForm(32 << 20); err == nil { // 32MB max memory
	// 	fmt.Println("👉 DATA FORM/TEXT:")
	// 	for key, values := range c.Request.MultipartForm.Value {
	// 		fmt.Printf("   %s: %v\n", key, values)
	// 	}

	// 	fmt.Println("\n👉 DATA FILES/GAMBAR:")
	// 	for key, files := range c.Request.MultipartForm.File {
	// 		fmt.Printf("   %s (%d file):\n", key, len(files))
	// 		for _, file := range files {
	// 			fmt.Printf("      - Nama: %s, Ukuran: %d bytes\n", file.Filename, file.Size)
	// 		}
	// 	}
	// } else {
	// 	fmt.Println("❌ Gagal membaca multipart form:", err)
	// }
	// fmt.Println("======================================")

	db := config.DB

	// 1. Mulai DB Transaction
	tx := db.Begin()

	// 2. Parsing Data Header
	bengkelID, _ := strconv.Atoi(c.PostForm("bengkel_id"))
	biayaBruto, _ := strconv.ParseFloat(c.PostForm("biaya_bruto"), 64)
	biayaNetto, _ := strconv.ParseFloat(c.PostForm("biaya_netto"), 64)
	diskon, _ := strconv.ParseFloat(c.PostForm("diskon"), 64)
	tglTrx, _ := time.Parse("2006-01-02", c.PostForm("tanggal_transaksi"))

	// 🟢 PERBAIKAN LOGIKA DATE POINTER: Jika kosong, biarkan tetap nil (NULL di database)
	var tglMasukPtr *time.Time
	if val := c.PostForm("tanggal_masuk"); val != "" {
		if t, err := time.Parse("2006-01-02", val); err == nil {
			tglMasukPtr = &t
		}
	}

	var tglSelesaiPtr *time.Time
	if val := c.PostForm("tanggal_selesai"); val != "" {
		if t, err := time.Parse("2006-01-02", val); err == nil {
			tglSelesaiPtr = &t
		}
	}

	noPolisi := c.PostForm("no_polisi")
	noNota := c.PostForm("no_nota")

	// SULAP Header
	utils.MultiToUpper(&noPolisi, &noNota)

	transaction := models.Transaction{
		NoPolisi:         noPolisi,
		BengkelID:        uint(bengkelID),
		BiayaBruto:       biayaBruto,
		BiayaNetto:       biayaNetto,
		Diskon:           diskon,
		TanggalTransaksi: tglTrx,
		TanggalMasuk:     tglMasukPtr,   // Menggunakan pointer aman
		TanggalSelesai:   tglSelesaiPtr, // Menggunakan pointer aman
		NoNota:           utils.OptionalString(noNota),
		Km:               utils.OptionalString(c.PostForm("km")),
	}

	if err := tx.Create(&transaction).Error; err != nil {
		tx.Rollback()
		code, message := utils.HandleDBError(err)
		c.JSON(code, gin.H{"error": message})
		return
	}

	// 3. Siapkan Folder Penyimpanan (uploads/trx_ID)
	folderPath := fmt.Sprintf("uploads/trx_%d", transaction.ID)
	if err := os.MkdirAll(folderPath, os.ModePerm); err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat folder foto"})
		return
	}

	form, _ := c.MultipartForm()

	// 4. Handle MULTIPLE Foto Nota Utama
	fileNotas := form.File["foto_nota[]"]
	var notaPaths []string

	for _, fileNota := range fileNotas {
		ext := filepath.Ext(fileNota.Filename)
		filename := "nota_" + uuid.New().String() + ext
		savePath := filepath.Join(folderPath, filename)

		if err := c.SaveUploadedFile(fileNota, savePath); err == nil {
			// 🟢 PERBAIKAN 1: Gunakan filepath.ToSlash agar Windows \\ berubah jadi /
			notaPaths = append(notaPaths, filepath.ToSlash(savePath))
		}
	}

	if len(notaPaths) > 0 {
		transaction.FotoNota = notaPaths
		tx.Save(&transaction)
	}

	// 5. Handle Array Detail Service
	sNames := c.PostFormArray("service_name[]")
	sPrices := c.PostFormArray("service_price[]")
	sFiles := form.File["foto_service[]"]

	for i, name := range sNames {
		if name == "" {
			continue
		}

		utils.MultiToUpper(&name)

		price, _ := strconv.ParseFloat(sPrices[i], 64)
		detail := models.DetailService{
			TransactionID: transaction.ID,
			NamaService:   name,
			Biaya:         price,
		}

		if i < len(sFiles) {
			ext := filepath.Ext(sFiles[i].Filename)
			filename := fmt.Sprintf("svc_%s%s", uuid.New().String(), ext)
			savePath := filepath.Join(folderPath, filename)
			c.SaveUploadedFile(sFiles[i], savePath)

			// 🟢 PERBAIKAN 2: Gunakan filepath.ToSlash
			detail.FotoService = filepath.ToSlash(savePath)
		}

		if err := tx.Create(&detail).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal simpan detail service ke-" + strconv.Itoa(i+1)})
			return
		}
	}

	// 6. Handle Array Detail Sparepart
	spNames := c.PostFormArray("sparepart_name[]")
	spCodes := c.PostFormArray("sparepart_code[]")
	spPrices := c.PostFormArray("sparepart_price[]")
	spQtys := c.PostFormArray("sparepart_qty[]")
	spFiles := form.File["foto_sparepart[]"]
	sparepart_batas_km := c.PostFormArray("sparepart_batas_km[]")
	sparepart_batas_waktu := c.PostFormArray("sparepart_batas_waktu[]")
	sparepart_remind := c.PostFormArray("sparepart_remind[]")

	for i, name := range spNames {
		if name == "" {
			continue
		}

		var codeStr string
		if i < len(spCodes) {
			codeStr = spCodes[i]
		}

		utils.MultiToUpper(&name, &codeStr)

		price, _ := strconv.ParseFloat(spPrices[i], 64)
		qty, _ := strconv.ParseFloat(spQtys[i], 64)

		detail := models.DetailSparepart{
			TransactionID: transaction.ID,
			NamaSparepart: name,
			KodeSparepart: utils.OptionalString(codeStr),
			HargaSatuan:   price,
			Qty:           qty,
			BatasKM:       utils.OptionalString(sparepart_batas_km[i]),
			BatasWaktu:    utils.OptionalString(sparepart_batas_waktu[i]),
			Remind:        sparepart_remind[i] == "true",
		}

		if i < len(spFiles) {
			ext := filepath.Ext(spFiles[i].Filename)
			filename := fmt.Sprintf("part_%s%s", uuid.New().String(), ext)
			savePath := filepath.Join(folderPath, filename)
			c.SaveUploadedFile(spFiles[i], savePath)

			// 🟢 PERBAIKAN 3: Gunakan filepath.ToSlash
			detail.FotoSparepart = filepath.ToSlash(savePath)
		}

		if err := tx.Create(&detail).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal simpan detail sparepart ke-" + strconv.Itoa(i+1)})
			return
		}
	}

	// 7. Selesai & Commit
	tx.Commit()

	db.Preload("DetailsService").Preload("DetailsSparepart").First(&transaction, transaction.ID)
	c.JSON(http.StatusOK, gin.H{
		"message": "Transaksi berhasil disimpan dengan seluruh detail",
		"data":    transaction,
	})
}

// GetTransaction godoc
// @Summary Ambil semua list transaksi
// @Description Mengambil data bengkel.
// @Tags Transaction
// @Produce json
// @Param page query int false "Halaman ke berapa (default 1)" default(1)
// @Param limit query int false "Jumlah data per halaman (bisa 5, 10, 15, dll)" default(10)
// @Param sort_by query string false "Kolom untuk sorting (tanggal_transaksi, created_at)" default(tanggal_transaksi)
// @Param order query string false "Arah urutan (asc, desc)" default(desc)
// @Success 200 {object} models.Transaction
// @Router /api/transaction [get]
func GetTransaction(c *gin.Context) {
	var transactions []models.Transaction
	sortBy := c.DefaultQuery("sort_by", "tanggal_transaksi")
	order := c.DefaultQuery("order", "desc")

	query := config.DB

	allowedSortColumns := map[string]bool{
		"tanggal_transaksi": true,
		"created_at":        true,
	}

	if !allowedSortColumns[sortBy] {
		sortBy = "tanggal_transaksi"
	}
	if order != "asc" && order != "desc" {
		order = "desc"
	}
	orderClause := fmt.Sprintf("%s %s", sortBy, order)
	query = query.
		Preload("Kendaraan").
		Preload("Bengkel").
		Preload("DetailsService").
		Preload("DetailsSparepart").
		Order(orderClause)
	result, err := utils.ExecuteWithPagination(c, query, &models.Transaction{}, &transactions)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data "})
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetTransactionByID godoc
// @Summary Ambil detail satu transaksi berdasarkan ID beserta semua detailnya
// @Description Mengambil data transaksi lengkap termasuk riwayat service dan sparepartnya menggunakan ID dari URL
// @Tags Transaction
// @Produce json
// @Param id path int true "ID Transaksi"
// @Success 200 {object} models.Transaction "Data transaksi ditemukan"
// @Failure 404 {object} map[string]string "Data transaksi tidak ditemukan"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /api/transaction/{id} [get]
func GetTransactionByID(c *gin.Context) {
	// 1. Ambil ID dari parameter URL (misal: /api/transactions/5)
	id := c.Param("id")

	var transaction models.Transaction

	// 2. Cari ke database dan tarik juga data relasi detailnya menggunakan Preload
	err := config.DB.
		Preload("Kendaraan").
		Preload("Bengkel").
		Preload("DetailsService").
		Preload("DetailsSparepart").
		First(&transaction, "id = ?", id).
		Error

	// 3. Jika data tidak ditemukan di database
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Transaksi dengan ID " + id + " tidak ditemukan",
		})
		return
	}

	// 4. Jika sukses, kembalikan objek transaksi utuh
	c.JSON(http.StatusOK, transaction)
}

// DeleteTransaction godoc
// @Summary Hapus Transaksi Bersih (Clean Delete)
// @Description Menghapus header transaksi, detail service, detail sparepart, serta menyapu bersih folder foto terkait di uploads.
// @Tags Transaction
// @Produce json
// @Param id path int true "ID Transaksi"
// @Success 200 {object} map[string]string "Transaksi dan seluruh aset berhasil dihapus"
// @Failure 404 {object} map[string]string "Transaksi tidak ditemukan"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /api/transaction/{id} [delete]
func DeleteTransaction(c *gin.Context) {
	db := config.DB
	idStr := c.Param("id")
	transactionID, _ := strconv.Atoi(idStr)

	// 1. Pastikan data transaksi memang ada sebelum dihapus
	var transaction models.Transaction
	if err := db.First(&transaction, transactionID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Transaksi tidak ditemukan"})
		return
	}

	// 2. Mulai DB Transaction
	tx := db.Begin()

	// 3. Hapus Data Anak Pertama: Detail Service
	if err := tx.Where("transaction_id = ?", transaction.ID).Delete(&models.DetailService{}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus detail service"})
		return
	}

	// 4. Hapus Data Anak Kedua: Detail Sparepart
	if err := tx.Where("transaction_id = ?", transaction.ID).Delete(&models.DetailSparepart{}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus detail sparepart"})
		return
	}

	// 5. Hapus Data Induk: Transaction Header
	if err := tx.Delete(&transaction).Error; err != nil {
		tx.Rollback()
		code, message := utils.HandleDBError(err)
		c.JSON(code, gin.H{"error": message})
		return
	}

	// 6. Sapu Bersih Folder Foto (uploads/trx_ID) beserta isinya
	folderPath := fmt.Sprintf("uploads/trx_%d", transaction.ID)
	// os.RemoveAll akan menghapus folder beserta seluruh file di dalamnya tanpa sisa
	if err := os.RemoveAll(folderPath); err != nil {
		tx.Rollback() // Rollback DB jika gagal menghapus file fisik demi integritas data
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus folder foto fisik transaksi"})
		return
	}

	// 7. Jika semua langkah di atas sukses, Commit!
	tx.Commit()

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("Transaksi ID %d beserta seluruh detail dan file foto berhasil dihapus bersih", transaction.ID),
	})
}
