Panduan Lengkap Instalasi & Deployment
Dokumentasi ini mencakup langkah-langkah setup dari nol, mulai dari instalasi Go di komputer lokal hingga proses deployment ke server production Linux menggunakan PM2 dan Nginx.

BAGIAN 1: Setup di Komputer Lokal (Development)
1. Install Go (Golang)
Unduh installer Go versi terbaru (minimal v1.20) di situs resmi: golang.org/dl.

Jalankan installer tersebut sesuai OS Anda (Windows/macOS/Linux) dan ikuti petunjuk Next sampai selesai.

Verifikasi instalasi dengan membuka Terminal / Command Prompt, lalu ketik:

Bash
go version
2. Jalankan Project Pertama Kali
Masuk ke folder project Anda melalui terminal:

Bash
cd path/to/service-record
Unduh semua dependency/library Go yang dibutuhkan aplikasi:

Bash
go mod tidy
Buat file .env di root folder dan sesuaikan DSN database lokal Anda:

Code snippet
PORT=8080
DB_DSN=root:password_lokal@tcp(127.0.0.1:3306)/service_record_dev?charset=utf8mb4&parseTime=True&loc=Local
Jalankan aplikasi langsung untuk testing:

Bash
go run main.go
BAGIAN 2: Setup Server Linux Mentah (Production)
Hubungkan SSH ke server Anda (misal Ubuntu Server), lalu lakukan instalasi software dasar berikut:

1. Update Server & Install Node.js + PM2
Aplikasi Go kita akan diatur oleh PM2, sehingga server membutuhkan Node.js.

Bash
# Update package list Linux
sudo apt update && sudo apt upgrade -y

# Install Node.js (Versi LTS)
curl -fsSL https://deb.nodesource.com/setup_20.x | sudo -E bash -
sudo apt install -y nodejs

# Install PM2 secara global
sudo npm install pm2 -g
2. Install & Setup MySQL/MariaDB (Jika DB di Server yang Sama)
Jika database Anda berada di server yang berbeda, lewati langkah ini.

Bash
sudo apt install mariadb-server -y

# Amankan instalasi database (Atur password root di sini)
sudo mysql_secure_installation

# Masuk ke MySQL untuk membuat database aplikasi
sudo mysql -u root -p
Di dalam prompt MySQL, jalankan perintah berikut:

SQL
CREATE DATABASE service_record_prod;
EXIT;
3. Install Nginx (Web Server / Reverse Proxy)
Bash
sudo apt install nginx -y
sudo systemctl start nginx
sudo systemctl enable nginx
BAGIAN 3: Proses Build & Deploy ke Server
Langkah 1: Compile Code di Komputer Lokal (Windows)
Jangan lakukan go build di server karena akan memakan CPU server. Lakukan Cross-Compile dari laptop Windows Anda agar menghasilkan file binary Linux:

Bash
# Jalankan ini di Terminal CMD/PowerShell laptop Anda:
set GOOS=linux
set GOARCH=amd64
go build -o service-record-api main.go
Hasilnya adalah file tanpa ekstensi bernama service-record-api.

Langkah 2: Upload File ke Server
Gunakan aplikasi FTP seperti WinSCP atau FileZilla. Upload file-file berikut ke direktori server Anda (Misal: /var/www/service-record/):

File binary service-record-api

File ecosystem.config.js

File .env.production

Langkah 3: Atur Hak Akses Folder di Server Linux
Masuk kembali ke SSH server, lalu masuk ke folder project dan berikan hak akses agar server bisa mengeksekusi file binary dan membuat folder foto:

Bash
cd /var/www/service-record/

# Buat folder untuk asset foto dan log PM2 jika belum ada
mkdir -p uploads logs

# Berikan izin eksekusi pada binary Go
chmod +x service-record-api

# Berikan izin baca-tulis penuh pada folder uploads dan logs
chmod -R 775 uploads logs
Langkah 4: Konfigurasi .env.production
Pastikan isi file .env.production di server sudah mengarah ke database production yang benar:

Code snippet
PORT=8080
DB_DSN=user_prod:password_prod@tcp(127.0.0.1:3306)/service_record_prod?charset=utf8mb4&parseTime=True&loc=Local
Langkah 5: Jalankan Aplikasi Menggunakan PM2
Nyalakan aplikasi Anda di server dengan perintah berikut:

Bash
# Jalankan menggunakan profil production
pm2 start ecosystem.config.js --env production

# Agar PM2 otomatis hidup kembali jika server tiba-tiba mati lampu/reboot
pm2 startup
⚠️ Penting: Setelah mengetik pm2 startup, Linux akan memunculkan satu baris perintah panjang yang diawali kata sudo env PATH.... Copy baris tersebut, paste di terminal, lalu tekan Enter.

Kunci konfigurasi PM2 saat ini agar permanen:

Bash
pm2 save
BAGIAN 4: Konfigurasi Nginx (Membuka Akses Publik)
Agar frontend bisa menembak API dan membaca gambar di folder uploads, kita harus mengonfigurasi Nginx sebagai jembatan.

Buat file konfigurasi baru di Nginx:

Bash
sudo nano /etc/nginx/sites-available/service-record
Paste konfigurasi berikut (Sesuaikan server_name dengan domain atau IP Server Anda):

Nginx
server {
    listen 80;
    server_name api.bengkelkamu.com; # Ubah ke domain/IP server

    # Jalur proxy ke aplikasi Go (Gin)
    location / {
        proxy_pass http://127.0.0.1:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    }

    # Jalur statis untuk mengakses foto/nota langsung dari browser/frontend
    location /uploads/ {
        alias /var/www/service-record/uploads/;
        expires 30d;
        add_header Cache-Control "public";
    }
}
Simpan file (Ctrl+O, lalu Ctrl+X).

Aktifkan konfigurasi dan restart Nginx:

Bash
sudo ln -s /etc/nginx/sites-available/service-record /etc/nginx/sites-enabled/
sudo nginx -t # Pastikan muncul tulisan "syntax is ok"
sudo systemctl restart nginx
BAGIAN 5: Lembar Ringkasan Perintah (Cheat Sheet)
Mengelola Aplikasi (PM2)
Melihat log/print data dari Go secara live: pm2 logs service-record-api

Melihat status & penggunaan RAM/CPU: pm2 status

Restart aplikasi (Lakukan ini setiap kali Anda mengganti file binary baru): pm2 restart service-record-api

Mematikan aplikasi temporer: pm2 stop service-record-api

Mengelola Web Server (Nginx)
Cek error konfigurasi Nginx: sudo nginx -t

Restart Nginx: sudo systemctl restart nginx

Melihat log error Nginx jika gambar tidak muncul: sudo tail -f /var/log/nginx/error.log