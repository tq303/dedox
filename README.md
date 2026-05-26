# ddx

![version](https://img.shields.io/github/v/release/tq303/ddx) ![build](https://github.com/tq303/ddx/actions/workflows/release.yml/badge.svg) ![language](https://img.shields.io/badge/built%20with-Go-00ADD8) ![license](https://img.shields.io/badge/license-none-lightgrey)

CLI tool that extracts text from documents and outputs it to stdout.

---

## Install

**macOS (Apple Silicon)**
```bash
curl -L https://github.com/tq303/ddx/releases/latest/download/ddx-darwin-arm64 -o /usr/local/bin/ddx && chmod +x /usr/local/bin/ddx
```

**macOS (Intel)**
```bash
curl -L https://github.com/tq303/ddx/releases/latest/download/ddx-darwin-amd64 -o /usr/local/bin/ddx && chmod +x /usr/local/bin/ddx
```

**Linux**
```bash
curl -L https://github.com/tq303/ddx/releases/latest/download/ddx-linux-amd64 -o /usr/local/bin/ddx && chmod +x /usr/local/bin/ddx
```

**npm / yarn**
```bash
npm install ddx
yarn add ddx
```

**Go**
```bash
go install github.com/tq303/ddx@latest
```

**Local development**
```bash
make install
```

---

## Usage

```bash
ddx [file]
```

Supported file types: `.txt`, `.pdf`, `.docx`, `.xlsx`, `.pptx`

### Examples

```bash
ddx report.pdf
ddx notes.docx | grep "keyword"
ddx data.xlsx > output.txt
```

---

## Principles

- Non-destructive — never modifies the source file
- No runtime dependencies — single compiled binary
- Composable — output to stdout, works in pipelines like any Unix tool
- Fast — no LLM, no network, all local processing
- Simple — no config needed to get started
