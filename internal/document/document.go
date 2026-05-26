// Package document provides text transform utilities
package document

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
)

const (
	FileTypePdf        = ".pdf"
	FileTypeDocx       = ".docx"
	FileTypeExcel      = ".xlsx"
	FileTypePowerPoint = ".pptx"
)

func Read(filePath string) ([]string, error) {
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
		fmt.Printf("unsupported file type: %s", filepath.Ext(filePath))
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
