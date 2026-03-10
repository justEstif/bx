// Package parse extracts boxes and connections from rough ASCII/Unicode diagrams.
package parse

import (
	"strings"
)

// Diagram holds parsed diagram elements.
type Diagram struct {
	Boxes       []ParsedBox
	Connections []ParsedConnection
}

// ParsedBox is a box extracted from ASCII input.
type ParsedBox struct {
	Label string
	Row   int // top-left row in original input
	Col   int // top-left col in original input
	Width int // original width (may be wrong)
	Height int
}

// ParsedConnection is a connection extracted from ASCII input.
type ParsedConnection struct {
	FromBox int // index into Boxes
	ToBox   int // index into Boxes
	// Points along the connection path (row, col pairs)
	Points [][2]int
	Arrow  bool // true if has arrowhead
}

// isStructural returns true if the rune is a box-drawing or ASCII structural char.
func isStructural(r rune) bool {
	switch r {
	case '+', '-', '|', '/', '\\',
		'┌', '┐', '└', '┘', '─', '│', '├', '┤', '┬', '┴', '┼',
		'╔', '╗', '╚', '╝', '═', '║', '╠', '╣', '╦', '╩', '╬',
		'╭', '╮', '╰', '╯',
		'┏', '┓', '┗', '┛', '━', '┃':
		return true
	}
	return false
}

// isHorizontal returns true if the rune represents a horizontal line.
func isHorizontal(r rune) bool {
	switch r {
	case '-', '─', '═', '━':
		return true
	}
	return false
}

// isVertical returns true if the rune represents a vertical line.
func isVertical(r rune) bool {
	switch r {
	case '|', '│', '║', '┃':
		return true
	}
	return false
}

// isCorner returns true if the rune is a corner or junction character.
func isCorner(r rune) bool {
	switch r {
	case '+',
		'┌', '┐', '└', '┘', '├', '┤', '┬', '┴', '┼',
		'╔', '╗', '╚', '╝', '╠', '╣', '╦', '╩', '╬',
		'╭', '╮', '╰', '╯',
		'┏', '┓', '┗', '┛':
		return true
	}
	return false
}

// isArrow returns true if the rune is an arrow character.
func isArrow(r rune) bool {
	switch r {
	case '>', '<', 'v', 'V', '^', '▶', '◀', '▼', '▲':
		return true
	}
	return false
}

// Parse parses raw ASCII/Unicode text and extracts boxes and connections.
func Parse(input string) *Diagram {
	lines := toRuneGrid(input)
	if len(lines) == 0 {
		return &Diagram{}
	}

	d := &Diagram{}
	visited := makeVisited(lines)

	// Phase 1: Find boxes by scanning for top-left corners
	for r := 0; r < len(lines); r++ {
		for c := 0; c < len(lines[r]); c++ {
			if visited[r][c] {
				continue
			}
			ch := lines[r][c]
			if isCorner(ch) || ch == '+' {
				if box, ok := tryExtractBox(lines, r, c); ok {
					d.Boxes = append(d.Boxes, box)
					markBoxVisited(visited, box)
				}
			}
		}
	}

	// Phase 2: Find connections (horizontal/vertical lines between boxes)
	d.Connections = findConnections(lines, d.Boxes)

	return d
}

func toRuneGrid(input string) [][]rune {
	rawLines := strings.Split(input, "\n")
	// Remove trailing empty line from final newline
	if len(rawLines) > 0 && rawLines[len(rawLines)-1] == "" {
		rawLines = rawLines[:len(rawLines)-1]
	}
	grid := make([][]rune, len(rawLines))
	for i, line := range rawLines {
		grid[i] = []rune(line)
	}
	return grid
}

func makeVisited(lines [][]rune) [][]bool {
	v := make([][]bool, len(lines))
	for i := range lines {
		v[i] = make([]bool, len(lines[i]))
	}
	return v
}

