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

<br>

## Modul PIC 2 — Master Data & Pelaporan Barang

Modul ini mencakup pengelolaan data master (kategori dan lokasi) serta siklus hidup pelaporan barang hilang (*lost*) dan ditemukan (*found*) beserta lampiran fotonya.

### Kontrak Endpoint

| Method | Endpoint | Akses | Deskripsi |
| :--- | :--- | :--- | :--- |
| `GET` | `/api/v1/categories` | Public / Auth | Menampilkan daftar seluruh kategori barang. |
| `GET` | `/api/v1/categories/:id` | Public / Auth | Menampilkan detail kategori berdasarkan ID. |
| `POST` | `/api/v1/categories` | Admin | Menambahkan kategori baru (`name` unik wajib diisi). |
| `PUT` | `/api/v1/categories/:id` | Admin | Memperbarui nama atau deskripsi kategori. |
| `DELETE` | `/api/v1/categories/:id` | Admin | Menghapus kategori (dilindungi constraint jika masih dipakai laporan). |
| `GET` | `/api/v1/locations` | Public / Auth | Menampilkan daftar seluruh lokasi kampus. |
| `GET` | `/api/v1/locations/:id` | Public / Auth | Menampilkan detail lokasi berdasarkan ID. |
| `POST` | `/api/v1/locations` | Admin | Menambahkan lokasi baru (`name` unik wajib diisi). |
| `PUT` | `/api/v1/locations/:id` | Admin | Memperbarui nama atau deskripsi lokasi. |
| `DELETE` | `/api/v1/locations/:id` | Admin | Menghapus lokasi (dilindungi constraint jika masih dipakai laporan). |
| `GET` | `/api/v1/reports` | Public / Auth | Daftar laporan dengan paginasi (`page`, `limit`) dan filter (`report_type`, `category_id`, `location_id`, `status`, `search`). |
| `GET` | `/api/v1/reports/:id` | Public / Auth | Detail lengkap laporan beserta lampiran foto (`item_images`). |
| `POST` | `/api/v1/reports` | Auth (Semua Role) | Membuat laporan *lost* atau *found* (status awal selalu `open`, dicatat ke `report_histories`). |
| `PUT` / `PATCH` | `/api/v1/reports/:id` | Owner / Admin | Memperbarui laporan (hanya dapat diubah saat status `open`). |
| `PATCH` | `/api/v1/reports/:id/close` | Owner / Admin | Menutup laporan (transisi status `open` → `closed`, dicatat ke `report_histories`). |
| `POST` | `/api/v1/reports/:id/images` | Owner / Admin | Mengunggah URL foto barang (maksimal 5 foto per laporan). |

### Aturan Bisnis & Integritas Data

1. **Integritas Master Data:** Nama kategori dan lokasi bersifat unik (`SQLSTATE 23505` menghasilkan respons `409 Conflict`). Data master yang masih memiliki relasi di tabel `item_reports` tidak dapat dihapus (`SQLSTATE 23503` menghasilkan respons `409 Conflict`).
2. **Audit Log Transaksi:** Pembuatan laporan baru (`open`) dan penutupan laporan (`closed`) dieksekusi dalam satu transaksi atomik database (`pgx.Tx`) bersama penyisipan riwayat ke `report_histories`.
3. **Keamanan & Otorisasi:** Modifikasi laporan (`update`, `close`, `add image`) dibatasi hanya untuk pemilik laporan (*owner*) atau user dengan role `admin` (akses tidak sah menghasilkan `403 Forbidden`).
4. **Pencegahan SQL Injection:** Seluruh query pencarian dan penyaringan paginasi menggunakan *parameterized queries* (`$1, $2, ...`) tanpa konkatenasi string input langsung.
