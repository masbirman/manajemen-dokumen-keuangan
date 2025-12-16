# Panduan Deployment ke Dokploy (Production)

Dokumen ini berisi langkah-langkah deployment aplikasi **Manajemen Dokumen Keuangan** ke server Dokploy dengan arsitektur **Terpisah** (App & Database).

## Arsitektur
- **Service 1 (Database)**: PostgreSQL (Managed by Dokploy).
- **Service 2 (Application)**: Frontend + Backend (Managed by Docker Compose).

## Langkah 1: Persiapan Database

1.  Buka Project Anda di Dokploy.
2.  Klik **"Create Service"** -> Pilih **"Database"** -> **PostgreSQL**.
3.  Beri nama (misal: `dokumen-db`).
4.  Setelah dibuat, masuk ke tab **Environment** atau detail database untuk mendapatkan kredensial:
    - **Host** (Internal Name, misal: `postgresql-dokumen-db` atau IP internal)
    - **Database Name**
    - **User**
    - **Password**
    - **Port** (Default: 5432)

## Langkah 2: Deployment Aplikasi

1.  Kembali ke Project, klik **"Create Service"** -> Pilih **"Compose"**.
2.  Beri nama (misal: `dokumen-app`).
3.  **Repository**: Hubungkan ke GitHub repo ini.
4.  **Branch**: `main`
5.  **Compose Path**: `docker-compose.prod.yml`
6.  Masuk ke tab **Environment** dan masukkan kredensial Database dari Langkah 1:
    ```env
    DB_HOST=nama_internal_host_database
    DB_PORT=5432
    DB_USER=user_database_anda
    DB_PASSWORD=password_database_anda
    DB_NAME=nama_database_anda
    JWT_SECRET=buat_random_string_panjang
    ```

## Langkah 3: Deploy & Domain

1.  Klik **Deploy** pada Service Aplikasi.
2.  Sistem otomatis akan:
    - Build Backend (Go).
    - Build Frontend (Vue).
    - Setup Nginx & Routing.
3.  **Domain Otomatis**:
    File `docker-compose.prod.yml` sudah dilengkapi label Traefik untuk domain `dokumen.keudisdiksulteng.web.id`.
    - Pastikan DNS A Record domain tersebut sudah mengarah ke IP Server Dokploy Anda.
    - Tunggu beberapa saat, Traefik akan otomatis generate SSL certificate.

Selesai! Aplikasi Anda berjalan aman dengan database terpisah. 🚀
