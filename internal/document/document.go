// Package document provides text transform utilities
package document

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const (
	FileTypePdf        = ".pdf"
	FileTypeDocx       = ".docx"
	FileTypeExcel      = ".xlsx"
	FileTypePowerPoint = ".pptx"
)

func ReadHttpFile(url string) ([]string, error) {
	resp, err := http.Get(url)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	ext := filepath.Ext(url)

	f, err := os.CreateTemp("/tmp", "ddx-http.*"+ext)

	if err != nil {
		return nil, err
	}

	defer f.Close()

	if _, err := io.Copy(f, resp.Body); err != nil {
		return nil, err
	}

	return ReadFile(f.Name())
}

func Read(filePath string) ([]string, error) {
	if strings.HasPrefix(filePath, "http") {
		return ReadHttpFile(filePath)
	}

	return ReadFile(filePath)
}

func ReadFile(filePath string) ([]string, error) {
	switch filepath.Ext(filePath) {
	case FileTypePdf:
		return ReadPdfFile(filePath)
	case FileTypeDocx:
		return ReadDocxFile(filePath)
	case FileTypeExcel:
		return ReadXlsxFile(filePath)
	case FileTypePowerPoint:
		return ReadPowerPointFile(filePath)
	default:
		return ReadTextFile(filePath)
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

func Filter(lines []string) {
	seen := map[string]bool{}

	for _, line := range lines {
		if !seen[line] {
			seen[line] = true
			fmt.Println(line)
		}
	}
}
