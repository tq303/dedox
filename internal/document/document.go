// Package document provides text transform utilities
package document

import (
	"archive/zip"
	"bufio"
	"bytes"
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ledongthuc/pdf"
)

const (
	FileTypeText = ".txt"
	FileTypePdf  = ".pdf"
	FileTypeDocx = ".docx"
)

func Read(filePath string) ([]string, error) {
	switch filepath.Ext(filePath) {
	case FileTypeText:
		return ReadTextFile(filePath)
	case FileTypePdf:
		return ReadPdfFile(filePath)
	case FileTypeDocx:
		return ReadDocxFile(filePath)
	default:
		return nil, fmt.Errorf("unsupported file type: %s", filepath.Ext(filePath))
	}
}

func ReadTextFile(filePath string) ([]string, error) {
	fileContents, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	defer func() {
		if cerr := fileContents.Close(); cerr != nil && err == nil {
			err = cerr
		}
	}()

	scanner := bufio.NewScanner(fileContents)
	lines := []string{}

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}

func ReadPdfFile(filePath string) ([]string, error) {
	f, r, err := pdf.Open(filePath)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	var buf bytes.Buffer
	b, err := r.GetPlainText()
	if err != nil {
		return nil, err
	}

	buf.ReadFrom(b)

	return strings.Split(buf.String(), "\n"), nil
}

func ReadDocxFile(filePath string) ([]string, error) {
	r, err := zip.OpenReader(filePath)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	for _, f := range r.File {
		if f.Name == "word/document.xml" {
			rc, err := f.Open()
			if err != nil {
				return nil, err
			}
			defer rc.Close()

			var doc struct {
				Body struct {
					Paragraphs []struct {
						Props struct {
							Style struct {
								Val string `xml:"val,attr"`
							} `xml:"pStyle"`
						} `xml:"pPr"`
						Runs []struct {
							Text string `xml:"t"`
						} `xml:"r"`
					} `xml:"p"`
				} `xml:"body"`
			}

			if err := xml.NewDecoder(rc).Decode(&doc); err != nil {
				return nil, err
			}

			var lines []string
			for _, p := range doc.Body.Paragraphs {
				var line string
				for _, r := range p.Runs {
					line += r.Text
				}
				style := p.Props.Style.Val
				if strings.HasPrefix(style, "Heading") {
					lines = append(lines, "")
					lines = append(lines, line)
					lines = append(lines, "")
				} else {
					lines = append(lines, line)
				}
			}
			return lines, nil
		}
	}

	return nil, fmt.Errorf("no document.xml found in docx")
}

func Filter(lines []string) {
	seen := map[string]bool{}

	for _, line := range lines {
		if !seen[line] {
			seen[line] = true
		}
	}

	for line := range seen {
		fmt.Printf("%s", line)
	}
}
