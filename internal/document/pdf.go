// Package document provides text transform utilities
package document

import (
	"math"

	"github.com/ledongthuc/pdf"
)

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

	if len(lines) == 0 {
		return ReadJpegBase64(filePath)
	}

	return lines, nil
}
