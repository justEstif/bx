// Package layout snaps parsed boxes to a clean grid while preserving spatial intent.
package layout

import (
	"sort"

	"github.com/justEstif/bx/parse"
	"github.com/justEstif/bx/render"
)

const (
	hGap = 4 // horizontal gap between boxes
	vGap = 2 // vertical gap between boxes (for connection lines)
)

// LayoutResult contains positioned boxes and connections ready to render.
type LayoutResult struct {
	Boxes       []render.Box
	Connections []render.Connection
	Width       int
	Height      int
}

// Snap takes a parsed diagram and produces a cleanly laid out version.
// It preserves relative spatial arrangement: boxes that were left stay left,
// boxes that were above stay above.
func Snap(d *parse.Diagram, style render.Style) LayoutResult {
	if len(d.Boxes) == 0 {
		return LayoutResult{}
	}

	// Assign each box a render.Box with the correct style
	rboxes := make([]render.Box, len(d.Boxes))
	for i, pb := range d.Boxes {
		rboxes[i] = render.Box{
			Label: pb.Label,
			Style: style,
		}
	}

	// Sort boxes into rows and columns based on original positions.
	// Group boxes by approximate row (boxes within 3 rows of each other = same row).
	type posBox struct {
		idx    int
		origR  int
		origC  int
	}
	pbs := make([]posBox, len(d.Boxes))
	for i, pb := range d.Boxes {
		pbs[i] = posBox{idx: i, origR: pb.Row, origC: pb.Col}
	}

	// Assign row groups: boxes whose original rows overlap or are within
	// a small gap belong to the same visual row.
	sort.Slice(pbs, func(i, j int) bool { return pbs[i].origR < pbs[j].origR })
	rowGroups := [][]posBox{{pbs[0]}}
	for i := 1; i < len(pbs); i++ {
		lastGroup := rowGroups[len(rowGroups)-1]
		// Find the bottom-most row of any box in the current group
		groupBottom := 0
		for _, pb := range lastGroup {
			bottom := d.Boxes[pb.idx].Row + d.Boxes[pb.idx].Height
			if bottom > groupBottom {
				groupBottom = bottom
			}
		}
		if pbs[i].origR < groupBottom {
			rowGroups[len(rowGroups)-1] = append(lastGroup, pbs[i])
		} else {
			rowGroups = append(rowGroups, []posBox{pbs[i]})
		}
	}

	// Sort each row group by column, then assign positions
	positioned := map[int]bool{} // track which boxes have been positioned
	curRow := 0
	for _, group := range rowGroups {
		sort.Slice(group, func(i, j int) bool { return group[i].origC < group[j].origC })

		// Find max height in this row
		maxH := 0
		for _, pb := range group {
			h := rboxes[pb.idx].Height()
			if h > maxH {
				maxH = h
			}
		}

		// Place boxes preserving their relative column positions.
		// Use the first box's original col as a reference.
		if len(group) == 1 {
			// Single box in row: try to center it under its connected box
			pb := group[0]
			rboxes[pb.idx].Row = curRow // set row BEFORE findBestCol
			rboxes[pb.idx].Col = findBestCol(d, rboxes, pb.idx, positioned)
			positioned[pb.idx] = true
		} else {
			curCol := 0
			for _, pb := range group {
				rboxes[pb.idx].Row = curRow
				rboxes[pb.idx].Col = curCol
				positioned[pb.idx] = true
				curCol += rboxes[pb.idx].Width() + hGap
			}
		}
		curRow += maxH + vGap
	}

	// Compute grid dimensions
	maxW, maxH := 0, 0
	for _, b := range rboxes {
		r := b.Row + b.Height()
		c := b.Col + b.Width()
		if r > maxH {
			maxH = r
		}
		if c > maxW {
			maxW = c
		}
	}

	// Build box pointer list for collision avoidance (after positions are set)
	allBoxPtrs := make([]*render.Box, len(rboxes))
	for i := range rboxes {
		allBoxPtrs[i] = &rboxes[i]
	}

	// Build connections
	var rconns []render.Connection
	for _, pc := range d.Connections {
		from := &rboxes[pc.FromBox]
		to := &rboxes[pc.ToBox]
		fromSide, toSide := inferSides(from, to)
		rconns = append(rconns, render.Connection{
			From:     from,
			To:       to,
			FromSide: fromSide,
			ToSide:   toSide,
			AllBoxes: allBoxPtrs,
		})
	}

	return LayoutResult{
		Boxes:       rboxes,
		Connections: rconns,
		Width:       maxW + 2,
		Height:      maxH + 2,
	}
}

