# Tugas Besar 2 IF2211 Strategi Algoritma 2024/2025

### *Pemanfaatan Algoritma BFS dan DFS dalam Pencarian Recipe pada Permainan Little Alchemy 2*

### Penjelasan singkat algoritma yang diimplementasikan
#### Breadth First Search (BFS)
Algoritma ini melakukan pencarian resep suatu elemen secara melebar dengan menelusuri semua simpul pada level yang sama terlebih dahulu sebelum lanjut ke level yang lebih dalam. Artinya, semua kemungkinan resep akan dicoba dari level bawah terlebih dahulu. Hal ini memastikan bahwa solusi pertama yang ditemukan adalah yang memiliki tier terendah. Algoritma ini juga dapat menemukan semua jalur pembuatan elemen karena setiap kombinasi bahan akan dievaluasi dan dikombinasikan secara eksplisit.

#### Depth First Search (DFS)
Algoritma ini melakukan pencarian resep suatu elemen dengan cara traversal setiap cabang pada simpul yang merupakan representasi alternatif resep dari suatu 
elemen hingga mencapai elemen dasar dengan aturan DFS yakni secara mendalam ke suatu cabang terlebih dahulu.
Apabila cabang tidak memenuhi aturan maka akan dilakukan backtracking ke cabang lain dari simpul
diatasnya. Algoritma ini diimplementasikan dengan memanfaatkan rekursif.


### Requirement
1. OS Windows atau Linux
2. Sudah terinstall bahasa pemrograman Go disarankan 1.24.3
3. Sudah terinstall goquery
4. Sudah terinstall Node.js disarankan v22.15.0.

### Cara menjalankan program
1. Clone repository
   ```sh
   git clone https://github.com/iqbalhaidr/Tubes2_TheAlchemist.git
   ```
2. Navigasi ke directory web
   ```sh
   cd ./Tubes2_TheAlchemistweb/web/my-app/
   ```
3. Unduh semua dependency
   ```sh
   npm install
   ```
4. Jalankan server
   ```sh
   npm start
   ```
5. Buka terminal baru dan navigasi ke parent directory
   ```sh
   cd ./Tubes2_TheAlchemistweb/
   ```
6. Jalankan program main.go
   ```sh
   go run main.go
   ```

### Anggota Kelompok
| NIM      | Nama                            |
| -------- | ------------------------------- |
| 13523081 | Jethro Jens Norbert Simatupang  |
| 13523111 | Muhammad Iqbal Haidar           |
| 15223090 | Ignacio Kevin Alberiann         |
