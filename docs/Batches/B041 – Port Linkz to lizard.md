# B041 – Port Linkz to lizard (M11)

Repo: `linkz` · Plan: [[Linkz – All Lizard Plan]] · Needs: lizard-script 5396133 (B039 + B040, pinned in `../.lizard-stable`)

## Goal
Linkz runs as a lizard program: `lizardscript run .` (or `lizardscript build . -o linkz`) does everything `go build -o linkz . && ./linkz …` does today — `edit` (splash → editor), `render`, `paint`, `trace`, `calibrate`, `mcp`, the agent backdoor (now `lizardscript drive`) — with the Go host left in place but unused, so both can be compared. B042 deletes the Go.

## Read first
`docs/Linkz – All Lizard Plan.md` (where each Go piece goes), lizard-script `docs/TERRARIUM.md` (window, surfaces, pen, img, io, lizard, fs, os, automation, built programs) and `proposals/0002-modules.md` (`use`), and these mapping notes from the language builders:

- **Loading**: `lib/lib.go`'s order becomes explicit `use` lines at the top of each file. Names are live (shared mutable globals keep working across files); names are **not** re-exported, so each file uses every file it needs; `_names` are private.
- **`src.run`** → `lizard.run(code, {"engine": "vm", "steps": …, "globals": {…libraries…}, "name": …, "on_statement": f, "on_layer": g})`; read `drawing`/`canvas` from `r.globals`. Line tagging via `on_statement(line, column, source)` filtered by positions from `lizard.parse`. `library_docs` → `lizard.docs`. The scanner (`host/scan.go`) is rebuilt in lizard on `lizard.parse` spans + `lizard.tokens` (comments are tokens).
- **fb → window + surfaces**: UI buffer `window.open(w/2, h/2, title, {"scale": 2, …})`, canvas `window.surface(w, "canvas", rect, 1)`; every drawing call takes the target first (`fb.spans(list, offset, color)` → `spans(t, list, color, offset)`, `fill(t, [x, y, w, h], c)`, `offscreen` → `img.new`, `keep` → a variable, `composite(t, img, mode, op)`, `fb.dirty` → `window.dirty(t, rect)`, `save_png` → `img.save(t, path, rect)`, `png_base64` → `img.png_base64`). Consider a tiny lizard `fb` shim only if it keeps the port readable; measure its cost.
- **Events**: `canvas_pos` → `surface_pos` + `e.surface == "canvas"`; `wheel` → `scroll`, `text` → `type`, buttons "left"/"right"/"middle"; `move.samples` carry pressure/raw/t; `wants_hover` → `window.set(w, {"hover": …})`; `needs_paint` → `on_demand` + `window.needs_paint(w)`.
- **Splash → editor**: open with `transparent: true`, frameless; then `window.set(w, {"frameless": false, "floating": false, "scale": 2, "size": …, "min_size": (320, 210), "resizable": true, "center": true})` (see lizard-script `examples/terrarium/splash.liz`).
- **Backdoor**: `window.automation(w, {"dump": f, "stats": g, "commands": {...}})`; drive file `<config dir>/lizard/drive-<pid>.json`; testers use `LIZARD_AUTOMATION=1` and an isolated HOME; Stats reports `tick_ms`/`compose_ms` (keep Linkz's per-phase numbers in the dump/stats hooks).
- **MCP**: `io.read_line`/`io.write`/`io.error`; log with `io.error` only; wrap each request in `lizard.run` with `steps`.
- **OS**: Trash `fs.trash`, Finder `os.reveal`, memory `os.memory`, file info `fs.list(dir, {"info": true})`, settings/machine files via `fs` + `lizard.run`, locale `os.locale`.
- **Data files** (font, sprites, splash art, fixtures): `bundle.txt` + `fs.join(os.program_dir(), …)`.
- **Known gap**: `img.load` doesn't read `bundle:` paths yet (only `fs` and `use` do) — fix it in lizard-script's `terrarium/img.go` with `bundle.ReadFile` (small, allowed in this batch) or avoid needing it.
- Pixels may differ by 1 where alpha is partial: regenerate soft-edge goldens once, explained in the commit.

## Build
1. `main.liz` (CLI on `os.args()`) and `lib/app/*.liz` (app loop, launcher/splash, settings, machine calibration cache, paint/render/trace commands, MCP stdio, home OS ops), every lib file with its `use` lines.
2. Tests in lizard: every Go test's cases (host, ui, gfx, draw, mcp, trace, scene/golden, incremental, perf benchmarks as `*_bench.liz`) as `lizardscript test` programs, using headless windows for the editor. Keep the Go tests running too until B042.
3. `scripts/smoke.sh` gets a lizard variant (`lizardscript run . render …`, an edit session driven headless); a parity script compares outputs of the Go host and the lizard port (rendered PNGs, MCP transcripts, trace output) on the fixtures.
4. Perf: the M10 benchmarks on the lizard port vs the Go host; within 10% is the bar (report and explain any gap).

## Done when
`lizardscript run . edit` opens the real editor (developer checks by hand; report how), all commands work, `lizardscript test` green on both engines, parity script clean, perf table. Clean code; LEP 0001 style; strings unchanged (en + es). Don't commit. Don't touch ../.lizard-stable or go.work; for the img.load fix, edit ../lizard-script and say so (the orchestrator commits it and moves the pin).