// findBestCol finds the best column for a single box in a row by centering it
// under/over boxes it's connected to that are already positioned.
func findBestCol(d *parse.Diagram, rboxes []render.Box, idx int, positioned map[int]bool) int {
	// Find connected boxes that have already been positioned
	var connectedCols []int
	for _, conn := range d.Connections {
		other := -1
		if conn.FromBox == idx {
			other = conn.ToBox
		} else if conn.ToBox == idx {
			other = conn.FromBox
		}
		if other >= 0 && positioned[other] {
			center := rboxes[other].Col + rboxes[other].Width()/2
			connectedCols = append(connectedCols, center)
		}
	}
	if len(connectedCols) > 0 {
		avg := 0
		for _, c := range connectedCols {
			avg += c
		}
		avg /= len(connectedCols)
		// Center our box on that average
		col := avg - rboxes[idx].Width()/2
		if col < 0 {
			col = 0
		}
		return col
	}
	return 0
}

// inferSides determines which sides of the boxes the connection should attach to.
// Priority: if boxes don't overlap vertically (different row groups), use vertical connection.
// If they don't overlap horizontally, use horizontal connection.
func inferSides(from, to *render.Box) (render.Direction, render.Direction) {
	// Check if boxes are in different rows (no vertical overlap)
	fromBottom := from.Row + from.Height()
	toBottom := to.Row + to.Height()
	verticallyDisjoint := to.Row >= fromBottom || from.Row >= toBottom

	// Check if boxes are in different columns (no horizontal overlap)
	fromRight := from.Col + from.Width()
	toRight := to.Col + to.Width()
	horizontallyDisjoint := to.Col >= fromRight || from.Col >= toRight

	if verticallyDisjoint && !horizontallyDisjoint {
		// Same column area, different rows → vertical
		if to.Row > from.Row {
			return render.Down, render.Up
		}
		return render.Up, render.Down
	}
	if horizontallyDisjoint && !verticallyDisjoint {
		// Same row area, different columns → horizontal
		if to.Col > from.Col {
			return render.Right, render.Left
		}
		return render.Left, render.Right
	}

	// Both disjoint — prefer vertical if the row gap is larger or equal
	if verticallyDisjoint && horizontallyDisjoint {
		rowGap := 0
		if to.Row >= fromBottom {
			rowGap = to.Row - fromBottom
		} else {
			rowGap = from.Row - toBottom
		}
		colGap := 0
		if to.Col >= fromRight {
			colGap = to.Col - fromRight
		} else {
			colGap = from.Col - toRight
		}

		if rowGap >= colGap {
			if to.Row > from.Row {
				return render.Down, render.Up
			}
			return render.Up, render.Down
		}
		if to.Col > from.Col {
			return render.Right, render.Left
		}
		return render.Left, render.Right
	}

	// Overlapping — use center distance
	fromCR := from.Row + from.Height()/2
	fromCC := from.Col + from.Width()/2
	toCR := to.Row + to.Height()/2
	toCC := to.Col + to.Width()/2
	dr := toCR - fromCR
	dc := toCC - fromCC
	if dr < 0 {
		dr = -dr
	}
	if dc < 0 {
		dc = -dc
	}
	if dr >= dc {
		if to.Row > from.Row {
			return render.Down, render.Up
		}
		return render.Up, render.Down
	}
	if to.Col > from.Col {
		return render.Right, render.Left
	}
	return render.Left, render.Right
}
