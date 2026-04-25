# Tubes2_IZZAT
Tugas Besar 2 Strategi Algoritma Pemanfaatan Algoritma BFS dan DFS dalam Mekanisme Penelusuran CSS pada pohon Document Object Model 

## Cara Menjalankan

Jalankan backend dari root repository:

```bash
go run ./src/backend
```

Backend berjalan di `http://localhost:5175` dan menyediakan endpoint `POST /api/traversal`.

Jalankan frontend pada terminal lain:

```bash
cd src/frontend
npm install
npm run dev
```

Secara default frontend akan memanggil backend di `http://localhost:5175`. Jika backend dijalankan pada alamat lain, set environment variable `VITE_API_BASE_URL`.

## Menjalankan dengan Docker

Requirement:

- Docker Desktop atau Docker Engine dengan Docker Compose

Jalankan dari root repository:

```bash
docker compose up --build
```

Setelah container aktif:

- Frontend: `http://localhost:8080`
- Backend API: `http://localhost:5175`
- Health check backend: `http://localhost:5175/api/health`

Untuk menghentikan container:

```bash
docker compose down
```

Setup ini sudah cukup untuk packaging lokal dan kebutuhan bonus Docker pada laporan, serta bisa dijadikan dasar untuk deployment ke VM.
