// Package document provides text transform utilities
package document

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/tq303/dedox/internal/parse"
)

func contentTypeExt(url string, contentType string) string {
	if strings.HasPrefix(contentType, "text/html") {
		return ".html"
	}

	if strings.HasPrefix(contentType, "application/octet-stream") {
		return ".pdf"
	}

	return filepath.Ext(url)
}

func ReadHttpFile(url string) (string, error) {
	resp, err := http.Get(url)

	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	ext := contentTypeExt(url, resp.Header.Get("Content-Type"))

	f, err := os.CreateTemp("/tmp", "ddx-http.*"+ext)

	if err != nil {
		return "", err
	}

	defer f.Close()

	if _, err := io.Copy(f, resp.Body); err != nil {
		return "", err
	}

	return f.Name(), nil
}

func ReadHtmlFile(filePath string) ([]string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	defer func() {
		if cerr := f.Close(); cerr != nil && err == nil {
			err = cerr
		}
	}()

	return parse.HtmlToMarkdown(f)
}
