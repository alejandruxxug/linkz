# B066 – The brush creator window

Model: Opus. Rules: `_v11-rules.md`. The approved design is the spec: R005 frames 1, 3, 4 (`docs/Art/Reviews/R005/R005-01-creator-coder-shape-layoutA.png`, `-03-…grain`, `-04-…es-pressure`) and the R005 section of `docs/Linkz – Style Guide.md` (metrics, icons, the approval answers 1, 4, 5). Read only those, plus `lib/ui/pen_window.liz`, `pen_window_input.liz`, `pen_props.liz`, `pen_pad.liz`, `pen_edit.liz` and the `pens.liz` header (the keys).

- The creator **replaces** the R002 pen settings pop-up: ⚙ opens it on the Pressure tab (the R002 curve editor, reused as is), stabilization moves into the Stroke tab. Layout A: tabs in a column on the left (Stroke, Taper, Shape, Grain, Dynamics, Pressure, Color, Rendering), a wide live test pad below, name field, **Save brush** and **Export…**.
- Every control writes the **flat keys** `lib/draw/pens.liz` documents (`render`, `spacing`, `shape`, `grain` hash, …); the tabs only group them. In coder mode the window shows the hash it writes (as the frame does).
- The test pad redraws live while tuning (stamped brushes too: B060 is in; pad strokes never reach the document).
- Save brush calls `brush_library_save(pen)` and Export… writes a `.lzbrush` via `brush_write` (B065 builds `lib/app/brush_files.liz` in parallel: call those names; if it isn't there yet when you test, stub them in your test only and say so).
- Shape and grain pickers list what `lib/draw` supports (B067 adds spray, leaf, bitmap, paper, canvas, chalk, hatch in parallel: read the list from one function in `lib/draw/stamps.liz`, e.g. `stamp_shapes()` / `stamp_grains()`, so new ones appear on their own; "import image" for bitmap shapes is in the frame but may show disabled until B067 lands — say which).
- en + es (frame 4 is the Spanish check); one undo step per Save; Esc closes without saving; the old pop-up's tests keep passing or move to the creator with the same coverage.
- Rig tests per tab; drive the real app (isolated HOME) and put shots next to frames 1/3/4 in `docs/Art/Checks/v1.2/creator/`, listing every difference from the frames.
- Others are in `lib/app/brush_files.liz` (B065) and `lib/draw/stamps*.liz` (B067): stay out; list shared edits.
