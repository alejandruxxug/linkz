# B067 – More stamp shapes and grains

Model: Sonnet. Rules: `_v11-rules.md`. Read only: `lib/draw/stamps.liz`, `stamp_color.liz`, `svg_stamps.liz`, the `pens.liz` header, `lib/draw/tests/stamps.liz`, and the shape/grain pickers in R005 frames 1 and 3 (`docs/Art/Reviews/R005/R005-01-…png`, `R005-03-…grain.png`) for what they show.

- New shapes: `"spray"` (a cloud of tiny dots per stamp), `"leaf"`, and **bitmap** shapes: `"shape": {"bitmap": <run-length data>, "w", "h"}` (the stamp is the bitmap's alpha, scaled to the stamp size).
- New grains: `"kind"` `"paper"`, `"canvas"`, `"chalk"`, `"hatch"` (procedural, deterministic) and `{"kind": "bitmap", "data", "w", "h"}`. Grain stays per stamp as B060 decided, unless a cheap per-pixel way exists with today's `fb` (if not, write the language need).
- One function each that lists what exists, with en/es names: `stamp_shapes()` and `stamp_grains()` (the creator, B066, builds its pickers from them).
- `img`-to-bitmap helper for the creator's "import image": `stamp_bitmap_from_image(path, size)` (Terrarium `img` to load and shrink; alpha from luminance when the image has none).
- Tests + a PNG golden per new shape/grain, both engines; SVG output for each (bitmaps as `<image>` or as paths — say which and why). The bench for a 300-point stroke for each new shape.
- Stay in `lib/draw/` (and its tests); others are in `lib/ui/pen_*` and `lib/app/brush_files.liz`.
