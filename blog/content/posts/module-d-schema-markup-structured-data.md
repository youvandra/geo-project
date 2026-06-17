---
title: "Module D: Schema Markup & Structured Data untuk GEO"
date: 2026-06-18
draft: false
tags: ["geo", "schema", "json-ld", "structured-data", "modul-d"]
categories: ["module-d"]
---

Structured data adalah bahasa yang dipahami mesin. Untuk GEO, schema markup membantu AI memahami konteks dan hubungan konten kita.

<!--more-->

## Kenapa Schema Penting untuk GEO?

AI search engines menggunakan structured data untuk:

1. **Entity recognition** — mengenali orang, tempat, organisasi
2. **Relationship mapping** — memahami hubungan antar entity
3. **Answer extraction** — mengambil jawaban dari FAQ/HowTo
4. **Authority signals** — validasi kredibilitas konten

## Schema Types yang Paling Relevan

### Article
Untuk blog post dan artikel:

```json
{
  "@context": "https://schema.org",
  "@type": "Article",
  "headline": "Belajar GEO dari A-Z",
  "author": {
    "@type": "Person",
    "name": "kucnigplaygame"
  },
  "datePublished": "2026-06-18",
  "description": "Panduan lengkap Generative Engine Optimization"
}
```

### FAQPage
Untuk Q&A — sangat efektif untuk AI search:

```json
{
  "@context": "https://schema.org",
  "@type": "FAQPage",
  "mainEntity": [{
    "@type": "Question",
    "name": "Apa itu GEO?",
    "acceptedAnswer": {
      "@type": "Answer",
      "text": "GEO adalah Generative Engine Optimization..."
    }
  }]
}
```

### HowTo
Untuk tutorial dan panduan langkah-demi-langkah.

### Person & Organization
Untuk author authority dan brand signaling.

## Praktik dengan `geo schema`

Tool kita bisa generate schema markup dengan mudah:

```bash
# Article schema
geo schema Article headline="Belajar GEO" author="kucnigplaygame" --html

# Person schema  
geo schema Person name="kucnigplaygame" jobTitle="GEO Practitioner" --html

# Organization schema
geo schema Organization name="GEO Project" url="https://github.com/youvandra/geo-project" --html
```

Flag `--html` menghasilkan output siap pakai:

```html
<script type="application/ld+json">
{
  "@context": "https://schema.org",
  "@type": "Article",
  "headline": "Belajar GEO",
  "author": "kucnigplaygame"
}
</script>
```

## Best Practices

1. **Satu halaman bisa punya multiple schemas** — Article + Person + Organization
2. **Gunakan `@id`** untuk entity linking antar schema
3. **Update dateModified** setiap konten diupdate (AI suka fresh content)
4. **Validasi** dengan Google Rich Results Test sebelum publish

## Tools

| Tool | Kegunaan |
|------|----------|
| `geo schema <type> key=val` | Generate JSON-LD |
| Google Rich Results Test | Validasi schema |
| Schema.org | Referensi tipe schema |

**Next:** Module E — Citation, Authority & Factual Accuracy.
