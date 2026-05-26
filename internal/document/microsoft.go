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
								Sz   struct {
									Val int `xml:"val,attr"`
								} `xml:"sz"`
							} `xml:"rPr"`
							Text string `xml:"t"`
						} `xml:"r"`
					} `xml:"p"`
				} `xml:"body"`
			}

			if err := xml.NewDecoder(rc).Decode(&doc); err != nil {
				return nil, err
			}

			// find base font size: most common sz across all runs
			szCount := map[int]int{}
			for _, p := range doc.Body.Paragraphs {
				for _, run := range p.Runs {
					if run.Props.Sz.Val > 0 {
						szCount[run.Props.Sz.Val]++
					}
				}
			}
			baseSz := 28 // fallback: 14pt
			maxCount := 0
			for sz, count := range szCount {
				if count > maxCount {
					maxCount = count
					baseSz = sz
				}
			}

			// rank distinct sizes above base: largest = h1, next = h2, etc.
			var headingSizes []int
			for sz := range szCount {
				if sz > baseSz {
					headingSizes = append(headingSizes, sz)
				}
			}
			sort.Sort(sort.Reverse(sort.IntSlice(headingSizes)))
			szToLevel := map[int]int{}
			for i, sz := range headingSizes {
				if i < 6 {
					szToLevel[sz] = i + 1
				}
			}

			var out strings.Builder
			out.WriteString("<html><body>")
			for _, p := range doc.Body.Paragraphs {
				var text strings.Builder
				maxSz := 0
				for _, run := range p.Runs {
					t := run.Text
					if run.Props.Bold != nil {
						t = "<strong>" + t + "</strong>"
					}
					text.WriteString(t)
					if run.Props.Sz.Val > maxSz {
						maxSz = run.Props.Sz.Val
					}
				}
				style := p.Props.Style.Val
				level := 0
				if strings.HasPrefix(style, "Heading") {
					fmt.Sscanf(style[len("Heading"):], "%d", &level)
				} else if l, ok := szToLevel[maxSz]; ok {
					level = l
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
						NvSpPr struct {
							NvPr struct {
								Ph struct {
									Type string `xml:"type,attr"`
								} `xml:"ph"`
							} `xml:"nvPr"`
						} `xml:"nvSpPr"`
						TxBody struct {
							Paragraphs []struct {
								Props struct {
									Level int `xml:"lvl,attr"`
								} `xml:"pPr"`
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
			phType := sp.NvSpPr.NvPr.Ph.Type
			isTitle := phType == "title" || phType == "ctrTitle"
			isSubtitle := phType == "subTitle"

			if isTitle || isSubtitle {
				tag := "h2"
				if isSubtitle {
					tag = "h3"
				}
				for _, p := range sp.TxBody.Paragraphs {
					var line strings.Builder
					for _, run := range p.Runs {
						line.WriteString(run.Text)
					}
					if line.Len() > 0 {
						fmt.Fprintf(&out, "<%s>%s</%s>", tag, line.String(), tag)
					}
				}
			} else {
				// body content: nested list based on paragraph level
				depth := -1
				for _, p := range sp.TxBody.Paragraphs {
					var line strings.Builder
					for _, run := range p.Runs {
						line.WriteString(run.Text)
					}
					if line.Len() == 0 {
						continue
					}
					lvl := p.Props.Level
					for depth < lvl {
						out.WriteString("<ul>")
						depth++
					}
					for depth > lvl {
						out.WriteString("</ul>")
						depth--
					}
					out.WriteString("<li>" + line.String() + "</li>")
				}
				for depth >= 0 {
					out.WriteString("</ul>")
					depth--
				}
			}
		}
	}

	out.WriteString("</body></html>")
	return parse.HtmlToMarkdown(strings.NewReader(out.String()))
}
