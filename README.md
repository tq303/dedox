# ddx

![version](https://img.shields.io/github/v/release/tq303/dedox) ![build](https://github.com/tq303/dedox/actions/workflows/release.yml/badge.svg) ![language](https://img.shields.io/badge/built%20with-Go-00ADD8) ![license](https://img.shields.io/badge/license-none-lightgrey)

CLI tool that extracts text from documents and outputs it to stdout.

---

## Install

**macOS / Linux**
```bash
curl -sL https://raw.githubusercontent.com/tq303/dedox/main/scripts/install.sh | sh
```

**Windows** (run as administrator)
```bat
curl -sL https://raw.githubusercontent.com/tq303/dedox/main/scripts/install.bat -o install.bat && install.bat
```

**npm / yarn**
```bash
npm install @tq303/dedox
yarn add @tq303/dedox
```

**Go**
```bash
go install github.com/tq303/dedox@latest
```

**Local development**
```bash
make install
```

**Link to a Node.js project (no publish needed)**

Build the binary into the npm package, register it globally, then link from your project:
```bash
# in ddx/
make dev                        # builds into npm/bin/ddx
cd npm && yarn link             # registers @tq303/dedox globally

# in your project/
yarn link "@tq303/dedox"        # symlinks node_modules to local ddx/npm/
```

After any Go change, just `make dev` — the linked project picks it up immediately.

---

## Usage

```bash
ddx [file]
```

Supported file types: `.pdf`, `.docx`, `.xlsx`, `.pptx`, `.html`, `.rtf`, `.jpg`, `.txt`

### Examples

```bash
ddx report.pdf
ddx notes.docx | grep "keyword"
ddx data.xlsx > output.txt
```

### Filters

Apply one or more named filters with `--filter` (repeatable, applied in order):

| Filter | What it does |
|---|---|
| `pii` | Redacts emails, phone numbers, SSNs, credit card numbers |
| `urls` | Redacts URLs |
| `ip` | Redacts IPv4 addresses |
| `boilerplate` | Strips page numbers, copyright lines, "Confidential" stamps |
| `norm` | Collapses multiple spaces and consecutive blank lines |
| `uniq` | Removes duplicate lines |

```bash
ddx report.pdf --filter boilerplate --filter norm
ddx contract.docx --filter pii --filter urls > redacted.txt
ddx logs.txt --filter ip | grep "error"
```

### Node.js

```js
const ddx = require('@tq303/dedox')

const text = ddx('report.pdf')
const redacted = ddx('contract.docx', { filters: ['pii', 'norm'] })
```

---

## Principles

- Non-destructive — never modifies the source file
- No runtime dependencies — single compiled binary
- Composable — output to stdout, works in pipelines like any Unix tool
- Fast — no LLM, no network, all local processing
- Simple — no config needed to get started
