# 🗓️ Linkz: Build Plan

## ⏳ Waiting on you

| What | Why it's blocking | Where to look | Since |
|---|---|---|---|
| OK to commit the extensions-log section + CLAUDE.md in lizard-script | Language work starts there | lizard-script `git status` | 2026-11-02 |

## 🌙 Overnight plan (2026-11-03) — done, see [[2026-11-03 Overnight]]
Order: finish B007b (editor view) + B015 (home) → **B007c autocomplete (mandatory)** → **B009** pencil/brush/anchors (v1 finish line) → v1.1 **engine** parts that need no UI approval (gradients + steps and gradient spans, blend-mode compositing, MAX OBJ calibration) → M7 MCP (agents draw by writing lizard) → M9 release *prep* only (GoReleaser config + install script; nothing published). v1.1 **UI** (R002) waits for the developer's approval in Figma. Report: `docs/Reports/2026-11-03 Overnight.md`.

## 🧱 Milestones

| # | Milestone | Done when |
|---|---|---|
| M0 ✅ | Foundation | `linkz render examples/hello.lzdraw` writes an SVG made entirely by lizard libraries; `go test ./...` green |
| M1 | Language basics | lizard-script has string escapes, dot access and floats, each with tests and an extensions-log row; `svg.liz` uses `"` and `\n` |
| M2 | Modules | `use "shapes"` works; `lib.Order` is gone; libraries call each other by name |
| M3 | Editor | `linkz serve` opens a page: code on the left, live SVG on the right, errors on their line; open/save `.lzdraw` |
| M4 | Math + more shapes | `math` (sin, cos, sqrt, random), `path`, `polygon`, `group`, `color`, `transform` |
| M5 | Raster export | A rasterizer in lizard; PNG export from the editor and `linkz render --png`; first speed measurement logged |
| M6 | Speed | Bytecode compiler + VM in lizard-script; rasterizer re-measured. Brief: [[B029 – Bytecode VM]] (started 2026-11-03) |
| M7 | Agents | `linkz mcp`: run / render (PNG) / docs / save, protocol written in lizard |
| M8 | Tracing | `trace.image(...)` turns a PNG into editable shapes |
| M9 | Release | GoReleaser binaries for macOS/Windows/Linux + install script; the reveal |

## 🧱 v1.1 (after the v1 finish line)
| # | Batch | Done when |
|---|---|---|
| B016 | Right-click menu + Fill editor + gradients (`linear`/`radial` values in lib/draw, gradient spans in the host), **stepped gradients** (`steps = n`) and hard stops | Approved R002 frames; a gradient made in the editor renders on canvas, in SVG and in PNG |
| B017 | Palette sidebar (HSV picker, recent, document palette as `let colors = {...}`) **+ the color popover** next to any clicked color (variables → writes `colors.sky` as a reference, recent colors, compact HSV + hex, save as variable; one undo step) | Changing a named color updates every shape using it |
| B018 | Artist mode (layers with thumbnails, code drawer) | Switching modes keeps the same file and selection |
| B020 | Layer blend modes (Normal, Multiply, Screen, Overlay, Darken, Lighten, Add) + opacity: `blend("multiply")` / `opacity(0.8)` as a layer's first lines; per-layer buffers blended by the host; SVG `mix-blend-mode` | A Multiply layer darkens what's below in the canvas, PNG and SVG |
| B021 | **Recommended MAX OBJ COUNT** derived at runtime: host `os.memory()` (total/free RAM) + a ~200 ms lizard calibration (ms per object for a standard mix, span-cache bytes per object) → min(memory cap, repaint-budget cap) with headroom, cached per machine/app version; New drawing dialog shows it and writes `canvas.max_objects`; status bar `objects n / max` (amber ≥ 80%, ember over); opening on a slower machine warns without editing the file | The estimate is reproducible on this Mac, shown in the dialog and status bar, and rises when speed improves |
| B026 | **Pen settings pop-up + test pad** (R002 18b/19b): floating window, curve editor, sliders, input overrides, scribble pad with live pressure, Save as new pen; **stabilization** (pulled string, live while drawing, also used by the Lasso) | Tuning a pen restyles the pad's strokes live; pad strokes never reach the document; undo covers saved pens |
| B027 | **Select/move + Lasso** (select and lasso fill): selection model shared with the code rows, marquee, move as literal rewrites (one undo step), lasso fill writes a `polygon` row | Moving 2 selected shapes rewrites only their literals; a lasso fill adds one filled row; refused on expression positions with the reason |
| B028 | **Color wheel + harmonies** in the color popover and fill editor: ring + triangle, harmony modes, swatches → variables | Complementary of #3fa34d is its 180° hue; "add all as variables" writes one `colors` entry per swatch, one undo step |
| B019 | First-run tutorial | Runs once, skippable, en + es, `tutorial_done` in settings |

## 🗃️ Later (parked ideas)
| Idea | Why parked |
|---|---|
| Brush creator (Procreate-style) + `.lzbrush` import/export + Procreate/ABR import | See [[Linkz – Brushes]]; after v1.1 |
| Mini videos in pen/tool tooltips | Needs a solid tooltip + animation structure first |
| Own Cocoa window layer (drop Ebitengine) | After v1 |
| Tablet pressure tested on real hardware | No tablet here |
| Pressure on Windows/Linux | macOS first. Windows keeps the stub (mouse, pressure 1.0); real pens need Windows Ink / `WM_POINTER` (`GetPointerPenInfo`), via `golang.org/x/sys/windows`, no cgo |
| Windows/Linux tested on real machines | `scripts/crossbuild.sh` only proves they compile and vet; clipboard, Recycle Bin, reveal and locale need a run on a Windows PC and a Linux desktop |

## 📦 Batches

Briefs live in `docs/Batches/B<nnn> – <name>.md`.

## ✅ Progress log

| Date | What | Result |
|---|---|---|
| 2026-11-03 | Overnight: autocomplete, gradients + blend in the editor, gradient variables, object budget, MCP (+ trace), image tracing (M8), Windows/Linux + Windows Ink, release prep, QA pass + fixes, R002 designs | See [[2026-11-03 Overnight]] |
| 2026-11-03 | **v1 finish line (B009)** | Pencil → editable `stroke([...])`, pressure brush with 5 tunable pens, anchor drag/insert/delete rewriting single literals, shape drag, honest refusals, one undo step per gesture. Tested live through the backdoor; anchor drag ~1 ms/event, 4–14 ms/paint |
| 2026-11-02 | M0 Foundation | `hello.lzdraw` → 6 circles + a line in SVG; host tests green; lizard libraries do all drawing |
