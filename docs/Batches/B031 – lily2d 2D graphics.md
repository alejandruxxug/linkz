# B031 – `lily2d`: 2D graphics for every lizard program

Repo: `../lizard-script` · Lane: lizard (+ Go only for per-pixel work) · Needs: Terrarium surfaces (B040, done)

> "can you build the 2d and 3d libraries? name them lily2d and lily3d" — the developer, 2026-11-03

## Goal
`use "lily2d"` gives any lizard program anti-aliased 2D vector graphics that draw into a window, a surface or an image:

```
use "lily2d"
let w = window.open(600, 400, "lily")
lily2d.circle(w, (300, 200), 50, #f2a65a)
lily2d.path(w, [(10, 10), (200, 40), (120, 180)], {"fill": lily2d.linear([(0, #f2a65a), (1, #3b2a5c)], 90deg), "stroke": #1a1418, "width": 3})
```

## Standard modules (new mechanism, this batch owns it)
Big libraries are not installed into every program (startup cost); they're **standard modules** found by `use` after the program's own folder and `LIZARD_PATH`: embedded `.liz` files in `terrarium/std/` (`terrarium/std/lily2d.liz`, and sub-files in `terrarium/std/lily2d/`). Add the lookup in `evaluator/modules.go` (+ tests), document it in proposals/0002 (amendment), TERRARIUM.md, guide chapter 18. lily3d (B032, running in parallel) puts its files in `terrarium/std/lily3d*` and relies on this.

## Contents (port and generalize Linkz's `../linkz/lib/draw/*.liz` and `../linkz/lib/gfx/*.liz`, read-only)
- Shapes: `circle`, `ellipse`, `rect` (rounded corners), `line`, `polyline`, `polygon`, `path` (points with Catmull-Rom smoothing and corners, or explicit cubic/quadratic Bézier commands), `arc`, `star`, `text` (the bitmap font, with alignment).
- Paint: solid colors, `linear`/`radial` gradients (stops, steps, hard stops, shape-relative and absolute), opacity; drawing through the Terrarium's gradient spans; blend layers through `push_layer`/`pop_layer`.
- Strokes: width, joins, caps, dashes; `brush(points, pressures, color, pen)` with the pens (ink, pencil, marker, brush, calligraphy) and the pressure curve block.
- Geometry: Bézier evaluation and flattening, bounds, transforms (`translate/rotate/scale` as a 2D matrix, applied to shapes), hit-testing (point in shape), simplification.
- Rasterizer: scanline fill with anti-aliasing (the `lib/gfx/raster.liz` algorithm) producing spans, a cache keyed by shape, and `lily2d.svg(shapes)` export.
- A retained mode: `lily2d.scene()` with `add/remove/draw(target)` that only repaints changed regions (the Linkz tiles + spatial index idea, generalized), so apps get it for free.

## Clean code
One concept per file under `terrarium/std/lily2d/`, every public function documented with `-~>` (signature, what it does, example). No Linkz/editor concepts (no line tags, rows, layers-as-code). LEP 0001 style.

## Done when
Lizard tests (`lizardscript test`, golden PNGs) on both engines, matching Linkz's renders of the same shapes pixel for pixel (or within 1 on soft alpha, explained); `examples/lily2d/` (shapes gallery, a gradient poster, a pressure doodle pad, an animated scene); docs `docs/LILY2D.md` (every function); guide chapter 19 "lily2d" (en + es; edit the guide index only for your row); benchmarks in `bench/` if one compares fairly. Don't commit.
