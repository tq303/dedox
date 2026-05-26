package parse

import (
	"regexp"
	"strings"
)

type FilterFunc func([]string) []string

var Filters = map[string]FilterFunc{
	"pii":         redactPII,
	"urls":        redactURLs,
	"ip":          redactIPs,
	"boilerplate": stripBoilerplate,
	"normalize":   normalizeWhitespace,
	"uniq":        uniq,
}

func uniq(lines []string) []string {
	seen := map[string]bool{}
	out := lines[:0:len(lines)]
	for _, line := range lines {
		if !seen[line] {
			seen[line] = true
			out = append(out, line)
		}
	}
	return out
}

var (
	reEmail      = regexp.MustCompile(`[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}`)
	rePhone      = regexp.MustCompile(`(\+?1[\s\-.]?)?\(?\d{3}\)?[\s\-.]?\d{3}[\s\-.]?\d{4}`)
	reSSN        = regexp.MustCompile(`\b\d{3}-\d{2}-\d{4}\b`)
	reCCNum      = regexp.MustCompile(`\b(?:\d[ \-]?){13,16}\b`)
	reURL        = regexp.MustCompile(`https?://[^\s]+`)
	reIP         = regexp.MustCompile(`\b\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}\b`)
	rePageNum    = regexp.MustCompile(`(?i)^\s*page\s+\d+(\s+of\s+\d+)?\s*$`)
	reMultiSpace = regexp.MustCompile(`[ \t]{2,}`)
)

var boilerplatePatterns = []*regexp.Regexp{
	rePageNum,
	regexp.MustCompile(`(?i)^\s*confidential\s*$`),
	regexp.MustCompile(`(?i)all rights reserved`),
	regexp.MustCompile(`(?i)©\s*\d{4}`),
	regexp.MustCompile(`(?i)copyright\s+\d{4}`),
}

func redactPII(lines []string) []string {
	out := make([]string, len(lines))
	for i, line := range lines {
		line = reEmail.ReplaceAllString(line, "[EMAIL]")
		line = rePhone.ReplaceAllString(line, "[PHONE]")
		line = reSSN.ReplaceAllString(line, "[SSN]")
		line = reCCNum.ReplaceAllString(line, "[CC]")
		out[i] = line
	}
	return out
}

func redactURLs(lines []string) []string {
	out := make([]string, len(lines))
	for i, line := range lines {
		out[i] = reURL.ReplaceAllString(line, "[URL]")
	}
	return out
}

func redactIPs(lines []string) []string {
	out := make([]string, len(lines))
	for i, line := range lines {
		out[i] = reIP.ReplaceAllString(line, "[IP]")
	}
	return out
}

func stripBoilerplate(lines []string) []string {
	out := lines[:0:len(lines)]
	for _, line := range lines {
		match := false
		for _, re := range boilerplatePatterns {
			if re.MatchString(line) {
				match = true
				break
			}
		}
		if !match {
			out = append(out, line)
		}
	}
	return out
}

func normalizeWhitespace(lines []string) []string {
	out := make([]string, 0, len(lines))
	blankRun := 0
	for _, line := range lines {
		line = strings.TrimRight(reMultiSpace.ReplaceAllString(line, " "), " \t")
		if line == "" {
			blankRun++
			if blankRun <= 1 {
				out = append(out, line)
			}
		} else {
			blankRun = 0
			out = append(out, line)
		}
	}
	return out
}
