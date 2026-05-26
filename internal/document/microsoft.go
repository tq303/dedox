// Package microsoft provides text transform utilities
package document

import (
	"archive/zip"
	"encoding/xml"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/tq303/ddx/internal/parse"
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
							Props struct {
								Bold *struct{} `xml:"b"`
							} `xml:"rPr"`
							Text string `xml:"t"`
						} `xml:"r"`
					} `xml:"p"`
				} `xml:"body"`
			}

			if err := xml.NewDecoder(rc).Decode(&doc); err != nil {
				return nil, err
			}

			var out strings.Builder
			out.WriteString("<html><body>")
			for _, p := range doc.Body.Paragraphs {
				var text strings.Builder
				for _, run := range p.Runs {
					t := run.Text
					if run.Props.Bold != nil {
						t = "<strong>" + t + "</strong>"
					}
					text.WriteString(t)
				}
				style := p.Props.Style.Val
				level := 0
				if strings.HasPrefix(style, "Heading") {
					fmt.Sscanf(style[len("Heading"):], "%d", &level)
				}
				if level >= 1 && level <= 6 {
					fmt.Fprintf(&out, "<h%d>%s</h%d>", level, text.String(), level)
				} else {
					out.WriteString("<p>" + text.String() + "</p>")
				}
			}
			out.WriteString("</body></html>")
			return parse.HtmlToMarkdown(strings.NewReader(out.String()))
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

			// resolve cell values and find max width
			rows := make([][]string, len(sheet.Rows))
			maxCols := 0
			for i, row := range sheet.Rows {
				cols := make([]string, len(row.Cells))
				for j, cell := range row.Cells {
					val := cell.Value
					if cell.Type == "s" {
						idx := 0
						fmt.Sscanf(val, "%d", &idx)
						if idx < len(sharedStrings) {
							val = sharedStrings[idx]
						}
					}
					cols[j] = val
				}
				rows[i] = cols
				if len(cols) > maxCols {
					maxCols = len(cols)
				}
			}

			var out strings.Builder
			out.WriteString("<html><body><table>")
			for i, cols := range rows {
				tag := "td"
				if i == 0 {
					tag = "th"
				}
				out.WriteString("<tr>")
				for j := 0; j < maxCols; j++ {
					val := ""
					if j < len(cols) {
						val = cols[j]
					}
					fmt.Fprintf(&out, "<%s>%s</%s>", tag, val, tag)
				}
				out.WriteString("</tr>")
			}
			out.WriteString("</table></body></html>")
			return parse.HtmlToMarkdown(strings.NewReader(out.String()))
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

	var out strings.Builder
	out.WriteString("<html><body>")

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
				var line strings.Builder
				for _, run := range p.Runs {
					line.WriteString(run.Text)
				}
				if line.Len() > 0 {
					out.WriteString("<p>" + line.String() + "</p>")
				}
			}
		}
	}

	out.WriteString("</body></html>")
	return parse.HtmlToMarkdown(strings.NewReader(out.String()))
}
