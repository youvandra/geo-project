---
title: "Module B: Topic Research & Entity Clustering untuk GEO"
date: 2026-06-18
draft: false
tags: ["geo", "topic-research", "entities", "modul-b"]
categories: ["module-b"]
---

Topic research untuk GEO berbeda dengan SEO. Kita perlu memahami **entity relationships** dan **semantic clusters**, bukan sekadar keyword density.

<!--more-->

## Entity-Based Research

AI search engines menggunakan **Knowledge Graphs** (Google KG, Wikidata) untuk memahami hubungan antar konsep. Konten yang terhubung dengan entity graph lebih mungkin dikutip.

### Cara Riset Topik untuk GEO:

1. **Start dengan core topic** — apa yang mau kita bahas?
2. **Expand ke entities** — entity apa yang terkait?
3. **Identify questions** — pertanyaan apa yang sering muncul?
4. **Map relationships** — bagaimana entity saling terhubung?

## Praktik dengan `geo topic`

Tool `geo topic` kita bisa digunakan untuk menganalisis suatu topik:

```bash
geo topic "Generative Engine Optimization" --json
```

Outputnya akan menunjukkan:
- **Description** — ringkasan topik dari Wikipedia
- **Entities** — entity terkait (links dari Wikipedia)
- **Questions** — pertanyaan GEO yang relevan
- **Subtopics** — sub-topik dari struktur artikel

### Contoh Analisis

Untuk topik "Search Engine Optimization", entity yang muncul antara lain:

- Google
- Bing
- PageRank
- Backlink
- Featured snippets
- BERT
- Semantic search

Entity-entity ini bisa jadi **content pillars** untuk strategi GEO kita.

## Topic Clusters untuk GEO

Buat struktur topik seperti ini:

```
Core Topic: GEO
├── Sub: RAG Fundamentals
│   ├── Entity: Retrieval-Augmented Generation
│   ├── Entity: Vector embeddings
│   └── Entity: Context window
├── Sub: Content Optimization
│   ├── Entity: Structured data
│   ├── Entity: Schema markup
│   └── Entity: Citation signals
├── Sub: Monitoring
│   ├── Entity: Citation tracking
│   ├── Entity: AI answer audit
│   └── Entity: Visibility score
```

Semakin dalam entity relationship-nya, semakin besar kemungkinan AI mengutip konten kita.

## Tools untuk Riset

| Tool | Fungsi |
|------|--------|
| `geo topic <topic>` | Analisis topik + entity |
| Wikipedia API | Sumber entity & deskripsi |
| Wikidata Query | Entity relationship graph |

**Next:** Module C — Bikin Konten yang RAG-Friendly.
