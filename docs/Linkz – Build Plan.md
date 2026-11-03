# 🗓️ Linkz: Build Plan

## ⏳ Waiting on you

| What | Why it's blocking | Where to look | Since |
|---|---|---|---|
| OK to commit the extensions-log section + CLAUDE.md in lizard-script | Language work starts there | lizard-script `git status` | 2026-11-02 |

## 🧱 Milestones

| # | Milestone | Done when |
|---|---|---|
| M0 ✅ | Foundation | `linkz render examples/hello.lzdraw` writes an SVG made entirely by lizard libraries; `go test ./...` green |
| M1 | Language basics | lizard-script has string escapes, dot access and floats, each with tests and an extensions-log row; `svg.liz` uses `"` and `\n` |
| M2 | Modules | `use "shapes"` works; `lib.Order` is gone; libraries call each other by name |
| M3 | Editor | `linkz serve` opens a page: code on the left, live SVG on the right, errors on their line; open/save `.lzdraw` |
| M4 | Math + more shapes | `math` (sin, cos, sqrt, random), `path`, `polygon`, `group`, `color`, `transform` |
| M5 | Raster export | A rasterizer in lizard; PNG export from the editor and `linkz render --png`; first speed measurement logged |
| M6 | Speed | Bytecode compiler + VM in lizard-script; rasterizer re-measured |
| M7 | Agents | `linkz mcp`: run / render (PNG) / docs / save, protocol written in lizard |
| M8 | Tracing | `trace.image(...)` turns a PNG into editable shapes |
| M9 | Release | GoReleaser binaries for macOS/Windows/Linux + install script; the reveal |

## 🧱 v1.1 (after the v1 finish line)
| # | Batch | Done when |
|---|---|---|
| B016 | Right-click menu + Fill editor + gradients (`linear`/`radial` values in lib/draw, gradient spans in the host) | Approved R002 frames; a gradient made in the editor renders on canvas, in SVG and in PNG |
| B017 | Palette sidebar (HSV picker, recent, document palette as `let colors = {...}`) | Changing a named color updates every shape using it |
| B018 | Artist mode (layers with thumbnails, code drawer) | Switching modes keeps the same file and selection |
| B019 | First-run tutorial | Runs once, skippable, en + es, `tutorial_done` in settings |

## 🗃️ Later (parked ideas)
| Idea | Why parked |
|---|---|
| Mini videos in pen/tool tooltips | Needs a solid tooltip + animation structure first |
| Own Cocoa window layer (drop Ebitengine) | After v1 |
| Tablet pressure tested on real hardware | No tablet here |
| Pressure on Windows/Linux | macOS first |

## 📦 Batches

Briefs live in `docs/Batches/B<nnn> – <name>.md`.

## ✅ Progress log

| Date | What | Result |
|---|---|---|
| 2026-11-02 | M0 Foundation | `hello.lzdraw` → 6 circles + a line in SVG; host tests green; lizard libraries do all drawing |
