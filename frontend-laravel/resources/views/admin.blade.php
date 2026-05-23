<!DOCTYPE html>
<html lang="id">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <meta name="csrf-token" content="{{ csrf_token() }}">
    <title>Admin TIKI Denpasar</title>
    <link href="https://fonts.googleapis.com/css2?family=Plus+Jakarta+Sans:wght@500;600;700;800&display=swap" rel="stylesheet">
    <style>
        :root {
            --blue:#0047ff;
            --blue-dark:#0737b6;
            --red:#e31b2f;
            --yellow:#ffd400;
            --ink:#10213f;
            --muted:#667085;
            --line:#dce6f2;
            --soft:#f5f8fc;
            --white:#fff;
            --green:#059669;
            --amber:#d97706;
        }

        * { box-sizing: border-box; }
        body {
            margin: 0;
            font-family: "Plus Jakarta Sans", Arial, sans-serif;
            background: var(--soft);
            color: var(--ink);
        }

        .layout {
            min-height: 100vh;
            display: grid;
            grid-template-columns: 260px 1fr;
        }

        .sidebar {
            background: var(--blue);
            color: var(--white);
            padding: 28px 22px;
            display: flex;
            flex-direction: column;
            gap: 24px;
        }

        .brand {
            display: flex;
            gap: 12px;
            align-items: center;
            font-size: 20px;
            font-weight: 800;
        }

        .brand-mark {
            width: 48px;
            height: 34px;
            border-radius: 8px;
            background: var(--white);
            position: relative;
        }

        .brand-mark::before {
            content: "";
            position: absolute;
            left: 8px;
            top: 7px;
            width: 20px;
            height: 20px;
            border-radius: 6px;
            background: var(--yellow);
            box-shadow: 16px 0 0 var(--red);
        }

        .role {
            color: rgba(255,255,255,.72);
            font-size: 12px;
            font-weight: 800;
            text-transform: uppercase;
            margin-top: 8px;
        }

        .nav {
            display: grid;
            gap: 8px;
        }

        .nav a, .logout {
            color: var(--white);
            text-decoration: none;
            border: 1px solid rgba(255,255,255,.22);
            border-radius: 8px;
            padding: 13px 14px;
            font: inherit;
            font-weight: 800;
            background: rgba(255,255,255,.09);
            cursor: pointer;
            text-align: left;
        }

        .nav a.active { background: var(--white); color: var(--blue); }
        .logout-form { margin-top: auto; }
        .logout { width: 100%; text-align: center; }

        .main {
            padding: 30px;
        }

        .topline {
            display: flex;
            align-items: center;
            justify-content: space-between;
            gap: 18px;
            margin-bottom: 20px;
        }

        h1, h2, h3 { margin-top: 0; letter-spacing: 0; }
        h1 { font-size: 32px; margin-bottom: 8px; }
        p { color: var(--muted); line-height: 1.65; margin: 0; }

        .operator {
            background: var(--white);
            border: 1px solid var(--line);
            border-radius: 8px;
            padding: 12px 14px;
            color: var(--muted);
            font-size: 13px;
            font-weight: 800;
        }

        .stats {
            display: grid;
            grid-template-columns: repeat(4, minmax(0, 1fr));
            gap: 14px;
            margin-bottom: 18px;
        }

        .stat, .card {
            background: var(--white);
            border: 1px solid var(--line);
            border-radius: 8px;
            box-shadow: 0 16px 42px rgba(16, 33, 63, .06);
        }

        .stat { padding: 18px; }
        .stat span {
            display: block;
            color: var(--muted);
            font-size: 12px;
            font-weight: 800;
            text-transform: uppercase;
            margin-bottom: 8px;
        }
        .stat strong { font-size: 30px; }

        .workspace {
            display: grid;
            grid-template-columns: minmax(0, 1.1fr) minmax(320px, .9fr);
            gap: 18px;
        }

        .card { padding: 24px; }
        .badge {
            display: inline-block;
            background: rgba(0,71,255,.09);
            color: var(--blue);
            padding: 7px 10px;
            border-radius: 8px;
            font-size: 12px;
            font-weight: 800;
            margin-bottom: 12px;
        }

        label {
            display: block;
            margin: 16px 0 8px;
            font-size: 13px;
            font-weight: 800;
        }

        input, select {
            width: 100%;
            height: 54px;
            padding: 0 15px;
            border: 1px solid var(--line);
            border-radius: 8px;
            background: var(--soft);
            font: inherit;
            outline: none;
        }

        input:focus, select:focus {
            border-color: var(--blue);
            background: var(--white);
            box-shadow: 0 0 0 4px rgba(0,71,255,.1);
        }

        .grid { display: grid; grid-template-columns: 1fr 1fr; gap: 12px; }

        .btn {
            width: 100%;
            height: 54px;
            border: 0;
            border-radius: 8px;
            background: var(--blue);
            color: var(--white);
            font-weight: 800;
            margin-top: 18px;
            cursor: pointer;
        }

        .message {
            display: none;
            margin-top: 16px;
            padding: 14px;
            border-radius: 8px;
            background: #ecfdf5;
            border: 1px solid #bbf7d0;
            color: var(--green);
            font-weight: 800;
            line-height: 1.5;
        }

        .queue {
            display: grid;
            gap: 12px;
        }

        .queue-row {
            display: grid;
            grid-template-columns: 1fr auto;
            gap: 12px;
            padding: 14px;
            border: 1px solid var(--line);
            border-radius: 8px;
            background: var(--soft);
        }

        .queue-row strong { display: block; margin-bottom: 4px; }
        .pill {
            align-self: start;
            border-radius: 8px;
            padding: 6px 9px;
            background: rgba(217,119,6,.12);
            color: var(--amber);
            font-size: 12px;
            font-weight: 800;
        }

        @media (max-width: 940px) {
            .layout, .workspace, .stats { grid-template-columns: 1fr; }
            .sidebar { min-height: auto; }
        }

        @media (max-width: 640px) {
            .main { padding: 20px; }
            .topline { align-items: flex-start; flex-direction: column; }
            .grid { grid-template-columns: 1fr; }
        }
    </style>
