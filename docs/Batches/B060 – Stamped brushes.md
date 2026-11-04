# B060 – Stamped brushes (engine, headless)

Model: Sonnet. Rules: `_v11-rules.md`. Read only: `docs/Linkz – Brushes.md` (the parameter table and "Stamped brushes"), `lib/draw/pens.liz`, the `brush()` part of `lib/draw/shapes.liz`, `lib/draw/svg.liz`, and the raster entry points you need in `lib/gfx/raster.liz`.

> "Brush creator would allow to parametrize brushes Procreate style." — the developer

Grow a pen into a Procreate-style brush, step 1 of the Brushes doc's order. No UI in this batch.

- A pen hash may say `"render": "stamped"` (default `"vector"`, today's behavior, unchanged). A stamped stroke is many stamps placed along the stroke.
- New optional keys, documented in the `pens.liz` header like the existing ones: `spacing` (share of the width between stamps), `jitter` (position), `shape` (procedural: `"round"`, `"square"`, `"flat"`, `"star"`; leave a hook for a bitmap later), `roundness`, `rotation` (`"fixed"`, `"follow"`, `"random"` plus an angle), `scatter`, `count` (stamps per step), `grain` as a hash (`{"kind": "noise", "scale", "depth", "moving"}`); a plain number keeps today's meaning, `size_jitter`/`opacity_jitter`, `color_jitter` (`[h, s, b]`), and `speed` dynamics for size and opacity. Taper and the `pressure` block apply to stamps too.
- Everything is deterministic: random values come from a seed made from the stroke's points, so a file always renders the same.
- Renders in the canvas, the PNG export and SVG. For SVG, the stamps become a `<g>` of shapes; keep it readable.
- Speed matters. Cache the stamps per stroke, and measure a 300-point stamped stroke at 2k shapes in `tests/editor/editor_bench.liz`. If it's slow, report which language feature would fix it; don't hack around it.
- Add 3–4 example stamped pens (e.g. `pens.chalk`, `pens.spray`, `pens.ink_dots`, `pens.star_trail`) with en/es names.
- Tests: `lib/draw/tests/` for spacing, determinism, taper, each shape and rotation, SVG output, plus a PNG golden. Both engines.

Report: files changed, any shared edits, bench numbers before and after, and the language needs you found.
