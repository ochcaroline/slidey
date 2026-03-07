package main

import (
	"fmt"
	"os"
	"strings"

	"golang.org/x/term"
)

// Presenter manages the terminal presentation loop.
type Presenter struct {
	meta     Metadata
	slides   []string
	current  int
	oldState *term.State
}

func NewPresenter(meta Metadata, slides []string) (*Presenter, error) {
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return nil, err
	}
	return &Presenter{meta: meta, slides: slides, oldState: oldState}, nil
}

func (p *Presenter) Close() {
	term.Restore(int(os.Stdin.Fd()), p.oldState)
	fmt.Print("\033[H\033[2J")
}

func (p *Presenter) hasTitleSlide() bool {
	return p.meta.Title != ""
}

func (p *Presenter) total() int {
	if p.hasTitleSlide() {
		return len(p.slides) + 1
	}
	return len(p.slides)
}

// slideIndex returns the index into p.slides for the current position.
func (p *Presenter) slideIndex() int {
	if p.hasTitleSlide() {
		return p.current - 1
	}
	return p.current
}

func (p *Presenter) draw() {
	w, h, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		w, h = 80, 24
	}

	fmt.Print("\033[H\033[2J")

	if p.hasTitleSlide() && p.current == 0 {
		p.drawTitleSlide(w, h)
	} else {
		p.drawContentSlide(w, h)
	}

	p.drawChrome(w, h)
}

// drawTitleSlide renders the auto-generated centered title/description slide.
func (p *Presenter) drawTitleSlide(w, h int) {
	// Vertical center: place title at ~40% down.
	titleRow := h * 2 / 5
	for range titleRow {
		fmt.Print("\r\n")
	}

	if p.meta.Title != "" {
		// Bold + medium violet red (#C71585).
		title := "\033[1;38;2;199;21;133m" + p.meta.Title + "\033[0m"
		fmt.Print(centerText(title, len(p.meta.Title), w) + "\r\n")
	}

	if p.meta.Description != "" {
		fmt.Print("\r\n")
		fmt.Print(centerText(p.meta.Description, len(p.meta.Description), w) + "\r\n")
	}
}

// drawContentSlide renders a normal markdown slide.
func (p *Presenter) drawContentSlide(w, h int) {
	_ = h
	out := renderSlide(p.slides[p.slideIndex()], w)
	fmt.Print("\r\n\r\n")
	fmt.Print(out)
}

// drawChrome draws the persistent UI elements:
//
//	top-right:    slide counter
//	bottom-left:  author
//	bottom-right: title
func (p *Presenter) drawChrome(w, h int) {
	counter := fmt.Sprintf("[%d/%d]  ", p.current+1, p.total())
	// Top-right.
	fmt.Printf("\033[1;%dH%s", w-len(counter)+1, counter)

	// Bottom-left: author.
	if p.meta.Author != "" {
		fmt.Printf("\033[%d;3H\033[2m%s\033[0m", h, p.meta.Author)
	}

	// Bottom-right: title.
	if p.meta.Title != "" {
		col := max(w-len(p.meta.Title)-2, 1)
		fmt.Printf("\033[%d;%dH\033[2m%s\033[0m", h, col, p.meta.Title)
	}
}

// centerText returns s padded with spaces to be visually centered in a field
// of the given terminal width. visualLen is the printable rune count (no ANSI).
func centerText(s string, visualLen, width int) string {
	pad := max((width-visualLen)/2, 0)
	return strings.Repeat(" ", pad) + s
}

func (p *Presenter) next() {
	if p.current < p.total()-1 {
		p.current++
		p.draw()
	}
}

func (p *Presenter) prev() {
	if p.current > 0 {
		p.current--
		p.draw()
	}
}

// Run starts the presentation event loop and blocks until the user quits.
func (p *Presenter) Run() {
	p.draw()

	buf := make([]byte, 1)
	for {
		os.Stdin.Read(buf)
		switch buf[0] {
		case 'j', ' ':
			p.next()
		case 'k':
			p.prev()
		case 'q', 3: // q or Ctrl-C
			return
		}
	}
}
