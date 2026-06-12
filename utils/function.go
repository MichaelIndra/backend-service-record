package utils

import (
	"math"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func OptionalString(val string) *string {
	if val == "" {
		return nil
	}
	return &val
}

func MultiToUpper(strs ...*string) {
	for _, str := range strs {
		if str != nil {
			*str = strings.ToUpper(*str)
		}
	}
}

type PaginationResult struct {
	TotalData int64       `json:"total_data" example:"100"`
	TotalPage int         `json:"total_page" example:"10"`
	Page      int         `json:"page" example:"1"`
	Limit     int         `json:"limit" example:"10"`
	Data      interface{} `json:"data"`
}

func ExecuteWithPagination(c *gin.Context, db *gorm.DB, model interface{}, dest interface{}) (PaginationResult, error) {
	// 1. Ambil query param ?page= dan ?limit= (limit adalah jumlah data: 5, 10, 15, dll)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10")) // Default 10 jika frontend tidak kirim

	// Validasi batas bawah
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	// Proteksi: Maksimal data per halaman dibatasi 100 agar server aman
	if limit > 100 {
		limit = 100
	}

	// 2. Hitung total seluruh data asli di DB
	var totalData int64
	db.Model(model).Count(&totalData)

	// 3. Hitung total halaman menggunakan pembulatan ke atas (Ceil)
	totalPage := int(math.Ceil(float64(totalData) / float64(limit)))
	if totalPage == 0 {
		totalPage = 1
	}

	// 4. Jalankan Query ke DB dengan Offset dan Limit
	offset := (page - 1) * limit
	err := db.Offset(offset).Limit(limit).Find(dest).Error

	// 5. Bungkus jadi object siap pakai
	result := PaginationResult{
		TotalData: totalData,
		TotalPage: totalPage,
		Page:      page,
		Limit:     limit,
		Data:      dest,
	}

	return result, err
}
