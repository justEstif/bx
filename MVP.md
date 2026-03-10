# MVP: bx — Terminal Diagram Tool

## What

A CLI tool that takes broken/rough ASCII diagrams and re-renders them with **geometrically perfect alignment** using Unicode box-drawing characters.

**Core features:**

- Parse rough ASCII/Unicode diagrams (boxes, arrows, labels)
- Understand the _intent_ — which boxes exist, what connects to what, rough spatial layout
- Re-render with mathematically correct alignment on a character grid
- Output to stdout (terminal) or markdown codeblock

**What it is NOT:**

- Not a PNG/SVG renderer — output is always monospace text
- Not a full diagramming IDE — no interactive editing
- Not trying to replace Mermaid/D2 — those are input formats, this is an output renderer
- Not a TUI app — it's a pipe-friendly CLI tool

## Why

**The problem:** Every AI agent, every developer, and even careful humans produce ASCII diagrams with misaligned lines. Vertical lines drift, box corners don't match, arrows don't land on connection points. This happens because text editors and LLMs have no concept of spatial geometry — they place characters sequentially with no grid constraint.

**Who it's for:**

- Developers who use AI coding agents and want clean diagrams in terminal/markdown
- Anyone who documents with ASCII diagrams and is tired of manual alignment fixes

**Why existing solutions fall short:**

| Tool                     | Problem                                                         |
| ------------------------ | --------------------------------------------------------------- |
| `mermaid-ascii`          | Renders from Mermaid DSL only, can't fix existing ASCII         |
| `markdown-diagram-fixer` | Patches specific issues, doesn't re-render from a layout engine |
| `flow-wire-diagram`      | Repairs alignment defects but limited scope                     |
| `graph-easy`             | Old Perl tool, ugly output, DOT input only                      |
| Raw LLM output           | The cause of the problem, not the solution                      |

The gap: **no tool treats terminal diagrams as a layout problem.** They all work at the text level. This tool works at the geometry level — parse structure, compute layout, render perfectly.

## How

### Architecture

```
INPUT (broken ASCII/text)
  → Parser (extract boxes, labels, connections, rough positions)
  → Layout Engine (snap to grid, compute exact positions, route connections)
  → Renderer (place box-drawing characters on character grid)
  → OUTPUT (perfect Unicode text)
```

### Key Technical Challenges

1. **Parsing rough ASCII** — Based on ditaa's proven approach:
   - **Structural characters:** `+`, `-`, `|`, `/`, `\` (ASCII) and `┌─│└` etc. (Unicode) are drawing characters — everything else is text
   - **Box detection:** Flood-fill connected regions of structural characters to find closed rectangular shapes
   - **Labels:** Text inside a closed box region = that box's label
   - **Connections:** `-` and `|` sequences that touch a box on at least one end = connections. `>`, `v`, `^`, `<` at endpoints = arrow direction
   - **Ambiguity rule:** If a structural sequence doesn't connect to any box, leave it untouched. Don't mangle what you can't confidently parse.
   - **Test data:** The [CHI 2024 ASCII diagram dataset](https://asciidiagrams.github.io/) has 2,156 real diagrams from Chromium, Linux, LLVM, TensorFlow — use these as test cases.
   - **Test oracle:** [terrylica/ascii-diagram-validator](https://github.com/terrylica/cc-skills/tree/main/plugins/doc-tools/skills/ascii-diagram-validator) is a linter that detects alignment errors in box-drawing diagrams. Use it to validate `bx` output — after `bx fix`, the validator should report zero errors.

2. **Layout engine** — Given N boxes with approximate positions and connections between them, compute exact grid positions that:
   - Preserve the user's intended spatial arrangement (don't rearrange everything)
   - Ensure all boxes are properly sized for their content
   - Route connections without overlapping boxes
   - Keep everything aligned to a character grid

3. **Rendering** — Map the layout to actual characters. Pick correct box-drawing glyphs for corners, intersections, and line segments. Handle arrows (▶ ▼ ◀ ▲).

### Tech Stack

- **Language:** Go (single binary, fast, no runtime deps)
- **No external dependencies** for core rendering — this is pure grid math
- **Testing:** Golden file tests — input broken ASCII, expected output perfect ASCII

### Data Model

```
Diagram
  ├── Box[]
  │     ├── id: string
  │     ├── label: string
  │     ├── position: {row, col}  // grid position (computed)
  │     ├── size: {width, height} // computed from label + padding
  │     └── style: single | double | rounded | bold
  │
  ├── Connection[]
  │     ├── from: box_id
  │     ├── to: box_id
  │     ├── direction: left | right | up | down
  │     ├── label?: string
  │     └── style: solid | dashed
  │
  └── Grid (character matrix for final render)
