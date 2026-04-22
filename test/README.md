# Test Case DOM Traversal

Folder ini berisi bahan uji untuk frontend dan backend.

## Cara uji manual lewat frontend

1. Jalankan backend:

```bash
go run ./src/backend
```

2. Jalankan frontend:

```bash
cd src/frontend
npm.cmd run dev
```

3. Buka halaman traversal, pilih input `HTML`, lalu isi dengan konten dari `test/sample-dom.html`.
4. Coba selector pada tabel berikut.

| No | Selector | Algoritma | Limit | Expected matched | Yang diuji |
| --- | --- | --- | --- | ---: | --- |
| 1 | `*` | BFS | Semua | 29 | Universal selector |
| 2 | `div` | BFS | Semua | 2 | Tag selector |
| 3 | `.panel` | DFS | Semua | 2 | Class selector |
| 4 | `#content` | BFS | Semua | 1 | ID selector |
| 5 | `button.primary` | DFS | Semua | 1 | Tag + class |
| 6 | `.btn.primary` | BFS | Semua | 1 | Multi-class |
| 7 | `button[type=submit]` | BFS | Semua | 1 | Tag + attribute |
| 8 | `[data-kind=promo]` | DFS | Semua | 1 | Attribute tanpa tag |
| 9 | `nav .external` | BFS | Semua | 1 | Descendant combinator |
| 10 | `.card .badge` | DFS | Semua | 2 | Descendant combinator |
| 11 | `main > section` | BFS | Semua | 2 | Child combinator |
| 12 | `ul > li` | DFS | Semua | 4 | Child combinator |
| 13 | `li + li` | BFS | Semua | 3 | Adjacent sibling combinator |
| 14 | `li ~ li` | DFS | Semua | 3 | General sibling combinator |
| 15 | `.item` | BFS | Top 2 | 2 | Limit hasil |

## Cara uji otomatis lewat API

Pastikan backend sudah berjalan di `http://localhost:5175`, lalu dari root repository jalankan:

```powershell
powershell -ExecutionPolicy Bypass -File .\test\run-api-tests.ps1
```

Jika backend berjalan di alamat lain:

```powershell
powershell -ExecutionPolicy Bypass -File .\test\run-api-tests.ps1 -ApiBase "http://localhost:5175"
```
