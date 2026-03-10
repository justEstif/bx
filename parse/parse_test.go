package parse

import (
	"testing"
)

func TestParseSingleBox(t *testing.T) {
	input := `+--------+
| Server |
+--------+`

	d := Parse(input)
	if len(d.Boxes) != 1 {
		t.Fatalf("expected 1 box, got %d", len(d.Boxes))
	}
	if d.Boxes[0].Label != "Server" {
		t.Errorf("expected label 'Server', got '%s'", d.Boxes[0].Label)
	}
}

func TestParseUnicodeBox(t *testing.T) {
	input := `┌────────┐
│ Server │
└────────┘`

	d := Parse(input)
	if len(d.Boxes) != 1 {
		t.Fatalf("expected 1 box, got %d", len(d.Boxes))
	}
	if d.Boxes[0].Label != "Server" {
		t.Errorf("expected label 'Server', got '%s'", d.Boxes[0].Label)
	}
}

func TestParseTwoBoxesWithHorizontalConnection(t *testing.T) {
	input := `+-----+       +-----+
|  A  |------>|  B  |
+-----+       +-----+`

	d := Parse(input)
	if len(d.Boxes) != 2 {
		t.Fatalf("expected 2 boxes, got %d", len(d.Boxes))
	}
	if d.Boxes[0].Label != "A" {
		t.Errorf("box 0 label: expected 'A', got '%s'", d.Boxes[0].Label)
	}
	if d.Boxes[1].Label != "B" {
		t.Errorf("box 1 label: expected 'B', got '%s'", d.Boxes[1].Label)
	}
	if len(d.Connections) != 1 {
		t.Fatalf("expected 1 connection, got %d", len(d.Connections))
	}
	if !d.Connections[0].Arrow {
		t.Error("expected arrow connection")
	}
}

func TestParseTwoBoxesVertical(t *testing.T) {
	input := `+-----+
|  A  |
+-----+
   |
   v
+-----+
|  B  |
+-----+`

	d := Parse(input)
	if len(d.Boxes) != 2 {
		t.Fatalf("expected 2 boxes, got %d", len(d.Boxes))
	}
	if len(d.Connections) != 1 {
		t.Fatalf("expected 1 connection, got %d", len(d.Connections))
	}
	if !d.Connections[0].Arrow {
		t.Error("expected arrow connection")
	}
}

func TestParseBrokenMisalignedBoxes(t *testing.T) {
	// Slightly misaligned but still parseable
	input := `+----------+       +----------+
|  Frontend |----->|   API    |
+----------+       +----------+`

	d := Parse(input)
	if len(d.Boxes) < 2 {
		t.Fatalf("expected at least 2 boxes, got %d", len(d.Boxes))
	}
}

func TestParseThreeBoxes(t *testing.T) {
	input := `+---+     +---+
| A |---->| B |
+---+     +---+
            |
            v
          +---+
          | C |
          +---+`

	d := Parse(input)
	if len(d.Boxes) != 3 {
		t.Fatalf("expected 3 boxes, got %d", len(d.Boxes))
	}
	if len(d.Connections) < 2 {
		t.Fatalf("expected at least 2 connections, got %d", len(d.Connections))
	}
}

func TestParseBoxWithSplitTopEdge(t *testing.T) {
	input := `+-------------------+
| bx test diagram   |
+---------+---------+
          |
          v
+---------+---------+
| install verified  |
+-------------------+`

	d := Parse(input)
	if len(d.Boxes) != 2 {
		t.Fatalf("expected 2 boxes, got %d", len(d.Boxes))
	}
	if d.Boxes[0].Label != "bx test diagram" {
		t.Errorf("box 0 label: expected 'bx test diagram', got %q", d.Boxes[0].Label)
	}
	if d.Boxes[1].Label != "install verified" {
		t.Errorf("box 1 label: expected 'install verified', got %q", d.Boxes[1].Label)
	}
	if len(d.Connections) != 1 {
		t.Fatalf("expected 1 connection, got %d", len(d.Connections))
	}
	if !d.Connections[0].Arrow {
		t.Error("expected arrow connection")
	}
}
