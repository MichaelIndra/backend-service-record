package utils

import (
	"fmt"
	"net/http"
	"strings"
)

func HandleDBError(err error) (int, string) {
	msg := err.Error()

	if strings.Contains(msg, "Duplicate entry") {
		parts := strings.Split(msg, "'")
		if len(parts) > 1 {
			return http.StatusConflict, fmt.Sprintf("Data '%s' sudah terdaftar di sistem. Gunakan nilai lain.", parts[1])
		}
		return http.StatusConflict, "Data tersebut sudah ada di database."
	}

	if strings.Contains(msg, "a foreign key constraint fails") {
		return http.StatusBadRequest, "Gagal simpan: Data referensi (ID Bengkel atau No Polisi) tidak ditemukan."
	}
	return http.StatusInternalServerError, "Terjadi kesalahan internal pada server."
}
