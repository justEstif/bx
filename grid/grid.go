// Package grid provides a character grid for rendering diagrams.
package grid

import (
	"strings"
)

// Grid is a 2D character matrix for rendering diagrams.
type Grid struct {
	cells  [][]rune
	width  int
	height int
}

// New creates a grid of the given dimensions filled with spaces.
func New(width, height int) *Grid {
	cells := make([][]rune, height)
	for r := range cells {
		cells[r] = make([]rune, width)
		for c := range cells[r] {
			cells[r][c] = ' '
		}
	}
	return &Grid{cells: cells, width: width, height: height}
}

// Set places a rune at (row, col). Out-of-bounds writes are ignored.
func (g *Grid) Set(row, col int, ch rune) {
	if row >= 0 && row < g.height && col >= 0 && col < g.width {
		g.cells[row][col] = ch
	}
}

// Get returns the rune at (row, col), or ' ' if out of bounds.
func (g *Grid) Get(row, col int) rune {
	if row >= 0 && row < g.height && col >= 0 && col < g.width {
		return g.cells[row][col]
	}
	return ' '
}

// Width returns the grid width.
func (g *Grid) Width() int { return g.width }

// Height returns the grid height.
func (g *Grid) Height() int { return g.height }

// String renders the grid as a string, trimming trailing blank lines
// and trailing spaces on each line.
func (g *Grid) String() string {
	lines := make([]string, g.height)
	lastNonEmpty := -1
	for r := 0; r < g.height; r++ {
		line := strings.TrimRight(string(g.cells[r]), " ")
		lines[r] = line
		if line != "" {
			lastNonEmpty = r
		}
	}
	if lastNonEmpty < 0 {
		return ""
	}
	return strings.Join(lines[:lastNonEmpty+1], "\n") + "\n"
}
