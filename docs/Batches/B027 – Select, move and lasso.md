# B027 – Select/move tool and Lasso (select + fill) (v1.1, wave 1)

Design: R003-1 select/move, R003-2 lasso select, R003-3 lasso fill; Style Guide R003 metrics (tool order, Select, Lasso). Rules: `docs/Batches/_v11-rules.md`.

## Build
- **Selection model** shared by canvas and code rows: click selects (its row highlights), Shift-click adds/removes, Esc clears, rows ↔ canvas in sync.
- **Select tool** (V, first in the toolbar): marquee selects shapes inside (Alt: touched), drag moves the selection (ghost at the old place, `move −30, +40` tooltip), arrow keys nudge 1 px (Shift 10), bounding box with 8 scale handles (Shift keeps the ratio) and a rotation handle. Moves/scales/rotations rewrite the literals of every selected shape in one undo step; shapes whose position isn't a literal are refused with the reason.
- **Lasso tool** (L) with an options strip: **Select** (freehand loop; shapes inside become the selection, then move like Select) and **Fill** (the loop becomes a new `polygon([...], color)` row in the current color and layer, simplified points). Marching-ants loop while drawing. Uses B026's `stabilize` (its own file) when present; otherwise a tiny stand-in behind the same interface, reported.
