# slidey

Terminal markdown slide presenter.

## Usage

```
go build -o slidey .
./slidey slides.md
```

## Navigation

| Key | Action |
|-----|--------|
| `j` / `space` | Next slide |
| `k` | Previous slide |
| `q` / `Ctrl-C` | Quit |

## Slide format

Slides are separated by `---` in your markdown file. YAML frontmatter and HTML comments are stripped automatically.

```markdown
# Slide one

Some content

---

# Slide two

More content
```
