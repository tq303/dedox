# DDox (`ddx`) — Plan

A single binary that reads any file and outputs clean, normalized text to stdout.

---

## Phase 1 — Plain Text ✓

**Goal:** Read a `.txt` file and output normalized text to stdout.

- [x] Read `.txt` file and output to stdout
- [ ] Strip extra whitespace and blank lines
- [ ] Normalize line endings (CRLF → LF)
- [ ] Fix common encoding issues (UTF-8 normalization)
- [ ] Basic flag: `--prune` to strip repeated blank lines and trailing whitespace
- Output to stdout only

**Usage:**
```
ddx file.txt
ddx file.txt --prune
```

---

## Phase 2 — PDF ✓

**Goal:** Read a `.pdf` file and output clean normalized text to stdout.

- [x] Extract text from PDF using `ledongthuc/pdf`
- [ ] Strip page numbers, headers, footers via heuristics
- [ ] Normalize whitespace and line breaks post-extraction
- [ ] Handle common PDF encoding issues
- [ ] `--prune` flag applies here too

**Usage:**
```
ddx file.pdf
ddx file.pdf --prune
```

---

## Phase 3 — DOCX (later)

**Goal:** Read a `.docx` file and output clean normalized text to stdout.

- Extract text from DOCX XML structure
- Strip style metadata, keep content
- Normalize output consistent with Phase 1 and 2

---

## Phase 4 — Output Styles

**Goal:** Control output format via a `--style` / `-s` flag.

- `txt` — plain text (default)
- `md` — markdown, preserving headings, tables, lists, bold/italic where possible

**Usage:**
```
ddx file.docx -s md
ddx file.docx -s txt
```

---

## Principles

- Non-destructive — never modifies the source file
- No runtime dependencies — single compiled binary
- Composable — works in pipelines like any Unix tool
- Fast — no LLM, no network, all local processing
- Simple — no config needed to get started