</head>
<body>
    <div class="layout">
        <aside class="sidebar">
            <div>
                <div class="brand">
                    <span class="brand-mark" aria-hidden="true"></span>
                    <span>TIKI ADMIN</span>
                </div>
                <div class="role">Hub Denpasar</div>
            </div>
            <nav class="nav" aria-label="Menu admin">
                <a class="active" href="/admin">Manifest</a>
                <a href="#queue">Prioritas</a>
                <a href="#summary">Ringkasan</a>
            </nav>
            <form class="logout-form" method="POST" action="/logout">
                @csrf
                <button class="logout" type="submit">Logout</button>
            </form>
        </aside>

        <main class="main">
            <div class="topline">
                <div>
                    <h1>Manifest Operasional</h1>
                    <p>Update status paket masuk, sortir, kurir, DEX, dan delivered untuk area Bali.</p>
                </div>
                <div class="operator">Admin: {{ session('user_name', 'ADMIN HUB DPS') }}</div>
            </div>

            <section class="stats" id="summary" aria-label="Ringkasan operasional">
                <div class="stat"><span>Paket Masuk</span><strong>128</strong></div>
                <div class="stat"><span>Sortir</span><strong>34</strong></div>
                <div class="stat"><span>Kurir</span><strong>76</strong></div>
                <div class="stat"><span>DEX</span><strong>9</strong></div>
            </section>

            <section class="workspace">
                <div class="card">
                    <span class="badge">FORM INTERNAL</span>
                    <h2>Update Paket Tunggal</h2>
                    <p>Data yang disimpan dari form ini dikirim ke backend Go lokal.</p>

                    <form id="manifestForm">
                        <label for="receipt">Nomor Resi</label>
                        <input id="receipt" type="text" placeholder="Contoh: TKI-DEN-2026001" required>

                        <div class="grid">
                            <div>
                                <label for="status">Status</label>
                                <select id="status">
                                    <option value="Arrived">Tiba di Hub Denpasar</option>
                                    <option value="Sorting">Sedang Disortir</option>
                                    <option value="With Courier">Dibawa Kurir</option>
                                    <option value="DEX">DEX, dikirim besok</option>
                                    <option value="Delivered">Diterima</option>
                                </select>
                            </div>
                            <div>
                                <label for="location">Area Lokasi</label>
                                <input id="location" type="text" placeholder="Cth: Denpasar Barat">
                            </div>
                        </div>

                        <button class="btn" type="submit">Simpan Update Manifest</button>
                    </form>
                    <div class="message" id="message"></div>
                </div>

                <aside class="card" id="queue">
                    <span class="badge">MONITORING</span>
                    <h2>Antrian Prioritas</h2>
                    <div class="queue">
                        <div class="queue-row">
                            <div><strong>TKI-DEN-2026001</strong><p>Sanur, dispatch sebelum 16:00</p></div>
                            <span class="pill">Kurir</span>
                        </div>
                        <div class="queue-row">
                            <div><strong>TKI-DEN-2026048</strong><p>Teuku Umar, alamat perlu validasi</p></div>
                            <span class="pill">DEX</span>
                        </div>
                        <div class="queue-row">
                            <div><strong>TKI-DEN-2026087</strong><p>Denpasar Barat, sortir ulang</p></div>
                            <span class="pill">Sortir</span>
                        </div>
                    </div>
                </aside>
            </section>
        </main>
    </div>

    <script>
        document.getElementById('manifestForm').addEventListener('submit', async (event) => {
            event.preventDefault();
            const message = document.getElementById('message');
            message.style.display = 'block';
            message.textContent = 'Menyimpan update paket...';
            message.style.background = '#ecfdf5';
            message.style.borderColor = '#bbf7d0';
            message.style.color = '#059669';

            try {
                const response = await fetch('/manifests', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                        'X-CSRF-TOKEN': document.querySelector('meta[name="csrf-token"]').content,
                    },
                    body: JSON.stringify({
                        receipt: document.getElementById('receipt').value,
                        status: document.getElementById('status').value,
                        location: document.getElementById('location').value
                    })
                });
                const data = await response.json();
                if (!response.ok) throw new Error(data.message || 'Gagal menyimpan update');
                message.textContent = data.message;
                event.target.reset();
            } catch (error) {
                message.textContent = error.message;
                message.style.background = '#fff1f2';
                message.style.borderColor = '#fecdd3';
                message.style.color = '#be123c';
            }
        });
    </script>
</body>
</html>
