<!DOCTYPE html>
<html lang="id">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>TIKI Denpasar - Customer</title>
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
        }

        * { box-sizing: border-box; }
        body {
            margin: 0;
            font-family: "Plus Jakarta Sans", Arial, sans-serif;
            background: var(--soft);
            color: var(--ink);
        }

        .topbar {
            height: 78px;
            display: flex;
            align-items: center;
            justify-content: space-between;
            padding: 0 34px;
            background: var(--white);
            border-bottom: 1px solid var(--line);
            position: sticky;
            top: 0;
            z-index: 20;
        }

        .brand {
            display: flex;
            align-items: center;
            gap: 12px;
            color: var(--blue);
            font-size: 22px;
            font-weight: 800;
        }

        .brand-mark {
            width: 50px;
            height: 34px;
            border-radius: 8px;
            background: var(--blue);
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

        .nav-actions {
            display: flex;
            align-items: center;
            gap: 10px;
        }

        .session {
            color: var(--muted);
            font-size: 13px;
            font-weight: 800;
        }

        .logout-form { margin: 0; }
        .logout {
            height: 42px;
            padding: 0 14px;
            border: 1px solid var(--line);
            border-radius: 8px;
            background: var(--white);
            color: var(--blue);
            font: inherit;
            font-weight: 800;
            cursor: pointer;
        }

        .page {
            width: min(1180px, 100%);
            margin: 0 auto;
            padding: 28px 22px 56px;
        }

        .hero {
            display: grid;
            grid-template-columns: minmax(0, 1.05fr) minmax(360px, .95fr);
            min-height: 440px;
            border-radius: 8px;
            overflow: hidden;
            background: var(--blue);
            box-shadow: 0 28px 70px rgba(0, 71, 255, .18);
        }

        .hero-copy {
            padding: 42px;
            color: var(--white);
            display: flex;
            flex-direction: column;
            justify-content: center;
        }

        .tag {
            width: fit-content;
            background: var(--yellow);
            color: var(--ink);
            border-radius: 999px;
            padding: 8px 12px;
            font-size: 13px;
            font-weight: 800;
            margin-bottom: 18px;
        }

        h1, h2, h3 { margin-top: 0; letter-spacing: 0; }
        h1 {
            font-size: clamp(38px, 5vw, 62px);
            line-height: 1.02;
            margin-bottom: 16px;
        }

        .hero-copy p {
            color: rgba(255,255,255,.82);
            line-height: 1.75;
            margin: 0;
            max-width: 540px;
        }

        .hero-image {
            min-height: 440px;
            background:
                linear-gradient(90deg, rgba(0,71,255,.18), rgba(0,71,255,0)),
                url("https://images.unsplash.com/photo-1586528116311-ad8dd3c8310d?auto=format&fit=crop&w=1000&q=80") center/cover;
        }

        .tools {
            display: grid;
            grid-template-columns: 1.1fr .9fr;
            gap: 18px;
            margin-top: 20px;
        }

        .card {
            background: var(--white);
            border: 1px solid var(--line);
            border-radius: 8px;
            padding: 26px;
            box-shadow: 0 16px 42px rgba(16, 33, 63, .06);
        }

        .tabs {
            display: flex;
            gap: 10px;
            margin-bottom: 18px;
        }

        .tab {
            border: 1px solid var(--line);
            background: var(--soft);
            color: var(--muted);
            border-radius: 8px;
            padding: 12px 16px;
            font-weight: 800;
            cursor: pointer;
        }

        .tab.active {
            background: var(--blue);
            border-color: var(--blue);
            color: var(--white);
        }

        .panel { display: none; }
        .panel.active { display: block; }
        h2 { font-size: 24px; margin-bottom: 16px; }
        p { color: var(--muted); line-height: 1.65; }

        .grid {
            display: grid;
            grid-template-columns: repeat(2, minmax(0, 1fr));
            gap: 12px;
        }

        input {
            width: 100%;
            height: 54px;
            padding: 0 15px;
            border: 1px solid var(--line);
            border-radius: 8px;
            background: var(--soft);
            font: inherit;
            outline: none;
        }

        input:focus {
            border-color: var(--blue);
            background: var(--white);
            box-shadow: 0 0 0 4px rgba(0,71,255,.1);
        }

        .btn {
            width: 100%;
            height: 54px;
            border: 0;
            border-radius: 8px;
            background: var(--blue);
            color: var(--white);
            font-weight: 800;
            cursor: pointer;
            margin-top: 12px;
        }

        .result {
            display: none;
            margin-top: 18px;
            border: 1px solid var(--line);
            background: var(--soft);
            border-radius: 8px;
            padding: 18px;
        }

        .badge {
            display: inline-block;
            margin-bottom: 10px;
            padding: 7px 10px;
            border-radius: 8px;
            background: rgba(5, 150, 105, .12);
            color: var(--green);
            font-size: 12px;
            font-weight: 800;
        }

        .timeline {
            margin-top: 18px;
            padding-left: 18px;
            border-left: 3px solid var(--line);
        }

        .step { margin-bottom: 16px; }
        .step strong {
            display: block;
            color: var(--muted);
            font-size: 12px;
            margin-bottom: 4px;
        }

        .side-list {
            display: grid;
            gap: 12px;
        }

        .service {
            display: grid;
            grid-template-columns: 54px 1fr;
            gap: 14px;
            align-items: center;
            padding: 14px;
            border: 1px solid var(--line);
            border-radius: 8px;
            background: var(--soft);
        }

        .service-code {
            width: 54px;
            height: 44px;
            border-radius: 8px;
            display: grid;
            place-items: center;
            color: var(--white);
            background: var(--blue);
            font-weight: 800;
        }

        .service:nth-child(2) .service-code { background: var(--red); }
        .service:nth-child(3) .service-code { background: var(--yellow); color: var(--ink); }
        .service strong { display: block; margin-bottom: 3px; }
        .service span { color: var(--muted); font-size: 13px; }

        .branches {
            margin-top: 20px;
            display: grid;
            grid-template-columns: repeat(3, minmax(0, 1fr));
            gap: 14px;
        }

        .branch {
            background: var(--white);
            border: 1px solid var(--line);
            border-radius: 8px;
            padding: 20px;
        }

        .branch h3 { color: var(--blue); font-size: 17px; margin-bottom: 8px; }
        .branch p { margin: 0; font-size: 14px; }

        @media (max-width: 900px) {
            .hero, .tools { grid-template-columns: 1fr; }
            .hero-image { min-height: 260px; }
        }

        @media (max-width: 640px) {
            .topbar { padding: 0 18px; }
            .session { display: none; }
            .hero-copy, .card { padding: 22px; }
            .grid, .branches { grid-template-columns: 1fr; }
            .tabs { flex-direction: column; }
        }
    </style>
</head>
<body>
    <header class="topbar">
        <div class="brand">
            <span class="brand-mark" aria-hidden="true"></span>
            <span>TIKI DENPASAR</span>
        </div>
        <div class="nav-actions">
            <span class="session">Customer: {{ session('user_name', 'Akun Customer') }}</span>
            <form class="logout-form" method="POST" action="/logout">
                @csrf
                <button class="logout" type="submit">Logout</button>
            </form>
        </div>
    </header>

    <main class="page">
        <section class="hero">
            <div class="hero-copy">
                <div class="tag">#PAKETMUDUNIAKU</div>
                <h1>Kirim paket dari Denpasar jadi lebih gampang.</h1>
                <p>Lacak resi, cek estimasi ongkir, dan temukan cabang TIKI Bali dalam satu dashboard customer.</p>
            </div>
            <div class="hero-image" aria-hidden="true"></div>
        </section>

        <section class="tools">
            <div class="card">
                <div class="tabs" aria-label="Fitur pengiriman">
                    <button class="tab active" data-tab="track">Cek Resi</button>
                    <button class="tab" data-tab="rate">Cek Ongkir</button>
                </div>

                <section id="track" class="panel active">
                    <h2>Lacak Kiriman</h2>
                    <input id="receiptInput" type="text" placeholder="Masukkan nomor resi">
                    <button class="btn" id="trackButton">Lacak Resi</button>
                    <div class="result" id="trackingResult"></div>
                </section>

                <section id="rate" class="panel">
                    <h2>Tarif Pengiriman</h2>
                    <div class="grid">
                        <input id="originInput" type="text" placeholder="Kota asal" value="Denpasar">
                        <input id="destinationInput" type="text" placeholder="Kota tujuan">
                    </div>
                    <input id="weightInput" type="number" min="1" placeholder="Berat kg" style="margin-top: 12px;">
                    <button class="btn" id="rateButton">Cek Ongkir</button>
                    <div class="result" id="rateResult"></div>
                </section>
            </div>

            <aside class="card">
                <h2>Produk Populer</h2>
                <div class="side-list">
                    <div class="service">
                        <div class="service-code">ONS</div>
                        <div><strong>Over Night Service</strong><span>Estimasi tiba esok hari.</span></div>
                    </div>
                    <div class="service">
                        <div class="service-code">REG</div>
                        <div><strong>Regular Service</strong><span>Pengiriman harian paling fleksibel.</span></div>
                    </div>
                    <div class="service">
                        <div class="service-code">ECO</div>
                        <div><strong>Economy Service</strong><span>Ongkir hemat untuk paket ringan.</span></div>
                    </div>
                </div>
            </aside>
        </section>

        <section class="branches" aria-label="Cabang TIKI Bali">
            <article class="branch">
                <h3>TIKI Denpasar Hub</h3>
                <p>Jl. Kapten Cok Agung Tresna No.22, Dangin Puri Klod.</p>
            </article>
            <article class="branch">
                <h3>TIKI Teuku Umar</h3>
                <p>Jl. Teuku Umar No.200, Dauh Puri Kauh, Denpasar Barat.</p>
            </article>
            <article class="branch">
                <h3>TIKI Sanur</h3>
                <p>Jl. Danau Buyan No.74, Sanur, Denpasar Selatan.</p>
            </article>
        </section>
    </main>

    <script>
        const API_BASE = 'http://127.0.0.1:5000';

        document.querySelectorAll('.tab').forEach((tab) => {
            tab.addEventListener('click', () => {
                document.querySelectorAll('.tab').forEach((item) => item.classList.remove('active'));
                document.querySelectorAll('.panel').forEach((panel) => panel.classList.remove('active'));
                tab.classList.add('active');
                document.getElementById(tab.dataset.tab).classList.add('active');
            });
        });

        document.getElementById('trackButton').addEventListener('click', async () => {
            const receipt = document.getElementById('receiptInput').value.trim();
            const result = document.getElementById('trackingResult');
            if (!receipt) {
                alert('Masukkan nomor resi dulu.');
                return;
            }

            result.style.display = 'block';
            result.innerHTML = '<p>Memuat data tracking...</p>';
            try {
                const response = await fetch(`${API_BASE}/api/tracking/${encodeURIComponent(receipt)}`);
                const data = await response.json();
                if (!response.ok) throw new Error(data.message || 'Gagal melacak resi');
                result.innerHTML = `
                    <span class="badge">${data.status}</span>
                    <h3>RESI: ${data.receipt}</h3>
                    <p>Lokasi terakhir: ${data.location}. Estimasi tiba: ${data.estimate}.</p>
                    <div class="timeline">
                        ${data.timeline.map((step) => `<div class="step"><strong>${step.date}</strong>${step.status}</div>`).join('')}
                    </div>
                `;
            } catch (error) {
                result.innerHTML = `<p>${error.message}</p>`;
            }
        });

        document.getElementById('rateButton').addEventListener('click', async () => {
            const origin = document.getElementById('originInput').value.trim();
            const destination = document.getElementById('destinationInput').value.trim();
            const weight = document.getElementById('weightInput').value.trim();
            const result = document.getElementById('rateResult');

            if (!origin || !destination || !weight) {
                alert('Lengkapi kota asal, tujuan, dan berat.');
                return;
            }

            result.style.display = 'block';
            result.innerHTML = '<p>Menghitung ongkir...</p>';
            try {
                const params = new URLSearchParams({ origin, destination, weight });
                const response = await fetch(`${API_BASE}/api/rates?${params}`);
                const data = await response.json();
                if (!response.ok) throw new Error(data.message || 'Gagal menghitung ongkir');
                result.innerHTML = `
                    <span class="badge">${data.service}</span>
                    <h3>Rp ${data.price.toLocaleString('id-ID')}</h3>
                    <p>${data.origin} ke ${data.destination}, ${data.weight_kg} kg. Estimasi ${data.estimate}.</p>
                `;
            } catch (error) {
                result.innerHTML = `<p>${error.message}</p>`;
            }
        });
    </script>
</body>
</html>
