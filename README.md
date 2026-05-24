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

## Examples

```bash
ddx file.txt
ddx file.pdf
ddx file.docx

## TODOs

1. html parser
2. spreadsheet parser
3. Find common file types. e.g. html

```
