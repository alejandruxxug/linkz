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

| `use` cycles are an error, but the editor's files call each other (13 files of `lib/ui` form one cycle: input calls painting, painting calls layout, layout reads the app). The Go host loaded them into one environment | B041: `lizardscript run . edit`, the `tests/ui` and `tests/editor` suites | None in the port: it runs on a lizardscript with the one-line change below (a `use` of a module that is still loading gets its live, partly loaded scope; `use ... as` of one stays an error). Measured: the editor loads in ~40 ms and works, pixel for pixel as the Go host | **Blocking for the editor**; changes what a cyclic program does (the cycle tests of `evaluator/modules_test.go`), so it **needs the team's agreement** (LEP 0002 amendment) |
| `lizard.run` can't give the code live modules: `globals` are copies, so a program that changes a library's variable (`raster_cache_limit = 4`, `here = ...`) changes only its copy | B041: `linkz paint` programs, drawings, the test runners | Tests that change a library use it themselves (`use` lines in `lib/*/tests`); a paint program can't | High: a `"use": ["draw/shapes", ...]` option (the code sees those modules alive, loaded once per run) would also make each drawing run cheaper |
| Each drawing run parses and compiles the 8 drawing libraries again (~5 ms) and its `on_statement` hook is called for every statement, the libraries' too (~0.45 µs each, +240 ms on 10k shapes) | B041: `src.run` (6 ms fixed per run vs ~1 ms in Go), `run_ms_per_object` 3x in calibrate | `drawing.liz` filters by source | High: a parse cache for modules, and an `on_statement` for the run's own source only |
| Names through modules are slower than in one environment: the editor's parse phase is ~20 ms per edit in the port vs 2 to 5 ms in the Go host, the same code on the same VM | B041 perf table (`tests/editor/editor_bench.liz`) | None | High (B043 VM overhaul) |
| No seconds since 1970 or UTC offset in `time.now()` | Home screen "2 h ago", calibrate's `measured_at` | `lib/app/clock.liz`: measures the offset once from a file's mtime | Medium (`time.now().unix`, `.offset`) |
| No content hash (`fs.hash`, SHA-256) | Thumbnail names | `lib/app/fingerprint.liz`, two polynomial hashes in lizard | Low |
| `blit` has no size (scaled copy) | Home thumbnails | `host_fb.liz` scales once with the nearest pixel, kept | Low |

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

| Parser error recovery cascades (a bad token on one line yields a fake error on a later valid line) | Editor error rows (QA BUG-6) | Editor marks later parse errors as "may be caused by the error above" | Medium — the real fix (recover at the next statement on a new line) changes the team's `syntax_several_errors.expected`, so it **needs the team's agreement** |

## 📓 Log

| Date | Need | What happened |
|---|---|---|
| 2026-11-02 | Escapes, floats, dot, `use`, raise | Found while writing the first `svg.liz` / `shapes.liz` |
| 2026-11-02 | Speed | Rasterizer timings logged (see *Speed*); no new language need found. No `continue`, `+=` or `acos` was missed enough to matter. |
| 2026-11-03 | `use` cycles, live modules in `lizard.run`, run and lookup speed, unix time, hash, scaled blit | Found porting Linkz to lizard (B041): see the rows above |
