# 🦎 Linkz: the move to all-lizard (M11)

> "Okay finish the transition for the app to be built fully on lizard." — the developer, 2026-11-03

## Goal

The Linkz repository contains **no Go**. Linkz is a lizard-script program (`main.liz` + `lib/`), run by `lizardscript` and shipped as one executable built by `lizardscript build`. Every Go line Linkz needs today moves **into the language**: the Terrarium (general modules any lizard program can use) or the lizard-script tools. The per-pixel exception (gradients, blending, compositing) becomes general Terrarium primitives.

Today: ~7,200 lines of Go in `main.go`, `debugcli.go`, `host/` (without tests) and ~15,000+ lines of lizard in `lib/`.

## Where each Go piece goes

| Linkz Go today | Becomes | Batch |
|---|---|---|
| `lib/lib.go` (embeds and loads `lib/*.liz` in order) | a module system in the language: `use "draw/shapes"` | B039 |
| `src.parse`, `src.run`, `src.library_docs`, `scan.go`, `attribution.go`, `src.here` | a `lizard` Terrarium module: tokens and syntax trees with start/end positions, sandboxed runs with a step limit and source names, doc comments; the scanner rewritten in lizard on top of it | B039 (language), B041 (scanner in lizard) |
| `stdio.go` (`io`: stdin lines, stdout/stderr) | Terrarium `io` | B039 |
| `osprims.go`, `thumbs.go`, `sys_*.go`, `mem*.go`, `trashinfo.go` (Trash, reveal in Finder/Explorer, mtime/size, memory) | Terrarium `fs` / `os` additions | B039 |
| `fb.go`, `buf.go`, `layers.go`, `gradient.go`, `blend.go` (UI buffer at 2×, canvas buffer, offscreen/keep/blit/clip, gradient spans, blend layers, compose) | Terrarium `img`/`window` surfaces: named buffers with a scale, clip, blit, gradient spans, blend layers, PNG encode/base64 | B040 |
| `host/window` (Ebitengine: frameless/transparent splash, switching to a decorated window, min size, focus/leave, wants-hover) | Terrarium `window` options and events | B040 |
| `host/pressure` (per-drag native samples, curves) | Terrarium `pen` (curves exist; add drag samples) | B040 |
| `debug.go`, `debugcli.go`, `backdoorfile.go` (the agent backdoor) | Terrarium `window` automation: any lizard window app can opt in to be driven by `lizardscript drive` (list, click, drag with pressure, type, screenshot, call a function, stats) — same per-pid token file | B040 |
| `img.go` (decode, shrink) | Terrarium `img` (exists) | — |
| `app.go`, `launcher.go`, `splash.go`, `settings.go`, `machine.go`, `paint.go`, `host.go`, `engine.go`, `conv.go` | lizard: `main.liz` and `lib/app/*.liz` | B041 |
| `main.go` (CLI: edit, render, paint, trace, debug, calibrate, mcp) | `main.liz` with `os.args()` | B041 |
| Go tests (`host/*_test.go`, perf, scene, incremental, ui/gfx/draw test runners) | lizard tests run by `lizardscript test` (programs + `.expected`, benchmarks with `time`), headless windows | B039 (runner), B041 (port) |
| GoReleaser, install scripts | `lizardscript build . -o linkz` per OS, same release workflow | B042 |

## Batches

1. **B039 – lizard-script: modules, reflection, io, OS, test runner, build** (language + tools).
2. **B040 – Terrarium: surfaces, blending, window parity, drag samples, automation** (runs in parallel with B039; file ownership split in the briefs).
3. **B041 – Linkz port**: `main.liz`, the app loop, splash, settings, MCP, trace, calibrate, the scanner in lizard, tests in lizard. Both hosts live side by side until parity.
4. **B042 – Delete the Go host**: parity checks (every old test's case covered, smoke, perf within 10% of the Go host), `lizardscript build`, release scripts, docs. Then the repo is `.liz` + docs only.

## Rules

- General first: a Terrarium feature must make sense for any lizard program, not just Linkz. No Linkz names in lizard-script.
- Performance can't regress: the M10 benchmarks (`docs/Batches/fixtures/`) must stay within 10% of the Go host's numbers, measured both ways.
- Clean code; LEP 0001 style; English + Spanish strings unchanged.
