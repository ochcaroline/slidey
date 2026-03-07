package slides

import (
	"regexp"
	"strings"
)

var htmlComment = regexp.MustCompile(`(?s)<!--.*?-->`)

// Metadata holds values parsed from the YAML frontmatter.
type Metadata struct {
	Title       string
	Description string
	Author      string
}

// ParseSlides strips YAML frontmatter and HTML comments, splits on "---"
// separators, and returns the metadata together with the non-empty slides.
func ParseSlides(content string) (Metadata, []string) {
	fm, rest := parseFrontmatter(content)
	meta := parseMetadata(fm)

	var slides []string
	for s := range strings.SplitSeq(rest, "\n---\n") {
		s = strings.TrimSpace(htmlComment.ReplaceAllString(s, ""))
		if s != "" {
			slides = append(slides, s)
		}
	}
	return meta, slides
}

// parseMetadata reads key: value lines from the raw frontmatter block.
func parseMetadata(fm string) Metadata {
	var m Metadata
	for line := range strings.SplitSeq(fm, "\n") {
		k, v, ok := strings.Cut(line, ":")
		if !ok {
			continue
		}
		v = strings.TrimSpace(v)
		switch strings.TrimSpace(k) {
		case "title":
			m.Title = v
		case "description":
			m.Description = v
		case "author":
			m.Author = v
		}
	}
	return m
}

// parseFrontmatter extracts YAML frontmatter and returns it together with the
// remaining content. Returns empty fm if no frontmatter is found.
func parseFrontmatter(s string) (fm string, rest string) {
	if !strings.HasPrefix(s, "---\n") {
		return "", s
	}
	end := strings.Index(s[4:], "\n---")
	if end == -1 {
		return "", s
	}
	return s[4 : 4+end], strings.TrimSpace(s[4+end+4:])
}
