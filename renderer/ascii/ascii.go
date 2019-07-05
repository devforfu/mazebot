package ascii

import (
	"fmt"
	"io"
	"log"
	"maze/generator"
	"os"
	"strconv"
	"strings"
)

type Renderer struct {
	*RenderingOptions
}

type CharsMap map[generator.CellType]string

type RenderingOptions struct {
	ShowBoundary bool
	ShowNumbers bool
	OutputStream io.Writer
	Chars CharsMap
	BadCharReplacement string
}

var DefaultChars = map[generator.CellType]string {
	generator.Empty: " ",
	generator.Wall: "#",
	generator.Start: "@",
	generator.Exit: "G",
	generator.Visited: ".",
}


var DefaultOptions = RenderingOptions{
	ShowBoundary: 		true,
	OutputStream: 		os.Stdout,
	BadCharReplacement: "•",
	Chars:				DefaultChars,
}

var DefaultRenderer = Renderer{&DefaultOptions}

func (r *Renderer) Render(m *generator.Maze) {
	mazeStr := r.toString(m)
	_, _ = fmt.Fprintf(r.OutputStream, mazeStr)
}

func (r *Renderer) toString(m *generator.Maze) string {
	if m.Map == nil { return "Maze is empty!" }

	var b strings.Builder
	b.WriteString(fmt.Sprintf("Maze ID=%s (%dx%d)\n", m.ID, m.Size.X, m.Size.Y))

	nCols, nRows := m.Size.X, m.Size.Y
	if r.ShowBoundary {
		if r.ShowNumbers {
			for i := 0; i < nCols; i++ {
				b.WriteString(strconv.Itoa(i))
			}
		}
		b.WriteString(fmt.Sprintf("+%s+\n", strings.Repeat("-", nCols)))
	}

	for j := 0; j < nRows; j++ {
		if r.ShowBoundary { b.WriteString("|") }
		for i := 0; i < nCols; i++ {
			if item, ok := r.Chars[m.Map[i][j]]; !ok {
				log.Printf("Unexpected cell type encountered at (%d, %d)", i, j)
				b.WriteString(r.BadCharReplacement)
			} else {
				b.WriteString(item)
			}
		}
		if r.ShowBoundary {
			b.WriteString("|")
			b.WriteString("\n")
		}
	}

	if r.ShowBoundary {
		b.WriteString(fmt.Sprintf("+%s+\n", strings.Repeat("-", nCols)))
	}
	return b.String()
}
