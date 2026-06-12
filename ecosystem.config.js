const fs = require('fs');
const path = require('path');

// Fungsi cerdas untuk membaca file .env tanpa perlu install library tambahan
function parseEnvFile(fileName) {
  const filePath = path.join(__dirname, fileName);
  if (!fs.existsSync(filePath)) {
    return {};
  }

  const content = fs.readFileSync(filePath, 'utf-8');
  const envObj = {};

  content.split('\n').forEach(line => {
    const trimmed = line.trim();
    // Lewati jika baris kosong atau berupa komentar (#)
    if (!trimmed || trimmed.startsWith('#')) return;

    // Split berdasarkan tanda '=' pertama untuk memisahkan Key dan Value
    const [key, ...valueParts] = trimmed.split('=');
    if (key) {
      // Gabungkan kembali value (jika di dalam value ada tanda '=' seperti pada DSN)
      let value = valueParts.join('=').trim();
      // Hapus tanda kutip jika pembungkus string ada di file .env
      value = value.replace(/^["']|["']$/g, '');
      
      envObj[key.trim()] = value;
    }
  });

  return envObj;
}

module.exports = {
  apps: [
    {
      name: 'service-record-api',
      script: './service-record-api',
      exec_mode: 'fork',
      instances: 1,
      watch: false,
      max_memory_restart: '250M',
      error_file: './logs/err.log',
      out_file: './logs/out.log',
      log_date_format: 'YYYY-MM-DD HH:mm:ss Z',

      // 🟢 OTOMATIS Membaca dari file .env
      env: parseEnvFile('.env'),

      // 🟢 OTOMATIS Membaca dari file .env.production
      env_production: parseEnvFile('.env.production')
    }
  ]
};