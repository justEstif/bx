package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/justEstif/bx/grid"
	"github.com/justEstif/bx/layout"
	"github.com/justEstif/bx/parse"
	"github.com/justEstif/bx/render"
)

func TestGoldenFiles(t *testing.T) {
	inputs, err := filepath.Glob("testdata/*.input.txt")
	if err != nil {
		t.Fatal(err)
	}
	for _, inputPath := range inputs {
		name := strings.TrimSuffix(filepath.Base(inputPath), ".input.txt")
		goldenPath := strings.Replace(inputPath, ".input.txt", ".golden.txt", 1)

		t.Run(name, func(t *testing.T) {
			inputData, err := os.ReadFile(inputPath)
			if err != nil {
				t.Fatal(err)
			}
			goldenData, err := os.ReadFile(goldenPath)
			if err != nil {
				t.Fatalf("missing golden file %s: %v", goldenPath, err)
			}

			got := fix(string(inputData), render.Single)
			want := string(goldenData)
			if got != want {
				t.Errorf("mismatch for %s\ngot:\n%s\nwant:\n%s", name, got, want)
			}
		})
	}
}

func fix(input string, style render.Style) string {
	d := parse.Parse(input)
	if len(d.Boxes) == 0 {
		return input
	}
	result := layout.Snap(d, style)
	g := grid.New(result.Width, result.Height)
	for i := range result.Connections {
		result.Connections[i].Draw(g)
	}
	for i := range result.Boxes {
		result.Boxes[i].Draw(g)
	}
	return g.String()
}
