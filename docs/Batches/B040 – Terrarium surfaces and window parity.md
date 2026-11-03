# B040 – Terrarium: surfaces, blending, window parity, drag samples, automation (M11)

Repo: `../lizard-script` · Plan: [[Linkz – All Lizard Plan]] · Parallel with B039 (see ownership there)

Generalize what Linkz's Go host does into Terrarium modules any lizard window app can use. Read `../linkz/host/fb.go`, `buf.go`, `layers.go`, `gradient.go`, `blend.go`, `paint.go`, `launcher.go`, `splash.go`, `host/window/*`, `host/pressure/*`, `debug.go`, `debugcli.go`, `backdoorfile.go`, and `../linkz/docs/Linkz – Tools & Setup.md` (backdoor reference). Don't change ../linkz.

## Build
1. **Surfaces**: named pixel buffers with a scale (a UI buffer shown at 2× and a full-resolution canvas composited into a rect of it, as Linkz does), `offscreen`/`keep`/`blit`/`clip`/`unclip`/`target`, dirty-rect tracking and compose of each dirty rect separately (B037's fix), PNG encode to a file and to base64. Same drawing functions on images, windows and surfaces.
2. **Per-pixel primitives** (the logged exception, now general): `span_gradient`/`spans_gradient` (linear/radial, steps, hard stops — port `gradient.go`), blend layers `push_layer(mode, opacity)`/`pop_layer` with normal, multiply, screen, overlay, darken, lighten, add (port `blend.go`), composite.
3. **Window parity**: frameless + transparent windows, switching chrome at run time (decorated, size, center: the splash → editor switch), minimum size, focus/leave/enter events, wants-hover, fps, the "only repaint when asked" mode (`needs_paint`).
4. **Pen**: native per-drag samples (every sample between frames, with pressure and time) in events, as Linkz's `dragSamples`; the curves already exist.
5. **Automation (the agent backdoor, general)**: a window app opts in (`window.automation(true)` or an env var); it writes the per-pid token file (0600) like Linkz does; `lizardscript drive [--pid N] <Command>` sends List, Dump (the app provides it via a lizard function), Screenshot, Click, DoubleClick, Drag (`--pressure`, `--full`), Scroll, Key, Type, Resize, Focus, Leave, Enter, Stats, Wait, Quit, `Call Name args…` (calls a lizard function the app registers). Port `debug.go`/`debugcli.go`. Note: `lizardscript drive` needs a line in `main.go`, which B039 owns — write the exact patch in your report.
6. **Headless mode** for tests: `window.open(..., {"headless": true})` runs ticks without a real window (fixed size, injected events, screenshots), so app tests and benchmarks run in `go test`/`lizardscript test` and CI.

## Done when
Every primitive tested (golden pixels against the Linkz implementations: gradients, blend modes, compose), examples (`examples/terrarium/layers.liz` showing blend modes; `splash.liz` showing a frameless transparent window that becomes a normal one), automation tested headless, docs in TERRARIUM.md + guide chapter 17 (en + es), cross-builds. Clean code. Don't commit.
