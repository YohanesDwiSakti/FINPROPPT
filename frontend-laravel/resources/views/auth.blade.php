<!DOCTYPE html>
<html lang="id">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>TIKI Denpasar - Masuk</title>
    <link href="https://fonts.googleapis.com/css2?family=Plus+Jakarta+Sans:wght@500;600;700;800&display=swap" rel="stylesheet">
    <style>
        :root {
            --blue: #0047ff;
            --red: #e31b2f;
            --yellow: #ffd400;
            --ink: #10213f;
            --muted: #667085;
            --line: #dce6f2;
            --soft: #f5f8fc;
            --white: #ffffff;
            --danger: #be123c;
        }

        * { box-sizing: border-box; }

        body {
            margin: 0;
            min-height: 100vh;
            font-family: "Plus Jakarta Sans", Arial, sans-serif;
            background:
                radial-gradient(circle at top left, rgba(255, 212, 0, .18), transparent 28%),
                linear-gradient(160deg, #f7faff 0%, #edf4ff 100%);
            color: var(--ink);
        }

        .page {
            min-height: 100vh;
            display: grid;
            place-items: center;
            padding: 24px;
        }

        .auth-shell {
            width: min(460px, 100%);
            text-align: center;
        }

        .logo-wrap {
            display: grid;
            justify-items: center;
            gap: 18px;
            margin-bottom: 30px;
        }

        .logo-mark {
            width: 96px;
            height: 68px;
            border-radius: 18px;
            background: var(--blue);
            position: relative;
            box-shadow: 0 20px 44px rgba(0, 71, 255, .24);
        }

        .logo-mark::before {
            content: "";
            position: absolute;
            left: 18px;
            top: 16px;
            width: 34px;
            height: 34px;
            border-radius: 10px;
            background: var(--yellow);
            box-shadow: 28px 0 0 var(--red);
        }

        .brand-title {
            margin: 0;
            font-size: clamp(32px, 8vw, 46px);
            line-height: 1;
            letter-spacing: 0;
            font-weight: 800;
            color: var(--blue);
        }

        .brand-subtitle {
            margin: 0;
            color: var(--muted);
            font-size: 15px;
            line-height: 1.6;
        }

        .entry-card,
        .form-card {
            background: var(--white);
            border: 1px solid var(--line);
            border-radius: 18px;
            padding: 26px;
            box-shadow: 0 24px 70px rgba(16, 33, 63, .1);
        }

        .entry-actions {
            display: grid;
            gap: 12px;
        }

        .entry-btn,
        .submit {
            width: 100%;
            height: 56px;
            border: 0;
            border-radius: 12px;
            font: inherit;
            font-weight: 800;
            cursor: pointer;
        }

        .entry-btn.primary,
        .submit {
            background: var(--blue);
            color: var(--white);
            box-shadow: 0 16px 32px rgba(0, 71, 255, .22);
        }

        .entry-btn.secondary {
            background: var(--soft);
            color: var(--blue);
            border: 1px solid var(--line);
        }

        .hint {
            margin: 18px 0 0;
            color: var(--muted);
            font-size: 13px;
            line-height: 1.55;
        }

        .form-card {
            display: none;
            text-align: left;
        }

        .form-card.active { display: block; }
        .entry-card.hidden { display: none; }

        .form-head {
            display: flex;
            align-items: flex-start;
            justify-content: space-between;
            gap: 16px;
            margin-bottom: 22px;
        }

        h2 {
            margin: 0 0 8px;
            font-size: 28px;
            line-height: 1.15;
            letter-spacing: 0;
        }

        .form-head p {
            margin: 0;
            color: var(--muted);
            line-height: 1.55;
            font-size: 14px;
        }

        .back-btn {
            border: 1px solid var(--line);
            background: var(--soft);
            color: var(--blue);
            border-radius: 10px;
            padding: 9px 12px;
            font: inherit;
            font-size: 13px;
            font-weight: 800;
            cursor: pointer;
            white-space: nowrap;
        }

        .alert, .note {
            padding: 13px 14px;
            border-radius: 12px;
            font-size: 14px;
            font-weight: 700;
            margin-bottom: 14px;
            line-height: 1.5;
        }

        .alert { background: #fff1f2; color: var(--danger); border: 1px solid #fecdd3; }
        .note { background: #eff6ff; color: var(--blue); border: 1px solid #bfdbfe; }

        label {
            display: block;
            margin: 14px 0 8px;
            font-size: 13px;
            font-weight: 800;
        }

        input, select {
            width: 100%;
            height: 54px;
            padding: 0 15px;
            border: 1px solid var(--line);
            border-radius: 12px;
            background: var(--soft);
            color: var(--ink);
            font: inherit;
            outline: none;
        }

        input:focus, select:focus {
            border-color: var(--blue);
            background: var(--white);
            box-shadow: 0 0 0 4px rgba(0, 71, 255, .1);
        }

        .submit {
            margin-top: 20px;
        }

        @media (max-width: 520px) {
            .page { padding: 18px; }
            .entry-card, .form-card { padding: 22px; }
            .logo-mark { width: 84px; height: 60px; }
            .logo-mark::before {
                left: 16px;
                top: 14px;
                width: 30px;
                height: 30px;
                box-shadow: 24px 0 0 var(--red);
            }
            .form-head {
                flex-direction: column-reverse;
            }
        }
    </style>
</head>
<body>
    <main class="page">
        <section class="auth-shell">
            <div class="logo-wrap">
                <div class="logo-mark" aria-hidden="true"></div>
                <h1 class="brand-title">TIKI DENPASAR</h1>
            </div>

            <div class="entry-card" id="entryCard">
                <div class="entry-actions">
                    <button class="entry-btn primary" type="button" data-open="loginCard">Login</button>
                    <button class="entry-btn secondary" type="button" data-open="registerCard">Register</button>
                </div>
            </div>

            <div class="form-card" id="loginCard">
                <div class="form-head">
                    <div>
                        <h2>Masuk Akun</h2>
                    </div>
                    <button class="back-btn" type="button">Kembali</button>
                </div>

                @if ($mode === 'login' && $errors->any())
                    <div class="alert">{{ $errors->first() }}</div>
                @endif

                <form method="POST" action="/login">
                    @csrf
                    <label for="login-email">Email</label>
                    <input id="login-email" name="email" type="email" placeholder="nama@email.com" value="{{ old('email') }}" required>

                    <label for="login-password">Password</label>
                    <input id="login-password" name="password" type="password" placeholder="Masukkan password" required>

                    <label for="login-role">Role Akses</label>
                    <select id="login-role" name="role" required>
                        <option value="customer" @selected(old('role') === 'customer')>Customer</option>
                        <option value="admin" @selected(old('role') === 'admin')>Admin</option>
                    </select>

                    <button class="submit" type="submit">Login</button>
                </form>
            </div>

            <div class="form-card" id="registerCard">
                <div class="form-head">
                    <div>
                        <h2>Register Akun</h2>
                    </div>
                    <button class="back-btn" type="button">Kembali</button>
                </div>

                @if ($mode === 'register' && $errors->any())
                    <div class="alert">{{ $errors->first() }}</div>
                @endif

                <form method="POST" action="/register">
                    @csrf
                    <label for="register-name">Nama</label>
                    <input id="register-name" name="name" type="text" placeholder="Nama lengkap" value="{{ old('name') }}" required>

                    <label for="register-email">Email</label>
                    <input id="register-email" name="email" type="email" placeholder="nama@email.com" value="{{ old('email') }}" required>

                    <label for="register-password">Password</label>
                    <input id="register-password" name="password" type="password" placeholder="Masukkan password" required>

                    <button class="submit" type="submit">Register</button>
                </form>
            </div>
        </section>
    </main>

    <script>
        const entryCard = document.getElementById('entryCard');
        const cards = document.querySelectorAll('.form-card');

        function clearAuthFields() {
            ['login-email', 'login-password', 'register-name', 'register-email', 'register-password'].forEach((id) => {
                const input = document.getElementById(id);
                if (input) input.value = '';
            });

            const role = document.getElementById('login-role');
            if (role) role.value = 'customer';
        }

        function showCard(id) {
            clearAuthFields();
            entryCard.classList.add('hidden');
            cards.forEach((card) => card.classList.remove('active'));
            document.getElementById(id).classList.add('active');
        }

        function showEntry() {
            clearAuthFields();
            cards.forEach((card) => card.classList.remove('active'));
            entryCard.classList.remove('hidden');
        }

        document.querySelectorAll('[data-open]').forEach((button) => {
            button.addEventListener('click', () => showCard(button.dataset.open));
        });

        document.querySelectorAll('.back-btn').forEach((button) => {
            button.addEventListener('click', showEntry);
        });

        @if ($errors->any())
            showCard('{{ $mode === 'register' ? 'registerCard' : 'loginCard' }}');
        @endif
    </script>
</body>
</html>
