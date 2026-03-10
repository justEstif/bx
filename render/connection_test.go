package render

import (
	"testing"

	"github.com/justEstif/bx/grid"
)

func TestHorizontalConnection(t *testing.T) {
	b1 := Box{Label: "A", Row: 0, Col: 0, Style: Single}
	b2 := Box{Label: "B", Row: 0, Col: 12, Style: Single}
	conn := Connection{From: &b1, To: &b2, FromSide: Right, ToSide: Left}

	g := grid.New(20, 3)
	conn.Draw(g)
	b1.Draw(g)
	b2.Draw(g)

	got := g.String()
	t.Logf("horizontal:\n%s", got)

	// Arrow should point right (▶) since it's entering B from the left
	ar, ac := anchorPoint(&b2, Left)
	if g.Get(ar, ac) != '▶' {
		t.Errorf("expected ▶ at (%d,%d), got %c", ar, ac, g.Get(ar, ac))
	}
}

func TestVerticalConnection(t *testing.T) {
	b1 := Box{Label: "Top", Row: 0, Col: 0, Style: Single}
	b2 := Box{Label: "Bot", Row: 6, Col: 0, Style: Single}
	conn := Connection{From: &b1, To: &b2, FromSide: Down, ToSide: Up}

	g := grid.New(10, 10)
	conn.Draw(g)
	b1.Draw(g)
	b2.Draw(g)

	got := g.String()
	t.Logf("vertical:\n%s", got)

	// Arrow should point down (▼) since it's entering Bot from the top
	ar, ac := anchorPoint(&b2, Up)
	if g.Get(ar, ac) != '▼' {
		t.Errorf("expected ▼ at (%d,%d), got %c", ar, ac, g.Get(ar, ac))
	}
}

func TestLShapedConnection(t *testing.T) {
	b1 := Box{Label: "Src", Row: 0, Col: 0, Style: Single}
	b2 := Box{Label: "Dst", Row: 5, Col: 15, Style: Single}
	conn := Connection{From: &b1, To: &b2, FromSide: Right, ToSide: Up}

	g := grid.New(25, 9)
	conn.Draw(g)
	b1.Draw(g)
	b2.Draw(g)

	got := g.String()
	t.Logf("L-shaped:\n%s", got)

	// Arrow should point down (▼) since it's entering Dst from the top
	ar, ac := anchorPoint(&b2, Up)
	if g.Get(ar, ac) != '▼' {
		t.Errorf("expected ▼ at (%d,%d), got %c", ar, ac, g.Get(ar, ac))
	}
}
