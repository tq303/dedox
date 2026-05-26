// Package document provides text transform utilities
package document

import (
	"bytes"
	"encoding/base64"
	"os"
)

var jpgStart = []byte{0xFF, 0xD8, 0xFF}
var jpgEnd = []byte{0xFF, 0xD9}

func ReadJpegBase64(filePath string) ([]string, error) {
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
