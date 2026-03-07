package main

import (
"fmt"
"os"

"golang.org/x/term"
)

// Presenter manages the terminal presentation loop.
type Presenter struct {
slides   []string
current  int
oldState *term.State
}

func NewPresenter(slides []string) (*Presenter, error) {
oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
if err != nil {
return nil, err
}
return &Presenter{slides: slides, oldState: oldState}, nil
}

func (p *Presenter) Close() {
term.Restore(int(os.Stdin.Fd()), p.oldState)
fmt.Print("\033[H\033[2J")
}

func (p *Presenter) draw() {
w, h, err := term.GetSize(int(os.Stdout.Fd()))
if err != nil {
w, h = 80, 24
}

out := renderSlide(p.slides[p.current], w)

fmt.Print("\033[H\033[2J")
fmt.Print("\r\n\r\n") // top padding
fmt.Print(out)
fmt.Printf("\033[%d;1H  [%d/%d]", h, p.current+1, len(p.slides))
}

func (p *Presenter) next() {
if p.current < len(p.slides)-1 {
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
