---
title: "Module E: Citation, Authority & Factual Accuracy"
date: 2026-06-18
draft: false
tags: ["geo", "citation", "authority", "entity-linking", "modul-e"]
categories: ["module-e"]
---

AI search engines sangat selektif dalam memilih sumber. Faktor terbesarnya? **Authority** dan **factual accuracy**.

<!--more-->

## Kenapa Citation Penting?

Model RAG (Retrieval-Augmented Generation) bekerja dengan:

1. **Retrieve** — mengambil dokumen relevan dari knowledge base
2. **Rank** — memprioritaskan sumber dengan authority tinggi
3. **Generate** — menyusun jawaban berdasarkan sumber terpilih

Tanpa citation yang kuat, konten kita tidak akan terpilih di tahap **Rank**.

## Entity Linking dengan `geo entity`

Tool `geo entity` menghubungkan teks kita ke Wikipedia/Knowledge Graph:

```bash
echo "Google and Perplexity AI are search engines." | geo entity --stdin
```

Output:
```
Found 4 linked entities
  Google (high) — Google LLC is an American multinational...
  Perplexity AI (high) — Perplexity AI is an American...
```

### Kenapa Ini Penting?

- Entity yang terverifikasi (via Wikipedia) meningkatkan **trust signal**
- AI engines lebih percaya konten yang entity-nya match dengan Knowledge Graph
- Setiap linked entity adalah **anchor point** untuk AI memahami konteks

## Citation Quality Signals

### What Makes a Good Citation?

| Kriteria | Contoh Baik | Contoh Buruk |
|----------|-------------|--------------|
| Source | jurnal, .edu, .gov | blog random |
| Recency | ≤ 2 tahun | ≥ 5 tahun |
| Relevance | langsung terkait | tangent |
| Specificity | data spesifik | opini umum |

### Cara Menulis Citation untuk GEO

❌ "Menurut penelitian..."
✅ "Menurut penelitian Princeton University (2024), GEO dapat meningkatkan visibility hingga 40% di AI-generated answers..."

## Authority Building untuk GEO

### On-Page Signals
- Author bio dengan schema Person
- Organization schema dengan logo
- Publication date yang jelas
- Update history

### Off-Page Signals
- Wikipedia mentions (sangat kuat)
- Links dari .edu/.gov
- Citations di publikasi industri
- Social proof dari platform terpercaya

## Factual Accuracy Checker

Sebelum publikasi, verifikasi:

1. **Entity names** — bener ejaannya? (Google, not Googl)
2. **Dates & numbers** — valid? (2024, not 2004)
3. **Claims** — ada sumbernya? (jangan asbun)
4. **Relationships** — entity A benar terkait dengan entity B?

## Tools

| Tool | Fungsi |
|------|--------|
| `geo entity <text>` | Link entity ke Knowledge Graph |
| `geo score <text>` | Citation quality check |
| Wikipedia API | Verifikasi entity |
| Wikidata | Entity relationship lookup |

**Next:** Module F — Technical GEO.
