// Package document provides text transform utilities
package document

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"golang.org/x/net/html"
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

func sliceHidden(attr []html.Attribute) bool {
	return slices.ContainsFunc(attr, func(a html.Attribute) bool {
		return a.Key == "hidden" || (a.Key == "aria-hidden" && a.Val == "true")
	})
}

func ReadHtmlFile(filePath string) ([]string, error) {
	fileContents, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	defer func() {
		if cerr := fileContents.Close(); cerr != nil && err == nil {
			err = cerr
		}
	}()

	tokenizer := html.NewTokenizer(fileContents)
	skipTags := []string{"script", "style", "head", "img", "footer", "nav"}

	var lines []string
	skipTag := false
	hiddenDepth := 0
	headingLevel := 0
	inListItem := false
	inBold := false
	inItalic := false
	linkHref := ""

	for {
		tt := tokenizer.Next()
		switch tt {
		case html.ErrorToken:
			return lines, nil

		case html.StartTagToken:
			tok := tokenizer.Token()
			tag := tok.Data

			if slices.Contains(skipTags, tag) {
				skipTag = true
				continue
			}
			if sliceHidden(tok.Attr) {
				hiddenDepth++
				continue
			}

			switch tag {
			case "h1":
				headingLevel = 1
			case "h2":
				headingLevel = 2
			case "h3":
				headingLevel = 3
			case "h4":
				headingLevel = 4
			case "h5":
				headingLevel = 5
			case "h6":
				headingLevel = 6
			case "li":
				inListItem = true
			case "strong", "b":
				inBold = true
			case "em", "i":
				inItalic = true
			case "a":
				for _, attr := range tok.Attr {
					if attr.Key == "href" {
						linkHref = attr.Val
					}
				}
			case "p":
				lines = append(lines, "")
			}

		case html.EndTagToken:
			tag := tokenizer.Token().Data

			if slices.Contains(skipTags, tag) {
				skipTag = false
				continue
			}
			if hiddenDepth > 0 {
				hiddenDepth--
				continue
			}

			switch tag {
			case "h1", "h2", "h3", "h4", "h5", "h6":
				headingLevel = 0
				lines = append(lines, "")
			case "li":
				inListItem = false
			case "strong", "b":
				inBold = false
			case "em", "i":
				inItalic = false
			case "a":
				linkHref = ""
			case "p":
				lines = append(lines, "")
			}

		case html.TextToken:
			if skipTag || hiddenDepth > 0 {
				continue
			}
			text := strings.TrimSpace(tokenizer.Token().Data)
			if text == "" {
				continue
			}

			if inBold && headingLevel == 0 {
				text = "**" + text + "**"
			}
			if inItalic {
				text = "_" + text + "_"
			}
			if linkHref != "" {
				text = "[" + text + "](" + linkHref + ")"
			}
			if headingLevel > 0 {
				text = strings.Repeat("#", headingLevel) + " " + text
			} else if inListItem {
				text = "- " + text
			}

			lines = append(lines, text)
		}
	}
}
