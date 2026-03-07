package renderer

import (
	"os"
	"strings"

	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/glamour/ansi"
	"github.com/charmbracelet/glamour/styles"
	"golang.org/x/term"
)

func ptr[T any](v T) *T { return &v }

// slideStyle returns a glamour StyleConfig tuned for slide presentation:
// controlled left/right margin, compact vertical spacing, no extra document
// padding (we handle top padding ourselves).
//
// Heading hierarchy (no font sizes in terminals, so we use visual weight):
//
//	H1 — large background block, bold, high-contrast        (title)
//	H2 — bold, bright colour, "── " prefix + underline      (section)
//	H3 — bold, muted colour, "  › " prefix                  (subsection)
func slideStyle() ansi.StyleConfig {
	var style ansi.StyleConfig
	if term.IsTerminal(int(os.Stdout.Fd())) {
		style = styles.DarkStyleConfig
	} else {
		style = styles.LightStyleConfig
	}

	style.Document.Margin = ptr(uint(3))
	style.Document.BlockPrefix = ""
	style.Document.BlockSuffix = ""

	// Clear the catch-all heading style so each level is fully independent.
	style.Heading = ansi.StyleBlock{}

	// H1: bold, high-contrast background, single-space padding.
	style.H1 = ansi.StyleBlock{
		StylePrimitive: ansi.StylePrimitive{
			Prefix:          "  ",
			Suffix:          "  ",
			Color:           ptr("255"),
			BackgroundColor: ptr("#C71585"),
			Bold:            ptr(true),
			BlockSuffix:     "\n",
		},
	}

	// H2: bold, mid-tone background, two-space padding.
	style.H2 = ansi.StyleBlock{
		StylePrimitive: ansi.StylePrimitive{
			Prefix:          "    ",
			Suffix:          " ",
			Color:           ptr("255"),
			BackgroundColor: ptr("#9A0F98"),
			Bold:            ptr(true),
			BlockSuffix:     "\n",
		},
	}

	// H3: bold, subtle background, three-space padding.
	style.H3 = ansi.StyleBlock{
		StylePrimitive: ansi.StylePrimitive{
			Prefix:          "       ",
			Suffix:          " ",
			Color:           ptr("255"),
			BackgroundColor: ptr("#6A0572"),
			Bold:            ptr(true),
			BlockSuffix:     "\n",
		},
	}

	style.CodeBlock.Margin = ptr(uint(0))
	style.BlockQuote.Margin = ptr(uint(0))

	// Syntax highlighting: github-dark palette.
	style.CodeBlock.Chroma = &ansi.Chroma{
		Text:                ansi.StylePrimitive{Color: ptr("#e6edf3")},
		Error:               ansi.StylePrimitive{Color: ptr("#f85149")},
		Comment:             ansi.StylePrimitive{Color: ptr("#8b949e"), Italic: ptr(true)},
		CommentPreproc:      ansi.StylePrimitive{Color: ptr("#8b949e"), Bold: ptr(true)},
		Keyword:             ansi.StylePrimitive{Color: ptr("#ff7b72")},
		KeywordReserved:     ansi.StylePrimitive{Color: ptr("#ff7b72")},
		KeywordNamespace:    ansi.StylePrimitive{Color: ptr("#ff7b72")},
		KeywordType:         ansi.StylePrimitive{Color: ptr("#79c0ff")},
		Operator:            ansi.StylePrimitive{Color: ptr("#ff7b72"), Bold: ptr(true)},
		Punctuation:         ansi.StylePrimitive{Color: ptr("#e6edf3")},
		Name:                ansi.StylePrimitive{Color: ptr("#e6edf3")},
		NameBuiltin:         ansi.StylePrimitive{Color: ptr("#79c0ff")},
		NameTag:             ansi.StylePrimitive{Color: ptr("#7ee787")},
		NameAttribute:       ansi.StylePrimitive{Color: ptr("#79c0ff")},
		NameClass:           ansi.StylePrimitive{Color: ptr("#f0883e"), Bold: ptr(true)},
		NameFunction:        ansi.StylePrimitive{Color: ptr("#d2a8ff"), Bold: ptr(true)},
		NameConstant:        ansi.StylePrimitive{Color: ptr("#79c0ff"), Bold: ptr(true)},
		NameDecorator:       ansi.StylePrimitive{Color: ptr("#d2a8ff"), Bold: ptr(true)},
		LiteralNumber:       ansi.StylePrimitive{Color: ptr("#a5d6ff")},
		LiteralString:       ansi.StylePrimitive{Color: ptr("#a5d6ff")},
		LiteralStringEscape: ansi.StylePrimitive{Color: ptr("#79c0ff")},
		GenericDeleted:      ansi.StylePrimitive{Color: ptr("#ffa198")},
		GenericEmph:         ansi.StylePrimitive{Color: ptr("#e6edf3"), Italic: ptr(true)},
		GenericInserted:     ansi.StylePrimitive{Color: ptr("#56d364")},
		GenericStrong:       ansi.StylePrimitive{Color: ptr("#e6edf3"), Bold: ptr(true)},
		GenericSubheading:   ansi.StylePrimitive{Color: ptr("#79c0ff"), Bold: ptr(true)},
	}

	return style
}

// RenderSlide renders a markdown slide to a terminal-ready string.
// All newlines are converted to \r\n for raw-mode terminals.
func RenderSlide(markdown string, width int) string {
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
