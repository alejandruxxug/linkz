# B009 – Pencil and two-way canvas

Milestone: M3 (the v1 finish line) · Repo: `linkz` · Lane: lizard (`lib/ui/`) · Needs: B007

Read first: `docs/Linkz – Editor Design.md` → "v1 finish line" and "Two-way canvas".

## Goal
Draw a line with the mouse (or trackpad/tablet pen) and edit it afterwards like a vector line in an art program. Every edit is a code edit.

## Build
- **Pencil ✎**: drag on the canvas → points recorded (in canvas coordinates) with live preview; on release `geometry.simplify` (tolerance ~1.5 px, keep endpoints) and insert one row into the selected layer (or as a loose layer at the top if none is selected): `stroke([(40, 300), (120, 210), …], #1d1d1f, 4)`, coordinates rounded to integers. Select it.
- **Brush 🖌** (same gesture, pressure): if the events carry real pressure (`pen.source` is "trackpad" or "tablet"), write `brush([...points], [...widths], color)` with widths = base width × pressure at each kept point; with a plain mouse, derive pressure from speed (slower = thicker, smoothed) so it still feels alive. One undo step per stroke.
- **Anchors**: a selected shape shows its `anchors` (from the B008 shape record). Drag an anchor → rewrite that literal: `arg` (a point or number argument) or `arg` + `item` (one point inside an array literal) via the spans from `src.parse`; re-run while dragging (coalesced moves); one undo step per drag. Radius/corner handles rewrite the number/size argument.
- **Move a whole shape**: drag the shape body → its `move` anchor (circle center, rect pos) or, for stroke/polygon/brush, every point in the array shifted.
- **Refuse honestly**: if the targeted argument isn't a literal (e.g. `horizon + i * 20`), don't move; flash the row and show "position comes from `horizon + i * 20`" in the status bar.
- **Insert/delete anchors**: double-click on a selected stroke → insert a point at `nearest_on_path` between its neighbours; select an anchor (click) + Delete/Backspace → remove that item (min 1 point).
- **Hover cursor feedback** isn't possible yet (no cursor API); highlight the anchor under the mouse instead (request hover moves from the host while a shape is selected).

## Test (tester, backdoor)
`Call AddShape stroke`… no: draw it for real: `Drag` across the canvas with 20 steps → `Dump` shows a new stroke row with ≥2 points; `Screenshot`; `Drag` one anchor → `Dump` shows only that point's text changed; double-click the curve → one more point; Delete an anchor; `Drag … --pressure 0.2,1.0` with the brush → a brush row with rising widths; drag a circle → its center literal changes; drag a shape whose position is an expression → refused, text unchanged; Undo restores each step.

## Done when
The test passes, the screenshots look like a vector app (smooth curve, visible anchors), `go test ./...` green, and drag latency is reported from `Stats`.
