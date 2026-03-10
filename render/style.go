// Package render draws diagram elements onto a grid.
package render

// Style defines the box-drawing characters for a box style.
type Style struct {
	TopLeft     rune
	TopRight    rune
	BottomLeft  rune
	BottomRight rune
	Horizontal  rune
	Vertical    rune
}

var (
	Single = Style{
		TopLeft: '┌', TopRight: '┐',
		BottomLeft: '└', BottomRight: '┘',
		Horizontal: '─', Vertical: '│',
	}
	Double = Style{
		TopLeft: '╔', TopRight: '╗',
		BottomLeft: '╚', BottomRight: '╝',
		Horizontal: '═', Vertical: '║',
	}
	Rounded = Style{
		TopLeft: '╭', TopRight: '╮',
		BottomLeft: '╰', BottomRight: '╯',
		Horizontal: '─', Vertical: '│',
	}
	Bold = Style{
		TopLeft: '┏', TopRight: '┓',
		BottomLeft: '┗', BottomRight: '┛',
		Horizontal: '━', Vertical: '┃',
	}
)
