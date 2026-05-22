package main

import (
	"github.com/tq303/ddx/internal/file"
)

func main() {
	lines, err := file.Read("./test.txt")
	if err != nil {
		return
	}

	file.Filter(lines)
}
