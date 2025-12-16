# Panduan Deployment ke Dokploy (Production)

Dokumen ini berisi langkah-langkah deployment aplikasi **Manajemen Dokumen Keuangan** ke server Dokploy dengan arsitektur **Terpisah** (App & Database).

## Arsitektur
- **Service 1 (Database)**: PostgreSQL (Managed by Dokploy).
- **Service 2 (Application)**: Frontend + Backend (Managed by Docker Compose).

---

## Langkah 1: Persiapan Database

1.  Buka Project Anda di Dokploy.
2.  Klik **"Create Service"** -> Pilih **"Database"** -> **PostgreSQL**.
3.  Beri nama (misal: `dokumen-db`).
4.  Setelah dibuat, masuk ke tab **Environment** atau detail database untuk mendapatkan kredensial:
    - **Host** (Internal Name, misal: `postgresql-dokumen-db-xxxxx`)
    - **Database Name**
    - **User**
    - **Password**
    - **Port** (Default: 5432)

> ⚠️ **PENTING**: Catat **Internal Host Name** dari database, bukan External IP.

---

## Langkah 2: Deployment Aplikasi

1.  Kembali ke Project, klik **"Create Service"** -> Pilih **"Compose"**.
2.  Beri nama (misal: `dokumen-app`).
3.  **Repository**: Hubungkan ke GitHub repo ini.
4.  **Branch**: `main`
5.  **Compose Path**: `docker-compose.prod.yml`
6.  Masuk ke tab **Environment** dan masukkan kredensial Database dari Langkah 1:

    ```env
    # Database Configuration (WAJIB DIISI!)
    DB_HOST=nama_internal_host_database
    DB_PORT=5432
    DB_USERNAME=user_database_anda
    DB_PASSWORD=password_database_anda
    DB_DATABASE=nama_database_anda
    
    # Security
    JWT_SECRET=buat_random_string_panjang_minimal_32_karakter
    ```

> ⚠️ **PERHATIAN Environment Variable Names:**
> - Gunakan `DB_USERNAME` (bukan `DB_USER`)
> - Gunakan `DB_DATABASE` (bukan `DB_NAME`)
> - Nama variable ini harus sesuai dengan yang dibaca oleh backend (`config.go`)

---

## Langkah 3: Deploy & Domain

1.  Klik **Deploy** pada Service Aplikasi.
2.  Sistem otomatis akan:
    - Build Backend (Go) dengan binary optimized.
    - Build Frontend (Vue) dengan Nginx.
    - Setup internal networking & Routing.
3.  **Domain Otomatis**:
    File `docker-compose.prod.yml` sudah dilengkapi label Traefik untuk domain `dokumen.keudisdiksulteng.web.id`.
    - Pastikan DNS A Record domain tersebut sudah mengarah ke IP Server Dokploy Anda.
    - Tunggu beberapa saat, Traefik akan otomatis generate SSL certificate.

---

## Troubleshooting

### Error: "hostname resolving error"
**Penyebab**: `DB_HOST` tidak valid atau database tidak berada di network yang sama.

**Solusi**:
1. Pastikan database service sudah running di Dokploy.
2. Gunakan **Internal Host Name** dari database (bukan External IP).
3. Verifikasi kedua service menggunakan network `dokploy-network`.

### Error: "404 page not found"
**Penyebab**: Frontend nginx tidak bisa reach backend atau routing rusak.

**Solusi**:
1. Pastikan backend container sudah running (cek logs).
2. Verifikasi environment variables sudah benar.
3. Re-deploy aplikasi untuk rebuild containers.

### Error: "failed to connect to database"
**Penyebab**: Kredensial database salah atau port tidak tepat.

**Solusi**:
1. Periksa kembali environment variables:
   - `DB_USERNAME` (bukan `DB_USER`)
   - `DB_DATABASE` (bukan `DB_NAME`)
2. Pastikan password tidak mengandung karakter khusus yang perlu di-escape.
3. Verifikasi port database (default PostgreSQL: 5432).

---

## Health Check Endpoints

Untuk verifikasi service berjalan dengan baik:

- **Backend**: `GET /api/health`
- **Frontend/Nginx**: `GET /nginx-health`

---

## File Konfigurasi

| File | Deskripsi |
|------|-----------|
| `docker-compose.prod.yml` | Orchestration production |
| `backend/Dockerfile.prod` | Multi-stage build untuk Go backend |
| `frontend/Dockerfile.prod` | Multi-stage build untuk Vue + Nginx |
| `frontend/nginx.conf` | Nginx config dengan proxy ke backend |

---

Selesai! Aplikasi Anda berjalan aman dengan database terpisah. 🚀
