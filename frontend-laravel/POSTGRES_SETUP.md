# PostgreSQL Setup

Mesin ini belum punya `psql` dan PHP belum memuat ekstensi `pdo_pgsql`, jadi migrasi PostgreSQL belum bisa dijalankan dari sini. Setelah PostgreSQL dan ekstensi PHP `pdo_pgsql` aktif, pakai langkah ini.

1. Buat database:

```sql
CREATE DATABASE finproppt;
```

2. Salin konfigurasi:

```powershell
Copy-Item .env.postgres.example .env
php artisan key:generate
```

3. Sesuaikan `DB_USERNAME` dan `DB_PASSWORD` di `.env`.

4. Jalankan migrasi dan seed admin privat:

```powershell
php artisan migrate --seed
```

Seeder membuat akun privat:

- Admin: `admin@tiki.test` / `admin123`
- Customer demo: `customer@tiki.test` / `customer123`

Register publik hanya membuat akun `customer`; akun `admin` harus dibuat dari seeder/database.
