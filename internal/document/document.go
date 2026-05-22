// Package document provides text transform utilities
package document

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ledongthuc/pdf"
)

const (
	FileTypeText = ".txt"
	FileTypePdf  = ".pdf"
	// FileTypeWord = ".docx"
)

func Read(filePath string) ([]string, error) {
	switch filepath.Ext(filePath) {
	case FileTypeText:
		return ReadTextFile(filePath)
	case FileTypePdf:
		return ReadPdfFile(filePath)
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

	var buf bytes.Buffer
	b, err := r.GetPlainText()
	if err != nil {
		return nil, err
	}

	buf.ReadFrom(b)

	return strings.Split(buf.String(), "\n"), nil
}

func Filter(lines []string) {
	seen := map[string]bool{}

	for _, line := range lines {
		if !seen[line] {
			seen[line] = true
		}
	}

	for line := range seen {
		fmt.Printf("%s", line)
	}
}