```

### Usage

```bash
# Fix a broken diagram from stdin
cat rough-diagram.txt | bx fix

# Fix diagrams inside a markdown file (rewrites fenced codeblocks)
bx fix --markdown README.md

# Pipe from an agent
echo "draw me an architecture diagram" | pi | bx fix


```

## Competition

| Tool                     | Input              | Output              | Layout Engine? | Fixes broken ASCII? |
| ------------------------ | ------------------ | ------------------- | -------------- | ------------------- |
| `mermaid-ascii`          | Mermaid DSL        | ASCII               | Yes            | No                  |
| `mermaid-ascii-diagrams` | Mermaid in MD      | Unicode             | Yes            | No                  |
| `beautiful-mermaid`      | Mermaid DSL        | ASCII/SVG           | Yes            | No                  |
| `markdown-diagram-fixer` | Broken ASCII       | Fixed ASCII         | No (patches)   | Sort of             |
| `flow-wire-diagram`      | Broken ASCII in MD | Fixed ASCII         | No (repairs)   | Sort of             |
| `graph-easy`             | DOT                | ASCII               | Yes            | No                  |
| `boxes`                  | Text               | Boxed text          | No             | No                  |
| **This tool**            | **Broken ASCII**   | **Perfect Unicode** | **Yes**        | **Yes**             |

**The differentiator:** This is the only tool that combines _parsing broken input_ with a _real layout engine_ that outputs _geometrically correct Unicode box-drawing diagrams_.

## MVP Milestones

1. **Box renderer** — Given a list of boxes with labels and explicit positions, render them perfectly on a character grid with Unicode box-drawing characters. No parsing, no connections. Just prove the grid math works.

2. **Connection routing** — Add arrows/lines between boxes. Handle horizontal, vertical, and L-shaped connections. Correct glyphs for corners and intersections.

3. **ASCII parser** — Parse rough/broken ASCII input to extract boxes and connections. Detect `+--+`, `|`, `---`, `-->` patterns. This is the hard milestone.

4. **Layout snapping** — Take parsed boxes with approximate positions and snap them to a clean grid. Preserve spatial intent (what's left of what, what's above what) but fix alignment.

5. **CLI interface** — stdin/stdout pipe, `diag fix` command, `--markdown` flag for processing fenced codeblocks in markdown files.

6. **Styles** — Support rounded (`╭╮`), double (`╔╗`), bold (`┏┓`) box styles. Auto-detect from input or allow override via flag.

## Open Questions

- [x] Go or Rust? → **Go**
- [x] Should it preserve the user's spatial arrangement exactly, or allow minor repositioning for cleaner output? → **Keep the user's layout, just fix alignment**
- [x] How to handle ambiguous input? → **Follow ditaa's model: `+`, `-`, `|` are structural characters. Text inside closed regions = labels. `---` is a connection only if it touches a box on at least one end. When in doubt, leave it alone — don't mangle what you can't parse.**
- [x] Labeled connections (text on arrows)? → **Post-MVP**
- [x] Mermaid input support? → **No. ASCII/Unicode input only.**
- [x] What's the tool name? → **`bx`**
