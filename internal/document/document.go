// Package document provides text transform utilities
package document

import (
	"archive/zip"
	"bufio"
	"encoding/xml"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"

	"github.com/ledongthuc/pdf"
)

const (
	FileTypeText  = ".txt"
	FileTypePdf   = ".pdf"
	FileTypeDocx  = ".docx"
	FileTypeExcel = ".xlsx"
)

func Read(filePath string) ([]string, error) {
	switch filepath.Ext(filePath) {
	case FileTypeText:
		return ReadTextFile(filePath)
	case FileTypePdf:
		return ReadPdfFile(filePath)
	case FileTypeDocx:
		return ReadDocxFile(filePath)
	case FileTypeExcel:
		return ReadXlsxFile(filePath)
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

	var lines []string
	for i := 1; i <= r.NumPage(); i++ {
		page := r.Page(i)
		content := page.Content()

		var currentY float64
		var currentLine string

		for _, text := range content.Text {
			if currentLine == "" {
				currentY = text.Y
			}
			if math.Abs(text.Y-currentY) > 1 {
				lines = append(lines, currentLine)
				currentLine = text.S
				currentY = text.Y
			} else {
				currentLine += text.S
			}
		}
		if currentLine != "" {
			lines = append(lines, currentLine)
		}
	}

	return lines, nil
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

func ReadXlsxFile(filePath string) ([]string, error) {
	r, err := zip.OpenReader(filePath)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	// parse shared strings first
	var sharedStrings []string
	for _, f := range r.File {
		if f.Name == "xl/sharedStrings.xml" {
			rc, err := f.Open()
			if err != nil {
				return nil, err
			}
			defer rc.Close()

			var ss struct {
				Items []struct {
					Text string `xml:"t"`
				} `xml:"si"`
			}
			if err := xml.NewDecoder(rc).Decode(&ss); err != nil {
				return nil, err
			}
			for _, item := range ss.Items {
				sharedStrings = append(sharedStrings, item.Text)
			}
			break
		}
	}

	for _, f := range r.File {
		if f.Name == "xl/worksheets/sheet1.xml" {
			rc, err := f.Open()
			if err != nil {
				return nil, err
			}
			defer rc.Close()

			var sheet struct {
				Rows []struct {
					Cells []struct {
						Type  string `xml:"t,attr"`
						Value string `xml:"v"`
					} `xml:"c"`
				} `xml:"sheetData>row"`
			}

			if err := xml.NewDecoder(rc).Decode(&sheet); err != nil {
				return nil, err
			}

			var lines []string
			for _, row := range sheet.Rows {
				var cols []string
				for _, cell := range row.Cells {
					val := cell.Value
					if cell.Type == "s" {
						idx := 0
						fmt.Sscanf(val, "%d", &idx)
						if idx < len(sharedStrings) {
							val = sharedStrings[idx]
						}
					}
					cols = append(cols, val)
				}
				lines = append(lines, strings.Join(cols, "\t"))
			}
			return lines, nil
		}
	}

	return nil, fmt.Errorf("no sheet1.xml found in %s", filePath)
}

func ReadPowerPointFile(filePath string) ([]string, error) {
	r, err := zip.OpenReader(filePath)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	var test = make([]string, 0)

	return test, nil
}

func Filter(lines []string) {
	seen := map[string]bool{}

	for _, line := range lines {
		if !seen[line] {
			seen[line] = true
			fmt.Println(line)
		}
	}
}
