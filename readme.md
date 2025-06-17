# 📦 FM Technical Test API

Project ini merupakan backend RESTful API yang dibuat dengan **Go menggunakan Echo Framework**. API ini dibuat untuk keperluan technical test, dengan arsitektur bertingkat: **handler → service → repository**. 
Fitur utama mencakup autentikasi JWT, absensi karyawan (clock-in dan clock-out), serta upload gambar.

---

## 📚 Dokumentasi API

Swagger UI tersedia di: 

`` http://localhost:50001/swagger/index.html ``


Untuk mengakses endpoint yang dilindungi, gunakan token JWT dengan format:

``Authorization: Bearer <token_kamu>``


Kamu bisa klik tombol **Authorize** di Swagger UI untuk memasukkan token tersebut.

---

## 📁 Struktur Folder

- `app/` – Konfigurasi aplikasi dan inisialisasi database
- `handler/` – HTTP handler untuk Echo
- `service/` – Logika bisnis
- `repository/` – Akses database (PostgreSQL)
- `upload/image/` – Folder penyimpanan file gambar yang diunggah (nanti akan terbuat saat upload image)

---

## 🧭 Daftar Endpoint API

### 🔐 `/api/v1/user` – Autentikasi User

| Endpoint | Method | Deskripsi |
|----------|--------|-----------|
| `/register` | POST | Mendaftarkan user baru sekaligus data karyawan. |
| `/login`    | POST | Login dan mendapatkan token JWT untuk akses endpoint selanjutnya. |

📌 **Kenapa seperti ini?**
- Proses registrasi dibuat gabungan karena user dan karyawan didaftarkan bersamaan.
- JWT digunakan agar sistem tetap stateless dan mudah diintegrasikan ke frontend/mobile.

---

### 🕒 `/api/v1/employee` – Absensi Karyawan

> Semua endpoint berikut **butuh autentikasi** JWT.

| Endpoint | Method | Deskripsi                                                |
|----------|--------|----------------------------------------------------------|
| `/clock_in`  | POST | Mengirim data clock-in termasuk lokasi dan foto absensi. |
| `/clock_out` | POST | Mengirim data clock-out dengan lokasi dan foto absensi saat keluar.      |
| `/absence/log` | GET | Mengambil riwayat absensi user yang sedang login.        |

📌 **Kenapa seperti ini?**
- Lokasi dipakai untuk validasi kehadiran karyawan di area(radius) kerja yang ditentukan.
- Token JWT digunakan untuk mengidentifikasi siapa yang melakukan clock-in/out tanpa perlu kirim email di request.

---

### 🖼 `/api/v1/image` – Upload Gambar

> Semua endpoint berikut **butuh autentikasi** JWT.

| Endpoint | Method | Deskripsi |
|----------|--------|-----------|
| `/save` | POST | Mengunggah gambar menggunakan `multipart/form-data`, menyertakan `jenis` (contoh: avatar, checkin, dll). |

📌 **Kenapa seperti ini?**
- File dikelompokkan berdasarkan jenis-nya agar rapi dalam folder `upload/image/{jenis}/...`.
- File yang berhasil diupload bisa langsung diakses secara publik via URL statis.

---

## 🌐 Akses File Statis

Setelah gambar berhasil diunggah, kamu bisa akses file-nya di:

`` http://{{url}}/uploads/avatar/1721251820.jpg ``
