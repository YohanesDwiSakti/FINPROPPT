# Supabase Setup

Tambahkan environment variables ini di Vercel Project Settings:

```text
SUPABASE_URL=https://xxxx.supabase.co
SUPABASE_SERVICE_ROLE_KEY=isi_service_role_key
```

Jalankan SQL ini di Supabase SQL Editor:

```sql
create table if not exists public.app_users (
  id uuid primary key default gen_random_uuid(),
  name text not null,
  email text unique not null,
  password text not null,
  role text not null check (role in ('customer', 'admin')),
  created_at timestamptz not null default now()
);

create table if not exists public.manifests (
  id uuid primary key default gen_random_uuid(),
  receipt text unique not null,
  status text not null,
  location text,
  updated_by text,
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now()
);

insert into public.app_users (name, email, password, role)
values
  ('Admin Hub Denpasar', 'admin@tiki.test', 'admin123', 'admin'),
  ('Customer Demo', 'customer@tiki.test', 'customer123', 'customer')
on conflict (email) do nothing;
```

Catatan: password ini masih plaintext untuk demo. Untuk production, pindahkan auth ke Supabase Auth atau hash password di backend.
