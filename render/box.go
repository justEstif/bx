package render

import (
	"github.com/justEstif/bx/grid"
)

// Box represents a box to render on the grid.
type Box struct {
	Label string
	Row   int // top-left row position
	Col   int // top-left col position
	Style Style
}

// Width returns the rendered width: 1 padding on each side + border.
// e.g. label "API" → │ API │ → width = 3 + 2 padding + 2 border = 7
func (b *Box) Width() int {
	return len([]rune(b.Label)) + 4 // 2 border + 2 padding
}

// Height returns the rendered height: always 3 (top border, content, bottom border).
func (b *Box) Height() int {
	return 3
}

// Draw renders the box onto the grid.
func (b *Box) Draw(g *grid.Grid) {
	w := b.Width()
	s := b.Style
	r, c := b.Row, b.Col

	// Top border
	g.Set(r, c, s.TopLeft)
	for i := 1; i < w-1; i++ {
		g.Set(r, c+i, s.Horizontal)
	}
	g.Set(r, c+w-1, s.TopRight)

	// Middle: vertical + padding + label + padding + vertical
	g.Set(r+1, c, s.Vertical)
	g.Set(r+1, c+1, ' ')
	for i, ch := range []rune(b.Label) {
		g.Set(r+1, c+2+i, ch)
	}
	g.Set(r+1, c+w-2, ' ')
	g.Set(r+1, c+w-1, s.Vertical)

	// Bottom border
	g.Set(r+2, c, s.BottomLeft)
	for i := 1; i < w-1; i++ {
		g.Set(r+2, c+i, s.Horizontal)
	}
	g.Set(r+2, c+w-1, s.BottomRight)
}
