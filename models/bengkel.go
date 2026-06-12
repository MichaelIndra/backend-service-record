package models

import "time"

type Bengkel struct {
	ID           uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Nama         string    `gorm:"not null" json:"nama"`
	JenisService string    `gorm:"not null" json:"jenis_service"`
	Tipe         string    `gorm:"type:varchar(50);not null;default:'MOBIL BBM'" json:"type" example:"MOBIL BBM"`
	NoTelp       *string   `json:"no_telp"`
	Provinsi     string    `gorm:"not null" json:"provinsi"`
	KotaKab      string    `gorm:"not null" json:"kota_kab"`
	Kecamatan    string    `gorm:"not null" json:"kecamatan"`
	Kelurahan    string    `gorm:"not null" json:"kelurahan"`
	Alamat       string    `gorm:"not null" json:"alamat"`
	Aktif        bool      `gorm:"default:true" json:"aktif"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
