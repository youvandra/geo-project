---
title: "Module F: Technical GEO — Sitemap, API & LLM Crawling"
date: 2026-06-18
draft: false
tags: ["geo", "technical", "sitemap", "crawling", "modul-f"]
categories: ["module-f"]
---

Technical GEO adalah infrastruktur yang memudahkan AI crawlers mengakses dan memahami konten kita. Ini fondasi yang sering diabaikan.

<!--more-->

## Bagaimana AI Crawlers Bekerja?

Berbeda dengan search engine tradisional, AI crawlers (seperti GPTBot, ClaudeBot, Google-Extended) punya behavior khusus:

1. **Rate limited** — lebih lambat dari Googlebot
2. **Content-focused** — kurang tertarik pada CSS/JS
3. **Structure-aware** — sangat bergantung pada heading hierarchy
4. **Freshness-sensitive** — lebih suka konten terbaru

## Sitemap untuk AI Crawlers

### Best Practices

- **Referensi entity** di URL structure
- **Prioritaskan** halaman dengan authority tinggi
- **Lastmod** harus akurat (AI sensitive ke freshness)
- **Changefreq** sesuai update frequency

### Generate dengan `geo sitemap`

```bash
geo sitemap https://example.com ./blog/content -o sitemap.xml
```

Output:
```xml
<?xml version="1.0" encoding="UTF-8"?>
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
  <url>
    <loc>https://example.com/posts/hello-geo</loc>
    <lastmod>2026-06-18</lastmod>
    <changefreq>weekly</changefreq>
    <priority>0.7</priority>
  </url>
</urlset>
```

Tool ini otomatis:
- Scan content directory untuk file `.md`, `.html`
- Generate priority berdasarkan depth (top-level = 0.9, deeper = lower)
- Set lastmod dari file modification time
- Exclude direktori node_modules, .git, dll

## Robots.txt untuk AI

Tambahkan aturan spesifik untuk AI crawlers:

```
User-agent: GPTBot
Allow: /
Sitemap: https://example.com/sitemap.xml

User-agent: Google-Extended
Allow: /
Sitemap: https://example.com/sitemap.xml
```

## Content Delivery untuk LLM

### Struktur API Response
Jika punya API, format response yang LLM-friendly:

```json
{
  "title": "Judul Artikel",
  "content": "Isi artikel...",
  "entities": ["Entity1", "Entity2"],
  "last_updated": "2026-06-18",
  "source": "https://example.com/article"
}
```

### Schema.org sebagai API Alternatif
JSON-LD yang embed di halaman juga bisa berfungsi sebagai "API" untuk AI.

## Performance untuk AI Crawlers

- **First Contentful Paint < 1.5s** — AI crawlers juga punya timeout
- **Server-side rendering** — konten harus visible di HTML source
- **No JavaScript dependency** — jangan sembunyiin konten di JS

## Tools

| Tool | Fungsi |
|------|--------|
| `geo sitemap <url> <dir>` | Generate sitemap.xml |
| Google Search Console | Monitor crawling |
| Robots.txt Tester | Validasi akses crawler |

**Next:** Module G — Monitoring & Analytics AI Answers.
