package render

import (
	"github.com/estifanos/bx/grid"
)

// Direction indicates the side of a box a connection attaches to.
type Direction int

const (
	Right Direction = iota
	Left
	Down
	Up
)

// ArrowHead runes: the arrow points *toward* the target box,
// so entering from the Left means the arrow traveled Right (▶).
var arrowHeads = map[Direction]rune{
	Left:  '▶', // entering target from its left = traveled right
	Right: '◀', // entering target from its right = traveled left
	Up:    '▼', // entering target from its top = traveled down
	Down:  '▲', // entering target from its bottom = traveled up
}

// Connection draws a line between two boxes.
type Connection struct {
	From     *Box
	To       *Box
	FromSide Direction // which side of From the line leaves
	ToSide   Direction // which side of To the line arrives
	AllBoxes []*Box    // all boxes in the diagram (for collision avoidance)
}

// anchorPoint returns the (row, col) where a connection attaches to a box side.
func anchorPoint(b *Box, side Direction) (int, int) {
	midRow := b.Row + 1 // vertical center (height is always 3)
	midCol := b.Col + b.Width()/2
	switch side {
	case Right:
		return midRow, b.Col + b.Width() // one past the right border
	case Left:
		return midRow, b.Col - 1 // one before the left border
	case Down:
		return b.Row + 3, midCol // one below the bottom border
	case Up:
		return b.Row - 1, midCol // one above the top border
	}
	return midRow, midCol
}

// Draw renders the connection onto the grid.
// Supports straight horizontal, straight vertical, and L-shaped routes.
func (cn *Connection) Draw(g *grid.Grid) {
	r1, c1 := anchorPoint(cn.From, cn.FromSide)
	r2, c2 := anchorPoint(cn.To, cn.ToSide)

	// Determine if this is horizontal, vertical, or L-shaped
	if r1 == r2 {
		cn.drawHorizontal(g, r1, c1, c2)
	} else if c1 == c2 {
		cn.drawVertical(g, c1, r1, r2)
	} else {
		cn.drawLShaped(g, r1, c1, r2, c2)
	}
}

func (cn *Connection) drawHorizontal(g *grid.Grid, row, c1, c2 int) {
	start, end := c1, c2
	if c1 > c2 {
		start, end = c2, c1
	}
	for c := start; c <= end; c++ {
		g.Set(row, c, '─')
	}
	// Place arrowhead at the destination end
	ar, ac := anchorPoint(cn.To, cn.ToSide)
	g.Set(ar, ac, arrowHeads[cn.ToSide])
}

func (cn *Connection) drawVertical(g *grid.Grid, col, r1, r2 int) {
	start, end := r1, r2
	if r1 > r2 {
		start, end = r2, r1
	}
	for r := start; r <= end; r++ {
		g.Set(r, col, '│')
	}
	ar, ac := anchorPoint(cn.To, cn.ToSide)
	g.Set(ar, ac, arrowHeads[cn.ToSide])
}

// isInsideBox returns true if (r, c) is within the bounds of any box.
func (cn *Connection) isInsideBox(r, c int) bool {
	boxes := cn.AllBoxes
	if len(boxes) == 0 {
		boxes = []*Box{cn.From, cn.To}
	}
	for _, b := range boxes {
		if r >= b.Row && r < b.Row+b.Height() && c >= b.Col && c < b.Col+b.Width() {
			return true
		}
	}
	return false
}

// drawLShaped draws an L-shaped route.
// For vertical-primary connections (Down/Up from source), go vertical first then horizontal.
// For horizontal-primary connections (Left/Right from source), go horizontal first then vertical.
func (cn *Connection) drawLShaped(g *grid.Grid, r1, c1, r2, c2 int) {
	var bendR, bendC int

	if cn.FromSide == Down || cn.FromSide == Up {
		// Vertical first: go from (r1,c1) down/up to (r2,c1), then horizontal to (r2,c2)
		bendR, bendC = r2, c1

		// Vertical segment
		vStart, vEnd := r1, bendR
		if r1 > bendR {
			vStart, vEnd = bendR, r1
		}
		for r := vStart; r <= vEnd; r++ {
			g.Set(r, c1, '│')
		}

		// Horizontal segment
		hStart, hEnd := bendC, c2
		if bendC > c2 {
			hStart, hEnd = c2, bendC
		}
		for c := hStart; c <= hEnd; c++ {
			g.Set(bendR, c, '─')
		}
	} else {
		// Horizontal first: go from (r1,c1) to (r1,c2), then vertical to (r2,c2)
		bendR, bendC = r1, c2

		// Horizontal segment
		hStart, hEnd := c1, bendC
		if c1 > bendC {
			hStart, hEnd = bendC, c1
		}
		for c := hStart; c <= hEnd; c++ {
			g.Set(r1, c, '─')
		}

		// Vertical segment
		vStart, vEnd := bendR, r2
		if bendR > r2 {
			vStart, vEnd = r2, bendR
		}
		for r := vStart; r <= vEnd; r++ {
			g.Set(r, c2, '│')
		}
	}

	// Corner glyph at bend point
	g.Set(bendR, bendC, cn.cornerGlyph(bendR, c1, r1, c2))

	// Arrowhead at destination
	ar, ac := anchorPoint(cn.To, cn.ToSide)
	g.Set(ar, ac, arrowHeads[cn.ToSide])
}

// cornerGlyph determines the corner character at a bend point.
// bendR, bendC is the bend location. We need to know which directions
// the lines go from the bend to pick the right glyph.
func (cn *Connection) cornerGlyph(bendR, bendC, otherR, otherC int) rune {
	// Determine which two directions meet at the bend
	hasUp := otherR < bendR || (cn.ToSide == Up && bendR > otherR)
	hasDown := otherR > bendR || (cn.ToSide == Down && bendR < otherR)
	hasLeft := otherC < bendC
	hasRight := otherC > bendC

	// Also check the other direction from From/To anchors
	ar, ac := anchorPoint(cn.To, cn.ToSide)
	if ar < bendR {
		hasUp = true
	}
	if ar > bendR {
		hasDown = true
	}
	if ac < bendC {
		hasLeft = true
	}
	if ac > bendC {
		hasRight = true
	}

	fr, fc := anchorPoint(cn.From, cn.FromSide)
	if fr < bendR {
		hasUp = true
	}
	if fr > bendR {
		hasDown = true
	}
	if fc < bendC {
		hasLeft = true
	}
	if fc > bendC {
		hasRight = true
	}

	_ = ar
	_ = ac
	_ = fr
	_ = fc

	switch {
	case hasDown && hasRight:
		return '┌'
	case hasDown && hasLeft:
		return '┐'
	case hasUp && hasRight:
		return '└'
	case hasUp && hasLeft:
		return '┘'
	default:
		return '┼'
	}
}
