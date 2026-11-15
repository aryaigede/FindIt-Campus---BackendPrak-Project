FindIt-Campus---BackendPrak-Project
FindIt Campus - Sistem Pelaporan Barang Hilang dan Ditemukan di Lingkungan Kampus

## Aktor & Hak Akses

| Role | Hak Akses Utama |
| :--- | :--- |
| **Mahasiswa / Pengguna** | Registrasi & login; kelola profil; buat, ubah, atau batalkan laporan miliknya; cari laporan; ajukan klaim; pantau status klaim. |
| **Petugas Keamanan** | Lihat laporan; tinjau klaim; setujui/tolak klaim disertai catatan; catat serah terima barang. |
| **Admin** | Kelola pengguna & role; kelola kategori & lokasi; moderasi laporan; akses seluruh data dan riwayat status. |

<br>

## Fitur Utama

1. **Autentikasi & Profil:** Registrasi dan login berbasis JWT beserta manajemen profil. *(Implementasi utuh pada tahap boilerplate)*
2. **Laporan Barang:** Pembuatan laporan *lost* atau *found* lengkap dengan kategori, lokasi, tanggal kejadian, dan foto.
3. **Pencarian & Filter:** Daftar laporan dengan filter jenis laporan, kategori, lokasi, dan status (`open`, `claimed`, `returned`, `closed`).
4. **Detail Laporan:** Informasi lengkap barang beserta lampiran foto.
5. **Pengajuan Klaim:** Pengajuan klaim oleh pengguna pada laporan *found* berstatus `open` dengan bukti kepemilikan.
6. **Verifikasi Klaim:** Peninjauan klaim oleh petugas (`approved` / `rejected`) beserta catatan verifikasi.
7. **Serah Terima Barang:** Pencatatan petugas penyerah, penerima, waktu pengembalian, dan catatan.
8. **Riwayat Status:** Audit log setiap perubahan status laporan beserta aktor yang melakukannya.

<br>

## Fitur Utama

1. **Autentikasi & Profil:** Registrasi dan login berbasis JWT beserta manajemen profil. *(Implementasi utuh pada tahap boilerplate)*
2. **Laporan Barang:** Pembuatan laporan *lost* atau *found* lengkap dengan kategori, lokasi, tanggal kejadian, dan foto.
3. **Pencarian & Filter:** Daftar laporan dengan filter jenis laporan, kategori, lokasi, dan status (`open`, `claimed`, `returned`, `closed`).
4. **Detail Laporan:** Informasi lengkap barang beserta lampiran foto.
5. **Pengajuan Klaim:** Pengajuan klaim oleh pengguna pada laporan *found* berstatus `open` dengan bukti kepemilikan.
6. **Verifikasi Klaim:** Peninjauan klaim oleh petugas (`approved` / `rejected`) beserta catatan verifikasi.
7. **Serah Terima Barang:** Pencatatan petugas penyerah, penerima, waktu pengembalian, dan catatan.
8. **Riwayat Status:** Audit log setiap perubahan status laporan beserta aktor yang melakukannya.

<br>
