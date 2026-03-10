package render

import (
	"testing"

	"github.com/estifanos/bx/grid"
)

func TestSingleBox(t *testing.T) {
	b := Box{Label: "API", Row: 0, Col: 0, Style: Single}
	g := grid.New(b.Width(), b.Height())
	b.Draw(g)

	want := "┌─────┐\n│ API │\n└─────┘\n"
	got := g.String()
	if got != want {
		t.Errorf("got:\n%s\nwant:\n%s", got, want)
	}
}

func TestRoundedBox(t *testing.T) {
	b := Box{Label: "Server", Row: 0, Col: 0, Style: Rounded}
	g := grid.New(b.Width(), b.Height())
	b.Draw(g)

	want := "╭────────╮\n│ Server │\n╰────────╯\n"
	got := g.String()
	if got != want {
		t.Errorf("got:\n%s\nwant:\n%s", got, want)
	}
}

func TestDoubleBox(t *testing.T) {
	b := Box{Label: "Auth", Row: 0, Col: 0, Style: Double}
	g := grid.New(b.Width(), b.Height())
	b.Draw(g)

	want := "╔══════╗\n║ Auth ║\n╚══════╝\n"
	got := g.String()
	if got != want {
		t.Errorf("got:\n%s\nwant:\n%s", got, want)
	}
}

func TestTwoBoxes(t *testing.T) {
	b1 := Box{Label: "Frontend", Row: 0, Col: 0, Style: Single}
	b2 := Box{Label: "Backend", Row: 0, Col: 16, Style: Single}

	g := grid.New(30, 3)
	b1.Draw(g)
	b2.Draw(g)

	want := "┌──────────┐    ┌─────────┐\n│ Frontend │    │ Backend │\n└──────────┘    └─────────┘\n"
	got := g.String()
	if got != want {
		t.Errorf("got:\n%s\nwant:\n%s", got, want)
	}
}
