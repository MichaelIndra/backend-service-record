# Service Record API

Service Record API adalah layanan backend berbasis **Go (Golang)** menggunakan framework **Gin Gonic** dan **GORM** untuk mencatat riwayat servis kendaraan, manajemen bengkel, serta pelacakan penggantian sparepart lengkap dengan sistem penyimpanan dokumentasi foto/nota.

Lihat riwayat perubahan project di [CHANGELOG.md](CHANGELOG.md).

## Fitur Utama
- **Multi-Upload File**: Upload nota utama, foto komponen servis, dan suku cadang dalam satu request `multipart/form-data`.
- **Database Transaction**: Menjamin konsistensi data (jika upload file atau simpan detail gagal, database otomatis di-*rollback*).
- **Clean Asset Management**: Menghapus data transaksi otomatis akan membersihkan file gambar fisik di server tanpa menyisakan sampah penyimpanan.
- **Auto Upper-Case**: Otomatis mengubah input penting (No Polisi, No Nota, Kode Sparepart) menjadi huruf kapital melalui utility khusus.

---

## Panduan Instalasi (Lokal / Development)

### Prasyarat
- Go 1.20 atau versi di atasnya
- MySQL atau MariaDB

### Langkah-langkah
1. **Clone & Masuk ke Direktori Project**
   ```bash
   cd service-record