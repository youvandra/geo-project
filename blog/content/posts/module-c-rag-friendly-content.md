---
title: "Module C: Bikin Konten yang RAG-Friendly"
date: 2026-06-18
draft: false
tags: ["geo", "rag", "content-structure", "modul-c"]
categories: ["module-c"]
---

Retrieval-Augmented Generation (RAG) adalah cara AI search mengambil informasi. Konten yang dioptimasi untuk RAG punya struktur khusus yang memudahkan AI mengekstrak jawaban.

<!--more-->

## Prinsip RAG-Friendly Content

### 1. Hierarchical Structure
Gunakan heading yang jelas (H1 > H2 > H3). AI menggunakan heading sebagai **navigation signal** untuk memahami topik per section.

```markdown
# Topik Utama
## Sub-topik A
### Detail A1
### Detail A2
## Sub-topik B
```

### 2. Direct Q&A Format
AI search sering menampilkan **direct answers**. Format Q&A membantu:

```markdown
## Apa itu GEO?
Generative Engine Optimization adalah...

## Bagaimana cara kerja GEO?
GEO bekerja dengan...
```

### 3. Comprehensive Coverage
Konten yang tipis (thin content) jarang dikutip. Pastikan konten mencakup:

- Definisi
- Cara kerja
- Manfaat
- Tantangan
- Best practices
- Contoh nyata

### 4. Entity-Rich Writing
Sertakan entity names dalam teks, bukan cuma kata ganti. AI menggunakan entity sebagai **signal relevance**.

❌ "Platform ini membantu optimasi..."
✅ "Google's AI Overviews membantu optimasi konten untuk generative search..."

## Scoring dengan `geo score`

Gunakan tool kita untuk mengukur RAG-friendliness:

```bash
cat artikel.md | geo score --stdin
```

Output menunjukkan skor di 6 dimensi:

| Dimensi | Skor Maks | Target Ideal |
|---------|-----------|--------------|
| Structure | 25 | ≥ 18 |
| Q&A Coverage | 20 | ≥ 12 |
| Entity Richness | 15 | ≥ 10 |
| Citation Quality | 15 | ≥ 10 |
| Schema Readiness | 15 | ≥ 9 |
| Readability | 10 | ≥ 7 |

### Cara Interpretasi

- **≥ 80:** Excellent — siap untuk AI consumption
- **60-79:** Good — perlu minor improvement
- **40-59:** Fair — perlu restrukturisasi
- **< 40:** Poor — rewrite diperlukan

## Checklist Konten GEO

- [ ] Heading hierarchy jelas
- [ ] Ada Q&A section minimal 3 pertanyaan
- [ ] Entity density ≥ 10%
- [ ] Minimal 3 citation/sumber
- [ ] Ada bullet/numbered list
- [ ] Rata-rata kalimat ≤ 20 kata

**Next:** Module D — Schema Markup & Structured Data untuk GEO.
