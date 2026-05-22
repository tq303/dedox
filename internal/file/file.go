// Package file provides text transform utitilies
package file

import (
	"bufio"
	"fmt"
	"os"
)

const (
	FormatTxt = "txt"
)

func Read() []string {
	fileContents, err := os.Open("./test.txt")
	if err != nil {
		return []string{}
	}

	defer fileContents.Close()

	scanner := bufio.NewScanner(fileContents)
	lines := []string{}

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return []string{}
	}

	return lines
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
