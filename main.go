package main

import (
	"github.com/tq303/ddx/internal/file"
)

func main() {
	text := file.Read()
	file.Filter(text)
}
