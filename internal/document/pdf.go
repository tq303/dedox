// Package document provides text transform utilities
package document

import (
	"math"
	"sort"
	"strings"

	pdflib "github.com/ledongthuc/pdf"
)

type pdfPara struct {
	text string
	size float64
}

func ReadPdfFile(filePath string) ([]string, error) {
	f, r, err := pdflib.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var allParas []pdfPara

	for i := 1; i <= r.NumPage(); i++ {
		texts := r.Page(i).Content().Text

		// Reading order: top-to-bottom (Y descending), left-to-right (X ascending).
		sort.SliceStable(texts, func(a, b int) bool {
			if math.Abs(texts[a].Y-texts[b].Y) > 1 {
				return texts[a].Y > texts[b].Y
			}
			return texts[a].X < texts[b].X
		})

		// pass 1: collect visual lines with Y position and max font size
		type rawLine struct {
			text string
			size float64
			y    float64
		}
		var pageLines []rawLine
		var currentY, nextX, maxSize float64
		var currentLine string

		flushLine := func() {
			if strings.TrimSpace(currentLine) != "" {
				pageLines = append(pageLines, rawLine{text: currentLine, size: maxSize, y: currentY})
			}
			currentLine = ""
			maxSize = 0
		}

		for _, t := range texts {
			if currentLine == "" {
				currentY = t.Y
				currentLine = t.S
				nextX = t.X + t.W
				maxSize = t.FontSize
				continue
			}
			if math.Abs(t.Y-currentY) > 1 {
				flushLine()
				currentLine = t.S
				currentY = t.Y
				nextX = t.X + t.W
				maxSize = t.FontSize
			} else {
				if t.X-nextX > t.FontSize*0.2 && !strings.HasSuffix(currentLine, " ") {
					currentLine += " "
				}
				currentLine += t.S
				nextX = t.X + t.W
				if t.FontSize > maxSize {
					maxSize = t.FontSize
				}
			}
		}
		flushLine()

		// pass 2: group lines into paragraphs.
		// break on font size change OR large Y gap (> 1.5x line height)
		var paraLines []string
		var paraSz float64
		flushPara := func() {
			if len(paraLines) > 0 {
				allParas = append(allParas, pdfPara{
					text: strings.Join(paraLines, " "),
					size: paraSz,
				})
				paraLines = nil
				paraSz = 0
			}
		}
		for j, line := range pageLines {
			if j > 0 {
				prev := pageLines[j-1]
				yGap := prev.y - line.y
				sizeChanged := math.Round(line.size) != math.Round(prev.size)
				if sizeChanged || yGap > prev.size*1.5 {
					flushPara()
				}
			}
			paraLines = append(paraLines, line.text)
			if line.size > paraSz {
				paraSz = line.size
			}
		}
		flushPara()
	}

	if len(allParas) == 0 {
		return ReadJpegBase64(filePath)
	}

	// find base font size: size with the most total characters (body text dominates)
	szChars := map[float64]int{}
	for _, p := range allParas {
		szChars[math.Round(p.size)] += len(p.text)
	}
	baseSz := 12.0
	maxChars := 0
	for sz, n := range szChars {
		if n > maxChars || (n == maxChars && sz < baseSz) {
			maxChars = n
			baseSz = sz
		}
	}

	// rank sizes above base as heading levels h1–h6
	var headingSizes []float64
	for sz := range szChars {
		if sz > baseSz {
			headingSizes = append(headingSizes, sz)
		}
	}
	sort.Slice(headingSizes, func(a, b int) bool { return headingSizes[a] > headingSizes[b] })
	szToLevel := map[float64]int{}
	for i, sz := range headingSizes {
		if i < 6 {
			szToLevel[sz] = i + 1
		}
	}

	// emit markdown directly — avoids double blank lines from the HTML pipeline
	var out []string
	for _, p := range allParas {
		if len(out) > 0 {
			out = append(out, "")
		}
		text := strings.TrimSpace(p.text)
		if level, ok := szToLevel[math.Round(p.size)]; ok {
			out = append(out, strings.Repeat("#", level)+" "+text)
		} else {
			out = append(out, text)
		}
	}
	return out, nil
}
