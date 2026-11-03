# B008 – Drawing libraries

Milestone: M4 · Repo: `linkz` · Lane: lizard (`lib/draw/`) · Needs: lizard-script B003 (points, colors) · Parallel with B005/B006

## Goal
The drawing API a `.lzdraw` file uses, written in lizard with points and colors. Every shape becomes a **path** that the rasterizer (B006), the SVG writer and hit-testing all share.

## The shape record (the contract with B006/B007/B009 — keep it exactly)
```
{"kind": "circle", "args": [...the call's values...],
 "path": [["M", p], ["L", p], ["C", c1, c2, p], ["Z"]],     -~> absolute points, floats ok
 "fill": color or null, "stroke": color or null, "width": number,
 "widths": [numbers] or null, "opacities": [numbers] or null, -~> brush only, computed by the pen
 "anchors": [{"pos": p, "arg": i, "item": j or null, "role": "move"|"radius"|"corner"|"point"}],
 "line": n, "layer": ["sea", "waves"]}                       -~> set by draw() from src.here() when the host provides it; null otherwise
```
`anchors` say which literal a handle edits: argument index `arg`, and `item` for an element of an array argument (one point of a stroke).

## API (`lib/draw/shapes.liz`; positions are points, colors are COLOR)
| Call | Draws | Anchors |
|---|---|---|
| `circle(center, r, fill)` | circle (4 cubic Béziers) | center (move); center + (r, 0) (radius) |
| `ellipse(center, radii, fill)` | ellipse | center; center + radii (corner) |
| `rect(pos, size, fill)` | rectangle | pos (move); pos + size (corner) |
| `line(a, b, color, width)` | stroked segment | a; b |
| `polygon(points, fill)` | closed polygon | every point |
| `stroke(points, color, width)` | **smooth curve through the points** (Catmull–Rom → cubic Béziers, open ends, round caps) | every point |
| `brush(points, pressures, color, pen)` | smooth stroke shaped by a **pen style** from the raw pressures (0–1) | every point |
| `canvas` | `{"size": (512, 512), "background": #ffffff}`, a variable the drawing reassigns | — |
`stroke`/`brush` with 1 point → a dot; 2 → a straight segment. Wrong argument types → a lizard runtime error with a clear message (the language has no `error()` builtin yet; if you need one, report it as a language need instead of faking it).

## Doc comments (tooltips)
Every public function gets a doc comment right above it, in this exact format (the editor shows it as a tooltip; B005's `src.parse` reads it):
```
-~> Draws a filled circle.
-~> center: the middle of the circle, a point
-~> r: the radius in pixels
-~> fill: the color inside
-~> example: circle((256, 256), 40, #3fa34d)
function circle(center, r, fill) { ... }
```
Pens get one too (above each `pens.*` entry): what it feels like and which keys it uses.

## Pen styles (`lib/draw/pens.liz`)
A pen is a hash that turns pressure (and direction) into width and opacity per point; the stroke stores raw pressures, so changing the pen re-styles it. Defaults: `pens.ink` (pressure → width strongly, tapered ends), `pens.pencil` (thin; pressure → opacity more than width; slightly grainy edge from deterministic noise), `pens.marker` (constant width, semi-transparent, flat), `pens.brush` (wide width range, soft taper, fades at low pressure), `pens.calligraphy` (width from the angle between stroke direction and a nib angle). Shape: `{"name": "ink", "min": 1, "max": 10, "taper": 0.3, "opacity": [1, 1], "nib": null}` (document every key); users can write their own. `pens.list` for the editor's picker.

## Also
- `lib/draw/geometry.liz`: Catmull–Rom → Bézier, Bézier flattening (tolerance 0.25 px) → polylines, path bounds, point-in-path (nonzero), distance from a point to a polyline, `nearest_on_path(path, p)` (for inserting an anchor), and Ramer–Douglas–Peucker `simplify(points, tolerance)`. B006's rasterizer imports the flattening from here.
- `lib/draw/svg.liz`: rewrite with real `"` escapes and `\n`; a `<path d=…>` for every shape in **paint order** (bottom row first, so rows higher in the file end up in front; inside a layer the same rule); strokes with `stroke`/`stroke-width`, round caps/joins; `brush` as its filled outline.
- `host.Render` and `lib.Order` updated for `lib/draw/` (the orchestrator owns `lib/lib.go`; list the order you need in your report and edit it only if nobody else is).
- `examples/`: rewrite `hello.lzdraw` with the new API; add `layers.lzdraw` (two layers, nested objects, a loop; **only if lizard-script's `layer` statement exists yet — check `../lizard-script/docs/GRAMMAR.md`; otherwise skip and say so**) and `strokes.lzdraw`.

## Tests
`lib/draw/tests/*.liz` + `.expected`, run by a Go test like lizard-script's `tests/` (copy that harness style): geometry with known answers, simplify, point-in-path, the SVG of each example. `linkz render` works on every example.

## Done when
`go test ./...` green; the example SVGs look right (render to PNG with `qlmanage -t -s 512 -o <tmp> file.svg` and Read them).
