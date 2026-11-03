# B016 – Right-click shape menu + fill editor (v1.1, wave 1)

Design: R002 shape menu and fill editor frames (Style Guide R002 table and metrics; design log rows "shape menu", "fill editor"). Rules: `docs/Batches/_v11-rules.md`.

## Build
- **Shape menu** on right-click / two-finger click on a shape (canvas) or its row: Rename, Duplicate, Bring forward/backward/front/back (moves the row), Hide/Show, Lock, Fill…, Stroke…, Blend ▸ (wired by B020; until then disabled with a tooltip), Delete (danger). Keyboard navigable, Esc closes.
- **Fill editor**, floating beside the shape: kind tabs (none / solid / linear / radial), the gradient bar with draggable stops (add by clicking the bar, remove by dragging off or Delete, hard stops, `steps`), angle / center controls, the selected stop's color through the existing color picker (B017/B028 upgrade it later), live preview on the canvas; writes the fill literal into the shape's row (`gradient(stops, angle)` / `radial_gradient(stops)` / pinned `linear`/`radial`, as `lib/draw` supports) — one undo step per gesture.
- Expression fills (a variable, a computed value) are read-only with "this fill is code: edit it in the row" and an "inline as literal" action.
