package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/estifanos/bx/grid"
	"github.com/estifanos/bx/layout"
	"github.com/estifanos/bx/parse"
	"github.com/estifanos/bx/render"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "fix":
		runFix()
	case "help", "--help", "-h":
		printUsage()
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n", os.Args[1])
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Fprintln(os.Stderr, `bx — fix broken ASCII diagrams

Usage:
  bx fix [OPTIONS]              Read from stdin, output fixed diagram
  bx fix [OPTIONS] <file>       Read from file, output fixed diagram

Options:
  --style <style>   Box style: single (default), rounded, double, bold
  --markdown        Process fenced codeblocks in a markdown file

Examples:
  cat diagram.txt | bx fix
  bx fix --style rounded diagram.txt
  bx fix --markdown README.md`)
}

func runFix() {
	style := render.Single
	var inputFile string
	markdown := false

	args := os.Args[2:]
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--style":
			if i+1 < len(args) {
				i++
				style = parseStyle(args[i])
			}
		case "--markdown":
			markdown = true
		default:
			inputFile = args[i]
		}
	}

	var input string
	if inputFile != "" {
		data, err := os.ReadFile(inputFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error reading %s: %v\n", inputFile, err)
			os.Exit(1)
		}
		input = string(data)
	} else {
		data, err := io.ReadAll(os.Stdin)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error reading stdin: %v\n", err)
			os.Exit(1)
		}
		input = string(data)
	}

	if markdown {
		fmt.Print(processMarkdown(input, style))
	} else {
		fmt.Print(fixDiagram(input, style))
	}
}

func fixDiagram(input string, style render.Style) string {
	d := parse.Parse(input)
	if len(d.Boxes) == 0 {
		return input // nothing to fix, return as-is
	}

	result := layout.Snap(d, style)
	g := grid.New(result.Width, result.Height)

	// Draw connections first, then boxes on top (boxes overwrite connection overlap)
	for i := range result.Connections {
		result.Connections[i].Draw(g)
	}
	for i := range result.Boxes {
		result.Boxes[i].Draw(g)
	}

	return g.String()
}

func processMarkdown(input string, style render.Style) string {
	lines := strings.Split(input, "\n")
	var out strings.Builder
	inBlock := false
	var blockLines []string

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if !inBlock && (trimmed == "```" || strings.HasPrefix(trimmed, "```")) {
			// Check if this looks like a diagram codeblock
			inBlock = true
			blockLines = nil
			out.WriteString(line + "\n")
			continue
		}
		if inBlock && trimmed == "```" {
			// End of block — process it
			block := strings.Join(blockLines, "\n")
			fixed := fixDiagram(block, style)
			out.WriteString(fixed)
			if !strings.HasSuffix(fixed, "\n") {
				out.WriteString("\n")
			}
			out.WriteString(line + "\n")
			inBlock = false
			continue
		}
		if inBlock {
			blockLines = append(blockLines, line)
		} else {
			out.WriteString(line + "\n")
		}
	}

	return out.String()
}

func parseStyle(s string) render.Style {
	switch strings.ToLower(s) {
	case "rounded":
		return render.Rounded
	case "double":
		return render.Double
	case "bold":
		return render.Bold
	default:
		return render.Single
	}
}
