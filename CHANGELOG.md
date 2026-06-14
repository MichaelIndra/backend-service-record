# Changelog

Semua perubahan penting pada project **Service Record API** akan dicatat di file ini.

## [1.0.1] - 2026-06-13

### Fixed
- **CORS**: update cors ketika dihit dari tempat lain

## [1.0.0] - 2026-06-12

### Added
- **Fitur Transaksi Lengkap (`POST /api/transaction`)**: Mendukung pembuatan transaksi baru dengan *multiple upload* untuk foto nota utama, foto detail service, dan foto detail sparepart secara bersamaan.
- **Fitur Clean Delete (`DELETE /api/transaction/:id`)**: Menghapus transaksi secara menyeluruh dari database (termasuk relasi detail service & sparepart) sekaligus menghapus folder fisik gambar di direktori `uploads/` secara otomatis (*Hard Delete*).
- **Fitur Filter Bengkel (`GET /api/bengkel/getBengkelByStatus/:status`)**: Mengambil daftar banyak bengkel (array/slice JSON) berdasarkan status aktif/tidak.
- **Dokumentasi Swagger (Go-Swagger)**: Menambahkan anotasi swagger pada endpoint transaksi dan bengkel untuk kemudahan integrasi frontend.

### Fixed
- **GORM AutoMigrate Batching Issue**: Memisahkan proses `AutoMigrate` menjadi 2 tahap (Tabel Master dahulu, disusul Tabel Transaksi) untuk menghindari error *Strict Foreign Key Checks* pada MySQL/MariaDB.
- **Custom Foreign Key Mapping**: Memperbaiki pemetaan relasi `Kendaraan` di struct `Transaction` menggunakan `references:NoPolisi` agar GORM tidak memaksa tipe data integer (`bigint`).
- **Windows Path Separator Bug**: Mengimplementasikan `filepath.ToSlash()` pada sistem upload gambar agar path yang tersimpan di database selalu menggunakan `/` (standar URL/Linux) meskipun server dijalankan di lingkungan Windows (`\\`).
- **Safe Date Pointer Validation**: Memperbaiki logika *parsing* tanggal opsional (`TanggalMasuk` dan `TanggalSelesai`). Jika frontend mengirimkan string kosong, database akan mengisinya dengan `NULL`, bukan nilai default Go (`0001-01-01`) yang memicu error database.