package models

import "time"

type Kendaraan struct {
	NoPolisi   string    `gorm:"primaryKey;type:varchar(20);not null" json:"no_polisi"`
	Nama       string    `gorm:"not null" json:"nama"`
	Manufaktur string    `gorm:"not null" json:"manufaktur"`
	Jenis      string    `gorm:"not null" json:"jenis"`
	Tahun      string    `gorm:"not null" json:"tahun"`
	Warna      string    `gorm:"not null" json:"warna"`
	Aktif      bool      `gorm:"default:true" json:"aktif"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `Json:"updated_at"`
}
