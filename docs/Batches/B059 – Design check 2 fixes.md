# B059 – Fixes from design check 2 (small tasks)

Source: `docs/Reviews/2026-11-04 design check 2.md` (each item: what, where, expected, actual; shots in `docs/Art/Checks/v1.2/`). Rules: `_v11-rules.md`. One agent per task; each reads only its review items and the files named. Mark each fixed item ✓ (one line + a `-fixed` shot) in the review.

| Task | Items | Files |
|---|---|---|
| **F1** card placement | 🔴 the artist hover card covers the hovered layer's canvas box (R004-3/4); the card preview draws layers in Normal (proposal in the Style Guide R004 metrics) | `lib/ui/layer_card.liz`, `layer_iso.liz` |
| **F2** brush feedback | live brush ink hairlines at segment joins (R003-4); stabilizer feedback: "ink follows the string" tooltip, `pressure 0.71` status readout, dotted cursor path, cursor square (R003-4) | `lib/ui/view_tools.liz`, `stabilize_view.liz`, `tools_stroke_stab.liz` |
| **F3** rows | hidden rows: amber eye-off, no raw `-~> [hidden]` comment shown (R004-7); clicking a row's eye doesn't scroll the code to line 1; a "nothing selected" state on open so the "hover a layer to preview it" hint shows (R004-1), en + es | `lib/ui/editor_rows.liz`, `view_panels.liz`, `view_input.liz` (eye click), `artist_layers.liz` |
| **F4** palette | the palette stays on `colors.sky` after a variable edit: it must let go on selecting another row/layer and on opening another file (R002-8); the artist compact palette shouldn't open collapsed because the coder sidebar was | `lib/ui/palette_view.liz`, `palette_compact.liz`, `color_target.liz`, `color_state.liz` |
| **F5** properties | gradient fills in properties: a gradient chip, a `steps` field and the hint instead of cut-off text (R002-2/3/9); `let colors` row properties: one swatch field per entry (R002-4, R003-5) | `lib/ui/editor_props.liz`, `view_layout.liz` (props build), `fill_model.liz`, `color_vars.liz` |
| **F6** trace preview | the original half drawn on the same paper as the traced half and at screen resolution (R003-6) | `lib/ui/trace_view.liz`, `trace_model.liz` |
| **F7** small items | stale status-bar cursor position; completion popup over the PROPERTIES header; opacity slider label; `layer` lowercase in the layer header; Spanish wheel caption wraps to 3 lines | `view_panels.liz` (status), `complete_view.liz`, `blend_props.liz`, lang files, `color_popover.liz` |

F1, F2, F4 and F6 can run at the same time (separate files). F3 → F5 → F7 one after another (they share panel files).
