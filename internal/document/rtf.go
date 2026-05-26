package document

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func ReadRtfFile(filePath string) ([]string, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	return parseHTML(strings.NewReader(rtfToHTML(data)))
}

func rtfToHTML(data []byte) string {
	var out strings.Builder
	out.WriteString("<html><body>")

	i := 0
	n := len(data)
	groupDepth := 0
	skipDepth := 0
	skipping := false

	headingLevel := 0
	isBullet := false

	type frame struct{ bold bool }
	stack := []frame{{}}

	var para strings.Builder

	// Table state
	inTable := false
	var tableRow []string

	emit := func() {
		text := strings.TrimSpace(para.String())
		para.Reset()
		if text == "" {
			return
		}
		if headingLevel >= 1 && headingLevel <= 6 {
			out.WriteString(fmt.Sprintf("<h%d>%s</h%d>", headingLevel, text, headingLevel))
		} else if isBullet {
			out.WriteString("<li>" + text + "</li>")
		} else {
			out.WriteString("<p>" + text + "</p>")
		}
		headingLevel = 0
		isBullet = false
	}

	for i < n {
		if skipping {
			switch data[i] {
			case '{':
				groupDepth++
			case '}':
				groupDepth--
				if groupDepth < skipDepth {
					skipping = false
				}
			}
			i++
			continue
		}

		switch data[i] {
		case '{':
			groupDepth++
			stack = append(stack, frame{stack[len(stack)-1].bold})
			i++

			// peek for skip groups and bullet detection
			peek := i
			for peek < n && data[peek] == ' ' {
				peek++
			}
			if peek < n && data[peek] == '\\' {
				peek++
				if peek < n && data[peek] == '*' {
					// check for bullet list before skipping
					snippet := string(data[peek:min(peek+60, n)])
					if strings.Contains(snippet, "pnlvlblt") {
						isBullet = true
					}
					skipping = true
					skipDepth = groupDepth
					stack = stack[:len(stack)-1]
					i = peek + 1
					continue
				}
				wordEnd := peek
				for wordEnd < n && data[wordEnd] >= 'a' && data[wordEnd] <= 'z' {
					wordEnd++
				}
				firstWord := string(data[peek:wordEnd])
				if firstWord == "pntext" || firstWord == "fonttbl" || firstWord == "colortbl" || firstWord == "stylesheet" {
					skipping = true
					skipDepth = groupDepth
					stack = stack[:len(stack)-1]
					continue
				}
			}

		case '}':
			if len(stack) > 1 {
				stack = stack[:len(stack)-1]
			}
			groupDepth--
			i++

		case '\\':
			i++
			if i >= n {
				break
			}

			// escaped literal chars
			if data[i] == '\\' || data[i] == '{' || data[i] == '}' {
				para.WriteByte(data[i])
				i++
				continue
			}
			// hex escape \'XX
			if data[i] == '\'' {
				i += 3
				continue
			}
			// newline escape
			if data[i] == '\n' || data[i] == '\r' {
				i++
				continue
			}

			// read control word
			wordStart := i
			for i < n && data[i] >= 'a' && data[i] <= 'z' {
				i++
			}
			word := string(data[wordStart:i])

			// read optional number
			numStr := ""
			if i < n && (data[i] == '-' || (data[i] >= '0' && data[i] <= '9')) {
				numStart := i
				if data[i] == '-' {
					i++
				}
				for i < n && data[i] >= '0' && data[i] <= '9' {
					i++
				}
				numStr = string(data[numStart:i])
			}

			// consume delimiter space
			if i < n && data[i] == ' ' {
				i++
			}

			switch word {
			case "trowd":
				inTable = true
				tableRow = tableRow[:0]
			case "intbl":
				inTable = true
			case "cell":
				if inTable {
					tableRow = append(tableRow, strings.TrimSpace(para.String()))
					para.Reset()
				}
			case "row":
				if inTable {
					out.WriteString("<tr>")
					for _, cell := range tableRow {
						out.WriteString("<td>" + cell + "</td>")
					}
					out.WriteString("</tr>")
					tableRow = tableRow[:0]
					inTable = false
				}
			case "pard":
				if !inTable {
					emit()
				}
				stack[len(stack)-1].bold = false
			case "s":
				if numStr != "" {
					level, _ := strconv.Atoi(numStr)
					if level >= 1 && level <= 6 {
						headingLevel = level
					}
				}
			case "b":
				if numStr == "0" {
					stack[len(stack)-1].bold = false
				} else {
					stack[len(stack)-1].bold = true
				}
			case "par":
				if !inTable {
					emit()
				}
			case "u":
				if numStr != "" {
					code, _ := strconv.ParseInt(numStr, 10, 32)
					if code < 0 {
						code += 65536
					}
					para.WriteRune(rune(code))
					// skip replacement char \'XX
					if i < n && data[i] == '\\' && i+1 < n && data[i+1] == '\'' {
						i += 4
					}
				}
			}

		case '\n', '\r':
			i++

		default:
			textStart := i
			for i < n && data[i] != '{' && data[i] != '}' && data[i] != '\\' && data[i] != '\n' && data[i] != '\r' {
				i++
			}
			text := string(data[textStart:i])
			if strings.TrimSpace(text) != "" {
				if stack[len(stack)-1].bold {
					para.WriteString("<strong>" + text + "</strong>")
				} else {
					para.WriteString(text)
				}
			}
		}
	}

	emit()
	out.WriteString("</body></html>")
	return out.String()
}
