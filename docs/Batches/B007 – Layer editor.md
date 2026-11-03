# B007 – Layer editor

Milestone: M3 Editor · Repo: `linkz` · Lane: lizard (`lib/ui/`) · Needs: B005, B006, B008, lizard-script B004 (layers, text/list helpers)

**Build from the approved Figma frames** (file "Linkz", page ✅ Approved / 🧪 Reviews → R001 – Editor views; read them with the Figwright MCP `get_design_context` for exact sizes and colors, and re-read before building since the developer may have edited them). The palette is **Linkz Dark** (Style Guide), not Lizard 16. Each window also writes its own backdoor file (`~/.linkz/debug-<pid>.json`), and `linkz debug` targets the newest one unless `--pid` is given, so test runs never quit the developer's window (host change; the orchestrator assigns it).

Read first: `docs/Linkz – Editor Design.md` (the whole thing: rows, rules, creating, properties, the host contract) and `docs/Linkz – Style Guide.md`.

## Goal
`linkz edit drawing.lzdraw` shows the real editor: the code is the layers panel and the properties panel, the canvas shows the drawing, and the file saves.

## Build (all lizard, `lib/ui/*.liz`; replace B005's placeholder `app.liz`)
- **Layout**: left panel (rows + properties at the bottom), right canvas area; resizable split; toolbar row (○ □ ╱ ⬠ ～ ✎ T ƒ ▤ layer, export PNG, save). **8-bit style**: follow `docs/Linkz – Style Guide.md` exactly (Lizard 16 palette, bevels, 8×8 icons, no anti-aliasing). The UI draws into the half-resolution UI buffer (font at scale 1 there, ×2 on screen); the drawing goes into the full-resolution canvas target (`fb.canvas`, `fb.target`) and uses `canvas_pos` from events. Make the 8×8 icons as lizard data in `lib/ui/icons.liz`; have the `design` agent review them.
- **Rows** from `src.parse`: one row per line, indentation shown as depth, kind icon, name (trailing comment or callee), code text with simple syntax colors (keywords, numbers, strings, colors as swatches, comments). Layers collapsible (▾/▸, editor state only). Error rows marked red with the message on hover/selection.
- **Text editing** per row: caret, click to place, arrows (Up/Down move between rows keeping the column), Home/End, Backspace/Delete (joining rows at the edges), typing, Cmd+C/X/V of the selected row(s) (in-app clipboard if the host has none), Tab/Shift+Tab indent/outdent by 2 spaces.
- **Creating**: Enter → new row below at the same depth; Shift+Enter → above; an empty row shows the inline tool strip; picking a tool writes the template (centered on the canvas, Style Guide colors), runs, selects it, focuses the first property. `layer` tool writes `layer layer1` with one empty indented row.
- **Properties panel**: for the selected call row, one field per argument, labelled by the lizard function's parameter names (ask `src.parse`/the drawing env for the function's parameters); number fields (type, or drag horizontally to scrub), point fields (two numbers), color fields (swatch + hex text), expression fields (code text). Editing rewrites exactly that argument's span in the row (spans from `src.parse`) and re-runs. Layer rows: name, hidden, locked (writing the trailing comment).
- **Canvas**: runs `src.run` after every edit (debounced to the next paint); paints shapes in paint order (bottom row first) with `raster.shape`, centered, canvas background, a subtle border; errors don't blank the canvas (keep the last good drawing, dimmed). Click a shape → topmost hit (point-in-path, or near a stroke) selects its row and expands its layer; the selected shape gets an outline and its anchor handles (drawn; dragging is B009).
- **Tooltips**: hovering (host hover moves while the pointer is over the panel) for ~400 ms over a callee name in a row, a property field, a tool in the toolbar or a pen in the picker shows a small tooltip box drawn in lizard: the doc summary, the parameter's line (for a field), and the example in code colors. Docs come from `src.library_docs()` and, for the user's own functions, from `src.parse` (their own `-~>` comments). Never covers the hovered thing; disappears on move/click. No doc → no tooltip.
- **Languages**: follow Editor Design → "Languages". `lib/ui/i18n.liz` with `t(key, [values])`, `lib/ui/lang/en.liz` and `lib/ui/lang/es.liz` (natural Spanish, not word-for-word), the toolbar language switch, `settings.liz` load/save via `os.config_dir()`, `os.locale()` as the default. Spanish tooltips for the drawing functions in `lib/draw/docs/es.liz`. A Go test fails if any `t("…")` key used in `lib/ui` is missing from `en.liz`, and lists keys missing from `es.liz`.
- **Files**: Cmd+S saves (atomic via host), title shows `•` when dirty; Export PNG renders the drawing at canvas size into an offscreen buffer and saves `name.png` next to the file; Cmd+Z / Shift+Cmd+Z undo/redo (snapshot stack of the text + caret).
- **Backdoor**: `app.debug_state()` returns rows (line, depth, kind, name, rect on screen, flags, error), selection, caret, properties (name, value text, rect), canvas rect, dirty, file; `app.debug_commands`: `SelectRow n`, `SetProp name value`, `AddShape kind`, `ToggleLayer n`, `Caret line col`.
- **Speed**: repaint only dirty regions (row area / properties / canvas separately); cache rendered row spans; keep a paint under ~50 ms for a 100-line file and report the real numbers from `Stats`.

## Test (tester, through the backdoor; no Chrome)
Create a drawing from empty only with `Call AddShape`, `Type`, `Key`, `SetProp`; check `Dump` and `Screenshot` at each step; save; reopen with `Open`; export PNG and Read it; errors (type a bad line) show on the row and keep the last canvas.

## Done when
The test passes end to end, `go test ./...` green, screenshots look like the Editor Design mockup, paint timings reported.
