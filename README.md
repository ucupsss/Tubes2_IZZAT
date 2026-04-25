# Tubes2_IZZAT

Tugas Besar 2 Strategi Algoritma Pemanfaatan Algoritma BFS dan DFS dalam Mekanisme Penelusuran CSS pada pohon Document Object Model.

## Penjelasan Singkat Algoritma

### BFS

Breadth First Search (BFS) melakukan penelusuran pohon secara melebar, yaitu memeriksa semua node pada suatu level terlebih dahulu sebelum melanjutkan ke level berikutnya. Pada project ini, BFS digunakan untuk mencari node DOM yang sesuai dengan selector CSS dengan pendekatan level-order traversal.

### DFS

Depth First Search (DFS) melakukan penelusuran pohon secara mendalam, yaitu menelusuri satu cabang hingga sedalam mungkin sebelum kembali ke cabang sebelumnya. Pada project ini, DFS digunakan untuk mencari node DOM yang sesuai dengan selector CSS dengan pendekatan depth-first traversal.

## Requirement

Program ini terdiri dari backend dan frontend, sehingga kebutuhan utamanya adalah:

- Go untuk menjalankan backend
- Node.js dan npm untuk menjalankan frontend
- Browser modern untuk membuka antarmuka web

Versi yang disarankan:

- Go 1.22 atau lebih baru
- Node.js 20 atau lebih baru
- npm 10 atau lebih baru

## Instalasi

### Backend

Dari root repository:

```bash
go run ./src/backend
```

Backend akan berjalan di `http://localhost:5175` dan menyediakan endpoint `POST /api/traversal`.

### Frontend

Masuk ke folder frontend lalu install dependency:

```bash
cd src/frontend
npm install
```

## Cara Menjalankan Program

Jalankan backend dari root repository:

```bash
go run ./src/backend
```

Jalankan frontend pada terminal lain:

```bash
cd src/frontend
npm run dev
```

Secara default frontend akan memanggil backend di `http://localhost:5175`. Jika backend dijalankan pada alamat lain, atur environment variable `VITE_API_BASE_URL`.

## Command Build dan Compile

### Menjalankan backend

```bash
go run ./src/backend
```

### Menjalankan frontend dalam mode development

```bash
cd src/frontend
npm run dev
```

### Build frontend

```bash
cd src/frontend
npm run build
```
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

## Author

| NIM | Nama | Pembagian Tugas |
| --- | --- | --- |
|13524014  |Yusuf Faishal Listyardi  |Implementasi algoritma BFS&DFS, Implementasi DOM Parser dan graph, Implementasi LCA binary lifting, Integrasi dengan docker  |
|13524081  |Alya Nur Rahmah  |Frontend, Implementasi animasi pada penelusuran pohon, Implementasi multithreading  |
|13524137  |Reysha Syafitri Mulya Ramadhan  |Implementasi web scrapping dan CSS selector, editor video  |
