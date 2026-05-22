// Package file provides text transform utitilies
package file

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
)

const (
	FormatTxt = "txt"
)

func Read(filePath string) ([]string, error) {
	switch filepath.Ext(filePath) {
	case ".txt":
		return ReadTextFile(filePath)
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
