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
| A line that starts with `[` continues the line before it as an index: `let y0 = row * side` then `[x0, y0]` on the next line is read as `side[x0, y0]` (a syntax error) | `lib/gfx/tiles.liz`, `spatial.liz`, `lib/ui/view_scene.liz` (B035) | `return [...]` | Low |
| Garbage collection of boxed values: in the editor benchmark (B035, `host/perf_test.go`) Go's GC mark is ~45% of all CPU samples; a 10k-shape drawing keeps ~250 MB live (every number boxed, every span a 3-element array, every parsed row a hash). Per-op costs measured on the VM: an empty `for` over a list ~70 ns per turn, `raster.same_look` (9 field `===`) ~1.7 µs, one literal `circle(...)` row in `src.run` ~11 µs | Every big drawing: the run, the rows, the span cache | Linkz avoids the work (tiles, single-row re-run); the cost per object stays | High, see B033 (profiler, quickening) and packed number arrays |

| ~~`use` cycles are an error~~ | B041: the 13 files of `lib/ui` form one cycle | — | **Done** (LEP 0002 Amendment 2, pin f34d1ac): the port runs on the stock `lizardscript` |
| ~~`lizard.run` can't give the code live modules~~ | `linkz paint` programs, drawings | — | **Done** (`"use"` option): every drawing run (editor, render, paint, MCP) uses the drawing libraries alive, through `lib/app/drawing_bridge.liz` (B042) |
| ~~Each drawing run parses and compiles the drawing libraries again; `on_statement` fires for the libraries' statements too~~ | `src.run` fixed cost, calibrate | — | **Done** (module cache, `hook_source`): a small run is ~0.4 ms (Go host ~1.2 ms), the 10k fixture ~125 ms (Go host ~120 ms) (B042) |
| ~~Names through modules are slower than in one environment~~ | B041 perf table | — | **Done** (B043 VM overhaul in the pin): the parse phase that is left is the scanner, see the B042 rows |
| ~~No seconds since 1970 or UTC offset in `time.now()`~~ | Home screen, calibrate | — | **Done** (`time.now().epoch`, `.utc_offset`): `lib/app/clock.liz` removed (B042) |
| ~~No content hash~~ | Thumbnail names | — | **Done** (`fs.digest`): thumbnails are named by SHA-256 again, like the Go host; `lib/app/fingerprint.liz` removed (B042) |
| ~~`blit` has no size~~ | Home thumbnails | — | **Done** (`blit` `scale`): the nearest-pixel scaler in `host_fb.liz` removed (B042) |
| `lizard.run` gives only the first syntax error | A broken edit: the editor runs the whole text to confirm it is broken | `lizard.run` (~15 ms on the 10k fixture) then `lizard.check` for the list (~24 ms) | High: an `errors` list in the result (every syntax error) halves a broken run; `10k/broken_keystroke` is 36 ms vs 24 ms in the Go host |
| `lizard.tokens` makes a hash per token (~0.7 µs each) | The scanner (`lib/app/scan.liz`) | Linkz copies type/line/column into arrays once per scan | High: the 10k fixture has 210k tokens (150 to 240 ms); a 5 KB region (the editor's rescan after an edit) ~0.7 ms of a ~5 ms parse phase. A columnar form (`{"type": [...], "line": [...], "column": [...], "text": [...]}`) would cut both |
| VM speed for token loops | The scanner: the same algorithm as `host/scan.go` | None | High: a 5 KB region rescans in ~5 ms vs ~1.5 to 2.7 ms in the Go host's `parse` phase; it is the gap left in `small_edit`, `enter` and `paste_20` on 10k and `keystroke` on 2k (B042 perf table) |
| `on_statement` can't tell a function body's statements from the drawing's own | Shape tags (a shape drawn in a helper gets the line of the call) | `drawing.liz` parses the drawing (`lizard.parse`, 60 to 150 ms on 500 KB) to list the statements outside functions, only when the text has the word `function` | Medium: an option to leave function bodies out of `on_statement` (or a statements-only `lizard.parse`) |
| `lizardscript check` can't be told the host's names | `lib/` and its tests use what the host puts around them (`fb`, `src`, `settings`, the prelude's stand-ins) | None: 383 unknown-name and 35 member findings stay | Medium: a `check --names fb,src,...` flag or a project file (`lizard.check` already takes `names`); the other way is for the libraries to `use` host modules instead of names put around them |
| `check` says `unused-module` for a module whose variable is only assigned | `lib/app/drawing_bridge.liz` (`use "../draw/here"` then `here = ...`) | None | Low: a write through a module is a use of it |

## ⏱ Speed

Measured with `go run . paint lib/gfx/tests/timing.liz out.png 512x512` (Apple silicon laptop, tree-walking evaluator, B006). The rasterizer is lizard (`lib/gfx/raster.liz`); `fb.spans` is Go.

| What | Time | Notes |
|---|---|---|
| `raster.fill` of a 512x512 circle (4 Beziers, 500 rows) | ~7 ms | |
| `raster.stroke` of a 200-point polyline, width 8 (2159 spans) | ~50 ms | quads plus round wedges, all edges scanned together |
| `raster.brush` of 40 points (312 samples, 1551 spans) | ~68 ms | about 8 samples a segment |
| `fb.spans` of those 3 span lists | ~0.6 ms | host side is not the cost |

So a typical shape costs a few ms and is rasterized once (`raster.shape` caches by path and widths); redoing a long stroke or brush on every mouse move would cost 50 to 70 ms a frame, so the editor should cache or cut samples while dragging. Faster sorting of crossings or less per-edge work in lizard would help; Go is not an option (rule 1).

| Tokens don't know which file they came from | `src.run` tells library errors from user errors | The host pads library sources with 100,000 blank lines and treats any line ≥ 100,000 as library | Medium (a `source` name on tokens/errors removes the hack) |
| Optional/default parameters | `t()` + `tf()` in i18n | Two functions | B014 in progress |
| try/catch | Failed file writes kill the call | `fs.exists` guards | B014 in progress |

| ~~Parser error recovery cascades~~ | Editor error rows (QA BUG-6) | — | **Done** (recovery by line in the pin): the editor's "may be caused by the error above" marking is gone, every error is its own (B042) |

## 📓 Log

| Date | Need | What happened |
|---|---|---|
| 2026-11-02 | Escapes, floats, dot, `use`, raise | Found while writing the first `svg.liz` / `shapes.liz` |
| 2026-11-02 | Speed | Rasterizer timings logged (see *Speed*); no new language need found. No `continue`, `+=` or `acos` was missed enough to matter. |
| 2026-11-03 | `use` cycles, live modules in `lizard.run`, run and lookup speed, unix time, hash, scaled blit | Found porting Linkz to lizard (B041): see the rows above |
| 2026-11-03 | Pin f34d1ac: cycles, live `use`, module cache, `hook_source`, epoch, digest, scaled blit, recovery | B042: all used, the workarounds removed. New: first syntax error only, token hashes, scanner speed, function bodies in `on_statement`, check without host names |
