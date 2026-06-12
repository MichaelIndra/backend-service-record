package models

import (
	"time"
)

type Transaction struct {
	ID               uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	NoPolisi         string     `gorm:"not null;type:varchar(20)" json:"no_polisi"`
	NoNota           *string    `json:"no_nota"`
	Km               *string    `json:"km"`
	BengkelID        uint       `gorm:"not null" json:"bengkel_id"`
	BiayaBruto       float64    `gorm:"not null" json:"biaya_bruto"`
	Ppn              float64    `json:"ppn"`
	Diskon           float64    `json:"diskon"`
	BiayaNetto       float64    `gorm:"not null" json:"biaya_netto"`
	TanggalMasuk     *time.Time `json:"tanggal_masuk"`
	TanggalSelesai   *time.Time `json:"tanggal_selesai"`
	TanggalTransaksi time.Time  `gorm:"not null" json:"tanggal_transaksi"`
	FotoNota         []string   `gorm:"type:text;serializer:json" json:"foto_nota"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
	// Relasi
	DetailsService   []DetailService   `gorm:"foreignKey:TransactionID" json:"details_service"`
	DetailsSparepart []DetailSparepart `gorm:"foreignKey:TransactionID" json:"details_sparepart"`

	Kendaraan Kendaraan `gorm:"foreignKey:NoPolisi;references:NoPolisi" json:"kendaraan"`
	Bengkel   Bengkel   `gorm:"foreignKey:BengkelID" json:"Bengkel"`
}

type DetailService struct {
	ID            uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	TransactionID uint      `json:"transaction_id"`
	NamaService   string    `gorm:"not null" json:"nama_service"`
	Biaya         float64   `gorm:"not null" json:"biaya"`
	FotoService   string    `json:"foto_service"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type DetailSparepart struct {
	ID            uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	TransactionID uint      `json:"transaction_id"`
	NamaSparepart string    `gorm:"not null" json:"nama_sparepart"`
	KodeSparepart *string   `json:"kode_sparepart"`
	FotoSparepart string    `json:"foto_sparepart"`
	BatasKM       *string   `json:"batas_km"`
	BatasWaktu    *string   `json:"batas_waktu"`
	Remind        bool      `gorm:"default:false" json:"remind"`
	Qty           float64   `gorm:"not null" json:"qty"`
	HargaSatuan   float64   `gorm:"not null" json:"harga_satuan"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
