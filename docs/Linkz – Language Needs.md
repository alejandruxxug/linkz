# 🦎 Linkz: Language Needs

Every time Linkz hits something lizard-script can't do (or does slowly), it goes here. When it's built, the change lands in the lizard-script repo and gets a row in its `docs/DESIGN.md` → *Extensions beyond the course design*.

## 🧱 Gaps

| Gap | Hit by | Workaround now | Priority |
|---|---|---|---|
| No `\"` or `\n` escapes in strings | `lib/svg.liz` | Single-quoted XML attributes; a literal line break inside a string | High (JSON and the MCP need `"`) |
| No floats | Curves, rotation, color math, 3D | Integer coordinates only | High |
| No dot access (`shapes.circle`) | Library API | Flat global names | High |
| No `use` / modules | Every library | The host loads `lib/` in a fixed order into one environment | High |
| No way to raise an error from lizard | `svg_shape` on an unknown type | Writes an SVG comment instead | Medium |
| No error handling (`try`) | The MCP server must survive bad requests | None yet | Medium (before the MCP) |
| Tree-walking speed | Rasterizer, tracing | Measured in B006 (table below): fine for one shape, slow for a long stroke redone every frame. Spans are cached per shape | Medium, see *Speed* |
| No string builder / bytes | PNG and JSON output | `+` on strings | Medium |

## ⏱ Speed

Measured with `go run . paint lib/gfx/tests/timing.liz out.png 512x512` (Apple silicon laptop, tree-walking evaluator, B006). The rasterizer is lizard (`lib/gfx/raster.liz`); `fb.spans` is Go.

| What | Time | Notes |
|---|---|---|
| `raster.fill` of a 512x512 circle (4 Beziers, 500 rows) | ~7 ms | |
| `raster.stroke` of a 200-point polyline, width 8 (2159 spans) | ~50 ms | quads plus round wedges, all edges scanned together |
| `raster.brush` of 40 points (312 samples, 1551 spans) | ~68 ms | about 8 samples a segment |
| `fb.spans` of those 3 span lists | ~0.6 ms | host side is not the cost |

So a typical shape costs a few ms and is rasterized once (`raster.shape` caches by path and widths); redoing a long stroke or brush on every mouse move would cost 50 to 70 ms a frame, so the editor should cache or cut samples while dragging. Faster sorting of crossings or less per-edge work in lizard would help; Go is not an option (rule 1).

## 📓 Log

| Date | Need | What happened |
|---|---|---|
| 2026-11-02 | Escapes, floats, dot, `use`, raise | Found while writing the first `svg.liz` / `shapes.liz` |
| 2026-11-02 | Speed | Rasterizer timings logged (see *Speed*); no new language need found. No `continue`, `+=` or `acos` was missed enough to matter. |
