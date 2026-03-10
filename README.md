# bx

Fix broken ASCII diagrams → perfect Unicode box-drawing.

```
cat rough.txt | bx fix
```

## Before → After

```
+----------+       +----------+        ┌──────────┐    ┌─────┐
|  Frontend |----->|   API    |        │ Frontend │───▶│ API │
+----------+       +----------+        └──────────┘    └─────┘
                       |                                  │
                       v                                  ▼
                  +---------+                         ┌──────┐
                  |  Auth   |                         │ Auth │
                  +---------+                         └──────┘
```

## Install

### Go

```bash
go install github.com/justEstif/bx@latest
```

### mise

```bash
mise use -g go:github.com/justEstif/bx@latest
mise reshim
```

### From source

```bash
git clone https://github.com/justEstif/bx.git
cd bx && go build -o bx .
```

## Usage

```bash
# Fix from stdin
cat rough.txt | bx fix

# Fix a file
bx fix diagram.txt

# Choose style: single (default), rounded, double, bold
echo '...' | bx fix --style rounded

# Fix diagrams inside markdown code blocks
bx fix --markdown README.md
```

## Styles

```
single (default)    rounded           double            bold
┌──────┐            ╭──────╮          ╔══════╗          ┏━━━━━━┓
│ Auth │            │ Auth │          ║ Auth ║          ┃ Auth ┃
└──────┘            ╰──────╯          ╚══════╝          ┗━━━━━━┛
```

## How it works

1. **Parse** — detect boxes, connections, and labels from rough ASCII/Unicode input
2. **Snap** — align to a clean grid, preserving spatial layout
3. **Render** — redraw with perfect Unicode box-drawing characters

Tolerates misaligned corners, overflowing labels, broken edges, and mixed ASCII/Unicode input.

## License

MIT
