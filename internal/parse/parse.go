package parse

import (
	"io"
	"slices"
	"strings"

	"golang.org/x/net/html"
)

func sliceHidden(attr []html.Attribute) bool {
	return slices.ContainsFunc(attr, func(a html.Attribute) bool {
		return a.Key == "hidden"
	})
}

func parseMeta(tok html.Token) (name, property, content string) {
	for _, attr := range tok.Attr {
		switch attr.Key {
		case "name":
			name = strings.ToLower(attr.Val)
		case "property":
			property = strings.ToLower(attr.Val)
		case "content":
			content = attr.Val
		}
	}
	return
}

func metaFallback(metaTitle, metaDescription string) []string {
	var lines []string
	if metaTitle != "" {
		lines = append(lines, "# "+metaTitle, "")
	}
	for _, l := range strings.Split(metaDescription, "\n") {
		lines = append(lines, l)
	}
	return lines
}

func HtmlToMarkdown(r io.Reader) ([]string, error) {
	tokenizer := html.NewTokenizer(r)
	skipTags := []string{"script", "style", "head", "img", "footer", "nav"}

	var lines []string
	var metaTitle, metaDescription string
	skipTag := false
	hiddenDepth := 0
	headingLevel := 0
	inListItem := false
	listDepth := 0
	listItemDepth := 0
	inBold := false
	inItalic := false
	linkHref := ""

	// Table state
	inCell := false
	isHeaderRow := false
	headerRowDone := false
	var currentRow []string
	var currentCellBuf strings.Builder

	for {
		tt := tokenizer.Next()
		switch tt {
		case html.ErrorToken:
			if metaDescription != "" {
				nonEmpty := 0
				for _, l := range lines {
					if strings.TrimSpace(l) != "" {
						nonEmpty++
					}
				}
				if nonEmpty < 5 {
					return metaFallback(metaTitle, metaDescription), nil
				}
			}
			return lines, nil

		case html.SelfClosingTagToken:
			token := tokenizer.Token()
			if token.Data == "meta" {
				name, property, content := parseMeta(token)
				if content != "" {
					if (name == "description" || property == "og:description") && metaDescription == "" {
						metaDescription = content
					} else if (name == "title" || property == "og:title") && metaTitle == "" {
						metaTitle = content
					}
				}
			}

		case html.StartTagToken:
			token := tokenizer.Token()
			tag := token.Data

			if tag == "meta" {
				name, property, content := parseMeta(token)
				if content != "" {
					if (name == "description" || property == "og:description") && metaDescription == "" {
						metaDescription = content
					} else if (name == "title" || property == "og:title") && metaTitle == "" {
						metaTitle = content
					}
				}
				continue
			}

			if slices.Contains(skipTags, tag) {
				skipTag = true
				continue
			}
			if hiddenDepth > 0 || sliceHidden(token.Attr) {
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
			case "ul", "ol":
				listDepth++
			case "li":
				inListItem = true
				listItemDepth = listDepth - 1
			case "strong", "b":
				inBold = true
			case "em", "i":
				inItalic = true
			case "a":
				for _, attr := range token.Attr {
					if attr.Key == "href" {
						linkHref = attr.Val
					}
				}
			case "p":
				lines = append(lines, "")
			case "table":
				headerRowDone = false
			case "tr":
				currentRow = currentRow[:0]
				isHeaderRow = false
			case "th", "td":
				inCell = true
				if tag == "th" {
					isHeaderRow = true
				}
				currentCellBuf.Reset()
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
			case "ul", "ol":
				if listDepth > 0 {
					listDepth--
				}
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
			case "th", "td":
				currentRow = append(currentRow, strings.TrimSpace(currentCellBuf.String()))
				inCell = false
			case "tr":
				if len(currentRow) > 0 {
					lines = append(lines, "| "+strings.Join(currentRow, " | ")+" |")
					if isHeaderRow && !headerRowDone {
						sep := make([]string, len(currentRow))
						for i := range sep {
							sep[i] = "---"
						}
						lines = append(lines, "| "+strings.Join(sep, " | ")+" |")
						headerRowDone = true
					}
				}
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

			if inCell {
				if currentCellBuf.Len() > 0 {
					currentCellBuf.WriteString(" ")
				}
				currentCellBuf.WriteString(text)
				continue
			}

			if headingLevel > 0 {
				text = strings.Repeat("#", headingLevel) + " " + text
			} else if inListItem {
				text = strings.Repeat("  ", listItemDepth) + "- " + text
			}

			lines = append(lines, text)
		}
	}
}
