package main

import (
	"fmt"
	"os"

	"github.com/ochcaroline/slidey/internal/presenter"
	"github.com/ochcaroline/slidey/internal/slides"
)

// version is set at build time via -ldflags "-X main.version=x.y.z"
var version = "dev"

func main() {
	if len(os.Args) == 2 && (os.Args[1] == "--version" || os.Args[1] == "-v") {
		fmt.Println("slidey", version)
		return
	}

	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: slidey <file.md>")
		os.Exit(1)
	}

	data, err := os.ReadFile(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	meta, slides := slides.ParseSlides(string(data))
	if len(slides) == 0 && meta.Title == "" {
		fmt.Fprintln(os.Stderr, "no slides found")
		os.Exit(1)
	}

	p, err := presenter.NewPresenter(meta, slides)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer p.Close()

	p.Run()
}
