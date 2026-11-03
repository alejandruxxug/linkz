# B031 – `gfx`: 2D graphics in the Terrarium

Milestone: std · Repo: `../lizard-script` · Lane: lizard (+ the Go span primitives) · Needs: B030 (window, img), B029 (speed)

> "I would like a graphics library and a 3d graphics library. Since we are doubling down on this" — the developer, 2026-11-03

## Goal
Every lizard program can draw anti-aliased 2D graphics into a `window` or an `img`:
`gfx.circle(target, (300, 200), 50, #f2a65a)`, paths, Béziers, strokes, pressure brushes, gradients, blend modes, text, SVG export.

## Source
Move (copy, then generalize) Linkz's `lib/draw/*.liz` (shapes, geometry, pens, gradient, blend, svg) and `lib/gfx/*.liz` (raster, font, calibrate's cost model) into `terrarium/gfx/*.liz` (embedded). The per-pixel Go primitives (`span`, `spans`, `span_gradient`, blend compositing) become functions on Terrarium targets. No Linkz names, no editor concepts (line tags, layers-as-rows) in the library: those stay in Linkz.

## Clean code (the developer's ask)
One concept per file, small functions, every public function documented with `-~>` (signature, what it does, an example), no dead code, no editor-specific hacks. Someone new should be able to read `terrarium/gfx/` top to bottom.

## Done when
Tests (lizard tests with expected output + Go tests for the primitives), `examples/terrarium/shapes.liz` renders a PNG identical to Linkz's render of the same shapes, docs section in `docs/TERRARIUM.md` and a chapter in `docs/learn/` (EN + ES). Linkz migrates to it in a later batch.
