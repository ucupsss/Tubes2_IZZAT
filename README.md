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
