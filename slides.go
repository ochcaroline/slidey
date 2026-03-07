package main

import (
"regexp"
"strings"
)

var htmlComment = regexp.MustCompile(`(?s)<!--.*?-->`)

// parseSlides strips YAML frontmatter and HTML comments, then splits the
// content on "---" separators, returning only non-empty slides.
func parseSlides(content string) []string {
content = stripFrontmatter(content)
var slides []string
for _, s := range strings.Split(content, "\n---\n") {
s = strings.TrimSpace(htmlComment.ReplaceAllString(s, ""))
if s != "" {
slides = append(slides, s)
}
}
return slides
}

// stripFrontmatter removes YAML frontmatter (--- ... ---) from the top of content.
func stripFrontmatter(s string) string {
if !strings.HasPrefix(s, "---\n") {
return s
}
end := strings.Index(s[4:], "\n---")
if end == -1 {
return s
}
return strings.TrimSpace(s[4+end+4:])
}
