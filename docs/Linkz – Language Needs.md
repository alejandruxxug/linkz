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
| ~~`lizard.run` gives only the first syntax error~~ | A broken edit | — | **Done** (`r.errors`, a2c3768): a broken run parses once; `lib/app/syntax.liz` removed; 10k/broken_keystroke ~20 ms (Go host ~22 ms) |
| ~~`lizard.tokens` makes a hash per token~~ | The scanner | — | **Done** (`{"columns": true}`, a2c3768): `lib/app/scan.liz` reads the columns, the same output on 96 files |
| ~~VM speed for token loops~~ | The scanner | — | **Done** (index caches, `push`, string `==` in the VM, a2c3768): the parse phase of an edit is ~3.5 ms on 10k; every editor benchmark but the first `small_edit` round is within 10% of the Go host (B042) |
| ~~`on_statement` can't tell a function body's statements from the drawing's own~~ | Shape tags | — | **Done** (`hook_functions: false`, a2c3768): no parse before a run |
| ~~`lizardscript check` can't be told the host's names~~ | `lib/` and its tests | — | **Done** (`lizard-check.liz`, a2c3768): `lib/lizard-check.liz`, `lib/ui/tests/` and `lib/gfx/tests/` list theirs; Linkz's own additions moved out of the Terrarium's hashes into `linkz` (host_linkz.liz). `lizardscript check main.liz lib tests tools`: no findings |
| ~~`check` says `unused-module` for a module whose variable is only assigned~~ | `drawing_bridge.liz` | — | **Done** (a2c3768) |

## ⏱ Speed

Measured with `go run . paint lib/gfx/tests/timing.liz out.png 512x512` (Apple silicon laptop, tree-walking evaluator, B006). The rasterizer is lizard (`lib/gfx/raster.liz`); `fb.spans` is Go.

| What | Time | Notes |
|---|---|---|
| `raster.fill` of a 512x512 circle (4 Beziers, 500 rows) | ~7 ms | |
| `raster.stroke` of a 200-point polyline, width 8 (2159 spans) | ~50 ms | quads plus round wedges, all edges scanned together |
| `raster.brush` of 40 points (312 samples, 1551 spans) | ~68 ms | about 8 samples a segment |
| `fb.spans` of those 3 span lists | ~0.6 ms | host side is not the cost |

So a typical shape costs a few ms and is rasterized once (`raster.shape` caches by path and widths); redoing a long stroke or brush on every mouse move would cost 50 to 70 ms a frame, so the editor should cache or cut samples while dragging. Faster sorting of crossings or less per-edge work in lizard would help; Go is not an option (rule 1).

**Tracing speed (v3 step 4).** The trace now fans its independent work out to `task`s (see *Tracing v3*), and three things in the language would make that simpler or let it reach the editor:

| Need | Why | Today |
|---|---|---|
| A task can start tasks | The Trace image dialog runs the whole trace in one `task` (so the window keeps drawing), and inside a task `task` is the error `not available inside a task`, so the trace's own fan-out is off in the editor | The editor traces serially (about 40 percent slower than `linkz trace` on a big picture) |
| `os.cpus()` (the processor count) | The number of pieces to cut the work in, and the choice to stay serial on one processor, are guesses | A fixed 12 pieces; tasks beyond the free processors wait their turn |
| A list of n copies: `list.fill(n, value)` | Every grid of the tracer is built by n calls of `push`, a good part of its time | `push` in a loop |

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
| 2026-11-03 | Pin a2c3768: `r.errors`, column tokens, VM token loops, `hook_functions`, check project files | B042: all used; the Go host deleted (M11 done) |
| 2026-11-04 | Tracing speed: tasks inside tasks, `os.cpus()`, `list.fill` | v3 step 4: the trace fans out to tasks outside a task only; see the table under *Speed* |
