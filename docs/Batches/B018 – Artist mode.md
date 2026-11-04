# B018 – Artist mode (v1.1, wave 3)

Design: R002-05 artist mode, R002-06 artist code drawer, R002-08 blend in artist mode, R002-19 / 19b pen settings (artist, ES), R002-25 color popover (artist, ES); Style Guide design log row "artist mode" and the approval note. Rules: `docs/Batches/_v11-rules.md`.

> "Include a non coder mode that will let the user onto a more Photoshop or Clip Studio Paint layout. Meaning that the code part can be stashed away." — the developer

## Build
- **A mode switch** (`</> Code` / `artist mode` in the header as drawn; remembered in settings as `let mode = "artist"|"coder"`; same file, same selection, same undo history when switching).
- **Artist layout**: vertical tool column on the left (select, lasso, shapes, pencil, brush, text, layer…), canvas in the center, right column = **compact palette** (R002-25 / the B017 sidebar in its compact form) + **layers panel with thumbnails** (one row per layer and per shape inside it; eye/lock; drag to reorder = row moves in the code; blend dropdown + opacity in the layers header as R002-08), **context bar** under the tabs for the selected tool or shape (Aseprite-style: the same properties the coder props panel edits, in a horizontal strip).
- **Code drawer**: the code panel slides in from the bottom/side as R002-06 (toggle with a button and ⌘/), fully functional, so an artist can peek at what their drawing is.
- Layer thumbnails: rendered small from the scene (cached per layer version, updated incrementally — no full re-render per frame).
- Every popup built in waves 1–2 (shape menu, fill editor, pen settings, color popover, blend dropdown, New drawing dialog) works in artist mode and is placed sensibly for the new layout (pen settings shows the artist variant without the code block, as 19/19b).
- The M10 benchmarks must hold in artist mode too: add artist-mode cases to `tests/editor/editor_bench.liz` (drag, keystroke-equivalent property edit, layer reorder at 2k/10k).
