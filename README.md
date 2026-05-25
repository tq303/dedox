# ddx - Dedox

![version](https://img.shields.io/github/v/release/tq303/ddx) ![build](https://github.com/tq303/ddx/actions/workflows/release.yml/badge.svg) ![language](https://img.shields.io/badge/built%20with-Go-00ADD8) ![license](https://img.shields.io/badge/license-none-lightgrey)

CLI file parser

---

## Install

**macOS (Apple Silicon)**
```bash
curl -L https://github.com/tq303/ddx/releases/latest/download/ddx-darwin-arm64 -o /usr/local/bin/ddx && chmod +x /usr/local/bin/ddx
```

**Linux**
```bash
curl -L https://github.com/tq303/ddx/releases/latest/download/ddx-linux-amd64 -o /usr/local/bin/ddx && chmod +x /usr/local/bin/ddx
```

Or with Go:
```bash
go install github.com/tq303/ddx@latest
```

**Local development:**
```bash
make install
```

## Usage

```
Parse and output text from supported file types: .txt, .pdf, .docx

Usage:
  ddx [file] [flags]

Flags:
  -h, --help   help for ddx
```

---

## Principles

- Non-destructive — never modifies the source file
- No runtime dependencies — single compiled binary
- Composable — works in pipelines like any Unix tool
- Fast — no LLM, no network, all local processing
- Simple — no config needed to get started

## Examples

```bash
ddx file.pdf
```

## TODOs

1. New files types
  - HTML
  - PPT
2. Test bigger versions of files

```
