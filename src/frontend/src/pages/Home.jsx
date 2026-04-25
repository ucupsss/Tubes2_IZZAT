import BFS from "../assets/BFS.png";
import DFS from "../assets/DFS.png";

export default function Home() {
  return (
    <>
    <h1 className="page-title">Teori</h1>
    <div className="card">
        <h2>Mengenai Penelusuran Graf</h2>
        <p>
          Penelusuran Graf (Graph Traversal) merupakan sebuah prosedur sistematis untuk mengunjungi setiap simpul (vertex) dan memeriksa setiap sisi (edge) dalam suatu struktur data graf. Tujuan utama dari proses ini adalah untuk memastikan bahwa setiap bagian dari graf telah diakses setidaknya satu kali, yang kemudian menjadi fondasi bagi algoritma yang lebih kompleks seperti pencarian jalur, verifikasi konektivitas, dan analisis jaringan. Dua metode fundamental yang paling umum digunakan dalam penelusuran graf adalah Breadth-First Search (BFS) dan Depth-First Search (DFS), di mana keduanya memiliki karakteristik dan mekanisme kerja yang berbeda dalam mengeksplorasi struktur data tersebut.
        </p>
      </div>
      <div className="card">
        <h2>Mengenai BFS</h2>
        <p>
          Breadth-First Search (BFS) adalah algoritma penelusuran yang mengutamakan lebar atau cakupan horizontal dalam prosesnya. Algoritma ini bekerja dengan cara mengunjungi semua simpul tetangga dari simpul awal secara berurutan sebelum melanjutkan pencarian ke simpul di tingkat (level) berikutnya yang lebih dalam. Secara teknis, BFS menggunakan struktur data Queue yang bersifat First-In-First-Out (FIFO) untuk menyimpan urutan simpul yang akan dikunjungi. Karakteristik utama dari BFS adalah kemampuannya untuk menjamin penemuan jalur terpendek (shortest path) pada graf yang tidak memiliki bobot, karena simpul-simpul diperiksa berdasarkan jarak terdekatnya dari titik asal.
        </p>
        <img src={BFS} alt="BFS" style={{ width: "800px", height: "400px"}}/>
      </div>
           <div className="card">
        <h2>Mengenai DFS</h2>
        <p>
          Depth-First Search (DFS) adalah algoritma penelusuran yang memprioritaskan kedalaman pada saat penelusuran. Mekanisme kerja DFS adalah dengan menelusuri satu jalur cabang hingga mencapai titik terjauh (simpul daun) sebelum akhirnya melakukan proses backtracking untuk mengeksplorasi cabang lain yang belum dikunjungi. Dalam implementasinya, DFS memanfaatkan struktur data tumpukan (Stack) yang bersifat Last-In-First-Out (LIFO), baik secara eksplisit maupun melalui pemanggilan fungsi rekursif. DFS sangat efektif digunakan dalam penyelesaian masalah yang berkaitan dengan deteksi siklus, serta pencarian solusi dalam ruang pencarian yang luas seperti pada permainan labirin.
        </p>
        <img src={DFS} alt="DFS"  style={{ width: "800px", height: "400px"}} />
      </div>
    </>
  );
}