package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: slidey <file.md>")
		os.Exit(1)
	}

	data, err := os.ReadFile(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	meta, slides := parseSlides(string(data))
	if len(slides) == 0 && meta.Title == "" {
		fmt.Fprintln(os.Stderr, "no slides found")
		os.Exit(1)
	}

	p, err := NewPresenter(meta, slides)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer p.Close()

	p.Run()
}
