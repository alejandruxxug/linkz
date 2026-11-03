# B006 – Rasterizer and text

Milestone: M3 · Repo: `linkz` · Lane: lizard (`lib/gfx/`) · Needs: B005 (`fb` primitives, `linkz paint`), B006a (font data), B008's `geometry.liz` (flattening)

## Goal
Lizard turns paths into pixels and draws text, using only the host's `fb.span` / `fb.spans` / `fb.fill`.

## `lib/gfx/raster.liz`
- `raster.fill(path, rule)` → span list `[[y, x0, x1], ...]` (x0/x1 FLOAT for soft horizontal edges); `rule` "nonzero" or "evenodd". Scanline with an active edge table over the flattened path.
- `raster.stroke(path, width)` → spans of the stroke outline: offset both sides along normals, round joins and caps (short arc segments), filled nonzero.
- `raster.brush(path_points, widths)` → spans of a variable-width smooth stroke (half-width interpolated along the curve between points).
- `raster.shape(shape)` → spans for a B008 shape record (fill and/or stroke/brush), **cached** by the shape's path + width(s) so an unchanged shape is never re-rasterized.
- `raster.draw(spans, offset, color)` → `fb.spans`.
- Vertical anti-aliasing optional for v1; horizontal soft ends required. No gaps between adjacent shapes, no off-by-one on the last row, shapes partly outside the buffer clip cleanly.

## `lib/gfx/font.liz`
- Uses `lib/gfx/font_data.liz` (B006a). `font.draw(text, pos, color, scale)` (integer scale 1–3, nearest), `font.measure(text, scale)`, `font.height(scale)`. Glyph span lists are pre-scaled and cached per scale; one `fb.spans` call per glyph. Non-ASCII → a box glyph.

## Tests
- `linkz paint` programs in `lib/gfx/tests/` drawing known scenes (filled rect, circle, a star evenodd vs nonzero, a stroked zigzag, a brush, text at scale 1 and 2) to PNG; a Go test checks spot pixels with independently computed expectations (inside filled, just outside empty, edge coverage partial, glyph pixels set).
- Timing: ms to rasterize a 200-point stroke and a 512×512 circle; add them to the speed table in `docs/Linkz – Language Needs.md`.

## Done when
`go test ./...` green; the test PNGs look right (Read them); timings logged.
