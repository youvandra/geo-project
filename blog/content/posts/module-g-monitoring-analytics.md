---
title: "Module G: Monitoring & Analytics untuk GEO"
date: 2026-06-18
draft: false
tags: ["geo", "monitoring", "tracking", "analytics", "modul-g"]
categories: ["module-g"]
---

Kita sudah bikin konten yang dioptimasi untuk GEO. Tapi bagaimana kita tahu strategi kita berhasil? Tanpa monitoring, kita buta.

<!--more-->

## Apa yang Perlu Dimonitor?

GEO monitoring berbeda dengan SEO analytics. Metrik yang relevan:

### 1. Citation Rate
Seberapa sering konten kita muncul sebagai sumber di AI answers?

- **Manual check:** "site:example.com" + query di ChatGPT/Perplexity
- **Automated:** Belum ada tool publik, tapi bisa pakai API

### 2. Pageview Trends
Gunakan `geo tracker` untuk monitor pageview sebagai proxy interest:

```bash
geo tracker "Generative Engine Optimization" --history
```

Output:
```
Topic:              Generative Engine Optimization
Trend:              increasing (+145 views)
Current Daily Views: 3219
```

Tool ini menyimpan history di `~/.geo-tracker/`, jadi kita bisa lihat tren dari waktu ke waktu.

### 3. Entity Visibility
Apakah entity kita muncul di Knowledge Graph? Cek dengan:

```bash
geo entity "Nama Brand/Topik kita" --json
```

### 4. Content Score
Pantau skor GEO konten kita secara berkala:

```bash
cat artikel-updated.md | geo score --stdin
```

## Metrik GEO

| Metrik | Alat Ukur | Target |
|--------|-----------|--------|
| RAG Score | `geo score` | ≥ 70/100 |
| Entity Count | `geo entity` | ≥ 5 per artikel |
| Pageviews | `geo tracker` | Stabil/meningkat |
| Schema Coverage | Manual | ≥ 1 schema per halaman |
| Citation Rate | Manual check | Muncul di AI answers |

## Tracking dengan `geo tracker`

Tool tracker kita menggunakan Wikipedia pageview API sebagai proxy untuk mengukur interest/tren suatu topik.

**Fitur:**
- Cek pageviews harian
- Simpan history lokal
- Analisis tren (increasing/stable/decreasing)
- Output JSON untuk integrasi

```bash
geo tracker "Search engine optimization" --json --history
```

## Dashboard Monitoring

Backend dashboard juga menyediakan interface visual:

1. Buka `http://localhost:8080/tracker`
2. Masukkan topik
3. Lihat pageview & tren langsung

## Weekly GEO Audit Checklist

- [ ] Cek citation rate di ChatGPT + Perplexity (manual sample)
- [ ] Score konten baru dengan `geo score`
- [ ] Track pageviews topik kunci dengan `geo tracker`
- [ ] Update entity linking di artikel existing
- [ ] Review schema markup masih valid

**Next:** Module H — Advanced GEO.
