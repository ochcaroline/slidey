package main

import (
"os"
"strings"

"github.com/charmbracelet/glamour"
"github.com/charmbracelet/glamour/ansi"
"github.com/charmbracelet/glamour/styles"
"golang.org/x/term"
)

func uintPtr(u uint) *uint { return &u }

func boolPtr(b bool) *bool   { return &b }
func strPtr(s string) *string { return &s }

// slideStyle returns a glamour StyleConfig tuned for slide presentation:
// controlled left/right margin, compact vertical spacing, no extra document
// padding (we handle top padding ourselves).
//
// Heading hierarchy (no font sizes in terminals, so we use visual weight):
//   H1 — large background block, bold, high-contrast        (title)
//   H2 — bold, bright colour, "── " prefix + underline      (section)
//   H3 — bold, muted colour, "  › " prefix                  (subsection)
func slideStyle() ansi.StyleConfig {
	var style ansi.StyleConfig
	if term.IsTerminal(int(os.Stdout.Fd())) {
		style = styles.DarkStyleConfig
	} else {
		style = styles.LightStyleConfig
	}

	style.Document.Margin = uintPtr(3)
	style.Document.StylePrimitive.BlockPrefix = ""
	style.Document.StylePrimitive.BlockSuffix = ""

	// Clear the catch-all heading style so each level is fully independent.
	style.Heading = ansi.StyleBlock{}

	// H1: bold, high-contrast background, single-space padding.
	style.H1 = ansi.StyleBlock{
		StylePrimitive: ansi.StylePrimitive{
			Prefix:          " ",
			Suffix:          " ",
			Color:           strPtr("255"),
			BackgroundColor: strPtr("63"),
			Bold:            boolPtr(true),
			BlockSuffix:     "\n",
		},
	}

	// H2: bold, mid-tone background, two-space padding.
	style.H2 = ansi.StyleBlock{
		StylePrimitive: ansi.StylePrimitive{
			Prefix:          "  ",
			Suffix:          "  ",
			Color:           strPtr("255"),
			BackgroundColor: strPtr("25"),
			Bold:            boolPtr(true),
			BlockSuffix:     "\n",
		},
	}

	// H3: bold, subtle background, three-space padding.
	style.H3 = ansi.StyleBlock{
		StylePrimitive: ansi.StylePrimitive{
			Prefix:          "   ",
			Suffix:          "   ",
			Color:           strPtr("255"),
			BackgroundColor: strPtr("23"),
			Bold:            boolPtr(true),
			BlockSuffix:     "\n",
		},
	}

	style.CodeBlock.StyleBlock.Margin = uintPtr(0)
	style.BlockQuote.Margin = uintPtr(0)

	return style
}

// renderSlide renders a markdown slide to a terminal-ready string.
// All newlines are converted to \r\n for raw-mode terminals.
func renderSlide(markdown string, width int) string {
r, err := glamour.NewTermRenderer(
glamour.WithStyles(slideStyle()),
glamour.WithWordWrap(width-6), // account for 3-char margin on each side
)
if err != nil {
r, _ = glamour.NewTermRenderer(glamour.WithAutoStyle())
}

out, err := r.Render(markdown)
if err != nil {
out = markdown
}

// In raw mode \n is line-feed only; normalise to \r\n so the cursor
// returns to column 0 on each new line.
return strings.ReplaceAll(strings.TrimLeft(out, "\n"), "\n", "\r\n")
}
