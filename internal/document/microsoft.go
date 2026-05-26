// Package microsoft provides text transform utilities
package document

import (
	"archive/zip"
	"encoding/xml"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

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

	var slideFiles []*zip.File
	for _, f := range r.File {
		if strings.HasPrefix(f.Name, "ppt/slides/slide") && strings.HasSuffix(f.Name, ".xml") {
			slideFiles = append(slideFiles, f)
		}
	}

	sort.Slice(slideFiles, func(i, j int) bool {
		numI, _ := strconv.Atoi(strings.TrimSuffix(strings.TrimPrefix(slideFiles[i].Name, "ppt/slides/slide"), ".xml"))
		numJ, _ := strconv.Atoi(strings.TrimSuffix(strings.TrimPrefix(slideFiles[j].Name, "ppt/slides/slide"), ".xml"))
		return numI < numJ
	})

	var lines []string

	for _, f := range slideFiles {

		rc, err := f.Open()
		if err != nil {
			return nil, err
		}
		defer rc.Close()

		var slide struct {
			CSld struct {
				SpTree struct {
					Shapes []struct {
						TxBody struct {
							Paragraphs []struct {
								Runs []struct {
									Text string `xml:"t"`
								} `xml:"r"`
							} `xml:"p"`
						} `xml:"txBody"`
					} `xml:"sp"`
				} `xml:"spTree"`
			} `xml:"cSld"`
		}

		if err := xml.NewDecoder(rc).Decode(&slide); err != nil {
			return nil, err
		}

		for _, sp := range slide.CSld.SpTree.Shapes {
			for _, p := range sp.TxBody.Paragraphs {
				var line string
				for _, run := range p.Runs {
					line += run.Text
				}
				lines = append(lines, line)
			}
		}
	}

	return lines, nil
}
