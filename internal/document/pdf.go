// Package document provides text transform utilities
package document

import (
	"bytes"
	"encoding/base64"
	"math"
	"os"

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
		return ReadPdfImage(filePath)
	}

	return lines, nil
}

var jpgStart = []byte{0xFF, 0xD8, 0xFF}
var jpgEnd = []byte{0xFF, 0xD9}

func ReadPdfImage(filePath string) ([]string, error) {
	pdf, err := os.ReadFile(filePath)

	if err != nil {
		return nil, err
	}

	images := []string{}

	offset := 0

	for {
		start := bytes.Index(pdf[offset:], jpgStart)
		if start == -1 {
			break
		}

		start += offset

		end := bytes.Index(pdf[start:], jpgEnd)
		if end == -1 {
			break
		}

		end = start + end + 2

		jpeg := pdf[start:end]

		encoding := "data:image/jpeg;base64," + base64.StdEncoding.EncodeToString(jpeg)

		images = append(images, encoding)

		offset = end
	}

	return images, nil
}