func markBoxVisited(visited [][]bool, box ParsedBox) {
	for r := box.Row; r < box.Row+box.Height && r < len(visited); r++ {
		for c := box.Col; c < box.Col+box.Width && c < len(visited[r]); c++ {
			visited[r][c] = true
		}
	}
}

func getCell(lines [][]rune, r, c int) rune {
	if r < 0 || r >= len(lines) {
		return ' '
	}
	if c < 0 || c >= len(lines[r]) {
		return ' '
	}
	return lines[r][c]
}

// tryExtractBox attempts to find a rectangular box starting at (r, c) as top-left corner.
// It's tolerant of slight misalignment: if the top edge gives a width, it also checks
// content rows for vertical bars that might be slightly off.
func tryExtractBox(lines [][]rune, r, c int) (ParsedBox, bool) {
	// Walk down from top-left to find bottom-left corner.
	bottomRow := -1
	for tr := r + 1; tr < len(lines); tr++ {
		ch := getCell(lines, tr, c)
		if isCorner(ch) || ch == '+' {
			bottomRow = tr
			break
		}
		if !isVertical(ch) {
			break
		}
	}
	if bottomRow < 0 || bottomRow-r < 2 { // minimum box height
		return ParsedBox{}, false
	}

	height := bottomRow - r + 1

	// Walk right along the top edge looking for a candidate top-right corner.
	// Some malformed boxes contain stray '+' characters on the top edge, so we
	// keep scanning until we find the first candidate that also has a matching
	// bottom-right corner.
	actualRight := -1
	for tc := c + 1; tc < len(lines[r]); tc++ {
		ch := getCell(lines, r, tc)
		if !isHorizontal(ch) && !isCorner(ch) && ch != ' ' && !isArrow(ch) {
			break
		}
		if !isCorner(ch) && ch != '+' {
			continue
		}
		if tc-c < 3 { // minimum box width
			continue
		}

		// Verify bottom-right corner near this candidate (tolerance for misalignment).
		brCol := -1
		for delta := -2; delta <= 2; delta++ {
			bc := tc + delta
			brch := getCell(lines, bottomRow, bc)
			if isCorner(brch) || brch == '+' {
				brCol = bc
				break
			}
		}
		if brCol < 0 {
			continue
		}

		candidateRight := max(tc, brCol)

		// Verify bottom edge (from c+1 to candidateRight-1, tolerant).
		bottomValid := true
		for bc := c + 1; bc < candidateRight; bc++ {
			bch := getCell(lines, bottomRow, bc)
			if !isHorizontal(bch) && !isCorner(bch) && bch != ' ' && !isArrow(bch) {
				bottomValid = false
				break
			}
		}
		if !bottomValid {
			continue
		}

		actualRight = candidateRight
		break
	}
	if actualRight < 0 {
		return ParsedBox{}, false
	}

	width := actualRight - c + 1

	// Extract label: scan content rows, find vertical bars on each side (with tolerance)
	var labelParts []string
	for tr := r + 1; tr < bottomRow; tr++ {
		// Find the left vertical bar (should be at c, but check c±1)
		leftV := -1
		for delta := -1; delta <= 1; delta++ {
			if isVertical(getCell(lines, tr, c+delta)) {
				leftV = c + delta
				break
			}
		}
		// Find the right vertical bar (near actualRight)
		rightV := -1
		for delta := -2; delta <= 2; delta++ {
			tc := actualRight + delta
			if isVertical(getCell(lines, tr, tc)) {
				rightV = tc
				break
			}
		}

		if leftV >= 0 && rightV >= 0 && rightV > leftV {
			var part []rune
			for tc := leftV + 1; tc < rightV; tc++ {
				ch := getCell(lines, tr, tc)
				part = append(part, ch)
			}
			s := strings.TrimSpace(string(part))
			if s != "" {
				labelParts = append(labelParts, s)
			}
		}
	}
	label := strings.Join(labelParts, " ")

	return ParsedBox{
		Label:  label,
		Row:    r,
		Col:    c,
		Width:  width,
		Height: height,
	}, true
}

