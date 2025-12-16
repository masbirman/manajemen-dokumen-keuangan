# Panduan Deployment ke Dokploy (VPS)

Dokumen ini berisi langkah-langkah untuk men-deploy aplikasi **Manajemen Dokumen Keuangan** ke server VPS menggunakan **Dokploy**.

## Persiapan File

Saya telah menyiapkan 4 file kunci untuk deployment production:

1.  `backend/Dockerfile.prod`: Konfigurasi build backend yang ringan (Go -> Alpine).
2.  `frontend/Dockerfile.prod`: Konfigurasi build frontend + Nginx server.
3.  `frontend/nginx.conf`: Konfigurasi server Nginx untuk handle routing Vue & Proxy API.
4.  `docker-compose.prod.yml`: File orkestrasi utama.

## Langkah-Langkah Deployment di Dokploy

### 1. Push Code ke GitHub
Pastikan semua file terbaru (termasuk file deployment yang baru dibuat) sudah di-push ke repository GitHub Anda.

```bash
git add .
git commit -m "Chore: Add deployment configuration files"
git push origin main
```

### 2. Setup Project di Dokploy

1.  Login ke Dashboard Dokploy Anda.
2.  Buat **Project Baru** (misal: `manajemen-dokumen`).
3.  Masuk ke Project tersebut.

### 3. Deploy Menggunakan Docker Compose (Recommended)

Cara termudah adalah menggunakan fitur "Compose" di Dokploy untuk men-handle stack (Frontend + Backend + DB) sekaligus.

1.  Klik **"Compose"** di menu service.
2.  **Service Name**: Bebas (misal: `main-stack`).
3.  **Repository**: Pilih repository GitHub project ini.
4.  **Branch**: `main`.
5.  **Compose Path**: Masukkan path ke file produksi kita:
    `./docker-compose.prod.yml`
    *(Dokploy akan membaca file ini alih-alih docker-compose.yml biasa)*.

### 4. Konfigurasi Environment (Opsional)

Jika Anda ingin mengubah password database atau secret key:
1.  Masuk ke tab **Environment**.
2.  Tambahkan variabel:
    ```env
    DB_USER=dokumen_user_prod
    DB_PASSWORD=password_super_rahasia
    DB_NAME=dokumen_keuangan_prod
    JWT_SECRET=rahasia_jwt_yang_panjang_dan_acak
    ```
    *(Pastikan value ini sama di service Backend dan DB)*.

### 5. Deploy

Klik tombol **Deploy**. Dokploy akan:
1.  Pull repository.
2.  Build image backend (Go).
3.  Build image frontend (Vue -> Nginx).
4.  Menjalankan database Postgres.
5.  Menjalankan semuanya dalam satu network.

### 6. Akses Domain

Setelah status "Running":
1.  Di Dokploy, buka service `frontend` (atau port 80 yang di-expose).
2.  Tambahkan **Domain** Anda (misal: `dokumen.domainanda.com`).
3.  Aktifkan **HTTPS/SSL** (Let's Encrypt) di tab Domain.
4.  Buka domain tersebut di browser.

Aplikasi Anda kini live di production! 🚀
