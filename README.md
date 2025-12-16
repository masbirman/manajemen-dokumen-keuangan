# Manajemen Dokumen Keuangan

Sistem Manajemen Dokumen Keuangan - Full Web Application

## Tech Stack

- **Backend**: Go (Fiber framework)
- **Frontend**: Vue.js 3 + Vite + TypeScript + Tailwind CSS
- **Database**: PostgreSQL 16
- **State Management**: Pinia
- **HTTP Client**: Axios

## Quick Start with Docker

### Prerequisites

- Docker Desktop installed and running
- Git

### Running the Application

1. Start all services:

```bash
docker-compose up -d --build
```

2. Check service status:

```bash
docker-compose ps
```

3. View logs:

```bash
docker-compose logs -f
```

### Access Points

- **Frontend**: http://localhost:5173
- **Backend API**: http://localhost:8000
- **Health Check**: http://localhost:8000/api/health
- **PostgreSQL**: localhost:5432

### Verification Steps

After running `docker-compose up -d --build`, verify:

1. **All services are running**:

```bash
docker-compose ps
# Should show: db, backend, frontend all "Up"
```

2. **Backend API is accessible**:

```bash
curl http://localhost:8000/api/health
# Should return: {"status":"ok","message":"Dokumen Keuangan API is running","database":"connected"}
```

3. **Frontend is accessible**:

   - Open http://localhost:5173 in browser
   - Should see "Manajemen Dokumen Keuangan" page with API status

4. **Database is connected**:

```bash
docker-compose exec db psql -U dokumen_user -d dokumen_keuangan -c "\dt"
# Should connect successfully (empty table list is expected initially)
```

### Stopping Services

```bash
docker-compose down
```

### Stopping and Removing Volumes

```bash
docker-compose down -v
```

## Development

### Backend (Go)

The backend uses Air for hot reload. Any changes to `.go` files will automatically rebuild.

### Frontend (Vue.js)

The frontend uses Vite dev server with HMR. Changes are reflected immediately.

## Project Structure

```
├── backend/
│   ├── app/
│   │   ├── http/controllers/
│   │   ├── http/middleware/
│   │   ├── models/
│   │   ├── repositories/
│   │   └── services/
│   ├── config/
│   ├── database/
│   │   ├── migrations/
│   │   └── seeders/
│   ├── routes/
│   └── storage/
├── frontend/
│   ├── src/
│   │   ├── components/
│   │   ├── router/
│   │   ├── services/
│   │   ├── stores/
│   │   └── views/
│   └── public/
└── docker-compose.yml
```