// findConnections scans for line segments between boxes.
func findConnections(lines [][]rune, boxes []ParsedBox) []ParsedConnection {
	var conns []ParsedConnection

	for i := range boxes {
		for j := range boxes {
			if i == j {
				continue
			}
			if conn, ok := findHorizontalConn(lines, boxes, i, j); ok {
				conns = append(conns, conn)
			}
			if conn, ok := findVerticalConn(lines, boxes, i, j); ok {
				conns = append(conns, conn)
			}
		}
	}

	return dedupConnections(conns)
}

func findHorizontalConn(lines [][]rune, boxes []ParsedBox, fromIdx, toIdx int) (ParsedConnection, bool) {
	from := boxes[fromIdx]
	to := boxes[toIdx]

	// Check if 'to' is to the right of 'from' and they overlap vertically
	if to.Col <= from.Col+from.Width {
		return ParsedConnection{}, false
	}

	// Find overlapping rows (the content rows of each box)
	overlapStart := max(from.Row, to.Row)
	overlapEnd := min(from.Row+from.Height-1, to.Row+to.Height-1)

	for row := overlapStart; row <= overlapEnd; row++ {
		startCol := from.Col + from.Width
		endCol := to.Col - 1

		if startCol > endCol {
			continue
		}

		// Check if there's a continuous horizontal line or arrow path
		valid := true
		hasArrow := false
		for c := startCol; c <= endCol; c++ {
			ch := getCell(lines, row, c)
			if isHorizontal(ch) || ch == ' ' || isVertical(ch) {
				continue
			}
			if isArrow(ch) {
				hasArrow = true
				continue
			}
			valid = false
			break
		}

		// Need at least one non-space drawing char
		if !valid {
			continue
		}
		hasDrawing := false
		for c := startCol; c <= endCol; c++ {
			ch := getCell(lines, row, c)
			if isHorizontal(ch) || isArrow(ch) {
				hasDrawing = true
				break
			}
		}
		if !hasDrawing {
			continue
		}

		return ParsedConnection{
			FromBox: fromIdx,
			ToBox:   toIdx,
			Arrow:   hasArrow,
		}, true
	}
	return ParsedConnection{}, false
}

func findVerticalConn(lines [][]rune, boxes []ParsedBox, fromIdx, toIdx int) (ParsedConnection, bool) {
	from := boxes[fromIdx]
	to := boxes[toIdx]

	// Check if 'to' is below 'from' and they overlap horizontally
	if to.Row <= from.Row+from.Height {
		return ParsedConnection{}, false
	}

	overlapStart := max(from.Col, to.Col)
	overlapEnd := min(from.Col+from.Width-1, to.Col+to.Width-1)

	for col := overlapStart; col <= overlapEnd; col++ {
		startRow := from.Row + from.Height
		endRow := to.Row - 1

		if startRow > endRow {
			continue
		}

		valid := true
		hasArrow := false
		hasDrawing := false
		for r := startRow; r <= endRow; r++ {
			ch := getCell(lines, r, col)
			if isVertical(ch) {
				hasDrawing = true
				continue
			}
			if isArrow(ch) {
				hasArrow = true
				hasDrawing = true
				continue
			}
			if ch == ' ' {
				continue
			}
			valid = false
			break
		}

		if valid && hasDrawing {
			return ParsedConnection{
				FromBox: fromIdx,
				ToBox:   toIdx,
				Arrow:   hasArrow,
			}, true
		}
	}
	return ParsedConnection{}, false
}

func dedupConnections(conns []ParsedConnection) []ParsedConnection {
	type key struct{ a, b int }
	seen := map[key]bool{}
	var result []ParsedConnection
	for _, c := range conns {
		k := key{min(c.FromBox, c.ToBox), max(c.FromBox, c.ToBox)}
		if !seen[k] {
			seen[k] = true
			result = append(result, c)
		}
	}
	return result
}


