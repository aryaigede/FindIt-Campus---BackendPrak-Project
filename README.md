FindIt-Campus---BackendPrak-Project
FindIt Campus - Sistem Pelaporan Barang Hilang dan Ditemukan di Lingkungan Kampus

## Aktor & Hak Akses

| Role                     | Hak Akses Utama                                                                                                                 |
| :----------------------- | :------------------------------------------------------------------------------------------------------------------------------ |
| **Mahasiswa / Pengguna** | Registrasi & login; kelola profil; buat, ubah, atau batalkan laporan miliknya; cari laporan; ajukan klaim; pantau status klaim. |
| **Petugas Keamanan**     | Lihat laporan; tinjau klaim; setujui/tolak klaim disertai catatan; catat serah terima barang.                                   |
| **Admin**                | Kelola pengguna & role; kelola kategori & lokasi; moderasi laporan; akses seluruh data dan riwayat status.                      |

<br>

## Fitur Utama

1. **Autentikasi & Profil:** Registrasi dan login berbasis JWT beserta manajemen profil. _(Implementasi utuh pada tahap boilerplate)_
2. **Laporan Barang:** Pembuatan laporan _lost_ atau _found_ lengkap dengan kategori, lokasi, tanggal kejadian, dan foto.
3. **Pencarian & Filter:** Daftar laporan dengan filter jenis laporan, kategori, lokasi, dan status (`open`, `claimed`, `returned`, `closed`).
4. **Detail Laporan:** Informasi lengkap barang beserta lampiran foto.
5. **Pengajuan Klaim:** Pengajuan klaim oleh pengguna pada laporan _found_ berstatus `open` dengan bukti kepemilikan.
6. **Verifikasi Klaim:** Peninjauan klaim oleh petugas (`approved` / `rejected`) beserta catatan verifikasi.
7. **Serah Terima Barang:** Pencatatan petugas penyerah, penerima, waktu pengembalian, dan catatan.
8. **Riwayat Status:** Audit log setiap perubahan status laporan beserta aktor yang melakukannya.

<br>

## Rancangan Basis Data

| Tabel              | Kolom Utama                                                                                                                                                 | Relasi & Keterangan                                                                                  |
| :----------------- | :---------------------------------------------------------------------------------------------------------------------------------------------------------- | :--------------------------------------------------------------------------------------------------- |
| `roles`            | `id`, `name`, `created_at`, `updated_at`                                                                                                                    | Menyimpan daftar role sistem (`admin`, `security`, `user`).                                          |
| `users`            | `id`, `role_id`, `name`, `email`, `password_hash`, `student_id`, `phone`, `is_active`, `created_at`, `updated_at`                                           | `role_id` → `roles.id`.                                                                              |
| `categories`       | `id`, `name`, `description`, `created_at`, `updated_at`                                                                                                     | Data master kategori barang.                                                                         |
| `locations`        | `id`, `name`, `description`, `created_at`, `updated_at`                                                                                                     | Data master lokasi kampus.                                                                           |
| `item_reports`     | `id`, `user_id`, `category_id`, `location_id`, `report_type`, `title`, `description`, `custody_note`, `event_date`, `status`, `created_at`, `updated_at`    | Relasi ke `users`, `categories`, dan `locations`; `custody_note` mencatat posisi pengamanan barang.  |
| `item_images`      | `id`, `report_id`, `image_url`, `created_at`                                                                                                                | `report_id` → `item_reports.id`.                                                                     |
| `claims`           | `id`, `report_id`, `claimant_id`, `reviewed_by`, `proof_description`, `proof_image_url`, `status`, `review_note`, `reviewed_at`, `created_at`, `updated_at` | Relasi ke `item_reports` dan `users`; dibatasi 1 klaim `approved` per laporan.                       |
| `item_returns`     | `id`, `report_id`, `claim_id`, `handed_over_by`, `received_by`, `returned_at`, `notes`                                                                      | `report_id` → `item_reports.id` (unik), `claim_id` → `claims.id` (opsional/unik untuk alur `found`). |
| `report_histories` | `id`, `report_id`, `actor_id`, `old_status`, `new_status`, `note`, `created_at`                                                                             | Mencatat transisi status pada `item_reports`.                                                        |
